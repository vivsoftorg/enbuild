> **STATUS: IMPLEMENTED in chart 0.1.0-p1ccm-trunk.0** (rebased onto trunk
> `0.0.49-p1ccm-trunk.4`) — applied: 6 dead apps + 3 subchart deps removed, the
> public value surface cut to ~25–30 `[CUSTOMER]` knobs, and all secrets
> converted to `existingSecret` references. The on-trunk Keycloak SSO + PKI
> ClusterIssuer + nginx gzip/websocket work (trunk `.4`) are **preserved** and
> Keycloak/PKI are now first-class `[CUSTOMER]`/`[INTERNAL]` value sections
> (superseding this doc's earlier "off-trunk orphan" treatment of `keycloak.*`/
> `pki.*`). NOTE: `clusterRpcTimeoutMs` chart default is **`10000`** — that is the
> intentional 2026-06-11 latency-P0 fail-fast value, NOT a regression (this doc's
> earlier "fix to 30000" was based on the stale pre-`.4` comment and was dropped).
> Operator instructions: `docs/OPERATOR-DEPLOYMENT-GUIDE.md`. Environment
> override: `examples/enbuild/values-vendor13-ib.yaml`. Render verified:
> resource-identical to the live trunk-`.4` chart except the two intended secret
> changes (chart-rendered `image-pull-secret`/`mongo-secrets` → operator Secrets).

# ENBUILD Hub Chart — Opinionation & Surface-Reduction Recommendation (CORRECTED, trunk-verified)

> Every path in this document was grep-verified against the TRUNK chart at
> `/Users/csanchez/code/enbuild-helm.trunk/charts/enbuild` on 2026-06-17. Paths the
> trunk chart does **not** consume are called out explicitly and never recommended as
> "configure/hardcode here." This supersedes the first-pass synthesis, which conflated
> this trunk chart with the off-trunk live chart-line and used several wrong key names.

---

## 1. Executive summary + opinionation philosophy

The `enbuild` hub chart (TRUNK version `0.0.49-p1ccm-trunk.2`, appVersion `1.0.30`) is the
ENBUILD control plane: backend gRPC API (`enbuildBk`), console UI (`enbuildUi`), mq-consumer
launch worker (`enbuildConsumer`), `enbuildUser`, an optional AI proxy, MongoDB, RabbitMQ, plus
a pile of optional bundled apps. It is deployed LIVE as release `enbuild-ib` in namespace
`enbuild` on the hardened vendor13-ib (P1 / Iron Bank) cluster, where the operator currently
supplies **943 lines** of override values.

That number is the whole problem: the chart ships internal wiring, dormant-subchart config, and
plaintext credentials as public configuration, so operators must re-state dozens of
platform-internal constants (and re-paste secrets) just to stand it up.

**Philosophy: the platform decides how it works; the customer only binds where it runs.**
The hub owns every internal service name, port, cross-component URL, datastore topology, console
cache/snapshot tuning, install-agent pin, and feature-internal default. The customer binds a
small set of *environment* facts: external DNS/gateway/ingress, image registry + per-service
tags (air-gap rebind), pre-created credential **secret refs** (never literals), storageClass,
capacity, the agent gRPC routing for this hub, and a couple of console iframe URLs. Target public
surface ≈ **25–30 leaf paths**, down from the hundreds exposed today.

---

## 2. Headline findings

1. **The scale problem.** vendor13-ib overrides **943 lines**. The majority are (a) plaintext
   credentials, (b) re-statements of platform-internal wiring the chart should own, and (c)
   fully-populated config for **dormant** optional subcharts (open-webui/ollama, jupyterhub
   profileList, loki-stack, CTF, bolt) that are gated OFF and never render. This is accreditation
   cruft and attack surface, not configuration.

2. **CHART-vs-LIVE BOUNDARY (first-class finding).** A whole class of live values has **zero
   consumers in this trunk chart** — grep across `templates/` + `_helpers.tpl` returns nothing
   for: top-level `pki.*`, top-level `keycloak.*`, `enbuildBk.keycloak.*`,
   `enbuildBk.authMechanism`, `enbuildBk.adminEmails`, `enbuildBk.mtls.*`, and
   `enbuildBk.gitlabConnection.*`. These are **operator-owned / off-trunk** (fed by the
   out-of-band Keycloak operator and the off-trunk PKI chart-line per the chart-line-divergence
   record). You cannot `_helpers`-constant or allow-list a value the chart never reads. Where the
   feature *should* exist in this chart (the gRPC/pull-agent mTLS fabric; the KEYCLOAK_* env
   wiring), the action is **"PORT the templates first,"** not "configure here." See bucket C.

3. **The MOST load-bearing per-env knobs are the agent gRPC routing — and they ARE consumed.**
   `enbuildBk.grpcEnabled`, `enbuildBk.grpcVirtualService.{enabled,host,additionalHosts,gateway}`
   are consumed by `enbuild-bk-grpc-virtualservice.yaml` (:1, :23, :33, :37) and
   `enbuild-bk-grpc-service.yaml` (:1). Likewise `enbuildBk.installAgent.{hubUrl,tlsServerName,
   hubInternalIps,agentImageTag,existingSecret}` are consumed by `enbuild-bk.yaml` (:130, :132,
   :134, :136, :195). These are the canonical per-environment bindings for this hub and belong on
   the allow-list — the first pass never classified them.

4. **`clusterRpcTimeoutMs` live override REGRESSES the chart fix.** Chart default is `"30000"`
   (`values.yaml:227`), consumed at `enbuild-bk.yaml:77`; the live override sets `10000`, the very
   value the default was raised away from to fix the "Lost connection to the cluster" bug.
   SIMPLIFY/correctness — drop the override, keep the 30000 default.

