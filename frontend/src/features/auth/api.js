// 认证相关 API
import http from '@/core/api/request'

export default {
  // 登录
  login: (username, password) => 
    http.post('/auth/login', { username, password }),
  
  // 登出
  logout: () => 
    http.post('/logout'),
  
  // 刷新令牌
  refresh: (refreshToken) => 
    http.post('/auth/refresh', { refresh_token: refreshToken }),
  
  // 获取当前用户信息
  getProfile: () => 
    http.get('/profile')
}

