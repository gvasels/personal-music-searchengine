# S3 Bucket for Media Assets (audio files and cover art)
# Uses Intelligent-Tiering for automatic cost optimization

data "aws_caller_identity" "current" {}

resource "aws_s3_bucket" "media" {
  bucket = "${data.aws_caller_identity.current.account_id}-${local.name_prefix}-media"
}

resource "aws_s3_bucket_versioning" "media" {
  bucket = aws_s3_bucket.media.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "media" {
  bucket = aws_s3_bucket.media.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "media" {
  bucket = aws_s3_bucket.media.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# CORS configuration for browser uploads
resource "aws_s3_bucket_cors_configuration" "media" {
  bucket = aws_s3_bucket.media.id

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST", "GET", "HEAD"]
    allowed_origins = [
      "http://localhost:5173",
      "http://localhost:3000",
      "https://d8wn3lkytn5qe.cloudfront.net",
      "https://d1xxw2bv6ilv0c.cloudfront.net"
    ]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}

# Intelligent-Tiering configuration for automatic cost optimization
resource "aws_s3_bucket_intelligent_tiering_configuration" "media" {
  bucket = aws_s3_bucket.media.id
  name   = "EntireBucket"

  tiering {
    access_tier = "ARCHIVE_ACCESS"
    days        = 90
  }

  tiering {
    access_tier = "DEEP_ARCHIVE_ACCESS"
    days        = 180
  }
}

# Bucket policy to allow Bedrock async invoke to read/write
resource "aws_s3_bucket_policy" "media_bedrock" {
  bucket = aws_s3_bucket.media.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "BedrockAsyncInvokeAccess"
        Effect = "Allow"
        Principal = {
          Service = "bedrock.amazonaws.com"
        }
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:GetBucketLocation"
        ]
        Resource = [
          aws_s3_bucket.media.arn,
          "${aws_s3_bucket.media.arn}/*"
        ]
      }
    ]
  })
}

# Lifecycle rules
resource "aws_s3_bucket_lifecycle_configuration" "media" {
  bucket = aws_s3_bucket.media.id

  # Rule for incomplete multipart uploads
  rule {
    id     = "abort-incomplete-multipart"
    status = "Enabled"

    filter {}

    abort_incomplete_multipart_upload {
      days_after_initiation = 7
    }
  }

  # Rule for temporary uploads folder - delete unprocessed after 7 days
  rule {
    id     = "cleanup-temp-uploads"
    status = "Enabled"

    filter {
      prefix = "uploads/"
    }

    expiration {
      days = 7
    }
  }

  # Transition all objects to Intelligent-Tiering after upload
  rule {
    id     = "intelligent-tiering-transition"
    status = "Enabled"

    filter {
      prefix = "media/"
    }

    transition {
      days          = 0
      storage_class = "INTELLIGENT_TIERING"
    }
  }
}

# S3 Bucket for Search Indexes (Nixiesearch)
resource "aws_s3_bucket" "search_indexes" {
  bucket = "${data.aws_caller_identity.current.account_id}-${local.name_prefix}-search-indexes"
}

resource "aws_s3_bucket_versioning" "search_indexes" {
  bucket = aws_s3_bucket.search_indexes.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "search_indexes" {
  bucket = aws_s3_bucket.search_indexes.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

resource "aws_s3_bucket_public_access_block" "search_indexes" {
  bucket = aws_s3_bucket.search_indexes.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Search indexes bucket outputs defined in main.tf
