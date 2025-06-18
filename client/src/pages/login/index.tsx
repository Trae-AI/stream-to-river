// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import React from 'react';
import { View } from '@tarojs/components';
import './index.less';
import CustomNavbar from '@/components/custom-navbar';
import LoginForm from '@/components/login-form';

const Login = () => {
  return (
    <View className="container">
      <CustomNavbar />
      <View className="full">
        <LoginForm />
      </View>
    </View>
  );
};

export default Login;
