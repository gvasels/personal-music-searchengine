# GitHub Actions OIDC Provider and IAM Role
# Enables keyless authentication from GitHub Actions to AWS

# GitHub OIDC Identity Provider
resource "aws_iam_openid_connect_provider" "github" {
  url             = "https://token.actions.githubusercontent.com"
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = ["6938fd4d98bab03faadb97b34396831e3780aea1", "1c58a3a8518e8759bf075b76b750d4f2df264fcd"]

  tags = {
    Name = "github-actions-oidc"
  }
}

# IAM Role for GitHub Actions
resource "aws_iam_role" "github_actions" {
  name        = "${local.name_prefix}-github-actions"
  description = "IAM role for GitHub Actions CI/CD deployments"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Federated = aws_iam_openid_connect_provider.github.arn
        }
        Action = "sts:AssumeRoleWithWebIdentity"
        Condition = {
          StringLike = {
            "token.actions.githubusercontent.com:sub" = "repo:gvasels/personal-music-searchengine:*"
          }
          StringEquals = {
            "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
          }
        }
      }
    ]
  })

  tags = {
    Name = "${local.name_prefix}-github-actions"
  }
}

# S3 Frontend Deployment Policy
resource "aws_iam_policy" "github_actions_s3_deploy" {
  name        = "${local.name_prefix}-github-actions-s3-deploy"
  description = "Allows GitHub Actions to deploy frontend assets to S3"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "S3FrontendDeployment"
        Effect = "Allow"
        Action = [
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:ListBucket",
          "s3:GetObject",
          "s3:GetBucketLocation"
        ]
        Resource = [
          "arn:aws:s3:::${local.name_prefix}-frontend",
          "arn:aws:s3:::${local.name_prefix}-frontend/*"
        ]
      }
    ]
  })
}

# CloudFront Invalidation Policy
resource "aws_iam_policy" "github_actions_cloudfront" {
  name        = "${local.name_prefix}-github-actions-cloudfront"
  description = "Allows GitHub Actions to create CloudFront invalidations"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "CloudFrontInvalidation"
        Effect = "Allow"
        Action = [
          "cloudfront:CreateInvalidation",
          "cloudfront:GetInvalidation",
          "cloudfront:ListInvalidations"
        ]
        Resource = "arn:aws:cloudfront::${data.aws_caller_identity.current.account_id}:distribution/*"
      }
    ]
  })
}

# ECR Push Policy
resource "aws_iam_policy" "github_actions_ecr" {
  name        = "${local.name_prefix}-github-actions-ecr"
  description = "Allows GitHub Actions to push container images to ECR"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "ECRGetAuthToken"
        Effect = "Allow"
        Action = [
          "ecr:GetAuthorizationToken"
        ]
        Resource = "*"
      },
      {
        Sid    = "ECRPushPull"
        Effect = "Allow"
        Action = [
          "ecr:BatchCheckLayerAvailability",
          "ecr:GetDownloadUrlForLayer",
          "ecr:BatchGetImage",
          "ecr:PutImage",
          "ecr:InitiateLayerUpload",
          "ecr:UploadLayerPart",
          "ecr:CompleteLayerUpload"
        ]
        Resource = [
          "arn:aws:ecr:${var.aws_region}:${data.aws_caller_identity.current.account_id}:repository/${local.name_prefix}-*"
        ]
      }
    ]
  })
}

# Lambda Update Policy
resource "aws_iam_policy" "github_actions_lambda" {
  name        = "${local.name_prefix}-github-actions-lambda"
  description = "Allows GitHub Actions to update Lambda function code"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "LambdaUpdateCode"
        Effect = "Allow"
        Action = [
          "lambda:UpdateFunctionCode",
          "lambda:GetFunction",
          "lambda:GetFunctionConfiguration"
        ]
        Resource = [
          "arn:aws:lambda:${var.aws_region}:${data.aws_caller_identity.current.account_id}:function:${local.name_prefix}-*"
        ]
      }
    ]
  })
}

# OpenTofu State Access Policy
resource "aws_iam_policy" "github_actions_tofu_state" {
  name        = "${local.name_prefix}-github-actions-tofu-state"
  description = "Allows GitHub Actions to access OpenTofu state"

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "S3StateAccess"
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ]
        Resource = [
          "arn:aws:s3:::${local.name_prefix}-tofu-state",
          "arn:aws:s3:::${local.name_prefix}-tofu-state/*"
        ]
      },
      {
        Sid    = "DynamoDBLockAccess"
        Effect = "Allow"
        Action = [
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:DeleteItem"
        ]
        Resource = "arn:aws:dynamodb:${var.aws_region}:${data.aws_caller_identity.current.account_id}:table/${local.name_prefix}-tofu-lock"
      }
    ]
  })
}

# Attach all policies to the GitHub Actions role
resource "aws_iam_role_policy_attachment" "github_actions_s3_deploy" {
  role       = aws_iam_role.github_actions.name
  policy_arn = aws_iam_policy.github_actions_s3_deploy.arn
}

resource "aws_iam_role_policy_attachment" "github_actions_cloudfront" {
  role       = aws_iam_role.github_actions.name
  policy_arn = aws_iam_policy.github_actions_cloudfront.arn
}

resource "aws_iam_role_policy_attachment" "github_actions_ecr" {
  role       = aws_iam_role.github_actions.name
  policy_arn = aws_iam_policy.github_actions_ecr.arn
}

resource "aws_iam_role_policy_attachment" "github_actions_lambda" {
  role       = aws_iam_role.github_actions.name
  policy_arn = aws_iam_policy.github_actions_lambda.arn
}

resource "aws_iam_role_policy_attachment" "github_actions_tofu_state" {
  role       = aws_iam_role.github_actions.name
  policy_arn = aws_iam_policy.github_actions_tofu_state.arn
}

# Data source for current AWS account
data "aws_caller_identity" "current" {}

# Outputs
output "github_actions_role_arn" {
  description = "ARN of the IAM role for GitHub Actions"
  value       = aws_iam_role.github_actions.arn
}

output "github_oidc_provider_arn" {
  description = "ARN of the GitHub OIDC identity provider"
  value       = aws_iam_openid_connect_provider.github.arn
}
