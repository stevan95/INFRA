{{- define "kibana-labels" }}
  k8s-app: kibana-logging
  kubernetes.io/cluster-service: "true"
  addonmanager.kubernetes.io/mode: Reconcile
{{- end }}

{{- define "kibana-labels-selectors" }}
  k8s-app: kibana-logging
{{- end }}