5. **Worst security finding: a published at-rest encryption key.** Chart default
   `enbuildBk.encryption_key: "encryption_key"` and the live release both use the literal
   `encryption_key`, rendered verbatim as `ENCRYPTION_KEY` into the backend Secret
   (`backend-secret.yaml:18`). The hub's at-rest crypto key is a well-known shipped constant in a
   hardened P1 environment.

6. **Multiple live, committed plaintext credentials** (full table §7). Iron Bank registry creds
   (`cms553` / `VxoMdavv...`), GitLab PAT (`glpat-6PcSe...`), Mongo `SuperSecret`, RabbitMQ
   `lamba` / `SuperSecret`, Grafana `V1v$oftR0ck5`, Keycloak `clientSecret` +
   `ChangeMe123!`, AI `api_key: dummy`, JupyterHub `authenticator_class: dummy` (auth bypass).
   **Rotate `cms553` and the GitLab PAT now.**

7. **gRPC port is a footgun.** `enbuildBk.grpcPort` feeds the gRPC Service `port`
   (`enbuild-bk-grpc-service.yaml:23`, `default 8443`) but `targetPort` and the container port are
   hardcoded `8443` (`enbuild-bk-grpc-service.yaml:24`, `enbuild-bk.yaml:243`). Any non-8443 value
   silently breaks routing → HARDCODE.

8. **Bare-install lands wrong.** `enbuildBk.resources` has no default in the values block, so a
   clean install renders empty resources (`enbuild-bk.yaml:258`, `toYaml`). The opinionated chart
   must ship sane defaults.

9. **Six bundled apps are dead weight.** CTF, Bolt, JupyterHub, open-webui/ollama, and loki-stack
   are gated OFF in defaults AND live and never render. Headlamp is the only optional subchart ON
   in live (`lightning_features.operations_lightning.headlamp: true`). Drop the five, keep Headlamp.

10. **Subchart passthrough explosion.** RabbitMQ alone contributes a large fraction of the 943
    lines as pure bitnami passthrough; the parent consumes only `auth.username/password`, `host`,
    `env`, `queue_prefix` (+ `enabled`). The hub must pin upstream internals and expose nothing.

---

## 3. Three buckets (do not confuse them)

### (A) THIS CHART's customer allow-list — paths with REAL trunk consumers

Every path below was confirmed consumed by a trunk template. Items marked **(port template X
first)** have no consumer today and are *aspirational* surface — list them as future knobs, do not
tell an operator to set them on the current chart.

