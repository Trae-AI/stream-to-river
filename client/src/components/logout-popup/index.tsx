// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { View, Text } from '@tarojs/components';
import './index.less'
import useUserStore from '@/store/user';
import AuthService from '@/service/auth-service';
import Taro from '@tarojs/taro';

interface LogoutPopupProps {
  visible: boolean;
  onClose: () => void;
}

export const LogoutPopup = ({ visible, onClose }: LogoutPopupProps) => {

  const { userInfo, logout } = useUserStore();
  const authService = AuthService.getInstance();

  const handleLogout = async () => {
    try {
      await authService.logout();
      logout();
      Taro.redirectTo({ url: '/pages/login/index' });
      onClose();
    } catch (error) {
      console.error('Logout error:', error);
      logout();
      onClose();
    }
  };

  return (
    <>
      {visible && (
        <View className='popup-container'>
          <View className='popup-mask' onClick={onClose} />
          <View className='popup-content'>
            <View className='popup-header'>
              <Text className='popup-title'>👋{' '}{userInfo?.username}</Text>
            </View>

            <View className='popup-actions'>
              <View className='action-item danger' onClick={handleLogout}>
                <Text>退出登录</Text>
              </View>

              <View className='action-item' onClick={onClose}>
                <Text>取消</Text>
              </View>
            </View>

            <View className='safe-line'></View>
          </View>
        </View>
      )}
    </>
  );
};

export default LogoutPopup;
