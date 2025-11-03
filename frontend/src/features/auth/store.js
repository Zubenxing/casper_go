import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/core/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const refreshToken = ref(localStorage.getItem('refreshToken') || '')
  const user = ref(null)

  const isLoggedIn = computed(() => !!token.value)

  // 登录
  async function login(username, password) {
    const response = await api.auth.login(username, password)
    if (response.code === 0) {
      token.value = response.data.access_token
      refreshToken.value = response.data.refresh_token
      user.value = response.data.user
      
      localStorage.setItem('token', token.value)
      localStorage.setItem('refreshToken', refreshToken.value)
      localStorage.setItem('user', JSON.stringify(user.value))
      
      return true
    }
    return false
  }

  // 登出
  async function logout() {
    try {
      await api.auth.logout()
    } catch (error) {
      console.error('登出失败', error)
    } finally {
      token.value = ''
      refreshToken.value = ''
      user.value = null
      localStorage.removeItem('token')
      localStorage.removeItem('refreshToken')
      localStorage.removeItem('user')
    }
  }

  // 获取用户信息
  async function fetchUserInfo() {
    const response = await api.auth.getProfile()
    if (response.code === 0) {
      user.value = response.data
      localStorage.setItem('user', JSON.stringify(user.value))
    }
  }

  // 初始化（从 localStorage 恢复）
  function init() {
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      user.value = JSON.parse(savedUser)
    }
  }

  return {
    token,
    refreshToken,
    user,
    isLoggedIn,
    login,
    logout,
    fetchUserInfo,
    init
  }
})

