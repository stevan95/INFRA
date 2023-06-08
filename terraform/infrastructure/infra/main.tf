module "ecr" {
  source = "../../modules/ecr"

  repository_name = "stevan-repo"

  repository_read_write_access_arns = ["arn:aws:iam::022865527167:role/ecr_readwrite_role"]
  repository_lifecycle_policy = jsonencode({
    rules = [
      {
        rulePriority = 1,
        description  = "Keep last 30 images",
        selection = {
          tagStatus     = "tagged",
          tagPrefixList = ["v"],
          countType     = "imageCountMoreThan",
          countNumber   = 30
        },
        action = {
          type = "expire"
        }
      }
    ]
  })

  tags = {
    Terraform   = "true"
    Environment = "dev"
  }
}