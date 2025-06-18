// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { FC, useState, useMemo, useEffect, useRef } from "react";
import Taro from "@tarojs/taro";
import { View, Text } from "@tarojs/components";
import TestChoices from "@/components/test-choices";
import TestFilling from "@/components/test-filling";
import TestListening from "@/components/test-listening";
import ReviewPuzzleService, {
  ReviewQuestion,
  OptionItem,
  QuestionType,
  AnswerPuzzleReq,
} from "@/service/review-service";
import WordService from "@/service/word-service";

import "./index.less";
import { ReactComponent as GreatIcon } from '@/assets/images/great.svg';
import { ReactComponent as SmileIcon } from '@/assets/images/smile.svg';
import { Navbar } from "./navbar";

const VocabularyReviewComponent: FC<{
  visible: boolean;
  onClose: () => void;
}> = ({ visible, onClose }) => {
  const [puzzles, setPuzzles] = useState<ReviewQuestion[]>([]);
  const [currentIndex, setCurrentIndex] = useState(-1);
  const [loading, setLoading] = useState(true);
  const [currentCorrectAnswerId, setCurrentCorrectAnswerId] =
    useState<number>(-1);
  // undefined means unanswered, true means correct answer, false means wrong answer
  const [currentAnswerCorrect, setCurrentAnswerCorrect] = useState<
    boolean | undefined
  >(undefined);
  const [currentWord, setCurrentWord] = useState<string>("");
  const currentPuzzle = puzzles[currentIndex];

  // Child component refs
  const testChoicesRef = useRef<{
    reviewAnswer: (correctAnswerId: number) => void;
  }>(null);
  const testFillingRef = useRef<{ reviewAnswer: (word: string) => void }>(null);
  const testListeningRef = useRef<{
    reviewAnswer: (correctAnswerId: number) => void;
  }>(null);

  const handleCheckAnswer = async () => {
    let correctAnswerId = currentCorrectAnswerId;
    if (correctAnswerId === -1) {
      const puzzleReq = {
        answer_id: -1,
        word_id: currentPuzzle.word_id.toString(),
        question_type: currentPuzzle.question_type,
        filled_name: "",
      };
      const res = await ReviewPuzzleService.getInstance().submitPuzzleAnswer(
        puzzleReq as any
      );
      correctAnswerId = res.correct_answer_id;
      setCurrentCorrectAnswerId(correctAnswerId);
    }

    setCurrentAnswerCorrect(true);

    // Call child component's reviewAnswer method to set correct answer
    if (currentPuzzle.question_type === QuestionType.FILL_IN_BLANK) {
      // Fill-in-blank: Set correct word
      testFillingRef.current?.reviewAnswer(currentWord);
    } else if (
      currentPuzzle.question_type === QuestionType.CHOOSE_CN ||
      currentPuzzle.question_type === QuestionType.CHOOSE_EN
    ) {
      testChoicesRef.current?.reviewAnswer(correctAnswerId);
    } else if (currentPuzzle.question_type === QuestionType.PRONOUNCE_CHOOSE) {
      testListeningRef.current?.reviewAnswer(correctAnswerId);
    }
  };

  useEffect(() => {
    if (currentIndex >= 0 && currentIndex < puzzles.length) {
      const currentPuzzle = puzzles[currentIndex];
      if (
        currentPuzzle &&
        (currentPuzzle.question_type === QuestionType.FILL_IN_BLANK ||
          currentPuzzle.question_type === QuestionType.PRONOUNCE_CHOOSE)
      ) {
        WordService.getInstance()
          .queryWordByID(currentPuzzle.word_id)
          .then((resp) => {
            setCurrentWord(resp.word_name);
          });
      }
    }
  }, [currentIndex, puzzles]);

  const cleanUp = () => {
    setPuzzles([]);
    setCurrentIndex(-1);
    setLoading(true);
    setCurrentCorrectAnswerId(-1);
    setCurrentAnswerCorrect(undefined);
    setCurrentWord("");
  };

  useEffect(() => {
    (async () => {
      try {
        setLoading(true);
        const response =
          await ReviewPuzzleService.getInstance().getReviewPuzzles();
        setPuzzles(response.questions);

        if (response.questions.length > 0) {
          setCurrentIndex(0);
        }
      } catch (error) {
        console.error("获取题目失败:", error);
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  const handleAnswer = (answer: any) => {
    if (!currentPuzzle) return;

    var puzzleReq = {};

    if (
      currentPuzzle.question_type === QuestionType.CHOOSE_CN ||
      currentPuzzle.question_type === QuestionType.CHOOSE_EN ||
      currentPuzzle.question_type === QuestionType.PRONOUNCE_CHOOSE
    ) {
      puzzleReq = {
        answer_id: answer,
        word_id: currentPuzzle.word_id.toString(),
        question_type: currentPuzzle.question_type,
        filled_name: "",
      };
    } else {
      puzzleReq = {
        answer_id: -1,
        word_id: currentPuzzle.word_id.toString(),
        question_type: currentPuzzle.question_type,
        filled_name: answer,
      };
    }

    ReviewPuzzleService.getInstance()
      .submitPuzzleAnswer(puzzleReq as AnswerPuzzleReq)
      .then((response) => {
        setCurrentCorrectAnswerId(response.correct_answer_id);
        setCurrentAnswerCorrect(response.is_correct);
        autoNext(response.is_correct);
      })
      .catch((error) => {
        console.error("提交答案失败:", error);
      });
  };

  const lockNextRef = useRef<boolean>(false);

  const autoNext = (is_correct: boolean) => {
    if (is_correct) {
      lockNextRef.current = true;
      setTimeout(() => {
        lockNextRef.current = false;
        handleNext();
      }, 800);
    }
  };

  const handleNext = () => {
    if (lockNextRef.current) return;

    if (currentIndex < puzzles.length - 1) {
      setCurrentIndex((prev) => prev + 1);
      setCurrentCorrectAnswerId(-1);
      setCurrentAnswerCorrect(undefined);
    } else {
      Taro.showToast({
        title: "已完成所有题目",
        icon: "success",
      });
    }
  };

  const convertOptionsToChoices = (options: OptionItem[]) => {
    return options.map((option) => ({
      id: option.answer_list_id.toString(),
      text: option.description,
    }));
  };

  const currentFillingBlanks = useMemo(() => {
    const currentPuzzle = puzzles[currentIndex];
    if (
      currentPuzzle &&
      currentPuzzle.question_type === QuestionType.FILL_IN_BLANK &&
      "show_info" in currentPuzzle
    ) {
      return currentPuzzle.show_info.map((character) =>
        character === "_" ? "" : character
      );
    }
    return [];
  }, [puzzles, currentIndex]);

  const renderPuzzle = () => {
    const currentPuzzle = puzzles[currentIndex];
    if (!currentPuzzle) return null;

    switch (currentPuzzle.question_type) {
      case QuestionType.CHOOSE_CN:
      case QuestionType.CHOOSE_EN:
        return (
          <TestChoices
            ref={testChoicesRef}
            type={
              currentPuzzle.question_type
                ? "word-to-meaning"
                : "meaning-to-word"
            }
            question={currentPuzzle.question}
            choices={convertOptionsToChoices(currentPuzzle.options)}
            onSelect={(choiceId) => handleAnswer(choiceId)}
            correctChoiceId={currentCorrectAnswerId}
            key={currentPuzzle.question + currentPuzzle.question_type}
          />
        );

      case QuestionType.FILL_IN_BLANK:
        if (
          !("show_info" in currentPuzzle) ||
          currentPuzzle.show_info.length === 0
        )
          if (currentIndex < puzzles.length - 1) {
            setCurrentIndex((prev) => prev + 1);
            return <View></View>;
          } else {
            return (
              <View className="vocabulary-review__empty">
                <Text>暂无题目</Text>
              </View>
            );
          }
        return (
          <TestFilling
            visible={visible}
            ref={testFillingRef}
            question={currentPuzzle.question}
            fillingBlanks={currentFillingBlanks}
            onComplete={(userInput) => handleAnswer(userInput)}
            isCorrect={currentAnswerCorrect}
            correctWord={currentWord}
            key={currentPuzzle.word_id + currentPuzzle.question_type}
          />
        );

      case QuestionType.PRONOUNCE_CHOOSE:
        return (
          <TestListening
            visible={visible}
            ref={testListeningRef}
            audioUrl={currentPuzzle.question}
            choices={convertOptionsToChoices(currentPuzzle.options)}
            onSelect={(choiceId) => handleAnswer(choiceId)}
            correctChoiceId={currentCorrectAnswerId}
            key={currentPuzzle.word_id + currentPuzzle.question_type}
          />
        );

      default:
        return <View>未知题目类型</View>;
    }
  };

  if (loading) {
    return null;
  }

  if (puzzles.length === 0) {
    return <View
      className="vocabulary-review"
      style={{ display: visible ? "block" : "none" }}
    >
      <Navbar back={onClose} />
      <View className="vocabulary-review__empty">
        <GreatIcon className="great-icon" />
        <SmileIcon className="smile-icon" />
        <Text className="vocabulary-review__empty__text">积少成多，积流成江</Text>
        <Text className="vocabulary-review__empty__text">你今天已经学完所有的单词</Text>
      </View>
    </View>
  };

  return (
    <View
      className="vocabulary-review"
      style={{ display: visible ? "block" : "none" }}
    >
      <Navbar back={onClose} />
      {/* Question content area */}
      <View className="vocabulary-review__content">{renderPuzzle()}</View>

      {/* Bottom control area */}
      <View className="vocabulary-review__footer">
        {/* Progress display in bottom left */}
        <View className="vocabulary-review__progress">
          <Text className="vocabulary-review__progress-text">
            {currentIndex + 1} / {puzzles.length}
          </Text>
        </View>

        {/* Navigation buttons in the middle */}
        {currentAnswerCorrect === undefined ||
        (currentAnswerCorrect === false &&
          currentPuzzle.question_type === QuestionType.FILL_IN_BLANK) ? (
          <View className="vocabulary-review__link">
            <Text
              className="vocabulary-review__link-text"
              onClick={handleCheckAnswer}
            >
              查看答案
            </Text>
          </View>
        ) : (
          <View className="vocabulary-review__navigation">
            {currentIndex < puzzles.length - 1 && (
              <Text
                className="vocabulary-review__nav-btn vocabulary-review__nav-btn--next"
                onClick={handleNext}
              >
                下一题
              </Text>
            )}
            {currentIndex === puzzles.length - 1 && (
              <Text
                className="vocabulary-review__nav-btn vocabulary-review__nav-btn--finish"
                onClick={() => {
                  cleanUp();
                  onClose();
                }}
              >
                完成练习
              </Text>
            )}
          </View>
        )}
      </View>
    </View>
  );
};

export default VocabularyReviewComponent;
