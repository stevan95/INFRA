terraform {
    source = "../../../modules/azure/vnet"

}

retryable_errors = [
  "a regex to match the error",
  "another regex"
]

retry_max_attempts = 3
retry_sleep_interval_sec = 15

include "root" {
    path = find_in_parent_folders()
}
