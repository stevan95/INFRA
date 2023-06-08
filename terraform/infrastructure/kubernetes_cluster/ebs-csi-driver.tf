data "template_file" "ebs-csi-driver-policy" {
  template = "${file("ebs-csi-driver-policy.json")}"
}

resource "aws_iam_policy" "ebs-csi-driver-policy" {
  name        = "ebs-csi-driver-policy"
  description = "Policy for ebs-csi-controller-sa"
  policy = "${data.template_file.ebs-csi-driver-policy.rendered}"
} 

module "aws_ebs_csi_driver_irsa_role" {
  source  = "terraform-aws-modules/iam/aws//modules/iam-role-for-service-accounts-eks"
  version = "5.3.1"

  role_name = "aws-ebs-csi-driver-controller"

  role_policy_arns = {
    policy = "${aws_iam_policy.ebs-csi-driver-policy.arn}"
  }

  oidc_providers = {
    ex = {
      provider_arn               = module.eks.oidc_provider_arn
      namespace_service_accounts = ["kube-system:ebs-csi-controller-sa"]
    }
  }
}

resource "helm_release" "aws-ebs-csi-driver" {
  repository            = "https://kubernetes-sigs.github.io/aws-ebs-csi-driver"
  name                  = "aws-ebs-csi-driver"
  chart                 = "aws-ebs-csi-driver"
  version               = "2.18.0"
  values = [
    "${file("values.yaml")}"
  ]
  namespace = "kube-system"
}

resource "kubernetes_storage_class" "aws-ebs-csi-driver" {
  metadata {
    name = "cloud-ssd"
  }
  
  storage_provisioner    = "kubernetes.io/aws-ebs"
  volume_binding_mode    = "WaitForFirstConsumer"
  reclaim_policy         = "Retain"
  allow_volume_expansion = true

  parameters = {
    type = "gp2"
  }
}
