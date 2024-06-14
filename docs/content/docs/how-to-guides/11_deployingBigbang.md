---
title: "Deploying Bigbang"
description: "Steps to Deploy Bigbang"
summary: "Steps to Deploy Bigbang"
draft: false
menu:
  docs:
    parent: "docs/how-to-guides/"
    identifier: ""
weight: 203
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

## Prerequisites
Before you begin, ensure that you have the following prerequisites in place:
- You have created the KMS encryption key to encrypt your cluster and have the ARN of the KMS key handy.
- You have deployed the Kubernetes and have access to `kubeconfig` file.
- All your worker nodes have instance profile with a policy to use KMS:decrypt

## Deploy BigBang

-  Login to Enbuild -[Enbuild](https://enbuild.vivplatform.io)
-  Click on the **Home**
-  Select the **Platform One BigBang**
-  At the SOPS tab provide the `sops.yaml` created on SOPS prerequisite section.

<picture><img src="/images/how-to-guides/SOPS.png" alt="SOPS"></img></picture>

-  At the REPO tab , provide the
   <ol type="a"> 
     <li> Registry URL — The container registry from where you are pulling the images for flux deployment. </li> <li> Registry Username - The container registry username to pull flux images </li>   <li>Registry Password - The container registry password to pull flux images</li> <li>Repository Username -  The gitlab repository username to pull the BigBang Helm chart. We have cloned the chart at <a href="https://gitlab.com/enbuild-staging/charts/bigbang.git" target="_blank">chart</a></lI> <li>Repository Password - The gitlab repository password to pull the BigBang Helm chart. (We have cloned the chart at <a href="https://gitlab.com/enbuild-staging/charts/bigbang.git" target="_blank">chart</a>) </li>
     </ol>

     <picture><img src="/images/how-to-guides/repoCreds.png" alt="Repo Credentials"></img></picture>

- Next, choose the component **Repo** from the **settings** category and click on the **SECRETS** tab and provide the `registryCredentials`and `git credentials`  this is basically used by BigBang Helm chart to pull the container images and cloning the dependant helm charts used by bigbang.

<picture><img src="/images/how-to-guides/repoSecrets.png" alt="Repo Secrets"></img></picture>

  - The values of these will be same as previous section.

```
  registryCredentials:
  registry: registry.gitlab.com
  username: registry_username
  password: registry_password
  email: ""
git:
  credentials:
  username: repository_usernane
  password: registry_password
```

-  Similarly you can check other components and edit the values of the component deployment. If you feel the value is sensitive you can add that in secrets tab, so that enbuild will encrypt it using the KMS key provided before committing to the git repo.

<picture><img src="/images/how-to-guides/securityTools.png" alt="Security Tools"></img></picture>

<picture><img src="/images/how-to-guides/SSO.png" alt="SSO"></img></picture>

<picture><img src="/images/how-to-guides/CI-CD_Tools.png" alt="SSO"></img></picture>

One important component setting while deploying BigBang is

`domain: dev.bigbang.mil` which is present in Settings → Repo → Values. This defines the istio ingress domain on which the bigbang applications will be available.

<picture><img src="/images/how-to-guides/repoValues.png" alt="Repo Values"></img></picture>

You also have to provide the right tls certificate and key for the same domain defined above in the Component → Service Mesh → Istio → Secrets tab. So that you can access the bigbang applications in browser without any security/certificate warning.

<picture><img src="/images/how-to-guides/istioSecrets.png" alt="Istio secrets"></img></picture>

- After providing all the input values, provide the name for your deployment proceed to Infrastructure section, and provide your 
  -    kubeconfig file - Paste your `kubeconfig` file
  -    Select AWS as your cloud and provide your AWS credentials. These AWS credentials are used to encrypt the secrets using SOPS. So make sure the IAM user of these credentials is having `kms:Encrypt` and `kms:Decrypt` permissions.

<picture><img src="/images/how-to-guides/AWS-Bigbang.png" alt="AWS-Bigbang"></img></picture>

  Once all inputs are provided click on **Deploy Stack**

  <picture><img src="/images/how-to-guides/deployStackButton.png" alt="Deploy StacK BUTTON"></img></picture>