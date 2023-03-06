resource "aws_instance" "ansible-host" {
  ami           = "ami-006dcf34c09e50022"
  instance_type = "t3.micro"
  subnet_id     = var.public-subnet-id

  associate_public_ip_address = true

  security_groups = [var.sg-ansible-id]

  user_data = <<-EOF
  #!/bin/bash 

  sudo apt update 
  sudo apt-add-repository -y ppa:ansible/ansible
  sudo apt-get -y install ansible
  ansible --version

  sudo apt update
  sudo apt install python3-boto3
  EOF

  tags = var.wazuh_tags
}