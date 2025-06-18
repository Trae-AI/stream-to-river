// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { View, ScrollView, Text } from "@tarojs/components";
import React, {
  useState,
  useEffect,
  useCallback,
  useRef,
  forwardRef,
  useImperativeHandle,
} from "react";
import Taro from "@tarojs/taro";

import { Word } from "@/type/word";
import { WordService } from "@/service/word-service";

import { Tag } from "@/type/tag";

import MiniWordCard from "@/components/mini-word-card";
import useWordListStore from "@/store/word-list";
import useSystemTagsStore from "@/store/tag";

import "./index.less";

export interface WordListProps {
  isEnableScroll?: boolean;
}

export interface WordListRef {
  refresh: () => Promise<void>;
}

const WordList = forwardRef<WordListRef, WordListProps>(
  ({ isEnableScroll = true }, ref) => {
    const [words, setWords] = useState<Word[]>([]);
    const [loading, setLoading] = useState(false);
    const [hasMore, setHasMore] = useState(true);
    const [offset, setOffset] = useState<number>(0);
    const { tags } = useSystemTagsStore();
    const [tagMap, setTagMap] = useState<Map<number, Tag>>(new Map());

    const divRef = useRef<HTMLDivElement>(null);
    const obsRef = useRef<any>(null);

    const { setWordList, setInitialIndex } = useWordListStore();

    // Load word list
    const loadWords = useCallback(
      async (isRefresh = false, customOffset?: number) => {
        if (loading) return;

        setLoading(true);

        try {
          // Use custom offset or current offset
          const currentOffset =
            customOffset !== undefined ? customOffset : offset;
          const response = await WordService.getInstance().getWordList(
            currentOffset
          );

          if (isRefresh) {
            setWords(response.words_list);
            setOffset(response.words_list.length);
          } else {
            setWords((prev) => [...prev, ...response.words_list]);
            setOffset((prev) => prev + response.words_list.length);
          }

          // Check if there is more data
          setHasMore(response.words_list.length > 0);

          if (response.words_list.length === 0) {
            obsRef?.current?.disconnect?.();
            obsRef.current = null;
          }
        } catch (error) {
          console.error("Failed to load word list:", error);
          Taro.showToast({
            title: "加载失败，请重试",
            icon: "none",
          });
        } finally {
          setLoading(false);
        }
      },
      [loading, offset]
    );

    // Expose refresh method to external components
    const refresh = useCallback(async () => {
      // Reset offset to 0
      setOffset(0);
      setHasMore(true);
      // Reload from offset=0
      await loadWords(true, 0);
      // Reload tags
      // await loadTags();
    }, [words, loadWords]);

    // Use useImperativeHandle to expose methods to parent component
    useImperativeHandle(
      ref,
      () => ({
        refresh,
      }),
      [refresh]
    );

    useEffect(() => {
      if (hasMore) {
        const observer = new IntersectionObserver((entries) => {
          entries.forEach((entry) => {
            // When the observed element enters viewport
            if (entry.isIntersecting) {
              loadWords(); // Function to load more data
            }
          });
        });
        obsRef.current = observer;
        observer.observe(divRef.current as Element);
      }
      return () => {
        obsRef?.current?.disconnect?.();
      };
    }, [hasMore, loadWords]);

    useEffect(() => {
      const tMap = new Map<number, Tag>();
      tags.forEach((tag) => {
        tMap.set(tag.tag_id, tag);
      });
      setTagMap(tMap);
    }, [tags]);

    const handleWordRowClick = (index: number) => {
      setWordList(words);
      setInitialIndex(index);
      Taro.navigateTo({ url: "/pages/words_study/index" });
    };

    return (
      <ScrollView
        className="word-list-scroll"
        scrollY={isEnableScroll}
        refresherEnabled
        refresherTriggered={loading}
        // onScrollToLower={onScrollToLower}
        // onScroll={onScroll}
        // lowerThreshold={200}
        scrollWithAnimation
      >
        {words.map((word, index) => (
          <View key={`${word.word_id}-${index}`}>
            <MiniWordCard
              word={word}
              tags={tagMap}
              onClick={() => handleWordRowClick(index)}
            />
          </View>
        ))}

        <div
          ref={divRef}
          id="obs-element"
          style={{ width: 100, height: 1, background: "transparent" }}
        />

        {loading && (
          <View className="foot-container">
            <View className="foot-hint">加载中...</View>
          </View>
        )}

        {!hasMore && words.length > 0 && (
          <View className="foot-container">
            <Text className="foot-hint">已加载全部单词</Text>
          </View>
        )}

        {words.length === 0 && !loading && (
          <View className="foot-container">
            <Text className="foot-hint">暂无单词数据</Text>
          </View>
        )}
      </ScrollView>
    );
  }
);

WordList.displayName = "WordList";

export default WordList;
