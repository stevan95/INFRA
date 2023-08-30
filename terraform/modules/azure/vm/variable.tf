variable "vm_eni_name" {
  type = string
  description = "Name for the virtual interface"
}

variable "resource_group_name" {
  type = string
  description = "Name of the resource group"
}

variable "resource_group_location" {
  type = string
  description = "Region where resource group is located"
}

variable "azurerm_subnet_id" {
  type = string
  description = "Azure Subnet ID"
}

variable "vm_name" {
  type = string
  description = "Name of the virtual machine"
}

variable "create_public_ip" {
  type = bool
  description = "Set to true if you want to attach public ip to your vm"
  default = false  
}
