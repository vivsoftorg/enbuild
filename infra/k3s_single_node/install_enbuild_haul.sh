#!/usr/bin/env bash

set -eu -o pipefail

ENBUILD_HELM_CHART_VERSION=0.0.20
ENBUILD_HAULER_URL="https://enbuild-haul.s3.us-east-1.amazonaws.com/enbuild-${ENBUILD_HELM_CHART_VERSION}.tar.zst"

echo "Starting script execution."

# Wait till kubectl get nodes returns the node
echo "Waiting for the node to be ready."
while ! kubectl get nodes; do
  sleep 1
done

# Download and load the Enbuild Hauler package
echo "Downloading Enbuild Hauler package from ${ENBUILD_HAULER_URL}."
curl -O ${ENBUILD_HAULER_URL}
echo "Loading Enbuild Hauler package."
hauler store load enbuild-${ENBUILD_HELM_CHART_VERSION}.tar.zst
hauler store serve registry &
sleep 5
echo "Enbuild Hauler package loaded and registry serving started."

# Get the private IP of the node
echo "Retrieving the private IP of the node."
PRIVATE_IP=$(hostname -I | awk '{print $1}')
echo "Private IP is ${PRIVATE_IP}."

# Wait until the registry is reachable
echo "Waiting for the registry at ${PRIVATE_IP}:5000 to be reachable."
while ! nc -z $PRIVATE_IP 5000; do
  sleep 1
done
echo "Registry is now reachable."

# Pull and install the Enbuild Helm chart
echo "Pulling Enbuild Helm chart."
helm pull --plain-http oci://${PRIVATE_IP}:5000/hauler/enbuild --version ${ENBUILD_HELM_CHART_VERSION}

# Creating configuration files needed for installation
echo "Creating quick_install_hauler.yaml configuration file."
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

echo "Creating /etc/rancher/k3s/registries.yaml configuration file."
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

echo "Restarting k3s service."
systemctl restart k3s
echo "k3s service restarted."

# Upgrade and install the Helm release
echo "Upgrading and installing Helm release enbuild."
helm upgrade --install --namespace enbuild enbuild --plain-http oci://${PRIVATE_IP}:5000/hauler/enbuild --version ${ENBUILD_HELM_CHART_VERSION} -f quick_install_hauler.yaml --create-namespace
echo "Helm release enbuild has been upgraded and installed."

echo "Script execution finished successfully."
