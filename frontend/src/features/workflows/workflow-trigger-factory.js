/**
 * 工作流触发器工厂
 * 根据工作流类型动态创建对应的触发器策略
 */

import { WORKFLOW_TYPES } from './workflow-configs'
import {
  SimpleTrigger,
  FileUploadTrigger,
  FormInputTrigger,
  WebhookTrigger,
  BatchTrigger
} from './workflow-triggers'

/**
 * 工作流触发器工厂类
 */
export class WorkflowTriggerFactory {
  /**
   * 触发器策略映射表
   * 支持动态扩展新的触发器类型
   */
  static strategies = {
    [WORKFLOW_TYPES.SIMPLE]: SimpleTrigger,
    [WORKFLOW_TYPES.UPLOAD_FILES]: FileUploadTrigger,
    [WORKFLOW_TYPES.FORM_INPUT]: FormInputTrigger,
    // 扩展的触发器类型
    'webhook': WebhookTrigger,
    'batch': BatchTrigger
  }
  
  /**
   * 创建触发器实例
   * @param {string} type - 工作流类型
   * @returns {WorkflowTriggerStrategy} 触发器实例
   * @throws {Error} 当类型不支持时抛出错误
   */
  static createTrigger(type) {
    const StrategyClass = this.strategies[type]
    
    if (!StrategyClass) {
      const supportedTypes = Object.keys(this.strategies).join(', ')
      throw new Error(`不支持的工作流类型: ${type}。支持的类型: ${supportedTypes}`)
    }
    
    return new StrategyClass()
  }
  
  /**
   * 注册新的触发器策略
   * @param {string} type - 触发器类型
   * @param {Class} StrategyClass - 策略类（必须继承 WorkflowTriggerStrategy）
   */
  static registerStrategy(type, StrategyClass) {
    if (!StrategyClass.prototype || !StrategyClass.prototype.execute) {
      throw new Error('策略类必须继承 WorkflowTriggerStrategy 并实现 execute 方法')
    }
    
    this.strategies[type] = StrategyClass
    console.log(`✅ 已注册新的触发器策略: ${type}`)
  }
  
  /**
   * 获取支持的触发器类型列表
   * @returns {Array<string>} 支持的类型列表
   */
  static getSupportedTypes() {
    return Object.keys(this.strategies)
  }
  
  /**
   * 检查是否支持指定类型
   * @param {string} type - 要检查的类型
   * @returns {Boolean} 是否支持
   */
  static isSupported(type) {
    return type in this.strategies
  }
  
  /**
   * 获取触发器信息
   * @param {string} type - 触发器类型
   * @returns {Object|null} 触发器信息
   */
  static getTriggerInfo(type) {
    if (!this.isSupported(type)) {
      return null
    }
    
    const StrategyClass = this.strategies[type]
    const instance = new StrategyClass()
    
    return {
      type,
      name: this.getTriggerName(type),
      description: this.getTriggerDescription(type),
      requiredFields: instance.getRequiredFields ? instance.getRequiredFields() : []
    }
  }
  
  /**
   * 获取触发器显示名称
   * @param {string} type - 触发器类型
   * @returns {string} 显示名称
   */
  static getTriggerName(type) {
    const names = {
      [WORKFLOW_TYPES.SIMPLE]: '简单执行',
      [WORKFLOW_TYPES.UPLOAD_FILES]: '文件上传',
      [WORKFLOW_TYPES.FORM_INPUT]: '表单输入',
      'webhook': 'Webhook 调用',
      'batch': '批量处理'
    }
    
    return names[type] || type
  }
  
  /**
   * 获取触发器描述
   * @param {string} type - 触发器类型
   * @returns {string} 描述信息
   */
  static getTriggerDescription(type) {
    const descriptions = {
      [WORKFLOW_TYPES.SIMPLE]: '无需参数，直接执行工作流',
      [WORKFLOW_TYPES.UPLOAD_FILES]: '上传文件作为工作流输入',
      [WORKFLOW_TYPES.FORM_INPUT]: '通过表单输入参数执行工作流',
      'webhook': '直接调用 Webhook URL 执行工作流',
      'batch': '批量处理多条数据'
    }
    
    return descriptions[type] || '自定义触发器'
  }
}