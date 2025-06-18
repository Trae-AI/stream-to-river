// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

export function extractStyleExamples(htmlString: string) {
  const parser = new DOMParser();
  const doc = parser.parseFromString(htmlString, 'text/html');

  const examples = doc.getElementsByClassName('style-text style-example');
  const results: {english: string; chinese: string}[] = [];

  for (let example of examples) {
      const words = example.getElementsByClassName('style-word');

      const englishText = Array.from(words)
          .map(word => word.textContent || '')
          .join(' ')
          .trim();

      const chineseElement = example.querySelector('.style-chn');
      const chineseText = chineseElement ? chineseElement.textContent || '' : '';

      results.push({
          english: englishText,
          chinese: chineseText
      });
  }

  return results;
}
