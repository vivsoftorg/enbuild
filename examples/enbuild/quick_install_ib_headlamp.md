# Iron Bank Headlamp Example

This example extends `quick_install_ib.yaml` for installs that also need Headlamp enabled with Iron Bank images.

It captures the settings that were validated on a disposable GovCloud EKS cluster:

- ENBUILD core services pulled from Iron Bank
- Headlamp pulled from Iron Bank
- Headlamp kubeconfig init container pulled from Iron Bank
- explicit `storageClass` for stateful dependencies when the cluster has no default StorageClass
- MQ runtime override for the Iron Bank MQ image so `npm` resolves `node` from `/usr/bin`

Use this example when you want:

- Iron Bank-backed ENBUILD images
- Headlamp enabled through `lightning_features.operations_lightning.headlamp`
- in-cluster Headlamp access through `/headlamp/`

Important notes:

1. StorageClass

If your target cluster does not define a default StorageClass, set `global.storageClass` to the class you want ENBUILD to use. The validated GovCloud install used `gp2`.

2. MQ runtime

The Iron Bank MQ image required the example's explicit command/args override so `npm` uses the working `/usr/bin/node` path.

3. Headlamp images

The example overrides both the main Headlamp image and the kubeconfig init container image to approved Iron Bank references.

4. Image pull secret name

The chart creates the image pull secret as `<release-name>-image-pull-secret`. This example assumes the Helm release name is `enbuild-ib`, so the subchart image pull secret references are set to `enbuild-ib-image-pull-secret`.

If you install with a different release name, update the RabbitMQ and Headlamp pull secret references in the example values accordingly.
