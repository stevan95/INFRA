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
