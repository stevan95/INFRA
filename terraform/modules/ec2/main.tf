locals {
  tags = {
    Name = var.tag_name
  }
}

resource "aws_security_group" "sg" {
  name        = var.sg_name
  description = var.description
  vpc_id      = var.vpc_id

  tags = local.tags
}

resource "aws_security_group_rule" "ingress-rules" {
    count = length(var.security-group-ingress-settings)
    type = "ingress"
    security_group_id = aws_security_group.sg.id
    from_port = var.security-group-ingress-settings[count.index].port
    to_port = var.security-group-ingress-settings[count.index].port
    protocol = var.security-group-ingress-settings[count.index].protocol
    cidr_blocks = var.security-group-ingress-settings[count.index].cidr_block == "" ? null : [var.security-group-ingress-settings[count.index].cidr_block]
    source_security_group_id = var.security-group-ingress-settings[count.index].sg_id == "" ? null :var.security-group-ingress-settings[count.index].sg_id 
}

resource "aws_security_group_rule" "egress-rules" {
    count = length(var.security-group-egress-settings)
    type = "egress"
    security_group_id = aws_security_group.sg.id
    from_port = var.security-group-egress-settings[count.index].port
    to_port = var.security-group-egress-settings[count.index].port
    protocol = var.security-group-egress-settings[count.index].protocol
    cidr_blocks = [var.security-group-egress-settings[count.index].cidr_block]  
}

resource "aws_instance" "ec2_instance" {
  ami           = var.ec2_configuration.ami
  instance_type = var.ec2_configuration.instance_type
  subnet_id     = var.ec2_configuration.subnet_id
  associate_public_ip_address = var.ec2_configuration.associate_public_ip_address
  key_name = "ssh_key_stevan"

  security_groups = [aws_security_group.sg.id]

  tags = local.tags
}
