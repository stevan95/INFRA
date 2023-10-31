data "aws_iam_policy_document" "management_eks_role_document" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["eks.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "management_eks_role" {
  name               = "management_eks_role"
  assume_role_policy = data.aws_iam_policy_document.management_eks_role_document.json
}

resource "aws_iam_role_policy_attachment" "management_eks_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.management_eks_role.name
}

#Worker Nodes IAM role

resource "aws_iam_role" "worker_node_role" {
  name = "worker_node_role"

  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy_attachment" "workerNode_AmazonEKSWorkerNodePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.worker_node_group_iam_role.name
}

resource "aws_iam_role_policy_attachment" "workerNode_AmazonEKS_CNI_Policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.worker_node_group_iam_role.name
}

resource "aws_iam_role_policy_attachment" "workerNode_AmazonEC2ContainerRegistryReadOnly" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.worker_node_group_iam_role.name
}
