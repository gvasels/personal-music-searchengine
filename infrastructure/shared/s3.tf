# S3 Bucket for Media Assets (audio files and cover art)
# Uses Intelligent-Tiering for automatic cost optimization

data "aws_caller_identity" "current" {}

resource "aws_s3_bucket" "media" {
  bucket = "${local.name_prefix}-media"
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
    allowed_origins = local.s3_cors_origins
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}

# Intelligent-Tiering: automatic tiers handle cost optimization
# Frequent Access → Infrequent Access (30 days) → Archive Instant Access (90 days)
# No opt-in Archive/Deep Archive tiers — Glacier Instant Retrieval is the lowest tier

# Note: Bucket policy for media is managed in infrastructure/backend/cloudfront.tf
# (combines CloudFront OAC + Bedrock access in a single policy)

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

  # Transition all objects to Intelligent-Tiering immediately
  # Covers media/, hls/, coverart/ and any other prefixes
  rule {
    id     = "intelligent-tiering-transition"
    status = "Enabled"

    filter {}

    transition {
      days          = 0
      storage_class = "INTELLIGENT_TIERING"
    }
  }
}

# S3 Bucket for Search Indexes (Nixiesearch)
resource "aws_s3_bucket" "search_indexes" {
  bucket = "${local.name_prefix}-search-indexes"
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

# Transition search index objects to Intelligent-Tiering
resource "aws_s3_bucket_lifecycle_configuration" "search_indexes" {
  bucket = aws_s3_bucket.search_indexes.id

  rule {
    id     = "intelligent-tiering-transition"
    status = "Enabled"

    filter {}

    transition {
      days          = 0
      storage_class = "INTELLIGENT_TIERING"
    }
  }
}

# Search indexes bucket outputs defined in main.tf
