# HLS-Based Audio Analysis Pipeline — Requirements

## Overview

Move the audio analysis pipeline trigger from the track-creator Lambda (which runs before HLS transcoding) to the transcode-complete Lambda (which runs after HLS files are ready). Update the audio-features Lambda to sample random HLS segments for BPM/key detection instead of the original uploaded file. Provide a data wipe script for fresh start.

## User Stories

### US-1: Unbiased Audio Analysis via HLS Sampling
**As a** music library owner,
**I want** BPM and key detection to analyze the actual HLS-transcoded audio segments,
**So that** the analysis validates what users hear when streaming and avoids bias from the original file format.

**Acceptance Criteria:**
- [ ] audio-features Lambda lists HLS segments at `hls/{userId}/{trackId}/`
- [ ] Picks a random quality level (96k, 192k, or 320k)
- [ ] Samples 5-8 random segments from that level
- [ ] Falls back to original file if no HLS segments exist
- [ ] Result includes `source` field ("hls" or "original")

### US-2: Audio Pipeline Triggers After HLS Completion
**As a** system operator,
**I want** the audio analysis pipeline to start only after HLS transcoding completes,
**So that** HLS segments are guaranteed to exist when audio-features runs.

**Acceptance Criteria:**
- [ ] transcode-complete Lambda triggers audio pipeline on COMPLETE status
- [ ] transcode-complete does NOT trigger audio pipeline on ERROR/CANCELED
- [ ] track-creator Lambda does NOT trigger audio pipeline for new tracks
- [ ] track-creator Lambda still triggers for duplicates with existing HLS (delta analysis)

### US-3: Fresh Start Data Wipe
**As a** music library owner,
**I want** a script to wipe all tracks, uploads, and media files,
**So that** I can start fresh and re-upload with the new analysis pipeline.

**Acceptance Criteria:**
- [ ] Script deletes all TRACK# and UPLOAD# items from DynamoDB
- [ ] Script deletes all S3 objects under media/, hls/, uploads/, coverart/
- [ ] Script has --dry-run and --confirm safety flags
- [ ] Script reports counts before and after deletion

## Non-Functional Requirements

- No increase in transcode-complete Lambda timeout (30s is sufficient for GetItem + StartExecution)
- Audio-features container image size increase minimal (ffmpeg-free only)
- Graceful degradation: if HLS segments missing, fall back to original file
- Best-effort audio pipeline trigger: failure to read track or start pipeline does not fail the transcode-complete response
