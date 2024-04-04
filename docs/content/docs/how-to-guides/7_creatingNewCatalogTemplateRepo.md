---
title: "Creting new catalog template repository"
description: "Steps to create a new catalog template repository in your Github/Gitlab account"
summary: "Creting new catalog template repository"
draft: false
menu:
  docs:
    parent: "docs/how-to-guides/deploying-enbuild-for-local-testing/"
    identifier: "newTemplate"
weight: 207
toc: true
---
## Overview

This guide outlines the steps required to create a new catalog template repository in your GitHub or Gitlab account. 


## Prerequisites

- You must have a [GitHub](https://github.com/) [Gitlab](https://gitlab.com/) account and have access to the orgnizations you are part of.
- Familiarity with [Git](https://git-scm.com/doc) and [Github Actions](https://docs.github.com/en/actions) is recommended but not required.

## Steps

1. **Create a new repository**

   Navigate to your GitHub dashboard and click on the '+' button in the top left corner, then select `New repository`. Enter a name for your catalog template repository initialize it with a README file, and make it public if you want your catalog templates to be publicly available.

2. **Add the images and infrastructure files**

   Add the following files to your repository:

   - `images/<catalog-name>.png` -- This is the image ENBUILD will disaply while showing the catalog in the UI.
   - `images/<component-name>.png` -- You have to create images for all the compoments for your catalog. This is the image ENBUILD will disaply while showing the component of catalog in the UI.
   - `infra/` , all your Terraform infrastructure files should go here
   - `scripts/`, all your scripts should go here
   - `values.yaml` -- This file will contain all the default values for your templates.

3. **Add the CI-CD Pipeline files**

If its a GitHub repo:
   - Create a file `.github/workflows/main.yml` and add the deployment steps to it.
   - Disable the github actions for this repo. Since this is a template repository, you don't want your templates being deployed when someone adds a commit to this repo.

If its a Gitlab repo:
   - Create a file `gitlab-ci.yml` and add the deployment steps to it

:zap: **Note:** The Gitlab CI expects the CI file name to be `.gitlab-ci.yml`, but since the repository we are operating on is a template repository and we do not want to run our pipeline locally for this Repository, hence the name of the file is changed to `gitlab-ci.yml` and ENBUILD will rename this to 
`.gitlab-ci.yml` in the target deployment repository of this catalog.

## Reference Catalog Repository

   See the [example Github catalog template repository](https://github.com/VivSoftOrg/iac-templates-cypress-test-catalog/)
   and  
   [example Gitlab catalog template repository](https://gitlab.com/enbuild-staging/iac-templates/cypress-test-catalog) 
   
   for an example of what the template repository should look like.



## Next Steps

- [Creating a new catalog entry in ENBUILD](https://enbuild-docs.vivplatform.io/docs/how-to-guides/adding-new-catalog-item-in-enbuild/)

