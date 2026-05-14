{{/*
Expand the name of the chart.
*/}}
{{- define "enbuild-helm-chart.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "enbuild-helm-chart.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "enbuild-helm-chart.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "enbuild-helm-chart.labels" }}
{{- $chart := include "enbuild-helm-chart.chart" . }}
{{- $selector := include "enbuild-helm-chart.selectorLabels" . }}
helm.sh/chart: {{ $chart }}
{{ $selector | nindent 0 }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "enbuild-helm-chart.selectorLabels" }}
{{- $name := include "enbuild-helm-chart.name" . }}
app.kubernetes.io/name: {{ $name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "enbuild-helm-chart.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "enbuild-helm-chart.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{- /*
  Build a Docker config.json containing one OR more registry auth entries,
  ready for inclusion in a kubernetes.io/dockerconfigjson Secret.

  Always includes the primary registry (.Values.global.image.registry +
  .registry_credentials). Optionally adds a second entry for the GitLab
  Container Registry if .Values.global.gitlabRegistryCredentials.username
  is non-empty.

  Why two registries: vendor13-ib (and similar prototype-iteration setups)
  occasionally need to pull container images from the GitLab staging
  registry (registry.gitlab.com/enbuild-staging/vivsoft-platform-ui/*)
  ahead of the Iron Bank rebuild cycle. Merging both registries' auth
  into ONE Secret means every existing Deployment that references
  `{{ .Release.Name }}-image-pull-secret` automatically gains access to
  both — no per-service plumbing change required.

  Both blocks read placeholder values from values.yaml (e.g.
  REGISTRY1_USER_NAME); operators supply real values at install time via
  --set or a separate untracked secrets values file. Real credentials are
  never committed.
*/}}
{{- define "imagePullSecret" }}
{{- $primary := dict "registry" .Values.global.image.registry "creds" .Values.global.image.registry_credentials }}
{{- $secondary := .Values.global.gitlabRegistryCredentials | default dict }}
{{- $auths := dict }}
{{- $_ := set $auths $primary.registry (dict "username" $primary.creds.username "password" $primary.creds.password "auth" (printf "%s:%s" $primary.creds.username $primary.creds.password | b64enc)) }}
{{- if and $secondary.username $secondary.password }}
{{- $gitlabHost := default "registry.gitlab.com" $secondary.registry }}
{{- $_ := set $auths $gitlabHost (dict "username" $secondary.username "password" $secondary.password "auth" (printf "%s:%s" $secondary.username $secondary.password | b64enc)) }}
{{- end }}
{{- dict "auths" $auths | toJson | b64enc }}
{{- end }}
