provider "aws" {
  region = "us-east-1"
}

module "wazuh-indexer" {
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
}

/*resource "aws_network_interface_sg_attachment" "sg_attachment-ansible" {
  security_group_id    = aws_security_group.allow_ssh_ansible.id
  network_interface_id = aws_instance.ansible-host.primary_network_interface_id
}*/

module "vpc" {
  source = "./modules/vpc"

  name                = var.name
  port                = var.port
  description         = var.description
  allow_all           = var.allow_all
  wazuh_tags          = local.wazuh_tags
  cidr_block          = var.cidr_block
  public_subnet_cidr  = local.public-subnet
  private_subnet_cidr = local.private-subnet
  common_tags         = local.common_tags
}

/*resource "aws_network_interface_sg_attachment" "sg_attachment-wazuh" {
  security_group_id    = aws_security_group.allow_ssh_wazuh.id
  network_interface_id = aws_instance.wazuh-host.primary_network_interface_id
}*/