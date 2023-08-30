resource "azurerm_resource_group" "azure_infra" {
  name     = "infra-azure-rg"
  location = "westus"
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

module "private_dns_zone" {
  source = "../../modules/azure/dns"

  private_zone_name   = "stevantest.org"
  resource_group_name = azurerm_resource_group.azure_infra.name
  vnet_id             = module.vnet.vnet_id
  private_zone_vnet_link = "stevaninfralink"
}

module "vm_instance" {
  source = "../../modules/azure/vm"

  vm_eni_name             = "testvm"
  resource_group_name     = azurerm_resource_group.azure_infra.name
  resource_group_location = azurerm_resource_group.azure_infra.location
  azurerm_subnet_id       = module.vnet.subnet_ids["infra-subnet-1"]
  vm_name                 = "testvm"
  create_public_ip        = true
}
