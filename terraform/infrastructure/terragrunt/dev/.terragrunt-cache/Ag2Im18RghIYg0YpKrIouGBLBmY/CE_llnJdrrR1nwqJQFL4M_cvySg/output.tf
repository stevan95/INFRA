output "vnet_id" {
  value = azurerm_virtual_network.azure_infra_network.id
}

output "subnet_ids" { value = { for subnet in azurerm_virtual_network.azure_infra_network.subnet : subnet.name => subnet.id } }
