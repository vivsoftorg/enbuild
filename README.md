# [ENBUILD-HELM-CHART](https://github.com/vivsoftorg/enbuild_helm_chart.git)


## Usage

This Helm chart repository enables you to install a ENBUILD
Helm chart directly from it into your Kubernetes cluster. Please refer to the
[ENBUILD documentation](https://enbuild-docs.vivplatform.io/o) for all
the additional details required.

```shell
# Let helm the command line tool know about a Helm chart repository
# that we decide to name enbuild.
❯ helm repo add vivsoft https://vivsoftorg.github.io/enbuild_helm_chart

# Update the Helm chart repository.
❯ helm repo update

# Search for the ENBUILD Helm chart in the enbuild Helm chart repository.
❯ helm search repo  vivsoft/enbuild
NAME              	CHART VERSION	APP VERSION	DESCRIPTION
vivsoft/enbuild	1.0.6        	1.0.6      	A Helm chart for ENBUILD

# Simplified example on how to install a Helm chart from a Helm chart repository
# named enbuild. See the Helm chart's documentation for additional details
# required.
❯ helm install vivsoft/enbuild --version 1.0.6
```
