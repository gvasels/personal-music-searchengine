# Requirements Document: Bedrock Access Gateway + Marengo Integration

## Introduction

This spec deploys the AWS Bedrock Access Gateway with extended support for TwelveLabs Marengo video embedding models. The gateway provides OpenAI-compatible APIs for all Bedrock models while adding custom endpoints for Marengo's multimodal video understanding capabilities. This enables text-to-video search, audio-to-video matching, and visual content analysis for music videos, concert footage, and DJ sets.

## Alignment with Product Vision

This directly supports:
- **AI Chatbot & Agent System** - Agents can reference and search video content
- **Creator Studio** - Video creators can analyze and search their content
- **Rights Management** - Visual content fingerprinting for copyright detection
- **Search & Discovery** - Unified multimodal search across audio and video

## Requirements

### Requirement 1: Bedrock Access Gateway Deployment

**User Story:** As a developer, I want an OpenAI-compatible API gateway for Bedrock, so that I can use standard SDKs and patterns to access Claude, Nova, and other foundation models.

#### Acceptance Criteria

1. WHEN the gateway is deployed THEN the system SHALL expose OpenAI-compatible endpoints at `/v1/chat/completions`, `/v1/embeddings`, `/v1/models`
2. WHEN a request is made with an API key THEN the system SHALL authenticate and route to appropriate Bedrock model
3. WHEN streaming is requested THEN the system SHALL return Server-Sent Events (SSE) responses
4. WHEN a model is specified (e.g., `claude-3-sonnet`) THEN the system SHALL map to corresponding Bedrock model ID
5. IF an invalid model is requested THEN the system SHALL return 400 with available models list
6. WHEN deployed THEN the system SHALL support both Lambda (cost-effective) and Fargate (low-latency) modes

### Requirement 2: Marengo Video Embedding API Extension

**User Story:** As a developer, I want to generate video embeddings via the gateway, so that I can build multimodal search across my video content.

#### Acceptance Criteria

1. WHEN `POST /v1/embeddings/video` is called with video S3 URI THEN the system SHALL invoke Marengo via Bedrock and return 1024-dimensional embeddings
2. WHEN a video exceeds 4 hours THEN the system SHALL return 400 with "Video exceeds maximum duration"
3. WHEN embeddings are requested THEN the system SHALL support `embedding_type`: "visual", "audio", "speech", or "combined"
4. WHEN a composed query is requested THEN the system SHALL accept `{text: "...", image_url: "..."}` for combined search
5. IF video processing fails THEN the system SHALL return 500 with Marengo error details
6. WHEN successful THEN the response SHALL include `{embedding: float[], model: "marengo-3.0", usage: {video_seconds: N}}`

### Requirement 3: Video Upload and Processing Pipeline

**User Story:** As a user, I want to upload videos that are automatically processed for embedding, so that I can search my video library using natural language.

#### Acceptance Criteria

1. WHEN a video is uploaded to S3 `videos/{userId}/{videoId}.{ext}` THEN the system SHALL trigger embedding pipeline
2. WHEN processing starts THEN the system SHALL create a processing record with status "processing"
3. WHEN Marengo returns embeddings THEN the system SHALL store vectors in OpenSearch Serverless (or Pinecone)
4. WHEN processing completes THEN the system SHALL update status to "indexed" with metadata
5. IF processing fails THEN the system SHALL update status to "failed" with error details
6. WHEN a video is deleted THEN the system SHALL remove corresponding embeddings from vector store

### Requirement 4: Vector Storage Integration

**User Story:** As a system, I need to store and query high-dimensional video embeddings efficiently, so that similarity search is fast and scalable.

#### Acceptance Criteria

