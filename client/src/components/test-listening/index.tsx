// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import {
  useState,
  useEffect,
  useRef,
  forwardRef,
  useImperativeHandle,
} from "react";
import { View, Text, Image } from "@tarojs/components";
import Taro from "@tarojs/taro";
import { ReactComponent as Bg } from "@/assets/icons/bg.svg";
import { ReactComponent as Wave0 } from "@/assets/icons/wave0.svg";
import { ReactComponent as Wave1 } from "@/assets/icons/wave1.svg";
import { ReactComponent as Wave2 } from "@/assets/icons/wave2.svg";
import "./index.less";

interface Choice {
  id: string;
  text: string;
}

interface TestListeningProps {
  // Audio file path
  audioUrl: string;
  // Question prompt text
  question?: string;
  // Choice list
  choices: Choice[];
  // Selection callback
  onSelect?: (choiceId: string) => void;
  // Whether a choice has been made
  correctChoiceId: number;
  visible?: boolean;
  // // Whether to show result
  // showResult?: boolean
}

const TestListening = forwardRef<
  { reviewAnswer: (correctAnswerId: number) => void },
  TestListeningProps
>(
  (
    {
      audioUrl,
      question,
      choices,
      onSelect,
      correctChoiceId,
      visible,
      // showResult = false
    },
    ref
  ) => {
    const [selectedId, setSelectedId] = useState<string>("");
    const [isPlaying, setIsPlaying] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const audioContextRef = useRef<any>(null);
    const innerAudioContextRef = useRef<any>(null);
    const playPromiseRef = useRef<Promise<void> | null>(null);
    const isComponentMountedRef = useRef(true);
    const [currentWave, setCurrentWave] = useState(2);
    const waveAnimationRef = useRef<NodeJS.Timeout | null>(null);
    const initialLoadCompleteRef = useRef(false);

    useImperativeHandle(ref, () => ({
      reviewAnswer: (correctAnswerId: number) => {
        setSelectedId(correctAnswerId.toString());
      },
    }));

    const waves = [Wave0, Wave1, Wave2];

    const animateWaves = () => {
      if (!isPlaying) return;
      waveAnimationRef.current = setTimeout(() => {
        setCurrentWave((prev) => (prev + 1) % 3);
        animateWaves();
      }, 400); // This time can be adjusted to control animation speed
    };

    useEffect(() => {
      if (isPlaying) {
        setCurrentWave(2); // Reset to wave form 2 when starting playback
        animateWaves();
      } else {
        if (waveAnimationRef.current) {
          clearTimeout(waveAnimationRef.current);
        }
        setCurrentWave(2); // Reset to wave form 2 when stopping playback
      }

      return () => {
        if (waveAnimationRef.current) {
          clearTimeout(waveAnimationRef.current);
        }
      };
    }, [isPlaying]);
    const playAudio = async () => {
      // Only show loading on initial load
      if (!initialLoadCompleteRef.current) {
        setIsLoading(true);
      }

      try {
        if (audioContextRef.current) {
          // Wait for previous playback Promise to complete
          if (playPromiseRef.current) {
            try {
              await playPromiseRef.current;
            } catch (error) {
              // Ignore interrupted playback errors
            }
          }

          // Pause current playback first
          if (!audioContextRef.current.paused) {
            audioContextRef.current.pause();
            await new Promise((resolve) => setTimeout(resolve, 100));
          }

          if (isComponentMountedRef.current) {
            // Reset playback position
            audioContextRef.current.currentTime = 0;

            // Start new playback
            playPromiseRef.current = audioContextRef.current.play();
            await playPromiseRef.current;
          }
        }
      } catch (error) {
        console.warn("Error playing audio:", error);
        if (isComponentMountedRef.current) {
          setIsLoading(false);
          setIsPlaying(false);

          // Only show error prompt for non-abort errors
          if (error.name !== "AbortError") {
            Taro.showToast({
              title: "请点击播放按钮",
              icon: "none",
            });
          }
        }
      }
    };
    useEffect(() => {
      isComponentMountedRef.current = true;
      // H5 environment
      audioContextRef.current = new Audio(audioUrl);

      audioContextRef.current.onloadstart = () => {
        if (isComponentMountedRef.current && !initialLoadCompleteRef.current) {
          setIsLoading(true);
        }
      };

      audioContextRef.current.oncanplay = () => {
        if (isComponentMountedRef.current) {
          setIsLoading(false);
          initialLoadCompleteRef.current = true; // Mark initial load as complete
        }
      };

      audioContextRef.current.onplay = () => {
        if (isComponentMountedRef.current) {
          setIsPlaying(true);
          setIsLoading(false);
        }
      };

      audioContextRef.current.onended = () => {
        if (isComponentMountedRef.current) {
          setIsPlaying(false);
          setIsLoading(false);
        }
      };

      audioContextRef.current.onpause = () => {
        if (isComponentMountedRef.current) {
          setIsPlaying(false);
          setIsLoading(false);
        }
      };

      audioContextRef.current.onerror = (err) => {
        console.error("Audio playback error:", err);
        if (isComponentMountedRef.current) {
          setIsPlaying(false);
          setIsLoading(false);
        }
      };
      if (visible) {
        setTimeout(() => {
          playAudio();
        }, 300);
      }

      return () => {
        isComponentMountedRef.current = false;
        initialLoadCompleteRef.current = false;

        // Clean up audio resources
        if (innerAudioContextRef.current) {
          try {
            innerAudioContextRef.current.stop();
            innerAudioContextRef.current.destroy();
          } catch (error) {
            console.warn("Error cleaning up mini program audio resources:", error);
          }
        }

        if (audioContextRef.current) {
          try {
            audioContextRef.current.pause();
            audioContextRef.current.currentTime = 0;
            audioContextRef.current.src = "";
          } catch (error) {
            console.warn("Error cleaning up H5 audio resources:", error);
          }
        }

        // Wait for playback Promise to complete
        if (playPromiseRef.current) {
          playPromiseRef.current.catch(() => {
            // Ignore interrupted playback errors
          });
        }
      };
    }, [audioUrl, visible]);

    const handleChoiceSelect = (choice: Choice) => {
      if (selectedId) return; // No further selection allowed after choice is made and result is shown
      onSelect?.(choice.id);
      setSelectedId(choice.id);
    };

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
    const WaveComponent = waves[currentWave];
    return (
      <View className="test-listening">
        {/* Audio playback area */}
        <View className="test-listening__question">
          <Text className="test-listening__question-text">{question}</Text>
        </View>

        {/* Audio play button and wave effect */}
        <View className="test-listening__play-button" onClick={playAudio}>
          <Bg className="test-listening__bg"></Bg>
          {isLoading ? (
            <View className="test-listening__loading-spinner"></View>
          ) : (
            <WaveComponent className="test-listening__wave" />
          )}
        </View>

        <Text className="test-listening__tips">
          请选择你听到的单词的中文释义
        </Text>

        {/* Options area */}
        <View className="test-listening__options">
          {choices
            .filter((item) => Boolean(item.text))
            .sort((a, b) => a.id.localeCompare(b.id))
            .map((choice) => (
              <View
                key={choice.id}
                className={getChoiceClassName(choice)}
                onClick={() => handleChoiceSelect(choice)}
              >
                <Text className="test-listening__option-text">
                  {choice.text}
                </Text>
              </View>
            ))}
        </View>
      </View>
    );
  }
);

export default TestListening;
