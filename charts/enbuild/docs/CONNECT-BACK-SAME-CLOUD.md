# Same-cloud spoke connect-back — the CA-trust + HUB_URL procedure

When you launch a spoke into the **same cloud** as the hub, the spoke's pull-agent dials
the hub's gRPC AgentGateway (`hub-grpc.<domain>:443`) through the hub's Istio
ingressgateway. For that TLS/gRPC handshake to complete, the agent must **trust the hub's
server certificate** and know **where to dial**. This is set in the console's Agent admin
settings **before** you launch the spoke.

> Skipping this is the single most common connect-back stall: the agent comes up `1/1
> Running`, dials forever, and logs `grpc.Dial ... context deadline exceeded`. That is a
> **TLS-trust** symptom, not a network one — don't go down the firewall rabbit hole.

## When is this a manual step?

| Hub cert | Manual CA step needed? |
|----------|------------------------|
| **Self-signed / private CA** (all evals; some prod) | **Yes** — you must give the agent the hub CA (below). |
| **Publicly-trusted** (Let's Encrypt / Google CAS / a CA already in the system trust store) | No — the agent's system CA pool already trusts it; set only `HUB_URL` + `TLS_SERVER_NAME`. |

## Procedure

**1. Deploy the hub and get the ingressgateway LB IP** (see
[`docs/DEPLOY-HUB-PER-CLOUD.md`](../../../docs/DEPLOY-HUB-PER-CLOUD.md) steps 3–4):

```bash
kubectl -n <ingressgateway-namespace> get svc <ingressgateway> \
  -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
```

**2. Set the hub domain** — `global.domain` (eval: `<ip>.nip.io`), `grpcVirtualService.host`
(`hub-grpc.<domain>`), and the keycloak hosts — and re-run helm so the gateway serves
those SANs.

**3. (Self-signed hub only) Extract the hub gateway CA:**

```bash
kubectl -n <istio-cert-namespace> get secret <gateway-tls-secret> \
  -o jsonpath='{.data.ca\.crt}' | base64 -d > hub-ca.crt
# (if the gateway cert has no ca.crt key, use the issuing CA you created for the Gateway cert)
```

**4. Set the Agent admin settings in the console** — **Admin Settings → Agent**:

| Field | Value |
|-------|-------|
| `HUB_URL` | `hub-grpc.<domain>:443` |
| `TLS_SERVER_NAME` | `hub-grpc.<domain>` |
| `HUB_CA_BUNDLE_PEM` | *(self-signed only)* paste the contents of `hub-ca.crt` |

These are injected into every subsequent launch. (The same can be set via the backend API
`PUT /api/v1/catalog/admin-settings/agent`.)

**5. Launch the spoke** (console: Catalog → template → Project → Launch). The catalog
delivers the hub CA to the agent automatically — the mechanism differs per cloud:

- **GKE** (`platform-one-gke/infra/src/enbuild-agent.tf`): the launch passes
  `ENBUILD_HUB_CA_BUNDLE_PEM` → Terraform creates the `enbuild-spoke-ca-bundle`
  ConfigMap and sets `agent.hubCaBundle.enabled=true` (plain-TLS server-cert trust when
  mTLS is off). No manual step on the spoke.
- **AKS** (`platform-one-azure` `install-agent`): the install script writes the hub CA
  bundle from the injected value and points the agent at it.

## Verify

```
# spoke agent log
"agent stream open"  hub=hub-grpc.<domain>:443  clusterId=<cluster>
# hub backend log
agent.connect cluster=<cluster> ...   + agent.heartbeat every ~10s
```

In the console **Fleet Health** the spoke shows **connected**, and **Config Health**
starts reporting its live pod/service/deployment counts.

## If it still stalls

- `context deadline exceeded` with the agent `Running` → CA trust. Confirm `HUB_CA_BUNDLE_PEM`
  is the **gateway's** CA and `TLS_SERVER_NAME` exactly matches the gRPC host (SNI is
  load-bearing; a wildcard cert won't match a 4-label host unless it's a SAN).
- Verify raw reachability first: from a spoke pod, `nc -z -w5 <lb-ip> 443` should succeed
  (private GKE clusters need Cloud NAT — the catalog provisions it).
- **Cross-cloud / cross-impact-level** hubs are a different problem (the hub's gateway may
  be an internal-only NLB) — that's out of scope here; keep the hub and its spokes in the
  same cloud.
