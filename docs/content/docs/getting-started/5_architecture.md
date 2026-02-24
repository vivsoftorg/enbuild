---
title: "ENBUILD Architecture"
description: "Architecture Diagram of ENBUILD"
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
draft: false
menu:
  docs:
    parent: ""
    identifier: "architecture"
weight: 105
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

## Architecture Digram

<picture><img src="/images/getting-started/enbuild-architecure.png" alt="Screenshot of ENBUILD Architecture"></img></picture>

### Frontend Service

The ENBUILD frontend service provides the ENBUILD User Interface (UI) to the end user.

### Backend Service

The ENBUILD backend service is responsible for integration with the CI/CD Provider.

### User Service

The ENBUILD user service manages the end-user's state, such as authentication, access, API access and role-based access control.

### ML Service

The ENBUILD ML (Machine Learning) service enables data scientists to quickly create feature sets and deploy models. An instance of Jupyter Notebook can also be created and accessed from this service.
***(This is a placeholder service for demo purposes for now and will be implemented in the future)***

### RabbitMQ Consumer Service

The ENBUILD RabbitMQ consumer service processes jobs in the work queue as well as integrates with the CI/CD Provider APIs.

### NoSQL Database

ENBUILD utilizes a NoSQL database to manage the application’s state across all microservices. ENBUILD provides a MongoDB instance out-of-the-box, but also can integrate with Cloud Service Provider NoSQL Databases such as Azure’s CosmosDB and AWS’ DocumentDB.

### Identity and Access Management

ENBUILD supports integration with Okta and Keycloak. Keycloak can act as an Identity Broker for other IdAM products such as Active Directory.

### CTF Backend Service

The ENBUILD CTF (Capture The Flag) Backend Service is part of the **Secure Lightning** feature set. It provides a comprehensive platform for hosting cybersecurity training competitions and challenges.

**Key Features:**
- countermeasure management (create, update, delete countermeasures)
- Repository management for projects.
- Real-time scoring and leaderboard
- RESTful API for frontend integration

**Technical Details:**
- Runs on container port 8000
- Uses MongoDB for persistent storage
- Configurable CORS origins for frontend integration
- Supports AWS regions for cloud-native deployments
- Can utilize node IAM roles for AWS access (hostNetwork mode)
