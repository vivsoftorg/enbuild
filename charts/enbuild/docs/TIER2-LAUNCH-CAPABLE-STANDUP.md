# Tier-2 standup — a hub that LAUNCHES and MANAGES spoke clusters

**Audience:** an engineer standing up ENBUILD in their own Platform One / GovCloud
environment, from this chart's released artifact alone.
**Goal:** a hub that not only runs the console (Tier 1), but **launches spoke
clusters and manages them** — the full product (Tier 2).

> **Tier 1** = control-plane eval: console + SSO + catalog browse, reachable by
> port-forward, no spokes. That is [`examples/values-quickstart.yaml`](../examples/values-quickstart.yaml)
> end-to-end — nothing below is needed for it.
> **Tier 2 (this doc)** = Tier 1 **plus** the launch credential, the hub↔spoke
> PKI, and the Istio-terminated gRPC edge that launched agents dial back to.

---

## 0. The one architectural fact everything below follows from

A launched spoke connects **back** to the hub over the hub's gRPC AgentGateway,
and **that gateway is designed to run behind Istio**:

- the chart ships a gRPC **VirtualService** (`<release>-enbuild-bk-grpc-vs`,
  path-prefix `/enbuild.v1.`) that binds to an **existing** Istio Gateway
  (`enbuildBk.grpcVirtualService.gateway`);
- the **Gateway terminates the TLS** the agents dial (SIMPLE mode, HTTP/2 —
  which carries gRPC transparently). The backend gRPC server does **not**
  present its own external TLS cert — it is plain HTTP/2 behind the mesh.

**What the chart deliberately does NOT ship:** the Istio `Gateway`/`Certificate`
pair and the hub CA. Every Platform One target already runs Istio (it *is* Big
Bang), so the chart binds to **your** gateway instead of bundling a second Istio
that would collide with Big Bang's (duplicate CRDs, conflicting webhooks). A bare
cluster with no Istio cannot expose a working agent-gateway by just adding a
LoadBalancer — install an Istio edge first (Appendix A).

## 1. Prerequisites (on the hub cluster)

