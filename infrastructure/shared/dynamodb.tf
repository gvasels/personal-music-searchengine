# DynamoDB Table - Single Table Design
resource "aws_dynamodb_table" "music_library" {
  name         = "MusicLibrary"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "PK"
  range_key    = "SK"

  # Primary key attributes
  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }

  # GSI1 attributes
  attribute {
    name = "GSI1PK"
    type = "S"
  }

  attribute {
    name = "GSI1SK"
    type = "S"
  }

  # Global Secondary Index 1 - For artist-based queries and tag lookups
  global_secondary_index {
    name            = "GSI1"
    hash_key        = "GSI1PK"
    range_key       = "GSI1SK"
    projection_type = "ALL"
  }

  # Point-in-time recovery
  point_in_time_recovery {
    enabled = true
  }

  # TTL for temporary data (uploads, etc.)
  ttl {
    attribute_name = "ExpiresAt"
    enabled        = true
  }

  # Server-side encryption
  server_side_encryption {
    enabled = true
  }

  # Tags
  tags = {
    Name = "MusicLibrary"
  }

  lifecycle {
    prevent_destroy = true
  }
}

# DynamoDB Table - Entity Reference
# ================================================================
# Entity Types and Key Patterns:
# ================================================================
#
# USER:
#   PK: USER#{userId}
#   SK: PROFILE
#
# TRACK:
#   PK: USER#{userId}
#   SK: TRACK#{trackId}
#   GSI1PK: USER#{userId}#ARTIST#{artist}
#   GSI1SK: TRACK#{trackId}
#
# ALBUM:
#   PK: USER#{userId}
#   SK: ALBUM#{albumId}
#   GSI1PK: USER#{userId}#ARTIST#{artist}
#   GSI1SK: ALBUM#{year}
#
# PLAYLIST:
#   PK: USER#{userId}
#   SK: PLAYLIST#{playlistId}
#
# PLAYLIST_TRACK:
#   PK: PLAYLIST#{playlistId}
#   SK: POSITION#{position}  (zero-padded: POSITION#00000001)
#
# UPLOAD:
#   PK: USER#{userId}
#   SK: UPLOAD#{uploadId}
#   GSI1PK: UPLOAD#STATUS#{status}
#   GSI1SK: {timestamp}
#
# TAG:
#   PK: USER#{userId}
#   SK: TAG#{tagName}
#
# TRACK_TAG:
#   PK: USER#{userId}#TRACK#{trackId}
#   SK: TAG#{tagName}
#   GSI1PK: USER#{userId}#TAG#{tagName}
#   GSI1SK: TRACK#{trackId}
#
# ================================================================
# Access Patterns:
# ================================================================
#
# 1. Get user profile:
#    Query: PK = USER#{userId}, SK = PROFILE
#
# 2. List all tracks for a user:
#    Query: PK = USER#{userId}, SK begins_with TRACK#
#
# 3. List tracks by artist:
#    Query GSI1: GSI1PK = USER#{userId}#ARTIST#{artist}
#
# 4. List all albums for a user:
#    Query: PK = USER#{userId}, SK begins_with ALBUM#
#
# 5. List albums by artist:
#    Query GSI1: GSI1PK = USER#{userId}#ARTIST#{artist}, GSI1SK begins_with ALBUM#
#
# 6. List all playlists:
#    Query: PK = USER#{userId}, SK begins_with PLAYLIST#
#
# 7. Get playlist tracks (ordered):
#    Query: PK = PLAYLIST#{playlistId}, SK begins_with POSITION#
#
# 8. List uploads by status:
#    Query GSI1: GSI1PK = UPLOAD#STATUS#{status}
#
# 9. List all tags for a user:
#    Query: PK = USER#{userId}, SK begins_with TAG#
#
# 10. Find tracks by tag:
#     Query GSI1: GSI1PK = USER#{userId}#TAG#{tagName}
#
# ================================================================
