provider "aws" {
  region = "us-east-1"
}

/*module "wazuh-indexer" {
  source            = "./modules/wazuh-indexer"
  vpc-id            = module.vpc.main-vpc-id
  private-subnet-id = module.vpc.private-subnet-id
  sg-wazuh-id       = module.vpc.sg-wazuh-id

  wazuh_tags = local.wazuh_tags
}

module "ansible_ec2" {
  source           = "./modules/ansible_ec2"
  vpc-id           = module.vpc.main-vpc-id
  public-subnet-id = module.vpc.public-subnet-id
  sg-ansible-id    = module.vpc.sg-ansible-id

  wazuh_tags = local.wazuh_tags
}*/

module "security-group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "4.17.1"

  name        = "Ansible-SG"
  description = "Allows SSH into port 22."
  vpc_id      = module.vpc.main-vpc-id

  ingress_with_cidr_blocks = [

    {
      from_port   = 22
      to_port     = 22
      protocol    = "tcp"
      description = "Allow SSH to ansible master node from anyware."
      cidr_blocks = "0.0.0.0/0"
    }
  ]

  egress_with_cidr_blocks = [
    {
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      cidr_blocks = "0.0.0.0/0"
    }
  ]

}

module "vpc" {
  source = "./modules/vpc"

  cidr_block              = var.cidr_block
  all_subnets             = var.subnets
  enable_internet_gateway = var.enable_internet_gateway
  enable_nat_gateway      = var.enable_nat_gateway

  /*name                    = var.name
  port                    = var.port
  description             = var.description
  allow_all               = var.allow_all
  wazuh_tags              = local.wazuh_tags*/
  common_tags = local.common_tags
}
