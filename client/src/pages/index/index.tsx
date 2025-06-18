// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { View, ScrollView } from "@tarojs/components";
import { useMessageStore } from "@/store/chat";
import UserMessage from "@/components/user-message";
import BotMessage from "@/components/bot-message";
import SuggestedQuestions from "@/components/suggested-question";
import InputBar from "@/components/input-bar";
import WelcomeTips from "@/components/welcome-tips";
import { useEffect, useRef, useState } from "react";
import CustomNavbar from "@/components/custom-navbar";
import { v4 as uuid } from "uuid";

import { SuggestedQuestion } from "./types";
import "./index.less";
import Taro from "@tarojs/taro";
import { STORAGE_WORD_KEY } from "@/consts";

const Message = ({
  messageId,
  isGenerating,
}: {
  messageId: string;
  isGenerating: boolean;
}) => {
  const message = useMessageStore((state) => state.messageMap[messageId]);

  if (message.creator === "user") {
    return <UserMessage message={message.text} type={message.type} />;
  }

  return (
    <View className="message-bot-item">
      <BotMessage
        markdownText={message.text}
        linkWordList={message.linkWordList || []}
        isGenerating={isGenerating}
        isDone={message.isDone}
      />
    </View>
  );
};

const MessageList = ({ isGenerating }: { isGenerating: boolean }) => {
  const messageLinks = useMessageStore((state) => state.messageLinks);

  return (
    <>
      {messageLinks.map((messageId, index) => (
        <Message
          key={messageId}
          messageId={messageId}
          isGenerating={isGenerating && index === messageLinks.length - 1}
        />
      ))}
    </>
  );
};

export default function Index() {
  const {
    isGenerating,
    sendMessage,
    stopGenerating,
    sendImage,
    setConversationId,
  } = useMessageStore();
  const messageLinks = useMessageStore((state) => state.messageLinks);
  useEffect(() => {
    Taro.setStorageSync(STORAGE_WORD_KEY, "");
  }, []);
  useEffect(() => {
    const cid = uuid();
    setConversationId(cid);
  }, []);

  const needScrollToBottomRef = useRef<boolean>(true);

  const [suggestedQuestionList] = useState<SuggestedQuestion[]>([
    {
      id: "1",
      text: "发布会用英文怎么说",
    },
    {
      id: "2",
      text: "participate 和 anticipate 有什么关系",
    },
    {
      id: "3",
      text: "用英文探讨：如何用AI提升研发效率？",
    },
  ]);
  const [clientHeight, setClientHeight] = useState<number>(0);
  const pauseScrollRef = useRef<boolean>(false);
  useEffect(() => {
    Taro.createSelectorQuery()
      .select(".chat-area")
      .boundingClientRect((rect) => {
        const clientHeight = (rect as Taro.NodesRef.BoundingClientRectCallbackResult).height;
        setClientHeight(clientHeight);
      })
      .exec();
  }, []);

  const goToBottom = () => {
    Taro.createSelectorQuery()
      .select(".chat-area")
      .fields({
        node: true,
      })
      .exec((res) => {
        const scrollView = res?.[0];
        if (scrollView?.node) {
          scrollView.node.scrollTop = scrollView.node.scrollHeight;
        }
      });
  };

  useEffect(() => {
    if (!isGenerating) {
      return;
    }
    const next = () => {
      if (pauseScrollRef.current) {
        return;
      }
      if (needScrollToBottomRef.current) {
        goToBottom();
      }
    };
    const timer = setInterval(next, 100);
    return () => {
      clearInterval(timer);
    };
  }, [isGenerating]);
  useEffect(() => {
    goToBottom();
  }, []);
  const send = (text: string) => {
    sendMessage(text);
    goToBottom();
  };
  return (
    <View className="index-inner">
      <CustomNavbar />
      <ScrollView
        scrollY={true}
        onTouchStart={() => {
          pauseScrollRef.current = true;
        }}
        onTouchEnd={() => {
          pauseScrollRef.current = false;
        }}
        onScroll={(e) => {
          if (
            e.detail?.scrollHeight - e.detail?.scrollTop - clientHeight <
            50
          ) {
            needScrollToBottomRef.current = true;
          } else {
            needScrollToBottomRef.current = false;
          }
        }}
        className="chat-area"
      >
        {messageLinks.length > 0 ? (
          <MessageList isGenerating={isGenerating} />
        ) : (
          <View className="welcome-view">
            <WelcomeTips />
            <SuggestedQuestions questions={suggestedQuestionList} send={send} />
          </View>
        )}
      </ScrollView>
      <View className="input-bar">
        <InputBar
          hasHistoryMessage={messageLinks.length > 0}
          isGenerating={isGenerating}
          send={send}
          sendImage={sendImage}
          cancel={stopGenerating}
        />
      </View>
    </View>
    // </View>
  );
}
