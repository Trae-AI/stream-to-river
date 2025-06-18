// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { PropsWithChildren } from "react";
import Taro, { useLaunch } from "@tarojs/taro";
import { View, Image } from "@tarojs/components";
import { ReactComponent as TopSVG } from "@/assets/images/top.svg";
import useSystemTagsStore from "@/store/tag";
import TagsService from "@/service/tag-service";
import { useInitUser } from "./hooks/init-user";

import "./app.less";

function App({ children }: PropsWithChildren<any>) {
  const { setSystemTags } = useSystemTagsStore();
  useLaunch(() => {
    TagsService.getInstance().getTags().then((res) => {
      setSystemTags(res.tags);
    })

    // prevent long press context menu
    window.oncontextmenu = function () {
      return false;
    };
    // global taro object
    window.Taro = Taro;
  });

  useInitUser();

  return (
    <View className="container">
      <TopSVG className="top-svg" />
      {children}
    </View>
  );
}

export default App;
