#!/usr/bin/env bash

set -eu -o pipefail

K3S_VERSION="v1.30.2+k3s2"
HELM_VERSION="v3.16.1"

echo "Starting script execution..."

echo "Making sure the kubectl is working properly..."
kubectl get nodes

# Update package list and install jq
echo "Updating package list and installing jq..."
apt-get update && apt-get install -y jq
echo "Package list updated and jq installed."

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
echo "Script execution completed."
 
