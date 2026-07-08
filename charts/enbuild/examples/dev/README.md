# examples/dev — throwaway eval layers (NOT for delivery)

These are **development/test convenience layers**, split out of the production
overlays on 2026-07-05 (`80b17bf`) during the Azure/GCP Phase-A/B bring-up. They
are **not** part of the customer delivery set and must never be layered on a real
deployment.

| file | what it is |
|---|---|
| `values-aks-eval.yaml` | throwaway eval layer applied *after* `../values-aks.yaml` |
| `values-gke-eval.yaml` | throwaway eval layer applied *after* `../values-gke.yaml` |

Each swaps production defaults for eval shortcuts: a `nip.io` domain (no DNS),
**`authMechanism: local` + break-glass `adminEmails`** (auth effectively open),
and the bundled demo realm. That open-auth posture is exactly why they cannot ship.

**Why they are effectively obsolete:** once the bundled Keycloak becomes turnkey
(domain-only, generated admin password), an engineer can stand up a *real-auth*
eval on Azure/GCP with just the production overlay + a `nip.io` domain — no
local-auth shortcut needed:

```
helm ... -f examples/values-aks.yaml --set global.domain=<LB-IP>.nip.io
```

The delivery decision set is the three overlays one level up:
`values-quickstart.yaml` (local port-forward), `values-aks.yaml`, `values-gke.yaml`.
Provenance for these eval layers is preserved in git history + the
`p1-cluster-mgmt/docs/delivery/2026-07-0{5,6}-hub-on-{aks,gke}-phase-*.md` logs.
