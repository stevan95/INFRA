output "endpoint" {
  value = aws_eks_cluster.kubernetes_cluster.endpoint
}

output "kubeconfig-certificate-authority-data" {
  value = aws_eks_cluster.kubernetes_cluster.certificate_authority[0].data
}

output "oidc-issuer" {
  value = aws_eks_cluster.kubernetes_cluster.identity[0].oidc[0].issuer
}
