// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import React, { useState, useEffect, useLayoutEffect } from "react";
import { View, Text } from "@tarojs/components";

import useWordListStore from "@/store/word-list";
import { useDrag } from "@use-gesture/react";

import WordCard from "@/components/word-card";
import CustomNavbar from "@/components/custom-navbar";

import "./index.less";
import Taro from "@tarojs/taro";

const WordStudyPage: React.FC<{}> = () => {
  const { wordList, initialIndex } = useWordListStore();
  const [isAnimating, setIsAnimating] = useState(false);
  const [currentIndex, setCurrentIndex] = useState(initialIndex);
  useLayoutEffect(() => {
    const pages = Taro.getCurrentPages() || [];
    if (pages.length <= 1) {
      Taro.redirectTo({
        url: "/pages/study_main/index",
      });
    }
  }, []);
  useEffect(() => {
    if (wordList.length > 0) {
      setCurrentIndex(initialIndex);
    }
  }, [wordList, initialIndex]);

  const bind = useDrag(
    ({
      movement: [mx],
      direction: [xDir],
      velocity: [vx],
      cancel,
      canceled,
    }) => {
      // If animating, cancel the gesture
      if (isAnimating) {
        cancel();
        return;
      }

      // Set trigger threshold
      const threshold = 50;
      const velocityThreshold = 0.2;

      // Determine whether to switch based on swipe distance or velocity
      const shouldTrigger =
        Math.abs(mx) > threshold || Math.abs(vx) > velocityThreshold;

      if (shouldTrigger && !canceled) {
        setIsAnimating(true);

        if (xDir > 0 && currentIndex > 0) {
          // Swipe right, show previous card
          setCurrentIndex(currentIndex - 1);
        } else if (xDir < 0 && currentIndex < wordList.length - 1) {
          // Swipe left, show next card
          setCurrentIndex(currentIndex + 1);
        }

        setTimeout(() => {
          setIsAnimating(false);
        }, 300);

        cancel();
      }
    },
    {
      axis: "x", // Only respond to horizontal drag
      filterTaps: true, // Filter tap events
      rubberband: true, // Enable rubber band effect
      bounds: { left: -100, right: 100 }, // Set drag boundaries
    }
  );

  const getCardStyle = (index) => {
    const offset = index - currentIndex;
    const isVisible = Math.abs(offset) <= 1;

    if (!isVisible) {
      return {
        display: "none",
      };
    }

    const translateX = offset * 90;
    const scale = offset === 0 ? 1 : 0.9;
    const zIndex = offset === 0 ? 10 : 5;
    const opacity = offset === 0 ? 1 : 0.6;

    return {
      transform: `translateX(${translateX}%) scale(${scale})`,
      zIndex,
      opacity,
      transition: isAnimating ? "all 0.3s ease-out" : "none",
    };
  };

  if (!wordList || wordList.length === 0) {
    return (
      <View>
        <CustomNavbar />
        <View className="words-study">
          <View className="words-study__content">
            {/* No words */}
          </View>
        </View>
      </View>
    );
  }

  return (
    <>
      <CustomNavbar />
      <View className="words-study">
        { /* @ts-ignore */ }
        <View
          className="words-study__content"
          {...bind()} // Apply gesture binding
          style={{ touchAction: "pan-y" }} // Allow vertical scroll, disable horizontal scroll
        >
          <View className="words-study__cards-container">
            {wordList.map((word, index) => (
              <View
                key={index}
                className="words-study__card-wrapper"
                style={getCardStyle(index)}
              >
                <WordCard word={word} />
              </View>
            ))}
          </View>
        </View>
        <Text className="words-study__progress">
          {currentIndex + 1} / {wordList.length}
        </Text>
      </View>
    </>
  );
};

export default WordStudyPage;
