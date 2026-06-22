# ENBUILD Helm Chart

Installs the [ENBUILD hub platform](https://gitlab.com/enbuild-staging/vivsoft-platform-ui) —
the control plane for ENBUILD's hub-and-spoke multi-cluster manager: the backend
gRPC API, the console UI, the mq launch worker, the user service, MongoDB and
RabbitMQ, the embedded Headlamp cluster console, and (optionally) self-hosted
Keycloak SSO and the hub/spoke PKI ClusterIssuer.

> **This chart is opinionated (0.1.0+).** The platform decides *how* it runs; you
> bind only *where* it runs. The customer-facing surface is ~30 values — the
> entries tagged `[CUSTOMER]` in [`values.yaml`](values.yaml). Everything tagged
> `[INTERNAL]` is platform wiring with one correct value; don't override it.
> Full deploy walkthrough + the exact `kubectl` Secret commands:
> [`docs/OPERATOR-DEPLOYMENT-GUIDE.md`](../../docs/OPERATOR-DEPLOYMENT-GUIDE.md).
> Why the surface is small: [`docs/VALUES-OPINIONATION-AUDIT.md`](../../docs/VALUES-OPINIONATION-AUDIT.md).

> **Branches & release channels:** P1 CCM consumers track the long-lived trunk
> branch `feat/p1ccm-enbuild-helm-trunk`, published as GitHub Release assets
> tagged `enbuild-trunk-<version>` — separate from the `main` channel. See
> [`docs/P1CCM-TRUNK.md`](../../docs/P1CCM-TRUNK.md).
>
> **Hitting a problem?** See [TROUBLESHOOTING.md](TROUBLESHOOTING.md).

---

## 1. Prerequisites — create the bootstrap Secrets first

Credentials are **never** placed in Helm values. Pre-create these Kubernetes
Secrets in the release namespace **before** installing. (Operational credentials —
GitLab/SCM connections for catalogs, cloud creds for cluster launches — are
entered later **in the ENBUILD admin UI**, not in this chart.)

| Secret (default name) | Keys | Powers | Required |
|---|---|---|---|
| `enbuild-encryption-key` | `ENCRYPTION_KEY` | at-rest data encryption | yes |
| `enbuild-mongo` | `MONGO_INITDB_ROOT_USERNAME/PASSWORD/DATABASE`, `MONGO_SERVER` | MongoDB auth | yes |
| `enbuild-rabbitmq` | `rabbitmq-password`, `rabbitmq-erlang-cookie`, `RABBIT_MQ_CONNECTION_STRING` | broker auth | yes (bundled broker) |
| `enbuild-install-agent-creds` | `GITLAB_TOKEN`, `ENBUILD_REPO1_USER`, `ENBUILD_REPO1_TOKEN` | agent installs onto spokes + Iron Bank pulls | yes |
| image pull secret (`global.imagePullSecretName`) | `.dockerconfigjson` | pulling images from a private/Iron Bank registry | yes (private registry) |
| `<release>-keycloak-secrets` | `realm-enbuild.json`, `KC_BOOTSTRAP_ADMIN_PASSWORD` | self-hosted Keycloak SSO | only if `keycloak.enabled` |

`kubectl create secret` commands for each are in
[`docs/OPERATOR-DEPLOYMENT-GUIDE.md`](../../docs/OPERATOR-DEPLOYMENT-GUIDE.md).

## 2. Installing

**P1 CCM (trunk channel):**

```shell
# download the published chart artifact
gh release download enbuild-trunk-<version> \
  --repo vivsoftorg/enbuild --pattern '*.tgz' -D /tmp

# install / upgrade
helm upgrade --install enbuild-ib /tmp/enbuild-<version>.tgz \
  --namespace enbuild --create-namespace \
  --values my-values.yaml
```

Verify the reverse-proxy chain end to end after install:

```shell
helm test <release> -n <namespace>
```

## 3. Configuration — the customer surface

Set only what binds the platform to your environment; everything else has a
correct default baked in. The complete, authoritative list is the `[CUSTOMER]`
entries in [`values.yaml`](values.yaml). The essential ones:

| Path | What it binds | Example |
|---|---|---|
| `global.domain` | external DNS suffix for hub hosts + cert SANs | `apps.example.mil` |
| `global.istio.enabled` / `.gateway` | mesh edge (or use `global.ingress.*`) | `true` / `istio-system/main` |
| `global.image.registry` | registry first-party images pull from (air-gap) | `registry1.dso.mil` |
| `global.imagePullSecretName` | pre-created `dockerconfigjson` Secret for private pulls | `private-registry` |
| `global.storageClass` | StorageClass for stateful deps (if no cluster default) | `gp2` |
| `enbuildUi.hostname` | console FQDN = `hostname.domain` | `enbuild` |
| `enbuildBk.grpcVirtualService.host` / `.gateway` | how spoke agents reach this hub's gRPC | `enbuild.apps.example.mil` / `istio-gateway/public-ingressgateway` |
| `enbuildBk.installAgent.hubUrl` / `.existingSecret` | `host:443` agents dial back to + the admin-creds Secret | `enbuild.apps.example.mil:443` / `enbuild-install-agent-creds` |
| `enbuildBk.encryptionKey.existingSecret` | at-rest key Secret | `enbuild-encryption-key` |
| `enbuildBk.messaging.existingSecret` | broker connection Secret | `enbuild-rabbitmq` |
| `mongodb.enabled` / `.existingSecret` / `.mongo_endpoint_override` | bundle vs external Mongo (always required either way) + creds | `false` / `enbuild-mongo` / external URI |
| `rabbitmq.auth.existingPasswordSecret` / `.existingErlangSecret` | broker creds | `enbuild-rabbitmq` |
| `keycloak.enabled` | deploy bundled SSO (else use an external IdP) | `true` |
| `pki.recreateHubIssuer` | one-shot restore of a missing hub mTLS ClusterIssuer | `false` |
| `<svc>.image.tag` | pin a service image; empty tracks the chart appVersion | unset |

A complete filled-out example for a real environment:
[`examples/enbuild/values-vendor13-ib.yaml`](../../examples/enbuild/values-vendor13-ib.yaml).

## 4. What this chart deploys

backend (gRPC API) · console UI · mq launch worker · user service · MongoDB
(bundled single-node *or* point at external/HA) · RabbitMQ · Headlamp cluster
console (optional, `lightning_features.operations_lightning.headlamp`) · Keycloak
SSO (optional, `keycloak.enabled`) · hub PKI ClusterIssuer (optional,
`pki.recreateHubIssuer`).

> Removed in 0.1.0: JupyterHub, open-webui/Ollama, the AI proxy, CTF, Bolt, and
> loki-stack (they were disabled in every environment). Deploy any of those as
> separate charts if needed.

## 5. Uninstalling

```shell
helm uninstall <release> --namespace <namespace>
```

Stateful data (MongoDB/RabbitMQ volumes, operator Secrets) is intentionally not
removed; clean those up separately if you mean to fully tear down.
