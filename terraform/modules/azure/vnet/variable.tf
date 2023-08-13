variable "vnet_name" {
  type = string
  description = "Virtual Network name"
}

variable "resource_group_name" {
  type = string
  description = "Name of the resource group"
}

variable "resource_group_location" {
  type = string
  description = "Region where resource group is located"
}

variable "cidr_address_space" {
  type = list(string)
  description = "CIDR block range"
}

variable "subnet_address_prefix" {
  type = list(object({
    name = string
    address_prefix = string
  }))
  description = "Names of subnets with their prefixes"
}
