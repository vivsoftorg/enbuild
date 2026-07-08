# ENBUILD hub — customer prototype deploy bundle

This folder is a self-contained bundle for deploying the ENBUILD hub in a
**customer environment** and demonstrating the in-console **Headlamp** cluster
browser working end-to-end (view a managed/imported spoke's Kubernetes resources
from the hub console).

Pick ONE values file for your posture. They form a ladder from "just works" to
"fully hardened":

| File | Posture | Use when |
|------|---------|----------|
| **`values-demo.yaml`** ⭐ | **Permissive (`authMechanism=local`).** Headlamp browses EVERY managed cluster out-of-the-box under an admin fallback — no per-user auth. **Most bulletproof / zero-config-drift.** | **Recommended quick-start.** You want a guaranteed-working demo in the customer environment with no OIDC to get wrong. Live-validated: Headlamp browses greenfield + imported spokes in the browser. |
| `values-prototype.yaml` | Real Keycloak SSO login + real per-user Headlamp OIDC, permissive tenancy. | You have a real, browser-AND-pod-reachable Keycloak FQDN and want per-user auth (not just permissive). |
| `values-production.example.yaml` | The above + strict per-user tenancy enforcement. | Production hardening. |
| **`values-images.yaml`** ⚠️ | **Pinned image tags. REQUIRED with every install** — the chart defaults tags to a non-existent appVersion, so without this a fresh install ImagePullBackOffs. BE tag includes the Headlamp fixes. | **Always** (`-f values-images.yaml` in addition to your posture file). |
| `HEADLAMP-AUTH-HARDENING.md` | Full status, the blockers (with fixes), the cutover plan. **Read this.** | Always. |

## Honest posture / DISCLOSURE (state this to the customer)

The **recommended quick-start (`values-demo.yaml`) is intentionally permissive**:
`authMechanism=local` means the backend does **not** enforce per-user
authorization — the console's Headlamp K8s proxy serves reads under a single
admin fallback identity, so anyone who can reach the console can browse cluster
resources. **This is a deliberate, disclosed prototype choice** — it guarantees
Headlamp works everywhere with no environment-specific auth wiring to break.

The **production-secure path is already implemented and documented** (real
Keycloak SSO + Headlamp OIDC + strict per-user tenancy — `values-production.example.yaml`
+ HEADLAMP-AUTH-HARDENING.md §4). Nothing about hardening is an open question;
it is a values switch away when enforced authorization is required.

Why permissive is the *quick-start* and not OIDC: Headlamp's OIDC does a
SERVER-SIDE token exchange from inside the cluster, so it needs Keycloak
reachable from the browser **and** the in-cluster Headlamp pod (a real FQDN).
Until that exists, local mode is the no-surprises default.

## Prerequisites

1. A Kubernetes cluster with the ENBUILD chart repo added (`helm repo add enbuild …`)
   or the chart available locally, chart version **`0.1.0-p1ccm-trunk.35`** or newer
   (nginx redirect fix landed in `.34`; create-wire + release-name portability in `.35`).
2. DNS + TLS for two hostnames: the console (`enbuild-ib.<domain>`) and Keycloak
   (`keycloak.<domain>`). If you run Istio, the chart creates the VirtualServices;
   otherwise expose the console/Keycloak via your own ingress and set
   `global.istio.enabled=false`.
3. The bootstrap secrets (encryption key + Mongo). A helper is provided in the
   chart: `charts/enbuild/scripts/create-bootstrap-secrets.sh` (creates
   `enbuild-ib-encryption-key` and the Mongo secret in the `enbuild` namespace).
4. **Brownfield cluster import (Beat 5) — agent connect-back needs the hub's PKI, which
   the chart provisions automatically.** On chart **`.42`+** the default
   `enbuildBk.agentMtls.spokeCertEnroll: "true"` makes the hub mint a per-spoke mTLS
   client cert off the chart-provisioned `enbuild-hub-issuer` cert-manager ClusterIssuer
   and hand its **live CA** to the imported agent, so the spoke trusts the hub gRPC
   gateway's self-signed cert and connects back. **Nothing to set on `.42`+.** Requirements
   it depends on (all chart-provided): cert-manager installed + `pki.provisionHubCA` (creates
   `enbuild-hub-issuer` + `enbuild-ca-tls`). On an OLDER chart, set
   `enbuildBk.agentMtls.spokeCertEnroll: "true"` explicitly. **Do NOT rely on the legacy manual
   `HUB_CA_BUNDLE_PEM` AdminSettings field** — it goes stale per hub and was the imported-agent
   `context deadline exceeded` connect-back bug (failure-registry **G7**; launched spokes were
   always fine — this makes imported spokes fine too, with no operator step). Requires backend
   image ≥ `74ca99071` (pinned in `values-images.yaml`).

