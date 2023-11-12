variable "provider_arn" {
  type = string
  description = "OIDC IAM provider ARN"
}

variable "endpoint" {
  type = string
  description = "Cluster endpint string"
}

variable "cluster_id" {
  type = string
  description = "Cluster id"
}

variable "kubeconfig_cert" {
  type = string
  description = "Kubernetes cluster config"
}
