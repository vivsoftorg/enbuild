---
title: "How ENBUILD works"
description: "Key details about ENBUILD."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
draft: false
menu:
  docs:
    parent: "docs"
    identifier: "howEnbuildWorks"
weight: 104
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

ENBUILD operates as a bridge between developers and the Continuous Integration/Continuous Deployment (CI/CD) systems, specifically with popular platforms like GitHub and GitLab. This tool simplifies the entire software development process, making it accessible and efficient for users with various technical backgrounds.

### GitHub Integration

When working with GitHub, ENBUILD seamlessly communicates with the CI/CD provider through REST and GraphQL APIs. It starts by creating a Repository Project, essentially a structured workspace, using a template repository project. The user interacts with ENBUILD through an intuitive user interface (UI), providing inputs that guide the customization of files. These files are then updated based on the user's preferences. Once the configuration is complete, ENBUILD executes the predefined workflow or pipeline, automating the deployment process without the need for intricate manual steps.

<picture><img src="/images/getting-started/gitHubTemplateConfig.png" alt="Screenshot of GitHub Repo Configuration for Template Repository"></img></picture>

### GitLab Integration

Similarly, when integrating with GitLab, ENBUILD initiates the process by creating a Repository Project within the GitLab environment. Just like with GitHub, the user's inputs through the ENBUILD UI guide the customization of files within the project. ENBUILD updates these files accordingly, ensuring that the project aligns with the user's specifications. Subsequently, the tool triggers the execution of the workflow or pipeline, seamlessly automating the development and deployment steps.

In essence, ENBUILD acts as a facilitator, taking the complexities out of setting up projects and workflows. It empowers users to customize their development environments through a straightforward UI, enabling a smooth and efficient interaction with CI/CD systems like GitHub and GitLab. The result is an accelerated and simplified software development process, allowing developers to focus on building innovative solutions without getting bogged down by technical intricacies.
