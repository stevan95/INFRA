output "main-vpc-id" {
  value = aws_vpc.main-vpc.id
}

output "all_subnets" {
  value = aws_subnet.subnets.*.id
}
