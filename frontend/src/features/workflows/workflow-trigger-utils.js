/**
 * 工作流触发器工具函数
 * 提供触发器策略的注册、管理和扩展功能
 */

import { WorkflowTriggerFactory } from './workflow-trigger-factory'
import { WORKFLOW_TYPES } from './workflow-configs'

/**
 * 工作流触发器管理器
 * 提供统一的触发器策略管理接口
 */
export class WorkflowTriggerManager {
  /**
   * 获取所有可用的触发器类型
   * @returns {Array<Object>} 触发器类型列表
   */
  static getAvailableTriggers() {
    const types = WorkflowTriggerFactory.getSupportedTypes()
    return types.map(type => ({
      type,
      name: WorkflowTriggerFactory.getTriggerName(type),
      description: WorkflowTriggerFactory.getTriggerDescription(type),
      info: WorkflowTriggerFactory.getTriggerInfo(type)
    }))
  }
  
  /**
   * 注册自定义触发器策略
   * @param {string} type - 触发器类型
   * @param {Class} StrategyClass - 策略类
   * @param {Object} options - 配置选项
   */
  static registerCustomTrigger(type, StrategyClass, options = {}) {
    try {
      WorkflowTriggerFactory.registerStrategy(type, StrategyClass)
      
      if (options.name) {
        // 可以扩展工厂方法来支持自定义名称和描述
        console.log(`✅ 已注册自定义触发器: ${type} - ${options.name}`)
      }
      
      return true
    } catch (error) {
      console.error(`❌ 注册自定义触发器失败:`, error)
      return false
    }
  }
  
  /**
   * 验证工作流配置
   * @param {Object} config - 工作流配置
   * @returns {Object} 验证结果
   */
  static validateWorkflowConfig(config) {
    const result = {
      valid: true,
      errors: [],
      warnings: []
    }
    
    if (!config.type) {
      result.valid = false
      result.errors.push('工作流类型不能为空')
      return result
    }
    
    if (!WorkflowTriggerFactory.isSupported(config.type)) {
      result.valid = false
      result.errors.push(`不支持的工作流类型: ${config.type}`)
      return result
    }
    
    // 根据类型进行特定验证
    switch (config.type) {
      case WORKFLOW_TYPES.UPLOAD_FILES:
        if (!config.fields || config.fields.length === 0) {
          result.warnings.push('文件上传类型工作流建议配置上传字段')
        }
        break
        
      case WORKFLOW_TYPES.FORM_INPUT:
        if (!config.fields || config.fields.length === 0) {
          result.warnings.push('表单输入类型工作流建议配置输入字段')
        }
        break
        
      case WORKFLOW_TYPES.WEBHOOK:
        if (!config.webhookUrl) {
          result.warnings.push('Webhook 类型工作流建议配置 Webhook URL')
        }
        break
        
      case WORKFLOW_TYPES.BATCH:
        if (!config.batchSize || config.batchSize <= 0) {
          result.warnings.push('批量处理类型工作流建议配置合理的批次大小')
        }
        break
    }
    
    return result
  }
  
