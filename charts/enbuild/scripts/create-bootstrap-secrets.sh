#!/usr/bin/env bash
# =============================================================================
# ENBUILD hub — bootstrap-secret generator (clean-room "deploy from scratch")
# =============================================================================
# The hub chart references its credentials via *.existingSecret and creates NONE
# of them (secrets are deliberately kept out of Helm values + git). This script
# generates fresh random credentials and creates every Secret the quickstart
# values file expects, so a brand-new engineer goes from "empty namespace" to
# "ready to helm install" in one command — no secret hand-crafting.
#
# It is IDEMPOTENT and SAFE TO RE-RUN: an already-existing Secret is LEFT AS-IS
# (rotating Mongo/RabbitMQ passwords would break the running StatefulSets whose
# on-disk data was initialised with the old password). Use FORCE=1 to recreate.
#
# What it creates (names assume RELEASE=enbuild-ib, NAMESPACE=enbuild):
#   enbuild-ib-mongo-secrets     MONGO_INITDB_ROOT_USERNAME/PASSWORD/DATABASE/SERVER   (mongodb.existingSecret)
#   enbuild-ib-encryption-key    ENCRYPTION_KEY                                         (enbuildBk.encryptionKey.existingSecret)
#   enbuild-ib-messaging         RABBIT_MQ_CONNECTION_STRING                            (enbuildBk.messaging.existingSecret)
#   enbuild-ib-rabbitmq-creds    rabbitmq-password / rabbitmq-erlang-cookie             (rabbitmq.auth.existing*Secret)
#   enbuild-ib-image-pull-secret .dockerconfigjson (registry1 + gitlab)                 (global.imagePullSecretName)
#   enbuild-ib-install-agent     GITLAB_TOKEN/ENBUILD_REPO1_USER/ENBUILD_REPO1_TOKEN    (enbuildBk.installAgent.existingSecret) [optional]
#   enbuild-ib-reaper-svc        REAPER_SVC_AWS_ACCESS_KEY_ID/_SECRET_ACCESS_KEY/_ACCOUNT_ID (enbuildBk.reaperSvc.existingSecret) [optional, cluster destroy]
#   enbuild-ib-observability     SIEM_AUTH_HEADER/PROMETHEUS_TOKEN/LOKI_TOKEN (subset)  (enbuildBk.observability.existingSecret) [optional, SOO §1.5]
#   enbuild-ib-export-signing    SIEM_SIGNING_KEY (ECDSA P-256 PEM, auto-generated)     (enbuildBk.exportSigning.existingSecret) [recommended, CCM-32 Auditable]
#
# The RabbitMQ broker password and the backend's RABBIT_MQ_CONNECTION_STRING are
# generated ONCE and written to BOTH secrets, so they can never diverge.
#
# Required tools: kubectl, openssl (or /dev/urandom), base64.
# Required context: kubectl must already target the destination cluster.
# =============================================================================
set -euo pipefail

NAMESPACE="${NAMESPACE:-enbuild}"
RELEASE="${RELEASE:-enbuild-ib}"
FORCE="${FORCE:-0}"

# --- registry credentials (image pulls). REQUIRED for private/Iron Bank images.
#     A clean-room engineer supplies the entitlements they were granted. ---
REPO1_REGISTRY="${REPO1_REGISTRY:-registry1.dso.mil}"      # Iron Bank — mongo/rabbitmq/keycloak/headlamp
REPO1_USER="${REPO1_USER:-}"
REPO1_TOKEN="${REPO1_TOKEN:-}"
GITLAB_REGISTRY="${GITLAB_REGISTRY:-registry.gitlab.com}"  # ENBUILD app images (BE/FE/UI/MQ/user)
GITLAB_USER="${GITLAB_USER:-}"
GITLAB_TOKEN="${GITLAB_TOKEN:-}"

# --- install-agent credentials (OPTIONAL). Only needed for catalog LAUNCHES /
#     agent installs; a login-functional hub does not require them. ---
IA_GITLAB_TOKEN="${IA_GITLAB_TOKEN:-}"
IA_REPO1_USER="${IA_REPO1_USER:-${REPO1_USER}}"
IA_REPO1_TOKEN="${IA_REPO1_TOKEN:-${REPO1_TOKEN}}"

