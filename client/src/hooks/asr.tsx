// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { useEffect, useRef, useState } from "react";
import { AsrService } from "@/service/asr-service";
let t = 0;
export const useAudioRecording = () => {
  const audioContext = useRef<AudioContext | null>(null);
  const audioStream = useRef<MediaStream | null>(null);
  const audioChunks = useRef<Float32Array[]>([]);
  const processor = useRef<ScriptProcessorNode | null>(null);
  const source = useRef<MediaStreamAudioSourceNode | null>(null);
  const asrService = AsrService.getInstance();

  const cleanupAudioResources = () => {
    if (processor.current && audioContext.current) {
      processor.current.disconnect();
      processor.current = null;
    }

    if (source.current) {
      source.current.disconnect();
      source.current = null;
    }

    if (audioStream.current) {
      audioStream.current.getTracks().forEach((track) => track.stop());
      audioStream.current = null;
    }

    if (audioContext.current) {
      if (audioContext.current.state !== "closed") {
        audioContext.current.close();
      }
      audioContext.current = null;
    }

    audioChunks.current = [];
  };

  const startRecording = async () => {
    t = Date.now();
    try {
      // Clean up previous resources
      cleanupAudioResources();

      // Create new audio context
      audioContext.current = new AudioContext({
        sampleRate: 16000,
      });

      // Get microphone stream
      audioStream.current = await navigator.mediaDevices.getUserMedia({
        audio: {
          channelCount: 1,
          sampleRate: 16000,
          sampleSize: 16,
        },
      });

      // Create audio source
      source.current = audioContext.current.createMediaStreamSource(
        audioStream.current
      );
      processor.current = audioContext.current.createScriptProcessor(
        4096,
        1,
        1
      );

      audioChunks.current = [];

      // Process audio data
      processor.current.onaudioprocess = (e) => {
        const inputData = e.inputBuffer.getChannelData(0);
        audioChunks.current.push(new Float32Array(inputData));
      };

      // Connect nodes
      source.current.connect(processor.current);
      processor.current.connect(audioContext.current.destination);

      // In startRecording, after processor.current.connect(audioContext.current.destination);
      if (audioContext.current && audioContext.current.state === "suspended") {
        await audioContext.current.resume();
      }

      return true;
    } catch (error) {
      console.error("Failed to start recording:", error);
      cleanupAudioResources();
      return false;
    }
  };

  const stopRecording = async (): Promise<string> => {
    return new Promise((resolve, reject) => {
      try {
        if (!audioContext.current) {
          throw new Error("No audio context");
        }

        // Merge all audio data
        const length = audioChunks.current.reduce(
          (acc, chunk) => acc + chunk.length,
          0
        );
        const audioData = new Float32Array(length);
        let offset = 0;

        for (const chunk of audioChunks.current) {
          audioData.set(chunk, offset);
          offset += chunk.length;
        }

        // Convert to 16-bit integer
        const pcmData = new Int16Array(audioData.length);
        for (let i = 0; i < audioData.length; i++) {
          const s = Math.max(-1, Math.min(1, audioData[i]));
          pcmData[i] = s < 0 ? s * 0x8000 : s * 0x7fff;
        }

        // Create WAV header
        const wavHeader = new ArrayBuffer(44);
        const view = new DataView(wavHeader);

        // WAV header format
        const writeString = (
          view: DataView,
          offset: number,
          string: string
        ) => {
          for (let i = 0; i < string.length; i++) {
            view.setUint8(offset + i, string.charCodeAt(i));
          }
        };

        writeString(view, 0, "RIFF"); // RIFF identifier
        view.setUint32(4, 36 + pcmData.length * 2, true); // file length
        writeString(view, 8, "WAVE"); // WAVE identifier
        writeString(view, 12, "fmt "); // fmt chunk
        view.setUint32(16, 16, true); // length of fmt chunk
        view.setUint16(20, 1, true); // PCM format
        view.setUint16(22, 1, true); // mono
        view.setUint32(24, 16000, true); // sample rate
        view.setUint32(28, 16000 * 2, true); // byte rate
        view.setUint16(32, 2, true); // block align
        view.setUint16(34, 16, true); // bits per sample
        writeString(view, 36, "data"); // data chunk
        view.setUint32(40, pcmData.length * 2, true); // data length

        // Merge header and audio data
        const wavBlob = new Blob([wavHeader, pcmData], { type: "audio/wav" });

        // Clean up resources
        cleanupAudioResources();

        // Call ASR service
        asrService
          .recognizeAudio(wavBlob, {
            userId: "current-user-id",
            format: "wav",
            sampleRate: 16000,
            channels: 1,
          })
          .then(resolve)
          .catch(reject);
      } catch (error) {
        reject(error);
      }
    });
  };

  const cancelRecording = () => {
    cleanupAudioResources();
  };

  useEffect(() => {
    return () => {
      cancelRecording();
    };
  }, []);

  return {
    startRecording,
    stopRecording,
    cancelRecording,
  };
};
