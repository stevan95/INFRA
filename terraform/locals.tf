locals {
  private-subnet = cidrsubnet(var.cidr_block, 8, 0)
  public-subnet  = cidrsubnet(var.cidr_block, 8, 1)
}

locals {
  common_tags = {
    Maintainer = "Stevan"
    Project    = "TestInfra"
    Resource   = "Main Cloud Infra"
  }
}

locals {
  wazuh_tags = {
    Maintainer = "Stevan"
    App        = "Wazuh"
  }
}