// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

// services/asr-service.ts
import { ServerConfig } from "@/service/server-config";

export type AudioFormat = "wav" | "mp3" | "ogg" | "pcm";

interface AsrResponse {
  code: number;
  message: string;
  data: {
    reqid: string;
    code: number;
    message: string;
    sequence: number;
    result: {
      text: string;
      confidence: number;
    }[];
  };
}

interface RecognizeOptions {
  userId: string;
  format: AudioFormat;
  sampleRate?: number;
  channels?: number;
}

class AsrService {
  private static instance: AsrService;

  private constructor() {}

  public static getInstance(): AsrService {
    if (!AsrService.instance) {
      AsrService.instance = new AsrService();
    }
    return AsrService.instance;
  }

  public async recognizeAudio(
    audioBlob: Blob,
    options: RecognizeOptions
  ): Promise<string> {
    const url = ServerConfig.getInstance().getFullUrl("/api/asrrecognize");

    const params = new URLSearchParams();
    params.append("format", options.format);

    try {
      const response = await fetch(`${url}?${params.toString()}`, {
        method: "POST",
        body: audioBlob,
        headers: {
          "Content-Type": "application/octet-stream",
        },
      });

      if (!response.ok) {
        throw response;
      }

      const data = (await response.json()) as AsrResponse;

      if (data?.code !== 0) {
        throw data
      }

      if (data?.data?.code !== 1000) {
        throw data
      }

      return data.data.result.map((result) => result.text).join("");
    } catch (error) {
      console.error("Failed to recognize audio:", error);
      throw error;
    }
  }
}

export { AsrResponse, AsrService, RecognizeOptions };
export default AsrService;
