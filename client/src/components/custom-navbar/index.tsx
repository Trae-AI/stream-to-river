// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

// components/CustomNavbar/index.tsx
import { View, Text } from "@tarojs/components";
import Taro, { useRouter } from "@tarojs/taro";
import { useState } from "react";
import { ReactComponent as Back } from "@/assets/icons/back.svg";
import { ReactComponent as Avatar } from "@/assets/images/avatar.svg";
import "./index.less";
import useUserStore from "@/store/user";
import { LogoutPopup } from "../logout-popup";

const NAV_ITEMS = [
  { id: 0, title: "问答", url: "/pages/index/index", type: "main" },
  { id: 1, title: "学习", url: "/pages/study_main/index", type: "main" },
  {
    id: 2,
    title: "单词学习",
    url: "/pages/words_study/index",
    type: "sub-nav",
  },
  {
    id: 3,
    title: "登录/注册",
    url: "/pages/login/index",
  },
];

const CustomNavbar = () => {
  const { userInfo } = useUserStore();
  const [logoutVisible, setLogoutVisible] = useState(false);

  const router = useRouter();

  const currentPath = router.path.split("?")[0]; // e.g., "/pages/index/index"


  const handleTabClick = async (path: string) => {
    if (path !== currentPath) {
      Taro.redirectTo({
        url: path,
      });
    }
  };

  const currentNavItem = NAV_ITEMS.find((item) => item.url === currentPath);

  return (
    <View className="custom-navbar">
      {currentNavItem?.type === "main" ? (
        <View className="tab-container">
          <View>
            {NAV_ITEMS.filter((tab) => tab.type === "main").map((tab) => (
              <Text
                key={tab.id}
                className={`tab-item ${
                  currentPath === tab.url ? "active" : ""
                }`}
                onClick={() => handleTabClick(tab.url)}
              >
                {tab.title}
              </Text>
            ))}
          </View>
          {userInfo?.id ? (
            <Avatar className="avatar" onClick={() => setLogoutVisible(true)} />
          ) : null}
        </View>
      ) : (
        <View className="sub-nav-container">
          <Back
            className="sub-nav-btn"
            onClick={() => {
              const pages = Taro.getCurrentPages() || [];
              if (
                pages.length === 1 && currentPath === "/pages/words_study/index"
              ) {
                Taro.redirectTo({
                  url: "/pages/study_main/index",
                });
                return;
              }

              Taro.navigateBack();
            }}
          />
          <Text className="sub-nav-title">{currentNavItem?.title}</Text>
        </View>
      )}

      <LogoutPopup
        visible={logoutVisible}
        onClose={() => setLogoutVisible(false)}
      />
    </View>
  );
};

export default CustomNavbar;
