output "main-vpc-id" {
  value = aws_vpc.main-vpc.id
}

output "all_subnets" {
  value = aws_subnet.subnets.*.id
}

/*output "sg-ansible-id" {
  value = aws_security_group.security_group_ansible.id
}

output "sg-wazuh-id" {
  value = aws_security_group.security_group_wazuh.id
}*/



