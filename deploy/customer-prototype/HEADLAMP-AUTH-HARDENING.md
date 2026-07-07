# Headlamp cluster-view — authentication status & production hardening plan

**Audience:** ENBUILD delivery team + evaluating customer.
**Scope:** the in-console **Headlamp** cluster browser (view a managed/imported spoke's
Kubernetes resources from the ENBUILD hub console).
**Status date:** 2026-07-07 · trunk `ccm-p1` (backend) + `enbuild-helm@0.1.0-p1ccm-trunk.34` (chart).

---

## 1. Executive summary

**The Headlamp data path is proven end-to-end.** With a valid Keycloak user token, the
console → hub `K8sApiProxy` → agent gRPC → spoke kube-apiserver chain returns real
cluster data (verified live: an authenticated request listed **49 / 46 namespaces** on two
managed spokes). The agent's spoke RBAC, the tenancy checks, and the token verification are
all implemented and working.

**What is *not* turnkey yet is production authentication.** The chart ships with browser
OIDC for Headlamp and strict console auth **OFF by default**, so a fresh install runs in
**demonstration mode** (permissive) until an operator performs a small, well-defined
"cutover." This is a deliberate prototype posture, not a defect — but it must be stated
honestly and the production steps must be documented. That is what this file is.

**Recommended posture for the prototype delivery:** deploy in **permissive/demo mode**
(see `values-prototype.yaml`), which demonstrates real connectivity and cluster visibility
in the customer's environment, and treat the production auth cutover (Section 4) as the
Phase-2 / if-selected work.

---

## 2. How it works (so the gaps are legible)

```
 Console (React SPA)                    Hub (in-cluster)                     Spoke cluster
 ───────────────────                    ────────────────                     ─────────────
 "Open in Headlamp ↗"  ──new tab──▶  /headlamp/  (Headlamp pod)
   (plain <a> link;                     │  dynamic cluster = <slug>
    no token passed)                    │  context from ConfigMap
                                        │  headlamp-kubeconfigs
                                        ▼
                              Headlamp browser OIDC login  ◀── Keycloak (realm: enbuild,
                              (chart headlamp.config.oidc.*)     public client: enbuild-ui)
                                        │  attaches user bearer
                                        ▼
                     Hub  /api/v1/clusters/<id>/k8s-api/*   (KubeProxyController)
                       • RolesGuard: verify bearer (issuer EXACT-match, JWKS-first)
                       • strict gate: require keycloak-jwt actor  (CONSOLE_AUTH_STRICT)
                       • TenancyGuard: project membership from token 'groups' claim
                                        │  gRPC over agent's outbound mTLS stream
                                        ▼
                              enbuild-agent (on spoke)  ──▶  kube-apiserver
                       • executes as its OWN ServiceAccount (view ClusterRole)
                       • strips inbound Authorization; spoke never sees the user token
```

Key consequences:
- The **console does not bridge a token** to Headlamp — Headlamp runs **its own** OIDC
  login. So `headlamp.config.oidc.*` must be configured for real per-user auth.
- The **spoke identity is the agent's `view` SA**, not the user — so no per-user RBAC has to
  exist on the spoke for browsing to work (per-user spoke enforcement is optional Phase-2
  defense-in-depth).
- Reference: FE `frontend-p1ccm/.../ClusterDetail.tsx` (the `Open in Headlamp` anchor);
  BE `backend/microservices/enbuild/src/kubeProxy/kubeProxy.controller.ts` (proxy + guards);
  BE `.../installAgent/installAgent.service.ts:2224` (`reconcileHeadlamp` builds the
  `headlamp-kubeconfigs` ConfigMap); agent `agents/enbuild-agent/internal/clients/kube_proxy.go`.

---

## 3. Blockers — status, evidence, fix

