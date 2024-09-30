#!/bin/sh

set -eu

CLUSTER_NAME=${1:-enbuild}
DEBUG=${2:-false}
VALUES_FILE=${3:-"/tmp/enbuild/values.yaml"}

HELM_DEBUG=""

if [ "$DEBUG" == "true" ]; then
  set -x
  HELM_DEBUG="--debug"
fi

POWERSHELL_CMD='powershell.exe'

get_or_create_cluster() {
  clusterName=$1
  if ! k3d cluster list -o json | jq -e '.[] | select(.name=="'"$clusterName"'")' >/dev/null; then
    k3d cluster create "$clusterName" \
      --image 'docker.io/rancher/k3s:v1.28.9-k3s1' \
      --subnet '172.42.0.0/16' \
      --k3s-arg "--node-ip=172.42.0.3@server:0" \
      --k3s-arg "--disable=traefik@server:*" \
      --registry-create enbuild-registry.lan \
      --port "80:80@loadbalancer" --port "443:443@loadbalancer"
  else
    k3d cluster start "$clusterName"
    k3d kubeconfig merge "$clusterName" -d
  fi
}

install_or_upgrade_helm_charts() {
  if ! helm list -n enbuild -o json | jq -e '.[] | select(.name=="enbuild")' >/dev/null; then
    set -x
    helm upgrade --install --create-namespace ${HELM_DEBUG} --timeout=15m -n enbuild -f $VALUES_FILE --atomic \
      --set global.image.cert-pullPolicy=IfNotPresent \
      enbuild vivsoft/enbuild
  fi

  for i in $(seq 1 3); do
    helm upgrade --install --create-namespace ${HELM_DEBUG} --timeout=15m -n enbuild -f $VALUES_FILE --wait --atomic enbuild vivsoft/enbuild && break
    set +x
    echo "Install failed. Retrying in 10 seconds."
    0
  done

  set +x
}

setup_network() {
  case "$(uname -s)" in
  Darwin)
    echo "Setting up network. Please provide your password to run the sudo command"
    sudo ifconfig lo0 alias 172.42.0.3/32 up || true
    ;;
  Linux)
    if grep -qi microsoft /proc/version; then
      echo "Setting up network. Please provide your password to run the sudo command"
      sudo ip addr add 172.42.0.3/32 dev lo || true
      ${POWERSHELL_CMD} -Command "Start-Process powershell -Verb RunAs -ArgumentList \"netsh interface ipv4 add address name='Loopback Pseudo-Interface 1' address=172.42.0.3 mask=255.255.255.255 skipassource=true\""
    fi
    ;;
  esac
}

try_install_missing_deps() {
  if command -v sudo >/dev/null; then
    SUDO="sudo"
  else
    SUDO=""
  fi

  if command -v apt-get >/dev/null; then
    echo "Installing dependencies with apt"
    ${SUDO} apt-get update && ${SUDO} apt-get install -y jq grep sed curl iproute2
  elif command -v yum >/dev/null; then
    echo "Installing dependencies with yum"
    ${SUDO} yum update -y && ${SUDO} yum install -y jq grep sed curl iproute
  elif command -v pacman >/dev/null; then
    echo "Installing dependencies with pacman"
    ${SUDO} pacman -Sy && ${SUDO} pacman --noconfirm -S jq grep curl sed iproute
  elif command -v brew >/dev/null; then
    echo "Installing dependencies with brew"
    brew update && brew install jq grep curl helm k3d
  else
    echo "Cannot detect your package manager. Please install the following commands: jq grep curl sed iproute2 helm k3d"
    exit 1
  fi
}

install_deps() {
    if which docker >/dev/null; then
     echo "docker already installed"
  else
    echo "docker command is missing. Please use your package manager to install it"
    echo "https://docs.docker.com/engine/install/"
    exit 1
  fi

  docker_running=$( (docker ps -q >/dev/null && echo true ) || echo false )
  if "$docker_running" == "true"; then
     echo "docker is running"
  else
    echo "Docker is not running. Please start Docker before running this command"
    exit 1
  fi

  if which jq >/dev/null; then
     echo "jq already installed"
  else
    try_install_missing_deps
  fi

  if which grep >/dev/null; then
     echo "grep already installed"
  else
    try_install_missing_deps
  fi

  if which sed >/dev/null; then
     echo "sed already installed"
  else
    try_install_missing_deps
  fi

  if test -f /proc/version && grep -qi microsoft /proc/version; then
    if which ip >/dev/null; then
       echo "iproute already installed"
    else
      try_install_missing_deps
    fi
  fi

  if which curl >/dev/null; then
     echo "curl already installed"
  else
    try_install_missing_deps
  fi

  if which k3d >/dev/null; then
    echo "k3d already installed"
  else
    echo "Installing k3d"
    curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | TAG=v5.6.3 bash
  fi

  if which helm >/dev/null; then
     echo "helm already installed"
  else
    echo "Installing HELM"
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
  fi

  # Wsl powershell
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

  echo "All dependencies are installed"
}

cd "$(dirname "$(realpath "$0")")"

echo 'Checking and installing dependencies'

install_deps

echo 'Fetching ENBUILD values to setup your cluster'


mkdir -p "$(dirname "$VALUES_FILE")"
curl -s -L https://raw.githubusercontent.com/vivsoftorg/enbuild/refs/heads/main/examples/enbuild/quick_install.yaml > $VALUES_FILE
echo 'Helm values written into $VALUES_FILE'

echo 'Installing ENBUILD helm repositories'

helm repo add vivsoft https://vivsoftorg.github.io/enbuild
helm repo update vivsoft

echo "Creating $CLUSTER_NAME kube cluster"

get_or_create_cluster "$CLUSTER_NAME"

echo 'Installing ENBUILD helm charts'

install_or_upgrade_helm_charts

echo 'Configuring network'

setup_network
echo "---------------------------------------------------------"
echo "ENBUILD demo cluster is now installed !!!!"
echo "The kubeconfig is correctly set, so you can connect to it directly with kubectl or k9s from your local machine"
echo "To delete/stop/start your cluster, use k3d cluster $CLUSTER_NAME"
echo "To access the ENBUILD dashboard, use port-forward with below command"
echo "kubectl --namespace enbuild port-forward svc/enbuild-enbuild-ui 3000:80"
echo "---------------------------------------------------------"
