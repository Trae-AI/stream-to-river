// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

import { defineConfig, type UserConfigExport } from '@tarojs/cli'
import TsconfigPathsPlugin from 'tsconfig-paths-webpack-plugin'
import devConfig from './dev'
import prodConfig from './prod'
import path from 'path';

const cdnPrefix = process.env.CDN_OUTER_CN;
const cdnPublicPath = cdnPrefix
  ? '//' + cdnPrefix + '/' + process.env.CDN_PATH_PREFIX + '/'
  : '/';

// https://taro-docs.jd.com/docs/next/config#defineconfig-helper-function
export default defineConfig<'webpack5'>(async (merge) => {
  const baseConfig: UserConfigExport<'webpack5'> = {
    devtool: 'source-map',
    projectName: 'client',
    date: '2025-5-30',
    designWidth: 375,
    deviceRatio: {
      640: 2.34 / 2,
      750: 1,
      375: 2,
      828: 1.81 / 2
    },
    sourceRoot: 'src',
    outputRoot: 'dist',
    plugins: [],
    alias: {
      '@/assets': require('path').resolve(__dirname, '..', 'src/assets'),
    },
    copy: {
      patterns: [
        { from: 'static', to: 'dist/static' }
      ],
      options: {
      }
    },
    framework: 'react',
    compiler: 'webpack5',
    cache: {
      enable: false // Webpack persistent cache configuration, recommended to enable. For default configuration, please refer to: https://docs.taro.zone/docs/config-detail#cache
    },
    mini: {
      postcss: {
        pxtransform: {
          enable: true,
          config: {

          }
        },
        cssModules: {
          enable: false, // Default is false. Set to true if you need to use css modules feature
          config: {
            namingPattern: 'module', // Transform mode, value can be global/module
            generateScopedName: '[name]__[local]___[hash:base64:5]'
          }
        }
      },
      webpackChain(chain) {
        chain.resolve.plugin('tsconfig-paths').use(TsconfigPathsPlugin)
      },
    },
    h5: {
      imageUrlLoaderOption: { limit: false },
      publicPath: cdnPublicPath,
      staticDirectory: 'static',
      output: {
        filename: 'js/[name].[fullhash:8].js',
        chunkFilename: 'js/[name].[chunkhash:8].js',
      },
      htmlPluginOption: {
        favicon: path.join(__dirname, '../static/images/logo.ico'),
      },

      miniCssExtractPluginOption: {
        ignoreOrder: true,
        filename: 'css/[name].[contenthash].css',
        chunkFilename: 'css/[name].[chunkhash].css'
      },
      postcss: {
        autoprefixer: {
          enable: true,
          config: {}
        },
        cssModules: {
          enable: false,
          config: {
            namingPattern: 'module',
            generateScopedName: '[name]__[local]___[hash:base64:5]'
          }
        }
      },
      webpackChain(chain) {
        chain.resolve.plugin('tsconfig-paths').use(TsconfigPathsPlugin)
        // SVG processing rules (prioritize conversion to React components)
        chain.module
          .rule('svgr')
          .enforce('post')
          .test(/\.svg$/)
          .exclude
          .add(path.join(__dirname, '../src/assets/bg'))
          .end()
          .use('@svgr/webpack')
          .loader('@svgr/webpack')
          .options({
            icon: true,
            svgoConfig: {
              plugins: [

              ]
            }
          })
          .end()
          .use('file-loader')
          .loader('file-loader')
          .options({
            name: 'static/svg/[name].[hash:8].[ext]', // Output path and filename format
            publicPath: cdnPublicPath // Use CDN public path
          }).end()

        // Delete default image processing rules
        chain.module.rules.delete('image')
        // Processing rules for other image types (png/jpg/jpeg/gif/bpm)
        chain.module
          .rule('image')
          .test(/\.(png|jpe?g|gif|bpm|webp)(\?.*)?$/)
          .use('file-loader')
          .loader('file-loader')
          .options({
            name: 'static/images/[name].[hash:8].[ext]', // Output path and filename format
            publicPath: cdnPublicPath // Use CDN public path
          })
          .end()
        chain.module
          .rule('bg-svg')
          .enforce('pre')
          .test(/\.svg$/)
          .include
          .add(path.join(__dirname, '../src/assets/bg'))
          .end()
          .use('url-loader')
          .loader('url-loader')
          .options({
            name: 'static/svg-bg/[name].[hash:8].[ext]', // Output path and filename format
            publicPath: cdnPublicPath // Use CDN public path
          })
          .end()
      }
    },
    rn: {
      appName: 'taroDemo',
      postcss: {
        cssModules: {
          enable: false,
        }
      }
    }
  }


  if (process.env.NODE_ENV === 'development') {
    // Local development build configuration (no minification or obfuscation)
    return merge({}, baseConfig, devConfig)
  }
  // Production build configuration (minification and obfuscation enabled by default)
  return merge({}, baseConfig, prodConfig)
})