# --- teardown Reaper delete credentials (OPTIONAL). Only needed for cluster
#     DESTROY (the cloud-release gate + tag-then-delete ELB sweep). WITHOUT them,
#     teardown orphans the Istio gateway ELB. The BE enforces an account-scope
#     contract (teardown-creds.service.ts): if the delete creds are set, the
#     12-digit REAPER_SVC_AWS_ACCOUNT_ID is REQUIRED (it refuses an unscoped cloud
#     read + validates the id == the account the creds resolve to), so all three
#     are written together. The complete set may live here in a DEDICATED reaper-svc
#     Secret (enbuildBk.reaperSvc.existingSecret) OR in the install-agent Secret —
#     do NOT split it across both (chart NOTES validates each source as a whole).
#     Optional REAPER_RESOURCE_TAG_KEY/_VALUE narrow the ELB tag gate (must match
#     the IAM condition); omit both to use the enbuild-reaper-approved=true default. ---
REAPER_ACCESS_KEY_ID="${REAPER_ACCESS_KEY_ID:-}"
REAPER_SECRET_ACCESS_KEY="${REAPER_SECRET_ACCESS_KEY:-}"
REAPER_ACCOUNT_ID="${REAPER_ACCOUNT_ID:-}"            # 12-digit AWS account the creds resolve to (REQUIRED with the above)
REAPER_RESOURCE_TAG_KEY="${REAPER_RESOURCE_TAG_KEY:-}"
REAPER_RESOURCE_TAG_VALUE="${REAPER_RESOURCE_TAG_VALUE:-}"

# --- observability/SIEM bearer tokens (OPTIONAL, SOO §1.5). Only the SECRET
#     tokens go here; the non-secret endpoints are set in values
#     (enbuildBk.observability.*). Any subset may be supplied; the secret is
#     created only when at least one is set. ---
OBS_SIEM_AUTH_HEADER="${OBS_SIEM_AUTH_HEADER:-}"   # e.g. 'Bearer <p1-siem-token>'
OBS_PROMETHEUS_TOKEN="${OBS_PROMETHEUS_TOKEN:-}"   # hub-self Prometheus bearer
OBS_LOKI_TOKEN="${OBS_LOKI_TOKEN:-}"               # Loki bearer (pod-log tailing)

log()  { printf '  %s\n' "$*"; }
section() { printf '\n== %s ==\n' "$*"; }

rand() { # rand <chars> -> URI-safe alphanumeric (no amqp:// / URI-breaking chars)
  local n="${1:-32}"
  if command -v openssl >/dev/null 2>&1; then
    # hex is alphanumeric + URI-safe; finite output → safe under `set -o pipefail`.
    openssl rand -hex "$(( (n + 1) / 2 ))" | cut -c "1-${n}"
  else
    # Subshell disables pipefail so head closing the urandom stream isn't fatal.
    ( set +o pipefail; LC_ALL=C tr -dc 'A-Za-z0-9' </dev/urandom | head -c "$n" )
  fi
}

secret_exists() { kubectl -n "$NAMESPACE" get secret "$1" >/dev/null 2>&1; }

apply_secret() { # apply_secret <name> <kubectl-create-secret-args...>
  local name="$1"; shift
  if secret_exists "$name" && [ "$FORCE" != "1" ]; then
    log "exists, kept: $name  (FORCE=1 to recreate)"
    return 0
  fi
  kubectl -n "$NAMESPACE" create secret "$@" --dry-run=client -o yaml | kubectl apply -f - >/dev/null
  log "created: $name"
}

section "ENBUILD hub bootstrap secrets → ns/$NAMESPACE (release $RELEASE)"
kubectl get namespace "$NAMESPACE" >/dev/null 2>&1 || { log "creating namespace $NAMESPACE"; kubectl create namespace "$NAMESPACE" >/dev/null; }

