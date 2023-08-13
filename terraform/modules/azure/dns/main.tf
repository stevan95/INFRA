resource "azurerm_private_dns_zone" "private_zone" {
  name                = var.private_zone_name
  resource_group_name = var.resource_group_name
}

resource "azurerm_private_dns_zone_virtual_network_link" "vnet_link" {
  name                  = var.private_zone_name
  resource_group_name   = var.resource_group_name
  private_dns_zone_name = azurerm_private_dns_zone.private_zone.name
  virtual_network_id    = var.vnet_id
}

resource "azurerm_private_dns_cname_record" "test_record" {
  name                = "google"
  zone_name           = azurerm_private_dns_zone.private_zone.name
  resource_group_name = var.resource_group_name
  ttl                 = 300
  record              = "www.google.com"
}
