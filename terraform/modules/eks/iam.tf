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
