// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import React, { useMemo, useEffect, useState } from "react";
import { View, Text, Image } from "@tarojs/components";
import WordPronounceComp from "@/components/word-pronunce";
import TagSelection from "@/components/tag-selection";

import { Sentence, Word } from "@/type/word";

import "./index.less";
import TagsService from "@/service/tag-service";
import WordService from "@/service/word-service";

interface WordCardProps {
  word: Word;
}

const generateProgressPath = (
  progress: number,
  centerX: number = 9,
  centerY: number = 9.5,
  radius: number = 8
): string => {
  // Ensure progress is between 0 and 1
  const normalizedProgress = Math.max(0, Math.min(1, progress));

  // Return empty path if progress is 0
  if (normalizedProgress === 0) {
    return "";
  }

  // Return complete circle if progress is 1
  if (normalizedProgress === 1) {
    return `M ${centerX - radius} ${centerY} A ${radius} ${radius} 0 1 1 ${
      centerX + radius
    } ${centerY} A ${radius} ${radius} 0 1 1 ${centerX - radius} ${centerY}`;
  }

  // Calculate angle (starting from top, clockwise)
  const angle = normalizedProgress * 2 * Math.PI - Math.PI / 2;

  // Calculate end coordinates
  const endX = centerX + radius * Math.cos(angle);
  const endY = centerY + radius * Math.sin(angle);

  // Determine if large arc flag is needed
  const largeArcFlag = normalizedProgress > 0.5 ? 1 : 0;

  // Generate path
  return `M ${centerX} ${centerY} L ${centerX} ${
    centerY - radius
  } A ${radius} ${radius} 0 ${largeArcFlag} 1 ${endX} ${endY} Z`;
};

const WordCard: React.FC<WordCardProps> = ({ word }) => {
  const progress = useMemo(() => {
    return word.level / word.max_level;
  }, [word.level, word.max_level]);
  const progressPath = generateProgressPath(progress);
  const [sentences, setSentences] = useState<Sentence[]>([]);

  useEffect(() => {
    if (!word.word_id) {
      return;
    }
    (async () => {
      const res = await WordService.getInstance().queryWordByName(word.word_name);
      setSentences(res.sentences);
    })()
  }, [word.word_id]);

  return (
    <View className="word-card-large">
      <Text className="word-text-large">{word.word_name}</Text>
      <View className="pronunciation-group-container">
        <WordPronounceComp word={word} />
      </View>
      <View className="content-area">
        <Text className="explain">释义：{word?.explains}</Text>
        {sentences?.length ? (
          <View className="word-large-card-description">
            示例：
            {sentences.map((item, index) => {
              return (
                <Text className="word-large-card-description-content" key={item.text}>
                  {item.text}
                </Text>
              );
            })}
          </View>
        ): null}
      </View>

      <View style={{ marginTop: "16px" }}></View>
      <View className="word-large-card-learning-target">
          学习目标：
      </View>
      <TagSelection
        initialTagId={word.tag_id}
        onTagChange={(tag) => {
          TagsService.getInstance().updateWordTag(word.word_id, tag.tag_id);
        }}
      />
    </View>
  );
};

export default WordCard;
