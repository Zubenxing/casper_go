// 公共格式化工具函数

/**
 * 格式化时间
 */
export const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

/**
 * 格式化日期
 */
export const formatDate = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleDateString('zh-CN')
}

/**
 * 获取剩余天数标签类型
 */
export const getDaysLeftType = (days) => {
  if (days < 0) return 'danger'
  if (days <= 30) return 'warning'
  return 'success'
}

/**
 * 获取状态标签类型
 */
export const getStatusType = (status) => {
  const types = { 
    1: 'success',  // 正常
    2: 'warning',  // 警告
    3: 'danger',   // 过期
    0: 'info'      // 禁用
  }
  return types[status] || 'info'
}

