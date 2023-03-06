data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "main-vpc" {
  cidr_block       = var.cidr_block
  instance_tenancy = "default"

  tags = var.common_tags
}

resource "aws_subnet" "private-subnet" {
  vpc_id     = aws_vpc.main-vpc.id
  cidr_block = var.private_subnet_cidr

  availability_zone = data.aws_availability_zones.available.names[0]

  tags = var.common_tags
}

resource "aws_subnet" "public-subnet" {
  vpc_id     = aws_vpc.main-vpc.id
  cidr_block = var.public_subnet_cidr

  availability_zone = data.aws_availability_zones.available.names[1]

  tags = var.common_tags
}

resource "aws_internet_gateway" "public-gateway" {
  vpc_id = aws_vpc.main-vpc.id

  tags = var.common_tags
}

resource "aws_eip" "nat-gateway-eip" {
  vpc = true

  tags = var.common_tags
}

resource "aws_nat_gateway" "nat-gateway" {
  allocation_id = aws_eip.nat-gateway-eip.id
  subnet_id     = aws_subnet.public-subnet.id

  depends_on = [aws_internet_gateway.public-gateway]

  tags = var.common_tags
}

resource "aws_route_table" "public-vpc-route" {
  vpc_id = aws_vpc.main-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.public-gateway.id
  }

  tags = var.common_tags
}

resource "aws_route_table" "private-vpc-route" {
  vpc_id = aws_vpc.main-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_nat_gateway.nat-gateway.id
  }

  tags = var.common_tags
}

resource "aws_route_table_association" "public" {
  depends_on     = [aws_subnet.public-subnet]
  route_table_id = aws_route_table.public-vpc-route.id
  subnet_id      = aws_subnet.public-subnet.id
}
resource "aws_route_table_association" "private" {
  depends_on     = [aws_subnet.private-subnet]
  route_table_id = aws_route_table.private-vpc-route.id
  subnet_id      = aws_subnet.private-subnet.id
}

resource "aws_security_group" "security_group_ansible" {
  name        = "Ansible-SG"
  description = var.description
  vpc_id      = aws_vpc.main-vpc.id

  ingress {
    description = var.description
    from_port   = var.port
    to_port     = var.port
    protocol    = "tcp"
    cidr_blocks = tolist(var.allow_all)
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = var.wazuh_tags
}

resource "aws_security_group" "security_group_wazuh" {
  name        = "Wazuh-SG"
  description = var.description
  vpc_id      = aws_vpc.main-vpc.id

  ingress {
    description = var.description
    from_port   = var.port
    to_port     = var.port
    protocol    = "tcp"
    security_groups = [aws_security_group.security_group_ansible.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = var.wazuh_tags
}