---
title: "Catalog Item Configuration"
description: "Various configuration details for configuring the ENBUILD Catalog Items."
summary: ""
date: 2023-09-07T16:04:48+02:00
lastmod: 2023-09-07T16:04:48+02:00
draft: false
menu:
  docs:
    parent: "docs"
    identifier: "enbuildCatalogItemManifestSchema"
weight: 8
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

The detailed configuration for each Catalog Item resides in separate JSON files, such as "aws_landing_zone.json." These files define crucial parameters for deployment using tools like Terraform or Helm.

### AWS Landing Zone Example:

```json
{
  "type": "terraform",
  "slug": "aws_landing_zone",
  // ... (other fields)
  "description": "The Landing Zone (IL4) module deploys an AWS infrastructure following AWS recommended best practices for CMS cloud.",
  "components": [
    // ... (list of components)
  ],
  "infrastructure": {
    // ... (infrastructure details)
  }
}
```

This snippet showcases essential details, including the deployment type, repository information, description, and specific components.

### Big Bang Deployment Example:

For more complex deployments, such as Platform One's Big Bang, a Helm-based deployment, the manifest file (`bigbang.json`) orchestrates a multitude of components within a Kubernetes cluster. This exemplifies the versatility of ENBUILD in managing diverse and intricate deployment scenarios.

```json
{
  "type": "helm",
  "slug": "bigbang",
  // ... (other fields)
  "description": "The BigBang module deploys standard BigBang components on top of an existing Kubernetes cluster. Ver 1.50",
  "components": [
    // ... (list of components)
  ],
  "infrastructure": {
    // ... (infrastructure details)
  },
  "configuration": [
    // ... (additional configuration details)
  ]
}
```

In summary, the ENBUILD Manifest plays a pivotal role in configuring and managing the diverse array of software deployments offered through the ENBUILD platform. Developers can leverage this structured approach to streamline their workflows, ensuring consistency, version control, and efficient collaboration.

### Catalog Item Manifest JSON Key Value Pairs:

| Name                | Description                                                                                                                                             | Example Input(s)                                         |
| :------------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------ | :------------------------------------------------------- |
| type                | The type of Catalog Item                                                                                                                                | terraform helm                                           |
| slug                | Slug of the Catalog Item                                                                                                                                |                                                          |
| name                | Name of the Catalog Item                                                                                                                                | AWS Landing Zone Big Bang                                |
| repository          | The reference URL of the template repository (what ENBUILD creates a project / repository from)                                                         | https://gitlab.com/enbuild-staging/iac-templates/bigbang |
| project_id          | This is specific to GitLab the project id of the repository in GitLab                                                                                   | 202                                                      |
| readme_file_path    | The path of the README file inside the template repository. Note: This will be displayed on the ENBUILD UI when the user clicks the Information button. | /docs/README.md                                          |
| values_folder_path  | The location to save the updated values                                                                                                                 |                                                          |
| secrets_folder_path | The location to save the updated secrets                                                                                                                |                                                          |
| ref                 | The branch of the template repository to build from                                                                                                     | main                                                     |
| sops                | Enable or disable Mozilla SOPS for the Catalog Item. If false the secret tab is not displayed                                                           | True                                                     |
| image_path          | The path of the image to display as the ENBUILD Catalog Item logo. This image needs to be available within the template project repository.             | /images/AKS.jpeg                                         |
| components          | What components the Catalog Item has.                                                                                                                   | See Table Below                                          |
| infrastructure      | What infrastructure is required and compatible with the Catalog Item.                                                                                   | See Table Below                                          |

### Component JSON Key Value Pairs:

| Name                | Description                                                                                                                                             | Example Input(s)                                         |
| :------------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------ | :------------------------------------------------------- |
| type                | The type of Catalog Item                                                                                                                                | terraform helm                                           |
| slug                | Slug of the Catalog Item                                                                                                                                |                                                          |
| name                | Name of the Catalog Item                                                                                                                                | AWS Landing Zone Big Bang                                |
| repository          | The reference URL of the template repository (what ENBUILD creates a project / repository from)                                                         | https://gitlab.com/enbuild-staging/iac-templates/bigbang |
| project_id          | This is specific to GitLab the project id of the repository in GitLab                                                                                   | 202                                                      |
| readme_file_path    | The path of the README file inside the template repository. Note: This will be displayed on the ENBUILD UI when the user clicks the Information button. | /docs/README.md                                          |
| values_folder_path  | The location to save the updated values                                                                                                                 |                                                          |
| secrets_folder_path | The location to save the updated secrets                                                                                                                |                                                          |
| ref                 | The branch of the template repository to build from                                                                                                     | main                                                     |
| image_path          | The path of the image to display as the ENBUILD Catalog Item logo. This image needs to be available within the template project repository.             | /images/AKS.jpeg                                         |
| tool_type           | The type of tool used to group common tools on the ENBUILD UI                                                                                           | Security                                                 |

### Infrastructure JSON Key Value Pairs:

| Name           | Description                      | Example Input(s)         |
| :------------- | :------------------------------- | :----------------------- |
| type           | The type of Catalog Item         | terraform helm           |
| slug           | Slug of the Catalog Item         |                          |
| name           | Name of the Infrastructure Input | AWS Service Principal ID |
| mandatory      | Is this input optional?          | true                     |
| showKubeConfig | ???                              | ???                      |
| selections     | ???                              | ???                      |
