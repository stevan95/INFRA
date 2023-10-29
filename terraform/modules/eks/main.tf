resource "aws_eks_cluster" "kubernetes_cluster" {
  
  name     = var.kubernetes_cluster_name
  role_arn = aws_iam_role.management_eks_role.arn
  version  = var.kubernetescluster_version

  vpc_config {
    subnet_ids = toset(var.vpc_subnets)
  }

  depends_on = [
    aws_iam_role_policy_attachment.management_eks_policy,
  ]

  tags = {
    env = var.enviroment
  }
}
