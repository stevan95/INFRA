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

resource "aws_eks_identity_provider_config" "identity_provider_config" {
  cluster_name = aws_eks_cluster.kubernetes_cluster.name

  oidc {
    //client_id                     = "[ ${join(", ", [for s in aws_iam_openid_connect_provider.iam_openid_connect.client_id_list : format("%q", s)])} ]"
    client_id = "sts.amazonaws.com"
    identity_provider_config_name = var.identity_provider_config_name
    issuer_url                    = "https://${aws_iam_openid_connect_provider.iam_openid_connect.url}"
  }

  tags = {
    env = var.enviroment
  }

  depends_on = [ aws_eks_cluster.kubernetes_cluster, aws_iam_openid_connect_provider.iam_openid_connect ]
}
