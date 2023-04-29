provider "aws" {
  region = "us-west-2"
}

module "ec2_wazuh_indexer" {
  source = "./modules/ec2"

  sg_name     = "Wazuh-SG"
  description = "Allow SSH from ansible host and allow trafic on port 9200"
  vpc_id      = module.vpc.vpc_id

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
    subnet_id     = module.vpc.private_subnet_ids[0]
    associate_public_ip_address = false

  }

  tag_name = "Wazuh"
}

module "ec2_ansible" {
  source = "./modules/ec2"

  sg_name     = "Ansible-SG"
  description = "Allow SSH from anyware"
  vpc_id      = module.vpc.vpc_id

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
    subnet_id     = module.vpc.public_subnet_ids[0]
    associate_public_ip_address = true

  }

  tag_name = "Ansible"
}

module "vpc" {
  source = "./modules/vpc"

  env             = "dev"
  azs             = ["us-east-1a", "us-east-1b"]
  private_subnets = ["10.0.0.0/19", "10.0.32.0/19"]
  public_subnets  = ["10.0.64.0/19", "10.0.96.0/19"]

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = 1
    "kubernetes.io/cluster/dev-demo"  = "owned"
  }

  public_subnet_tags = {
    "kubernetes.io/role/elb"         = 1
    "kubernetes.io/cluster/dev-demo" = "owned"
  }
}
