# ENBUILD Helm chart — Troubleshooting

Common install / upgrade issues, what causes them, and how to fix them. Each entry is "**symptom**" (what you see) → "**cause**" → "**fix**".

The first sections cover bugs that have already been fixed in the chart but may bite you on older installs. Later sections cover environment-side issues outside the chart that can look like chart bugs.

---

## Chart-side issues (fixed in chart >= 0.0.43)

### "UI loads but every API call returns React HTML / `Unexpected token '<'`"

**Symptom:** ENBUILD UI page loads in the browser. Any user action triggers a JavaScript exception:

```
Uncaught (in promise) SyntaxError: Unexpected token '<', "<html>
<h"... is not valid JSON
```

DevTools Network tab shows requests to `/enbuild-bk/api/...` returning 404 with an HTML body.

**Cause:** Older chart versions had `server_name _;` (only) on the nginx server block. The base nginx image ships another server block listening on the same port with `server_name localhost;`. nginx selects server blocks by Host header before falling back to `default_server`, so any request with `Host: localhost` (i.e. every `kubectl port-forward` or local proxy) matched the base block — which has no `/enbuild-bk/`, `/enbuild-user/`, `/headlamp/` proxy paths. API calls returned the React app's 404 HTML; the frontend tried to JSON-parse it and crashed.

**Fix:**
- **Chart >= 0.0.43**: nginx server block now claims `server_name localhost _;`. Upgrade.
- **Older chart, can't upgrade**: patch the live ConfigMap, then restart the UI deployment:
  ```
  kubectl -n enbuild get cm <release>-nginx-conf -o yaml > /tmp/cm.yaml
  sed -i 's|server_name _;|server_name localhost _;|' /tmp/cm.yaml
  kubectl apply -f /tmp/cm.yaml
  kubectl -n enbuild rollout restart deployment/<release>-enbuild-ui
  ```

**Verify:** `helm test <release> -n enbuild` — the bundled `test-ui-proxy` test exercises the nginx routing chain and will fail on this regression.

---

### "mq-consumer pod restarts every minute / `libatomic.so.1` errors"

**Symptom:** `kubectl get pod -l app=mq` shows a high `RESTARTS` count (often 30+ within an hour). `kubectl get events` shows:
```
Warning   Unhealthy   pod/<release>-enbuild-mq-...  Liveness probe failed: node: error while loading shared libraries: libatomic.so.1: cannot open shared object file: No such file or directory
```

**Cause:** Older chart versions ran the liveness probe as `node src/queue/testConnection.js` directly. The Iron Bank `enbuild-mq` image ships multiple `node` binaries on PATH; only `/usr/bin/node` is built against `libatomic.so.1`. Plain `node` resolved to the wrong binary, the probe failed, and the kubelet killed the pod every probe period (60s).

**Fix:**
- **Chart >= 0.0.43**: default probe is now a shell wrapper `PATH=/usr/bin:/usr/sbin:/bin:/sbin:/usr/local/bin /usr/bin/node src/queue/testConnection.js`. Upgrade.
- **Older chart, can't upgrade**: patch the deployment in place:
  ```
  kubectl -n enbuild patch deployment <release>-enbuild-mq --type=json -p '[
    {"op":"replace","path":"/spec/template/spec/containers/0/livenessProbe/exec/command","value":["/bin/sh","-c","PATH=/usr/bin:/usr/sbin:/bin:/sbin:/usr/local/bin /usr/bin/node src/queue/testConnection.js"]}
  ]'
  ```

If you've overridden `enbuildConsumer.livenessProbe` in your values file, keep the `/usr/bin/node` invocation (or the equivalent absolute path on your image).

---

## Install-time issues

### "ImagePullBackOff on enbuild service pods"

**Symptom:** Pods stuck in `ImagePullBackOff` immediately after `helm install`.

**Cause:** Iron Bank registry credentials weren't supplied. The chart creates the image pull secret from `global.image.registry_credentials.username` and `.password` — if these are empty or wrong, the secret can't authenticate to `registry1.dso.mil`.

**Fix:** Pass credentials at install time. Never commit them.
```
helm install <release> ./charts/enbuild --namespace enbuild --create-namespace \
  --values examples/enbuild/values-vendor13-ib.yaml \
  --set global.image.registry_credentials.username=<your registry1 username> \
  --set global.image.registry_credentials.password=<your registry1 password>
```

Test the credentials separately first:
```
docker login registry1.dso.mil -u <user>
docker pull registry1.dso.mil/ironbank/vivsoft/enbuild/backend:1.0.30-5101272
```

---

### "MongoDB pod fails to start with `<REPLACE_BEFORE_INSTALL>` errors"

**Symptom:** mongodb pod won't reach Ready. Logs show password validation failures.

**Cause:** Default `mongo_root_password` value in `examples/enbuild/values-vendor13-ib.yaml` is the sentinel `<REPLACE_BEFORE_INSTALL>`, which is invalid as a MongoDB password. This is intentional — it forces the operator to set a real password instead of inheriting an unsafe placeholder.

**Fix:** Pass the password at install time:
```
helm install <release> ./charts/enbuild ... \
  --set mongodb.mongo_root_password='<your secret>'
```

