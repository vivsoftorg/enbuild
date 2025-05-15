# ENBUILD CLI

ENBUILD CLI is a command-line tool designed to interact with the ENBUILD platform and help generate catalog templates, haul manifests, and manage resources.

## Prerequisites

Make sure you install the [yq cli](https://mikefarah.gitbook.io/yq) as the ENBUILD CLI uses it internally for creating catalog template values files.

## Installation

*Installation instructions to be added*

## Usage

```
enbuild [command] [subcommand] [flags]
```

## Available Commands

```
enbuild is a CLI to work with ENBUILD

Usage:
  enbuild [flags]
  enbuild [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Create resources for Enbuild
  demo        Try Enbuild on your local machine
  get         Get resources from Enbuild
  help        Help about any command
  version     Print installed version of the Enbuild CLI

Flags:
      --base-url string   API base URL for ENBUILD (or set env variable ENBUILD_BASE_URL) (default "https://enbuild-dev.vivplatform.io/enbuild-bk/")
      --debug             Enable debug output
  -h, --help              help for enbuild
      --token string      API token (or set env variable ENBUILD_API_TOKEN)
  -v, --version           Print the version and exit
```

Note: The `--base-url` and `--token` flags can be used with any command or subcommand.

## Command Details

### Create Resources

Create various resources for Enbuild such as templates and haul manifests.

```
Usage:
  enbuild create [command]

Available Commands:
  bigbang-template Create a BigBang ENBUILD Catalog template for given version
  haul             Create haul manifests
```

#### Create BigBang Template

Create a BigBang ENBUILD Catalog template for a specific version.

```
Usage:
  enbuild create bigbang-template [flags]

Flags:
  -v, --bb-version string   Specify the BigBang version (required)
```

#### Create Haul Manifests

Create haul manifests for various components.

```
Usage:
  enbuild create haul [command]

Available Commands:
  bigbang     Create a haul manifest file for BigBang
  enbuild     Create a haul manifest file for the ENBUILD Helm Chart
```

##### Create Haul for BigBang

Create a haul manifest file for BigBang.

```
Usage:
  enbuild create haul bigbang [flags]

Flags:
  -v, --bb-version string   Specify the BigBang version (required)
```

##### Create Haul for ENBUILD

Create a haul manifest file for the ENBUILD Helm Chart.

```
Usage:
  enbuild create haul enbuild [flags]

Flags:
  -v, --helm-chart-version string   Specify the ENBUILD Helm Chart version
```

### Get Resources

Get various resources from Enbuild such as catalogs, manifests, etc.

```
Usage:
  enbuild get [command]

Available Commands:
  catalogs    Get catalogs from Enbuild
```

#### Get Catalogs

Retrieve and list catalogs from the Enbuild platform.

```
Usage:
  enbuild get catalogs [flags]

Flags:
      --id string     Get catalog by ID
      --name string   Search catalogs by name
      --type string   Filter catalogs by type (e.g., terraform)
      --vcs string    Filter catalogs by VCS (e.g., github)
```

### Demo

Try Enbuild on your local machine.

```
Usage:
  enbuild demo [command]

Available Commands:
  destroy     Remove k3d cluster with ENBUILD installed on your local machine
  down        Uninstall ENBUILD on local k3d cluster and stop the k3d cluster your local machine
  up          Create a k3d kubernetes cluster with ENBUILD installed on your local machine
```

### Version

Print installed version of the Enbuild CLI.

```
Usage:
  enbuild version
```

## Examples

### Creating a BigBang Template

```bash
# Using flags
enbuild create bigbang-template --bb-version 2.5.0 --token your-api-token --base-url https://custom-enbuild.example.com/

# Using environment variables
export ENBUILD_API_TOKEN=your-api-token
export ENBUILD_BASE_URL=https://custom-enbuild.example.com/
enbuild create bigbang-template --bb-version 2.5.0
```

### Creating a Haul Manifest for BigBang

```bash
enbuild create haul bigbang --bb-version 2.5.0 --token your-api-token
```

### Creating a Haul Manifest for ENBUILD

```bash
enbuild create haul enbuild --helm-chart-version 0.1.0 --base-url https://custom-enbuild.example.com/
```

### Getting Catalogs

```bash
# List all catalogs
enbuild get catalogs --token your-api-token

# Get a specific catalog by ID
enbuild get catalogs --id 6638a128d6852d0012a27491 --token your-api-token

# Filter catalogs by type
enbuild get catalogs --type terraform --base-url https://custom-enbuild.example.com/
```

## Environment Variables

- `ENBUILD_API_TOKEN`: API token for authentication with ENBUILD
- `ENBUILD_BASE_URL`: Base URL for the ENBUILD API

## Version

```bash
enbuild version
```

Current version: v0.0.11