| Path | What it binds | Example (live) | Trunk consumer |
|---|---|---|---|
| `global.domain` | External DNS suffix for hub ingress/SANs | `staging.dso.mil` | ingress/VS templates |
| `global.istio.enabled` / `global.istio.gateway` | Mesh edge gateway | `true` / `istio-system/main` | VS/ingress templates |
| `global.ingress.enabled`/`.tls`/`.tls_secret`/`.className`/`.annotations` | Non-mesh ingress edge + TLS ref | `false` | `ingress-frontend.yaml` (see className bug §4) |
| `global.image.registry` | Air-gap / Iron Bank registry rebind | `registry1.dso.mil` | every image ref |
| `global.image.pullPolicy` | Pull behavior (air-gap → `IfNotPresent`) | `Always` | all Deployments |
| `global.imagePullSecretName` *(new)* | **Pre-created** docker-registry Secret name | — | replaces `_helpers.tpl` plaintext lane |
| `global.storageClass` | StorageClass for stateful deps | `gp2` | `mongodb.yaml:69` |
| `global.disable_tls_gitlab` | Self-signed GitLab trust | `false` | `backend-secret.yaml:13` |
| `enbuildUi.hostname` | Console external FQDN | `enbuild-ib...` | `virtual-service-frontend.yaml`, ingress |
| `enbuildBk.image.tag` / `enbuildUi.` / `enbuildConsumer.` / `enbuildUser.image.tag` | Per-service version pin / air-gap tag | SHA pins | image refs (fallback chain §4) |
| `enbuildConsumer.image.registry` | Per-service air-gap registry override | `registry.gitlab.com` | `enbuild-mq.yaml:44` |
| `enbuildBk.replicas` / `enbuildUi.replicas` / `enbuildConsumer.replicas` / `enbuildUser.replicas` | Capacity (NOTE: key is **`replicas`**, NOT `replicaCount`) | `1` | `enbuild-bk.yaml:10`, `ui:12`, `mq:10`, `user:10` |
| `enbuildBk.resources` / `enbuildConsumer.resources` | Capacity tuning (ONLY these two are wired) | requests cpu/mem | `enbuild-bk.yaml:258`, `enbuild-mq.yaml:102` |
| `enbuildBk.grpcEnabled` | Enable agent gRPC routing | `true` | `enbuild-bk-grpc-{service,virtualservice}.yaml:1`, `enbuild-bk.yaml:242` |
| `enbuildBk.grpcVirtualService.enabled` / `.host` / `.additionalHosts` / `.gateway` | **Agent gRPC ingress (the most load-bearing per-env knobs)** | `true` / `enbuild-ib-vendor13.staging.dso.mil` / `[...]` / `istio-gateway/public-ingressgateway` | `enbuild-bk-grpc-virtualservice.yaml:1,23,33,37` |
| `enbuildBk.installAgent.hubUrl` / `.tlsServerName` / `.hubInternalIps` | Hub address the agent dials back to | `enbuild-ib-vendor13.staging.dso.mil:443` | `enbuild-bk.yaml:130,132,134` |
| `enbuildBk.installAgent.existingSecret` | **Pre-created** install-agent creds Secret | `enbuild-install-agent-creds` | `enbuild-bk.yaml:195`, `NOTES.txt:16` |
| `enbuildBk.kubeProxyFallbackActor` | Audit actor for proxy fallback | `devops@platform.one.mil` | `enbuild-bk.yaml:66` |
| `enbuildBk.encryptionKey.existingSecret` *(new — replaces literal)* | **Pre-created** Secret holding `ENCRYPTION_KEY` | — | requires §4 secret-ref conversion |
| `enbuildBk.securityTooling.existingSecret` | **Pre-created** security-tooling creds Secret | `""` | `enbuild-bk.yaml:180` |
| `enbuildBk.exportSigning.existingSecret` | **Pre-created** SIEM signing-key Secret (prod path) | `""` | `enbuild-bk.yaml:163`, `backend-secret.yaml:78` |
| `enbuildUi.kiali_url` / `enbuildUi.kubecost_url` | Console iframe URLs (env-specific dashboards) | live-set | `enbuild-ui.yaml:31,37`, `frontend-config.yaml:27,29` |
| `mongodb.enabled` | The ONE Mongo topology switch | `false` (managed) | `mongodb.yaml`, secrets |
| `mongodb.mongo_endpoint_override` | External/managed Mongo URI (pw via `$(...)` env) | DocumentDB/cosmos URI | `enbuild-bk.yaml:51-55` |
| `mongodb.existingSecret` *(new — replaces literals)* | **Pre-created** Mongo root credential Secret | — | requires §4 secret-ref conversion |
| `mongodb.storageClassName` | Mongo PVC storageClass (falls back to global) | — | `mongodb.yaml:69` |
| `mongodb.image.repository` / `.tag` | Mongo image (air-gap rebind, when self-hosted) | — | `mongodb.yaml:26` |
| `rabbitmq.enabled` | RabbitMQ topology switch | `true` | Chart.yaml condition |
| `rabbitmq.auth.username` / `rabbitmq.host` / `rabbitmq.env` / `rabbitmq.queue_prefix` | The 4 leaves the PARENT consumes | — | `backend-secret.yaml:25,27,29,30` |
| `rabbitmq.auth.existingPasswordSecret` / `.existingErlangSecret` | **Pre-created** RabbitMQ Secrets | — | bitnami subchart |
| `rabbitmq.persistence.size` | RabbitMQ volume sizing (REAL; mongo has none) | `8Gi` | bitnami subchart |
| `rabbitmq.image.registry`/`.repository`/`.tag` | Air-gap rebind + version pin (chart default `3.11.13` is **stale** vs live `3.12.14`) | `registry1.dso.mil/.../rabbitmq` | bitnami subchart |
| `headlamp.image.registry`/`.repository`/`.tag` / `.spokes` / `.replicaCount` / `.resources` | Only retained optional subchart: image (air-gap), spoke topology, capacity | live-set | `headlamp-kubeconfigs-configmap.yaml` (spokes); rest = subchart |
| `enbuildBk.gitlabConnection.host` **(port template first)** | GitLab connection endpoint | — | *no consumer in trunk yet* |
| `enbuildBk.gitlabPat.existingSecret` | **Pre-created** Secret with the GitLab PAT | — | `backend-secret.yaml:35` (gates inline lane) |
| `enbuildBk.keycloak.connectionUrl`/`.realm`/`.clientId`/`.existingSecret` **(port template first)** | Keycloak SSO connection | `enbuild` / `enbuild` | *no consumer in trunk* (bucket C) |
| `enbuildBk.adminEmails` **(port template first)** | Bootstrap admin allowlist | — | *no consumer in trunk* (bucket C) |
| `pki.hubServerCert.dnsNames` **(port template first)** | External SAN(s) for hub gRPC cert | — | *no consumer in trunk* (bucket C) |

> `replicaCount` is the correct key ONLY for `rabbitmq.replicaCount` and `headlamp.replicaCount`
> (genuine bitnami/upstream passthrough). The four core services use **`replicas`**.

> There is **no `aws.region`** path anywhere in chart or live — the only region is
> `enbuildCTF.aws_region` (`backend-secret.yaml:68`), which dies with the CTF DROP. Do **not** ship
> "aws.region is configurable" guidance.

### (B) REMOVE / HARDCODE — internal wiring this trunk chart DOES consume

