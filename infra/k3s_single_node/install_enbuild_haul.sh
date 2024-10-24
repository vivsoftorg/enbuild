#!/usr/bin/env bash

set -eu -o pipefail

K3S_VERSION="v1.30.2+k3s2"
HELM_VERSION="v3.16.1"
ENBUILD_HELM_CHART_VERSION=0.0.20
ENBUILD_HAULER_URL="https://enbuild-haul.s3.us-east-1.amazonaws.com/enbuild-${ENBUILD_HELM_CHART_VERSION}.tar.zst"

echo "Starting script execution..."

# # Update package list and install jq
# echo "Updating package list and installing jq..."
# apt-get update && apt-get install -y jq
# echo "Package list updated and jq installed."

# Download and install Helm
echo "Downloading Helm version ${HELM_VERSION}..."
curl https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz -o /tmp/helm-${HELM_VERSION}-linux-amd64.tar.gz
echo "Extracting Helm archive..."
tar -xvzf /tmp/helm-${HELM_VERSION}-linux-amd64.tar.gz
echo "Installing Helm..."
cp -u linux-amd64/helm /usr/local/bin/
echo "Helm installed successfully."

# Download and install K3s
echo "Downloading and installing K3s version ${K3S_VERSION}..."
curl -sfL https://get.k3s.io | INSTALL_K3S_VERSION=$K3S_VERSION INSTALL_K3S_EXEC="server --disable=traefik --write-kubeconfig-mode=666" sh -
sleep 5
systemctl enable --now k3s
echo "K3s installed and service started."

# Set up kubectl alias and kubeconfig
echo "Setting up kubectl alias and kubeconfig files..."
echo "alias k='/usr/local/bin/kubectl'" >>/root/.bashrc
mkdir -p /home/ubuntu/.kube/ /root/.kube/
/usr/local/bin/kubectl config view --raw > /home/ubuntu/.kube/config
/usr/local/bin/kubectl config view --raw > /root/.kube/config
echo "kubectl setup completed."

# Install Hauler
echo "Installing Hauler..."
curl -sfL https://get.hauler.dev | bash
echo "Hauler installed successfully."

# Download and load the Enbuild Hauler package
echo "Downloading Enbuild Hauler package from ${ENBUILD_HAULER_URL}..."
curl -O ${ENBUILD_HAULER_URL}
echo "Loading Enbuild Hauler package..."
hauler store load enbuild-${ENBUILD_HELM_CHART_VERSION}.tar.zst
hauler store serve registry &
sleep 5
echo "Enbuild Hauler package loaded and registry serving started."

# Get the private IP of the node
echo "Retrieving the private IP of the node..."
PRIVATE_IP=$(hostname -I | awk '{print $1}')
echo "Private IP is ${PRIVATE_IP}."

# Wait until the registry is reachable
echo "Waiting for the registry at ${PRIVATE_IP}:5000 to be reachable..."
while ! nc -z $PRIVATE_IP 5000; do
  sleep 1
done
echo "Registry is now reachable."

# Pull and install the Enbuild Helm chart
echo "Pulling Enbuild Helm chart..."
helm pull --plain-http oci://${PRIVATE_IP}:5000/hauler/enbuild --version ${ENBUILD_HELM_CHART_VERSION}

# Creating configuration files needed for installation
echo "Creating quick_install_hauler.yaml configuration file..."
cat <<EOF > quick_install_hauler.yaml
global:
  image:
    registry: ${PRIVATE_IP}:5000
    pullPolicy: Always
rabbitmq:
  image:
    registry: ${PRIVATE_IP}:5000
    repository: bitnami/rabbitmq
    tag: 3.11.13-debian-11-r0
EOF

echo "Creating /etc/rancher/k3s/registries.yaml configuration file..."
cat <<EOF > /etc/rancher/k3s/registries.yaml
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
EOF

echo "Restarting k3s service..."
systemctl restart k3s
echo "k3s service restarted."

# Upgrade and install the Helm release
echo "Upgrading and installing Helm release enbuild..."
helm upgrade --install --namespace enbuild enbuild --plain-http oci://${PRIVATE_IP}:5000/hauler/enbuild --version ${ENBUILD_HELM_CHART_VERSION} -f quick_install_hauler.yaml --create-namespace
echo "Helm release enbuild has been upgraded and installed."

echo "Script execution finished successfully."
