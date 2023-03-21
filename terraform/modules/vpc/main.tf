data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "main-vpc" {
  cidr_block       = var.cidr_block
  instance_tenancy = "default"

  tags = var.common_tags
}

resource "aws_subnet" "public-subnets" {
  count = length(var.public_subnets)
  vpc_id     = aws_vpc.main-vpc.id
  cidr_block = var.public_subnets[count.index]

  availability_zone = data.aws_availability_zones.available.names[count.index]

  tags = var.common_tags
}

resource "aws_subnet" "private-subnets" {
  count = length(var.private_subnets)
  vpc_id     = aws_vpc.main-vpc.id
  cidr_block = var.private_subnets[count.index]

  availability_zone = data.aws_availability_zones.available.names[count.index]

  tags = var.common_tags
}

resource "aws_internet_gateway" "public-gateway" {
  count = var.enable_internet_gateway ? 1 : 0

  vpc_id = aws_vpc.main-vpc.id 

  tags = var.common_tags
}

resource "aws_eip" "nat-gateway-eip" {
  count = var.enable_nat_gateway ? 1 : 0

  vpc = true

  tags = var.common_tags
}

resource "aws_nat_gateway" "nat-gateway" {
  count = var.enable_nat_gateway ? 1 : 0

  allocation_id = aws_eip.nat-gateway-eip[0].id
  subnet_id     = var.public_subnets[0].id

  depends_on = [aws_internet_gateway.public-gateway]

  tags = var.common_tags
}

resource "aws_route_table" "public-vpc-route" {
  count = var.enable_internet_gateway ? 1 : 0
  vpc_id = aws_vpc.main-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.public-gateway[0].id
  }

  tags = var.common_tags
}

resource "aws_route_table" "private-vpc-route" {
  count = var.enable_nat_gateway ? 1 : 0

  vpc_id = aws_vpc.main-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_nat_gateway.nat-gateway[0].id
  }

  tags = var.common_tags
}

resource "aws_route_table_association" "public" {
  count = var.enable_internet_gateway ? 1 : 0

  depends_on     = [aws_subnet.public-subnets]
  route_table_id = aws_route_table.public-vpc-route[0].id
  subnet_id      = aws_subnet.public-subnets[0].id
}

resource "aws_route_table_association" "private" {
    count = var.enable_nat_gateway ? 1 : 0

  depends_on     = [aws_subnet.private-subnets]
  route_table_id = aws_route_table.private-vpc-route[0].id
  subnet_id      = aws_subnet.private-subnets[0].id
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