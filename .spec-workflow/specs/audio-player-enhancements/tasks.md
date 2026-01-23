# Tasks Document: Audio Analysis & Player Enhancements

## Task Overview
| Task | Description | Estimated Files |
|------|-------------|-----------------|
| 1 | Waveform Backend Generation | 2 |
| 2 | Waveform Frontend Component | 2 |
| 3 | Beat Grid Generation | 1 |
| 4 | Web Audio Processor Hook | 1 |
| 5 | Equalizer Component | 2 |
| 6 | Gapless Playback | 1 |
| 7 | Crossfade Manager | 2 |
| 8 | Playback Speed Control | 1 |
| 9 | Player Store Updates | 1 |
| 10 | Tests | 4 |

---

- [ ] 1. Add Waveform generation to backend processor
  - Files: backend/cmd/processor/metadata/waveform.go, backend/internal/models/track.go (modify)
  - Use FFmpeg to extract waveform peaks (100 samples/second)
  - Store waveform JSON in S3: waveforms/{trackId}.json
  - Update Track model with waveformUrl field
  - Add to metadata processing pipeline
  - Purpose: Generate waveform data during track processing
  - _Leverage: backend/cmd/processor/metadata/main.go, backend/internal/repository/s3.go_
  - _Requirements: 1.1, 1.2, 7.1, 7.2, 7.5_
  - _Prompt: Implement the task for spec audio-player-enhancements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Go Backend Developer with audio processing expertise | Task: Add waveform generation using FFmpeg. Extract 100 peaks/second as normalized floats. Store as JSON in S3. Update Track model with waveformUrl. Integrate into existing metadata processor | Restrictions: Keep waveform JSON under 50KB, use efficient FFmpeg command, handle errors gracefully | Success: Waveforms generated for uploaded tracks, JSON stored in S3, Track updated | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 2. Create Waveform display component
  - Files: frontend/src/components/player/Waveform.tsx, frontend/src/hooks/useWaveform.ts
  - Fetch waveform JSON from S3 (via CloudFront)
  - Render waveform using Canvas or SVG
  - Show playback progress with accent color
  - Handle hover for timestamp preview
  - Handle click for seek
  - Purpose: Visual waveform display in player
  - _Leverage: frontend/src/components/player/PlayerBar.tsx_
  - _Requirements: 1.3, 1.4, 1.5, 1.6, 1.7_
  - _Prompt: Implement the task for spec audio-player-enhancements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React/TypeScript Developer with canvas expertise | Task: Create Waveform component using Canvas. Fetch waveform data via hook. Render peaks, show progress overlay, handle hover (show time), click (seek). Fallback to progress bar if no data | Restrictions: Render in <100ms, responsive sizing, touch-friendly, accessible | Success: Waveform displays, progress shows, seek works, mobile-friendly | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 3. Add Beat Grid generation and display
  - File: backend/cmd/processor/metadata/beatgrid.go
  - Calculate beat timestamps from BPM detection
  - Store beat grid array in Track metadata
  - Display beat markers on waveform (frontend)
  - Emphasize downbeats (every 4th)
  - Purpose: Visual beat markers for DJ mixing
  - _Leverage: Existing BPM detection code_
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6_
  - _Prompt: Implement the task for spec audio-player-enhancements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Audio Analysis Developer | Task: Generate beat grid from BPM. Calculate beat timestamps starting from first detected beat. Store as array in Track.beatGrid. Display vertical markers on Waveform component with emphasized downbeats | Restrictions: Handle variable BPM (show warning), zoom-aware display, don't clutter at low zoom | Success: Beat markers align with audio, downbeats highlighted, variable BPM handled | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 4. Create Web Audio Processor hook
  - File: frontend/src/hooks/useAudioProcessor.ts
  - Initialize AudioContext
  - Create processing chain: source → EQ filters → analyzer → destination
  - Expose setEQ, setPlaybackRate, getAnalyzerData methods
  - Handle browser autoplay policies
  - Purpose: Web Audio API integration for advanced audio features
  - _Leverage: frontend/src/lib/store/playerStore.ts_
  - _Requirements: 5.4, 5.5, 6.4_
  - _Prompt: Implement the task for spec audio-player-enhancements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Web Audio API Developer | Task: Create useAudioProcessor hook. Initialize AudioContext on user interaction. Create filter chain with 3 BiquadFilterNodes (low/mid/high). Connect Howler audio to processing chain. Handle autoplay policies | Restrictions: Support all browsers, handle context state, clean up on unmount | Success: Audio routed through Web Audio, EQ works in real-time, no audio glitches | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 5. Create Equalizer component
  - Files: frontend/src/components/player/Equalizer.tsx, frontend/src/lib/store/eqStore.ts
  - 3-band EQ with sliders (Low 60Hz, Mid 1kHz, High 12kHz)
  - Range: -12dB to +12dB
  - Presets: Flat, Bass Boost, Vocal, Electronic
  - Persist settings to localStorage
  - Purpose: User-adjustable audio EQ
  - _Leverage: frontend/src/hooks/useAudioProcessor.ts_
  - _Requirements: 5.1, 5.2, 5.3, 5.5, 5.6_
  - _Prompt: Implement the task for spec audio-player-enhancements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React Developer | Task: Create Equalizer component with 3 vertical sliders (-12 to +12 dB). Add preset dropdown. Connect to useAudioProcessor for real-time adjustment. Persist to localStorage. Show visual feedback | Restrictions: Use DaisyUI components, accessible (keyboard adjustable), show dB values | Success: EQ adjusts audio in real-time, presets apply correctly, persists across sessions | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 6. Implement Gapless playback
  - File: frontend/src/lib/store/playerStore.ts (modify)
  - Preload next track 5 seconds before current ends
  - Use Howler HTML5 mode for gapless
  - Seamlessly transition (<10ms gap)
  - Handle manual skip (don't apply gapless)
  - Purpose: Seamless album/playlist playback
  - _Leverage: Howler.js documentation_
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6_
  - _Prompt: Implement the task for spec audio-player-enhancements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Audio Developer with Howler.js expertise | Task: Implement gapless playback in playerStore. Preload next track (create Howl instance) 5s before end. On track end, immediately play preloaded. Handle preload failure gracefully. Skip button bypasses gapless | Restrictions: HTML5 audio mode only, manage memory (unload old tracks), handle edge cases | Success: Albums play without gaps, preload doesn't cause stutters, fallback works | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 7. Create Crossfade manager
  - Files: frontend/src/lib/audio/crossfade.ts, frontend/src/lib/store/playerStore.ts (modify)
  - Implement equal-power crossfade curve
  - Configurable duration: 0, 1, 3, 5, 10 seconds
  - Start crossfade N seconds before track end
  - Handle skip during crossfade
  - Don't apply during repeat-one
  - Purpose: Smooth transitions between tracks
  - _Leverage: frontend/src/lib/store/playerStore.ts_
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7_
  - _Prompt: Implement the task for spec audio-player-enhancements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: Audio Developer | Task: Create CrossfadeManager class. Implement equal-power crossfade (sqrt curve). Support configurable duration. Handle edge cases: skip during fade, repeat-one mode, short tracks. Integrate with playerStore | Restrictions: Smooth curves (no clicks), handle concurrent fades, memory cleanup | Success: Crossfade sounds smooth, configurable, handles all edge cases | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 8. Add Playback speed control
  - File: frontend/src/components/player/SpeedControl.tsx
  - Speed options: 0.5x, 0.75x, 1x, 1.25x, 1.5x, 2x
  - Use Howler rate() method
  - Show current speed indicator when not 1x
  - Optional pitch preservation toggle
  - Purpose: Variable playback speed
  - _Leverage: frontend/src/lib/store/playerStore.ts_
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_
  - _Prompt: Implement the task for spec audio-player-enhancements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React Developer | Task: Create SpeedControl component with dropdown for speed selection. Call Howler rate() method. Show indicator when not 1x. Add to playerStore state. Duration display shows actual duration (not adjusted) | Restrictions: Smooth speed transitions, persist preference, clear visual indicator | Success: Speed changes in real-time, indicator visible, persists across sessions | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 9. Update Player Store with all new state
  - File: frontend/src/lib/store/playerStore.ts (modify)
  - Add state: waveformData, beatGrid, eq, crossfade, playbackSpeed, gapless
  - Add actions: setEQ, setCrossfade, setPlaybackSpeed, setGapless
  - Persist preferences to localStorage
  - Load preferences on init
  - Purpose: Central state management for player features
  - _Leverage: Existing playerStore.ts_
  - _Requirements: All_
  - _Prompt: Implement the task for spec audio-player-enhancements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React/Zustand Developer | Task: Extend playerStore with new state fields and actions. Persist EQ, crossfade, speed, gapless settings to localStorage. Load on init. Integrate with all new components | Restrictions: Maintain existing functionality, type-safe, efficient updates | Success: All new features integrated in store, persistence works, no regressions | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_

- [ ] 10. Write tests for audio features
  - Files: frontend/src/components/player/__tests__/Waveform.test.tsx, frontend/src/components/player/__tests__/Equalizer.test.tsx, frontend/src/lib/audio/__tests__/crossfade.test.ts, frontend/src/lib/store/__tests__/playerStore.test.ts (extend)
  - Test waveform rendering and interaction
  - Test EQ slider behavior and presets
  - Test crossfade timing and curves
  - Test playerStore state transitions
  - Purpose: Ensure audio features are reliable
  - _Leverage: Existing test patterns_
  - _Requirements: All_
  - _Prompt: Implement the task for spec audio-player-enhancements, first run spec-workflow-guide to get the workflow guide then implement the task: Role: React Test Engineer | Task: Write tests for Waveform (rendering, seek), Equalizer (sliders, presets), CrossfadeManager (timing, edge cases), playerStore (new state/actions). Mock Web Audio API and Howler | Restrictions: Mock audio APIs properly, test user interactions, accessibility tests | Success: 80%+ coverage, all tests pass, edge cases covered | After completion: Mark task in progress in tasks.md with [-], implement, log with log-implementation tool, then mark [x] complete_
