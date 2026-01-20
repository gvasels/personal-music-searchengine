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
    key            = "global/terraform.tfstate"
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

# Local values
locals {
  name_prefix = "${var.project_name}-${var.environment}"
}

# Outputs for other modules
output "state_bucket_name" {
  value = aws_s3_bucket.tofu_state.id
}

output "lock_table_name" {
  value = aws_dynamodb_table.tofu_lock.id
}

output "ecr_repository_urls" {
  value = {
    api         = aws_ecr_repository.api.repository_url
    processor   = aws_ecr_repository.processor.repository_url
    indexer     = aws_ecr_repository.indexer.repository_url
    nixiesearch = aws_ecr_repository.nixiesearch.repository_url
  }
}

output "lambda_execution_role_arn" {
  value = aws_iam_role.lambda_execution.arn
}
