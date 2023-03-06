output "main-vpc-id" {
  value = aws_vpc.main-vpc.id
}

output "private-subnet-id" {
  value = aws_subnet.private-subnet.id
}

output "public-subnet-id" {
  value = aws_subnet.public-subnet.id
}

output "sg-ansible-id" {
  value = aws_security_group.security_group_ansible.id
}

output "sg-wazuh-id" {
  value = aws_security_group.security_group_wazuh.id
}



