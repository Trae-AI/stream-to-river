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
