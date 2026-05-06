{{/*
Common labels for step-ca subchart resources.
*/}}
{{- define "step-ca.labels" -}}
app.kubernetes.io/name: step-ca
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: pki
app.kubernetes.io/part-of: enbuild
{{- if .Chart }}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" }}
{{- end }}
{{- end -}}

{{- define "step-ca.selectorLabels" -}}
app.kubernetes.io/name: step-ca
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}
