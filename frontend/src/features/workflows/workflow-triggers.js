/**
 * 工作流触发器策略模式实现
 * 支持多种工作流触发方式的扩展
 */

import { executeWorkflow } from './api'
import { ElMessage } from 'element-plus'

/**
 * 基础触发器策略接口
 */
export class WorkflowTriggerStrategy {
  /**
   * 执行工作流
   * @param {Object} workflow - 工作流对象
   * @param {Object} formData - 表单数据
   * @param {Object} config - 工作流配置
   * @returns {Promise<Object>} 执行结果
   */
  async execute(workflow, formData, config) {
    throw new Error('execute method must be implemented')
  }
  
  /**
   * 获取需要的字段
   * @returns {Array} 字段配置数组
   */
  getRequiredFields() {
    return []
  }
  
  /**
   * 验证表单数据
   * @param {Object} formData - 表单数据
   * @returns {Boolean} 是否有效
   */
  validate(formData) {
    return true
  }
  
  /**
   * 处理执行结果
   * @param {Object} result - 执行结果
   * @param {Object} workflow - 工作流对象
   * @returns {Object} 处理后的结果
   */
  handleResult(result, workflow) {
    return result
  }
}

/**
 * 简单触发器 - 无需参数直接执行
 */
export class SimpleTrigger extends WorkflowTriggerStrategy {
  async execute(workflow, formData, config) {
    console.log('简单触发器执行:', workflow.name)
    
    try {
      const res = await executeWorkflow(workflow.id, {})
      
      if (res.code === 500 || res.code !== 0) {
        throw new Error(res.message || res.msg || '执行失败')
      }
      
      if (!res.data) {
        throw new Error('返回数据为空')
      }
      
      return {
        success: true,
        data: res.data,
        message: '执行成功'
      }
    } catch (error) {
      console.error('简单触发器执行失败:', error)
      throw error
    }
  }
}

/**
 * 文件上传触发器 - 处理文件上传和下载
 */
export class FileUploadTrigger extends WorkflowTriggerStrategy {
  async execute(workflow, formData, config) {
    console.log('文件上传触发器执行:', workflow.name, formData)
    
    try {
      // 创建 FormData 对象
      const form = new FormData()
      for (const [key, value] of Object.entries(formData)) {
        if (value) {
          form.append(key, value)
        }
      }
      
      // 通过后端 API 执行工作流（后端会调用 Webhook）
      const res = await executeWorkflow(workflow.id, form)
      
      if (res.code === 500 || res.code !== 0) {
        throw new Error(res.message || res.msg || '执行失败')
      }
      
      if (!res.data) {
        throw new Error('返回数据为空')
      }
      
      // 检查是否是文件下载
      if (res.data.data && res.data.data.isFile) {
        console.log('检测到文件下载，开始处理...')
        await this.handleFileDownload(res.data, workflow)
        return {
          success: true,
          data: res.data,
          message: '✅ 文件已下载！'
        }
      }
      
      return {
        success: true,
        data: res.data,
        message: '执行成功'
      }
    } catch (error) {
      console.error('文件上传触发器执行失败:', error)
      throw error
    }
  }
  
  async handleFileDownload(data, workflow) {
    try {
      const result = data.data.result
      let blob
      
      // 尝试将结果转换为二进制
      if (result.startsWith('data:')) {
        // 如果是 base64 字符串
        const base64Data = result.split(',')[1]
        const binaryString = atob(base64Data)
        const bytes = new Uint8Array(binaryString.length)
        for (let i = 0; i < binaryString.length; i++) {
          bytes[i] = binaryString.charCodeAt(i)
        }
        blob = new Blob([bytes], { type: data.data.contentType || 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
      } else {
        // 直接作为二进制处理
        const bytes = new Uint8Array(result.length)
        for (let i = 0; i < result.length; i++) {
          bytes[i] = result.charCodeAt(i)
        }
        blob = new Blob([bytes], { type: data.data.contentType || 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
      }
      
      // 触发下载
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `${workflow.name}_${new Date().toISOString().slice(0, 10)}.xlsx`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
      
    } catch (error) {
      console.error('下载文件失败:', error)
      throw new Error('文件下载失败')
    }
  }
}

/**
 * 表单输入触发器 - 处理表单数据
 */
export class FormInputTrigger extends WorkflowTriggerStrategy {
  async execute(workflow, formData, config) {
    console.log('表单输入触发器执行:', workflow.name, formData)
    
    try {
      const res = await executeWorkflow(workflow.id, formData)
      
      if (res.code === 500 || res.code !== 0) {
        throw new Error(res.message || res.msg || '执行失败')
      }
      
      if (!res.data) {
        throw new Error('返回数据为空')
      }
      
      return {
        success: true,
        data: res.data,
        message: '执行成功'
      }
    } catch (error) {
      console.error('表单输入触发器执行失败:', error)
      throw error
    }
  }
  
  validate(formData) {
    // 基础验证逻辑
    return formData && Object.keys(formData).length > 0
  }
}

/**
 * Webhook 触发器 - 直接调用 Webhook URL
 */
export class WebhookTrigger extends WorkflowTriggerStrategy {
  async execute(workflow, formData, config) {
    console.log('Webhook 触发器执行:', workflow.name)
    
    try {
      // 获取 Webhook URL（可以从配置或工作流元数据中获取）
      const webhookUrl = config.webhookUrl || this.getWebhookUrl(workflow)
      
      if (!webhookUrl) {
        throw new Error('未找到 Webhook URL')
      }
      
      // 直接调用 Webhook
      const response = await fetch(webhookUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData)
      })
      
      if (!response.ok) {
        throw new Error(`Webhook 调用失败: ${response.status}`)
      }
      
      const result = await response.json()
      
      return {
        success: true,
        data: result,
        message: 'Webhook 执行成功'
      }
    } catch (error) {
      console.error('Webhook 触发器执行失败:', error)
      throw error
    }
  }
  
  getWebhookUrl(workflow) {
    // 从工作流元数据中提取 Webhook URL
    // 这里需要根据实际的工作流结构来获取
    return workflow.webhookUrl || null
  }
}

/**
 * 批量处理触发器 - 处理批量数据
 */
export class BatchTrigger extends WorkflowTriggerStrategy {
  async execute(workflow, formData, config) {
    console.log('批量处理触发器执行:', workflow.name)
    
    try {
      const batchData = formData.batchData || []
      const results = []
      
      // 分批处理
      const batchSize = config.batchSize || 10
      for (let i = 0; i < batchData.length; i += batchSize) {
        const batch = batchData.slice(i, i + batchSize)
        
        const res = await executeWorkflow(workflow.id, {
          batchData: batch,
          batchIndex: Math.floor(i / batchSize)
        })
        
        if (res.code === 500 || res.code !== 0) {
          throw new Error(`批次 ${Math.floor(i / batchSize)} 执行失败: ${res.message || res.msg}`)
        }
        
        results.push(res.data)
      }
      
      return {
        success: true,
        data: {
          results,
          total: batchData.length,
          batches: Math.ceil(batchData.length / batchSize)
        },
        message: '批量处理完成'
      }
    } catch (error) {
      console.error('批量处理触发器执行失败:', error)
      throw error
    }
  }
  
  validate(formData) {
    return formData.batchData && Array.isArray(formData.batchData) && formData.batchData.length > 0
  }
}