// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

export interface WordPronounce{
  phonetic: string;
  url: string;
}

export interface Sentence{
  text: string;
  audio_url: string;
}
export interface Word {
  word_id: number;
  word_name: string;
  description: string;
  explains: string;
  pronounce_uk: WordPronounce;
  pronounce_us: WordPronounce;
  tag_id: number;
  level: number;
  max_level: number;
  sentences: Sentence[];
}
