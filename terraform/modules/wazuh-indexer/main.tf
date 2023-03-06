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
  subnet_id     = var.private-subnet-id
  key_name      = aws_key_pair.generated_key.key_name

  security_groups = [var.sg-wazuh-id]

  tags = var.wazuh_tags
}