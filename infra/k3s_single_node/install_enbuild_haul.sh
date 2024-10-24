#!/usr/bin/env bash

set -eu -o pipefail

K3S_VERSION="v1.30.2+k3s2"
HELM_VERSION="v3.16.1"
ENBUILD_HELM_CHART_VERSION=0.0.20
ENBUILD_HAULER_URL="https://enbuild-haul.s3.us-east-1.amazonaws.com/enbuild-${ENBUILD_HELM_CHART_VERSION}.tar.zst"

# Update package list and install jq
apt-get update
apt-get install -y jq

# Download and install Helm
curl https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz -o /tmp/helm-${HELM_VERSION}-linux-amd64.tar.gz
tar -xvzf /tmp/helm-${HELM_VERSION}-linux-amd64.tar.gz
cp -u linux-amd64/helm /usr/local/bin/

# Download and install K3s
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=$K3S_VERSION INSTALL_K3S_EXEC="server --disable=traefik --write-kubeconfig-mode=666" sh -
# curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=$K3S_VERSION INSTALL_K3S_EXEC="server --disable=traefik --disable=local-storage" sh -
sleep 5
systemctl enable --now k3s
sleep 5

# Set up kubectl alias and kubeconfig
echo "alias k='/usr/local/bin/kubectl'" >>/root/.bashrc
mkdir -p /home/ubuntu/.kube/ /root/.kube/
/usr/local/bin/kubectl config view --raw > /home/ubuntu/.kube/config
/usr/local/bin/kubectl config view --raw > /root/.kube/config
sleep 5

# Install Hauler
curl -sfL https://get.hauler.dev | bash

# Download and load the Enbuild Hauler package
curl -O ${ENBUILD_HAULER_URL} 
hauler store load enbuild-${ENBUILD_HELM_CHART_VERSION}.tar.zst
hauler store serve registry &
sleep 15

# get the private IP of the node
PRIVATE_IP=$(hostname -I | awk '{print $1}')

# Pull and install the Enbuild Helm chart
helm pull --plain-http oci://${PRIVATE_IP}:5000/hauler/enbuild --version ${ENBUILD_HELM_CHART_VERSION}

echo """global:
  image:
    registry: ${PRIVATE_IP}:5000
    pullPolicy: Always
rabbitmq:
  image:
    registry: ${PRIVATE_IP}:5000
    repository: bitnami/rabbitmq
    tag: 3.11.13-debian-11-r0
""" > quick_install_hauler.yaml

echo """
mirrors:
  "${PRIVATE_IP}:5000":
    endpoint:
      - "http://${PRIVATE_IP}:5000"
  "registry.gitlab.com":
    endpoint:
      - "http://${PRIVATE_IP}:5000"

configs:
  "${PRIVATE_IP}:5000":
    tls:
      insecure_skip_verify: true

"""  > /etc/rancher/k3s/registries.yaml

systemctl restart k3s

helm upgrade --install --namespace enbuild enbuild --plain-http oci://${PRIVATE_IP}:5000/hauler/enbuild --version ${ENBUILD_HELM_CHART_VERSION}  -f quick_install_hauler.yaml --create-namespace

