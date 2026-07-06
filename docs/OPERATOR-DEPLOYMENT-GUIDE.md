# ENBUILD Hub — Operator Deployment Guide

How to deploy the opinionated ENBUILD hub chart (`enbuild` ≥ 0.1.0).

The platform is opinionated: it decides **how** it runs; you only bind **where**
it runs. Credentials are **never** placed in Helm values.

There are two distinct kinds of credentials — keep them separate:

- **Operational credentials** (GitLab/SCM connections for catalogs, cloud creds
  for cluster launches, etc.) are entered by the admin **inside the ENBUILD
  platform** (Connections / admin UI) and stored encrypted in the hub database.
  They are **not** part of this chart and require no Helm config.
- **Bootstrap infrastructure secrets** (below) are the few things the pods need
  to start *before* the platform is running — the at-rest encryption key, the
  database and broker credentials, and the image-pull secret. These are the only
  credentials the chart references, via pre-created Kubernetes Secrets.

- Customer-tunable values: see the `[CUSTOMER]`-tagged entries in
  `charts/enbuild/values.yaml` (≈25–30 knobs).
- Why the surface is small: `docs/VALUES-OPINIONATION-AUDIT.md`.
- Example environment override: `examples/enbuild/values-vendor13-ib.yaml`.

---

## 1. Bootstrap infrastructure secrets (create these FIRST)

The few credentials the pods need to start. Create them in the release namespace
**before** `helm upgrade --install`. Replace every `<PLACEHOLDER>` — do not commit
real values to git. (Operational SCM/cloud credentials are NOT here — those are
entered in the ENBUILD admin UI; see the note above.)

| Secret (default name) | Keys | Powers | Required? |
|---|---|---|---|
| `enbuild-install-agent-creds` | `GITLAB_TOKEN`, `ENBUILD_REPO1_USER`, `ENBUILD_REPO1_TOKEN` | Bootstrap token the backend uses to clone the agent chart repo when **installing the agent** onto a spoke, plus Iron Bank pulls. | **Yes** |
| `enbuild-encryption-key` | `ENCRYPTION_KEY` | At-rest encryption of stored secrets/state in the hub DB. | **Yes** |
| `enbuild-mongo` | `MONGO_INITDB_ROOT_USERNAME`, `MONGO_INITDB_ROOT_PASSWORD`, `MONGO_INITDB_DATABASE`, `MONGO_SERVER` | MongoDB authentication (catalog/launch state store). | **Yes** |
| `enbuild-rabbitmq` | `rabbitmq-password`, `rabbitmq-erlang-cookie`, `RABBIT_MQ_CONNECTION_STRING` | Broker auth (server + backend connection). | Yes (bundled broker) |
| `private-registry` (or your name) | `.dockerconfigjson` | Pulling first-party + dependency images from a private/Iron Bank registry. | Yes (private registry) |
| `enbuild-security-tooling` (optional) | Twistlock/Anchore/Falco env | hub-self security views. | Optional |
| `enbuild-export-signing` (optional) | `SIEM_SIGNING_KEY` | Signed `/audit/export-bundle`. | Optional |
| `<release>-keycloak-secrets` | `realm-enbuild.json`, `KC_BOOTSTRAP_ADMIN_PASSWORD` | Self-hosted Keycloak SSO (rendered when `keycloak.enabled=true`). Realm + admin password — never in git. | Yes when `keycloak.enabled` |

> Keycloak (`keycloak.enabled`) and the PKI ClusterIssuer (`pki.recreateHubIssuer`)
> are deployed by this chart. The Keycloak realm/admin Secret is created out-of-band:
> import your own realm export into a `<release>-keycloak-secrets` Secret, or — for
> demo/eval only — let the chart seed the bundled demo realm with
> `keycloak.demoRealm.enabled=true` (see `charts/enbuild/examples/values-quickstart.yaml`).
> `pki.recreateHubIssuer` restores the `enbuild-hub-issuer` CA ClusterIssuer pointing
> at the existing CA — it does not rotate certs (CA bootstrap: see
> `charts/enbuild/docs/TIER2-LAUNCH-CAPABLE-STANDUP.md` §3).

### Create commands (replace placeholders)