| # | Blocker | Evidence | Fix | Status |
|---|---------|----------|-----|--------|
| **1** | **nginx 301'd no-slash paths to a dead `:8080`** — broke the SSO login round-trip *and* Headlamp's `/headlamp/oidc-callback`. Looked exactly like an auth wall. | `curl /p1-ccm-console` → `301 http://localhost:8080/p1-ccm-console/`; `/headlamp` → `:8080/headlamp/`. Cause: nginx default `absolute_redirect on` + in-pod `listen 8080`, so slash-append redirects use the internal port, not the browser origin. | Add **`absolute_redirect off;`** to `charts/enbuild/templates/nginx-conf.yaml` (relative Location keeps the browser on its real origin). | ✅ **FIXED** in `0.1.0-p1ccm-trunk.34` |
| **2** | **Headlamp browser OIDC + strict console auth ship OFF.** Fresh install = permissive fallback (no real per-user auth), or if strict is turned on *without* OIDC = hard 401. This is the "issues authenticating." | `values.yaml headlamp.config.oidc.{clientID,issuerURL,scopes,callbackURL}` all default `""`; `enbuildBk.consoleAuthStrict` default `false`. Live hub had neither set (tokenless K8s-api call → `401 "User roles not defined"`). | Set the **four `headlamp.config.oidc.*` values together** + `consoleAuthStrict=true` + real FQDNs (Section 4). The chart's render-time guards enforce consistency once OIDC is non-empty. | 📋 **DOCUMENTED** — prototype runs permissive; production = Section 4 |
| **3** | **Catalog-CREATE does not wire a new cluster into Headlamp.** `reconcileHeadlamp` is only called by the **import** path, so a created cluster is invisible in Headlamp until a later import (which rebuilds the ConfigMap from *all* managed clusters) or a manual `headlamp.spokes` edit. | `installAgent.service.ts:1096` is the only caller (`reconcileHeadlamp` defined at `:2224`); catalog-create registers the agent via heartbeat (`agent-registry.service.ts` sets `managedCluster=true`) but never calls reconcile. | **Workaround (prototype):** import your test clusters (import fully wires Headlamp), or run any import to sweep created ones in. **Durable fix (Phase-2):** trigger a *content-guarded* reconcile on agent register/first-heartbeat (only patch+roll Headlamp when the rendered ConfigMap actually changes, to avoid roll churn). | 📋 **DOCUMENTED** + workaround |
| **4** | **`reconcileHeadlamp` hardcodes the hub service host + Headlamp deployment name to `enbuild-ib` / `enbuild`.** A non-`enbuild-ib` release name (or non-`enbuild` namespace) makes every spoke context point at a non-existent Service → Headlamp spoke browsing 404s. | `installAgent.service.ts:2264` writes `server: http://enbuild-ib-enbuild-backend.enbuild.svc.cluster.local/...`; the Headlamp rollout targets `deployments/enbuild-ib-headlamp`. | **Prototype:** install the release as **`enbuild-ib`** in namespace **`enbuild`** (zero-code — see README). **Durable fix (Phase-2):** derive host + deployment name from the release name / namespace env instead of the literals. | 📋 **DOCUMENTED** — pin release name |
| **5** | **ConfigMap OIDC `client-id` sourced from the confidential client env** (`KEYCLOAK_CLIENT_ID`) instead of the public PKCE client (`enbuild-ui`). | `installAgent.service.ts:2249`: `oidcClientId = process.env.KEYCLOAK_CLIENT_ID || 'enbuild-ui'`. Tolerated by the verifier (`azp`/`aud` ∈ {enbuild, enbuild-ui}) — advisory field only. | **Durable fix (Phase-2):** read `KEYCLOAK_PUBLIC_CLIENT_ID || 'enbuild-ui'`. Low severity; not an auth blocker. | 📋 **DOCUMENTED** — low severity |
| **6** | **Under strict, a freshly imported cluster must be project-tagged** or verified non-admins 403 (K8sApiProxy is excluded from the unscoped-viewer bypass). | tenancy-guard; ImportWizard already has the tag step (`tagStackProject`). | Operator completes the import wizard's project-tag step (already in the UI). | 📋 Operational note |
| **7** | **Air-gapped / IL5:** the import Job and agent pod pull the enbuild-agent chart + image from `registry.gitlab.com`. | `installAgent.service.ts` agent chart/image env; `agents/enbuild-agent/chart/values.yaml` image repo. | Mirror the enbuild-agent OCI chart + image and the enbuild-stack bootstrap chart to the customer registry; override `AGENT_OCI_REF` / `image.repository` / `pullSecrets`. | 📋 Air-gap note |

