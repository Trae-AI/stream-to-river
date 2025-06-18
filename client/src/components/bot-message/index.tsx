// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

import { useEffect, useState } from "react";
import Modal from "../word-tips";
import { Text } from "@tarojs/components";
import "./index.less";
import { posCalcRange } from "./utils";
import Taro from "@tarojs/taro";
import useUserStore from "@/store/user";
import { STORAGE_WORD_KEY } from "@/consts";
import {Calypso} from '@bytedance/calypso';

type LinkWord = {
  start: number;
  end: number;
  text: string;
};

const BotMessage = ({
  markdownText,
  linkWordList,
  isDone,
}: {
  markdownText: string;
  linkWordList: LinkWord[];
  isGenerating: boolean;
  isDone?: boolean;
}) => {
  const { userInfo } = useUserStore();

  const [currentWords, setCurrentWords] = useState("");
  useEffect(() => {
    const u = Taro.getStorageSync(STORAGE_WORD_KEY);
    if (userInfo && userInfo?.id && u) {
      Taro.setStorageSync(STORAGE_WORD_KEY, "");
      setCurrentWords(u);
    }
  }, [userInfo]);
  const posReplaceList = linkWordList
    .map((rep) => [rep.start, rep.end])
    .flat(2)
    .filter((i) => i !== undefined)
    .map((i) => Number(i));

  const posReplaceMap = posCalcRange(markdownText, posReplaceList);

  const linkWordListNewPos = linkWordList.map((linkWord) => {
    const { start, end } = linkWord;
    const startRange = posReplaceMap[start];
    const endRange = posReplaceMap[end];

    return {
      ...linkWord,
      start: startRange,
      end: endRange,
    };
  });

  const insertedElements = linkWordListNewPos.map((linkWord) => {
    const { start, end, text } = linkWord;
    return {
      range: [start, end] as [number, number],
      render: (raw) => {
        return (
          <Text
            onClick={() => setCurrentWords(text)}
            className="dashed-underline"
          >
            {raw}
          </Text>
        );
      },
    };
  });

  return (
    <>
      <Calypso
        markDown={markdownText}
        insertedElements={insertedElements}
        style={{ fontSize: "16Px", lineHeight: "20px" }}
        showIndicator={!isDone}
        smooth={!isDone}
      />
      <Modal
        visible={Boolean(currentWords)}
        onClose={() => {
          setCurrentWords("");
        }}
        word={currentWords}
        onOK={() => {
          setCurrentWords("");
        }}
      ></Modal>
    </>
  );
};

export default BotMessage;
