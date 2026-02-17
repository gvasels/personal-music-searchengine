# backend/cmd/processor/audio-features/

## Purpose
Python Lambda function that extracts technical audio features (BPM, key) using librosa.

## Files
- `handler.py` — Lambda handler; downloads audio from S3, detects BPM via beat tracking, detects key via chroma features with Krumhansl-Schmuckler profiles, returns Camelot notation
- `Dockerfile` — Container image build
- `requirements.txt` — Python dependencies (librosa, numpy, etc.)

## Key Details
- Processes 60-second audio samples for speed
- Uses /tmp for Numba cache
- Part of the Step Functions audio analysis pipeline
