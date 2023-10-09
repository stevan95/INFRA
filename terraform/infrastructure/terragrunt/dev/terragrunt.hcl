terraform {
    source = "../../../modules/azure/vnet"

    extra_arguments "conditional_vars" {
      commands = [
        "apply",
        "plan"
      ]

    optional_var_files = [
      "${get_terragrunt_dir()}/tfvars/dev1.tfvars",
      "${get_terragrunt_dir()}/tfvars/dev12.tfvars",
    ]
  }
}

retryable_errors = [
  "*"
]

retry_max_attempts = 2
retry_sleep_interval_sec = 15

include "root" {
    path = find_in_parent_folders()
}
