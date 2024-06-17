---
title: "Deploying Kubernetes"
description: "Steps to Deploy Kubernetes "
summary: "Steps to Deploy Kubernets"
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

# Kubernetes

Kubernetes is an open-source container orchestration system for automating software deployment, scaling, and management.Originally designed by Google, the project is now maintained by a worldwide community of contributors, and the trademark is held by the Cloud Native Computing Foundation.

# Deploy Kubernetes


-  Login to Enbuild -[Enbuild](https://enbuild.vivplatform.io)
-  Click on the **Home**
-  Select the **Kubernetes**
-  Choose the component **RKE2** from the **DISTRO** category and click on the **VALUES** tab and provide the `Credentials`

<picture><img src="/images/how-to-guides/RKE2.png" alt="RKE2"></img></picture>

## RKE2

RKE2, also known as RKE Government, is Rancher's next-generation Kubernetes distribution.

It is a fully conformant Kubernetes distribution that focuses on security and compliance within the U.S. Federal Government sector.

## Deploy RKE2

# VPC Configuration

- **create_vpc** (`true` or `false`)
  - **Description:** The **create_vpc** variable decides if a new Virtual Private Cloud (VPC) should be created.
  - **Default value:** `false`

- **vpc_cidr** (`string`)
  - **Description:** The **vpc_cidr** variable defines the IP address range for the VPC.
  - **Default value:** `"10.0.0.0/16"`

  # NAT Gateway Configuration

- **enable_nat_gateway** (`true` or `false`)
  - **Descriprion:** The **enable_nat_gateway** variable decides if a NAT Gateway should be enabled.
  - **Default value:** `true`

- **single_nat_gateway** (`true` or `false`)
  - **Description:**  The **single_nat_gateway** decides if only one NAT Gateway should be used.
  - **Default value:** `true`

- **vpc_id** (`string`)
  - **Description:**  The **vpc_id** is used to  Specify an existing VPC to use. Needed if `create_vpc` is `false`.
  - **Default value::** `"vpc-39b8da44"`

- **subnets** (`list of strings`)
  - **Description:** Lists the IDs of subnets within the VPC. Needed if `create_vpc` is `false`.
  - **Example:** `["subnet-5817463e", "subnet-f191cdd0"]`

# EC2 Instance Configuration

- **instance_type** (`string`)
  - **Description:**  The **instance_type**  specifies the type of EC2 instance to use.
  - **Default value:** `"t3.large"`

- **associate_public_ip_address** (`true` or `false`)
  - **Description:** The **associate_public_ip_address** decides if the instance should have a public IP address.
  - **Default value:** `true`

- **controlplane_internal** (`true` or `false`)
  - **What it does:** Decides if the control plane should be internal only.
  - **Default value:** `false`

- **servers** (`number`)
  - **Description:** Number of EC2 instances to create.
  - **Default value:** `1`

# Auto Scaling Group (ASG) Configuration

- **asg** (`object`)
  - **Description:**  The variable **asg**is used to the Auto Scaling Group (ASG).
  - **Properties:**
    - `min` (`number`): Minimum number of instances in the ASG.
      - **Default value:** `1`
    - `max` (`number`): Maximum number of instances in the ASG.
      - **Default value:** `10`
    - `desired` (`number`): Desired number of instances in the ASG.
      - **Default value:** `1`
    - `suspended_processes` (`list of strings`): List of processes to suspend.
      - **Default value:** `[]`
    - `termination_policies` (`list of strings`): List of termination policies.
      - **Default value:** `["Default"]`

# Block Device Mapping

- **block_device_mappings** (`object`)
  - **What it does:** Configuration for the block device (storage).
  - **Properties:**
    - `size` (`number`): Size of the volume in GB.
      - **Default value:** `50`
    - `type` (`string`): Type of the volume.
      - **Default value:** `"gp2"`

# Registry Mirror Configuration

- **create_registry1_mirror** (`true` or `false`)
  - **Description:**  **create_registry1_mirror** decides if a mirror for the `https://registry1.dso.mil` container registry should be created.
  - **Default value:** `false`

- **registry1_mirror_proxy_address** (`string`)
  - **Description:** **registry1_mirror_proxy_address**  variable is used to declare the proxy address for the registry1 mirror.
  - **Example:** `"http://44.210.192.97:5000"`

- After providing all the input values, provide the name for your deployment proceed to Infrastructure section,
- Select AWS as your cloud and provide your AWS credentials.

<picture><img src="/images/how-to-guides/AWS_Kubernetes.png" alt="AWS_Kubernetes"></img></picture>

Once all inputs are provided click on **Deploy Stack**

- Choose the component **EKS** from the **DISTRO** category and follow the previous steps.

## EKS

Amazon Elastic Kubernetes Service (Amazon EKS) is a managed Kubernetes service to run Kubernetes in the AWS cloud and on-premises data centers. In the cloud, Amazon EKS automatically manages the availability and scalability of the Kubernetes control plane nodes responsible for scheduling containers, managing application availability, storing cluster data, and other key tasks. With Amazon EKS, you can take advantage of all the performance, scale, reliability, and availability of AWS infrastructure, as well as integrations with AWS networking and security services. On-premises, EKS provides a consistent, fully-supported Kubernetes solution with integrated tooling and simple deployment to AWS Outposts, virtual machines, or bare metal servers.

## Deploy EKS

# VPC Configuration

- **create_vpc** (`true` or `false`)
  - **Description:** The **create_vpc** variable decides if a new Virtual Private Cloud (VPC) should be created.
  - **Default value:** `true`

- **vpc_cidr** (`string`)
  - **Description:** The **vpc_cidr** variable defines the IP address range for the VPC.
  - **Default value:** `"10.0.0.0/16"`

  If you don't want to create a new VPC, set `create_vpc` to `false` and provide the following variables:
- **vpc_id** (`string`)
  - **Description:** The **vpc_id** specifies an existing VPC to use. Needed if `create_vpc` is `false`.
  - **Example:** `"vpc-39b8da44"`

- **subnet_ids** (`list of strings`)
  - **Description:** Lists the IDs of subnets within the VPC. Needed if `create_vpc` is `false`.
  - **Example:** `["subnet-1242491c", "subnet-5817463e"]`

# EKS Cluster Configuration

- **cluster_name** (`string`)
  - **Description:** The **cluster_name** variable specifies the name of the EKS cluster.
  - **Default value:** `"juned-eks"`

- **cluster_version** (`string`)
  - **Description:** The **cluster_version** variable specifies the version of the EKS cluster.
  - **Default value:** `"1.29"`

- **cluster_endpoint_public_access** (`true` or `false`)
  - **Description:** The **cluster_endpoint_public_access** variable decides if the EKS cluster endpoint should be publicly accessible.
  - **Default value:** `true`

- **cluster_endpoint_private_access** (`true` or `false`)
  - **Description:** The **cluster_endpoint_private_access** variable decides if the EKS cluster endpoint should be privately accessible.
  - **Default value:** `false`

# EKS Node Groups Configuration

- **eks_node_groups_min_size** (`number`)
  - **Description:** The **eks_node_groups_min_size** variable specifies the minimum number of nodes in the EKS node group.
  - **Default value:** `1`

- **eks_node_groups_max_size** (`number`)
  - **Description:** The **eks_node_groups_max_size** variable specifies the maximum number of nodes in the EKS node group.
  - **Default value:** `5`

- **eks_node_groups_desired_size** (`number`)
  - **Description:** The **eks_node_groups_desired_size** variable specifies the desired number of nodes in the EKS node group.
  - **Default value:** `1`

# NAT Gateway Configuration

- **enable_nat_gateway** (`true` or `false`)
  - **Description:** The **enable_nat_gateway** variable decides if a NAT Gateway should be enabled.
  - **Default value:** `true`

- **single_nat_gateway** (`true` or `false`)
  - **Description:** The **single_nat_gateway** variable decides if only one NAT Gateway should be used.
  - **Default value:** `true`

# Kubernetes Configuration

- **create_kubeconfig** (`true` or `false`)
  - **Description:** The **create_kubeconfig** variable decides if a kubeconfig file should be created for the EKS cluster.
  - **Default value:** `true`

# EC2 Instance Configuration

- **instance_types** (`list of strings`)
  - **Description:** The **instance_types** variable specifies the types of EC2 instances to use for the EKS nodes.
  - **Default value:** `["t3.large"]`

# Registry Mirror Configuration

- **create_registry1_mirror** (`true` or `false`)
  - **Description:** The **create_registry1_mirror** variable decides if a mirror for the `https://registry1.dso.mil` container registry should be created.
  - **Default value:** `false`

- **registry1_mirror_proxy_address** (`string`)
  - **Description:** The **registry1_mirror_proxy_address** variable specifies the proxy address for the registry1 mirror.
  - **Example:** `"http://44.210.192.97:5000"`


<picture><img src="/images/how-to-guides/EKS.png" alt="EKS"></img></picture>