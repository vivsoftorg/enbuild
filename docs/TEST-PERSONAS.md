# CCM Test Personas — RBAC / Tenancy Demonstration

Non-secret evaluation artifact. These Keycloak users let an evaluator log in as each role and
**see what the roles do and how multi-tenancy is enforced** (test plan 3.2 SSO + the Management /
Security / Visibility objectives, and the 3.4 upgrade RBAC).

- Import file: `charts/enbuild/files/keycloak/ccm-test-personas.partial-import.json`
- Realm: `enbuild`
- **Shared test password:** `CcmRole!Demo2026`  ← documented eval credential; **change before any non-eval use.**

## Personas

| Username (email) | Keycloak group | Role | Use it to show |
|---|---|---|---|
| `ccm-admin@p1.mil` | `/enbuild/platform-admins` | platform-admin | sees/acts on **every** project (cross-project admin bypass) |
| `ccm-alpha-viewer@p1.mil` | `/enbuild/projects/alpha/viewer` | viewer | read-only on project **alpha** |
| `ccm-alpha-member@p1.mil` | `/enbuild/projects/alpha/member` | member | viewer + day-2 ops (cert renew, cred rotation, drift remediation) |
| `ccm-alpha-maintainer@p1.mil` | `/enbuild/projects/alpha/maintainer` | maintainer | member + **BB security patch + performance upgrade** |
| `ccm-alpha-owner@p1.mil` | `/enbuild/projects/alpha/owner` | owner | full control + **BB functional upgrade, K8s distro upgrade, capacity change** |
| `ccm-bravo-owner@p1.mil` | `/enbuild/projects/bravo/owner` | owner (project **bravo**) | **tenant isolation** — full power on bravo, **zero** access to alpha |
| `ccm-no-access@p1.mil` | _(none)_ | — | **denial** — authenticated but member of nothing → every project action denied |

## Capability matrix (grounded in `tenancy-guard.service.ts` `ACTION_MIN_ROLE`)

Roles are ranked `viewer(1) < member(2) < maintainer(3) < owner(4)`; an action is allowed when the
caller's role on the cluster's project is **>=** the action's minimum role. Unknown actions fail
closed (require owner).

| Capability | Min role | viewer | member | maintainer | owner | platform-admin |
|---|---|:--:|:--:|:--:|:--:|:--:|
| Read cluster info / metrics / cost / audit / network / security findings | viewer | ✅ | ✅ | ✅ | ✅ | ✅ |
| Cert renewal, credential rotation, drift remediation, NS-label reconcile, OIDC client register | member | ❌ | ✅ | ✅ | ✅ | ✅ |
| **BB security patch**, **BB performance upgrade** | maintainer | ❌ | ❌ | ✅ | ✅ | ✅ |
| **BB functional upgrade**, **K8s distro upgrade**, **capacity change**, **BB version rollout** | owner | ❌ | ❌ | ❌ | ✅ | ✅ |
| Any action on **a different project** | — | ❌ | ❌ | ❌ | ❌ | ✅ |

This is exactly the 3.4 Updates & Upgrades RBAC: a *member* cannot patch Big Bang, a *maintainer*
can apply security patches + performance upgrades, and only an *owner* (or platform-admin) can do a
functional/R-1 upgrade or a distro upgrade.

## How to apply

The personas live in the `enbuild` realm. After the hub's Keycloak is up:

```bash
# kcadm authenticated against the bundled Keycloak admin
kcadm.sh create partialImport -r enbuild \
  -s ifResourceExists=SKIP \
  -f charts/enbuild/files/keycloak/ccm-test-personas.partial-import.json
```

(or **Realm → Partial import** in the Keycloak admin console, selecting users + groups.)

## Prerequisite for end-to-end role behavior

Keycloak group membership is the tenancy **source of truth** (ADR-0022), but a project role only
takes effect once a **Project with the matching slug exists in CCM**:

1. Log in as `ccm-admin@p1.mil` and create projects **`alpha`** and **`bravo`** in the console.
2. Launch or import a cluster into each project.
3. Log in as each persona to observe the matrix above.

> Note: until the AUTH-1 fix lands (thread the verified `groups` claim into the cross-cluster gRPC
> `TenancyGuard`), per-project role on the **cluster-operations path** is read from the Mongo
> `Project.members` projection — so add the personas as project members in the console too (the
> console write dual-writes the Keycloak group). After AUTH-1, the verified `groups` claim alone
> drives both the console and the cluster-operations path.

## What to verify (maps to the test plan)

- **3.2 SSO Capabilities:** each persona logs in via Keycloak SSO; the console reflects their role.
- **Management / tenant isolation:** `ccm-alpha-*` users see only alpha; `ccm-bravo-owner` sees only
  bravo; `ccm-no-access` is denied everywhere; `ccm-admin` sees all.
- **3.4 Updates & Upgrades RBAC:** member is blocked from a BB upgrade; maintainer can security-patch;
  owner can run a functional / R-1 / distro upgrade.
- **Visibility / Monitoring / Security (3.5):** viewer can read metrics, cost, audit, and security
  findings but cannot mutate.
