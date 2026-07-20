#!/usr/bin/env bash
set -euo pipefail

chart_dir="${1:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"
secret_render="$(mktemp)"
backend_render="$(mktemp)"
operator_render="$(mktemp)"
trap 'rm -f "$secret_render" "$backend_render" "$operator_render"' EXIT

common_args=(
  --set global.evalMode=true
  --set global.istio.enabled=false
  --set global.image.registry_credentials.username=fake
  --set global.image.registry_credentials.password=fake
  --set enbuildBk.encryption_key=fake-encryption-key-for-render-only
  --set mongodb.enabled=true
  --set keycloak.enabled=true
  --set keycloak.hostname=http://localhost:8088
  --set enbuildBk.authMechanism=keycloak
  --set enbuildBk.keycloak.url=http://localhost:8088
  --set enbuildBk.keycloak.backendUrl=http://writer-test-keycloak:8080
  --set enbuildBk.grpcVirtualService.enabled=false
)

helm template writer-test "$chart_dir" "${common_args[@]}" \
  --set keycloak.demoRealm.enabled=true \
  --show-only templates/keycloak-demo-realm-secret.yaml > "$secret_render"

grep -q '^  KEYCLOAK_CLIENT_SECRET:' "$secret_render"
grep -q '^  KEYCLOAK_ADMIN_CLIENT_SECRET:' "$secret_render"

helm template writer-test "$chart_dir" "${common_args[@]}" \
  --set keycloak.demoRealm.enabled=true \
  --show-only templates/enbuild-bk.yaml > "$backend_render"

grep -q 'name: writer-test-keycloak-secrets' "$backend_render"

helm template writer-test "$chart_dir" "${common_args[@]}" \
  --set keycloak.demoRealm.enabled=false \
  --set keycloak.existingSecret=operator-realm \
  --set enbuildBk.keycloak.existingSecret=operator-writer \
  --show-only templates/enbuild-bk.yaml > "$operator_render"

grep -q 'name: operator-writer' "$operator_render"
if grep -q 'name: writer-test-keycloak-secrets' "$operator_render"; then
  echo 'demo Keycloak Secret must not override an explicit operator writer Secret' >&2
  exit 1
fi

echo 'Keycloak writer Secret render contract: PASS'
