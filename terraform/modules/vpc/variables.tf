variable "cidr_block" {
    type = any
}

variable "all_subnets" {
    type = any
}

variable "enable_internet_gateway" {
    type = bool
}

variable "enable_nat_gateway" {
    type = bool
}

variable "common_tags" {
    type = any
}

variable "name" {
    type = any
}

variable "description" {
    type = any
}

variable "port" {
    type = any
}

variable "allow_all" {
    type = any
}

variable "wazuh_tags" {
    type = any
}
