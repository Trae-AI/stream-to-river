// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { create } from 'zustand';

interface UserInfo {
  id: number;
  username: string;
  email: string;
}

interface UserState {
  userInfo: UserInfo | null;
  token: string | null;
  isLoggedIn: boolean;
  setUserInfo: (userInfo: UserInfo | null) => void;
  setToken: (token: string | null) => void;
  setLoginStatus: (status: boolean) => void;
  logout: () => void;
}

const useUserStore = create<UserState>((set) => ({
  userInfo: null,
  token: null,
  isLoggedIn: false,
  setUserInfo: (userInfo) => set({ userInfo }),
  setToken: (token) => set({ token }),
  setLoginStatus: (status) => set({ isLoggedIn: status }),
  logout: () => set({ userInfo: null, token: null, isLoggedIn: false }),
}));

export default useUserStore;