| Prereq | Why |
|---|---|
| **Istio with an ingress gateway** (e.g. Big Bang's `istio-gateway/public-ingressgateway`) | terminates TLS for the agent gateway + console/Keycloak edge |
| **cert-manager** | issues the hub CA + per-spoke client certs (mTLS) |
| **A default StorageClass** (+ EBS CSI on EKS) | Mongo/RabbitMQ PVCs (if bundled) |
| **Registry entitlements** | the ENBUILD app images + `registry1.dso.mil/ironbank/*` dependencies |
| **A DNS name** for the hub gRPC host, pointed at the gateway's load balancer | spokes dial it (`installAgent.hubUrl`) |
| **A GitLab token** that can create projects in your deployments group | prefer a **group access token** (Maintainer, `api`+`write_repository`) — SSO-exempt and portable |

## 2. Secrets the operator pre-creates

Run [`scripts/create-bootstrap-secrets.sh`](../scripts/create-bootstrap-secrets.sh)
for the base set (mongo, encryption key, rabbitmq, image pull), **plus the launch
credential** — the Secret that makes the hub Tier-2:

```bash
kubectl -n enbuild create secret generic enbuild-install-agent-creds \
  --from-literal=GITLAB_TOKEN='<GitLab token: create deploy repos + trigger pipelines>' \
  --from-literal=ENBUILD_REPO1_USER='<registry1.dso.mil pull user>' \
  --from-literal=ENBUILD_REPO1_TOKEN='<registry1.dso.mil pull token>' \
  --from-literal=REAPER_SVC_AWS_ACCESS_KEY_ID='<aws key for the teardown reaper>' \
  --from-literal=REAPER_SVC_AWS_SECRET_ACCESS_KEY='<aws secret>'
```

Without the `REAPER_SVC_*` keys teardowns still run, but the cloud-release gate
times out and every destroy orphans the spoke's Istio gateway ELB — see the
`enbuildBk.reaperSvc` comment in [`values.yaml`](../values.yaml).

## 3. Hub PKI (mTLS hub↔spoke) — one-time CA bootstrap

The chart only *recreates* the `enbuild-hub-issuer` ClusterIssuer against an
**already-existing** CA secret (`pki.recreateHubIssuer`; it never generates or
rotates the CA — regenerating would break every spoke that trusts it). On a
fresh cluster, bootstrap the CA once, out-of-band:

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: enbuild-selfsigned-bootstrap
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: enbuild-ca
  namespace: cert-manager          # cert-manager's cluster-resource namespace
spec:
  isCA: true
  commonName: enbuild-hub-ca
  secretName: enbuild-ca-tls       # must match pki.caSecretName
  issuerRef:
    name: enbuild-selfsigned-bootstrap
    kind: ClusterIssuer
---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: enbuild-hub-issuer         # must match pki.hubIssuerName
spec:
  ca:
    secretName: enbuild-ca-tls
```

With the issuer created out-of-band as above, keep `pki.recreateHubIssuer: false`.
(Set it `true` for one roll only if the ClusterIssuer object ever goes missing —
it re-points at the existing CA secret, nothing more.) The install-agent flow
then issues each spoke a `<cluster>-spoke-client` certificate off this issuer.

## 4. Values that make the hub launch + connect

On top of your Tier-1 values (quickstart or your own):

```yaml
enbuildBk:
  installAgent:
    existingSecret: enbuild-install-agent-creds
    hubUrl: <hub-grpc-host>:443          # what spoke agents dial back to
    tlsServerName: <cert-SAN>            # SNI override, only if the DNS host differs from the cert SAN
  grpcVirtualService:
    enabled: true
    gateway: <ns>/<gateway>              # YOUR Istio gateway, e.g. istio-gateway/public-ingressgateway
    host: <hub-grpc-host>
  # GitLab connection for the launch worker (a launch creates a deploy repo
  # in this group and runs its pipeline):
  gitlabConnection:
    host: https://gitlab.com/            # or your self-hosted GitLab
    namespaceId: "<deploy-group-id>"     # numeric id of the deployments group
```

Plus AWS credentials for the launch pipeline itself: group-level CI variables
(`AWS_ACCESS_KEY_ID`/`AWS_SECRET_ACCESS_KEY`) on the deployments group, **or**
the console's Admin Settings → platform AWS credentials. Pin every
`<svc>.image.tag` to a published build (the quickstart carries known-good pins —
the empty default resolves to the chart `appVersion`, which is **not** guaranteed
deployable).

## 5. DNS, install, verify

1. Create the DNS record: `<hub-grpc-host>` → the Istio gateway's load balancer.
2. Install **from the released chart artifact** (subchart dependencies bundled):
   ```bash
   gh release download enbuild-trunk-<version> \
     --repo vivsoftorg/enbuild --pattern '*.tgz' -D /tmp
   helm upgrade --install enbuild-ib /tmp/enbuild-<version>.tgz \
     -n enbuild --create-namespace -f <your-values>.yaml
   ```
3. Verify, before launching anything:
   - `istioctl proxy-config route <ingressgateway-pod> -n <gw-ns>` shows
     `/enbuild.v1.` routed to `<release>-enbuild-bk-grpc`;
   - `POST /clusters/install-agent` returns **400** (validation), not **503**
     (503 = the installAgent Secret is missing);
   - catalog tiles render ACTIVE in the console.

## 6. The two console gates every fresh hub hits (do these BEFORE the first launch)

1. **Create a Project and make it the active Project.** A fresh hub has zero
   Projects, so the header sits on "All Projects" (an aggregate view with an
   *empty* project id) and submitting the catalog wizard fails with
   `Launch failed — project '' not found`. Fix: **Projects → + New Project**
   (name/slug/impact level; platform-admin only), then select it in the header
   Project switcher. Launches are scoped to the active Project's id.
2. **Restart the backend after any `KEYCLOAK_ISSUER` change.** The backend reads
   its issuer **once at boot** via `envFrom` — a `helm upgrade` that changes it
   does *not* restart the pod. The stale pod keeps validating against the old
   issuer, producing a telltale **mixed pattern**: strict endpoints
   (`/api/v1/manifests`, `/stacks`, `/mission-apps`, `/admin-settings/*`) return
   401 while lenient ones (`/fleet-health`, `/projects`) return 200 — and the
   catalog hangs on "Loading templates…". Fix:
   `kubectl -n <ns> rollout restart deploy/<release>-enbuild-backend`.

With both gates cleared, the catalog wizard submits end-to-end: stack create →
deploy repo created → pipeline → Terraform apply → agent install → the new spoke
shows **connected** in fleet health.

---

## Appendix A — greenfield Istio edge (no Big Bang on the hub cluster)

Validated sequence for a bare cluster (e.g. fresh EKS) with no Istio:

1. **Install Istio** (match your fleet's version) via `istioctl install`, with an
   ingress gateway of `service.type: LoadBalancer` (internal NLB annotations as
   your network requires).
2. **Bootstrap the hub CA** (§3 above).
3. **Create the gateway + server cert:** a `Certificate` (issued off
   `enbuild-hub-issuer`, secret in `istio-system`, SAN = the gRPC host) and a
   `Gateway` (`selector: {istio: ingressgateway}`, port 8443, `tls.mode: SIMPLE`,
   `credentialName: <that secret>`, `hosts: ["*"]`).
4. **Install/roll the hub** with §4's values, `gateway:` pointing at the new
   Gateway and `host`/`hubUrl`/`tlsServerName` = the LB's DNS name.
5. **Put the backend in the mesh** so the gateway→backend hop works:
   `kubectl label ns enbuild istio-injection=enabled`, a
   `PeerAuthentication {mtls: {mode: PERMISSIVE}}` in the namespace, then
   `kubectl -n enbuild rollout restart deploy/<release>-enbuild-backend` (→ 2/2).
6. **Verify** per §5.3 — no spoke needed.

Because the gateway cert chains to the self-signed hub CA, launched agents must
either trust that CA or the launch must set the agent's insecure-TLS toggle
(`enbuild_agent_hub_insecure: true`) — eval only; use a real cert for production.
