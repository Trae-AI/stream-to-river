// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { ServerConfig } from '@/service/server-config';

class Image2TextService {
  private static instance: Image2TextService;

  private constructor() {}

  public static getInstance(): Image2TextService {
    if (!Image2TextService.instance) {
      Image2TextService.instance = new Image2TextService();
    }
    return Image2TextService.instance;
  }

  public async image2Text(base64Image: string): Promise<string> {
    const url = ServerConfig.getInstance().getFullUrl('/api/image2text');

    try {
      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          base64: base64Image
        })
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();

      if (data.code !== 200) {
        throw new Error(data.message || 'Image to text conversion failed');
      }

      return data.data as string;

    } catch (error) {
      console.error('Failed to convert image to text:', error);
      throw error;
    }
  }
}

export default Image2TextService;
