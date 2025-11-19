import request from '@/core/api/request'

// 获取工作流列表
export function getWorkflows() {
  return request({
    url: '/workflows',
    method: 'get'
  })
}

// 获取工作流详情
export function getWorkflow(id) {
  return request({
    url: `/workflows/${id}`,
    method: 'get'
  })
}

// 执行工作流
export function executeWorkflow(id, data = {}) {
  return request({
    url: `/workflows/${id}/execute`,
    method: 'post',
    data
  })
}

// 获取执行历史
export function getExecutions(limit = 20) {
  return request({
    url: '/workflows/executions',
    method: 'get',
    params: { limit }
  })
}

// 获取执行详情
export function getExecution(id) {
  return request({
    url: `/workflows/executions/${id}`,
    method: 'get'
  })
}

// 删除执行记录
export function deleteExecution(id) {
  return request({
    url: `/workflows/executions/${id}`,
    method: 'delete'
  })
}

// 激活工作流
export function activateWorkflow(id) {
  return request({
    url: `/workflows/${id}/activate`,
    method: 'post'
  })
}

// 停用工作流
export function deactivateWorkflow(id) {
  return request({
    url: `/workflows/${id}/deactivate`,
    method: 'post'
  })
}

// 检查 n8n 健康状态
export function checkHealth() {
  return request({
    url: '/workflows/health',
    method: 'get'
  })
}

// ========== 工作流配置管理 ==========

// 获取所有工作流配置
export function getWorkflowConfigs() {
  return request({
    url: '/workflows/configs',
    method: 'get'
  })
}

// 获取指定工作流配置
export function getWorkflowConfig(id) {
  return request({
    url: `/workflows/configs/${id}`,
    method: 'get'
  })
}

// 创建工作流配置
export function createWorkflowConfig(data) {
  return request({
    url: '/workflows/configs',
    method: 'post',
    data
  })
}

// 更新工作流配置
export function updateWorkflowConfig(id, data) {
  return request({
    url: `/workflows/configs/${id}`,
    method: 'put',
    data
  })
}

// 删除工作流配置
export function deleteWorkflowConfig(id) {
  return request({
    url: `/workflows/configs/${id}`,
    method: 'delete'
  })
}

// 从 n8n 同步工作流配置
export function syncWorkflowConfigs() {
  return request({
    url: '/workflows/configs/sync',
    method: 'post'
  })
}

