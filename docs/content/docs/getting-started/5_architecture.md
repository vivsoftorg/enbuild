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

<picture><img src="/images/getting-started/architecture.svg" alt="Screenshot of ENBUILD Login Screen"></img></picture>

### Frontend Service

The ENBUILD frontend service provides the ENBUILD User Interface (UI) to the end user.

### Backend Service

The ENBUILD backend service is responsible for integration with the CI/CD Provider.

### User Service

The ENBUILD user service manages the end-user's state, such as authentication, access, API access and role-based access control.

### ML Service

The ENBUILD ML (Machine Learning) service enables data scientists to quickly create feature sets and deploy models. An instance of Jupyter Notebook can also be created and accessed from this service. 
***(This is a placeholder service for demo purposes for now and will be implemented in the future)***

### Request Service

The ENBUILD Request service is demo service to enable linking multiple catalog items to one another and deploy them together.
***(This is a placeholder service for demo purposes for now and will be implemented in the future)***

### RabbitMQ Consumer Service

The ENBUILD RabbitMQ consumer service processes jobs in the work queue as well as integrates with the CI/CD Provider APIs.

### NoSQL Database

ENBUILD utilizes a NoSQL database to manage the application’s state across all microservices. ENBUILD provides a MongoDB instance out-of-the-box, but also can integrate with Cloud Service Provider NoSQL Databases such as Azure’s CosmosDB and AWS’ DocumentDB.

### Identity and Access Management

ENBUILD supports integration with Okta and Keycloak. Keycloak can act as an Identity Broker for other IdAM products such as Active Directory.
