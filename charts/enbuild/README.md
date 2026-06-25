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
| `<release>-keycloak-secrets` | `realm-enbuild.json`, `KC_BOOTSTRAP_ADMIN_PASSWORD` | bundled Keycloak SSO realm import | only if `keycloak.enabled` **and** you bring your own realm — `keycloak.demoRealm.enabled=true` auto-creates this for you (see §3a) |

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

## 3a. Console login & multi-tenancy (Keycloak SSO)

The console authenticates every user through **Keycloak** (OIDC Authorization-Code
+ PKCE, full-page redirect) and derives **per-project tenancy from Keycloak group
membership** — no platform-specific user store. The backend verifies tokens via
**JWKS** fetched from Keycloak in-cluster, so there is **no realm public key to
paste anywhere**. Pick one of two paths:

### Path A — bundled Keycloak + demo realm (deploy & log in in minutes)

For evaluating the platform. The chart deploys Keycloak *and* auto-imports a demo
realm with ready-made personas and projects — set **one host** and a flag:

```yaml
keycloak:
  enabled: true
  demoRealm:
    enabled: true                 # auto-seed the demo realm (personas + groups)
  hostname: https://kc.<domain>    # browser-reachable Keycloak URL (KC_HOSTNAME)
enbuildUi:
  hostname: enbuild                # console FQDN = enbuild.<global.domain>
  keycloak:
    url: https://kc.<domain>        # same host → the SPA's config.json
enbuildBk:
  keycloak:
    url: https://kc.<domain>        # same host → backend issuer auto-derives to <url>/realms/enbuild
    backendUrl: http://<release>-keycloak:8080   # in-cluster Keycloak service (JWKS/discovery)
```

`helm install`, open the console, log in (all personas password `ChangeMe123!`):

| Persona | Sees | Proves |
|---|---|---|
| `admin@p1.mil` | the whole fleet | platform admin (`/enbuild/platform-admins`) |
| `owner-bb@p1.mil` | **only** the `big-bang` project | project-scoped owner |
| `owner-ib@p1.mil` | **only** `iron-bank` | cross-project isolation (denied `big-bang`) |
| `viewer-ib@p1.mil` | `iron-bank`, read-only | viewer can read, not write |

> **Demo only.** `demoRealm` ships no signing keys/secrets (Keycloak generates the
> key; the backend uses JWKS) but the persona passwords + bootstrap admin are
> well-known defaults — **change them, or use Path B, before production.** Setting
> `enbuildBk.keycloak.existingSecret` (your own realm) disables the demo seed.

### Path B — bring your own Keycloak / IdP (production, incl. P1 SSO / CAC)

Point the same values at *your* Keycloak instead of the bundled one (set
`keycloak.enabled: false`), then in your realm create:

- a **public** client `enbuild-ui` — standard-flow + direct-access + PKCE(S256), a
  **group-membership protocol mapper** emitting **full group paths** in a `groups`
  claim (mandatory — without it tenancy sees no groups), and the console host in
  `redirectUris`/`webOrigins`;
- a **confidential** client `enbuild` for the backend service account;
- groups `/enbuild/platform-admins` and `/enbuild/projects/<project>/<role>`
  (`role` ∈ `viewer|member|maintainer|owner`), and assign your users.

The backend verifies via JWKS at `enbuildBk.keycloak.backendUrl` — nothing to pin.
This is the path you use to integrate ENBUILD with an existing enterprise IdP.

**Why it works the same everywhere:** one Keycloak (bundled or yours), the SPA
redirects to it at its own host, the backend verifies its tokens via JWKS, and
tenancy comes from the `groups` claim. The only hard rule is **host alignment** —
`keycloak.hostname`, `enbuildUi.keycloak.url`, and `enbuildBk.keycloak.url` must be
the same browser-facing Keycloak URL (the backend issuer is `<url>/realms/enbuild`,
matched exactly).

### Headlamp cluster console (OIDC)

The embedded **Headlamp** console browses spoke clusters by forwarding *the logged-in
user's* Keycloak token to the hub proxy — so under `CONSOLE_AUTH_STRICT` it needs the
same OIDC identity as the SPA. It is **off by OIDC default** (all four fields empty);
enable it by setting **all four together**:

```yaml
lightning_features:
  operations_lightning:
    headlamp: true                  # embedded cluster console (default on)
headlamp:
  config:
    oidc:
      clientID: enbuild-ui          # SAME public PKCE client the console uses
      scopes: "openid email groups"  # `groups` is REQUIRED — Headlamp tenancy reads it
      issuerURL: https://kc.<domain>/realms/enbuild     # MUST byte-match the backend KEYCLOAK_ISSUER
      callbackURL: https://enbuild.<domain>/headlamp/oidc-callback   # console-host + /headlamp/oidc-callback
```

Also register that `callbackURL` as a redirect URI on the realm's `enbuild-ui` client
(Path A's demo realm already includes it for the demo console host).

> **Fail-closed:** the chart refuses to install on the two silent misconfigurations
> that would otherwise break login with no error — a **partial** config (some but not
> all four fields set) and an `issuerURL` that doesn't equal the backend
> `KEYCLOAK_ISSUER` exactly. Leave all four empty to run Headlamp without SSO.

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
