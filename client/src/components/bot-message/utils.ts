// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
//
// SPDX-License-Identifier: MIT

import { cloneDeep, uniq } from 'lodash-es';

type PosMap = Record<string, number>;

export const posCalcRange = (text: string, posList: Array<number>): PosMap => {
  let segmentLen = 0;
  const resMap: PosMap = {};
  let list = cloneDeep(posList);
  list = uniq(list.sort((a, b) => a - b));
  for (let i = 0; i < text.length; i++) {
    if (segmentLen === list[0]) {
      resMap[list[0]] = i;
      list.shift();
    }
    if (
      (text[i].charCodeAt(0) & 0xfc00) === 0xd800 &&
      (text[i + 1].charCodeAt(0) & 0xfc00) === 0xdc00
    ) {
      segmentLen += calculateLength(`${text[i]}${text[i + 1]}`);
      i++;
    } else {
      segmentLen += calculateLength(text[i]);
    }
    if (!list.length) {
      break;
    }
  }
  return resMap;
};

// 辅助函数：计算UTF8字符串的字节长度
function calculateLength(string: string): number {
  let len = 0;
  let c = 0;
  for (let i = 0; i < string.length; ++i) {
    c = string.charCodeAt(i);
    if (c < 128) {
      len += 1;
    } else if (c < 2048) {
      len += 2;
    } else if (
      (c & 0xfc00) === 0xd800 &&
      (string.charCodeAt(i + 1) & 0xfc00) === 0xdc00
    ) {
      ++i;
      len += 4;
    } else {
      len += 3;
    }
  }
  return len;
}
