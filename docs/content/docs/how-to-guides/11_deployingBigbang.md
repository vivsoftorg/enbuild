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

# Bigbang

Big Bang is a declarative, continuous delivery tool for deploying DoD hardened and approved packages into a Kubernetes cluster. More details are found [Here](https://enbuild-docs.vivplatform.io/docs/references/platform-one-big-bang/)

## Prerequisites
Before you begin, ensure that you have the following prerequisites in place:
- You have created the KMS encryption key to encrypt your cluster and have the ARN of the KMS key handy.
- You have deployed the Kubernetes and have access to `kubeconfig` file.
- All your worker nodes have instance profile with a policy to use KMS:decrypt

## KMS encryption key

The KMS encryption provider uses an envelope encryption scheme to encrypt data in etcd. The data is encrypted using a data encryption key (DEK). The DEKs are encrypted with a key encryption key (KEK) that is stored and managed in a remote KMS

## Create the KMS encryption key

You can create manually or via automation. Its of type **Customer-managed keys** and make sure the person deploying the bigbang is added into the key policy to allow ` kms:Encrypt` and `kms:Decrypt`

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "1",
            "Effect": "Allow",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::986602297069:user/jmemon@vivsoft.io",
                    "arn:aws:iam::986602297069:user/tkaza@vivsoft.io"
                ]
            },
            "Action": [
                "kms:Update*",
                "kms:UntagResource",
                "kms:TagResource",
                "kms:ScheduleKeyDeletion",
                "kms:Revoke*",
                "kms:Put*",
                "kms:List*",
                "kms:Get*",
                "kms:Enable*",
                "kms:Disable*",
                "kms:Describe*",
                "kms:Delete*",
                "kms:Create*",
                "kms:CancelKeyDeletion"
            ],
            "Resource": "*"
        },
        {
            "Sid": "2",
            "Effect": "Allow",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::986602297069:user/jmemon@vivsoft.io",
                    "arn:aws:iam::986602297069:user/tkaza@vivsoft.io"
                ]
            },
            "Action": [
                "kms:ReEncrypt*",
                "kms:GenerateDataKey*",
                "kms:Encrypt",
                "kms:DescribeKey",
                "kms:Decrypt"
            ],
            "Resource": "*"
        },
        {
            "Sid": "3",
            "Effect": "Allow",
            "Principal": {
                "AWS": [
                    "arn:aws:iam::986602297069:role/demo-control-plane",
                    "arn:aws:iam::986602297069:role/demo-worker"
                ]
            },
            "Action": [
                "kms:DescribeKey",
                "kms:Decrypt"
            ],
            "Resource": "*"
        }
    ]
}
```

# KMS Policy Variables

## Policy Version

- **Version** (`string`)
  - **Description:** The version of the policy language.
  - **Default value:** `"2012-10-17"`

## Statements

### Statement 1

- **Sid** (`string`)
  - **Description:** The statement identifier for the first statement.
  - **Default value:** `"1"`

- **Effect** (`string`)
  - **Description:** The effect of the statement, which allows or denies the specified actions.
  - **Default value:** `"Allow"`

- **Principal** (`object`)
  - **Description:** The AWS principals (users) that the statement applies to.
  - **Properties:**
    - **AWS** (`list of strings`)
      - **Description:** List of AWS IAM user ARNs.
      - **Examples:**
        - `"arn:aws:iam::986602297069:user/jmemon@vivsoft.io"`
        - `"arn:aws:iam::986602297069:user/tkaza@vivsoft.io"`

- **Action** (`list of strings`)
  - **Description:** List of actions that are allowed.
  - **Default values:**
    - `kms:Update*`
    - `kms:UntagResource`
    - `kms:TagResource`
    - `kms:ScheduleKeyDeletion`
    - `kms:Revoke*`
    - `kms:Put*`
    - `kms:List*`
    - `kms:Get*`
    - `kms:Enable*`
    - `kms:Disable*`
    - `kms:Describe*`
    - `kms:Delete*`
    - `kms:Create*`
    - `kms:CancelKeyDeletion`

- **Resource** (`string`)
  - **Description:** The resources that the statement applies to.
  - **Default value:** `"*"`

### Statement 2

- **Sid** (`string`)
  - **Description:** The statement identifier for the second statement.
  - **Default value:** `"2"`

- **Effect** (`string`)
  - **Description:** The effect of the statement, which allows or denies the specified actions.
  - **Default value:** `"Allow"`

- **Principal** (`object`)
  - **Description:** The AWS principals (users) that the statement applies to.
  - **Properties:**
    - **AWS** (`list of strings`)
      - **Description:** List of AWS IAM user ARNs.
      - **Examples:**
        - `"arn:aws:iam::986602297069:user/jmemon@vivsoft.io"`
        - `"arn:aws:iam::986602297069:user/tkaza@vivsoft.io"`

- **Action** (`list of strings`)
  - **Description:** List of actions that are allowed.
  - **Default values:**
    - `kms:ReEncrypt*`
    - `kms:GenerateDataKey*`
    - `kms:Encrypt`
    - `kms:DescribeKey`
    - `kms:Decrypt`

- **Resource** (`string`)
  - **Description:** The resources that the statement applies to.
  - **Default value:** `"*"`

### Statement 3

- **Sid** (`string`)
  - **Description:** The statement identifier for the third statement.
  - **Default value:** `"3"`

- **Effect** (`string`)
  - **Description:** The effect of the statement, which allows or denies the specified actions.
  - **Default value:** `"Allow"`

- **Principal** (`object`)
  - **Description:** The AWS principals (roles) that the statement applies to.
  - **Properties:**
    - **AWS** (`list of strings`)
      - **Description:** List of AWS IAM role ARNs.
      - **Examples:**
        - `"arn:aws:iam::986602297069:role/demo-control-plane"`
        - `"arn:aws:iam::986602297069:role/demo-worker"`

- **Action** (`list of strings`)
  - **Description:** List of actions that are allowed.
  - **Default values:**
    - `kms:DescribeKey`
    - `kms:Decrypt`

- **Resource** (`string`)
  - **Description:** The resources that the statement applies to.
  - **Default value:** `"*"`

Once created and note down the ARN of the KMS key you created, and create the sops.yaml file in below format , changing the ` ADD_YOUR_KMS_KEY_ARN_HERE` with your actual `ARN`.

This file we will use whole deploying BigBang.

```
---
creation_rules:
  - kms: ADD_YOUR_KMS_KEY_ARN_HERE 
    encrypted_regex: "^(data|stringData)$"
