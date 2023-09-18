resource "azurerm_resource_group" "azure_infra" {
  name     = "infra-azure-rg"
  location = "westus"
}

resource "azurerm_virtual_network" "azure_infra_network" {
  name                = var.vnet_name
  location            = azurerm_resource_group.azure_infra.location
  resource_group_name = azurerm_resource_group.azure_infra.name
  address_space       = var.cidr_address_space

  dynamic "subnet" {
    for_each = var.subnet_address_prefix
    content {
      name           = subnet.value["name"]
      address_prefix = subnet.value["address_prefix"]
    }
  }
}
