locals {
  private-subnet = cidrsubnet(var.cidr_block, 8, 0)
  public-subnet  = cidrsubnet(var.cidr_block, 8, 1)
}
