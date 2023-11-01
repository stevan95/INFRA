variable "kubernetes_cluster_name" {
  type = string
  description = "Name of the kubernetes cluster"
}

variable "kubernetescluster_version" {
  type = string
  description = "Version of the kubernetes cluster"
}

variable "vpc_subnets" {
  type = list(string)
  description = "set of strings which represent subnet ids"
}

variable "enviroment" {
  type = string
  description = "Enviroment where cluster should be deployyed"
}

variable "node_group_name" {
  type = string
  description = "Name of worker nodes group"
}

variable "instance_types" {
  type = list(string)
  description = "set of strings which represent subnet ids"
}

variable "capacity_type" {
  type = string
  description = "Capacity type of ec2 instance (ON_DEMAND, SPOT)"
  default = "ON_DEMAND"
}

variable "disk_size" {
  type = number
  description = "Size of disk default is 40Gb"
  default = 40
}

variable "identity_provider_config_name" {
  type = string
  description = "Name of identity provider used for eks cluster"
}

