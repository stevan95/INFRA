output "private_key" {
  value     = tls_private_key.wazuh-ssh.private_key_pem
  sensitive = true
}