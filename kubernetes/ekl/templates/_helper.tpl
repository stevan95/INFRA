{{- define "fluentd-labels" }}
  k8s-app: fluentd-es
  kubernetes.io/cluster-service: "true"
  addonmanager.kubernetes.io/mode: Reconcile
{{- end }}
