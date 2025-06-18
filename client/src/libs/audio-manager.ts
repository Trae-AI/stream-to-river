// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

class AudioManager {
  private static instance: AudioManager;
  private currentAudio: HTMLAudioElement | null = null;
  private currentWordId: number | null = null;
  private playingCallbacks: Map<number, (isPlaying: boolean) => void> = new Map();

  static getInstance(): AudioManager {
    if (!AudioManager.instance) {
      AudioManager.instance = new AudioManager();
    }
    return AudioManager.instance;
  }

  // Register playback status callback
  registerCallback(wordId: number, callback: (isPlaying: boolean) => void) {
    this.playingCallbacks.set(wordId, callback);
  }

  // Unregister callback
  unregisterCallback(wordId: number) {
    this.playingCallbacks.delete(wordId);
  }

  // Play audio
  async playAudio(wordId: number, audioUrl: string): Promise<void> {
    try {
      // Stop currently playing audio
      this.stopCurrentAudio();

      // Create new audio instance
      const audio = new Audio(audioUrl);
      this.currentAudio = audio;
      this.currentWordId = wordId;

      // Set audio properties
      audio.preload = 'auto';
      audio.volume = 1.0;

      // Notify playback start
      this.notifyPlayingState(wordId, true);

      // Listen to audio events
      audio.addEventListener('ended', () => {
        this.handleAudioEnd();
      });

      audio.addEventListener('error', (e) => {
        console.error('Audio playback error:', e);
        this.handleAudioEnd();
      });

      // Play audio
      await audio.play();
    } catch (error) {
      console.error('Failed to play audio:', error);
      this.handleAudioEnd();
    }
  }

  // Stop current audio
  private stopCurrentAudio() {
    if (this.currentAudio) {
      this.currentAudio.pause();
      this.currentAudio.currentTime = 0;
      this.currentAudio = null;
    }
    if (this.currentWordId !== null) {
      this.notifyPlayingState(this.currentWordId, false);
      this.currentWordId = null;
    }
  }

  // Handle audio end
  private handleAudioEnd() {
    if (this.currentWordId !== null) {
      this.notifyPlayingState(this.currentWordId, false);
    }
    this.currentAudio = null;
    this.currentWordId = null;
  }

  // Notify playback state change
  private notifyPlayingState(wordId: number, isPlaying: boolean) {
    const callback = this.playingCallbacks.get(wordId);
    if (callback) {
      callback(isPlaying);
    }
  }

  // Check if currently playing
  isPlaying(wordId: number): boolean {
    return this.currentWordId === wordId && this.currentAudio !== null;
  }
}

export default AudioManager;
