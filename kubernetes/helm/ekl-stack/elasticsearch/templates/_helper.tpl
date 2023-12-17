{{- define "elasticsearch-labels" }}
  k8s-app: elasticsearch-logging
  kubernetes.io/cluster-service: "true"
  addonmanager.kubernetes.io/mode: Reconcile
{{- end }}

{{- define "elasticsearch-labels-selectors" }}
  k8s-app: elasticsearch-logging
{{- end }}