For a longer-lived install, prefer wiring to an externally-managed Kubernetes Secret using the mongodb subchart's `existingSecret` value (not enabled by default in the example).

---

## Backend-side issues (fix lives in `vivsoft-platform-ui`, not chart)

### "`GET /api/v1/stacks` returns 500 with `$regex has to be a string`"

**Symptom:** Direct API call to `GET /enbuild-bk/api/v1/stacks` returns:
```
{"statusCode":500,"message":"Internal server error"}
```
Backend logs show `MongoServerError: $regex has to be a string`.

**Cause:** Backend's `getAll` controller builds a Mongo filter with `$regex: query.search`. If `search` is `undefined` (no query param), the resulting `$regex: undefined` errors at the driver layer.

**Workaround:** Always include `?search=` (empty value is fine):
```
curl -H "Authorization: Bearer $TOKEN" \
  'http://localhost:3000/enbuild-bk/api/v1/stacks?search=&page=1&limit=10'
```

The ENBUILD UI itself always sends a `search` param so this only affects direct API consumers. The fix belongs in `vivsoft-platform-ui/backend/microservices/enbuild/src/stacks/stacks.service.ts` — guard the `$regex` clause to be added only when `search` is a non-empty string. Tracked separately.

---

## Cluster / environment issues outside the chart

### "Loki receives no logs from ENBUILD pods"

**Symptom:** ENBUILD pods are happily Running and producing logs, but Loki queries (e.g. `{namespace="enbuild"}`) return empty.

**Cause:** Big Bang's `fluentbit` HelmRelease may be in a failed state on the cluster — not an ENBUILD problem. Check first:
```
kubectl get hr -A | grep -v True
kubectl -n bigbang describe hr fluentbit
```

If `fluentbit` is failed, no namespace gets logs ingested — including ENBUILD. This was observed on the vendor13-ib install (April 2026); a fix from the cluster platform team was needed.

**Fix:** Owned by the cluster operator / Big Bang admin. ENBUILD logs are flowing into stdout normally; the missing piece is the cluster's log shipper. Ask the platform team to fix `fluentbit` (or whichever Big Bang log component is failing).

### "kind cluster on macOS Colima: kube-proxy CrashLoopBackOff with 'too many open files'"

**Symptom:** Local `kind create cluster` succeeds but `kube-proxy` is in `CrashLoopBackOff`. Logs:
```
"command failed" err="failed complete: too many open files"
```

**Cause:** Default inotify limits in the Colima VM are too low for kube-proxy.

**Fix:**
```
docker exec <kind-cluster-name>-control-plane sh -c \
  "sysctl -w fs.inotify.max_user_instances=8192 && sysctl -w fs.inotify.max_user_watches=524288"
kubectl -n kube-system delete pod -l k8s-app=kube-proxy
```

The new CI workflow `.github/workflows/helm-pr-checks.yml` applies the same fix automatically.

### "SSM port-forward to vendor13-ib drops every few hours"

**Symptom:** `kubectl get nodes` against the SSM-tunneled vendor13-ib cluster returns "connection refused" after several hours of working.

**Cause:** SSM session enters a half-dead state — the AWS-side reports the session as connected but the local socket has dropped.

**Fix:** Kill the tunnel process and restart it:
```
kill $(cat /tmp/vendor13-tunnel.pid)
./scripts/vendor13-ib-tunnel.sh > /tmp/vendor13-tunnel.log 2>&1 &
echo $! > /tmp/vendor13-tunnel.pid
```

Operator-side issue, not chart. Captured in the vendor13-ib runbook for the team.

---

## Future enhancements (deferred — not bugs)

### DNS hostname / Istio VirtualService for browser-via-VPN access

The chart supports Istio integration via `global.istio.enabled` and `enbuildUi.hostname`. The vendor13-ib install does not enable this yet — it uses `kubectl port-forward` for local access. When DNS for ENBUILD is provisioned (Route 53 record pointing at the cluster's public ingress ELB), enable Istio and set the hostname in the values file.

### IRSA / Pod Identity for the mq-consumer ServiceAccount

For installs where the mq-consumer needs to call AWS APIs directly (e.g., to discover or operate on AWS resources outside of user-driven catalog flows), annotate the ServiceAccount with an IAM role using IRSA or EKS Pod Identity. The cluster's OIDC provider is already registered for vendor13-ib (`arn:aws-us-gov:iam::219533807077:oidc-provider/oidc.eks.us-gov-west-1.amazonaws.com/id/48EA4C3B972522B4A77C3643EAFE1512`) and `eks:CreatePodIdentityAssociation` is in the role policy. Not configured today because end-users supply per-stack AWS credentials via the catalog UI.

---

## When in doubt

1. Run `helm test <release> -n <namespace>` — it specifically exercises the nginx proxy chain and would catch the most common regression class.
2. Check `kubectl get hr -A | grep -v True` for any failed Big Bang components — they often look like ENBUILD problems.
3. Compare against `examples/enbuild/quick_install.yaml` (non-IB, public images) to isolate IB-specific issues.