| Path | Why | How |
|---|---|---|
| `enbuildBk.grpcPort` | Service `port` must match hardcoded `targetPort`/containerPort `8443`; mismatch silently breaks routing | Template constant 8443 (`enbuild-bk-grpc-service.yaml:23-24`, `enbuild-bk.yaml:243`) |
| `enbuildBk.grpcServiceType`, `enbuildBk.service_type` | Internal gRPC/HTTP svc; NodePort/LB bypasses mesh/SSO | Fixed `ClusterIP` (`enbuild-bk-grpc-service.yaml:16`, `enbuild-bk.yaml:282`) |
| `enbuildBk.serviceAccount` | Installer SA name is release-derived | `_helpers`/`default (printf "%s-enbuild-bk-installer" ...)`; keep override only as escape hatch (`enbuild-bk.yaml:33`) |
| `enbuildBk.healthProbe.enabled` | Liveness/readiness must always render in prod | Constant `true` (`enbuild-bk.yaml:228,247`) |
| `enbuildBk.snapshotIngestEnabled` / `.snapshotReadFromStore` | Console snapshot pipeline = fixed prod default `true` | Bake `true` default (`enbuild-bk.yaml:93,95`) |
| `enbuildBk.kubeProxyCache.{schemaTtlSeconds,listTtlSeconds,maxEntries}` | Internal cache tuning, one correct value | Fixed defaults in `_helpers` (`enbuild-bk.yaml:104-114`) |
| `enbuildBk.installAgent.agentImageTag` | Fleet-moves-together (ADR-0034); chart-pinned | Fixed default; pin to a CalVer tag, not a SHA/feature ref (`enbuild-bk.yaml:136`) |
| `enbuildBk.authMechanism` | **No trunk consumer**; duplicate of `global.auth_plugin` | DROP / consolidate (also bucket C) |
| `enbuildBk.mtls.caSecretRef` | **No trunk consumer**; no analog feature | DROP (also bucket C) |
| `enbuildBk.exportSigning.privateKeyPem` | Inline-PEM signing key lane | Gate to dev/test; prod uses `existingSecret` (`backend-secret.yaml:78`) |
| `enbuildUi.service_type` / `enbuildUser.service_type` / `enbuildAI.service_type` / `enbuildCTF.service_type` | Internal proxied svcs reached via VS/Ingress/nginx | Fixed `ClusterIP` |
| `enbuildUi.node_port` | Inert once `service_type` locked to ClusterIP (read only under `eq service_type "NodePort"`) | DROP |
| `enbuildUi.repository` / `enbuildUser.repository` / `enbuildConsumer.repository` | Image path fixed; air-gap via registry | Template constant |
| `enbuildUi.resources` / `enbuildUser.resources` | **Hardcoded `resources: {}` in templates — INERT, not knobs** | Wire them first OR keep fixed; do NOT present as knobs (`enbuild-ui.yaml:44`, `enbuild-user.yaml:62`) |
| `enbuildUi.loki_url` / `grafana_url` | Relative reverse-proxy paths + fixed dashboard UIDs; reconcile the `loki_url`-vs-`grafana_url` env/configmap divergence (`enbuild-ui.yaml:35` vs `frontend-config.yaml:30`) | `_helpers` constants |
| `enbuildConsumer.command/args/livenessProbe` | Iron-Bank node-binary quirk; mirror the existing livenessProbe conditional | Template off image variant (subsumes live override) |
| `enbuildAI.ollama_endpoint`, `enbuildAI.serviceAccount.name` | Internal wiring / orphan (never read) | Constant / remove |
| `mongodb.mongo_root_username` | Credential material — should be in the Secret, not a values literal | Source from `mongodb.existingSecret` (`mongodb-secrets.yaml:9`) |
| `mongodb.type` | **CONSUMED at `mongodb.yaml:85** (`default "ClusterIP"`) — NOT an orphan**; internal datastore must not be LB/NodePort | HARDCODE fixed `ClusterIP` |
| `mongodb` StatefulSet topology (replicas/UID 998/GID 996/port 27017/`storage: 10Gi`) | Hand-rolled single-node; **`mongodb.persistence.size` does NOT exist (10Gi hardcoded `mongodb.yaml:75`)**; HA goes through enbuild-mongodb chart (ADR-0030) | Fixed |
| `global.auth_plugin`, `global.create_istio_virtual_service` | Undeclared orphans read by templates | Declare with fixed default |
| `global.ingress.classname` vs `className` | Latent bug: `ingress-frontend.yaml` reads BOTH; values declares only `classname` | Standardize on `className`, fix templates |
| `rabbitmq.fullnameOverride`/`service.type`/`service.ports`/`host`(:5672)/VS port 15672/`rbac.create`/clustering/`configuration` body | Parent consumes only 4 leaves; rest internal | Pin in chart; expose nothing |
| `rabbitmq.auth.securePassword: false` | Documented T1.2 correctness fix | Fixed |
| `headlamp.nameOverride`/`fullnameOverride`/`namespaceOverride`/`baseURL`/`inCluster`/`service.type`/`extraArgs`/`volumes`/`volumeMounts`/`clusterRoleBinding.create` | nginx proxies `{release}-headlamp`; static CM names can't template `.Release.Name` | Fixed |
| `headlamp.pluginsDir` | Live blanks it while volumeMounts mount theme there — live bug | Pin `/headlamp/plugins` |

### (C) OPERATOR-OWNED / OFF-TRUNK — NOT this chart's surface today

These have **zero consumers in this trunk chart** (verified: `grep -rln` over `templates/` +
`_helpers.tpl` returns nothing). They are merged from the off-trunk chart-line and out-of-band
operators (Keycloak operator → Secret `enbuild-ib-keycloak-secrets`; off-trunk PKI chart). Do
**not** tell anyone to `_helpers`-constant, allow-list, or "lock as required" a value this chart
never reads — a fresh `helm install` of THIS chart would demand inputs with no functional effect.

| Path | Status in trunk | Action |
|---|---|---|
| top-level `keycloak.*` (`enabled`, …) | Zero consumers — pure no-op orphan | **DROP from values** (gates nothing). Do not present as a knob. |
| `enbuildBk.keycloak.connectionUrl`/`.realm`/`.clientId`/`.clientSecret`/`.newUserPassword`/`.realmPublicKey`/`.adminBaseUrl`/`.backendUrl` | Zero consumers; Keycloak wired out-of-band | **PORT** the `KEYCLOAK_*` env wiring into `backend-secret.yaml`/`enbuild-bk.yaml` FIRST; then `connectionUrl`/`realm`/`clientId` become bucket-A knobs and `clientSecret`/`newUserPassword` move to `keycloak.existingSecret`. Until ported, leave live values untouched — they are silently ignored, not load-bearing. |
| `enbuildBk.authMechanism` | Zero consumers; real switch is `global.auth_plugin` (`backend-secret.yaml:10`) | DROP / consolidate onto `global.auth_plugin` when porting SSO |
| `enbuildBk.adminEmails` | Zero consumers | **PORT** the `ADMIN_EMAILS` env wiring first; then keep configurable (bucket A) |
| `enbuildBk.mtls.caSecretRef` | Zero consumers; no analog feature | DROP |
| `enbuildBk.gitlabConnection.host` | Zero consumers (the connection endpoint isn't wired on trunk) | **PORT** the consumer first; then bucket-A knob |
| top-level `pki.*` (`enabled`, `issuerName`, `caBundleConfigMap`, `caCert.secretName`, `hubServerCert.secretName`/`dnsNames`, durations, `istioAuthorizationPolicy.enabled`, `recreateHubIssuer`) | Zero consumers — no `kind: Certificate/Issuer/ClusterIssuer/AuthorizationPolicy`, no `cert-manager.io`, no `.Values.pki` anywhere in templates (only a stale comment in `nginx-conf.yaml`) | **KEEP THE FEATURE — PORT THE TEMPLATES.** mTLS PKI for the hub gRPC / pull-agent fabric is core. Port Certificate/Issuer/CA-bundle/AuthorizationPolicy gated on `pki.enabled`, add a `pki:` block to `values.yaml`. **Do NOT delete `pki.*` from the live override until ported** — it is the only thing pinning the live cert names. Note the live `enbuild-hub-issuer` ClusterIssuer is currently MISSING on-cluster (per the divergence record) — reconcile before enabling. After porting, only `pki.hubServerCert.dnsNames` is a customer knob (bucket A); the rest are HARDCODE. |

---

## 4. SIMPLIFY — keep the knob, fix the default

| Path | Problem today | Fix |
|---|---|---|
| `enbuildBk.clusterRpcTimeoutMs` | **Live override `10000` REGRESSES the chart default `"30000"`** (`values.yaml:227`) — `10000` re-triggers the "Lost connection to the cluster" bug the default was raised to fix. Consumed `enbuild-bk.yaml:77`. | Remove the live override; keep `30000`. Consider hardcoding the default and dropping the knob entirely. **Correctness flag.** |
| `enbuildBk.encryption_key` | Published literal `encryption_key` (default + live) → `ENCRYPTION_KEY` (`backend-secret.yaml:18`) | Required `enbuildBk.encryptionKey.existingSecret`; **no shipped default**; migration-gated (rotation breaks decryption — see §9 Phase 2) |
| `mongodb.mongo_root_password` / `mongodb.mongo_server` | Default `mongo_root_password`/live `SuperSecret`; `mongo_server` carries the dangerous **literal cosmosDB help-text sentence** (`mongodb-secrets.yaml:15`) as a value — a stomp footgun | Required `mongodb.existingSecret`; env-substitute the password in `mongo_endpoint_override` via `$(MONGO_INITDB_ROOT_PASSWORD)`; remove the help-text default — derive or require explicit |
| `rabbitmq.auth.password` / `erlangCookie` | Literals `SuperSecret` / `lamba` | `auth.existingPasswordSecret` / `auth.existingErlangSecret` |
| `rabbitmq.image.tag` | **Chart default `3.11.13-debian-11-r0` is STALE** vs live `3.12.14` from `registry1.dso.mil/ironbank/bitnami/rabbitmq` | Bump chart default to the Iron Bank tag actually running; keep configurable |
| `enbuildBk.keycloak.clientSecret` / `newUserPassword` | Live literals `enbuild-prototype-secret-CHANGE-ME` / `ChangeMe123!` (orphan today, but credentials) | Source from `keycloak.existingSecret` **when SSO wiring is ported** (bucket C); never inline |
| `enbuildAI.api_key` | Literal `dummy` → `OPENAI_API_KEY` (`backend-secret.yaml:20`, gated on `ai_lightning`) | Required existingSecret if AI retained; else dropped with feature |
| `global.image.pullPolicy` | Keep configurable (air-gap wants `IfNotPresent`); challenge ruled HARDCODE over-locked | Default `Always`, allow override |
| `enbuildBk.resources` | No default → empty resources on bare install (`enbuild-bk.yaml:258`) | Ship sane request/limit defaults |
| `loki-stack.grafana.adminPassword` | Live literal `V1v$oftR0ck5` | existingSecret (or removed with loki-stack drop) |
| `headlamp.clusterRoleBinding.clusterRoleName` | Default `cluster-admin` — hub privilege-escalation surface | Pin to a read-scoped `view`-style ClusterRole |
| `headlamp.persistence.enabled` | Default `true`, live `false`; Headlamp stateless here | Default `false` |
| `jupyterhub...authenticator_class` | Default `dummy` = auth bypass | If kept: `GenericOAuthenticator` → Keycloak; else dropped |

> `global.AppVersion` is **NOT dead.** Image tags resolve as
> `default .Chart.AppVersion (default .Values.global.AppVersion .Values.<svc>.image.tag)` (e.g.
> `enbuild-mq.yaml:44` and all 6 peers). `global.AppVersion` is the **middle fallback** when a
> per-service tag is unset; live sets it to `1.0.30`. Every live service currently pins an explicit
> tag so the fallback isn't hit today, but removing the key silently changes behavior for any
> service whose tag is later unset. **Keep it, or remove only with a deliberate migration note** —
> do not call it dead (the first pass double-listed it as both hardcode and drop, both wrong).

---

## 5. DROP — features / orphans to cut, with Chart.yaml cleanup

| Item | Verdict | Justification |
|---|---|---|
| **CTF** (`enbuildCTF.*`, `secure_lightning.ctf`, ctf blocks in `backend-secret.yaml:62-69` + `nginx-conf.yaml`) | DROP | Off in default+live; carries `debug:true`, dev CORS, `LOG_LEVEL=DEBUG`, and **the only `AWS_REGION`** into the shared backend Secret |
| **Bolt** (`enbuildBolt.*`, `develop_lightning.application`, `/bolt/` nginx + FE iframe) | DROP | `NODE_ENV=development`, public allowed-host; off everywhere |
| **JupyterHub** (`jupyterhub` dep + block + `develop_lightning.models`) | DROP from core | Heavyweight; embeds CDAO profiles w/ public github postStart clones + quay.io images; `authenticator_class: dummy`. Ship sanitized add-on if needed |
| **open-webui** (dep + `open-webui.*`) | DROP | Non-Iron-Bank repo (`helm.openwebui.com`); off everywhere; keep only `enbuildAI.ollama_endpoint` as managed-LLM bind point |
| **ollama / enbuildAI feature** | DROP (or hard off-by-default) | Dormant; live carries inert open-webui/ollama config (models/GPU/PV) as dead weight. If retained: keep `image.registry/tag` configurable + `api_key` as existingSecret |
| **loki-stack** (dep + `loki-stack.*`, `operations_lightning.monitoring`) | DROP | Loki+Promtail+Prometheus+Grafana duplicates Big Bang monitoring in P1; off everywhere; console uses `enbuildUi.grafana_url` iframes; plaintext Grafana pw inside |
| top-level `keycloak` | DROP | No-op orphan, zero trunk consumers (bucket C) |
| `enbuildConsumer.gitlab.*` / `github.*` (19 SCM keys) | DROP | Orphan-by-default; duplicate the connection path; emit plaintext PATs into the Secret if set (`backend-secret.yaml:33-57`); github has no existingSecret guard |
| `enbuildBk.mtls.caSecretRef` | DROP | No consumer, no analog (bucket C) |
| `enbuildAI.serviceAccount.name` / `enbuildCTF.serviceAccount.name` | DROP | Genuine orphans — SA name is always `{release}-...` |
| `enbuildUi.node_port` | DROP | Inert once `service_type` locked |
| `lightning_features.deploy_lightning.data_lightning` / `.infra_lightning` | DROP | `data_lightning` cosmetic-only FE flag; `infra_lightning` is the only flag defaulted `true` yet gates nothing (FE caps hardcoded true) |
| **Headlamp** | **KEEP** | Only optional subchart ON in live; fence passthrough, pin `clusterRoleName` to read-scope, `pluginsDir`, `persistence.enabled:false`; expose only image/spokes/replicaCount/resources |
| **pki** | **KEEP — port templates** (bucket C) | Core hub mTLS; do not delete from live until ported |

**Chart.yaml dependency cleanup:** remove `jupyterhub`, `open-webui`, and `loki-stack` from
`dependencies:` so their entire upstream value trees vanish from the surface. Keep `rabbitmq` and
`headlamp`. (Do not touch `mongodb.type`/`global.AppVersion` as "orphans to drop" — `mongodb.type`
is consumed and `global.AppVersion` is a live fallback; see corrections above.)

---

## 6. Subchart passthrough risk + fencing

The dependency subcharts (`rabbitmq` 11.13.0; the to-be-dropped `jupyterhub` 4.3.1, `open-webui`
3.1.6, `loki-stack` 2.10.2; `headlamp` 0.36.0) each expose their *entire* upstream value tree. A
customer can override **any** upstream key the hub never surfaces — rabbitmq clustering/TLS/LDAP,
jupyterhub `hub.cookieSecret`/storage/securityContext, headlamp OIDC — and silently destabilize
the platform. The live file proves it: RabbitMQ passthrough alone (full `rabbitmq.conf` body, LDAP
auth_backend, TLS subtree, metrics) is a large fraction of the 943 lines. Worse, `loki-stack` and
the off-trunk `pki` carry `{{ .Release.Name }}` / `{{ .Values.global.domain }}` strings *inside
subchart values* that Helm will not template (latent bugs masked only because those features are
OFF).

**Fence it:**
1. **Pin upstream internals in the chart's own `values.yaml`** under each subchart key
   (`service.type ClusterIP`, fullnameOverride family, ports, rbac, clustering profile,
   configuration body) so they are chart constants.
2. **Expose nothing from subcharts beyond the allow-list** — for survivors (rabbitmq, headlamp):
   only `image.registry/repository/tag`, `imagePullSecrets`, `persistence.size/storageClass`,
   `resources`, `replicaCount`, `headlamp.spokes`, and the existingSecret refs.
3. **Delete the dropped deps** (jupyterhub, open-webui, loki-stack) from `Chart.yaml`.
4. Document that unknown subchart overrides are unsupported; consider a CI lint that rejects
   operator values outside the allow-list.

---

## 7. Security red flags — plaintext secret defaults

| Path | Live committed value | Risk | Fix |
|---|---|---|---|
| `enbuildBk.encryption_key` | `encryption_key` (default + live identical) | Published at-rest crypto key in P1 | Required existingSecret; **migration** (rotation breaks decrypt) |
| `global.image.registry_credentials.username/password` | `cms553` / `VxoMdavv3DJs0ac3Y6JjXy69KVJh8gkj` | Leaked Iron Bank registry creds, b64'd into dockerconfigjson (`_helpers.tpl:89-99`) | Pre-created docker-registry Secret; **ROTATE cms553** |
| `global.gitlabRegistryCredentials.username/password` | `cmsanchezm` / `glpat-6PcSe4LCffqlb_fpk-...` | Leaked GitLab PAT (merged into same dockerconfigjson, `_helpers.tpl:91-96`) | Pre-created Secret; **ROTATE PAT immediately** |
| `rabbitmq.global.image.registry_credentials` (live ~519, ~694) | `cms553` / `VxoMdavv...` (duplicated in two subchart global blocks) | Leaked creds, duplicated | Remove per-subchart global blocks (helm merges parent global); **ROTATE cms553** |
| `mongodb.mongo_root_password` | `SuperSecret` (+ embedded in `mongo_endpoint_override` URI) | Committed DB credential | `mongodb.existingSecret`; env-substitute in URI |
| `mongodb.mongo_root_username` | `enbuild` | Credential material in values | Move into `mongodb.existingSecret` |
| `rabbitmq.auth.password` | `SuperSecret` | Committed broker credential | `auth.existingPasswordSecret` |
| `rabbitmq.auth.erlangCookie` | `lamba` | Clustering shared secret | `auth.existingErlangSecret` |
| `loki-stack.grafana.adminPassword` | `V1v$oftR0ck5` | Committed Grafana admin pw | existingSecret (or dropped with loki-stack) |
| `enbuildBk.keycloak.clientSecret` | `enbuild-prototype-secret-CHANGE-ME` | OIDC client secret literal (orphan today) | existingSecret only, when SSO ported |
| `enbuildBk.keycloak.newUserPassword` | `ChangeMe123!` | Bootstrap user pw literal (orphan today) | existingSecret only, when SSO ported |
| `enbuildAI.api_key` | `dummy` | Plaintext LLM key default → `OPENAI_API_KEY` | existingSecret if AI retained |
| `enbuildConsumer.github.token`/`manifest_token`/`gitlab.token` | (unset, but template emits plaintext if set; github paths have NO existingSecret guard) | Latent plaintext PAT lanes (`backend-secret.yaml:35,47,53`) | Drop blocks; existingSecret if reintroduced |
| `jupyterhub...authenticator_class` | `dummy` | Auth bypass (any user, no password) | GenericOAuthenticator → Keycloak, or drop |
| `headlamp.clusterRoleBinding.clusterRoleName` | `cluster-admin` | Full read/WRITE on hub cluster | Pin to read-scoped ClusterRole |
| `enbuildBk.exportSigning.privateKeyPem` | (empty default — good) | Inline-PEM audit-signing key lane exists | Gate inline lane to dev/test; prod uses existingSecret |

**ROTATION call-to-action (do first, out-of-band):** rotate the `cms553` Iron Bank registry
credential AND the `glpat-6PcSe...` GitLab PAT immediately — both are committed live and the
`cms553` key is duplicated across the parent and rabbitmq subchart global blocks.

---

## 8. Proposed LEAN `values.yaml` skeleton (public surface only — CORRECT key names)

```yaml
## ===== ENBUILD Hub — public surface. Everything else is locked in templates. =====
global:
  domain: ""                             # external DNS suffix (REQUIRED bind)
  disable_tls_gitlab: false
  istio:
    enabled: true
    gateway: istio-system/main
  ingress:                               # alternative to istio
    enabled: false
    tls: false
    tls_secret: ""
    className: ""                        # standardized (was classname bug)
    annotations: {}
  image:
    registry: registry1.dso.mil          # air-gap / Iron Bank rebind
    pullPolicy: Always                   # air-gap may set IfNotPresent
  imagePullSecretName: ""                # pre-created docker-registry Secret (REQUIRED for private/IB)
  storageClass: ""

