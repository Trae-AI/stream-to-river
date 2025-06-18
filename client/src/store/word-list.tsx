// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { create } from 'zustand';
import { Word } from '@/type/word';


interface WordListState {
  wordList: Array<Word>;
  initialIndex: number;
  setInitialIndex: (index: number) => void;
  setWordList: (wordList: Array<Word>) => void;
}

const useWordListStore = create<WordListState>((set) => ({
  wordList: [],
  initialIndex: 0,
  setInitialIndex: (index) => set({ initialIndex: index }),
  setWordList: (wordList) => set({ wordList }),
}));

export default useWordListStore;
