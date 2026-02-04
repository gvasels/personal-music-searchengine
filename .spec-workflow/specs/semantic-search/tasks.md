# Tasks - Semantic Search Enhancements

## Epic: Semantic Search with Bedrock Embeddings
**Status**: In Progress
**Stream**: C - Search Enhancements

---

## Group 1: Embedding Service

### Task 1.1: Create Embedding Service
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/embedding.go`
- `backend/internal/service/embedding_test.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `NewEmbeddingService(bedrockClient, modelID)` | Creates embedding service |
| `GenerateTrackEmbedding(ctx, track)` | Generate embedding for single track |
| `GenerateQueryEmbedding(ctx, query)` | Generate embedding for search query |
| `BatchGenerateEmbeddings(ctx, tracks)` | Generate embeddings for multiple tracks |
| `ComposeEmbedText(track)` | Create text to embed from metadata |

**Tests**:
| Test | Description |
|------|-------------|
| `TestComposeEmbedText_AllFields` | Full metadata produces correct text |
| `TestComposeEmbedText_MinimalFields` | Title-only track works |
| `TestComposeEmbedText_MaxLength` | Text truncated at 8000 chars |
| `TestGenerateTrackEmbedding_Success` | Returns 1024-dim embedding |
| `TestGenerateTrackEmbedding_BedrockError` | Handles Bedrock failure |
| `TestBatchGenerateEmbeddings_Success` | Batch of 25 tracks works |

**Acceptance Criteria**:
- [ ] Embedding service created with BedrockClient dependency
- [ ] Text composition includes title, artist, album, genre, tags, BPM, key
- [ ] Titan v2 model used for embeddings
- [ ] Batch processing respects rate limits
- [ ] Errors handled gracefully

### Task 1.2: Enhance Search Document with Embedding
**Status**: [ ] Pending
**Files**:
- `backend/internal/search/types.go`

**Functions**:
| Function | Description |
|----------|-------------|
| Add `Embedding []float32` to Document | Store vector in index |
| Add `BPM int` to Document | BPM for filtering |
| Add `KeyCamelot string` to Document | Camelot key for filtering |
| Add `Tags []string` to Document | Tags for filtering |

**Acceptance Criteria**:
- [ ] Document struct includes embedding field
- [ ] DJ metadata fields added (BPM, KeyCamelot)
- [ ] Tags added as array field

---

## Group 2: Semantic Search Integration

### Task 2.1: Add Semantic Search to Search Client
**Status**: [ ] Pending
**Files**:
- `backend/internal/search/client.go`
- `backend/internal/search/client_test.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `SemanticSearch(ctx, userID, query)` | Vector similarity search |
| `HybridSearch(ctx, userID, query)` | Combined keyword + vector search |

**Tests**:
| Test | Description |
|------|-------------|
| `TestSemanticSearch_Success` | Vector search returns ranked results |
| `TestHybridSearch_WeightedScoring` | Combined scores correct |
| `TestSemanticSearch_Filters` | BPM/key filters work |

**Acceptance Criteria**:
- [ ] Semantic search via vector similarity
- [ ] Hybrid mode combines BM25 + cosine scores
- [ ] Configurable weight parameters
- [ ] Works with existing filter system

### Task 2.2: Update Search Service for Semantic Search
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/search.go`
- `backend/internal/service/search_test.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `SemanticSearch(ctx, userID, req)` | Semantic search with embedding generation |
| `IndexTrackWithEmbedding(ctx, track)` | Index track with embedding |

**Tests**:
| Test | Description |
|------|-------------|
| `TestSemanticSearch_GeneratesQueryEmbedding` | Query embedded before search |
| `TestIndexTrackWithEmbedding_Success` | Track indexed with embedding |
| `TestSemanticSearch_FallbackToKeyword` | Keyword search if embedding fails |

**Acceptance Criteria**:
- [ ] Query embedding generated via EmbeddingService
- [ ] Track indexing includes embedding generation
- [ ] Fallback to keyword search if embedding fails
- [ ] Existing search functionality preserved

### Task 2.3: Add Semantic Search Handler
**Status**: [ ] Pending
**Files**:
- `backend/internal/handlers/search.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `(h *Handlers) SemanticSearch(c)` | POST /api/v1/search/semantic |
| Add `semantic` query param to `SearchSimple` | Enable hybrid mode |

**Tests**:
| Test | Description |
|------|-------------|
| `TestSemanticSearch_Success` | Handler returns semantic results |
| `TestSearchSimple_WithSemanticParam` | Hybrid mode via query param |

**Acceptance Criteria**:
- [ ] New POST /api/v1/search/semantic endpoint
- [ ] Existing GET /search supports `semantic=true` param
- [ ] Response includes semantic scores

---

## Group 3: Similar Tracks Feature

### Task 3.1: Create Camelot Key Utilities
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/camelot.go`
- `backend/internal/service/camelot_test.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `IsKeyCompatible(key1, key2)` | Check if keys mix harmonically |
| `GetCompatibleKeys(key)` | Return all compatible keys |
| `GetKeyTransition(from, to)` | Describe mixing transition |

**Tests**:
| Test | Description |
|------|-------------|
| `TestIsKeyCompatible_SameKey` | Same key is compatible |
| `TestIsKeyCompatible_Adjacent` | Adjacent keys compatible |
| `TestIsKeyCompatible_RelativeMajorMinor` | A/B switch compatible |
| `TestIsKeyCompatible_NotCompatible` | Opposite keys not compatible |
| `TestGetCompatibleKeys_AllKeys` | All 24 keys have 4 compatibles |