enbuildBk:
  image: { tag: "" }                     # unset -> global.AppVersion -> Chart.AppVersion
  replicas: 1                            # NOTE: key is `replicas`, not replicaCount
  resources: {}                          # templates ship sane defaults; override to tune
  grpcEnabled: true
  grpcVirtualService:                    # agent gRPC ingress — the load-bearing per-env knobs
    enabled: true
    host: ""                             # REQUIRED — external gRPC FQDN for this hub
    additionalHosts: []
    gateway: istio-gateway/public-ingressgateway
  installAgent:
    hubUrl: ""                           # REQUIRED — host:443 the agent dials back to
    tlsServerName: ""
    hubInternalIps: ""
    existingSecret: ""                   # pre-created install-agent creds Secret
  kubeProxyFallbackActor: ""             # audit actor for proxy fallback
  encryptionKey:
    existingSecret: ""                   # REQUIRED — Secret holding ENCRYPTION_KEY (was published literal)
  securityTooling:
    existingSecret: ""
  exportSigning:
    existingSecret: ""                   # prod signing-key Secret (inline PEM is dev/test only)
  gitlabPat:
    existingSecret: ""                   # REQUIRED — Secret holding the GitLab PAT
  # --- requires porting templates first (bucket C); not wired on trunk yet ---
  gitlabConnection: { host: "" }
  keycloak: { connectionUrl: "", realm: "", clientId: "", existingSecret: "" }
  adminEmails: []

