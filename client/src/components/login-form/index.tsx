// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import React, { useState } from 'react';
import { View, Input, Button, Text } from '@tarojs/components';
import Taro from '@tarojs/taro';
import AuthService from '@/service/auth-service';
import useUserStore from '@/store/user';
import './index.less';

interface LoginFormProps {
  onSuccess?: () => void;
}

const LoginForm: React.FC<LoginFormProps> = ({ onSuccess }) => {
  const [isLogin, setIsLogin] = useState(true);
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: ''
  });
  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const { setUserInfo, setToken, setLoginStatus } = useUserStore();
  const authService = AuthService.getInstance();

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.username.trim()) {
      newErrors.username = '用户名不能为空';
    }

    if (!isLogin && !formData.email.trim()) {
      newErrors.email = '邮箱不能为空';
    }

    if (!isLogin && formData.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      newErrors.email = '邮箱格式不正确';
    }

    if (!formData.password.trim()) {
      newErrors.password = '密码不能为空';
    }

    if (!isLogin && formData.password !== formData.confirmPassword) {
      newErrors.confirmPassword = '两次密码输入不一致';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async () => {
    if (!validateForm()) {
      return;
    }

    setLoading(true);
    try {
      let result;
      if (isLogin) {
        result = await authService.login({
          username: formData.username,
          password: formData.password
        });
      } else {
        result = await authService.register({
          username: formData.username,
          email: formData.email,
          password: formData.password
        });
      }

      setUserInfo({
        id: result.user.id,
        username: result.user.username,
        email: result.user.email
      });
      setToken(result.token);
      setLoginStatus(true);

      Taro.showToast({
        title: isLogin ? '登录成功' : '注册成功',
        icon: 'success'
      });

      setTimeout(() => {
        Taro.redirectTo({ url: '/pages/index/index' });
        onSuccess?.();
      }, 1000);

    } catch (error) {
      console.error('Auth error:', error);
      Taro.showToast({
        title: error.message || (isLogin ? '登录失败' : '注册失败'),
        icon: 'error'
      });
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (field: string, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    if (errors[field]) {
      setErrors(prev => ({ ...prev, [field]: '' }));
    }
  };

  return (
    <View className="login-form">
      <View className="form-header">
        <Text className="form-title">{isLogin ? '登录' : '注册'}</Text>
      </View>

      <View className="form-content">
        <View className="input-group">
          <Input
            className={`form-input ${errors.username ? 'error' : ''}`}
            placeholder="请输入用户名"
            value={formData.username}
            onInput={(e) => handleInputChange('username', e.detail.value)}
          />
          {errors.username && <Text className="error-text">{errors.username}</Text>}
        </View>

        {!isLogin && (
          <View className="input-group">
            <Input
              className={`form-input ${errors.email ? 'error' : ''}`}
              placeholder="请输入邮箱"
              value={formData.email}
              onInput={(e) => handleInputChange('email', e.detail.value)}
            />
            {errors.email && <Text className="error-text">{errors.email}</Text>}
          </View>
        )}

        <View className="input-group">
          <Input
            className={`form-input ${errors.password ? 'error' : ''}`}
            placeholder="请输入密码"
            password
            value={formData.password}
            onInput={(e) => handleInputChange('password', e.detail.value)}
          />
          {errors.password && <Text className="error-text">{errors.password}</Text>}
        </View>

        {!isLogin && (
          <View className="input-group">
            <Input
              className={`form-input ${errors.confirmPassword ? 'error' : ''}`}
              placeholder="请确认密码"
              password
              value={formData.confirmPassword}
              onInput={(e) => handleInputChange('confirmPassword', e.detail.value)}
            />
            {errors.confirmPassword && <Text className="error-text">{errors.confirmPassword}</Text>}
          </View>
        )}

        <Button
          className="submit-btn"
          onClick={handleSubmit}
          loading={loading}
          disabled={loading}
        >
          {loading ? '处理中...' : (isLogin ? '登录' : '注册')}
        </Button>

        <View className="switch-mode">
          <Text className="switch-text">
            {isLogin ? '还没有账号？' : '已有账号？'}
          </Text>
          <Text
            className="switch-link"
            onClick={() => {
              setIsLogin(!isLogin);
              setFormData({ username: '', email: '', password: '', confirmPassword: '' });
              setErrors({});
            }}
          >
            {isLogin ? '立即注册' : '立即登录'}
          </Text>
        </View>
      </View>
    </View>
  );
};

export default LoginForm;
