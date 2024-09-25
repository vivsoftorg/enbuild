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



delete_k3d_cluster() {
  clusterName=$1
  clusterExist=$(k3d cluster list -o json | jq '.[] | select(.name=="'"$clusterName"'") | .name')
  if [ -n "$clusterExist" ]
  then
    k3d cluster delete "$clusterName" || true
  fi
  docker network rm "k3d-${clusterName}" > /dev/null 2>&1 || true
  k3d registry delete qovery-registry.lan > /dev/null 2>&1 || true
}

teardown_network() {
  if [ "$(uname -s)" = 'Darwin' ]; then
    # MacOs
    set -e
    sudo ifconfig lo0 -alias 172.42.0.3/32 up > /dev/null 2>&1 || true
  elif grep -qi microsoft /proc/version; then
    # Wsl
    set -x
    echo "Removing network config, please provide your password to run the sudo command"
    sudo ip addr del 172.42.0.3/32 dev lo || true
    ${POWERSHELL_CMD} -Command "Start-Process powershell -Verb RunAs -ArgumentList \"netsh interface ipv4 delete address name='Loopback Pseudo-Interface 1' address=172.42.0.3\""
  fi
  set +x
}

# shellcheck disable=SC2046
# shellcheck disable=SC2086
cd "$(dirname $(realpath $0))"
echo "Removing $CLUSTER_NAME kube cluster"
delete_k3d_cluster "$CLUSTER_NAME"

echo 'Removing network config'
teardown_network

echo "Local demo cluster $CLUSTER_NAME deleted."