## Install

> **⚠️ Headlamp is ENABLED BY DEFAULT and is the whole point of this bundle.**
> The chart ships `lightning_features.operations_lightning.headlamp: true`. **Do NOT set it
> to `false` on a real deploy** — that turns off the cluster-browse feature (the Headlamp
> pod + its RBAC). Nothing in these values disables it; leave it on.

> **Release name / namespace are flexible** (blocker #4 fixed in `.35`): the backend derives
> the Headlamp hostnames from `{{ .Release.Name }}` / `{{ .Release.Namespace }}`, so any name
> works. Examples use `enbuild-ib` / `enbuild`; if you change them, keep the console/Keycloak
> FQDNs in the values in sync.

> **Running MORE THAN ONE hub in the SAME cluster?** Headlamp (and a few components) create
> **cluster-scoped** resources (e.g. a ClusterRoleBinding) named by the **release name** —
> so two hubs with the **same** release name in one cluster collide on those objects. The
> correct fix is a **DISTINCT release name per hub** (e.g. `enbuild-ib` and `enbuild-oobe`) —
> the backend reconcile is release-name-aware (Fix B), so any name works and Headlamp stays
> enabled on both. **Do NOT disable Headlamp to resolve a same-release-name collision.**
> (A customer running a single hub never hits this.)

```bash
# 1. Edit values-demo.yaml — replace every "example.mil" with your domain.
#    (Quick-start / recommended. For per-user auth use values-prototype.yaml instead.)

# 2. Namespace + bootstrap secrets
kubectl create namespace enbuild
charts/enbuild/scripts/create-bootstrap-secrets.sh -n enbuild   # or your own operator Secrets

# 3. Install (any release name/namespace works; enbuild-ib/enbuild used for consistency).
#    values-images.yaml is REQUIRED — it pins real image tags (chart default = non-existent appVersion).
helm install enbuild-ib enbuild/enbuild -n enbuild \
  -f deploy/customer-prototype/values-demo.yaml \
  -f deploy/customer-prototype/values-images.yaml
```

The chart runs **fail-closed render guards** at install time — if the Keycloak
URLs or the Headlamp OIDC issuer are inconsistent, `helm install` fails with a
specific message instead of deploying a silently-broken auth chain. If it
installs, the auth wiring is internally consistent.

## Verify Headlamp works (the acceptance test)

*(Steps below are for the recommended `values-demo.yaml` posture — LIVE-VALIDATED
2026-07-07 on the fresh eks-10 hub: Headlamp browsed a greenfield-created cluster's
46 namespaces and a second cluster's 125 workloads in the browser.)*

1. Browse to `https://enbuild-ib.<domain>/p1-ccm-console/` → you are redirected to
   Keycloak, log in, and land **back on the console** (no dead-port redirect — this
   is the nginx fix). Demo users are in the bundled realm (`admin@p1.mil` etc.);
   rotate these before production.
2. Onboard a test cluster **both ways**:
   - **Create** (greenfield) via the catalog — the cluster is swept into Headlamp
     automatically by the reconciler within ~5 min of the agent connecting (create-wire
     fix; no import needed).
   - **Import** (brownfield) via the Import wizard — wires Headlamp immediately.
3. Open the cluster's detail page → **"Open in Headlamp ↗"** → Headlamp **lists the
   spoke's namespaces / pods / workloads**. In `values-demo.yaml` (permissive) this
   works with no per-user token — the proxy serves reads under the admin fallback.
   (In `values-prototype.yaml`/production, Headlamp does its own OIDC login and the
   reads are per-user authorized.)

If step 3 shows resources, the console → hub proxy → agent → spoke path is proven.

## Hardening to production

Switch to `values-production.example.yaml` (`consoleAuthStrict=true`) and follow
Section 4 of `HEADLAMP-AUTH-HARDENING.md`. The three small code fixes (create-path
auto-wire, host de-hardcode, public client-id) are itemized there with file:line
and are offered as a follow-on change set.
