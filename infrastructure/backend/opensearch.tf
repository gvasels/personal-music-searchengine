# OpenSearch Serverless for Vector Storage
# k-NN enabled collection for Marengo video embeddings

# OpenSearch Serverless Collection
resource "aws_opensearchserverless_collection" "vectors" {
  name        = "${local.name_prefix}-vectors"
  type        = "VECTORSEARCH"
  description = "Vector storage for Marengo video embeddings and semantic search"

  depends_on = [
    aws_opensearchserverless_security_policy.encryption,
    aws_opensearchserverless_security_policy.network
  ]
}

# Encryption Policy (required for collection creation)
resource "aws_opensearchserverless_security_policy" "encryption" {
  name = "${local.name_prefix}-vectors-encryption"
  type = "encryption"

  policy = jsonencode({
    Rules = [
      {
        Resource     = ["collection/${local.name_prefix}-vectors"]
        ResourceType = "collection"
      }
    ]
    AWSOwnedKey = true
  })
}

# Network Policy (public access for Lambda)
resource "aws_opensearchserverless_security_policy" "network" {
  name = "${local.name_prefix}-vectors-network"
  type = "network"

  policy = jsonencode([
    {
      Description = "Public access for Lambda functions"
      Rules = [
        {
          Resource     = ["collection/${local.name_prefix}-vectors"]
          ResourceType = "collection"
        },
        {
          Resource     = ["collection/${local.name_prefix}-vectors"]
          ResourceType = "dashboard"
        }
      ]
      AllowFromPublic = true
    }
  ])
}

# Data Access Policy for Lambda
resource "aws_opensearchserverless_access_policy" "vectors" {
  name = "${local.name_prefix}-vectors-access"
  type = "data"

  policy = jsonencode([
    {
      Description = "Lambda access to vectors collection"
      Rules = [
        {
          Resource     = ["collection/${local.name_prefix}-vectors"]
          ResourceType = "collection"
          Permission = [
            "aoss:CreateCollectionItems",
            "aoss:DeleteCollectionItems",
            "aoss:UpdateCollectionItems",
            "aoss:DescribeCollectionItems"
          ]
        },
        {
          Resource     = ["index/${local.name_prefix}-vectors/*"]
          ResourceType = "index"
          Permission = [
            "aoss:CreateIndex",
            "aoss:DeleteIndex",
            "aoss:UpdateIndex",
            "aoss:DescribeIndex",
            "aoss:ReadDocument",
            "aoss:WriteDocument"
          ]
        }
      ]
      Principal = [
        aws_iam_role.bedrock_gateway.arn,
        local.lambda_role_arn
      ]
    }
  ])
}

# Add OpenSearch permissions to Bedrock Gateway Lambda role
resource "aws_iam_role_policy" "bedrock_gateway_opensearch" {
  name = "opensearch-access"
  role = aws_iam_role.bedrock_gateway.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "aoss:APIAccessAll"
        ]
        Resource = aws_opensearchserverless_collection.vectors.arn
      }
    ]
  })
}

# Add OpenSearch permissions to main Lambda role
resource "aws_iam_role_policy" "api_lambda_opensearch" {
  name = "opensearch-access"
  role = split("/", local.lambda_role_arn)[1]

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "aoss:APIAccessAll"
        ]
        Resource = aws_opensearchserverless_collection.vectors.arn
      }
    ]
  })
}

# Index Template (for reference - applied by application code)
# The actual index is created by the application on first use
# Index mapping for video embeddings:
# {
#   "mappings": {
#     "properties": {
#       "videoId": { "type": "keyword" },
#       "userId": { "type": "keyword" },
#       "embedding": {
#         "type": "knn_vector",
#         "dimension": 1024,
#         "method": {
#           "name": "hnsw",
#           "space_type": "cosinesimil",
#           "engine": "nmslib",
#           "parameters": {
#             "ef_construction": 256,
#             "m": 16
#           }
#         }
#       },
#       "title": { "type": "text" },
#       "artist": { "type": "text" },
#       "bpm": { "type": "float" },
#       "key": { "type": "keyword" },
#       "duration": { "type": "integer" },
#       "createdAt": { "type": "date" },
#       "tags": { "type": "keyword" }
#     }
#   },
#   "settings": {
#     "index": {
#       "knn": true,
#       "knn.algo_param.ef_search": 256
#     }
#   }
# }

# Outputs
output "opensearch_collection_endpoint" {
  description = "OpenSearch Serverless collection endpoint"
  value       = aws_opensearchserverless_collection.vectors.collection_endpoint
}

output "opensearch_collection_arn" {
  description = "OpenSearch Serverless collection ARN"
  value       = aws_opensearchserverless_collection.vectors.arn
}

output "opensearch_dashboard_endpoint" {
  description = "OpenSearch Serverless dashboard endpoint"
  value       = aws_opensearchserverless_collection.vectors.dashboard_endpoint
}
