#!/bin/bash
set -e

# exit if no arguments provided , we need one argument for helm_chart_version
if [ "$#" -ne 2 ]; then
    echo "Please provide the helm chart version and helm chart location"
    exit 1
fi


HELM_CHART_VERSION=$1
HAULER_FILE="/tmp/enbuild_${HELM_CHART_VERSION}_haul.yaml"
HELM_CHART_LOCATION=$2

# # make sure yq is installed if not install it 
# if ! command -v yq &> /dev/null
# then
#     echo "yq could not be found, installing it"
#     wget https://github.com/mikefarah/yq/releases/download/v4.2.0/yq_linux_amd64 -O /usr/bin/yq && chmod +x /usr/bin/yq
# fi

# make sure hauler is installed if not install it
if ! command -v hauler &> /dev/null
then
    echo "hauler could not be found, installing it"
    curl -sfL https://get.hauler.dev | bash
fi

# make sure the helm chart location is correct
if [ ! -d "$HELM_CHART_LOCATION" ]; then
    echo "Helm chart location is not correct"
    exit 1
fi



# Run the Helm command to get the list of images
# IMAGES=$(helm template $HELM_CHART_LOCATION | yq -N '..|.image? | select(.)' | sort -u)
IMAGES=$(helm template $HELM_CHART_LOCATION | grep -oP 'image: \K\S+' | sort -u)

# error out if no IMAGES is empty
if [ -z "$IMAGES" ]; then
    echo "No images found in the Helm chart"
    exit 1
fi


# Start creating the Hauler file
cat <<EOL > $HAULER_FILE
apiVersion: content.hauler.cattle.io/v1alpha1
kind: Charts
metadata:
  name: enbuild-chart-hauler
spec:
  charts:
    - name: enbuild
      repoURL: https://vivsoftorg.github.io/enbuild
      version: ${HELM_CHART_VERSION}
---
apiVersion: content.hauler.cattle.io/v1alpha1
kind: Images
metadata:
  name: enbuild-images-hauler 
spec:
  images:
EOL

# Append each image to the Hauler file
while IFS= read -r image; do
  echo "    - name: ${image}" >> $HAULER_FILE
  echo "      platform: linux/amd64" >> $HAULER_FILE
done <<< "$IMAGES"

echo "Hauler file generated: $HAULER_FILE"

/usr/local/bin/hauler store sync -f $HAULER_FILE
/usr/local/bin/hauler store save --filename enbuild-${HELM_CHART_VERSION}-haul.tar.zst
