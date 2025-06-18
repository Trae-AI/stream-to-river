// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import {
  useState,
  useEffect,
  useRef,
  forwardRef,
  useImperativeHandle,
} from "react";
import { View, Text, Input } from "@tarojs/components";
import "./index.less";

interface TestFillingProps {
  question: string;
  fillingBlanks: Array<string>;
  onComplete?: (userInput: string) => void;
  isCorrect?: boolean;
  correctWord?: string;
  visible: boolean;
}

const TestFilling = forwardRef<
  { reviewAnswer: (word: string) => void },
  TestFillingProps
>(
  (
    { question, fillingBlanks, onComplete, isCorrect, correctWord, visible },
    ref
  ) => {
    const [inputs, setInputs] = useState<Array<string>>(fillingBlanks);
    const inputRefs = useRef<any[]>([]);

    useEffect(() => {
      setInputs(fillingBlanks);
    }, [fillingBlanks]);

    useImperativeHandle(ref, () => ({
      reviewAnswer: (word: string) => {
        const wordArray = word.split("");
        const newInputs = [...fillingBlanks];
        wordArray.forEach((char, index) => {
          if (index < newInputs.length) {
            newInputs[index] = char;
          }
        });
        setInputs(newInputs);
      },
    }));

    // Focus on the first input box during initialization
    useEffect(() => {
      if (isCorrect === undefined && visible) {
        // Find the index of the first empty input box or input box with value "_"
        const firstEmptyIndex = inputs.findIndex(
          (input) => input.trim() === "" || input.trim() === "_"
        );

        if (firstEmptyIndex !== -1) {
          const targetInput = inputRefs.current[firstEmptyIndex];
          if (targetInput) { // Ensure the ref is populated
            setTimeout(() => {
              targetInput.focus();
            }, 100); // Delay might still be needed for some environments
          }
        }
      }
    }, [inputs, isCorrect, visible]); // Keep dependencies as is

    // Check if all input boxes are filled
    useEffect(() => {
      const allFilled = inputs.every(
        (input) => input.trim() !== "_" && input.trim() !== ""
      );
      if (allFilled && isCorrect === undefined) {
        const userAnswer = inputs.join("");
        onComplete?.(userAnswer);
      }
    }, [inputs, isCorrect, onComplete]);

    const handleInputChange = (index: number, value: string) => {
      // Only allow single letter input
      const letter = value.slice(-1).toLowerCase();
      if (letter && /^[a-zA-Z-]$/i.test(letter)) {
        const newInputs = [...inputs];
        newInputs[index] = letter;
        setInputs(newInputs);

        // Automatically jump to the next input box
        if (index < newInputs.length - 1) {
          let nextFocusTarget = -1;
          for (let i = index + 1; i < newInputs.length; i++) {
            if (newInputs[i].trim() === "" || newInputs[i].trim() === "_") {
              nextFocusTarget = i;
              break;
            }
          }
          if (nextFocusTarget !== -1) {
            setTimeout(() => {
              inputRefs.current[nextFocusTarget]?.focus();
            }, 50);
          }
        }
      } else if (value === "") {
        // Allow deletion
        const newInputs = [...inputs];
        newInputs[index] = "";
        setInputs(newInputs);
      }
    };

    const handleKeyDown = (index: number, e: any) => {
      // Handle backspace key
      if (e.detail.keyCode === 8 && inputs[index] === "" && index > 0) {
        // If current box is empty and backspace is pressed, move to previous box
        setTimeout(() => {
          inputRefs.current[index - 1]?.focus();
        }, 50);
      }
    };

    const getInputClassName = (index: number) => {
      if (isCorrect === false && inputs[index] !== correctWord?.[index]) {
        return "test-filling__input incorrect";
      } else {
        return "test-filling__input";
      }
    };

    const firstEmptyIndex = inputs.findIndex(
      (input) => input.trim() === "" || input.trim() === "_"
    );


    return (
      <View className="test-filling">
        <View className="test-filling__question">
          <Text className="test-filling__question-text">{question}</Text>
        </View>

        <View className="test-filling__hint">
          <Text className="test-filling__hint-text">填入单词的正确拼写</Text>
        </View>

        <View className="test-filling__inputs">
          {inputs.map((input, index) => (
            <View key={index} className="test-filling__input-wrapper">
              <Input
                ref={(el) => (inputRefs.current[index] = el)}
                className={getInputClassName(index)}
                value={input}
                maxlength={1}
                onInput={(e) => handleInputChange(index, e.detail.value)}
                //@ts-ignore
                onKeyDown={(e) => handleKeyDown(index, e)}
                disabled={isCorrect !== undefined}
                placeholder=""
                autoFocus={isCorrect === undefined && index === firstEmptyIndex}
              />
            </View>
          ))}
        </View>
      </View>
    );
  }
);

export default TestFilling;