```bash
NS=enbuild

# Admin credential bundle — GitLab PAT (catalog/GitOps/launches) + Iron Bank pull creds
kubectl -n "$NS" create secret generic enbuild-install-agent-creds \
  --from-literal=GITLAB_TOKEN='<gitlab PAT, read/write repos+registry>' \
  --from-literal=ENBUILD_REPO1_USER='<registry1.dso.mil user>' \
  --from-literal=ENBUILD_REPO1_TOKEN='<registry1.dso.mil token>'

# At-rest encryption key
kubectl -n "$NS" create secret generic enbuild-encryption-key \
  --from-literal=ENCRYPTION_KEY='<at-rest key>'

# MongoDB (external/managed example)
kubectl -n "$NS" create secret generic enbuild-mongo \
  --from-literal=MONGO_INITDB_ROOT_USERNAME='<mongo user>' \
  --from-literal=MONGO_INITDB_ROOT_PASSWORD='<mongo password>' \
  --from-literal=MONGO_INITDB_DATABASE='enbuild' \
  --from-literal=MONGO_SERVER='<host:port (unused when mongo_endpoint_override is set)>'

# RabbitMQ (bundled broker). rabbitmq-password / rabbitmq-erlang-cookie are read
# by the bitnami subchart; RABBIT_MQ_CONNECTION_STRING is read by the backend.
kubectl -n "$NS" create secret generic enbuild-rabbitmq \
  --from-literal=rabbitmq-password='<broker password>' \
  --from-literal=rabbitmq-erlang-cookie='<clustering cookie>' \
  --from-literal=RABBIT_MQ_CONNECTION_STRING='amqp://admin:<broker password>@<release>-rabbitmq:5672/'

# Private registry pull secret (skip if your cluster injects one via webhook)
kubectl -n "$NS" create secret docker-registry private-registry \
  --docker-server=registry1.dso.mil \
  --docker-username='<registry user>' \
  --docker-password='<registry token>'
```

> The keys `rabbitmq-password` / `rabbitmq-erlang-cookie` contain hyphens, so the
> backend's `envFrom` skips them as invalid env-var names (by design) and reads
> only `RABBIT_MQ_CONNECTION_STRING`. The bitnami subchart reads the hyphenated
> keys via `auth.existingPasswordSecret` / `auth.existingErlangSecret`.

---

## 2. Bind the environment values

Copy an example and edit only the `[CUSTOMER]` knobs:

```bash
cp examples/enbuild/values-vendor13-ib.yaml my-values.yaml
```

> The example above is the **AWS/Iron Bank** hub. **AKS and GKE** hub deploys
> instead layer the per-CSP overlays `examples/values-aks.yaml` /
> `examples/values-gke.yaml` (plus the matching `examples/values-<csp>-eval.yaml`
> for demo/eval), as in `helm ... -f examples/values-<csp>.yaml [-f examples/values-<csp>-eval.yaml]`.
> See the per-cloud entry point [`DEPLOY-HUB-PER-CLOUD.md`](DEPLOY-HUB-PER-CLOUD.md)
> for the full command and prerequisites. Spoke onboarding for those clouds is set
> up via the `platform-one-gke` / `platform-one-azure` catalog repos.

Minimum bindings:
- `global.domain`, edge (`global.istio.*` **or** `global.ingress.*`), `global.imagePullSecretName`, `global.storageClass`
- `enbuildUi.hostname`
- `enbuildBk.grpcVirtualService.host` (+ `gateway`) and `enbuildBk.installAgent.hubUrl`
- The `*.existingSecret` references to the Secrets from step 1
- `mongodb.enabled` / `mongodb.mongo_endpoint_override`, `rabbitmq.auth.username`
- Per-service `image.tag` (+ `image.registry` for air-gap)

---

## 3. Deploy

```bash
helm upgrade --install <release> ./charts/enbuild \
  --namespace enbuild --create-namespace \
  --values my-values.yaml \
  --kube-context <context>
```

Validate before applying to a live cluster:

```bash
helm diff upgrade <release> ./charts/enbuild -n enbuild -f my-values.yaml   # requires helm-diff
helm template <release> ./charts/enbuild -n enbuild -f my-values.yaml | less
```

---

## 3b. Observability & SIEM (SOO §1.5) — repeatable, resilient config

The three observability surfaces have **three different** config models. All are
declarable in values (GitOps-repeatable) so you never hand-edit the Deployment
(which a `helm upgrade` wipes). Secrets ride one operator-managed Secret; the
non-secret endpoints are plain values.

