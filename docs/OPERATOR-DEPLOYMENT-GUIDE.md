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
> are deployed by this chart. The Keycloak realm/admin Secret is created out-of-band
> (see `runbooks/enbuild-hub-reconcile-adopt.md`); `pki.recreateHubIssuer` restores
> the `enbuild-hub-issuer` CA ClusterIssuer pointing at the existing CA — it does
> not rotate certs.

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
