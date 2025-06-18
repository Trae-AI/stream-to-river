// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { useEffect, useRef, useState } from 'react';
import {debounce} from 'lodash-es';

const AutoSizeByCharCount = ({
  children,
  maxFontSize = 20,
  minFontSize = 12,
  className = "",
  style = {},
  updateDelay = 100
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const [fontSize, setFontSize] = useState(minFontSize);

  const adjustSize = debounce(() => {
    if (!containerRef.current) return;

    const text = containerRef.current.innerText;
    const lines = text.split('\n');

    const maxChars = Math.max(...lines.map(line => {
      return Array.from(line).reduce((count, char) => {
        return count + (/[\u4e00-\u9fa5]/.test(char) ? 1 : 0.5);
      }, 0);
    }));

    const baseFontSize = maxFontSize - ((maxChars - 1) * (maxFontSize - minFontSize) / 20);
    const newFontSize = Math.max(
      minFontSize,
      Math.min(
        maxFontSize,
        baseFontSize
      )
    );

    setFontSize(newFontSize);
  }, updateDelay);

  useEffect(() => {
    const resizeObserver = new ResizeObserver(adjustSize);
    resizeObserver.observe(containerRef.current!);

    adjustSize();

    return () => {
      resizeObserver.disconnect();
      adjustSize.cancel();
    };
  }, [children, maxFontSize, minFontSize]);

  return (
    <div
      ref={containerRef}
      className={className}
      style={{
        fontSize: `${fontSize}px`,
        lineHeight: '1.5',
        wordWrap: 'break-word',
        whiteSpace: 'pre-wrap',
        overflow: 'auto',
        ...style
      }}
    >
      {children}
    </div>
  );
};

export default AutoSizeByCharCount;
