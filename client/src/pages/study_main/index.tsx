// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

import { useState, useEffect, useMemo, useRef, useCallback } from 'react';
import { View, Text, Button, ScrollView } from "@tarojs/components";
import Taro, { useDidShow } from "@tarojs/taro";
import AuthService from '@/service/auth-service';
import useUserStore from "@/store/user";
import ReviewPuzzleService from "@/service/review-service";
import ProgressGauge from "@/components/progress-gauge";
import WordList, { WordListRef } from "@/components/word-list";
import { ReactComponent as ReviewIcon } from "@/assets/icons/review.svg";
import { LogoutPopup } from "@/components/logout-popup";

import CustomNavbar from "@/components/custom-navbar";
import { ReactComponent as LoadingIcon } from "@/assets/images/loading.svg";
import { ReactComponent as Logo } from "@/assets/images/logo.svg";
import { ReactComponent as Slogan } from "@/assets/images/slogan.svg";
import { ReactComponent as TraeLogo } from "@/assets/images/trae.svg";

import "./index.less";
import VocabularyReviewComponent from "@/components/vocabulary-review";
import { STORAGE_WORD_KEY } from "@/consts";

const StudyMain = () => {
  const { userInfo, setUserInfo, isLoggedIn, setLoginStatus } = useUserStore();
  const [hasData, setHasData] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [visible, setVisible] = useState(false);
  const [puzzleProgress, setPuzzleProgress] = useState({
    CompleteCount: 0,
    PendingCount: 0,
    TotalCount: 0,
    AllCompletedCount: 0,
  });
  const [gaugeHeight, setGaugeHeight] = useState(0);
  const [isTopBarShow, setIsTopBarShow] = useState(false);
  const gaugeRef = useRef<any>(null);
  const authService = AuthService.getInstance();

  const wordListRef = useRef<WordListRef>(null);
  const [showLogout, setShowLogout] = useState(false);

  useEffect(() => {
    // When the page is first opened, automatically clear the historical word records
    Taro.setStorageSync(STORAGE_WORD_KEY, "");
  }, []);

  const handleRefresh = () => {
    if (wordListRef.current) {
      wordListRef.current?.refresh();
    }
  };

  const dynamicContent = useMemo(() => {
    const percentage = Math.round(
      (puzzleProgress.CompleteCount / puzzleProgress.PendingCount) * 100
    );
    if (percentage >= 80) {
      return {
        percentage: percentage,
        gradientColors: ["#00CC6A", "#6BF18C"], // Green gradient
      };
    } else if (percentage >= 50) {
      return {
        percentage: percentage,
        gradientColors: ["#0578FF", "#CBE8FF"], // Blue gradient
      };
    } else {
      return {
        percentage: percentage,
        gradientColors: ["#FF9C1A", "#FFD836"], // Orange gradient
      };
    }
  }, [puzzleProgress]);

  const fetchData = async () => {
    setIsLoading(true);
    let user_id = -1;

    if (!userInfo?.id) {
      try {
        // Check if already logged in
        const loggedIn = await authService.isLoggedIn();
        if (!loggedIn) {
          setLoginStatus(false);
          setHasData(false);
          setIsLoading(false);
          return;
        }

        // Get user information
        const userData = await authService.getUserInfo();
        user_id = userData.id;
        setUserInfo({
          id: userData.id,
          username: userData.username,
          email: userData.email
        });
        setLoginStatus(true);

      } catch (e) {
        console.error('Failed to get user information:', e);
        setLoginStatus(false);
        setHasData(false);
        setIsLoading(false);
        return;
      }
    } else {
      user_id = userInfo?.id;
      setLoginStatus(true);
    }

    try {
      if (user_id > 0) {
        handleRefresh();
        ReviewPuzzleService.getInstance()
          .getReviewProgress()
          .then((res) => {
            setHasData(res.total_words_count > 0);
            setPuzzleProgress({
              CompleteCount: res.completed_review_count,
              PendingCount: res.pending_review_count,
              TotalCount: res.total_words_count,
              AllCompletedCount: res.all_completed_count,
            });
          })
          .catch((error) => {
            console.error("Error fetching user data:", error);
            setHasData(false);
          })
          .finally(() => {
            setIsLoading(false);
          });
      } else {
        setIsLoading(false);
      }
    } catch (e) {
      console.error('Failed to fetch data:', e);
      setIsLoading(false);
    }
  };

  useDidShow(() => {
    if (visible === true) {
      return;
    }
    fetchData();
  });

  const handleHideComponent = () => {
    fetchData();
    setVisible(false);
  };

  useEffect(() => {
    if (hasData && isLoggedIn) {
      setTimeout(() => {
        Taro.createSelectorQuery()
          .select(".gauge-container")
          .boundingClientRect((rect) => {
            if (rect) {
              setGaugeHeight((rect as any).height);
            }
          })
          .exec();
      }, 100);
    }
  }, [hasData, isLoggedIn]);

  // Scroll event handling
  const handleScroll = (e) => {
    const currentScrollTop = e.detail.scrollTop;

    if (gaugeHeight > 0) {
      // Calculate gauge visibility
      const gaugeVisibleHeight = Math.max(0, gaugeHeight - currentScrollTop);
      const visibilityPercentage = (gaugeVisibleHeight / gaugeHeight) * 100;

      // Switch to progress bar mode when gauge is more than 50% invisible
      const shouldShowTopBar = visibilityPercentage < 1;
      if (shouldShowTopBar != isTopBarShow) {
        setIsTopBarShow(shouldShowTopBar);
      }
    }
  };

  const handleAddWord = () => {
    // Navigate to add word page
    Taro.redirectTo({ url: "/pages/index/index" });
  };

  const handleLogin = () => {
    Taro.navigateTo({ url: "/pages/login/index" });
  };

  const handleStartReview = () => {
    // Start review
    setVisible(true);
  };

  const handleClickTraeLogo = useCallback(() => {
    window.open("https://www.trae.cn/");
  }, []);

  if (isLoading) {
    return (
      <>
        <CustomNavbar />
        <View className="loading-container">
          <LoadingIcon className="loading-container__image" />
          <Text className="loading-container__text">加载中</Text>
        </View>
      </>
    );
  }

  // State 1: Not logged in
  if (!isLoggedIn) {
    return (
      <>
        <CustomNavbar />
        <View className="vertical-center-main">
          <Logo className="logo"></Logo>

          <Slogan className="welcome-text"></Slogan>

          <View className="login-section">
            <Button className="login-btn" onClick={handleLogin}>
              立即登录
            </Button>
          </View>

          <View className="trae-logo-container" onClick={handleClickTraeLogo}>
            <TraeLogo className="trae-logo" />
          </View>
        </View>
      </>
    );
  }

  // State 2: Logged in but no data
  if (isLoggedIn && !hasData) {
    return (
      <>
        <CustomNavbar />
        <View className="vertical-center-main">
          <View className="no-data-container">
            <View className="logo"></View>
            <View className="text_area">
              <Text className="username">👋 Hi, {userInfo?.username}</Text>
              <Text className="welcome-title">
                欢迎来到"积流成江"，和"J"一起涨词量！
              </Text>
            </View>
            <Button className="add-word-btn" onClick={handleAddWord}>
              去查询添加
            </Button>
          </View>
          <Text className="logout-btn" onClick={() => setShowLogout(true)}>
            退出登录
          </Text>
          <LogoutPopup
            visible={showLogout}
            onClose={() => setShowLogout(false)}
          />
        </View>
      </>
    );
  }

  // State 3: Logged in with data
  const svgWidth = 100;
  const svgHeight = 12;
  const gradientColors = dynamicContent.gradientColors;

  return (
    <>
      <CustomNavbar />
      <View className="study-main-container">
        {isTopBarShow && (
          <View className="fixed-top-container">
            <View className="fixed-top-container__progress_container">
              <Text className="ixed-top-container__progress_container__completed-text">
                {puzzleProgress.CompleteCount}
              </Text>
              <Text className="fixed-top-container__progress_container__total-text">
                /{puzzleProgress.PendingCount}
              </Text>
              <svg
                className="fixed-top-container__progress_container__svg"
                width={svgWidth}
                height={svgHeight}
                viewBox={`0 0 ${svgWidth} ${svgHeight}`}
                preserveAspectRatio="xMidYMid meet"
              >
                <defs>
                  <linearGradient
                    id="topProgressGradient"
                    x1="0%"
                    y1="0%"
                    x2="100%"
                    y2="0%"
                  >
                    <stop offset="0%" stopColor={gradientColors[0]} />
                    <stop offset="100%" stopColor={gradientColors[1]} />
                  </linearGradient>
                </defs>
                <line
                  x1="0"
                  y1={svgHeight / 2}
                  x2="100"
                  y2={svgHeight / 2}
                  stroke="#FFFFFF"
                  strokeWidth="3"
                  strokeLinecap="round"
                />
                {/* url(#topProgressGradient) */}
                <line
                  x1="0"
                  y1={svgHeight / 2}
                  x2={((dynamicContent.percentage + 1) / 100) * svgWidth}
                  y2={svgHeight / 2}
                  stroke={gradientColors[1]}
                  strokeWidth="6"
                  strokeLinecap="round"
                />
              </svg>
            </View>

            <View className="fixed-top-container__review_container">
              <Button className="review-button-btn" onClick={handleStartReview}>
                <ReviewIcon className="review-button-btn-icon" />
                <Text className="review-button-btn-text">开始复习</Text>
              </Button>
            </View>
          </View>
        )}

        <ScrollView
          className="main-scroll-view"
          scrollY
          onScroll={handleScroll}
          scrollWithAnimation
          enhanced
          showScrollbar={false}
        >
          {/* Dashboard area */}
          <View className="gauge-container" ref={gaugeRef}>
            <ProgressGauge
              current={puzzleProgress.CompleteCount}
              total={puzzleProgress.PendingCount}
            />
          </View>

          <View className="total-progress-parent">
            <View className="total-progress-container">
              <Text className="total-progress-container-text">
                总单词 {puzzleProgress.TotalCount} | 已掌握{" "}
                {puzzleProgress.AllCompletedCount}
              </Text>
            </View>
          </View>

          {!isTopBarShow && (
            <View className="review-button-container">
              <Button className="review-button-btn" onClick={handleStartReview}>
                <ReviewIcon className="review-button-btn-icon" />
                <Text className="review-button-btn-text">开始复习</Text>
              </Button>
            </View>
          )}

          {/* Word list area */}
          <View className="word-list-container">
            <WordList isEnableScroll={false} ref={wordListRef} />
          </View>

          {/* Bottom spacer to ensure full scrolling */}
          <View className="bottom-spacer" />
        </ScrollView>
        <VocabularyReviewComponent
          visible={visible}
          onClose={handleHideComponent}
        />
      </View>
    </>
  );
};

export default StudyMain;
