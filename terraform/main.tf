provider "aws" {
  region = "us-east-1"
}

module "vpc" {
  source = "./modules/vpc"

  cidr_block = var.cidr_block
  public_subnet_cidr = local.public-subnet
  private_subnet_cidr = local.private-subnet
  common_tags = local.common_tags
}

resource "aws_security_group" "allow_ssh_ansible" {
  name        = "Allow SSH"
  description = "Allow SSH inbound traffic"
  vpc_id      = module.vpc.main-vpc-id

  ingress {
    description = "Allow SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.common_tags
}

resource "aws_instance" "ansible-host" {
  ami           = "ami-006dcf34c09e50022"
  instance_type = "t3.micro"
  subnet_id     = module.vpc.public-subnet-id

  associate_public_ip_address = true

  user_data = <<-EOF
  #!/bin/bash 

  sudo apt update 
  sudo apt-add-repository -y ppa:ansible/ansible
  sudo apt-get -y install ansible
  ansible --version

  sudo apt update
  sudo apt install python3-boto3
  EOF

  tags = local.common_tags
}

resource "aws_network_interface_sg_attachment" "sg_attachment-ansible" {
  security_group_id    = aws_security_group.allow_ssh_ansible.id
  network_interface_id = aws_instance.ansible-host.primary_network_interface_id
}

resource "tls_private_key" "wazuh-ssh" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "generated_key" {
  key_name   = "wazuh-key"
  public_key = tls_private_key.wazuh-ssh.public_key_openssh
}

resource "aws_instance" "wazuh-host" {
  ami           = "ami-006dcf34c09e50022"
  instance_type = "t3.medium"
  subnet_id     = module.vpc.private-subnet-id
  key_name      = aws_key_pair.generated_key.key_name

  tags = local.common_tags
}

resource "aws_security_group" "allow_ssh_wazuh" {
  name        = "allow_tls"
  description = "Allow SSH from Ansible Host"
  vpc_id      = module.vpc.main-vpc-id

  ingress {
    description     = "Allow SSH"
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    security_groups = [aws_security_group.allow_ssh_ansible.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = local.common_tags
}

resource "aws_network_interface_sg_attachment" "sg_attachment-wazuh" {
  security_group_id    = aws_security_group.allow_ssh_wazuh.id
  network_interface_id = aws_instance.wazuh-host.primary_network_interface_id
}