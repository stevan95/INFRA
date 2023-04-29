variable "cidr_block" {
  type    = string
  default = "10.0.0.0/16"
}

variable "subnets" {
  type    = list(any)
  default = ["10.0.0.0/24", "10.0.1.0/24"]
}

variable "enable_internet_gateway" {
  type    = bool
  default = true
}

variable "enable_nat_gateway" {
  type    = bool
  default = true
}
