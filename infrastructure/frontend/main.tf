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
    key            = "frontend/terraform.tfstate"
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

variable "custom_domain" {
  description = "Custom domain for CloudFront (e.g., music.vasels.com). Leave empty to use CloudFront default domain."
  type        = string
  default     = ""
}

variable "acm_certificate_arn" {
  description = "ACM certificate ARN for custom domain. Required if custom_domain is set."
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

data "terraform_remote_state" "backend" {
  backend = "s3"
  config = {
    bucket = "music-library-prod-tofu-state"
    key    = "backend/terraform.tfstate"
    region = "us-east-1"
  }
}

locals {
  name_prefix = "${var.project_name}-${var.environment}"
}

# Outputs
output "frontend_bucket_name" {
  description = "Frontend S3 bucket name"
  value       = aws_s3_bucket.frontend.id
}

output "frontend_bucket_arn" {
  description = "Frontend S3 bucket ARN"
  value       = aws_s3_bucket.frontend.arn
}

output "frontend_cloudfront_distribution_id" {
  description = "CloudFront distribution ID for cache invalidation"
  value       = aws_cloudfront_distribution.frontend.id
}

output "frontend_cloudfront_domain_name" {
  description = "CloudFront domain name for frontend access"
  value       = aws_cloudfront_distribution.frontend.domain_name
}

output "frontend_custom_domain" {
  description = "Custom domain for frontend (if configured)"
  value       = var.custom_domain != "" ? var.custom_domain : null
}

output "frontend_url" {
  description = "Frontend URL (custom domain or CloudFront)"
  value       = var.custom_domain != "" ? "https://${var.custom_domain}" : "https://${aws_cloudfront_distribution.frontend.domain_name}"
}
