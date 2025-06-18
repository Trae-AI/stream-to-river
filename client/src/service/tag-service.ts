// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 


import { ServerConfig } from '@/service/server-config';

import { Tag } from '@/type/tag';
import AuthService from './auth-service';

export interface GetTagsResp {
  tags: Tag[];
}

class TagsService {
  private static instance: TagsService;
  private constructor() {
  }

  public static getInstance(): TagsService {
    if (!TagsService.instance) {
      TagsService.instance = new TagsService();
    }
    return TagsService.instance;
  }

  public async getTags(): Promise<GetTagsResp> {
    const url = ServerConfig.getInstance().getFullUrl('/api/tags');
    const params = new URLSearchParams();
    const token = await AuthService.getInstance().getToken();
    if (!token) {
      throw new Error('未找到token');
    }

    try {
      const response = await fetch(`${url}?${params.toString()}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('Failed to fetch word list:', error);
      throw error;
    }
  }

  public async updateWordTag(wordId: number, tagId: number): Promise<void> {
    const url = ServerConfig.getInstance().getFullUrl('/api/word-tag');
    const body = {
      word_id: wordId,
      tag_id: tagId,
    }
    const token = await AuthService.getInstance().getToken();
    if (!token) {
      throw new Error('未找到token');
    }
    try {
      const response = await fetch(`${url}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(body),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      // FIMXME
    } catch (error) {
      console.error('Failed to fetch word list:', error);
      throw error;
    }
  }
}

export default TagsService;
