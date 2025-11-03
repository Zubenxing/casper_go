// 密码管理相关 API
import http from '@/core/api/request'

export default {
  // 创建账户
  create: (data) => 
    http.post('/passwords', data),
  
  // 获取账户列表
  getList: (params) => 
    http.get('/passwords', { params }),
  
  // 获取账户详情
  getDetail: (id) => 
    http.get(`/passwords/${id}`),
  
  // 获取密码（解密后）
  getPassword: (id) => 
    http.get(`/passwords/${id}/password`),
  
  // 更新账户
  update: (id, data) => 
    http.put(`/passwords/${id}`, data),
  
  // 删除账户
  delete: (id) => 
    http.delete(`/passwords/${id}`),
  
  // 获取分类列表
  getCategories: () => 
    http.get('/passwords/categories'),
  
  // 获取统计信息
  getStats: () => 
    http.get('/passwords/stats')
}

