resource "aws_eks_node_group" "worker_nede_group" {
  cluster_name    = aws_eks_cluster.kubernetes_cluster.name
  node_group_name = var.node_group_name
  node_role_arn   = aws_iam_role.worker_node_role.arn
  subnet_ids      = var.vpc_subnets
  capacity_type = var.capacity_type
  disk_size = var.disk_size
  instance_types = var.instance_types


  scaling_config {
    desired_size = 2
    max_size     = 3
    min_size     = 2
  }

  update_config {
    max_unavailable = 1
  }

  tags = {
    env = var.enviroment
  }

  depends_on = [
    aws_iam_role_policy_attachment.workerNode_AmazonEKSWorkerNodePolicy,
    aws_iam_role_policy_attachment.workerNode_AmazonEKS_CNI_Policy,
    aws_iam_role_policy_attachment.workerNode_AmazonEC2ContainerRegistryReadOnly,
  ]
}