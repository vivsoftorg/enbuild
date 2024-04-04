---
title: "Frequently Asked Questions"
description: "Frequently Asked Questions"
summary: ""
draft: false
menu:
  docs:
    parent: "/faq/"
    identifier: "faq"
weight: 401
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

<picture><img src="/images/emma/emma-green-hoodie.png" alt="Screenshot of Emma" width="20%" height="20%"></img></picture>
<br />
<br />

### About

> **Q: Who is behind this project?**

**A:** [VivSoft](https://www.vivsoft.io)! :elephant:<br/><br/>
VivSoft is focused on solving complex problems in the public sector using innovative technologies. VivSoft is working with business leaders in federal, state and local government to help mission owners accelerate innovation using DevSecOps, Cloud, AI/ML and Blockchain Technologies.<br/><br/>

> **Q: What problem does ENBUILD solve?**

**A:** ENBUILD's goal is to empower organizations to leverage Kubernetes and cloud-native technologies effectively while minimizing the complexity and overhead associated with deployment and management tasks. By offering pre-configured catalog items and promoting best practices, ENBUILD enables organizations to streamline their development workflows, reduce time to market, and achieve their product development goals more efficiently.

### Costs and Licensing Fees

> **Q: Is ENBUILD free to use?**

**A:** Yes, ENBUILD is free to use.

> **Q: What license does ENBUILD use?**

**A:** ENBUILD utilizes the Apache License, which is an open-source software license recognized for its flexibility and permissive nature. This license allows users to freely use, modify, and distribute the software, whether for personal, commercial, or open-source projects. <br/><br/>
For more detailed information about the Apache License and its implications, please refer to the official Apache Software Foundation website or consult the license file included with the ENBUILD software distribution.

### Security

> **Q: What dependencies does ENBUILD have?**

**A:** ENBUILD, being a Kubernetes native application, relies on certain dependencies to function optimally.
<br/><br/>These dependencies include:

- **Kubernetes Cluster**: ENBUILD requires a Kubernetes cluster to operate. This ensures that ENBUILD can leverage the scalability, resilience, and orchestration capabilities provided by Kubernetes, thereby enabling efficient deployment and management of containerized applications.

**Version Control System** (VCS) for ENBUILD Catalog Item Deployments: ENBUILD utilizes a Version Control System for managing and deploying ENBUILD Catalog Items. Currently, ENBUILD supports integration with two popular VCS platforms:

- **GitLab**: ENBUILD seamlessly integrates with GitLab, allowing users to leverage the robust features of GitLab for managing and versioning their ENBUILD Catalog Items. This includes support for both the Software as a Service (SaaS) offering of GitLab and self-hosted deployments of GitLab instances.

- **GitHub**: ENBUILD also supports integration with GitHub, enabling users to utilize GitHub's collaborative features and version control capabilities for ENBUILD Catalog Item deployments. Similar to GitLab, ENBUILD supports both the SaaS offering of GitHub and self-hosted deployments of GitHub Enterprise instances.

By supporting these Version Control Systems, ENBUILD provides users with flexibility and choice, allowing them to seamlessly integrate ENBUILD into their existing development workflows while leveraging the capabilities of their preferred VCS platform.

For further details on setting up ENBUILD dependencies and integration with specific Version Control Systems, please refer to the ENBUILD documentation or reach out to our support team for assistance.

### Deployment

> **Q: What types of Kubernetes distributions is ENBUILD compatible with?**

**A:** ENBUILD has been deployed and tested on the following distributions of Kubernetes:

- **Amazon EKS** (Elastic Kubernetes Service): ENBUILD is fully compatible with Amazon EKS, allowing users to deploy and manage ENBUILD on Amazon Web Services' managed Kubernetes service.
- **Azure AKS** (Azure Kubernetes Service): ENBUILD seamlessly integrates with Azure AKS, enabling users to deploy and run ENBUILD on Microsoft Azure's managed Kubernetes service.
- **Rancher K3D**: ENBUILD supports deployment on Rancher K3D, a lightweight Kubernetes distribution designed for local development and testing purposes.
- **Rancher RKE2**: ENBUILD is compatible with Rancher RKE2, an enterprise-grade Kubernetes distribution optimized for production workloads.

### Configuration

> **Q: What is an ENBUILD Catalog Item (CI)?**

**A:** An ENBUILD Catalog Item, often abbreviated as CI, serves as a standardized template project designed to streamline the deployment and management of various infrastructure components and applications within the ENBUILD ecosystem. These Catalog Items are meticulously crafted by developers to encapsulate pre-configured settings, best practices, and reusable components, providing users with a simplified and consistent approach to deploying complex infrastructure and applications.

ENBUILD Catalog Items are tailored to support a diverse range of use cases and technologies, including:

- **Terraform Infrastructure as Code** (IaC): Catalog Items for Terraform enable users to define and manage cloud infrastructure resources using Infrastructure as Code principles. These Catalog Items offer pre-defined templates and configurations for provisioning resources on popular cloud platforms such as AWS, Azure, and Google Cloud Platform, facilitating rapid deployment and automation of infrastructure provisioning tasks.

- **Kubernetes Distributions** (e.g., Amazon EKS, Azure AKS): Catalog Items for Kubernetes distributions provide users with pre-configured templates for deploying and managing Kubernetes clusters on various cloud providers, such as Amazon EKS (Elastic Kubernetes Service) and Azure AKS (Azure Kubernetes Service). These Catalog Items simplify the setup and configuration of Kubernetes clusters, enabling users to leverage the scalability and agility of Kubernetes for containerized application deployment and orchestration.

- **Helm Deployments** (e.g., Big Bang): Catalog Items for Helm deployments offer pre-packaged configurations and charts for deploying applications and services using Helm, a popular package manager for Kubernetes. These Catalog Items streamline the deployment of complex applications by providing ready-to-use Helm charts and configurations, reducing the time and effort required for setting up and configuring application environments.

By leveraging ENBUILD Catalog Items, developers gain access to a curated library of templates and configurations that serve as starting points for their projects. These Catalog Items not only accelerate the development and deployment process but also promote consistency, reliability, and best practices across projects within the ENBUILD ecosystem.

> **Q: What is an Version Control System (VCS)?**

**A:** A Version Control System (VCS) is a fundamental tool used in software development to manage and track changes to source code, configuration files, and other project assets over time. It allows multiple developers to collaborate on a project concurrently while maintaining a history of all modifications made to the project files.

In the context of ENBUILD, a Version Control System (VCS) plays a crucial role in managing and deploying ENBUILD Catalog Items, which are template projects designed to facilitate the deployment of infrastructure components and applications within the ENBUILD ecosystem. ENBUILD supports integration with popular VCS platforms such as GitLab and GitHub, providing users with the flexibility to version control their Catalog Items and streamline the deployment process.

Key features and benefits of using a Version Control System (VCS) include:

- **History Tracking**: VCS systems maintain a detailed history of changes made to project files, including who made the changes and when they were made. This allows developers to review past revisions, track the evolution of the project, and revert to previous versions if necessary.

- **Collaboration**: VCS systems enable seamless collaboration among team members by providing mechanisms for sharing and synchronizing changes to project files. Multiple developers can work on the same codebase simultaneously without risking conflicts or data loss.

- **Branching and Merging**: VCS systems support branching and merging workflows, allowing developers to create separate branches to work on specific features or fixes independently. Branches can later be merged back into the main codebase, ensuring a streamlined and organized development process.

- **Auditing and Compliance**: VCS systems offer auditing capabilities that help maintain compliance with regulatory requirements and internal policies. Organizations can track and monitor all changes made to project files, ensuring accountability and transparency in the development process.

By leveraging a Version Control System (VCS) such as GitLab or GitHub, ENBUILD users can effectively manage their Catalog Items, collaborate with team members, track changes, and maintain a consistent and reliable deployment workflow.

### Change Control

> **Q: What is the Release schedule?**

**A:** To Be Determined. :calendar:
