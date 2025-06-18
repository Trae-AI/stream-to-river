// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { View, Text, Image } from "@tarojs/components";
import React, { useState, useEffect, useRef } from "react";
import { Word } from "@/type/word";
import { Tag } from "@/type/tag";
import AudioManager from "@/libs/audio-manager";

import SpeakerSVG from "@/assets/icons/speaker.svg";

import "./index.less";

interface MiniWordCardProps {
  word: Word;
  tags: Map<number, Tag>;
  onClick: () => void;
}

const generateProgressPath = (
  progress: number,
  centerX: number = 9,
  centerY: number = 9.5,
  radius: number = 8
): string => {
  // Ensure progress is between 0 and 1
  const normalizedProgress = Math.max(0, Math.min(1, progress));

  // If progress is 0, return empty path
  if (normalizedProgress === 0) {
    return "";
  }

  // If progress is 1, return complete circle
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

const MiniWordCard: React.FC<MiniWordCardProps> = ({ word, tags, onClick }) => {
  const [isPlaying, setIsPlaying] = useState(false);
  const audioManager = useRef(AudioManager.getInstance());
  const [progress, setProgress] = useState<number>(0);
  const [progressPath, setProgressPath] = useState<string>("");

  useEffect(() => {
    // Register playback status callback
    audioManager.current.registerCallback(word.word_id, setIsPlaying);

    // Cleanup function
    return () => {
      audioManager.current.unregisterCallback(word.word_id);
    };
  }, [word.word_id]);

  const handlePlaySound = async (pronounceUrl: string) => {
    // Check if audio URL is valid
    if (!pronounceUrl || pronounceUrl.trim() === "") {
      console.warn("Audio URL is empty");
      return;
    }

    try {
      // If currently playing, stop
      if (isPlaying) {
        // Can add stop logic here, but usually let AudioManager handle it
        return;
      }

      // Play audio
      await audioManager.current.playAudio(word.word_id, pronounceUrl);
    } catch (error) {
      console.error("Failed to play audio:", error);
    }
  };

  // Get tag text
  const getTagText = (tagId: number) => {
    const tag = tags.get(tagId);
    return tag ? tag.tag_name : "积流成江";
  };

  useEffect(() => {
    // Register playback status callback
    const p = word.level / word.max_level;
    setProgress(p);
    setProgressPath(generateProgressPath(p));
  }, [word]);

  return (
    <View className="word-card">
      <View className="word-card-content">
        <Image
          src={SpeakerSVG}
          className={`word-card-speaker ${isPlaying ? "playing" : ""}`}
          onClick={() => handlePlaySound(word.pronounce_us.url)}
        />
        <View className="word-card-phonetic-container" onClick={onClick}>
          <View className="word-card-tag-container">
            <Text className="word-card-word">{word.word_name}</Text>
            <View className="word-card-tag-el">
              <Text className="word-card-tag-text">
                {getTagText(word.tag_id)}
              </Text>
            </View>
          </View>
          <Text className="word-card-explain">{word.explains}</Text>
        </View>
      </View>
    </View>
  );
};

export default MiniWordCard;