**Acceptance Criteria**:
- [ ] Camelot wheel implemented for all 24 keys
- [ ] Key compatibility check function
- [ ] Transition description for DJ UI

### Task 3.2: Create Similarity Service
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/similarity.go`
- `backend/internal/service/similarity_test.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `NewSimilarityService(searchClient, repo, embeddingSvc)` | Create service |
| `FindSimilarTracks(ctx, userID, trackID, opts)` | Find similar tracks |
| `FindMixableTracks(ctx, userID, trackID, opts)` | Find DJ-compatible tracks |
| `ComputeSimilarity(track1, track2, emb1, emb2)` | Calculate similarity score |
| `cosineSimilarity(vec1, vec2)` | Vector similarity calculation |

**Tests**:
| Test | Description |
|------|-------------|
| `TestFindSimilarTracks_ByEmbedding` | Embedding similarity works |
| `TestFindSimilarTracks_Combined` | Embedding + features combined |
| `TestFindMixableTracks_BPMRange` | BPM tolerance filter works |
| `TestFindMixableTracks_HarmonicKeys` | Key compatibility filter works |
| `TestComputeSimilarity_HighMatch` | Similar tracks get high score |
| `TestCosineSimilarity_Normalized` | Returns 0-1 range |

**Acceptance Criteria**:
- [ ] Similar tracks by embedding similarity
- [ ] Mixable tracks by BPM + key compatibility
- [ ] Combined similarity scoring
- [ ] Proper sorting by similarity score

### Task 3.3: Add Similar Tracks Handler
**Status**: [ ] Pending
**Files**:
- `backend/internal/handlers/tracks.go`

**Functions**:
| Function | Description |
|----------|-------------|
| `(h *Handlers) GetSimilarTracks(c)` | GET /api/v1/tracks/:id/similar |
| `(h *Handlers) GetMixableTracks(c)` | GET /api/v1/tracks/:id/mixable |

**Tests**:
| Test | Description |
|------|-------------|
| `TestGetSimilarTracks_Success` | Returns similar tracks |
| `TestGetSimilarTracks_NotFound` | 404 for unknown track |
| `TestGetMixableTracks_Success` | Returns DJ-compatible tracks |

**Acceptance Criteria**:
- [ ] GET /api/v1/tracks/:id/similar endpoint
- [ ] GET /api/v1/tracks/:id/mixable endpoint
- [ ] Proper error handling
- [ ] Response includes match reasons

---

## Group 4: Index Enhancement

### Task 4.1: Update Indexer Lambda for Embeddings
**Status**: [ ] Pending
**Files**:
- `backend/cmd/processor/indexer/main.go`

**Functions**:
| Function | Description |
|----------|-------------|
| Update `indexTrack` to generate embedding | Include embedding in index |
| Add `buildSearchDocumentWithEmbedding` | Create full document |

**Tests**:
| Test | Description |
|------|-------------|
| `TestIndexTrack_WithEmbedding` | Embedding included in index |
| `TestIndexTrack_EmbeddingFailure` | Graceful handling of failure |

**Acceptance Criteria**:
- [ ] Indexer generates embedding during track indexing
- [ ] Embedding failure doesn't block indexing
- [ ] DJ metadata (BPM, key) included in index

### Task 4.2: Index Rebuild with Embeddings
**Status**: [ ] Pending
**Files**:
- `backend/internal/service/search.go`

**Functions**:
| Function | Description |
|----------|-------------|
| Update `RebuildIndex` to include embeddings | Batch embedding generation |

**Tests**:
| Test | Description |
|------|-------------|
| `TestRebuildIndex_WithEmbeddings` | All tracks get embeddings |
| `TestRebuildIndex_PartialFailure` | Continues on individual failures |

**Acceptance Criteria**:
- [ ] Rebuild generates embeddings in batches
- [ ] Rate limiting respected
- [ ] Progress tracking for large libraries

---

## Summary

| Group | Tasks | Status |
|-------|-------|--------|
| Group 1: Embedding Service | 2 | Not Started |
| Group 2: Semantic Search Integration | 3 | Not Started |
| Group 3: Similar Tracks Feature | 3 | Not Started |
| Group 4: Index Enhancement | 2 | Not Started |
| **Total** | **10** | **0 Complete** |

---

## Test Plan Summary

### Unit Tests
| File | Tests |
|------|-------|
| `internal/service/embedding_test.go` | Embedding generation |
| `internal/service/camelot_test.go` | Key compatibility |
| `internal/service/similarity_test.go` | Similarity scoring |
| `internal/search/client_test.go` | Vector search |

### Integration Tests
| Test | Environment |
|------|-------------|
| Embedding generation | Mock Bedrock |
| Semantic search flow | LocalStack + Nixiesearch mock |
| Similar tracks | Full service integration |

---

## PR Checklist

After completing all tasks:
- [ ] All tests pass (`go test ./...`)
- [ ] Code builds (`go build ./...`)
- [ ] Linter passes (`golangci-lint run`)
- [ ] CLAUDE.md files updated for new packages
- [ ] CHANGELOG.md updated
- [ ] Semantic search works end-to-end
- [ ] Similar tracks feature works
- [ ] Create PR to main
