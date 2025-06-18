// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

import type { UserConfigExport } from "@tarojs/cli"

export default {
  devtool: 'source-map',
   logger: {
    quiet: false,
    stats: true
  },
  mini: {
  },
  h5: {
  },
} satisfies UserConfigExport<'webpack5'>
