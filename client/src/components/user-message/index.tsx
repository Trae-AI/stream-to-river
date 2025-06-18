// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { Image, Text, View } from "@tarojs/components";
import "./index.less";

const UserMessage = ({
  message,
  type,
}: {
  message: string;
  type: "text" | "image";
}) => {
  return type === "image" ? (
    <Image src={message} className="user-message-text auto-width-image"></Image>
  ) : (
    <Text className="user-message-text">{message}</Text>
  );
};
export default UserMessage;
