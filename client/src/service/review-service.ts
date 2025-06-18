// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { ServerConfig } from '@/service/server-config';
import AuthService from './auth-service';

export enum QuestionType {
  CHOOSE_CN = 1,    // Choose the correct Chinese meaning
  CHOOSE_EN = 2,    // Choose the correct English meaning
  PRONOUNCE_CHOOSE = 3, // Choose the correct Chinese meaning based on pronunciation
  FILL_IN_BLANK = 4,    // Fill in the blank question
}

export interface OptionItem {
  description: string;
  answer_list_id: number;
}

export interface ReviewQuestion {
  question: string;
  word_id: number;
  question_type: number;
  options: OptionItem[];
  show_info: string[];
}

export interface ReviewProgressResp {
  pending_review_count: number;
  completed_review_count: number;
  last_update_time: number;
  all_completed_count: number;
  total_words_count: number;
}

export interface ReviewListResp {
  total_num: string;
  questions: ReviewQuestion[];
}

export interface AnswerPuzzleReq {
  answer_id?: string;
  word_id: string;
  question_type: number;
  filled_name?: string;
}

export interface AnswerPuzzleResp {
  is_correct: boolean;
  message: string;
  correct_answer_id: number;
}

class ReviewPuzzleService {
  private static instance: ReviewPuzzleService;
  private constructor() {
  }

  public static getInstance(): ReviewPuzzleService {
    if (!ReviewPuzzleService.instance) {
      ReviewPuzzleService.instance = new ReviewPuzzleService();
    }
    return ReviewPuzzleService.instance;
  }


  public async getReviewPuzzles(): Promise<ReviewListResp> {
    const url = ServerConfig.getInstance().getFullUrl('/api/review-list');
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

  public async getReviewProgress(): Promise<ReviewProgressResp> {
    const url = ServerConfig.getInstance().getFullUrl('/api/review-progress');
    const token = await AuthService.getInstance().getToken();
    if (!token) {
      throw new Error('未找到token');
    }

    try {
      const response = await fetch(`${url}`, {
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

  public async submitPuzzleAnswer(req: AnswerPuzzleReq): Promise<AnswerPuzzleResp> {
    const url = ServerConfig.getInstance().getFullUrl('/api/answer');
    const params = new URLSearchParams();
    const body = {
      answer_id: Number(req.answer_id),
      word_id: Number(req.word_id),
      question_type: Number(req.question_type),
      filled_name: req.filled_name || '',
    }
    const token = await AuthService.getInstance().getToken();
    if (!token) {
      throw new Error('未找到token');
    }

    try {
      const response = await fetch(`${url}?${params.toString()}`, {
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

      return await response.json();
    } catch (error) {
      console.error('Failed to fetch word list:', error);
      throw error;
    }
  }
}

export default ReviewPuzzleService;
