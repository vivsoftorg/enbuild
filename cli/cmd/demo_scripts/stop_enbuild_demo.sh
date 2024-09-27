#!/bin/sh

set -eu

CLUSTER_NAME=${1:-enbuild}
DEBUG=${2:-false}
HELM_DEBUG=""

if [ "$DEBUG" == "true" ]; then
    set -x
    HELM_DEBUG="--debug"
fi

POWERSHELL_CMD='powershell.exe'
if test -f /proc/version && grep -qi microsoft /proc/version; then
  if which 'powershell.exe' >/dev/null; then
    echo "powershell is installed"
    POWERSHELL_CMD='powershell.exe'
  elif which '/mnt/c/Windows/System32/WindowsPowerShell/v1.0/powershell.exe' >/dev/null; then
    echo "powershell is installed"
    POWERSHELL_CMD='/mnt/c/Windows/System32/WindowsPowerShell/v1.0/powershell.exe'
  else
    echo "Cannot find powershell.exe, please be sure it is installed"
    exit 1
  fi
fi



stop_k3d_cluster() {
  clusterName=$1
  clusterExist=$(k3d cluster list -o json | jq '.[] | select(.name=="'"$clusterName"'") | .name')
  if [ -n "$clusterExist" ]
  then
    k3d cluster stop "$clusterName" || true
  fi
}

# shellcheck disable=SC2046
# shellcheck disable=SC2086
cd "$(dirname $(realpath $0))"
echo "Removing $CLUSTER_NAME kube cluster"
stop_k3d_cluster "$CLUSTER_NAME"

echo "Local demo k3d cluster $CLUSTER_NAME stopped."