// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { ServerConfig } from '@/service/server-config';
import { Word, Sentence, WordPronounce} from '@/type/word';
import AuthService from './auth-service';

interface GetWordListResponse {
  words_list: Word[];
  total: number;
}

export interface NewWordResp{
  new_word_name: string;
  description: string;
  explains: string;
  pronounce_us: WordPronounce;
  pronounce_uk: WordPronounce;
  sentences: Sentence[];
}

class WordService {
  private static instance: WordService;
  private constructor() {
  }

  public static getInstance(): WordService {
    if (!WordService.instance) {
      WordService.instance = new WordService();
    }
    return WordService.instance;
  }

  public async queryWordByName(word: string): Promise<NewWordResp> {
    const url = ServerConfig.getInstance().getFullUrl('/api/word-detail');
    const params = new URLSearchParams();
    params.append('word', word);

    try {
      const response = await fetch(`${url}?${params.toString()}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      return data as NewWordResp;
    } catch (error) {
      console.error('Failed to fetch word list:', error);
      throw error;
    }
  }

  public async queryWordByID (word_id: number): Promise<Word> {
    const url = ServerConfig.getInstance().getFullUrl('/api/word-query');
    const params = new URLSearchParams();
    params.append('word_id', word_id.toString());

    try {
      const response = await fetch(`${url}?${params.toString()}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const data = await response.json();
      if (data.status === -1) {
        throw new Error(data.error_msg);
      }

      return data.word as Word;
    } catch (error) {
      console.error('Failed to fetch word:', error);
      throw error;
    }
  }

  public async getWordList(offset?: number): Promise<GetWordListResponse> {
    const url = ServerConfig.getInstance().getFullUrl('/api/word-list');
    const params = new URLSearchParams();
    const token = await AuthService.getInstance().getToken();
    if (!token) {
      throw new Error('未找到token');
    }

    if (offset !== undefined) {
      params.append('offset', offset.toString());
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

      const data = await response.json();
      return data as GetWordListResponse;
    } catch (error) {
      console.error('Failed to fetch word list:', error);
      throw error;
    }
  }

  public async addWord(wordName: string, tag_id: number): Promise<any> {
    const url = ServerConfig.getInstance().getFullUrl('/api/word-add');
    const body = {
      word: wordName.trim(),
      tag_id: Number(tag_id),
    };
    const token = await AuthService.getInstance().getToken();
    if (!token) {
      throw new Error('未找到token');
    }
    try {
      const response = await fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(body),
      });

      if (!response.ok) {
        throw response;
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Failed to fetch word list:', error);
      throw error;
    }
  }
}

export { GetWordListResponse, WordService };
export default WordService;




