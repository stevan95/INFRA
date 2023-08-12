resource "azurerm_resource_group" "azure_infra" {
  name     = "infra-azure-rg"
  location = "East US"
}

module "vnet" {
  source = "../../modules/azure/vnet"

  vnet_name = "azure-infra-network"
  resource_group_name = azurerm_resource_group.azure_infra.name
  resource_group_location = azurerm_resource_group.azure_infra.location
  cidr_address_space = ["10.0.0.0/16"]

  subnet_address_prefix = [
    {
        "name"           = "infra-subnet-1"
        "address_prefix" = "10.0.1.0/24"
    },
    {
        "name"           = "infra-subnet-2"
        "address_prefix" = "10.0.2.0/24"
    }
  ]

}
