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

variable "tag_name" {
  type = string
}
