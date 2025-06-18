// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { View, Text, Image } from '@tarojs/components';
import React, { useState, useEffect, useCallback, useRef } from 'react';
import { Word, WordPronounce } from '@/type/word';
import SpeakerIcon from '@/assets/icons/speaker.svg'
import AudioManager from '@/libs/audio-manager';

import './index.less'

interface WordPronounceProps {
  word: Word;
}

export enum PronounceType {
  US = 'us',
  UK = 'uk',
}

interface PronounceProps {
  pronounce: WordPronounce;
  type: PronounceType;
  playCallback?: (url: string) => void;
}

const Pronounce: React.FC<PronounceProps> = ({
  pronounce,
  type,
  playCallback,
}) => {
  const phonetic = pronounce.phonetic && pronounce.phonetic.length > 0 ? `/${pronounce.phonetic}/` : '/暂无音标/';

  return (
    <View className='pronounce-container' onClick={() => { playCallback?.(pronounce.url) }}>
      {/* <Text className='pronounce-text'>{type === PronounceType.UK ? '英' : '美'}</Text> */}
      <Text className='pronounce-text'>{phonetic}</Text>
      <Image className='pronounce-icon' src={SpeakerIcon} onClick={() => { playCallback?.(pronounce.url) }} />
    </View>
  )
}

const WordPronounceComp: React.FC<WordPronounceProps> = ({
  word,
}) => {
  const [isPlaying, setIsPlaying] = useState(false);
  const audioManager = useRef(AudioManager.getInstance());


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
    if (!pronounceUrl || pronounceUrl.trim() === '') {
      console.warn('Audio URL is empty');
      return;
    }

    try {
      // If currently playing, then stop
      if (isPlaying) {
        // Can add stop logic here, but usually let AudioManager handle it
        return;
      }

      // Play audio
      await audioManager.current.playAudio(word.word_id, pronounceUrl);
    } catch (error) {
      console.error('Failed to play audio:', error);
    }
  };

  return (
    <View className='word-pronounce-container'>
      <Pronounce pronounce={word.pronounce_us} type={PronounceType.US} playCallback={handlePlaySound} />
      {/* <Pronounce pronounce={word.pronounce_uk} type={PronounceType.UK} playCallback={handlePlaySound} /> */}
    </View>
  )
}


export default WordPronounceComp;
