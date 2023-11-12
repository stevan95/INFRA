include "root" {
    path = find_in_parent_folders()
}

remote_state {
  backend = "s3"
  generate = {
    path      = "backend.tf"
    if_exists = "overwrite"
  }
  config = {
    bucket         = "stevan-terraform-state-bucket"
    key            = "${path_relative_to_include()}/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "lock-table"
  }
}

terraform {
  source = "../../../../modules/k8s_controllers"
}

dependency "vpc" {
  config_path = "../k8s_network"

  mock_outputs = {
    vpc_id = "id-mock-8734532123"
    private_subnet_ids = ["subnet-12345678", "subnet-87654321"]
    public_subnet_ids = ["subnet-4567890", "subnet-0453325"]
  }
}

dependency "eks_cluster" {
  config_path = "../k8s_cluster"

  mock_outputs = {
    cluster_id = "randomid"
    endpoint = "testendpoint" 
    kubeconfig-certificate-authority-data = "test-cert"
    oidc-issuer = "oidc-test-issuer"
    provider_arn = "provider-test-arn"
  }
}

inputs = {
  provider_arn = dependency.eks_cluster.outputs.provider_arn
  cluster_id = dependency.eks_cluster.outputs.cluster_id
  kubeconfig_cert = dependency.eks_cluster.outputs.kubeconfig-certificate-authority-data
  endpoint = dependency.eks_cluster.outputs.endpoint
}

retryable_errors = [
  "*"
]

retry_max_attempts = 2
retry_sleep_interval_sec = 10
  