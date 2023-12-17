{{- define "fluentd-labels" }}
  k8s-app: fluentd-es
  kubernetes.io/cluster-service: "true"
  addonmanager.kubernetes.io/mode: Reconcile
{{- end }}

{{- define "fluentd-selectors" }}
  k8s-app: fluentd-es
  version: v1.16-1
{{- end }}
