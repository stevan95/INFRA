resource "azurerm_virtual_network" "azure_infra_network" {
  name                = var.vnet_name
  location            = var.resource_group_location
  resource_group_name = var.resource_group_name
  address_space       = var.cidr_address_space

  dynamic "subnet" {
    for_each = var.subnet_address_prefix
    content {
      name           = subnet.value["name"]
      address_prefix = subnet.value["address_prefix"]
    }
  }
}
