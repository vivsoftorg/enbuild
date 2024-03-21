---
title: "Catalog Manifest Schema"
description: "Various configuration details for configuring the ENBUILD Manifest."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
draft: false
menu:
  docs:
    parent: "docs"
    identifier: "enbuildManifestSchema"
weight: 7
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

The ENBUILD Manifest serves as the backbone for configuring both ENBUILD and its Catalog Items, ensuring a smooth and efficient software development and deployment process. This configuration file, named `manifest.json`, is stored in a GitLab or GitHub repository, aligning with version control best practices and supporting a GitOps approach.

### Overview of manifest.json

The manifest file, devoid of the ".json" extension reference, organizes each Catalog Item that appears on the ENBUILD UI. It provides a structured view of the available items, facilitating easy management and customization.

```json
{
  "catalogs": [
    {
      "id": 1,
      "file_name": "aws_landing_zone"
    }
    //Additional Catalog Items...
  ]
}
```
