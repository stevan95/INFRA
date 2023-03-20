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
  default = false
}

/*variable "name" {
  type    = string
  default = "Allow SSH"
}

variable "description" {
  type    = string
  default = "Allow SSH from anywere"
}

variable "port" {
  type    = string
  default = "22"
}

variable "allow_all" {
  type    = list(string)
  default = ["0.0.0.0/0"]
}*/