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
    key            = "backend/terraform.tfstate"
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

# Data sources for shared resources
data "terraform_remote_state" "shared" {
  backend = "s3"
  config = {
    bucket = "music-library-prod-tofu-state"
    key    = "shared/terraform.tfstate"
    region = "us-east-1"
  }
}

data "terraform_remote_state" "global" {
  backend = "s3"
  config = {
    bucket = "music-library-prod-tofu-state"
    key    = "global/terraform.tfstate"
    region = "us-east-1"
  }
}

locals {
  name_prefix          = "${var.project_name}-${var.environment}"
  dynamodb_table_name  = data.terraform_remote_state.shared.outputs.dynamodb_table_name
  dynamodb_table_arn   = data.terraform_remote_state.shared.outputs.dynamodb_table_arn
  media_bucket_name    = data.terraform_remote_state.shared.outputs.media_bucket_name
  media_bucket_arn     = data.terraform_remote_state.shared.outputs.media_bucket_arn
  cognito_user_pool_id = data.terraform_remote_state.shared.outputs.cognito_user_pool_id
  lambda_role_arn      = data.terraform_remote_state.global.outputs.lambda_execution_role_arn
}

# Outputs
output "step_functions_arn" {
  value = aws_sfn_state_machine.upload_processor.arn
}

output "api_lambda_arn" {
  value = aws_lambda_function.api.arn
}

output "metadata_extractor_lambda_arn" {
  value = aws_lambda_function.metadata_extractor.arn
}

output "search_indexer_lambda_arn" {
  value = aws_lambda_function.search_indexer.arn
}
