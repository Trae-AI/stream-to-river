// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

// Option item interface
export interface OptionItem {
  description: string;
  answer_list_id: number;
}

// Review question interface
export interface ReviewQuestion {
  question: string;
  word_id: number;
  question_type: number;
  options: OptionItem[];
}

