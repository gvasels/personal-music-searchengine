# Requirements Document: Audio Analysis & Player Enhancements

## Introduction

This spec enhances the audio player with professional-grade features: waveform visualization, beat grid display, gapless playback, crossfade, and equalizer controls. These features support the DJ Studio and general music listening experience, building on the existing BPM/key detection capabilities.

## Alignment with Product Vision

This directly supports:
- **Creator Studio - DJ Module** - Beat grids, waveforms essential for DJing
- **Player & Playback** - Gapless, crossfade, EQ for premium listening
- **Audio Analysis** - Visual feedback for analyzed audio properties
- **Search & Discovery** - Similar tracks can use audio fingerprints

## Requirements

### Requirement 1: Waveform Generation and Display

**User Story:** As a user, I want to see a visual waveform of the track, so that I can identify song sections (intro, drop, outro) at a glance.

#### Acceptance Criteria

1. WHEN a track is processed THEN the system SHALL generate waveform data (peaks) at 100 samples/second
2. WHEN waveform data is generated THEN the system SHALL store as JSON in S3: `waveforms/{trackId}.json`
3. WHEN the player displays a track THEN the system SHALL render waveform as SVG/Canvas visualization
4. WHEN the user hovers over waveform THEN the system SHALL show timestamp at cursor position
5. WHEN the user clicks waveform THEN the system SHALL seek to that position
6. WHEN playback progresses THEN the system SHALL highlight played portion in accent color
7. IF waveform data is unavailable THEN the system SHALL display a progress bar fallback

### Requirement 2: Beat Grid Visualization

**User Story:** As a DJ, I want to see beat markers on the waveform, so that I can align tracks for mixing and identify downbeats.

#### Acceptance Criteria

1. WHEN BPM is detected THEN the system SHALL calculate beat positions throughout the track
2. WHEN beat grid is generated THEN the system SHALL store as array of timestamps in track metadata
3. WHEN waveform is displayed THEN the system SHALL overlay vertical beat markers at each beat
4. WHEN displaying beats THEN the system SHALL emphasize downbeats (every 4th beat) with different styling
5. WHEN BPM varies (live recordings) THEN the system SHALL show "Variable BPM" and skip beat grid
6. WHEN zooming waveform THEN the system SHALL show/hide beat markers based on zoom level

### Requirement 3: Gapless Playback

**User Story:** As a listener, I want seamless transitions between tracks in an album/playlist, so that continuous mixes and albums play without interruption.

#### Acceptance Criteria

1. WHEN a track nears end (5 seconds remaining) THEN the system SHALL preload the next track
2. WHEN current track ends THEN the system SHALL seamlessly start next track with <10ms gap
3. WHEN using Howler.js THEN the system SHALL use HTML5 audio with preload for gapless
4. IF next track fails to preload THEN the system SHALL fall back to standard transition with brief pause
5. WHEN gapless is enabled THEN the system SHALL show indicator in player UI
6. WHEN user manually skips THEN the system SHALL NOT apply gapless (immediate skip)

### Requirement 4: Crossfade Between Tracks

**User Story:** As a listener, I want configurable crossfade between tracks, so that transitions are smooth for party/background listening.

#### Acceptance Criteria

1. WHEN crossfade is enabled THEN the system SHALL fade out current track while fading in next track
2. WHEN configuring crossfade THEN the user SHALL set duration: 0 (off), 1s, 3s, 5s, 10s
3. WHEN crossfade starts THEN the system SHALL begin `crossfade_duration` seconds before track end
4. WHEN crossfade is active THEN the system SHALL apply equal-power crossfade curve
5. WHEN user skips during crossfade THEN the system SHALL immediately complete transition
6. WHEN shuffle is on THEN crossfade SHALL work with randomly selected next track
7. IF repeat-one is enabled THEN crossfade SHALL NOT apply (instant loop)

### Requirement 5: Equalizer Controls

**User Story:** As a listener, I want to adjust bass, mid, and treble, so that I can customize the sound to my headphones/speakers.

#### Acceptance Criteria

1. WHEN EQ is displayed THEN the system SHALL show 3-band EQ: Low (60Hz), Mid (1kHz), High (12kHz)
2. WHEN user adjusts a band THEN the system SHALL apply gain from -12dB to +12dB in real-time
3. WHEN EQ settings change THEN the system SHALL persist to localStorage per user
4. WHEN EQ is implemented THEN the system SHALL use Web Audio API BiquadFilterNode
5. IF Web Audio API unavailable THEN the system SHALL hide EQ and show "EQ not supported"
6. WHEN displaying EQ THEN the system SHALL include presets: Flat, Bass Boost, Vocal, Electronic

### Requirement 6: Playback Speed Control

**User Story:** As a user, I want to adjust playback speed, so that I can slow down for learning or speed up for previewing.

#### Acceptance Criteria

1. WHEN speed control is available THEN the system SHALL offer: 0.5x, 0.75x, 1x, 1.25x, 1.5x, 2x
2. WHEN speed changes THEN the system SHALL adjust playback rate in real-time without restart
3. WHEN speed is not 1x THEN the system SHALL show current speed indicator in player
4. WHEN using Howler.js THEN the system SHALL use `rate()` method for speed adjustment
5. IF pitch correction is needed THEN the system SHALL preserve pitch at altered speeds (optional toggle)
6. WHEN speed is changed THEN the system SHALL NOT affect duration display (show actual duration)

### Requirement 7: Audio Analysis Backend Enhancements

**User Story:** As a platform, I want to generate waveform and beat grid data during track processing, so that the player has rich visualization data.

#### Acceptance Criteria

1. WHEN a track is uploaded THEN the processor SHALL generate waveform peaks during transcoding
2. WHEN generating waveforms THEN the processor SHALL use FFmpeg's `showwavespic` or essentia library
3. WHEN BPM is detected THEN the processor SHALL also calculate beat onset timestamps
4. WHEN analysis completes THEN the processor SHALL update DynamoDB with `waveformUrl`, `beatGrid[]`
5. WHEN analysis fails THEN the processor SHALL mark track as "analysis_partial" and continue
6. WHEN re-processing THEN the system SHALL regenerate waveform/beatgrid and overwrite existing

## Non-Functional Requirements

### Code Architecture and Modularity
- **Single Responsibility Principle**: Waveform component, EQ component, Crossfade logic as separate modules
- **Modular Design**: Audio processing hooks should be reusable across player implementations
- **Dependency Management**: Web Audio API usage isolated in `useAudioProcessor` hook
- **Clear Interfaces**: Define `WaveformData`, `BeatGrid`, `EQSettings` types

### Performance
- Waveform rendering must complete in <100ms for any track length
- Crossfade must not cause audio glitches (use Web Audio scheduling)
- EQ adjustments must apply in <20ms (real-time feel)
- Waveform data must be <50KB per track (compressed JSON)

### Security
- Waveform/beat grid S3 URLs must use CloudFront signed URLs
- EQ presets must be validated (no XSS in preset names)
- Playback speed must be bounded (no arbitrary values from client)

### Reliability
- Missing waveform data must not break playback
- Web Audio context must handle browser autoplay policies
- EQ settings must degrade gracefully if Web Audio unavailable

### Usability
- Waveform must be touch-friendly (swipe to seek on mobile)
- EQ must show visual feedback (current levels)
- Crossfade settings should persist across sessions
- Speed indicator must be obvious but not intrusive
