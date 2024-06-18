---
title: "Exporting Catalog Item from ENBUILD"
description: "Exporting Catalog Item from ENBUILD"
summary: "Exporting Catalog Item from ENBUILD"
draft: false
menu:
  docs:
    parent: "docs/how-to-guides/deploying-enbuild-for-local-testing/"
    identifier: "exportCatalog"
weight: 210
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

You can export a catalog item from ENBUILD in JSON format. You can use this JSON to commit in the version control system for a backup. 
And then use this json to import to another ENBUILD instance.


1. Login to ENBUILD as admin user.
2. Navigate to **Manage Catalogs** tab and you will see the list of all active catalogs in the ENBUILD instance.
<picture><img src="/images/how-to-guides/exportCatalog.png" alt="CatalogList"></img></picture>
3. Select the Catalog that you want to export and then at the **Actions** coloumn select the export icon.