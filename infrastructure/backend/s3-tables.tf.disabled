# S3 Tables (Apache Iceberg) for Nixiesearch Indexing
# Provides columnar storage with efficient scan patterns

# S3 Table Bucket for Iceberg tables
resource "aws_s3_table_bucket" "search_index" {
  name = "${local.name_prefix}-search-index"

  maintenance_configuration {
    iceberg_unreferenced_file_removal {
      settings {
        unreferenced_days = 7
        non_current_days  = 3
      }
      status = "enabled"
    }
  }
}

# Table Bucket Policy
resource "aws_s3_table_bucket_policy" "search_index" {
  resource_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowLambdaAccess"
        Effect = "Allow"
        Principal = {
          AWS = [
            local.lambda_role_arn,
            aws_iam_role.bedrock_gateway.arn
          ]
        }
        Action = [
          "s3tables:*"
        ]
        Resource = [
          aws_s3_table_bucket.search_index.arn,
          "${aws_s3_table_bucket.search_index.arn}/*"
        ]
      }
    ]
  })
  table_bucket_arn = aws_s3_table_bucket.search_index.arn
}

# Namespace for music tracks
resource "aws_s3_table" "tracks" {
  name             = "tracks"
  namespace        = [local.name_prefix]
  table_bucket_arn = aws_s3_table_bucket.search_index.arn
  format           = "ICEBERG"

  maintenance_configuration {
    iceberg_compaction {
      settings {
        target_file_size_mb = 128
      }
      status = "enabled"
    }
    iceberg_snapshot_management {
      settings {
        min_snapshots_to_keep = 5
        max_snapshot_age_hours = 168 # 7 days
      }
      status = "enabled"
    }
  }
}

# Namespace for embeddings
resource "aws_s3_table" "embeddings" {
  name             = "embeddings"
  namespace        = [local.name_prefix]
  table_bucket_arn = aws_s3_table_bucket.search_index.arn
  format           = "ICEBERG"

  maintenance_configuration {
    iceberg_compaction {
      settings {
        target_file_size_mb = 256
      }
      status = "enabled"
    }
    iceberg_snapshot_management {
      settings {
        min_snapshots_to_keep = 5
        max_snapshot_age_hours = 168
      }
      status = "enabled"
    }
  }
}

# Add S3 Tables permissions to Lambda role
resource "aws_iam_role_policy" "lambda_s3tables" {
  name = "s3tables-access"
  role = split("/", local.lambda_role_arn)[1]

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3tables:GetTable",
          "s3tables:GetTableBucket",
          "s3tables:GetTableBucketPolicy",
          "s3tables:GetTableMetadataLocation",
          "s3tables:ListTableBuckets",
          "s3tables:ListTables",
          "s3tables:ListNamespaces",
          "s3tables:CreateTable",
          "s3tables:CreateNamespace",
          "s3tables:DeleteTable",
          "s3tables:DeleteNamespace",
          "s3tables:UpdateTableMetadataLocation",
          "s3tables:PutTableBucketPolicy",
          "s3tables:DeleteTableBucketPolicy"
        ]
        Resource = [
          aws_s3_table_bucket.search_index.arn,
          "${aws_s3_table_bucket.search_index.arn}/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ]
        Resource = [
          "arn:aws:s3:::${aws_s3_table_bucket.search_index.name}",
          "arn:aws:s3:::${aws_s3_table_bucket.search_index.name}/*"
        ]
      }
    ]
  })
}

# Add S3 Tables permissions to Bedrock Gateway Lambda role
resource "aws_iam_role_policy" "bedrock_gateway_s3tables" {
  name = "s3tables-access"
  role = aws_iam_role.bedrock_gateway.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "s3tables:GetTable",
          "s3tables:GetTableBucket",
          "s3tables:GetTableMetadataLocation",
          "s3tables:ListTableBuckets",
          "s3tables:ListTables",
          "s3tables:ListNamespaces",
          "s3tables:CreateTable",
          "s3tables:CreateNamespace",
          "s3tables:UpdateTableMetadataLocation"
        ]
        Resource = [
          aws_s3_table_bucket.search_index.arn,
          "${aws_s3_table_bucket.search_index.arn}/*"
        ]
      },
      {
        Effect = "Allow"
        Action = [
          "s3:GetObject",
          "s3:PutObject",
          "s3:DeleteObject",
          "s3:ListBucket"
        ]
        Resource = [
          "arn:aws:s3:::${aws_s3_table_bucket.search_index.name}",
          "arn:aws:s3:::${aws_s3_table_bucket.search_index.name}/*"
        ]
      }
    ]
  })
}

# Outputs
output "s3_table_bucket_arn" {
  description = "S3 Table Bucket ARN for search index"
  value       = aws_s3_table_bucket.search_index.arn
}

output "s3_table_bucket_name" {
  description = "S3 Table Bucket name"
  value       = aws_s3_table_bucket.search_index.name
}

output "tracks_table_arn" {
  description = "Tracks Iceberg table ARN"
  value       = aws_s3_table.tracks.arn
}

output "embeddings_table_arn" {
  description = "Embeddings Iceberg table ARN"
  value       = aws_s3_table.embeddings.arn
}