| Capability | Test cell | How to configure | Notes |
|---|---|---|---|
| **Performance — managed clusters** | 3.5-5 | **Nothing.** The connected spoke agent auto-discovers the cluster's Big Bang/kube-prometheus and routes metrics over the existing mTLS link. | Leave `enbuildBk.observability.prometheus.host` EMPTY. If the Performance tab shows synthetic, check BB monitoring is Running on that spoke + the agent is connected — **not** a hub setting. |
| **SIEM forwarding** | 3.5-1 | `enbuildBk.observability.siem.*` (endpoint/format/toggle) + `SIEM_AUTH_HEADER` in the Secret. Private-CA collector: `enbuildBk.observability.siem.caBundle.{existingConfigMap\|existingSecret,key}`. | Ships **OFF** until `forwardingEnabled: "true"` + an `endpoint`. Forwards the **solution audit stream** (CCM-32 platform actions), *not* aggregated cluster/pod logs — those ship from each cluster's Big Bang Fluentbit (see §3c). Admin UI `/admin/siem-settings` is the **live per-field override**; these values are the GitOps floor + the only air-gap-seedable path. If the collector's HTTPS cert is signed by a private/agency CA the backend doesn't already trust, the POST fails "unable to verify the first certificate" (no insecure-skip, by decision) — set `caBundle` (see Step 1b). |
| **Troubleshooting (pod logs)** | 3.5-4 | `enbuildBk.observability.loki.{host,tenant}` + `LOKI_TOKEN` in the Secret; optional `grafana.{base,lokiDatasource}` deep-link. | Without a host, in-console log tailing falls back to a synthetic zero-state + the Grafana deep-link. |

**Step 1 — create the bearer-token Secret** (any subset of the three keys; the
SIEM/Loki tokens are customer-furnished, so this is NOT auto-generated):

```bash
kubectl -n enbuild create secret generic enbuild-ib-observability \
  --from-literal=SIEM_AUTH_HEADER='Bearer <p1-siem-collector-token>' \
  --from-literal=LOKI_TOKEN='<loki-token>'
# (or: OBS_SIEM_AUTH_HEADER=... OBS_LOKI_TOKEN=... scripts/create-bootstrap-secrets.sh)
```

**Step 1b — (only if the SIEM collector uses a private/agency CA)** seed that CA
so the backend's Node runtime trusts the collector's TLS chain. Put the CA PEM in
a ConfigMap (or Secret) in the `enbuild` namespace, then point `caBundle` at it in
Step 2. There is **no** insecure-skip-verify path — this is by decision.

```bash
kubectl -n enbuild create configmap siem-collector-ca --from-file=ca.crt=./siem-ca.pem
```

When set, the chart mounts that PEM read-only at `/etc/enbuild/siem-ca/<key>` and
sets `NODE_EXTRA_CA_CERTS` to it. `NODE_EXTRA_CA_CERTS` is read at **process boot**,
so a helm upgrade that adds/changes `caBundle` rolls the backend pod automatically
(the env + volume are part of the pod template) — no manual restart. Leaving
`caBundle` empty (default) mounts nothing and is byte-for-byte today's behavior.

**Step 2 — set the non-secret endpoints in your values:**

```yaml
enbuildBk:
  observability:
    existingSecret: enbuild-ib-observability   # SIEM_AUTH_HEADER / PROMETHEUS_TOKEN / LOKI_TOKEN
    siem:
      forwardingEnabled: "true"
      endpoint: "https://siem.p1.example.mil/services/collector/raw"
      format: "json"                # ECS JSON (default); "cef" for Splunk/ArcSight
      caBundle:                      # ONLY if the collector uses a private/agency CA (Step 1b)
        existingConfigMap: siem-collector-ca   # or existingSecret: <name>
        key: ca.crt
    loki:
      host: "http://logging-loki-gateway.logging.svc.cluster.local"
      tenant: "enbuild"
    grafana:
      base: "https://grafana.p1.example.mil"
      lokiDatasource: "loki"
    # prometheus.host: only for HUB-SELF metrics; leave empty for spoke clusters.
```

**Repeatability lynchpin — pin the at-rest key.** Admin-UI settings (incl. SIEM
entered at `/admin/siem-settings`) are encrypted with `ENCRYPTION_KEY`. If that
key regenerates on a redeploy, those settings become unreadable and silently
revert. ALWAYS set `enbuildBk.encryptionKey.existingSecret` to a persistent
Secret (the chart NOTES warns at install if it is empty). Declaring SIEM/Loki/
Prometheus in values as above is unaffected by the key — that is the resilient,
audit-reviewable path and is recommended over admin-UI click-ops for production.

**Resilience built in:** SIEM POSTs retry on 5xx/network failure (queue + cron,
exponential backoff, dead-letter cap) and hot-reload from the admin UI within
~15s; the hub-self Prometheus client has a 5s per-attempt timeout + one bounded
retry; every surface fails safe to a labeled synthetic zero-state, never
fabricated data.

## 3c. Managed-resource logs → SIEM (SOO §1.5(i), "and managed resources")

