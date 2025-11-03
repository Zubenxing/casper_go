// 证书监控相关 API
import http from '@/core/api/request'

export default {
  // ========== 证书监控管理 ==========
  // 获取证书列表
  getList: (page = 1, pageSize = 10) => 
    http.get('/certificates', { params: { page, page_size: pageSize } }),
  
  // 获取证书详情
  getDetail: (id) => 
    http.get(`/certificates/${id}`),
  
  // 添加证书监控
  add: (url) => 
    http.post('/certificates', { url }),
  
  // 更新证书信息（刷新证书）
  update: (id) => 
    http.put(`/certificates/${id}`),
  
  // 更新客户名和备注
  updateInfo: (id, data) =>
    http.patch(`/certificates/${id}/info`, data),
  
  // 删除证书监控
  delete: (id) => 
    http.delete(`/certificates/${id}`),
  
  // 检查所有证书
  checkAll: () => 
    http.post('/certificates/check-all'),

  // ========== 证书工具 ==========
  // 生成 CSR
  generateCSR: (data) => 
    http.post('/certificates/tools/generate-csr', data),
  
  // 验证 CSR
  validateCSR: (data) => 
    http.post('/certificates/tools/validate-csr', data),
  
  // 验证证书
  validateCert: (data) => 
    http.post('/certificates/tools/validate-cert', data)
}