# 1) MongoDB root credentials --------------------------------------------------
# Bundled Mongo (mongodb.enabled=true) hardcodes envFrom <release>-mongo-secrets,
# so the name MUST be exactly this. The backend assembles MONGODB_ENDPOINT from
# these fields. MONGO_SERVER points at the bundled single-node service.
section "1/9 MongoDB credentials"
if secret_exists "${RELEASE}-mongo-secrets" && [ "$FORCE" != "1" ]; then
  log "exists, kept: ${RELEASE}-mongo-secrets"
else
  MONGO_PW="$(rand 32)"
  apply_secret "${RELEASE}-mongo-secrets" generic "${RELEASE}-mongo-secrets" \
    --from-literal=MONGO_INITDB_ROOT_USERNAME=enbuild \
    --from-literal=MONGO_INITDB_ROOT_PASSWORD="$MONGO_PW" \
    --from-literal=MONGO_INITDB_DATABASE=enbuild \
    --from-literal=MONGO_SERVER="${RELEASE}-mongo.${NAMESPACE}.svc.cluster.local:27017"
fi

# 2) At-rest encryption key ----------------------------------------------------
section "2/9 Backend at-rest ENCRYPTION_KEY"
apply_secret "${RELEASE}-encryption-key" generic "${RELEASE}-encryption-key" \
  --from-literal=ENCRYPTION_KEY="$(rand 48)"

# 3+4) RabbitMQ broker password + erlang cookie + backend connection string ----
# ONE password drives both the broker (rabbitmq-password) and the backend's
# RABBIT_MQ_CONNECTION_STRING — they can never diverge.
section "3/9 RabbitMQ credentials (broker) + 4/9 backend connection string"
if secret_exists "${RELEASE}-rabbitmq-creds" && [ "$FORCE" != "1" ]; then
  log "exists, kept: ${RELEASE}-rabbitmq-creds (reusing its password for messaging)"
  RMQ_PW="$(kubectl -n "$NAMESPACE" get secret "${RELEASE}-rabbitmq-creds" -o jsonpath='{.data.rabbitmq-password}' | base64 -d)"
else
  RMQ_PW="$(rand 32)"
  apply_secret "${RELEASE}-rabbitmq-creds" generic "${RELEASE}-rabbitmq-creds" \
    --from-literal=rabbitmq-password="$RMQ_PW" \
    --from-literal=rabbitmq-erlang-cookie="$(rand 32)"
fi
# messaging secret always reflects the current broker password
kubectl -n "$NAMESPACE" create secret generic "${RELEASE}-messaging" \
  --from-literal=RABBIT_MQ_CONNECTION_STRING="amqp://enbuild:${RMQ_PW}@${RELEASE}-rabbitmq:5672/" \
  --dry-run=client -o yaml | kubectl apply -f - >/dev/null
log "created/updated: ${RELEASE}-messaging"

# 5) Image pull secret (combined registries) -----------------------------------
section "5/9 Image-pull secret (${GITLAB_REGISTRY} + ${REPO1_REGISTRY})"
if secret_exists "${RELEASE}-image-pull-secret" && [ "$FORCE" != "1" ]; then
  log "exists, kept: ${RELEASE}-image-pull-secret"
elif [ -n "$REPO1_USER$REPO1_TOKEN$GITLAB_USER$GITLAB_TOKEN" ]; then
  auths=""
  add_auth() { # add_auth <registry> <user> <token>
    [ -z "$2$3" ] && return 0
    local b64; b64="$(printf '%s:%s' "$2" "$3" | base64 | tr -d '\n')"
    auths="${auths}${auths:+,}\"$1\":{\"username\":\"$2\",\"password\":\"$3\",\"auth\":\"$b64\"}"
  }
  add_auth "$GITLAB_REGISTRY" "$GITLAB_USER" "$GITLAB_TOKEN"
  add_auth "$REPO1_REGISTRY"  "$REPO1_USER"  "$REPO1_TOKEN"
  # Write the dockerconfigjson to a 0600 temp file (never world-readable in /tmp).
  dockercfg="$(mktemp "${TMPDIR:-/tmp}/.enbuild-dockercfg.XXXXXX")"
  chmod 600 "$dockercfg"
  printf '{"auths":{%s}}' "$auths" > "$dockercfg"
  kubectl -n "$NAMESPACE" create secret generic "${RELEASE}-image-pull-secret" \
    --type=kubernetes.io/dockerconfigjson \
    --from-file=.dockerconfigjson="$dockercfg" \
    --dry-run=client -o yaml | kubectl apply -f - >/dev/null
  rm -f "$dockercfg"
  log "created: ${RELEASE}-image-pull-secret"