```

### Variables of creation_rules

- **kms** (`string`)
  - **Description:** The Amazon Resource Name (ARN) of the KMS key to be used for encryption.
  - **Default value:** `"ADD_YOUR_KMS_KEY_ARN_HERE"`

- **encrypted_regex** (`string`)
  - **Description:** The regular expression pattern used to identify fields that should be encrypted.
  - **Default value:** `"^(data|stringData)$"`

## Create Kubernetes Cluster

It does not matter how you create the kubernetes cluster , but you should have access to kubeconfig file.

The cluster also have enough resources ( memory/cpu/) to run the bigbang nodes.

The cluster-api server should be publicly accessible so that public gitlab ci-cd can access it.

## Worker node instance Profile

All worker nodes in your cluster must have an instance profile , which have a policy to allowing the kms:decrypt and describe permission,

See, sample policy below , Change your KMS key ARN

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "",
            "Effect": "Allow",
            "Action": [
                "kms:DescribeKey",
                "kms:Decrypt"
            ],
            "Resource": "ADD_YOUR_KMS_KEY_ARN_HERE"
        }
    ]
}
```

## Create GPG Encryption Key

Generate the gpg key with name `bigbang-sops`

```
# Generate a GPG master key
# The GPG key fingerprint will be stored in the $fp variable
export fp=`gpg --quick-generate-key bigbang-sops rsa4096 encr | sed -e 's/ *//;2q;d;'`
gpg --quick-add-key ${fp} rsa4096 encr

echo ${fp}
```

# GPG Key Generation Script Variables

# `fp`

- **Type:** `string`
- **Description:** The fingerprint of the GPG master key.
- **How it's set:** This variable is set by running the `gpg --quick-generate-key` command and extracting the fingerprint using `sed`.
- **Default value:** N/A (set dynamically at runtime)
- **Example usage:**

  ```bash
  export fp=`gpg --quick-generate-key bigbang-sops rsa4096 encr | sed -e 's/ *//;2q;d;'`

Now create a secret in your cluster with SOPS private key for Big Bang to decrypt secrets at run time.

```
kubectl create namespace bigbang

