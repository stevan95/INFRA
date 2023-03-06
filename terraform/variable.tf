variable "cidr_block" {
  type    = string
  default = "10.0.0.0/16"
}

variable "name" {
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
}