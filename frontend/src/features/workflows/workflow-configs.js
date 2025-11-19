/**
 * 工作流配置中心
 * 配置已迁移到数据库，通过 API 动态加载
 * 此文件保留类型定义和默认配置作为后备
 */

import { getWorkflowConfigs } from './api'

export const WORKFLOW_TYPES = {
  SIMPLE: 'simple',           // 简单执行，无参数
  UPLOAD_FILES: 'upload_files', // 上传文件
  FORM_INPUT: 'form_input',     // 表单输入
  WEBHOOK: 'webhook',          // Webhook 调用
  BATCH: 'batch'               // 批量处理
}

// 默认配置（仅作为后备）
const defaultConfigs = {
  // GitHub Trending 工作流
  'My workflow 2': {
    type: WORKFLOW_TYPES.SIMPLE,
    description: '获取 GitHub 本周热门项目',
    executeLabel: '获取热门项目',
    resultType: 'html' // 结果展示类型
  },

  // Azure 账单对比工作流
  'Azure账单对比': {
    type: WORKFLOW_TYPES.UPLOAD_FILES,
    description: '上传本月和上月的 Azure 账单文件（CSV/Excel）进行对比分析',
    executeLabel: '开始对比',
    resultType: 'download', // 下载文件
    fields: [
      {
        name: 'current_month',
        label: '本月账单',
        type: 'file',
        accept: '.csv,.xlsx,.xls',
        required: true,
        description: '请上传本月的 Azure 账单文件（CSV/Excel）'
      },
      {
        name: 'last_month',
        label: '上月账单',
        type: 'file',
        accept: '.csv,.xlsx,.xls',
        required: true,
        description: '请上传上月的 Azure 账单文件（CSV/Excel）'
      }
    ]
  },

  // Azure 账单对比工作流（修正版）
  'Azure账单对比(修正版)': {
    type: WORKFLOW_TYPES.UPLOAD_FILES,
    description: '上传本月和上月的 Azure 账单文件（CSV/Excel）进行对比分析',
    executeLabel: '开始对比',
    resultType: 'download', // 下载文件
    fields: [
      {
        name: 'current_month',
        label: '本月账单',
        type: 'file',
        accept: '.csv,.xlsx,.xls',
        required: true,
        description: '请上传本月的 Azure 账单文件（CSV/Excel）'
      },
      {
        name: 'last_month',
        label: '上月账单',
        type: 'file',
        accept: '.csv,.xlsx,.xls',
        required: true,
        description: '请上传上月的 Azure 账单文件（CSV/Excel）'
      }
    ]
  },

  'Azure账单对比(CSV修复版)': {
    type: WORKFLOW_TYPES.UPLOAD_FILES,
    description: '上传本月和上月的 Azure 账单文件（CSV/Excel），导出对比结果为 CSV',
    executeLabel: '开始对比',
    resultType: 'download',
    fields: [
      {
        name: 'current_month',
        label: '本月账单',
        type: 'file',
        accept: '.csv,.xlsx,.xls',
        required: true,
        description: '请上传本月的 Azure 账单文件（CSV/Excel）'
      },
      {
        name: 'last_month',
        label: '上月账单',
        type: 'file',
        accept: '.csv,.xlsx,.xls',
        required: true,
        description: '请上传上月的 Azure 账单文件（CSV/Excel）'
      }
    ]
  },

  // 示例：表单输入类型的工作流
  // '数据查询': {
  //   type: WORKFLOW_TYPES.FORM_INPUT,
  //   description: '根据条件查询数据',
  //   executeLabel: '查询',
  //   resultType: 'json',
  //   fields: [
  //     {
  //       name: 'startDate',
  //       label: '开始日期',
  //       type: 'date',
  //       required: true
  //     },
  //     {
  //       name: 'endDate',
  //       label: '结束日期',
  //       type: 'date',
  //       required: true
  //     },
  //     {
  //       name: 'keyword',
  //       label: '关键词',
  //       type: 'text',
  //       required: false,
  //       placeholder: '请输入搜索关键词'
  //     }
  //   ]
  // }
}

// 缓存的配置（从后端加载）
let cachedConfigs = {}
let configsLoaded = false

/**
 * 加载工作流配置从后端
 * @returns {Promise<Object>} 配置字典
 */
export async function loadWorkflowConfigs() {
  try {
    const res = await getWorkflowConfigs()
    if (res.code === 0 && res.data) {
      // 转换为字典格式
      cachedConfigs = {}
      for (const config of res.data) {
        cachedConfigs[config.workflowName] = {
          type: config.type,
          description: config.description,
          executeLabel: config.executeLabel,
          resultType: config.resultType,
          fields: config.fields || []
        }
      }
      configsLoaded = true
      console.log('✅ 工作流配置已从后端加载:', Object.keys(cachedConfigs))
    }
  } catch (error) {
    console.warn('⚠️ 加载工作流配置失败，使用默认配置:', error)
    cachedConfigs = defaultConfigs
    configsLoaded = true
  }
  return cachedConfigs
}

/**
 * 根据工作流获取配置
 * @param {Object} workflow - 工作流对象
 * @returns {Object|null} 工作流配置
 */
export function getWorkflowConfig(workflow) {
  if (!workflow) return null

  // 优先从缓存的配置中获取
  let config = cachedConfigs[workflow.name] || defaultConfigs[workflow.name]

  if (config) {
    return { ...config, workflow }
  }

  // 如果没有配置，返回默认简单类型
  return {
    type: WORKFLOW_TYPES.SIMPLE,
    description: workflow.name,
    executeLabel: '执行工作流',
    resultType: 'auto', // 自动检测
    workflow
  }
}

/**
 * 判断工作流是否需要参数
 * @param {Object} workflow - 工作流对象
 * @returns {Boolean}
 */
export function workflowNeedsInput(workflow) {
  const config = getWorkflowConfig(workflow)
  return config.type !== WORKFLOW_TYPES.SIMPLE
}
