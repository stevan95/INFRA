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
  source = "../../../modules/vpc"
}

inputs = {
  env = "dev"
  vpc_cidr_block = "10.0.0.0/16"
  azs = ["us-east-1a", "us-east-1b"]

  private_subnets = ["10.0.0.0/19", "10.0.32.0/19"]
  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = "1"
  }

  public_subnets = ["10.0.64.0/19", "10.0.96.0/19"]
  public_subnet_tags = {
    "kubernetes.io/role/elb" = "1"
  }
}

retryable_errors = [
  "*"
]

retry_max_attempts = 3
retry_sleep_interval_sec = 10
  