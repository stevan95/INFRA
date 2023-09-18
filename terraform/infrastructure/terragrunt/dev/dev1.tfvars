vnet_name = "azure-infra-network1"
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