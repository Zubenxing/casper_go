import request from '@/core/api/request'

// ==================== Work APIs ====================

// 创建工作记录
export function createWork(data) {
  return request({
    url: '/works',
    method: 'post',
    data
  })
}

// 获取工作记录列表
export function getWorkList(params) {
  return request({
    url: '/works',
    method: 'get',
    params
  })
}

// 获取工作记录详情
export function getWorkDetail(id) {
  return request({
    url: `/works/${id}`,
    method: 'get'
  })
}

// 更新工作记录
export function updateWork(id, data) {
  return request({
    url: `/works/${id}`,
    method: 'put',
    data
  })
}

// 删除工作记录
export function deleteWork(id) {
  return request({
    url: `/works/${id}`,
    method: 'delete'
  })
}

// 获取工作统计
export function getWorkStats() {
  return request({
    url: '/works/stats',
    method: 'get'
  })
}

// ==================== Work Issue APIs ====================

// 创建问题记录
export function createWorkIssue(data) {
  return request({
    url: '/work-issues',
    method: 'post',
    data
  })
}

// 获取问题记录列表
export function getWorkIssueList(params) {
  return request({
    url: '/work-issues',
    method: 'get',
    params
  })
}

// 获取问题记录详情
export function getWorkIssueDetail(id) {
  return request({
    url: `/work-issues/${id}`,
    method: 'get'
  })
}

// 更新问题记录
export function updateWorkIssue(id, data) {
  return request({
    url: `/work-issues/${id}`,
    method: 'put',
    data
  })
}

// 删除问题记录
export function deleteWorkIssue(id) {
  return request({
    url: `/work-issues/${id}`,
    method: 'delete'
  })
}

// 获取问题统计
export function getWorkIssueStats() {
  return request({
    url: '/work-issues/stats',
    method: 'get'
  })
}

// 上传问题截图
export function uploadIssueImage(file) {
  const formData = new FormData()
  formData.append('file', file)
  return request({
    url: '/work-issues/upload',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

