locals {
  tags = {
    Name = var.tag_name
  }
}

data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "main-vpc" {
  cidr_block       = var.cidr_block
  instance_tenancy = "default"

  tags = local.tags
}

resource "aws_subnet" "subnets" {
  count = length(var.all_subnets)
  vpc_id     = aws_vpc.main-vpc.id
  cidr_block = var.all_subnets[count.index]

  availability_zone = data.aws_availability_zones.available.names[count.index]

  tags = local.tags
}

resource "aws_internet_gateway" "public-gateway" {
  count = var.enable_internet_gateway ? 1 : 0

  vpc_id = aws_vpc.main-vpc.id 

  tags = local.tags
}

resource "aws_eip" "nat-gateway-eip" {
  count = var.enable_nat_gateway ? 1 : 0

  vpc = true

  tags = local.tags
}

resource "aws_nat_gateway" "nat-gateway" {
  count = var.enable_nat_gateway ? 1 : 0

  allocation_id = aws_eip.nat-gateway-eip[0].id
  subnet_id     = aws_subnet.subnets[0].id

  depends_on = [aws_internet_gateway.public-gateway]

  tags = local.tags
}

resource "aws_route_table" "public-vpc-route" {
  count = var.enable_internet_gateway ? 1 : 0
  vpc_id = aws_vpc.main-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.public-gateway[0].id
  }

  tags = local.tags
}

resource "aws_route_table" "private-vpc-route" {
  count = var.enable_nat_gateway ? 1 : 0

  vpc_id = aws_vpc.main-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_nat_gateway.nat-gateway[0].id
  }

  tags = local.tags
}

resource "aws_route_table_association" "public" {
  count = var.enable_internet_gateway ? 1 : 0

  depends_on     = [aws_subnet.subnets[0]]
  route_table_id = aws_route_table.public-vpc-route[0].id
  subnet_id      = aws_subnet.subnets[0].id
}

resource "aws_route_table_association" "private" {
    count = var.enable_nat_gateway ? 1 : 0

  depends_on     = [aws_subnet.subnets[1]]
  route_table_id = aws_route_table.private-vpc-route[0].id
  subnet_id      = aws_subnet.subnets[1].id
}