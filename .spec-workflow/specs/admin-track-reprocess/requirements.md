# Requirements Document: Admin Track Reprocess Button

## Introduction

This feature adds a "Re-process AI" button for administrators to manually trigger full AI reprocessing on individual tracks. When clicked, the button triggers the audio analysis pipeline (BPM detection, musical key detection) and regenerates vector embeddings for semantic search. This is useful when AI analysis failed during initial upload, when analysis algorithms have been improved, or when embeddings need to be regenerated after embedding model updates.

## Requirements

### Requirement 1: Admin-Only Reprocess Button in Track List

**User Story:** As an admin, I want to see a "Re-process AI" button on tracks in the admin track view, so that I can trigger AI reprocessing without re-uploading the file.

#### Acceptance Criteria

1. WHEN an admin views the track list with "Show Uploaded By" enabled THEN the system SHALL display a "Re-process AI" button in the actions column for each track
2. IF a user is NOT an admin THEN the system SHALL NOT display the "Re-process AI" button
3. WHEN the button is clicked THEN the system SHALL display a loading indicator on the button
4. WHEN the reprocess completes successfully THEN the system SHALL show a success toast notification
5. IF the reprocess fails THEN the system SHALL show an error toast notification with the error message

### Requirement 2: Backend Reprocess Track Endpoint

**User Story:** As an admin, I want the reprocess action to trigger the full AI analysis pipeline on a track, so that BPM, musical key, and embeddings are regenerated.

#### Acceptance Criteria

1. WHEN a POST request is made to `/api/v1/tracks/:id/reprocess` with admin credentials THEN the system SHALL start the AI reprocessing pipeline
2. IF the requesting user is NOT an admin THEN the system SHALL return 403 Forbidden
3. IF the track does not exist THEN the system SHALL return 404 Not Found
4. WHEN the track is found THEN the system SHALL queue the track for audio analysis (BPM, key detection)
5. WHEN audio analysis completes THEN the system SHALL regenerate vector embeddings for the track
6. WHEN reprocessing starts THEN the system SHALL return 202 Accepted with a status object
7. WHEN reprocessing completes THEN the track's BPM, MusicalKey, KeyCamelot, and embedding fields SHALL be updated in DynamoDB

### Requirement 3: Reprocess Status Tracking

**User Story:** As an admin, I want to see the reprocess status while it's running, so that I know when the operation completes.

#### Acceptance Criteria

1. WHEN reprocessing starts THEN the system SHALL set `ReprocessStatus` to "processing" on the track
2. WHEN reprocessing completes successfully THEN the system SHALL set `ReprocessStatus` to "complete" and `ReprocessedAt` timestamp
3. IF reprocessing fails THEN the system SHALL set `ReprocessStatus` to "failed" with error details
4. WHEN viewing a track that is being reprocessed THEN the UI SHALL show a "Reprocessing..." indicator
5. WHEN reprocessing completes THEN the UI SHALL refresh the track data to show updated values

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: Handler, service, and model changes should be isolated
- **Modular Design**: Reprocess logic should be a standalone service method callable by other features
- **Clear Interfaces**: Follow existing `TrackService` pattern with `hasGlobal` for admin access

### Performance
- Reprocessing should not block the API response (202 Accepted, async processing)
- UI should remain responsive during reprocessing

### Security
- Only admins can trigger reprocessing (real-time DB role check, not just JWT)
- Rate limiting: Max 10 reprocess requests per minute per admin (prevent abuse)

### Reliability
- Failed reprocessing should not corrupt existing track data
- Partial failures (e.g., audio analysis succeeds but embedding fails) should be handled gracefully

### Usability
- Button should be clearly visible but not disruptive to normal track list usage
- Loading state should prevent accidental double-clicks
- Success/error feedback should be immediate and clear