else
  log "SKIPPED — no registry creds given. Set REPO1_USER/REPO1_TOKEN + GITLAB_USER/GITLAB_TOKEN,"
  log "or pre-create '${RELEASE}-image-pull-secret' yourself, or images will ImagePullBackOff."
fi

# 6) install-agent secret (OPTIONAL — only for catalog launches) ---------------
section "6/9 install-agent secret (optional — catalog launches)"
if [ -n "$IA_GITLAB_TOKEN" ]; then
  apply_secret "${RELEASE}-install-agent" generic "${RELEASE}-install-agent" \
    --from-literal=GITLAB_TOKEN="$IA_GITLAB_TOKEN" \
    --from-literal=ENBUILD_REPO1_USER="$IA_REPO1_USER" \
    --from-literal=ENBUILD_REPO1_TOKEN="$IA_REPO1_TOKEN"
else
  log "SKIPPED — no IA_GITLAB_TOKEN. Hub login/dashboard works without it;"
  log "catalog LAUNCHES need it (enbuildBk.installAgent.existingSecret)."
fi

# 7) teardown Reaper secret (OPTIONAL — cluster DESTROY / ELB reclaim) ----------
# Dedicated reaper-svc Secret (enbuildBk.reaperSvc.existingSecret) with the complete
# same-account delete-cred set. All three core keys go together: the BE HALTS the
# residual sweep if REAPER_SVC_AWS_ACCOUNT_ID is absent or mismatched (account-scope
# safety contract). Created only when REAPER_ACCESS_KEY_ID is set; the optional
# resource-tag pair is added only when BOTH are provided.
section "7/9 teardown Reaper secret (optional — cluster destroy / ELB reclaim)"
if [ -n "$REAPER_ACCESS_KEY_ID" ]; then
  if [ -z "$REAPER_ACCOUNT_ID" ]; then
    log "WARNING — REAPER_ACCESS_KEY_ID set but REAPER_ACCOUNT_ID is EMPTY. The BE"
    log "requires the 12-digit account id whenever delete creds are configured, or"
    log "teardown HALTS at the sweep. Set REAPER_ACCOUNT_ID and re-run."
  fi
  set -- generic "${RELEASE}-reaper-svc" \
    --from-literal=REAPER_SVC_AWS_ACCESS_KEY_ID="$REAPER_ACCESS_KEY_ID" \
    --from-literal=REAPER_SVC_AWS_SECRET_ACCESS_KEY="$REAPER_SECRET_ACCESS_KEY" \
    --from-literal=REAPER_SVC_AWS_ACCOUNT_ID="$REAPER_ACCOUNT_ID"
  if [ -n "$REAPER_RESOURCE_TAG_KEY" ] && [ -n "$REAPER_RESOURCE_TAG_VALUE" ]; then
    set -- "$@" --from-literal=REAPER_SVC_RESOURCE_TAG_KEY="$REAPER_RESOURCE_TAG_KEY" \
                --from-literal=REAPER_SVC_RESOURCE_TAG_VALUE="$REAPER_RESOURCE_TAG_VALUE"
  fi
  apply_secret "${RELEASE}-reaper-svc" "$@"
  log "set enbuildBk.reaperSvc.existingSecret=${RELEASE}-reaper-svc"
else
  log "SKIPPED — no REAPER_ACCESS_KEY_ID. Hub + launches work without it; cluster"
  log "DESTROY needs it (enbuildBk.reaperSvc.existingSecret) or ELBs orphan on teardown."
fi

