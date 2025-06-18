// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { Message } from "./types";
import { sleep, uuid } from "./utils";

export const sendMessage = async (message: string): Promise<Message> => {
  await sleep(5000);
  return {
    id: uuid(),
    text: message,
    creator: "bot",
    linkWordList: [
      {
        text: "doubao",
        range: [0, 6]
      }
    ]
  }
}