The §3b SIEM toggle ships the **solution** audit stream. The SOO also requires
**managed-resource/cluster logs** in the P1 SIEM. Those are shipped per managed
cluster by Big Bang's log collector (Fluentbit/Promtail) configured with the P1
SIEM as an output — set in the Big Bang catalog values, not this hub chart. See
the platform-one-bigbang catalog (`bigbang/envs/<env>/values/`) Fluentbit output
config. (Tracked as the remaining SOO §1.5(i) "Stream B" delivery item.)

---

## 3d. Spoke agent connect-back — hub CA trust (REQUIRED on a self-signed hub)

The hub's agent gRPC gateway serves a cert signed by the hub's **private CA**
(`enbuild-hub-issuer` / `CN=enbuild-hub-ca`, created by this chart's cert-manager
templates — ADR-0004). Every spoke agent (whether **launched** via the catalog or
**imported** via *Import existing cluster*) dials that gateway with cert
verification ON (`ENBUILD_AGENT_HUB_INSECURE=false`), so **the agent must trust that
CA or it silently fails the TLS dial** — the symptom is `grpc.Dial(<hub>): context
deadline exceeded` and the cluster never leaves "Stalled"/empty-inventory even though
`nc`/`curl` to the hub succeed (it's a cert-trust failure, not a network one).

**Operator step — set the hub CA source ONCE** so both onboarding paths distribute it
to spokes automatically (backend `catalog.service` for launch + `installAgent.service`
for import both read it and write the `enbuild-spoke-ca-bundle` ConfigMap + enable the
agent chart's `agent.hubCaBundle`):

```bash
# Base64 PEM of the hub gateway CA (from the chart-created cert secret):
CA=$(kubectl -n istio-system get secret enbuild-hub-grpc-tls -o go-template='{{index .data "ca.crt"}}')
# Preferred: set it in the encrypted AdminSettings via the admin UI/API
#   AdminSettings → agent → HUB_CA_BUNDLE_PEM = "$CA"   (survives rolls)
# Or as a backend env in your values (helm): enbuild-bk env ENBUILD_HUB_CA_BUNDLE_PEM="$CA"
```

**When to SKIP it:** a hub fronted by a **public/real cert** (e.g. a Big Bang
`public-ingressgateway`) needs nothing — leave the CA source unset and the agent uses
the system trust store (the provisioning is guarded off, byte-unchanged). Only set it
for a self-signed / private-CA hub gateway. Blast radius is narrow: the CA affects
*only* the agent→hub dial (the agent's K8s-API and other TLS use their own CAs), so
setting the wrong CA or setting it on a public-cert hub is the only failure mode —
don't.

---

## 4. Migrating the live vendor13-ib release (943-line override → lean)

The chart is backwards-compatible, so this is value-preserving if you create the
Secrets with the **current** values.

1. **Rotate the leaked credentials out-of-band first** (they were committed in
   the old override): the `cms553` Iron Bank registry credential and the
   `glpat-…` GitLab PAT. Put the rotated values into the Secrets above.
2. Create `enbuild-encryption-key`, `enbuild-mongo`, `enbuild-rabbitmq` using the
   **values currently in effect** so nothing breaks. Retrieve current values
   (do NOT commit them):
   ```bash
   kubectl -n enbuild get secret enbuild-ib-backend-secret \
     -o jsonpath='{.data.ENCRYPTION_KEY}' | base64 -d ; echo
   kubectl -n enbuild get secret enbuild-ib-backend-secret \
     -o jsonpath='{.data.RABBIT_MQ_CONNECTION_STRING}' | base64 -d ; echo
   ```
   - `ENCRYPTION_KEY`: **must equal** the current value — changing it re-keys
     at-rest data (treat any change as a separate data migration).
   - Mongo password: the real value is the one embedded in the current
     `mongo_endpoint_override` URI (it differs from the unused `mongo_root_password`
     field — a known footgun).
   - `enbuild-install-agent-creds` and `private-registry` already exist on the
     cluster; only rotate their contents.
3. `helm upgrade enbuild-ib` with `examples/enbuild/values-vendor13-ib.yaml`.
   The rendered core Deployments are byte-equivalent to the current release
   except: secrets sourced from the new Secrets, `imagePullSecrets: private-registry`,
   `CLUSTER_RPC_TIMEOUT_MS` corrected `10000`→`30000`, and the unused
   `image-pull-secret` / `mongo-secrets` chart Secrets no longer rendered.

---

## 5. What this chart no longer ships (removed in 0.1.0)

JupyterHub, open-webui/Ollama, the AI proxy, CTF, Bolt, and loki-stack were
removed (they were gated OFF in every environment). Headlamp is the only
optional bundled subchart (`lightning_features.operations_lightning.headlamp`).
If you need any removed capability, deploy it as a separate, supported chart.
