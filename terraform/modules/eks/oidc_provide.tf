data "tls_certificate" "cluster_certificates" {
  url = aws_eks_cluster.kubernetes_cluster.identity[0].oidc[0].issuer
}

resource "aws_iam_openid_connect_provider" "iam_openid_connect" {
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = [data.tls_certificate.cluster_certificates.certificates[0].sha1_fingerprint]
  url             = data.tls_certificate.cluster_certificates.url
}

data "aws_iam_policy_document" "assume_role_policy" {
  statement {
    actions = ["sts:AssumeRoleWithWebIdentity"]
    effect  = "Allow"

    condition {
      test     = "StringEquals"
      variable = "${replace(aws_iam_openid_connect_provider.iam_openid_connect.url, "https://", "")}:sub"
      values   = ["system:serviceaccount:kube-system:aws-node"]
    }

    principals {
      identifiers = [aws_iam_openid_connect_provider.iam_openid_connect.arn]
      type        = "Federated"
    }
  }
}

resource "aws_iam_role" "kubernetes_iam_sa_role" {
  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json
  name               = "kubernetes_iam_sa_role"
}