gpg --export-secret-key --armor ${fp} | kubectl create secret generic sops-gpg -n bigbang --from-file=bigbangkey.asc=/dev/stdin

kubectl get secret -n bigbang sops-gpg
```
The sops value for your BigBang Deployment will be as below,

```
---
creation_rules:
- encrypted_regex: '^(data|stringData)$'
  pgp: EEF17D87C3954A2AE9D406811D17192D335BBD12
```

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

# The different types of components avalable are listed below: 

# settings Tool

- **Repo** -Domain used for BigBang created exposed services, can be overridden by individual packages.

# Service Mesh Tools

- **Istio** - Istio extends Kubernetes to establish a programmable, application-aware network using the powerful Envoy service proxy. Working with both Kubernetes and traditional workloads, Istio brings standard, universal traffic management, telemetry, and security to complex deployments.
- **Istio Operator** - Instead of manually installing, upgrading, and uninstalling Istio, you can instead let the Istio operator manage the installation for you. This relieves you of the burden of managing different istioctl versions. Simply update the operator custom resource (CR) and the operator controller will apply the corresponding configuration changes for you.
- **Jaeger** - Distributed tracing observability platforms, such as Jaeger, are essential for modern software applications that are architected as microservices. Jaeger maps the flow of requests and data as they traverse a distributed system. These requests may make calls to multiple services, which may introduce their own delays or errors. Jaeger connects the dots between these disparate components, helping to identify performance bottlenecks, troubleshoot errors, and improve overall application reliability. Jaeger is 100% open source, cloud native, and infinitely scalable.
- **Kiali** - Kiali is a console for Istio service mesh. Kiali can be quickly installed as an Istio add-on, or trusted as a part of your production environment.

# Security Tools

- **NeuVector** - NeuVector is the only 100% open source, Zero Trust container security platform. Continuously scan throughout the container lifecycle. Remove security roadblocks. Bake in security policies at the start to maximize developer agility.

- **Fortify** - Fortify is a comprehensive application security (AppSec) platform developed by Micro Focus. It empowers organizations to proactively identify and address vulnerabilities throughout the entire software development lifecycle (SDLC). Think of it as a security shield woven into the fabric of your development process, helping you build secure software from the ground up.

- **TwistLock** - Twistlock, now known as Palo Alto Networks Prisma Cloud, is a comprehensive cloud-native security platform designed to protect containerized applications and serverless workloads across cloud environments
- **Anchore** - Anchore is a container security and compliance platform that helps organizations discover, analyze, and enforce security and compliance policies for containerized applications and images. It ensures that container images are free from vulnerabilities and meet security and compliance standards before they are deployed.
- **Vault** - Vaults work by encrypting each secret to help prevent unauthorized users from gaining access. They function mostly as an active storage container for secrets as well as an account management system for dealing with multiple privileged accounts across the company.

# SSO Tools

- **Auth Service** - An authentication service is an identity verification mechanism—similar to passwords—for apps, websites, or software systems.
- **Keycloak** - Keycloak is the standalone tool for identity and access management, which allows us to create a user database with custom roles and groups. We can use this information further to authenticate users within our application and secure parts of it based on predefined roles.

# Policy Management Tools

- **Kyverno** - Kyverno is a policy engine designed for Kubernetes platform engineering teams. It enables security, automation, compliance, and governance using policy-as-code.
- **Kyverno Reporter** - Monitoring and Observability Tool for the PolicyReport CRD with an optional UI.
- **Kyverno Policies** - Kyverno policies can validate, mutate, generate, and cleanup Kubernetes resources, and verify image signatures and artifacts to help secure the software supply chain.
- **Cluster Auditor** - Cluster Auditor is a tool that pulls constraints from the Kubernetes API, transforms them, and inserts them into Prometheus to be displayed in a Grafana Dashboard.
- **Gatekeeper** - Gatekeeper HQ is a vendor and contract lifecycle management platform designed for companies of all sizes. It offers a secure contract and vendor repository that stores every file, interaction, and piece of metadata relating to all your contract agreements. The platform provides features for vendor management, contract management, Kanban workflow, collaboration, and reporting, with the ability to extend functionality through additional modules and integration with over 220 third-party solutions.

# Logging Tools

- **ECK operator** - The ECK operator, or Elastic Cloud on Kubernetes, is built on the Kubernetes Operator pattern. It extends Kubernetes’ orchestration capabilities to support the setup and management of various Elastic Stack components on Kubernetes, including Elasticsearch, Kibana, APM Server, Enterprise Search, Beats, Elastic Agent, Elastic Maps Server, and Logstash
- **Fluentbit** - Fluent Bit is a super fast, lightweight, and highly scalable logging and metrics processor and forwarder. It is the preferred choice for cloud and containerized environments.
- **EFK stack** -The EFK stack is a popular open-source logging solution for Kubernetes environments. It stands for:
    - **Elasticsearch**: A search and analytics engine.
    - **Fluentd**: An open-source data collector for unified logging layer.
    - **Kibana**: A visualization dashboard for Elasticsearch.
- **Loki** - Grafana Loki is a set of open source components that can be composed into a fully featured logging stack. A small index and highly compressed chunks simplifies the operation and significantly lowers the cost of Loki.
- **Promtail** - Promtail is an agent which ships the contents of local logs to a private Grafana Loki instance or Grafana Cloud. It is usually deployed to every machine that runs applications which need to be monitored.

# Monitoring Tools

- **Monitoring** - Together, these tools provide a powerful platform for aggregating, analyzing, and visualizing logs in Kubernetes clusters, helping with monitoring, troubleshooting, and gaining insights from applications.
- **Thanos** - Thanos is based on Prometheus. With Thanos, Prometheus always remains as an integral foundation for collecting metrics and alerting using local data.
- **Tempo** -  Tempo is an open source, easy-to-use, and high-scale distributed tracing backend.
- **Sonarqube** - SonarQube is a self-managed, automatic code review tool that systematically helps you deliver Clean Code. As a core element of our Sonar solution.
- **Tempo** -  Tempo is an open source, easy-to-use, and high-scale distributed tracing backend.
- **Grafana** - Grafana is the open source analytics & monitoring solution for every database.
- **Metrics Server** - The Kubernetes Metrics Server is an aggregator of resource usage data in your cluster, and it isn't deployed by default in Amazon EKS clusters. 

# Development Tools

- **Sonarqube** - SonarQube is a self-managed, automatic code review tool that systematically helps you deliver Clean Code. As a core element of our Sonar solution.

- **Nexus** -  Nexus is a robust tool designed for managing and organizing artifacts throughout the software development lifecycle. 
- **Harbor** - Harbor is an open source registry that secures artifacts with policies and role-based access control, ensures images are scanned and free from vulnerabilities, and signs images as trusted. Harbor, a CNCF Graduated project, delivers compliance, performance, and interoperability to help you consistently and securely manage artifacts across cloud native compute platforms like Kubernetes and Docker.
- **GitLab**- GitLab is a web-based DevOps lifecycle tool that integrates various stages of the software development process into a single platform.

# CI/CD Tools

- **ArgoCD** - Argo CD is a declarative, GitOps continuous delivery tool for Kubernetes.
- **GitLab Runner** - GitLab Runner is an application that works with GitLab CI/CD to run jobs in a pipeline.

# Proxy Tool

- **Haproxy** - HAProxy is a free, open source high availability solution, providing load balancing and proxying for TCP and HTTP-based applications by spreading requests across multiple servers.

# Backup Tool:

- **Velero** - Velero is an open source tool to safely backup and restore, perform disaster recovery, and migrate Kubernetes cluster resources and persistent volumes.

# Object Storage Tools:

- **Minio** - MinIO is a high-performance, S3 compatible object store. It is built for
large scale AI/ML, data lake and database workloads. It is software-defined
and runs on any cloud or on-premises infrastructure. MinIO is dual-licensed
under open source GNU AGPL v3 and a commercial enterprise license.
- **Minio Operator** - The MinIO Operator installs a Custom Resource Definition (CRD) to support describing MinIO tenants as a Kubernetes object. See the MinIO Operator CRD Reference for complete documentation on the MinIO CRD.

# Communication Tools

- **Mattermost** - Mattermost is a secure collaboration platform for accelerating mission critical work in complex environments.
- **Metrics Server** - The Kubernetes Metrics Server is an aggregator of resource usage data in your cluster, and it isn't deployed by default in Amazon EKS clusters. 

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