---

## 4. Production auth cutover (Phase-2 / when hardened auth is required)

Perform these **as one atomic set** — the render guards (`headlamp-oidc-guard.yaml`,
`keycloak-url-guard.yaml`, `backend-secret.yaml`) fail the install on a *partial* or
*mismatched* set, but they cannot catch an all-empty OIDC block (that just silently leaves
Headlamp SSO off). See `values-production.example.yaml`.

1. `enbuildBk.consoleAuthStrict: true` — enforce verified tokens (no permissive fallback).
2. `enbuildBk.keycloak.url` **and** `enbuildUi.keycloak.url` **and** `keycloak.hostname` →
   the **same** real browser-facing Keycloak origin (`https://keycloak.<customer>` — never
   `localhost`). Guard: `keycloak-url-guard.yaml`.
3. `enbuildBk.keycloak.backendUrl` → in-cluster Keycloak URL (JWKS fetch source).
4. `enbuildUi.keycloak.clientId: enbuild-ui` (public PKCE client).
5. **All four** `headlamp.config.oidc.*`:
   - `clientID: enbuild-ui`
   - `issuerURL:` **byte-exact** to the backend `KEYCLOAK_ISSUER` (scheme/host/port/no
     trailing slash must all match — a trailing-slash difference is a silent 401)
   - `scopes: "openid email groups"` — the **`groups`** scope is required or every
     project-scoped browse 403s at the TenancyGuard
   - `callbackURL: https://<console-host>/headlamp/oidc-callback`
6. Keycloak realm `enbuild-ui` client must trust the console origin **and**
   `https://<console-host>/headlamp/oidc-callback`. The bundled demo realm uses a wildcard
   (`<console-url>/*`) that covers it once `keycloak.demoRealm.consoleUrl` (or
   `enbuildUi.hostname` + `global.domain`) is set to the real host; a bring-your-own realm
   with exact redirect URIs must register the callback explicitly.
7. Every user has ≥1 realm role and the correct `groups` (project membership +
   platform-admin group), or verified tokens still 401 (no role) / 403 (no group).

### Phase-2 code changes (small, precise — offered separately)

- **CREATE auto-wire (blocker #3):** add a content-guarded `reconcileHeadlamp` trigger on
  agent register / first heartbeat in `agent-registry.service.ts` (only patch + roll when the
  rendered `headlamp-kubeconfigs` changes).
- **De-hardcode (blocker #4):** in `installAgent.service.ts` derive the backend Service host
  (`:2264`) and Headlamp deployment name (the rollout target) from the release name /
  `ENBUILD_NS` env, with the current literals as defaults; wire the env in the chart from
  `{{ .Release.Name }}` / `{{ .Release.Namespace }}`.
- **Public client-id (blocker #5):** `installAgent.service.ts:2249` →
  `process.env.KEYCLOAK_PUBLIC_CLIENT_ID || 'enbuild-ui'`.

---

## 5. What already works (do not re-do)

- Console → K8sApiProxy → agent → spoke data path (proven: real namespaces returned).
- BE token verification: issuer exact-match pin, JWKS-first with pinned-key fallback,
  `azp`/`aud` accepted for both the confidential and public clients.
- Agent spoke RBAC: `rbac.headlampView.enabled=true` is the agent-chart default (binds the
  built-in `view` ClusterRole); the agent uses its own SA and strips inbound Authorization.
- `config.json` is ConfigMap-mounted (durably overridable via values — no image rebuild to
  retarget Keycloak).
- The **import** path fully and idempotently wires Headlamp and rolls the Headlamp pod.
- Chart fail-closed render guards catch most auth/URL misconfigurations at `helm install`.
- nginx dead-redirect (blocker #1) — **fixed** in `0.1.0-p1ccm-trunk.34`.
