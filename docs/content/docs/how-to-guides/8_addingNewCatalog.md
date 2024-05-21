---
title: "Adding new Catalog Item in ENBUILD"
description: "Adding new Catalog Item in ENBUILD"
summary: "Adding new Catalog Item in ENBUILD"
draft: false
menu:
  docs:
    parent: "docs/how-to-guides/deploying-enbuild-for-local-testing/"
    identifier: "newCatalog"
weight: 208
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

# Prerequisites
Before you begin, ensure that you have the following prerequisites in place:
- ENBUILD is [installed](https://enbuild-docs.vivplatform.io/docs/how-to-guides/deploying-enbuild-for-local-testing/) and you have `admin` [access](https://enbuild-docs.vivplatform.io/docs/how-to-guides/configuring-enbuild/#set-the-admin-password) to it. 
- You have created the Catalog templates repository for the catalog item (either in Gitlab/Github) 
- If the catalog template repository is private, you have aquired the token with readonly access to the repository.

# Adding a new repository

Navigate to the **Repositories** tab of the **Admin** panel and click on **Connect New Respository**. You will be presented with the following form:

<picture><img src="/images/how-to-guides/createCatalog.png" alt="Add New Catalog"></img></picture>

## Choose your VCS:
Choose the VCS of your catalog, either **Gitlab** or **Github**. In this example, we will use Gitlab.

## Set up Repository Credentials:

1. Enter the repository URL in the **Repository URL** field.

2. If the repository is private, enter the `Username` in the  **Username** field and the `Token` in the **Password** field.

3. Click on **Connect** check the connection of the repository.

There are two ways of adding a new catalog item to ENBUILD:
- [Adding a New Catalog Item Manually](#adding-a-new-catalog-item-manually)
- [Adding a New Catalog Item using exported JSON](#adding-a-new-catalog-item-using-exported-json)

# Adding a New Catalog Item Manually.
Login to ENBUILD as admin user.
Navigate to **Manage Catalogs** tab and click on **Create Catalog**. You will be presented with the following form:

<picture><img src="/images/how-to-guides/createCatalog.png" alt="Add New Catalog"></img></picture>

## Choose your VCS:
Choose the VCS of your catalog, either **Gitlab** or **Github**. In this example, we will use Gitlab.

## Set up Catalog Details:
1. Enter a **Name** for your catalog. This name is displayed to the user while they are browsing the ENBUILD.
2. Choose a **Type** for your catalog items. The type defines the type of the catalog. 
3. Choose the template  **Repository** for your catalog items. This is the template repository that you have created in previous step.
4. If the template repository is private, Check the `IsPrivate` button and provide the Readonly Access Token in the **Token** field.
5. Provide the **Project ID** of the template repository. 
6. Enter a **Readme File Path** this is the path of the README file that you want to disaply when user clicks on your catalog item.
7. Enter a **Values Folder Path** this is the directory of the  Values files for the components.
8. Choose the **Branch** for your catalog template to be used. (only required in Gitlab)
9. Enter an  **Image Path** of the image to be displayed for a catalog item in ENBUILD UI.
10. Select **Multi Select**. This feature allows users to select multiple components from same catagary in a catalog. If enabled the user can install/enable multiple components at once. Otherwise a single component is installed/enabled at a time for a component catagory.
11. Click on **Save And Continue** to proceed to add [Components for the Catalogs](#set-up-component-details)

## Set up Component Details

<picture><img src="/images/how-to-guides/createComponent.png" alt="Add New Component"></img></picture>

1. Enter a **Name** for your component.
2. Select a **Tool Type** for your component item. This is used to group components in a single group.
3. Provide component **Repository** for your component items. Ommit if the component and catalog are in same repository.
4. Provide the **Project ID** of the template repository. Ommit if the component and catalog are in same repository.
5. Choose the Branch as **Ref** for your component template to be used. (only required in Gitlab) Ommit if the component and catalog are in same repository.
6. Enter a **Variable File Path** this is the input variable file which ENBUILD will display for the components , and user will use to change the beheviour of the deployment.
7. Enter an  **Image Path** of the image to be displayed for a component item in ENBUILD UI.
8. Select **Mandatory**. If the Catalog is created enabling the **Multi Select** . Enabling this **Mandatory**  will set This component to be enabled by default in ENBUILD UI.
9. Click on **Save And Continue** to proceed to add [Infrastrcture for the Catalogs](#set-up-infrastrcture-details-for-the-catalogs)

## Set up Infrastrcture Details for the catalogs
Here you can select all the infrastructure providers that the catalog needs to be installed.
For now we support [AWS](https://docs.enbuild.io/en/latest/aws/) ,[Azure](https://docs.enbuild.io/en/latest/azure/)  [Oracle](https://docs.enbuild.io/en/latest/aws/) and  [KUBERNETES ](https://docs.enbuild.io/en/latest/aws/)

You can select 1 or multiple providers for your catalog, but you must have at least one provider selected for the catalog to be valid. 
And your catalog template should have provision to support all the selected providers. 

<picture><img src="/images/how-to-guides/catalogInfraSetup.png" alt="Setup Infrastrcture Details"></img></picture>

## Adding the Permission for the catalogs to ENBUILD ROLES

Once the Catalog is created , you need to add permissions to the catalog for a particuler role that can have access to this catalog.

And you can provide the permissions to each catalog as per your requirement for each catalog and roles.

<picture><img src="/images/how-to-guides/catalogPermissions.png" alt="Catalog Permissions"></img></picture>

# Adding a New Catalog Item using exported JSON

The other way of creating the Catalog Item is providing the raw catalog json file that you have exported from ENBUILD or uploading the json file containing the catalog information.

1. Login to ENBUILD as admin user.
2. Navigate to **Manage Catalogs** tab and click on **Create Catalog**. and then 
3. **Choose the VCS**  of your catalog, either **Gitlab** or **Github**. In this example, we will use Gitlab. 
4. Click on **Upload Json** button to to proceed to import the catalog using JSON 
<picture><img src="/images/how-to-guides/catalogImport.png" alt="Import Catalog"></img></picture>
5. Here you can provide the raw json of any catalog exported from ENBUILD UI. OR upload a file from your local machine having the catalog json.