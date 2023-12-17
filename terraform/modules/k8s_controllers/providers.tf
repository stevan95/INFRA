provider "helm" {
  kubernetes {
    host                   = var.endpoint
    cluster_ca_certificate = base64decode(var.kubeconfig_cert)
    exec {
      api_version = "client.authentication.k8s.io/v1beta1"
      args        = ["eks", "get-token", "--cluster-name", var.cluster_id]
      command     = "aws"
    }
  }
}

provider "kubernetes" {
  host                   = var.endpoint
  cluster_ca_certificate = base64decode(var.kubeconfig_cert)
  # token                  = data.aws_eks_cluster_auth.default.token

  exec {
    api_version = "client.authentication.k8s.io/v1beta1"
    args        = ["eks", "get-token", "--cluster-name", var.cluster_id]
    command     = "aws"
  }
}