// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

export type SuggestedQuestion = {
  id: string;
  text: string;
};

export type Message = {
  id: string;
  text: string;
  creator: "user" | "bot";
  linkWordList?: {
    text: string;
    range: [number, number];
  }[];
};
