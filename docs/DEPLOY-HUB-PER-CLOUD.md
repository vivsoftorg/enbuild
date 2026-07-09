# Deploy the ENBUILD hub-per-cloud (AKS / GKE) — start here

This is the **entry point** for standing up an ENBUILD hub on a commercial cloud
(Azure AKS or Google GKE) and onboarding same-cloud spokes. Follow it top to bottom;
each step links to the detailed doc for that piece.

> **Which path am I on?**
> - **AWS GovCloud / Big Bang (IL4+):** the hub runs inside Big Bang. Use the Big Bang
>   path — see [`charts/enbuild/docs/LAUNCH-CAPABLE-STANDUP.md`](../charts/enbuild/docs/LAUNCH-CAPABLE-STANDUP.md).
>   This document does **not** apply.
> - **Commercial AKS / GKE (this document):** the hub runs on a vanilla managed
>   Kubernetes cluster with an Istio edge. Spokes are launched into the **same cloud**
>   so connect-back works without crossing a cloud/impact-level boundary. Proven
>   end-to-end on both AKS and GKE.

The model: deploy the hub **in-network**, launch same-cloud spokes from the console
catalog, and the spoke pull-agent dials the hub's gRPC AgentGateway back through the
**same Istio ingressgateway** (HTTP/2 carries gRPC) — no AWS NLB, no IRSA.

---

## 0. Pick your cloud + overlay

The chart ships four per-CSP overlays under `charts/enbuild/examples/`. You always
layer the **production** overlay on top of the chart defaults; for a throwaway eval you
add the matching `-eval` overlay **after** it.

| Cloud | Production overlay | Throwaway eval overlay |
|-------|--------------------|------------------------|
| Azure AKS | `examples/values-aks.yaml` | `examples/dev/values-aks-eval.yaml` |
| Google GKE | `examples/values-gke.yaml` | `examples/dev/values-gke-eval.yaml` |

The eval overlays swap in a `nip.io` domain (no DNS setup), **local auth** (no external
Keycloak IdP), the bundled demo realm, and cheaper storage. **Never use an `-eval`
overlay for a real deployment** — it ships the hub with authentication effectively open.

---

## 1. Prerequisites (before you deploy)

1. **A managed Kubernetes cluster** (AKS or GKE) sized for the hub (~8 pods + Mongo/
   RabbitMQ/Keycloak). For onboarding spokes you also need spoke clusters in the **same
   cloud**.
2. **Istio** installed with an **ingressgateway Service of `type: LoadBalancer`** (AKS
   Istio add-on, GKE ASM, or upstream istioctl). Note its **namespace/name** — you set
   the `gateway:` refs in the overlay to match. The chart renders VirtualServices, not
   the Gateway; **you create the Istio `Gateway` object** (one, covering console +
   Keycloak + hub-grpc hosts).
3. **DNS** (or `nip.io` for an eval): records for `<console>.<domain>`,
   `keycloak.<domain>`, and `hub-grpc.<domain>` all pointing at the ingressgateway LB IP.
