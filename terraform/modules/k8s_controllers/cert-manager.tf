data "template_file" "cert_manager_policy" {
  template = "${file("cert-manager-policy.json")}"
}

resource "aws_iam_policy" "cert_manager_policy" {
  name        = "cert-manager-policy"
  description = "Policy for cert-manager serviceAccount"
  policy = "${data.template_file.cert_manager_policy.rendered}"
}

module "cert_manager_irsa_role" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"
  version = "5.3.1"

  role_name = "cert-manager-role"

  role_policy_arns = {
    policy = "${aws_iam_policy.cert_manager_policy.arn}"
  }

  oidc_providers = {
    ex = {
      provider_arn               = var.provider_arn
      namespace_service_accounts = ["cert-manager:cert-manager"]
    }
  }
}

resource "kubernetes_namespace" "cert_manager_namespace" {
  metadata {
    name = "cert-manager"
  }
}

resource "helm_release" "cert_manager" {
  repository            = "https://charts.jetstack.io"
  name                  = "cert-manager"
  chart                 = "cert-manager"
  version               = "1.13.3"
  values = [
    "${file("cert-manager-values.yaml")}"
  ]
  namespace = "cert-manager"

  depends_on = [
    kubernetes_namespace.cert_manager_namespace,
    helm_release.nginx_ingress,
  ]
}

resource "kubernetes_manifest" "clusterissuer" {
  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "ClusterIssuer"
    metadata   = {
      name = "letsencrypt-prod"
    }
    spec = {
      acme = {
        server             = "https://acme-v02.api.letsencrypt.org/directory"
        email              = "mstevan95@gmail.com"
        privateKeySecretRef = {
          name = "letsencrypt-prod"
        }
        solvers = [
          {
            dns01 = {
              route53 = {
                region       = "us-east-1"
                hostedZoneID = ""
              }
            }
          }
        ]
      }
    }
  }

  depends_on = [
    helm_release.cert_manager
  ]
}



