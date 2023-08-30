//Create Public IP
resource "azurerm_public_ip" "azure_vm_pip" {
  count               = var.create_public_ip ? 1 : 0
  name                = "${var.vm_name}_pip"
  resource_group_name = var.resource_group_name
  location            = var.resource_group_location
  allocation_method   = "Static"

  tags = {
    environment = "Test"
  }
}

//Crete ENI for Virtual Machine
resource "azurerm_network_interface" "vm_nic" {
  name                = "${var.vm_eni_name}-nic"
  location            = var.resource_group_location
  resource_group_name = var.resource_group_name

  ip_configuration {
    name                          = "${var.vm_eni_name}-nic"
    subnet_id                     = var.azurerm_subnet_id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = var.create_public_ip ? azurerm_public_ip.azure_vm_pip[0].id : null
  }
}

//Create Azure VM
resource "azurerm_virtual_machine" "azure_vm" {
  name                             = "${var.vm_name}"
  location                         = var.resource_group_location
  resource_group_name              = var.resource_group_name
  network_interface_ids            = [azurerm_network_interface.vm_nic.id]
  vm_size                          = "Standard_DS1_v2"
  delete_os_disk_on_termination    = true
  delete_data_disks_on_termination = true

  storage_image_reference {
    publisher = "Canonical"
    offer     = "0001-com-ubuntu-server-focal"
    sku       = "20_04-lts"
    version   = "latest"
  }

  storage_os_disk {
    name              = "${var.vm_name}_disk"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name  = var.vm_name
    admin_username = "testadmin"
    admin_password = "password123!"
  }

  os_profile_linux_config {
    disable_password_authentication = false
  }

  tags = {
    environment = "Test"
  }
}
