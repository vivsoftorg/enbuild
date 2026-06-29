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
section "1/6 MongoDB credentials"
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
section "2/6 Backend at-rest ENCRYPTION_KEY"
apply_secret "${RELEASE}-encryption-key" generic "${RELEASE}-encryption-key" \
  --from-literal=ENCRYPTION_KEY="$(rand 48)"

# 3+4) RabbitMQ broker password + erlang cookie + backend connection string ----
# ONE password drives both the broker (rabbitmq-password) and the backend's
# RABBIT_MQ_CONNECTION_STRING — they can never diverge.
section "3/6 RabbitMQ credentials (broker) + 4/6 backend connection string"
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
section "5/6 Image-pull secret (${GITLAB_REGISTRY} + ${REPO1_REGISTRY})"
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
section "6/6 install-agent secret (optional — catalog launches)"
if [ -n "$IA_GITLAB_TOKEN" ]; then
  apply_secret "${RELEASE}-install-agent" generic "${RELEASE}-install-agent" \
    --from-literal=GITLAB_TOKEN="$IA_GITLAB_TOKEN" \
    --from-literal=ENBUILD_REPO1_USER="$IA_REPO1_USER" \
    --from-literal=ENBUILD_REPO1_TOKEN="$IA_REPO1_TOKEN"
else
  log "SKIPPED — no IA_GITLAB_TOKEN. Hub login/dashboard works without it;"
  log "catalog LAUNCHES need it (enbuildBk.installAgent.existingSecret)."
fi

section "Done. Secrets in ns/$NAMESPACE:"
kubectl -n "$NAMESPACE" get secret | grep -E "^${RELEASE}-(mongo-secrets|encryption-key|messaging|rabbitmq-creds|image-pull-secret|install-agent)" || true
cat <<EOF

Next: helm upgrade --install ${RELEASE} . -n ${NAMESPACE} --create-namespace \\
        -f examples/values-quickstart.yaml
EOF