  /**
   * 获取工作流类型的配置模板
   * @param {string} type - 工作流类型
   * @returns {Object|null} 配置模板
   */
  static getConfigTemplate(type) {
    const templates = {
      [WORKFLOW_TYPES.SIMPLE]: {
        type: WORKFLOW_TYPES.SIMPLE,
        description: '简单执行工作流',
        executeLabel: '执行',
        resultType: 'auto'
      },
      
      [WORKFLOW_TYPES.UPLOAD_FILES]: {
        type: WORKFLOW_TYPES.UPLOAD_FILES,
        description: '文件上传处理工作流',
        executeLabel: '上传并执行',
        resultType: 'download',
        fields: [
          {
            name: 'file',
            label: '文件',
            type: 'file',
            accept: '*/*',
            required: true,
            description: '请选择要上传的文件'
          }
        ]
      },
      
      [WORKFLOW_TYPES.FORM_INPUT]: {
        type: WORKFLOW_TYPES.FORM_INPUT,
        description: '表单输入工作流',
        executeLabel: '提交',
        resultType: 'json',
        fields: [
          {
            name: 'input',
            label: '输入内容',
            type: 'text',
            required: true,
            placeholder: '请输入内容'
          }
        ]
      },
      
      [WORKFLOW_TYPES.WEBHOOK]: {
        type: WORKFLOW_TYPES.WEBHOOK,
        description: 'Webhook 调用工作流',
        executeLabel: '调用',
        resultType: 'json',
        webhookUrl: 'https://your-webhook-url.com/webhook',
        fields: []
      },
      
      [WORKFLOW_TYPES.BATCH]: {
        type: WORKFLOW_TYPES.BATCH,
        description: '批量处理工作流',
        executeLabel: '批量处理',
        resultType: 'json',
        batchSize: 10
      }
    }
    
    return templates[type] || null
  }
  
  /**
   * 创建默认配置
   * @param {string} workflowName - 工作流名称
   * @param {string} type - 工作流类型
   * @returns {Object} 默认配置
   */
  static createDefaultConfig(workflowName, type = WORKFLOW_TYPES.SIMPLE) {
    const template = this.getConfigTemplate(type)
    if (!template) {
      return {
        type: WORKFLOW_TYPES.SIMPLE,
        description: workflowName,
        executeLabel: '执行工作流',
        resultType: 'auto'
      }
    }
    
    return {
      ...template,
      description: template.description.replace('工作流', workflowName)
    }
  }
}

/**
 * 触发器工具函数
 */
export const TriggerUtils = {
  /**
   * 判断工作流是否需要用户输入
   * @param {Object} workflow - 工作流对象
   * @returns {Boolean} 是否需要输入
   */
  workflowNeedsInput(workflow) {
    if (!workflow || !workflow.config) return false
    
    const { type, fields } = workflow.config
    
    // 简单类型不需要输入
    if (type === WORKFLOW_TYPES.SIMPLE) return false
    
    // 有配置字段的类型需要输入
    if (fields && fields.length > 0) return true
    
    // 文件上传和表单输入类型默认需要输入
    return [WORKFLOW_TYPES.UPLOAD_FILES, WORKFLOW_TYPES.FORM_INPUT, WORKFLOW_TYPES.BATCH].includes(type)
  },
  
  /**
   * 获取工作流的执行模式
   * @param {Object} workflow - 工作流对象
   * @returns {String} 执行模式
   */
  getExecutionMode(workflow) {
    if (!workflow || !workflow.config) return 'auto'
    
    const { type } = workflow.config
    
    const modeMap = {
      [WORKFLOW_TYPES.SIMPLE]: '自动',
      [WORKFLOW_TYPES.UPLOAD_FILES]: '文件上传',
      [WORKFLOW_TYPES.FORM_INPUT]: '表单输入',
      [WORKFLOW_TYPES.WEBHOOK]: 'Webhook',
      [WORKFLOW_TYPES.BATCH]: '批量处理'
    }
    
    return modeMap[type] || '自动'
  },
  
  /**
   * 格式化执行结果
   * @param {Object} result - 执行结果
   * @returns {Object} 格式化后的结果
   */
  formatExecutionResult(result) {
    if (!result) return null
    
    const formatted = {
      id: result.id || '',
      status: result.status || 'unknown',
      workflowName: result.workflowName || '',
      startedAt: result.startedAt || '',
      stoppedAt: result.stoppedAt || '',
      data: result.data || null,
      message: result.message || ''
    }
    
    // 添加状态文本
    const statusTextMap = {
      'success': '成功',
      'error': '失败',
      'waiting': '等待中',
      'running': '运行中'
    }
    
    formatted.statusText = statusTextMap[formatted.status] || formatted.status
    
    return formatted
  }
}