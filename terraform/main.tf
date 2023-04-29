provider "aws" {
  region = "us-west-2"
}

module "ec2_wazuh_indexer" {
  source = "./modules/ec2"

  sg_name     = "Wazuh-SG"
  description = "Allow SSH from ansible host and allow trafic on port 9200"
  vpc_id      = module.vpc.main-vpc-id
  use_sg_id   = true

  security-group-ingress-settings = [
    {
      "port"              = "22"
      "protocol"          = "tcp"
      "cidr_block"        = ""
      "sg_id" = module.ec2_ansible.sg-id
    },
    {
      "port"              = "9200"
      "protocol"          = "tcp"
      "cidr_block"        = "0.0.0.0/0"
      "sg_id" = ""
    }
  ]

  security-group-egress-settings = [
    {
      "port"              = "0"
      "protocol"          = "-1"
      "cidr_block"        = "0.0.0.0/0"
    }
  ]

  ec2_configuration = {
    ami           = "ami-0ac64ad8517166fb1"
    instance_type = "t3.large"
    subnet_id     = module.vpc.all_subnets[1]
    associate_public_ip_address = false

  }

  tag_name = "Wazuh"
}

module "ec2_ansible" {
  source = "./modules/ec2"

  sg_name     = "Ansible-SG"
  description = "Allow SSH from anyware"
  vpc_id      = module.vpc.main-vpc-id

  security-group-ingress-settings = [
    {
      "port"              = "22"
      "protocol"          = "tcp"
      "cidr_block"        = "0.0.0.0/0"
      "sg_id" = ""
    }
  ]

  security-group-egress-settings = [
    {
      "port"              = "0"
      "protocol"          = "-1"
      "cidr_block"        = "0.0.0.0/0"
    }
  ]

  ec2_configuration = {
    ami           = "ami-0ac64ad8517166fb1"
    instance_type = "t3.micro"
    subnet_id     = module.vpc.all_subnets[0]
    associate_public_ip_address = true

  }

  tag_name = "Ansible"
}

module "vpc" {
  source = "./modules/vpc"

  cidr_block              = var.cidr_block
  all_subnets             = var.subnets
  enable_internet_gateway = var.enable_internet_gateway
  enable_nat_gateway      = var.enable_nat_gateway

  tag_name = "vpc-dev"
}
