terraform {
    source = "../../../modules/azure/vnet"

}



include "root" {
    path = find_in_parent_folders()
}