enbuildUi:
  hostname: ""                           # console external FQDN (REQUIRED)
  image: { tag: "" }
  replicas: 1
  kiali_url: ""                          # console iframe URL (env-specific)
  kubecost_url: ""                       # console iframe URL (env-specific)
  # resources: hardcoded {} in template — not a knob

enbuildConsumer:
  image: { registry: "", tag: "" }       # per-service air-gap registry override + tag
  replicas: 1
  resources: {}

enbuildUser:
  image: { tag: "" }
  replicas: 1
  # resources: hardcoded {} in template — not a knob

mongodb:
  enabled: false                         # the ONE topology switch (live runs managed/external)
  existingSecret: ""                     # REQUIRED when self-hosted — Mongo root credential
  mongo_endpoint_override: ""            # external/managed Mongo (pw via $(...) env)
  storageClassName: ""                   # falls back to global.storageClass
  image: { repository: "", tag: "" }
  # persistence.size does NOT exist — StatefulSet hardcodes 10Gi

rabbitmq:
  enabled: true
  image: { registry: "", repository: "", tag: "" }  # default tag bumped to live IB 3.12.14
  replicaCount: 1                        # genuine bitnami passthrough key
  resources: {}
  auth:
    username: enbuild
    existingPasswordSecret: ""           # REQUIRED
    existingErlangSecret: ""             # REQUIRED
  host: ""
  env: ""
  queue_prefix: ""
  persistence: { size: 8Gi }            # REAL (mongo has no size knob)

