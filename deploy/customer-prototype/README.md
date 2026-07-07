# ENBUILD hub — customer prototype deploy bundle

This folder is a self-contained bundle for deploying the ENBUILD hub in a
**customer environment** and demonstrating the in-console **Headlamp** cluster
browser working end-to-end (view a managed/imported spoke's Kubernetes resources
from the hub console).

| File | Purpose |
|------|---------|
| `values-prototype.yaml` | Recommended prototype deploy — real SSO login + real per-user Headlamp browsing, permissive tenancy. |
| `values-production.example.yaml` | Hardened target — same wiring with strict tenancy enforcement on. |
| `HEADLAMP-AUTH-HARDENING.md` | Full status: what works, the blockers (with fixes), and the production cutover plan. **Read this.** |

## Honest posture (what this demonstrates vs. what's deferred)

- **Demonstrated:** the hub deploys in your environment; operators log in via
  Keycloak SSO; they import/create clusters and **browse real Kubernetes
  resources on the spokes through Headlamp** — proving the full
  console → hub proxy → agent → spoke path.
- **Deferred to Phase 2 (if selected):** strict multi-tenant authorization
  enforcement, per-user spoke RBAC (impersonation), and a couple of small code
  hardening fixes. These are **designed and documented** in
  `HEADLAMP-AUTH-HARDENING.md` — not open questions.

The prototype runs with `consoleAuthStrict=false` (permissive tenancy). That is a
deliberate, disclosed choice for evaluation, not an oversight.

## Prerequisites

1. A Kubernetes cluster with the ENBUILD chart repo added (`helm repo add enbuild …`)
   or the chart available locally, chart version **`0.1.0-p1ccm-trunk.34`** or newer
   (the nginx redirect fix — blocker #1 — landed in `.34`).
2. DNS + TLS for two hostnames: the console (`enbuild-ib.<domain>`) and Keycloak
   (`keycloak.<domain>`). If you run Istio, the chart creates the VirtualServices;
   otherwise expose the console/Keycloak via your own ingress and set
   `global.istio.enabled=false`.
3. The bootstrap secrets (encryption key + Mongo). A helper is provided in the
   chart: `charts/enbuild/scripts/create-bootstrap-secrets.sh` (creates
   `enbuild-ib-encryption-key` and the Mongo secret in the `enbuild` namespace).

## Install

> **Release name / namespace are now flexible** (blocker #4 fixed in `.35`): the backend
> derives the Headlamp hostnames from `{{ .Release.Name }}` / `{{ .Release.Namespace }}`, so
> any name works. The examples below use `enbuild-ib` / `enbuild` for consistency with the
> demo realm and docs; if you change them, keep the console/Keycloak FQDNs in the values in
> sync.

```bash
# 1. Edit values-prototype.yaml — replace every "example.mil" with your domain.
#    Keep the three Keycloak URLs + the Headlamp issuerURL byte-identical.

# 2. Namespace + bootstrap secrets
kubectl create namespace enbuild
charts/enbuild/scripts/create-bootstrap-secrets.sh -n enbuild   # or your own operator Secrets

# 3. Install (any release name/namespace works; enbuild-ib/enbuild used for consistency)
helm install enbuild-ib enbuild/enbuild -n enbuild -f deploy/customer-prototype/values-prototype.yaml
```

The chart runs **fail-closed render guards** at install time — if the Keycloak
URLs or the Headlamp OIDC issuer are inconsistent, `helm install` fails with a
specific message instead of deploying a silently-broken auth chain. If it
installs, the auth wiring is internally consistent.

## Verify Headlamp works (the acceptance test)

1. Browse to `https://enbuild-ib.<domain>/p1-ccm-console/` → you are redirected to
   Keycloak, log in, and land **back on the console** (no dead-port redirect — this
   is the blocker-#1 fix in action). Demo users are in the bundled realm
   (`admin@p1.mil` etc.); rotate these before production.
2. Onboard a test cluster **both ways**:
   - **Import** (brownfield) via the Import wizard — wires Headlamp immediately; complete
     the project-tag step.
   - **Create** (greenfield) via the catalog — the cluster is swept into Headlamp
     automatically by the reconciler within ~5 minutes of the agent connecting (blocker #3
     fixed; no import needed). Give it a few minutes after the agent goes healthy.
3. Open the cluster's detail page → **"Open in Headlamp ↗"** → Headlamp completes
   its own OIDC login and **lists the spoke's namespaces/pods/etc.** That is the
   full path authenticated as the logged-in user.

If step 3 shows resources, connectivity + per-user cluster visibility are proven.

## Hardening to production

Switch to `values-production.example.yaml` (`consoleAuthStrict=true`) and follow
Section 4 of `HEADLAMP-AUTH-HARDENING.md`. The three small code fixes (create-path
auto-wire, host de-hardcode, public client-id) are itemized there with file:line
and are offered as a follow-on change set.
