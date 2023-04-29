variable "vpc_id" {
  type = string
}

variable "description" {
  type = string
}

variable "sg_name" {
  type = string
}

variable "tag_name" {
  type = string
}

variable "security-group-ingress-settings" {
  type = list(object({
    port = number
    protocol = string
    cidr_block = string
    sg_id = string
  }))
}

variable "security-group-egress-settings" {
  type = list(object({
    port = number
    protocol = string
    cidr_block = string
  }))
}

variable "ec2_configuration" {
  type = object({
    ami = string
    instance_type = string
    subnet_id = string
    associate_public_ip_address = bool
  })
}
