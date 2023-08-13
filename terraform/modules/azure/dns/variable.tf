variable "private_zone_name" {
  type = string
  description = "Name of the private zone"
}

variable "resource_group_name" {
  type = string
  description = "Name of the resource group"
}

variable "vnet_id" {
  type = string
  description = "ID of the vnet"
}

variable "private_zone_vnet_link" {
  type = string
  description = "Private zone vnet link name"
}