# 8) observability/SIEM bearer tokens (OPTIONAL — SOO §1.5) ---------------------
section "8/9 observability/SIEM secret (optional — SIEM/Loki/Prometheus tokens)"
if [ -n "$OBS_SIEM_AUTH_HEADER$OBS_PROMETHEUS_TOKEN$OBS_LOKI_TOKEN" ]; then
  # Only include the keys actually provided (any subset).
  set -- generic "${RELEASE}-observability"
  [ -n "$OBS_SIEM_AUTH_HEADER" ] && set -- "$@" --from-literal=SIEM_AUTH_HEADER="$OBS_SIEM_AUTH_HEADER"
  [ -n "$OBS_PROMETHEUS_TOKEN" ] && set -- "$@" --from-literal=PROMETHEUS_TOKEN="$OBS_PROMETHEUS_TOKEN"
  [ -n "$OBS_LOKI_TOKEN" ]       && set -- "$@" --from-literal=LOKI_TOKEN="$OBS_LOKI_TOKEN"
  apply_secret "${RELEASE}-observability" "$@"
  log "set enbuildBk.observability.existingSecret=${RELEASE}-observability + the"
  log "non-secret endpoints (siem.endpoint/loki.host/...) in your values."
else
  log "SKIPPED — no OBS_* tokens. Spoke metrics auto-route via the agent (no token);"
  log "set OBS_SIEM_AUTH_HEADER / OBS_LOKI_TOKEN / OBS_PROMETHEUS_TOKEN to wire those."
fi

# 8) audit-export signing key (CCM-32 Auditable) ------------------------------
# ECDSA P-256 private key the backend uses to sign the /audit/export-bundle
# checksums (auditExport.service.ts → SIEM_SIGNING_KEY). Without it the bundle
# returns signed:false with an honest verifyHint. Auto-generated + IDEMPOTENT:
# kept on re-run so the public key an auditor pinned stays stable (FORCE=1 to
# rotate). Wired via enbuildBk.exportSigning.existingSecret in all 3 postures.
section "9/9 audit-export signing key (SIEM_SIGNING_KEY, ECDSA P-256)"
if secret_exists "${RELEASE}-export-signing" && [ "$FORCE" != "1" ]; then
  log "exists, kept: ${RELEASE}-export-signing  (FORCE=1 to rotate the signing key)"
elif command -v openssl >/dev/null 2>&1; then
  signkey="$(mktemp "${TMPDIR:-/tmp}/.enbuild-signkey.XXXXXX")"; chmod 600 "$signkey"
  openssl ecparam -name prime256v1 -genkey -noout -out "$signkey"
  kubectl -n "$NAMESPACE" create secret generic "${RELEASE}-export-signing" \
    --from-file=SIEM_SIGNING_KEY="$signkey" \
    --dry-run=client -o yaml | kubectl apply -f - >/dev/null
  rm -f "$signkey"
  log "created: ${RELEASE}-export-signing  (enbuildBk.exportSigning.existingSecret)"
else
  log "SKIPPED — openssl not found; /audit/export-bundle returns signed:false (verifyHint)."
fi

section "Done. Secrets in ns/$NAMESPACE:"
kubectl -n "$NAMESPACE" get secret | grep -E "^${RELEASE}-(mongo-secrets|encryption-key|messaging|rabbitmq-creds|image-pull-secret|install-agent|reaper-svc|observability|export-signing)" || true
cat <<EOF

Next: helm upgrade --install ${RELEASE} . -n ${NAMESPACE} --create-namespace \\
        -f examples/values-<your-cloud>.yaml \\
        -f ../../deploy/customer-prototype/values-images.yaml \\
        [-f examples/dev/values-<your-cloud>-eval.yaml]

  Pick the overlay(s) for your target cloud (layer the image pins + the -eval overlay LAST):
    AKS:  -f examples/values-aks.yaml [-f examples/dev/values-aks-eval.yaml]
    GKE:  -f examples/values-gke.yaml [-f examples/dev/values-gke-eval.yaml]
    port-forward quickstart (no edge): -f examples/values-quickstart.yaml
  ALWAYS also: -f ../../deploy/customer-prototype/values-images.yaml (current validated image pins)

  Full deploy walkthrough (which overlay, connect-back, CI vars): docs/DEPLOY-HUB-PER-CLOUD.md
EOF
