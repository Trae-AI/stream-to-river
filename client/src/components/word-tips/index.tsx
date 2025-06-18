// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { Text, View } from "@tarojs/components";
import { useState, useEffect, useMemo } from "react";
import { WordService, NewWordResp } from "@/service/word-service";
import useUserStore from "@/store/user";
import Taro from "@tarojs/taro";
import TagSelection from "@/components/tag-selection";
import { extractStyleExamples, removeWhitespace } from "./utils";
import { Word } from "@/type/word";
import WordPronounceComp from "@/components/word-pronunce";
import { createPortal } from "react-dom";

import "./index.less";
import { STORAGE_WORD_KEY } from "@/consts";

function convertWord(word: NewWordResp): Word {
  return {
    word_id: 0,
    word_name: word.new_word_name,
    description: word.description,
    explains: word.explains,
    pronounce_uk: word.pronounce_uk,
    pronounce_us: word.pronounce_us,
    tag_id: 0,
    level: 0,
    max_level: 0,
    sentences: [],
  };
}

const Modal = ({
  visible,
  onClose,
  word,
  onOK,
}: {
  visible: boolean;
  onClose: () => void;
  word: string;
  onOK: () => void;
}) => {
  const [animationClass, setAnimationClass] = useState("");
  const { userInfo } = useUserStore();
  const [selectedTagId, setSelectedTagId] = useState<number>(-1);
  const [newWord, setNewWord] = useState<NewWordResp>();

  useEffect(() => {
    if (!word || word.trim().length == 0) {
      return;
    }
    const cleanWord = removeWhitespace(word);
    WordService.getInstance()
      .queryWordByName(cleanWord)
      .then((res) => {
        setNewWord(res);
      });
  }, [word]);

  const exampleList = useMemo(() => {
    return extractStyleExamples(newWord?.description || "");
  }, [newWord?.description]);

  const addWord = async (word: string, tag_id: number) => {
    if (tag_id < 0) {
      Taro.showToast({
        title: "请选择标签",
        icon: "none",
      });
      return;
    }

    const cleanWord = removeWhitespace(word);

    try {
      const res = await WordService.getInstance().addWord(cleanWord, tag_id);
      if (res?.BaseResp?.StatusCode === 0) {
        Taro.showToast({
          title: "添加成功",
          icon: "success",
        });
        onOK && onOK();
      } else if (res?.BaseResp?.StatusCode === 1) {
        Taro.showToast({
          title: "单词已被添加过",
          icon: "none",
        });
        onOK && onOK();
      } else {
        Taro.showToast({
          title: res?.error_msg || "添加失败",
          icon: "none",
        });
      }
    } catch (error: unknown) {
      const respError = error as Response;
      if (respError.status === 401) {
        onClose && onClose();
        Taro.showToast({
          title: "请先登录",
          icon: "none",
        });
        Taro.navigateTo({
          url: "/pages/login/index",
        });
        Taro.setStorageSync(STORAGE_WORD_KEY, word);
      } else {
        Taro.showToast({
          title: "添加失败",
          icon: "none",
        });
      }
    }
  };

  useEffect(() => {
    if (visible) {
      setAnimationClass("modal-show");
    } else {
      setAnimationClass("modal-hide");
    }
  }, [visible]);

  const handleMaskClick = () => {
    onClose && onClose();
  };

  return createPortal(
    visible ? (
      <View className={`modal-container ${animationClass}`}>
        <View className="modal-mask" onClick={handleMaskClick} />
        <View className="modal-content" onClick={(e) => e.stopPropagation()}>
          <Text className="word-text">{word}</Text>
          <View className="pronunciation-group-container">
            {newWord && <WordPronounceComp word={convertWord(newWord)} />}
          </View>
          <Text className="explain">{newWord?.explains}</Text>
          <View className="description">
            {exampleList.map((item, index) => {
              return (
                <Text className="description-content" key={item.english}>
                  {item.english}（{item.chinese}）
                </Text>
              );
            })}
          </View>
          <View className="learning-target-txt">
            学习目标：
          </View>
          <TagSelection onTagChange={(tag) => setSelectedTagId(tag.tag_id)} />
          <View style={{ marginTop: "20px" }}></View>
          <View className="button-group">
            <View className="button cancel-button" onClick={handleMaskClick}>
              取消
            </View>
            <View
              className="button add-button"
              onClick={() => addWord(word, selectedTagId)}
            >
              添加至学习
            </View>
          </View>
        </View>
      </View>
    ) : null,
    document.body
  );
};

export default Modal;