headlamp:                                 # only retained optional ops subchart
  image: { registry: "", repository: "", tag: "" }
  replicaCount: 1                        # genuine subchart passthrough key
  resources: {}
  spokes: []

pki:                                      # hub mTLS — TEMPLATES MUST BE PORTED FIRST (bucket C)
  enabled: true
  hubServerCert: { dnsNames: [] }        # external SAN(s) for hub gRPC cert

# DROPPED entirely: enbuildAI/ollama, open-webui, loki-stack, jupyterhub,
# enbuildBolt, enbuildCTF, top-level keycloak, enbuildConsumer.gitlab/github,
# enbuildBk.authMechanism, enbuildBk.mtls. No aws.region exists anywhere.
# All internal wiring (ports, service_type, internal URLs, auth_plugin,
# install-agent image pin, snapshot/cache/health tuning, subchart internals)
# lives in templates / _helpers.tpl.
```

---

## 9. Rollout / migration notes (LIVE `enbuild-ib` must not break)

The live release runs an off-trunk chart line, so move in phased, reversible steps.

0. **Phase 0 — rotate leaked creds now (out-of-band, no chart change):** rotate `cms553`
   registry creds and the `glpat-6PcSe...` PAT; create the pre-created Secrets (`enbuild-pull`,
   `enbuild-gitlab-pat`, `enbuild-encryption-key`, `enbuild-mongo`, `enbuild-rabbit`,
   `enbuild-install-agent-creds` already exists) with the **current working values** so each later
   switch to existingSecret is value-preserving.

1. **Phase 1 — drop the dead, zero-risk first.** Remove CTF, Bolt, JupyterHub, open-webui,
   loki-stack, top-level `keycloak`, `enbuildConsumer.gitlab/github`, `enbuildBk.authMechanism`,
   `enbuildBk.mtls`, and the orphan leaves (`node_port`, `data_lightning`, `infra_lightning`).
   All are OFF/never-render in live → bump `Chart.yaml` `version:`; no-op to running pods.
   Removing the dropped subchart deps also deletes their passthrough surface. **Leave
   `mongodb.type` and `global.AppVersion` in place** — `mongodb.type` is consumed and
   `global.AppVersion` is a live image-tag fallback.

2. **Phase 2 — `encryption_key` migration (highest risk).** This key decrypts stored data; you
   cannot rotate it casually. Create `enbuild-encryption-key` holding the **current** literal
   `encryption_key`, switch the chart to source `ENCRYPTION_KEY` from that existingSecret (value
   unchanged), verify backend decrypt works, *then* plan a true rotation as a separate
   data-migration. Ship with a migration note; never change the value in the same step as the
   source change.

3. **Phase 3 — fix the `clusterRpcTimeoutMs` regression.** Remove the live `10000` override so the
   chart `30000` default applies (restores the "Lost connection to the cluster" fix). Verify agent
   RPC stability on vendor13-ib.

4. **Phase 4 — flip remaining secrets to existingSecret refs** (Mongo password/username, RabbitMQ
   password/erlangCookie, registry creds; Keycloak client secret only after SSO porting). Because
   Phase 0 pre-created them with live values, each flip is value-preserving; remove literals after
   each flip verifies.

5. **Phase 5 — hardcode internal wiring into templates/`_helpers`.** ports, `service_type`
   ClusterIP, `grpcPort` 8443, internal URLs, `global.auth_plugin`/`create_istio_virtual_service`
   declared with fixed defaults, install-agent image pin (pin `agentImageTag` to a CalVer tag, not
   a SHA), snapshot/cache/health tuning, consumer command/args off image-variant. Verify the gRPC
   port stays 8443, mq-consumer still selects `run:mq:all`, and the
   `grpcVirtualService.host`/`gateway` knobs still route the agent.

6. **Phase 6 — port the off-trunk SSO + PKI (bucket C).** Add the `KEYCLOAK_*` / `ADMIN_EMAILS`
   env wiring to `backend-secret.yaml`/`enbuild-bk.yaml`, and port the Certificate / Issuer /
   CA-bundle / AuthorizationPolicy templates gated on `pki.enabled` with a `pki:` block in
   `values.yaml`. KEEP `connectionUrl`/`realm`/`clientId`/`adminEmails`/`dnsNames` configurable;
   secret-ref `clientSecret`/`newUserPassword`. Carry a migration note for the live cert names
   (`enbuild-hub-grpc-tls`, `enbuild-hub-issuer`, etc.) so the running gRPC fabric is undisturbed,
   and reconcile the **missing** `enbuild-hub-issuer` ClusterIssuer before enabling.
   **Until this phase lands, leave the live `pki.*` and `enbuildBk.keycloak.*` override values in
   place** — they are silently ignored by the chart but are the only thing pinning the
   operator-managed cert/SSO names.

7. **Key-name trap (critical).** When standardizing capacity, the core services use
   **`replicas`** (`enbuild-bk.yaml:10`, `enbuild-ui.yaml:12`, `enbuild-mq.yaml:10`,
   `enbuild-user.yaml:10`). Do **not** rename them to `replicaCount` in values without also
   patching the templates — a values-only rename silently forces every core service to 1 replica
   (Helm default lookup miss) on the live deploy. Only `rabbitmq.replicaCount` /
   `headlamp.replicaCount` are genuinely `replicaCount`.

8. **Headlamp gate caution.** If the `lightning_features` flattening renames
   `operations_lightning.headlamp`, carry a migration alias — it is the one optional flag ON in the
   live release.

9. **Per CLAUDE.md auto-publish:** ship each phase as edit chart → bump `version:` in `Chart.yaml`
   → commit/push; CI publishes the OCI package. Do not `helm push` manually. Validate each phase
   against vendor13-ib with `helm diff` before the operator rolls.
