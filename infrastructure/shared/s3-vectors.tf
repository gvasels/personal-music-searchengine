# S3 Vectors for similarity search
# Note: Using CloudFormation as Terraform doesn't support S3 Vectors natively yet

resource "aws_cloudformation_stack" "vectors" {
  name = "${local.name_prefix}-vectors"

  template_body = jsonencode({
    AWSTemplateFormatVersion = "2010-09-09"
    Description              = "S3 Vectors bucket and index for media embeddings"
    Resources = {
      VectorBucket = {
        Type = "AWS::S3Vectors::VectorBucket"
        Properties = {
          VectorBucketName = "${data.aws_caller_identity.current.account_id}-${local.name_prefix}-vectors"
        }
      }
      MediaEmbeddingsIndex = {
        Type      = "AWS::S3Vectors::Index"
        DependsOn = ["VectorBucket"]
        Properties = {
          VectorBucketName = "${data.aws_caller_identity.current.account_id}-${local.name_prefix}-vectors"
          IndexName        = "media-embeddings"
          DataType         = "float32"
          Dimension        = 512
          DistanceMetric   = "cosine"
        }
      }
      SearchEmbeddingsIndex = {
        Type      = "AWS::S3Vectors::Index"
        DependsOn = ["VectorBucket"]
        Properties = {
          VectorBucketName = "${data.aws_caller_identity.current.account_id}-${local.name_prefix}-vectors"
          IndexName        = "search-embeddings"
          DataType         = "float32"
          Dimension        = 1024
          DistanceMetric   = "cosine"
        }
      }
    }
    Outputs = {
      VectorBucketName = {
        Value       = "${data.aws_caller_identity.current.account_id}-${local.name_prefix}-vectors"
        Description = "S3 Vectors bucket name"
      }
      VectorBucketArn = {
        Value       = { "Fn::GetAtt" = ["VectorBucket", "VectorBucketArn"] }
        Description = "S3 Vectors bucket ARN"
      }
      IndexName = {
        Value       = "media-embeddings"
        Description = "Vector index name"
      }
      SearchIndexName = {
        Value       = "search-embeddings"
        Description = "Vector index name for search embeddings"
      }
    }
  })

  tags = {
    Name = "${local.name_prefix}-vectors"
  }
}

output "vector_bucket_name" {
  description = "S3 Vectors bucket name"
  value       = aws_cloudformation_stack.vectors.outputs["VectorBucketName"]
}

output "vector_bucket_arn" {
  description = "S3 Vectors bucket ARN"
  value       = aws_cloudformation_stack.vectors.outputs["VectorBucketArn"]
}

output "vector_index_name" {
  description = "Vector index name for media embeddings"
  value       = "media-embeddings"
}

output "search_vector_index_name" {
  description = "Vector index name for search text embeddings"
  value       = "search-embeddings"
}
