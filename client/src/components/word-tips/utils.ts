// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

export function extractStyleExamples(htmlString: string) {
  const parser = new DOMParser();
  const doc = parser.parseFromString(htmlString, 'text/html');

  const examples = doc.getElementsByClassName('style-text style-example');
  const results: {english: string; chinese: string}[] = [];

  for (let example of examples) {
      const words = example.getElementsByClassName('style-word');
      const boldTexts = example.getElementsByClassName('style-bold');

      let englishText = Array.from(words)
          .map(word => word.textContent || '')
          .join(' ');

      Array.from(boldTexts).forEach(bold => {
          englishText += ' ' + (bold.textContent || '');
      });

      englishText = englishText.trim();

      const chineseElement = example.querySelector('.style-chn');
      const chineseText = chineseElement ? chineseElement.textContent || '' : '';

      results.push({
          english: englishText,
          chinese: chineseText
      });
  }

  return results;
}


export function removeWhitespace(word: string): string {
  return word.replace(/\s/g, '');
}
