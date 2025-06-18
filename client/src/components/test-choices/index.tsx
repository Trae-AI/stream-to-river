// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import React, { useState, forwardRef, useImperativeHandle } from "react";
import { View, Text } from "@tarojs/components";
import "./index.less";
import AutoSizeByCharCount from "../auto-size-text";

interface Choice {
  id: string;
  text: string;
}

interface TestChoicesProps {
  type: "word-to-meaning" | "meaning-to-word";
  question: string;
  pronunciation?: {
    uk: string;
    us: string;
  };
  choices: Choice[];
  onSelect?: (choiceId: string) => void;
  correctChoiceId: number;
}

const TestChoices = forwardRef<{ reviewAnswer: (correctAnswerId: number) => void }, TestChoicesProps>((
  {
    type,
    question,
    pronunciation,
    choices,
    onSelect,
    correctChoiceId,
  },
  ref
) => {
  const [selectedId, setSelectedId] = useState<string>("");

  const handleChoiceSelect = (choice: Choice) => {
    if (selectedId) return;
    onSelect?.(choice.id);
    setSelectedId(choice.id);
  };

  useImperativeHandle(ref, () => ({
    reviewAnswer: (correctAnswerId: number) => {
      setSelectedId(correctAnswerId.toString());
    }
  }));

  const getChoiceClassName = (choice: Choice) => {
    const baseClass = "test-choices__option";

    if (!selectedId || selectedId.trim().length === 0) {
      return baseClass;
    }
    if (selectedId !== choice.id && correctChoiceId === parseInt(choice.id)) {
      return `${baseClass} ${baseClass}--correct`;
    }

    if (selectedId === choice.id) {
      if (correctChoiceId != -1) {
        if (parseInt(selectedId.trim()) === correctChoiceId) {
          return `${baseClass} ${baseClass}--correct`;
        } else {
          return `${baseClass} ${baseClass}--wrong`;
        }
      }
    }

    return baseClass;
  };

  return (
    <View className="test-choices">
      {/* 题目区域 */}
      <View className="test-choices__question">
        <Text className="test-choices__question-text">
          <AutoSizeByCharCount maxFontSize={40} minFontSize={18} style={{ paddingLeft: '16px', paddingRight: '16px'}}>
            {question}
          </AutoSizeByCharCount>
        </Text>

        {/* 发音信息（仅英文单词显示） */}
        {type === "word-to-meaning" && pronunciation && (
          <View className="test-choices__pronunciation">
            <View className="test-choices__pronunciation-item">
              <Text className="test-choices__pronunciation-label">英</Text>
              <Text className="test-choices__pronunciation-text">
                /{pronunciation.uk}/
              </Text>
              <View className="test-choices__speaker">🔊</View>
            </View>
            <View className="test-choices__pronunciation-item">
              <Text className="test-choices__pronunciation-label">美</Text>
              <Text className="test-choices__pronunciation-text">
                /{pronunciation.us}/
              </Text>
              <View className="test-choices__speaker">🔊</View>
            </View>
          </View>
        )}
      </View>

      {/* 提示文本 */}
      <View className="test-choices__hint">
        <Text className="test-choices__hint-text">
          {type === "word-to-meaning"
            ? "选择正确的中文释义"
            : "选择正确的英文单词"}
        </Text>
      </View>

      {/* 选项区域 */}
      <View className="test-choices__options">
        {choices.filter(item => Boolean(item.text)).sort((a, b) => a.id.localeCompare(b.id)).map((choice) => (
          <View
            key={choice.id}
            className={getChoiceClassName(choice)}
            onClick={() => handleChoiceSelect(choice)}
          >
            <Text className="test-choices__option-text">{choice.text}</Text>
          </View>
        ))}
      </View>
    </View>
  );
});

export default TestChoices;
