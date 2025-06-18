// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

export default defineAppConfig({
  pages: [
    'pages/index/index',
    'pages/study_main/index',
    'pages/words_list/index',
    'pages/words_study/index',
    'pages/login/index',
  ],
  window: {
    backgroundTextStyle: 'light',
    navigationStyle: 'custom',
  },
  permission: {
  },
  animation: false,
})
