# Tasks Document: Bedrock Access Gateway + Marengo Integration

## Task Overview
| Task | Description | Estimated Files |
|------|-------------|-----------------|
| 1 | Infrastructure (OpenTofu) | 3 |
| 2 | Gateway Lambda Handler | 2 |
| 3 | Bedrock Client | 1 |
| 4 | Marengo Client | 1 |
| 5 | Vector Store (OpenSearch) | 2 |
| 6 | Video Processor Lambda | 2 |
| 7 | Search API Handlers | 1 |
| 8 | S3 Tables Integration | 2 |
| 9 | Tests | 4 |

---

- [x] 1. Create Infrastructure for Bedrock Gateway
  - Files: infrastructure/backend/bedrock-gateway.tf, infrastructure/backend/opensearch.tf, infrastructure/backend/variables.tf
  - Deploy API Gateway with Lambda integration for gateway
  - Create OpenSearch Serverless collection for vector storage
  - Configure IAM roles for Bedrock access
  - Set up S3 bucket for video storage
  - Purpose: Cloud infrastructure for AI gateway
  - _Leverage: infrastructure/backend/api-gateway.tf for API Gateway patterns_
  - _Requirements: 1.1, 1.6, 4.1, 4.6_
  - _Prompt: Implement the task for spec bedrock-marengo-gateway, first run spec-workflow-guide to get the workflow guide then implement the task: Role: OpenTofu/Infrastructure Engineer | Task: Create OpenTofu modules for API Gateway with Lambda, OpenSearch Serverless with k-NN enabled, IAM roles for Bedrock invoke permissions, S3 bucket for videos | Restrictions: Use OpenTofu (not Terraform), follow existing patterns, enable encryption at rest | Success: Infrastructure deploys successfully, endpoints accessible, IAM permissions correct | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 2. Create Gateway Lambda Handler
  - Files: backend/cmd/gateway/main.go, backend/internal/handlers/gateway.go
  - Implement OpenAI-compatible endpoints: /v1/chat/completions, /v1/embeddings, /v1/models
  - Add /v1/embeddings/video for Marengo video embeddings
  - Handle API key authentication
  - Route requests to appropriate Bedrock models
  - Purpose: OpenAI-compatible API gateway
  - _Leverage: backend/cmd/api/main.go for Lambda handler patterns_
  - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 2.1_
  - _Prompt: Implement the task for spec bedrock-marengo-gateway, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer with AWS Lambda expertise | Task: Create Lambda handler with Echo for OpenAI-compatible API. Implement /v1/chat/completions (streaming SSE), /v1/embeddings (text), /v1/embeddings/video (Marengo), /v1/models. Translate OpenAI request format to Bedrock | Restrictions: Follow OpenAI API spec exactly, support streaming, handle model mapping | Success: OpenAI SDK can connect, chat completions stream, embeddings return vectors | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 3. Create Bedrock Client
  - File: backend/internal/clients/bedrock.go
  - Implement BedrockClient interface for model invocation
  - Support streaming with channels
  - Handle model ID mapping (claude-3-sonnet → anthropic.claude-3-sonnet-xxx)
  - Add retry logic and error handling
  - Purpose: Abstraction for Bedrock API calls
  - _Leverage: AWS SDK v2 patterns_
  - _Requirements: 1.1, 1.2, 1.3_
  - _Prompt: Implement the task for spec bedrock-marengo-gateway, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer with AWS SDK expertise | Task: Create BedrockClient that wraps AWS SDK Bedrock Runtime. Implement InvokeModel and InvokeModelWithResponseStream. Map OpenAI model names to Bedrock model IDs. Handle errors and retries | Restrictions: Use AWS SDK v2, support context cancellation, proper error wrapping | Success: Can invoke Claude, Nova models, streaming works, errors handled gracefully | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [x] 4. Create Marengo Client
  - File: backend/internal/clients/marengo.go
  - Implement MarengoClient interface for video embeddings
  - Call TwelveLabs Marengo via Bedrock
  - Support different embedding types (visual, audio, speech, combined)
  - Handle async job status for long videos
  - Purpose: Abstraction for Marengo video embedding API
  - _Leverage: backend/internal/clients/bedrock.go_
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6_
  - _Prompt: Implement the task for spec bedrock-marengo-gateway, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer with ML/embedding expertise | Task: Create MarengoClient that calls TwelveLabs Marengo via Bedrock. Generate 1024-dim video embeddings. Support visual/audio/speech/combined modes. Handle 4-hour video limit. Poll job status for async processing | Restrictions: Validate video duration before processing, handle timeouts, return structured embedding response | Success: Can embed videos up to 4 hours, different embedding types work, async jobs complete | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 5. Create Vector Store with OpenSearch
  - Files: backend/internal/vectorstore/opensearch.go, backend/internal/vectorstore/types.go
  - Implement VectorStore interface
  - Create index with HNSW k-NN mapping
  - Implement IndexVideo, SearchSimilar, DeleteVideo
  - Support metadata filtering
  - Purpose: Store and query video embeddings
  - _Leverage: OpenSearch Go client_
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6_
  - _Prompt: Implement the task for spec bedrock-marengo-gateway, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer with vector database expertise | Task: Create VectorStore using OpenSearch Serverless. Create index with knn_vector field (1024 dims, HNSW, cosine similarity). Implement CRUD and k-NN search with metadata filters | Restrictions: Use OpenSearch Go client, handle connection pooling, efficient batch operations | Success: Vectors indexed, k-NN search returns similar videos, filters work | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 6. Create Video Processor Lambda
  - Files: backend/cmd/processor/video/main.go, backend/internal/service/video.go
  - Lambda triggered by S3 video upload
  - Extract video metadata
  - Call Marengo for embedding
  - Store embedding in vector store
  - Update DynamoDB with processing status
  - Purpose: Async video processing pipeline
  - _Leverage: backend/cmd/processor/metadata/main.go for processor patterns_
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6_
  - _Prompt: Implement the task for spec bedrock-marengo-gateway, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer with async processing expertise | Task: Create video processor Lambda triggered by S3. Extract metadata (duration, format), call MarengoClient for embedding, store in VectorStore, update DynamoDB status. Handle failures with retry/dead-letter | Restrictions: Handle large videos (up to 4 hours), idempotent processing, proper error states | Success: Videos processed automatically, embeddings stored, status tracked | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 7. Create Search API Handlers
  - File: backend/internal/handlers/video_search.go
  - Implement POST /api/v1/search/videos (text-to-video)
  - Implement POST /api/v1/search/videos/by-audio (audio-to-video)
  - Implement POST /api/v1/search/videos/by-image (image-to-video)
  - Add pagination and filtering
  - Purpose: Video search API endpoints
  - _Leverage: backend/internal/handlers/search.go_
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_
  - _Prompt: Implement the task for spec bedrock-marengo-gateway, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer | Task: Create search handlers for multimodal video search. Text queries embed via Bedrock text embeddings then search VectorStore. Audio/image queries use Marengo to embed then search. Support pagination and metadata filters | Restrictions: Validate user ownership, apply similarity threshold (0.5), handle missing content | Success: All search modalities work, results ranked by similarity, proper pagination | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 8. Integrate S3 Tables for Nixiesearch
  - Files: backend/internal/search/s3tables.go, infrastructure/backend/s3-tables.tf
  - Create S3 Tables (Iceberg) for search index
  - Store embeddings with metadata for audio tracks
  - Embed metadata (title, artist, BPM, key, tags) via Bedrock embeddings
  - Integrate with existing Nixiesearch for hybrid search
  - Purpose: Improved search performance with semantic embeddings
  - _Leverage: backend/internal/search/nixiesearch.go_
  - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6_
  - _Prompt: Implement the task for spec bedrock-marengo-gateway, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer with data lake expertise | Task: Create S3 Tables (Iceberg format) for search index. Store audio track embeddings generated from metadata (title, artist, BPM, key) via Bedrock text embeddings. Integrate with Nixiesearch for hybrid keyword+semantic search | Restrictions: Use Apache Iceberg format, partition by userId, maintain consistency with Nixiesearch | Success: S3 Tables created, embeddings stored, hybrid search works | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 9. Write tests for Gateway and Video processing
  - Files: backend/internal/handlers/gateway_test.go, backend/internal/clients/bedrock_test.go, backend/internal/clients/marengo_test.go, backend/internal/vectorstore/opensearch_test.go
  - Unit tests with mocked AWS services
  - Integration tests with LocalStack
  - End-to-end search tests
  - Purpose: Ensure AI gateway and video processing reliability
  - _Leverage: Existing test patterns_
  - _Requirements: All_
  - _Prompt: Implement the task for spec bedrock-marengo-gateway, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Test Engineer | Task: Write tests for gateway handlers (OpenAI compatibility), Bedrock client (with mocks), Marengo client, VectorStore. Test streaming, error handling, search accuracy | Restrictions: Mock AWS services, use test containers for OpenSearch, test edge cases | Success: 80%+ coverage, all tests pass, mocks properly isolated | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_
