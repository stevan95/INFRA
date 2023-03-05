output "main-vpc-id" {
  value = aws_vpc.main-vpc.id
}

output "private-subnet-id" {
  value = aws_subnet.private-subnet.id
}

output "public-subnet-id" {
  value = aws_subnet.public-subnet.id
}