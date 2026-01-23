# CloudFront Distribution for Media Streaming
# Serves HLS streams and original files with signed URLs

# CloudFront Origin Access Control for S3
resource "aws_cloudfront_origin_access_control" "media" {
  name                              = "${local.name_prefix}-media-oac"
  description                       = "OAC for media bucket"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

# CloudFront Distribution
resource "aws_cloudfront_distribution" "media" {
  enabled             = true
  is_ipv6_enabled     = true
  comment             = "Media distribution for ${local.name_prefix}"
  default_root_object = ""
  price_class         = "PriceClass_100" # US, Canada, Europe

  # S3 Origin
  origin {
    domain_name              = "${local.media_bucket_name}.s3.${var.aws_region}.amazonaws.com"
    origin_id                = "S3-${local.media_bucket_name}"
    origin_access_control_id = aws_cloudfront_origin_access_control.media.id
  }

  # Default behavior - require signed URLs
  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "S3-${local.media_bucket_name}"

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 86400
    max_ttl                = 31536000

    # Require signed URLs for default
    trusted_key_groups = [aws_cloudfront_key_group.signing.id]
  }

  # HLS streaming behavior - optimized caching
  ordered_cache_behavior {
    path_pattern     = "/hls/*"
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "S3-${local.media_bucket_name}"

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600 # 1 hour for HLS segments
    max_ttl                = 86400

    # Require signed URLs
    trusted_key_groups = [aws_cloudfront_key_group.signing.id]

    # CORS headers
    response_headers_policy_id = aws_cloudfront_response_headers_policy.cors.id
  }

  # Media files behavior - longer cache
  ordered_cache_behavior {
    path_pattern     = "/media/*"
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "S3-${local.media_bucket_name}"

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 604800 # 1 week
    max_ttl                = 31536000

    # Require signed URLs
    trusted_key_groups = [aws_cloudfront_key_group.signing.id]

    # CORS headers
    response_headers_policy_id = aws_cloudfront_response_headers_policy.cors.id
  }

  # Cover art behavior - public, long cache
  ordered_cache_behavior {
    path_pattern     = "/covers/*"
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "S3-${local.media_bucket_name}"

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 604800 # 1 week
    max_ttl                = 31536000

    # Require signed URLs for cover art too (user's private library)
    trusted_key_groups = [aws_cloudfront_key_group.signing.id]
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  tags = {
    Name = "${local.name_prefix}-media-cdn"
  }
}

# CORS Response Headers Policy
resource "aws_cloudfront_response_headers_policy" "cors" {
  name    = "${local.name_prefix}-cors-policy"
  comment = "CORS policy for media streaming"

  cors_config {
    access_control_allow_credentials = false

    access_control_allow_headers {
      items = ["*"]
    }

    access_control_allow_methods {
      items = ["GET", "HEAD", "OPTIONS"]
    }

    access_control_allow_origins {
      items = ["http://localhost:5173", "https://music.example.com"]
    }

    access_control_max_age_sec = 3600

    origin_override = true
  }
}

# CloudFront Public Key for URL Signing
resource "aws_cloudfront_public_key" "signing" {
  name        = "${local.name_prefix}-signing-key"
  comment     = "Public key for signed URL verification"
  encoded_key = tls_private_key.cloudfront_signing.public_key_pem
}

# CloudFront Key Group
resource "aws_cloudfront_key_group" "signing" {
  name    = "${local.name_prefix}-signing-group"
  comment = "Key group for signed URLs"
  items   = [aws_cloudfront_public_key.signing.id]
}

# Generate RSA Key Pair for CloudFront Signing
resource "tls_private_key" "cloudfront_signing" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

# Store Private Key in Secrets Manager
resource "aws_secretsmanager_secret" "cloudfront_signing_key" {
  name                    = "${local.name_prefix}/cloudfront-signing-key"
  description             = "CloudFront URL signing private key"
  recovery_window_in_days = 7
}

resource "aws_secretsmanager_secret_version" "cloudfront_signing_key" {
  secret_id     = aws_secretsmanager_secret.cloudfront_signing_key.id
  secret_string = tls_private_key.cloudfront_signing.private_key_pem
}

# S3 Bucket Policy to allow CloudFront OAC access
resource "aws_s3_bucket_policy" "media_cloudfront" {
  bucket = local.media_bucket_name

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowCloudFrontOAC"
        Effect = "Allow"
        Principal = {
          Service = "cloudfront.amazonaws.com"
        }
        Action   = "s3:GetObject"
        Resource = "${local.media_bucket_arn}/*"
        Condition = {
          StringEquals = {
            "AWS:SourceArn" = aws_cloudfront_distribution.media.arn
          }
        }
      }
    ]
  })
}

# IAM Policy for Lambda to read CloudFront signing key
resource "aws_iam_role_policy" "lambda_cloudfront_signing" {
  name = "${local.name_prefix}-lambda-cloudfront-signing"
  role = split("/", local.lambda_role_arn)[1]

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "secretsmanager:GetSecretValue"
        ]
        Resource = aws_secretsmanager_secret.cloudfront_signing_key.arn
      }
    ]
  })
}

# Outputs
output "cloudfront_distribution_id" {
  description = "CloudFront distribution ID"
  value       = aws_cloudfront_distribution.media.id
}

output "cloudfront_domain_name" {
  description = "CloudFront domain name for media streaming"
  value       = aws_cloudfront_distribution.media.domain_name
}

output "cloudfront_signing_key_arn" {
  description = "Secrets Manager ARN for CloudFront signing key"
  value       = aws_secretsmanager_secret.cloudfront_signing_key.arn
}

output "cloudfront_key_pair_id" {
  description = "CloudFront key pair ID for signed URLs"
  value       = aws_cloudfront_public_key.signing.id
}
