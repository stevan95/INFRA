provider "aws" {
  region = "us-east-1"
}

module "wazuh-indexer" {
  source            = "./modules/wazuh-indexer"
  vpc-id            = module.vpc.main-vpc-id
  private-subnet-id = module.vpc.all_subnets[1]
  sg-wazuh-id       = module.vpc.sg-wazuh-id

  wazuh_tags = local.wazuh_tags
}

module "ansible_ec2" {
  source           = "./modules/ansible_ec2"
  vpc-id           = module.vpc.main-vpc-id
  public-subnet-id = module.vpc.all_subnets[0]
  sg-ansible-id    = module.vpc.sg-ansible-id

  wazuh_tags = local.wazuh_tags
}

module "vpc" {
  source = "./modules/vpc"

  cidr_block              = var.cidr_block
  all_subnets             = var.subnets
  enable_internet_gateway = var.enable_internet_gateway
  enable_nat_gateway      = var.enable_nat_gateway

  name                    = var.name
  port                    = var.port
  description             = var.description
  allow_all               = var.allow_all
  wazuh_tags              = local.wazuh_tags
  common_tags = local.common_tags
}
