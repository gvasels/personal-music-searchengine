terraform {
  required_version = ">= 1.8.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
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

variable "frontend_cloudfront_domain" {
  description = "Frontend CloudFront distribution domain for CORS (set after frontend deployment)"
  type        = string
  default     = ""
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
  name_prefix                = "${var.project_name}-${var.environment}"
  dynamodb_table_name        = data.terraform_remote_state.shared.outputs.dynamodb_table_name
  dynamodb_table_arn         = data.terraform_remote_state.shared.outputs.dynamodb_table_arn
  media_bucket_name          = data.terraform_remote_state.shared.outputs.media_bucket_name
  media_bucket_arn           = data.terraform_remote_state.shared.outputs.media_bucket_arn
  search_indexes_bucket_name = data.terraform_remote_state.shared.outputs.search_indexes_bucket_name
  search_indexes_bucket_arn  = data.terraform_remote_state.shared.outputs.search_indexes_bucket_arn
  cognito_user_pool_id       = data.terraform_remote_state.shared.outputs.cognito_user_pool_id
  lambda_role_arn            = data.terraform_remote_state.global.outputs.lambda_execution_role_arn
}

# Outputs
output "step_functions_arn" {
  description = "Upload processor state machine ARN"
  value       = aws_sfn_state_machine.upload_processor.arn
}

output "api_lambda_arn" {
  description = "API Lambda function ARN"
  value       = aws_lambda_function.api.arn
}

output "api_lambda_name" {
  description = "API Lambda function name"
  value       = aws_lambda_function.api.function_name
}

output "metadata_extractor_lambda_arn" {
  description = "Metadata extractor Lambda ARN"
  value       = aws_lambda_function.metadata_extractor.arn
}

output "cover_art_processor_lambda_arn" {
  description = "Cover art processor Lambda ARN"
  value       = aws_lambda_function.cover_art_processor.arn
}

output "track_creator_lambda_arn" {
  description = "Track creator Lambda ARN"
  value       = aws_lambda_function.track_creator.arn
}

output "file_mover_lambda_arn" {
  description = "File mover Lambda ARN"
  value       = aws_lambda_function.file_mover.arn
}

output "search_indexer_lambda_arn" {
  description = "Search indexer Lambda ARN"
  value       = aws_lambda_function.search_indexer.arn
}

output "upload_status_updater_lambda_arn" {
  description = "Upload status updater Lambda ARN"
  value       = aws_lambda_function.upload_status_updater.arn
}

output "api_gateway_id" {
  description = "API Gateway HTTP API ID"
  value       = aws_apigatewayv2_api.api.id
}

output "api_gateway_url" {
  description = "API Gateway invoke URL"
  value       = aws_apigatewayv2_stage.default.invoke_url
}

output "cognito_authorizer_id" {
  description = "API Gateway Cognito authorizer ID"
  value       = aws_apigatewayv2_authorizer.cognito.id
}

output "nixiesearch_lambda_arn" {
  description = "Nixiesearch Lambda function ARN"
  value       = aws_lambda_function.nixiesearch.arn
}

output "nixiesearch_lambda_name" {
  description = "Nixiesearch Lambda function name"
  value       = aws_lambda_function.nixiesearch.function_name
}
