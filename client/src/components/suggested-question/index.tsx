// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { Text, View } from "@tarojs/components";
import "./index.less";

type Question = {
  text: string;
  id: string;
};

type Props = {
  questions: Question[];
  send: (text: string) => void;
};

const SuggestedQuestions = (props: Props) => {

  const handleSend = (question: Question) => {
    props.send(question.text);
  }

  return (
    <View className='suggested-container'>
      <Text className='suggested-tips'>你是否有这些英语问题：</Text>
      {props.questions.map((question) => (
        <View key={question.id} className='suggested-question-container' onClick={() => handleSend(question)}>
          <Text className='suggested-question-text'>{question.text}</Text>
        </View>
      ))}
    </View>
  );
};

export default SuggestedQuestions;
