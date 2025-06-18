// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { create } from 'zustand';
import { Tag } from '@/type/tag';

interface SystemTagsState {
  tags: Array<Tag>;
  setSystemTags: (tags: Array<Tag>) => void;
}

const useSystemTagsStore = create<SystemTagsState>((set) => ({
  tags: [],
  setSystemTags: (tags) => set({ tags }),
}));

export default useSystemTagsStore;
