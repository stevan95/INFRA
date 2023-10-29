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
  source = "../../../../modules/eks"
}

dependency "vpc" {
  config_path = "../k8s_network"

  mock_outputs = {
    vpc_id = "id-mock-8734532123"
    private_subnet_ids = ["subnet-12345678", "subnet-87654321"]
    public_subnet_ids = ["subnet-4567890", "subnet-0453325"]
  }
}

inputs = {
    kubernetes_cluster_name = "dev_cluster01"
    kubernetescluster_version = "1.27"
    enviroment = "dev"
    vpc_subnets = dependency.vpc.outputs.private_subnet_ids
}

retryable_errors = [
  "*"
]

retry_max_attempts = 2
retry_sleep_interval_sec = 10
  