---
title: "Configuring ENBUILD CLI"
description: "Steps to Configure ENBUILD CLI"
summary: "Configuring ENBUILD CLI"
draft: true
menu:
  docs:
    parent: "docs/how-to-guides/deploying-enbuild-for-local-testing/"
    identifier: "configureEnbuildCLI"
weight: 202
toc: true
seo:
  title: "" # custom title (optional)
  description: "" # custom description (recommended)
  canonical: "" # custom canonical URL (optional)
  noindex: false # false (default) or true
---

Follow these step-by-step instructions to configure ENBUILD CLI.

## Prerequisites

Make sure you install the following dependencies.

1. [Docker](https://docs.docker.com/engine/install/)
    - Install docker by following these [steps](https://docs.docker.com/engine/install/).
    - Make sure that docker engine is running before going using the Enbuild CLI.

2. [Golang](https://go.dev/)
    - Install Go programming language by following these [steps](https://go.dev/doc/install).

3. [yq cli](https://mikefarah.gitbook.io/yq)
    - Install yq cli following these [steps](https://github.com/mikefarah/yq/#install).
    - Enbuild cli is using it internally for creating bigbang catalog template values file.


## Configuration

1. Clone the [Enbuild repository](https://github.com/vivsoftorg/enbuild.git)

    ``` bash
    git clone https://github.com/vivsoftorg/enbuild.git
    ```

2. Change your directory to `cli` in the enbuild repository

    ``` bash
    cd <path-to-the-above-cloned-enbuild-repository>/cli
    ```

3. Run the below command to build the `enbuild` cli

    ```bash
    go build
    ```

4. Add `enbuild` command to the PATH environment variable

    ```bash
    export <path-to-the-above-cloned-enbuild-repository>/cli
    ```

5. Verify that `enbuild` cli is ready to use by running these commands.

    ```bash
    enbuild -v
    ```
6. For more information on enbuild cli commands, please run

    ```bash
    enbuild -h
    ```