4. **A TLS cert on the Istio Gateway** covering all three hostnames as SANs (cert-manager
   + Let's Encrypt / Google CAS / an Azure- or DoD-PKI-issued cert). Istio terminates
   TLS; HTTP/2 then carries both the console and the spoke gRPC. For an eval a self-signed
   cert is fine (see the browser-trust note in the eval overlays).
5. **Image-pull entitlements** for `registry.gitlab.com` + `registry1.dso.mil`.
6. **☁️ CLOUD QUOTA — check before the first launch (this is the #1 avoidable failure):**
   - **GKE:** the launch region needs **Persistent Disk SSD (`SSD_TOTAL_GB`)** quota.
     Defaults are `pd-balanced`, 100 GB/node, up to 3 nodes/zone → ~300+ GB per cluster;
     a fresh project's 500 GB regional default is easily exhausted by the hub + a spoke.
     Raise it at *IAM & Admin → Quotas*, pick a region with headroom, or set
     `disk_type=pd-standard`.
   - **AKS:** an **entire VM family** (e.g. `Dsv5`) can be **quota-0** in a region even
     for small SKUs. Pre-flight: `az vm list-usage --location <region> -o table | grep -i <family>`.
   - Both are **account limits, not platform defects** — but they surface as cryptic
     apply-time errors if unchecked. See the catalog `OPERATOR-SETUP.md` for details.

---

## 2. Create the bootstrap secrets

Run the generator once per hub cluster — it produces every `enbuild-ib-*` secret the
overlays reference (Mongo, RabbitMQ, encryption-key, messaging, image-pull, install-agent):

```bash
charts/enbuild/scripts/create-bootstrap-secrets.sh          # RELEASE defaults to enbuild-ib
```

See [`docs/OPERATOR-DEPLOYMENT-GUIDE.md`](OPERATOR-DEPLOYMENT-GUIDE.md) §1 for the full
secret table.

> **Production Keycloak secret (manual):** the generator does **not** create the Keycloak
> realm secret. For the default `authMechanism: keycloak` path you must hand-craft
> `enbuild-ib-keycloak-secrets` (keys `realm-enbuild.json` + `KC_BOOTSTRAP_ADMIN_PASSWORD`)
> out-of-band before deploying. The **eval** overlays use `authMechanism: local` + the
> bundled demo realm, so they need no Keycloak secret.

---

## 3. Deploy the hub

Point the overlay's placeholders at your real values first (`REPLACE`, `example.com`,
gateway refs, `gitlabConnection`). Then:

```bash
# PRODUCTION (AKS shown; swap values-gke.yaml for GKE)
helm upgrade --install enbuild-ib <enbuild-chart>.tgz \
  -n enbuild --create-namespace \
  -f charts/enbuild/examples/values-aks.yaml

# THROWAWAY EVAL — layer the eval overlay AFTER the production one
helm upgrade --install enbuild-ib <enbuild-chart>.tgz \
  -n enbuild --create-namespace \
  -f charts/enbuild/examples/values-aks.yaml \
  -f charts/enbuild/examples/values-aks-eval.yaml
```

The chart auto-publishes to GitHub Releases (tag `enbuild-trunk-<version>`); pull the
`.tgz` from there, or `helm package charts/enbuild`.

---

## 4. Get the LB IP + finish DNS

```bash
kubectl -n <ingressgateway-namespace> get svc <ingressgateway> \
  -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
```

Point your DNS records (or set the eval `nip.io` domain to `<ip>.nip.io`) at that IP for
all three hosts, then re-run the helm command with the resolved `global.domain` /
`grpcVirtualService.host` / keycloak hosts.

## 5. Verify the hub

All hub pods `Running` in `enbuild`; the console loads at `https://<console>.<domain>/p1-ccm-console/`
and you can log in (Keycloak SSO in prod; the demo `admin@p1.mil` user in eval).

---

## 6. Set up the same-cloud catalog (for launching spokes)

The console launches spokes by forking a Terraform catalog repo per cluster. Do the
one-time catalog operator setup for your cloud:

- **GKE:** [`platform-one-gke/docs/OPERATOR-SETUP.md`](https://gitlab.com/enbuild-staging/iac-templates/platform-one-gke/-/blob/main/docs/OPERATOR-SETUP.md)
  — WIF (keyless) or SA-key auth, **and the REQUIRED group-level CI variables**
  (`GOOGLE_PROJECT`, `GCP_WORKLOAD_IDENTITY_PROVIDER`, `GCP_SERVICE_ACCOUNT`) on the
  `enbuild-deployments` GitLab group.
- **AKS:** [`platform-one-azure/docs/OPERATOR-SETUP.md`](https://gitlab.com/enbuild-staging/iac-templates/platform-one-azure/-/blob/main/docs/OPERATOR-SETUP.md)
  — service-principal auth (ARM_* injected per launch from console Platform Settings),
  node SKU + quota.

## 7. Launch a spoke

From the console: **Catalog → pick the AKS/GKE template → create a Project → Launch**.
The catalog CI provisions the cluster and installs the pull-agent in-band.

## 8. Verify connect-back

The spoke agent must trust the hub's cert and dial the gRPC endpoint. For a self-signed
(eval) hub cert this needs one manual step — see
[`charts/enbuild/docs/CONNECT-BACK-SAME-CLOUD.md`](../charts/enbuild/docs/CONNECT-BACK-SAME-CLOUD.md).
When it works you'll see the spoke as **connected** in the console's Fleet Health, with
live pod/service/deployment counts flowing to the hub.

---

## Where everything lives

| Piece | Location |
|-------|----------|
| Hub chart + overlays + secrets script | `charts/enbuild/` (this repo) |
| This entry point | `docs/DEPLOY-HUB-PER-CLOUD.md` |
| Hub deploy detail (secrets, Keycloak, Headlamp) | `docs/OPERATOR-DEPLOYMENT-GUIDE.md` |
| Same-cloud connect-back procedure | `charts/enbuild/docs/CONNECT-BACK-SAME-CLOUD.md` |
| GKE spoke catalog + operator setup | `platform-one-gke` repo |
| AKS spoke catalog + operator setup | `platform-one-azure` repo |
| AWS GovCloud / Big Bang path | `charts/enbuild/docs/LAUNCH-CAPABLE-STANDUP.md` |
