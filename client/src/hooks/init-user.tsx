// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { useEffect } from 'react';
import useUserStore from '@/store/user';
import Taro from '@tarojs/taro';
import AuthService from '@/service/auth-service';

export const useInitUser = () => {
  const { setUserInfo, setToken, setLoginStatus, logout } = useUserStore();
  const authService = AuthService.getInstance();

  useEffect(() => {
    const initUser = async () => {
      try {
        const isLoggedIn = await authService.isLoggedIn();
        if (!isLoggedIn) {
          Taro.navigateTo({ url: '/pages/login/index' });
          return;
        }

        const userInfo = await authService.getUserInfo();
        const token = await authService.getToken();

        setUserInfo({
          id: userInfo.id,
          username: userInfo.username,
          email: userInfo.email
        });
        setToken(token);
        setLoginStatus(true);

      } catch (error) {
        console.error('Init user failed:', error);
        await authService.clearToken();
        logout();
        Taro.navigateTo({ url: '/pages/login/index' });
      }
    };

    initUser();
  }, []);
};
