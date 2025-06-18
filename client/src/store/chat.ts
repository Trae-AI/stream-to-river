// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { create } from "zustand";
import { v4 as uuid } from "uuid";
import { fetchEventSource } from "@microsoft/fetch-event-source";
import { produce } from "immer";

export interface Message {
  text: string;
  creator: "user" | "bot";
  id: string;
  linkWordList?: Array<{ start: number; end: number; text: string }>;
  type: "image" | "text";
  isDone?: boolean;
}

interface MessageStoreSlice {
  messageLinks: string[];
  messageMap: { [key: string]: Message };
  isGenerating: boolean;
  addMessage: (message: Message) => void;
  sendMessage: (text: string) => void;
  stopGenerating: () => void;
  sendImage: (base64Image: string) => void;
  abortController: AbortController | null;
  conversationId: string;
  setConversationId: (conversationId: string) => void;
}

export const useMessageStore = create<MessageStoreSlice>((set, get) => ({
  messageLinks: [],
  messageMap: {},
  isGenerating: false,
  abortController: null,
  conversationId: "",

  setConversationId: (conversationId: string) => {
    set({ conversationId });
  },

  addMessage: (message: Message) => {
    set((state) => ({
      messageMap: {
        ...state.messageMap,
        [message.id]: message,
      },
    }));
  },

  stopGenerating: () => {
    const { abortController } = get();
    const botMessageId = get().messageLinks[get().messageLinks.length - 1];

    if (abortController) {
      abortController.abort();
      set(
        produce((state) => {
          // Set message completion status
          if (state.messageMap[botMessageId]) {
            state.messageMap[botMessageId].isDone = true;
          }
          state.isGenerating = false;
          state.abortController = null;
        })
      );
    }
  },

  sendImage: async (base64Image: string) => {
    const messageId = uuid();
    set(
      produce((state) => {
        state.messageLinks.push(messageId);
        state.messageMap[messageId] = {
          text: base64Image,
          creator: "user",
          id: messageId,
          type: "image",
        };
      })
    );
  },

  sendMessage: (text: string) => {
    // Create user message
    const messageId = uuid();
    const { conversationId, stopGenerating } = get();

    stopGenerating();

    // Create new AbortController
    const abortController = new AbortController();
    set({ abortController });

    set(
      produce((state) => {
        state.messageLinks.push(messageId);
        state.messageMap[messageId] = {
          text,
          creator: "user",
          id: messageId,
          type: "text",
        };
        state.isGenerating = true;
      })
    );

    // Create bot message placeholder
    const botMessageId = uuid();
    set(
      produce((state) => {
        state.messageLinks.push(botMessageId);
        state.messageMap[botMessageId] = {
          text: "",
          creator: "bot",
          id: botMessageId,
          linkWordList: [],
          type: "text",
          isDone: false, // Initialize as false
        };
      })
    );

    // Encode message content
    const encodedMessage = encodeURIComponent(text);

    fetchEventSource(
      `/api/chat?q=${encodedMessage}&conversation_id=${conversationId}`,
      {
        openWhenHidden: true,
        method: "GET",
        headers: {
          Accept: "text/event-stream",
        },
        signal: abortController.signal,
        onmessage: (event) => {
          // Handle server-sent messages
          if (event.event === "message") {
            set(
              produce((state) => {
                const currentMessage = state.messageMap[botMessageId];
                if (currentMessage) {
                  const jsonData = JSON.parse(event.data);
                  const msg = jsonData.msg || "";
                  let linkWordList = [];
                  if (jsonData.extra.meta_info) {
                    const items = JSON.parse(jsonData.extra.meta_info)[0]
                      ?.items;
                    linkWordList = items.map((item: any) => ({
                      start: item.start,
                      end: item.end,
                      text: item.text,
                      type: "text",
                    }));
                  }
                  currentMessage.linkWordList.push(...linkWordList);
                  currentMessage.text += msg;
                }
              })
            );
          }
        },
        onclose: () => {
          set(
            produce((state) => {
              // Set message completion status
              if (state.messageMap[botMessageId]) {
                state.messageMap[botMessageId].isDone = true;
              }
              state.isGenerating = false;
              state.abortController = null;
            })
          );
        },
        onerror: (err) => {
          if (err.name === "AbortError") {
            console.warn("Request aborted");
          } else {
            console.error("SSE error:", err);
          }
          set(
            produce((state) => {
              // Set message completion status
              if (state.messageMap[botMessageId]) {
                state.messageMap[botMessageId].isDone = true;
              }
              state.isGenerating = false;
              state.abortController = null;
            })
          );
          throw err;
        },
      }
    );
  },
}));
