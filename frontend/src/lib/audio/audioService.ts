import Hls from 'hls.js';

export interface AudioCallbacks {
  onTimeUpdate?: (currentTime: number, duration: number) => void;
  onEnded?: () => void;
  onError?: (error: Error) => void;
  onDurationChange?: (duration: number) => void;
  onLoadStart?: () => void;
  onCanPlay?: () => void;
}

/**
 * AudioService - Singleton service for HLS.js audio playback
 * Provides adaptive bitrate streaming with fallback to direct file playback
 */
export class AudioService {
  private static instance: AudioService | null = null;
  private audio: HTMLAudioElement | null = null;
  private hls: Hls | null = null;
  private callbacks: AudioCallbacks = {};
  private currentUrl: string | null = null;

  private constructor() {
    this.initAudioElement();
  }

  static getInstance(): AudioService {
    if (!AudioService.instance) {
      AudioService.instance = new AudioService();
    }
    return AudioService.instance;
  }

  static resetInstance(): void {
    if (AudioService.instance) {
      AudioService.instance.destroy();
      AudioService.instance = null;
    }
  }

  private initAudioElement(): void {
    this.audio = new Audio();
    this.audio.preload = 'metadata';

    this.audio.addEventListener('timeupdate', () => {
      this.callbacks.onTimeUpdate?.(
        this.audio?.currentTime ?? 0,
        this.audio?.duration ?? 0
      );
    });

    this.audio.addEventListener('ended', () => {
      this.callbacks.onEnded?.();
    });

    this.audio.addEventListener('error', () => {
      const error = new Error(
        `Audio error: ${this.audio?.error?.message || 'Unknown error'}`
      );
      this.callbacks.onError?.(error);
    });

    this.audio.addEventListener('durationchange', () => {
      if (this.audio?.duration && !isNaN(this.audio.duration)) {
        this.callbacks.onDurationChange?.(this.audio.duration);
      }
    });

    this.audio.addEventListener('loadstart', () => {
      this.callbacks.onLoadStart?.();
    });

    this.audio.addEventListener('canplay', () => {
      this.callbacks.onCanPlay?.();
    });
  }

  setCallbacks(callbacks: AudioCallbacks): void {
    this.callbacks = callbacks;
  }

  /**
   * Load a URL - automatically detects HLS vs direct file
   */
  load(url: string): void {
    if (!this.audio) {
      this.initAudioElement();
    }

    // Clean up existing HLS instance
    if (this.hls) {
      this.hls.destroy();
      this.hls = null;
    }

    this.currentUrl = url;

    // Check if URL is HLS (.m3u8)
    const isHLS = url.includes('.m3u8');

    if (isHLS && Hls.isSupported()) {
      // Use HLS.js for adaptive bitrate streaming
      this.hls = new Hls({
        enableWorker: true,
        lowLatencyMode: false,
        backBufferLength: 90,
      });

      this.hls.loadSource(url);
      this.hls.attachMedia(this.audio!);

      this.hls.on(Hls.Events.ERROR, (_event: string, data: { fatal: boolean; type: string; details: string }) => {
        if (data.fatal) {
          switch (data.type) {
            case Hls.ErrorTypes.NETWORK_ERROR:
              // Try to recover network error
              this.hls?.startLoad();
              break;
            case Hls.ErrorTypes.MEDIA_ERROR:
              // Try to recover media error
              this.hls?.recoverMediaError();
              break;
            default:
              // Cannot recover, destroy and notify
              this.callbacks.onError?.(
                new Error(`HLS fatal error: ${data.details}`)
              );
              break;
          }
        }
      });
    } else if (isHLS && this.audio?.canPlayType('application/vnd.apple.mpegurl')) {
      // Native HLS support (Safari)
      this.audio!.src = url;
    } else {
      // Direct file playback
      this.audio!.src = url;
    }
  }

  async play(): Promise<void> {
    if (!this.audio) return;

    try {
      await this.audio.play();
    } catch (error) {
      // Handle autoplay restrictions
      if (error instanceof Error && error.name === 'NotAllowedError') {
        this.callbacks.onError?.(
          new Error('Playback blocked. User interaction required.')
        );
      } else {
        throw error;
      }
    }
  }

  pause(): void {
    this.audio?.pause();
  }

  seek(time: number): void {
    if (this.audio && isFinite(time) && time >= 0) {
      this.audio.currentTime = Math.min(time, this.audio.duration || time);
    }
  }

  setVolume(level: number): void {
    if (this.audio) {
      this.audio.volume = Math.max(0, Math.min(1, level));
    }
  }

  getVolume(): number {
    return this.audio?.volume ?? 1;
  }

  getCurrentTime(): number {
    return this.audio?.currentTime ?? 0;
  }

  getDuration(): number {
    return this.audio?.duration ?? 0;
  }

  isPlaying(): boolean {
    return this.audio ? !this.audio.paused : false;
  }

  getCurrentUrl(): string | null {
    return this.currentUrl;
  }

  destroy(): void {
    if (this.hls) {
      this.hls.destroy();
      this.hls = null;
    }

    if (this.audio) {
      this.audio.pause();
      this.audio.src = '';
      this.audio.load();
      this.audio = null;
    }

    this.currentUrl = null;
    this.callbacks = {};
  }
}
