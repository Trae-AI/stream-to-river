// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { Text, View } from "@tarojs/components";
import { ReactComponent as Back } from "@/assets/icons/back.svg";
import { useCallback } from "react";

export const Navbar = ({ back }: {back: () => void}) => {

  const handleClickNavBack = useCallback(() => {
    back();
  }, []);

  return (
    <View className="custom-navbar">
      <View className="sub-nav-container">
        <Back className="sub-nav-btn" onClick={handleClickNavBack} />
        <Text className="sub-nav-title">单词复习</Text>
      </View>
    </View>
  );
};
