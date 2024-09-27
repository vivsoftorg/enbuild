# ENBUILD CLI
# ---------------

Make sure you install the [yq cli](https://mikefarah.gitbook.io/yq) as `enbuild` cli is using it internally for creating bigbang catalog template values file.

```
❯ enbuild -h
enbuild is a CLI to help generate the ENBUILD catalog templates

Usage:
  enbuild [command]

Available Commands:
  bigbang     bigbang
  completion  Generate the autocompletion script for the specified shell
  demo        Try Enbuild on your local machine
  help        Help about any command

Flags:
  -h, --help      help for enbuild
  -v, --version   Print the version number of enbuild CLI

Use "enbuild [command] --help" for more information about a command.

 rancher-desktop enbuild_helm_chart/cli cli-up ❯ enbuild -v
v0.0.8
```