1. WHEN embeddings are stored THEN the system SHALL use OpenSearch Serverless with k-NN plugin (or Pinecone)
2. WHEN storing vectors THEN the system SHALL include metadata: videoId, userId, duration, title, createdAt
3. WHEN querying THEN the system SHALL support k-NN search with configurable `k` (default 10)
4. WHEN filtering THEN the system SHALL support metadata filters (userId, duration range, date range)
5. WHEN a vector query returns THEN the response SHALL include similarity score (0-1) and video metadata
6. WHEN indexes are created THEN the system SHALL use HNSW algorithm with ef_construction=256, m=16

### Requirement 5: Text-to-Video Search

**User Story:** As a user, I want to search my videos using natural language queries, so that I can find specific moments without watching entire videos.

#### Acceptance Criteria

1. WHEN `POST /api/v1/search/videos` is called with text query THEN the system SHALL embed query and perform k-NN search
2. WHEN results are returned THEN the system SHALL include video ID, title, similarity score, and timestamp (if available)
3. WHEN searching THEN the system SHALL respect user's library (only search their videos)
4. IF no results match above threshold (0.5) THEN the system SHALL return empty array with message
5. WHEN search includes filters THEN the system SHALL apply metadata filters before k-NN
6. WHEN pagination is requested THEN the system SHALL support `limit` and `offset` parameters

### Requirement 6: Cross-Modal Search (Audio-to-Video, Image-to-Video)

**User Story:** As a DJ, I want to find videos that match an audio track or contain similar visuals to an image, so that I can sync visuals with my sets.

#### Acceptance Criteria

1. WHEN `POST /api/v1/search/videos/by-audio` is called with audio S3 URI THEN the system SHALL extract audio embedding and search video embeddings
2. WHEN `POST /api/v1/search/videos/by-image` is called with image URL THEN the system SHALL embed image and search video embeddings
3. WHEN audio-to-video search runs THEN the system SHALL use Marengo's audio embedding modality
4. WHEN image-to-video search runs THEN the system SHALL use Marengo's visual embedding modality
5. IF source content cannot be processed THEN the system SHALL return 400 with supported formats
6. WHEN results are returned THEN the system SHALL include matched modality (audio/visual) in response

### Requirement 7: S3 Tables for Search Index (Nixiesearch Enhancement)

**User Story:** As a platform operator, I want to use S3 Tables for Nixiesearch indexing, so that search operations are faster and more cost-effective.

#### Acceptance Criteria

1. WHEN Nixiesearch indexes are updated THEN the system SHALL write to S3 Tables (Apache Iceberg format)
2. WHEN search queries run THEN the system SHALL read from S3 Tables via Athena/Spark integration
3. WHEN indexes are compacted THEN the system SHALL use S3 Table maintenance jobs
4. IF S3 Tables are unavailable THEN the system SHALL fall back to direct Nixiesearch queries
5. WHEN index size grows THEN the system SHALL partition by userId and date
6. WHEN monitoring THEN the system SHALL track query latency, index size, and cost metrics

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: Gateway handler, Marengo client, Vector store client, Search service as separate modules
- **Modular Design**: Bedrock client should be swappable (for testing/mocking)
- **Dependency Management**: Use AWS SDK v2 for all AWS service integrations
- **Clear Interfaces**: Define `EmbeddingService`, `VectorStore`, `SearchService` interfaces

### Performance
- Gateway response time must be <100ms + model latency
- Video embedding must complete within 2x video duration (e.g., 1-hour video = 2-hour max processing)
- Vector search must return top-10 results in <200ms
- Text-to-video search end-to-end must complete in <500ms

### Security
- API keys must be stored in AWS Secrets Manager
- Video access must verify user ownership before embedding
- Gateway must use IAM roles for Bedrock access (no hardcoded credentials)
- All S3 access must use signed URLs with expiration

### Reliability
- Gateway must auto-scale based on request volume
- Failed video processing must retry 3 times with exponential backoff
- Vector store must have daily backups
- S3 Tables must have versioning enabled for rollback

### Usability
- Clear API documentation with OpenAPI spec
- Example requests/responses for all endpoints
- Progress webhook for long-running video processing
- Dashboard for monitoring processing queue
