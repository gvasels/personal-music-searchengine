# S3 Bucket for Media Assets (audio files and cover art)
# Uses Intelligent-Tiering for automatic cost optimization
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
    allowed_origins = ["http://localhost:5173", "https://music.example.com"]
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

# Lifecycle rules
resource "aws_s3_bucket_lifecycle_configuration" "media" {
  bucket = aws_s3_bucket.media.id

  # Rule for incomplete multipart uploads
  rule {
    id     = "abort-incomplete-multipart"
    status = "Enabled"

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

# Output for search indexes bucket
output "search_indexes_bucket_name" {
  value = aws_s3_bucket.search_indexes.id
}

output "search_indexes_bucket_arn" {
  value = aws_s3_bucket.search_indexes.arn
}
