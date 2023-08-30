output "public_ip_address" {
  value = azurerm_public_ip.azure_vm_pip[0].ip_address
}
