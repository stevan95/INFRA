data "aws_caller_identity" "current" {}
data "aws_partition" "current" {}

# Policy used by private repositories
data "aws_iam_policy_document" "repository" {

  dynamic "statement" {
    
    for_each = var.repository_type == "private" ? [1] : []

    content {
      sid = "PrivateReadOnly"

      principals {
        type = "AWS"
        identifiers = coalescelist(
          concat(var.repository_read_access_arns, var.repository_read_write_access_arns),
          ["arn:${data.aws_partition.current.partition}:iam::${data.aws_caller_identity.current.account_id}:root"],
        )
      }

      actions = [
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:BatchGetImage",
        "ecr:DescribeImageScanFindings",
        "ecr:DescribeImages",
        "ecr:DescribeRepositories",
        "ecr:GetDownloadUrlForLayer",
        "ecr:GetLifecyclePolicy",
        "ecr:GetLifecyclePolicyPreview",
        "ecr:GetRepositoryPolicy",
        "ecr:ListImages",
        "ecr:ListTagsForResource",
      ]
    }
  }

  dynamic "statement" {
    for_each = length(var.repository_read_write_access_arns) > 0 ? [var.repository_read_write_access_arns] : []

    content {
      sid = "ReadWrite"

      principals {
        type        = "AWS"
        identifiers = statement.value
      }

      actions = [
        "ecr:PutImage",
        "ecr:InitiateLayerUpload",
        "ecr:UploadLayerPart",
        "ecr:CompleteLayerUpload",
      ]
    }
  }
}

#Repository
resource "aws_ecr_repository" "this" {
  name                 = var.repository_name
  image_tag_mutability = var.repository_image_tag_mutability

  force_delete = var.repository_force_delete

  image_scanning_configuration {
    scan_on_push = var.repository_image_scan_on_push
  }

  tags = var.tags
}

resource "aws_ecr_repository_policy" "this" {
  repository = aws_ecr_repository.this.name
  policy     = data.aws_iam_policy_document.repository.json 
}

resource "aws_ecr_lifecycle_policy" "this" {
  repository = aws_ecr_repository.this.name
  policy     = var.repository_lifecycle_policy
}
