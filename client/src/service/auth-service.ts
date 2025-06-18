// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import Taro from '@tarojs/taro';
import { ServerConfig } from './server-config';

interface LoginRequest {
  username: string;
  password: string;
}

interface RegisterRequest {
  username: string;
  email: string;
  password: string;
}

interface User {
  id: number;
  username: string;
  email: string;
}

interface AuthResponse {
  user: User;
  token: string;
}

class AuthService {
  private static instance: AuthService;
  private serverConfig: ServerConfig;

  private constructor() {
    this.serverConfig = ServerConfig.getInstance();
  }

  public static getInstance(): AuthService {
    if (!AuthService.instance) {
      AuthService.instance = new AuthService();
    }
    return AuthService.instance;
  }

  async login(loginData: LoginRequest): Promise<AuthResponse> {
    try {
      const response = await Taro.request({
        url: this.serverConfig.getFullUrl('/api/login'),
        method: 'POST',
        data: loginData,
        header: {
          'Content-Type': 'application/json'
        }
      });

      if (response.statusCode === 200) {
        const authData = response.data as AuthResponse;
        await this.setToken(authData.token);
        return authData;
      } else {
        throw new Error(response.data || '登录失败');
      }
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  }

  async register(registerData: RegisterRequest): Promise<AuthResponse> {
    try {
      const response = await Taro.request({
        url: this.serverConfig.getFullUrl('/api/register'),
        method: 'POST',
        data: registerData,
        header: {
          'Content-Type': 'application/json'
        }
      });

      if (response.statusCode === 200) {
        const authData = response.data as AuthResponse;
        await this.setToken(authData.token);
        return authData;
      } else {
        throw new Error(response.data || '注册失败');
      }
    } catch (error) {
      console.error('Register error:', error);
      throw error;
    }
  }

  async getUserInfo(): Promise<User> {
    try {
      const token = await this.getToken();
      if (!token) {
        throw new Error('未找到token');
      }

      const response = await Taro.request({
        url: this.serverConfig.getFullUrl('/api/user'),
        method: 'GET',
        header: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      if (response.statusCode === 200) {
        return response.data as User;
      } else {
        throw new Error('获取用户信息失败');
      }
    } catch (error) {
      console.error('Get user info error:', error);
      throw error;
    }
  }

  async setToken(token: string): Promise<void> {
    try {
      await Taro.setStorageSync('jwt_token', token);
    } catch (error) {
      console.error('Set token error:', error);
      throw error;
    }
  }

  async getToken(): Promise<string | null> {
    try {
      return Taro.getStorageSync('jwt_token') || null;
    } catch (error) {
      console.error('Get token error:', error);
      return null;
    }
  }

  async clearToken(): Promise<void> {
    try {
      await Taro.removeStorageSync('jwt_token');
    } catch (error) {
      console.error('Clear token error:', error);
    }
  }

  async isLoggedIn(): Promise<boolean> {
    const token = await this.getToken();
    return !!token;
  }

  async logout(): Promise<void> {
    await this.clearToken();
  }
}

export default AuthService;
export type { LoginRequest, RegisterRequest, User, AuthResponse };
