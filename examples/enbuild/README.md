# Chart example values

The YAML files in this directory are topology and image overlays. They are not complete credentials, identity, networking or cluster-launch configurations. For the full AWS CCM platform, use the [repo1 IaC AWS guide](https://repo1.dso.mil/ccm/vivsoft/ccm-platform-iac/-/blob/main/terraform/hub-cluster/aws-eks/README.md) with its explicit verified image profile and your environment inputs. For the running 7077 application, use the [evaluator handoff](https://repo1.dso.mil/ccm/vivsoft/operator-guide/-/blob/main/delivery/7077/README.md); it requires no new deployment.

## Disposable local chart smoke

`quick_install.yaml` starts the core application with public images and local/demo authentication. It does not configure catalog launches, imports, strict SSO, Headlamp or spoke connectivity, and is not the 7077 deployment posture. Keep the existing 7077 HTTPS/Keycloak/Headlamp configuration.

On an isolated local cluster with a default StorageClass, from this repository root:

```sh
helm dependency build charts/enbuild
RELEASE=enbuild-ib NAMESPACE=enbuild bash charts/enbuild/scripts/create-bootstrap-secrets.sh
helm install enbuild-ib charts/enbuild --namespace enbuild --create-namespace \
  --values examples/enbuild/quick_install.yaml --wait --timeout 10m
kubectl -n enbuild port-forward service/enbuild-ib-enbuild-ui 8080:80
```

The bootstrap script generates local MongoDB, encryption and matching RabbitMQ/messaging credentials. It preserves existing Secrets by default; use a fresh namespace/cluster for this smoke. The fixture's Secret names require release `enbuild-ib` and namespace `enbuild`. Private registry, cloud and catalog tokens are not needed for this public-image smoke. The archived public Bitnami RabbitMQ image is a chart-layout fixture, not a production image recommendation.

## CI render checks

CI combines every YAML overlay with `charts/enbuild/ci/render-environment.values.yaml`, which provides synthetic required full-hub bindings solely for schema rendering. Those `.invalid` origins and Secret names are not deployable settings or working credentials. Helm must succeed and emit a nonempty manifest before kubeconform runs. The separate kind job installs the standalone quick profile with actual generated disposable Secrets and checks runtime availability, UI proxy paths and MQ probe stability.

A successful local smoke or synthetic render does not prove full CCM functionality, a new AWS hub deployment or strict browser authentication. Those require the environment's complete deployment procedure and acceptance tests.
