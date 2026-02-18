terraform {
  required_version = ">= 1.8.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket         = "music-library-prod-tofu-state"
    key            = "shared/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "music-library-prod-tofu-lock"
    encrypt        = true
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "personal-music-searchengine"
      Environment = var.environment
      ManagedBy   = "opentofu"
    }
  }
}

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "prod"
}

variable "project_name" {
  description = "Project name for resource naming"
  type        = string
  default     = "music-library"
}

variable "cognito_extra_callback_urls" {
  description = "Additional Cognito callback URLs beyond auto-computed ones"
  type        = list(string)
  default     = []
}

variable "cognito_extra_logout_urls" {
  description = "Additional Cognito logout URLs beyond auto-computed ones"
  type        = list(string)
  default     = []
}

variable "custom_domain" {
  description = "Custom domain for frontend (e.g., music.vasels.com). Used for S3 CORS origins."
  type        = string
  default     = ""
}

data "terraform_remote_state" "frontend" {
  backend = "s3"
  config = {
    bucket = "music-library-prod-tofu-state"
    key    = "frontend/terraform.tfstate"
    region = "us-east-1"
  }
}

locals {
  name_prefix = "${var.project_name}-${var.environment}"

  # Dynamic CORS origins from frontend outputs
  frontend_cf_domain     = try(data.terraform_remote_state.frontend.outputs.frontend_cloudfront_domain_name, "")
  frontend_custom_domain = try(data.terraform_remote_state.frontend.outputs.frontend_custom_domain, var.custom_domain)

  s3_cors_origins = compact([
    "http://localhost:5173",
    "http://localhost:3000",
    local.frontend_cf_domain != "" ? "https://${local.frontend_cf_domain}" : "",
    local.frontend_custom_domain != null && local.frontend_custom_domain != "" ? "https://${local.frontend_custom_domain}" : "",
  ])

  # Dynamic Cognito callback/logout URLs
  cognito_callback_urls = distinct(compact(concat(
    ["http://localhost:5173/callback"],
    [local.frontend_cf_domain != "" ? "https://${local.frontend_cf_domain}/callback" : ""],
    [local.frontend_custom_domain != null && local.frontend_custom_domain != "" ? "https://${local.frontend_custom_domain}/callback" : ""],
    var.cognito_extra_callback_urls,
  )))

  cognito_logout_urls = distinct(compact(concat(
    ["http://localhost:5173"],
    [local.frontend_cf_domain != "" ? "https://${local.frontend_cf_domain}" : ""],
    [local.frontend_custom_domain != null && local.frontend_custom_domain != "" ? "https://${local.frontend_custom_domain}" : ""],
    var.cognito_extra_logout_urls,
  )))
}

# Outputs
output "cognito_user_pool_id" {
  value = aws_cognito_user_pool.main.id
}

output "cognito_user_pool_arn" {
  value = aws_cognito_user_pool.main.arn
}

output "cognito_client_id" {
  value = aws_cognito_user_pool_client.web.id
}

output "cognito_domain" {
  value = aws_cognito_user_pool_domain.main.domain
}

output "dynamodb_table_name" {
  value = aws_dynamodb_table.music_library.name
}

output "dynamodb_table_arn" {
  value = aws_dynamodb_table.music_library.arn
}

output "media_bucket_name" {
  value = aws_s3_bucket.media.id
}

output "media_bucket_arn" {
  value = aws_s3_bucket.media.arn
}

output "search_indexes_bucket_name" {
  value = aws_s3_bucket.search_indexes.id
}

output "search_indexes_bucket_arn" {
  value = aws_s3_bucket.search_indexes.arn
}

# Cognito User Groups
output "cognito_admin_group_name" {
  value = aws_cognito_user_group.admin.name
}

output "cognito_artist_group_name" {
  value = aws_cognito_user_group.artist.name
}

output "cognito_subscriber_group_name" {
  value = aws_cognito_user_group.subscriber.name
}
