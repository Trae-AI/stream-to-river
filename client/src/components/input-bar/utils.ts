// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

export const compressImage = (img: HTMLImageElement): string => {
  const canvas = document.createElement('canvas');
  const MAX_WIDTH = 1280;
  const MAX_HEIGHT = 720;

  let width = img.width;
  let height = img.height;

  if (width > MAX_WIDTH || height > MAX_HEIGHT) {
    const ratio = Math.min(MAX_WIDTH / width, MAX_HEIGHT / height);
    width = Math.floor(width * ratio);
    height = Math.floor(height * ratio);
  }

  canvas.width = width;
  canvas.height = height;

  const ctx = canvas.getContext('2d');
  if (!ctx) throw new Error('Failed to get 2d context');

  ctx.imageSmoothingEnabled = true;
  ctx.imageSmoothingQuality = 'high';

  ctx.drawImage(img, 0, 0, width, height);
  return canvas.toDataURL('image/jpeg', 0.6);
};

export const getBase64FromUrl = (url: string): Promise<string> => {
  return new Promise((resolve, reject) => {
    const img = new Image();
    img.crossOrigin = 'anonymous';

    img.onload = () => {
      try {
        const compressedBase64 = compressImage(img);
        resolve(compressedBase64);
      } catch (error) {
        reject(error);
      }
    };

    img.onerror = reject;
    img.src = url;
  });
};
