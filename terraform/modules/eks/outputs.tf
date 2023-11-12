output "endpoint" {
  value = aws_eks_cluster.kubernetes_cluster.endpoint
}

output "cluster_id" {
  value = aws_eks_cluster.kubernetes_cluster.id
}

output "kubeconfig-certificate-authority-data" {
  value = aws_eks_cluster.kubernetes_cluster.certificate_authority[0].data
}

output "oidc-issuer" {
  value = aws_eks_cluster.kubernetes_cluster.identity[0].oidc[0].issuer
}

output "provider_arn" {
  value = aws_iam_openid_connect_provider.iam_openid_connect.arn
}
