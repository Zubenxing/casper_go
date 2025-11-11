<template>
  <div class="workflows-container">
    <el-card class="header-card" shadow="never">
      <div class="header-content">
        <div class="title-section">
          <h2>
            <el-icon style="vertical-align: middle; margin-right: 8px">
              <Promotion />
            </el-icon>
            AI 工作流
          </h2>
          <p class="subtitle">管理和执行自动化工作流，连接 400+ 服务</p>
        </div>
        <div class="action-section">
          <el-button
            type="primary"
            :icon="Refresh"
            @click="handleRefresh"
            :loading="loading"
          >
            刷新
          </el-button>
          <el-button
            :icon="Connection"
            @click="checkN8nHealth"
            :loading="healthChecking"
          >
            检查连接
          </el-button>
          <el-tag :type="healthStatus === 'healthy' ? 'success' : 'danger'">
            {{ healthStatus === 'healthy' ? 'n8n 已连接' : 'n8n 未连接' }}
          </el-tag>
        </div>
      </div>
    </el-card>

    <!-- 工作流列表 -->
    <el-card class="workflows-card" shadow="never" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>工作流列表 ({{ workflows.length }})</span>
          <el-input
            v-model="searchText"
            placeholder="搜索工作流..."
            :prefix-icon="Search"
            style="width: 300px"
            clearable
          />
        </div>
      </template>

      <div v-if="filteredWorkflows.length === 0" class="empty-state">
        <el-empty description="暂无工作流">
          <template #image>
            <el-icon :size="60" color="#909399">
              <DocumentCopy />
            </el-icon>
          </template>
          <el-button type="primary" @click="openN8nEditor">
            前往 n8n 创建工作流
          </el-button>
        </el-empty>
      </div>

      <div v-else class="workflows-grid">
        <el-card
          v-for="workflow in filteredWorkflows"
          :key="workflow.id"
          class="workflow-card"
          shadow="hover"
          :body-style="{ padding: '20px' }"
        >
          <div class="workflow-header">
            <div class="workflow-title">
              <el-icon :size="20" color="#409eff">
                <Operation />
              </el-icon>
              <h3>{{ workflow.name }}</h3>
            </div>
            <el-tag :type="workflow.active ? 'success' : 'info'" size="small">
              {{ workflow.active ? '已激活' : '未激活' }}
            </el-tag>
          </div>

          <div class="workflow-meta">
            <div class="meta-item">
              <el-icon><Clock /></el-icon>
              <span>更新于 {{ formatTime(workflow.updatedAt) }}</span>
            </div>
            <div class="meta-item">
              <el-icon><Grid /></el-icon>
              <span>{{ workflow.nodes?.length || 0 }} 个节点</span>
            </div>
          </div>

          <div class="workflow-tags" v-if="workflow.tags && workflow.tags.length > 0">
            <el-tag
              v-for="tag in workflow.tags"
              :key="tag.id"
              size="small"
              effect="plain"
              style="margin-right: 5px"
            >
              {{ tag.name }}
            </el-tag>
          </div>

          <div class="workflow-actions">
            <el-button
              type="primary"
              :icon="CaretRight"
              @click="showExecuteDialog(workflow)"
            >
              执行
            </el-button>
            <el-button
              :icon="View"
              @click="viewWorkflowDetail(workflow)"
            >
              详情
            </el-button>
            <el-dropdown @command="handleWorkflowAction" trigger="click">
              <el-button :icon="More" circle />
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item
                    :command="{ action: 'toggle', workflow }"
                    :icon="workflow.active ? VideoPause : VideoPlay"
                  >
                    {{ workflow.active ? '停用' : '激活' }}
                  </el-dropdown-item>
                  <el-dropdown-item
                    :command="{ action: 'edit', workflow }"
                    :icon="Edit"
                  >
                    在 n8n 中编辑
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-card>
      </div>
    </el-card>

    <!-- 执行结果展示 -->
    <el-card 
      v-if="executionResult" 
      class="result-card" 
      shadow="never"
    >
      <template #header>
        <div class="card-header">
          <span>
            <el-icon style="vertical-align: middle; margin-right: 8px">
              <CircleCheck v-if="executionResult.status === 'success'" />
              <CircleClose v-else />
            </el-icon>
            执行结果
          </span>
          <el-button 
            :icon="Close" 
            @click="executionResult = null"
            size="small"
            text
          >
            关闭
          </el-button>
        </div>
      </template>

      <div class="execution-result-content">
        <el-alert
          :type="executionResult.status === 'success' ? 'success' : 'error'"
          :closable="false"
          style="margin-bottom: 20px"
        >
          <template #title>
            <div style="font-size: 16px; font-weight: 600">
              {{ executionResult.status === 'success' ? '✅ 执行成功' : '❌ 执行失败' }}
            </div>
          </template>
          <div style="margin-top: 10px">
            <p><strong>工作流：</strong>{{ executionResult.workflowName || currentWorkflow?.name }}</p>
            <p><strong>执行ID：</strong>{{ executionResult.id }}</p>
            <p><strong>开始时间：</strong>{{ formatDateTime(executionResult.startedAt) }}</p>
            <p><strong>结束时间：</strong>{{ formatDateTime(executionResult.stoppedAt) }}</p>
          </div>
        </el-alert>

        <div v-if="executionResult.data" class="result-data-card">
          <h3 style="margin: 0 0 15px 0; font-size: 16px; color: #303133">
            <el-icon style="vertical-align: middle; margin-right: 5px">
              <Document />
            </el-icon>
            返回数据
            <el-tag 
              size="small" 
              style="margin-left: 10px"
              :type="isHTMLResult ? 'success' : 'info'"
            >
              {{ isHTMLResult ? 'HTML' : 'JSON' }}
            </el-tag>
          </h3>
          
          <!-- HTML 渲染模式 -->
          <div v-if="isHTMLResult" class="html-result-viewer">
            <el-scrollbar max-height="600px">
              <iframe 
                class="html-iframe"
                :srcdoc="extractHTML(executionResult.data)"
                frameborder="0"
                sandbox="allow-same-origin"
              ></iframe>
            </el-scrollbar>
          </div>
          
          <!-- JSON 展示模式 -->
          <el-scrollbar v-else max-height="500px">
            <div class="json-viewer">
              <pre>{{ formatJSON(executionResult.data) }}</pre>
            </div>
          </el-scrollbar>

          <div style="margin-top: 15px; text-align: right">
            <el-button 
              @click="copyResult"
              :icon="DocumentCopy"
            >
              {{ isHTMLResult ? '复制 HTML' : '复制 JSON' }}
            </el-button>
            <el-button 
              v-if="isHTMLResult"
              @click="viewHTMLInNewTab"
              :icon="View"
            >
              新窗口查看
            </el-button>
            <el-button 
              type="primary"
              @click="viewInN8n(executionResult)"
            >
              在 n8n 中查看详情
            </el-button>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 执行历史 -->
    <el-card class="executions-card" shadow="never" v-loading="executionsLoading">
      <template #header>
        <div class="card-header">
          <span>执行历史</span>
          <el-button
            :icon="Refresh"
            @click="fetchExecutions"
            size="small"
          >
            刷新
          </el-button>
        </div>
      </template>

      <el-table
        :data="executions"
        style="width: 100%"
        empty-text="暂无执行记录"
      >
        <el-table-column prop="workflowName" label="工作流" min-width="180" />
        <el-table-column label="状态" width="100">
          <template #default="scope">
            <el-tag
              :type="getStatusType(scope.row.status)"
              size="small"
            >
              {{ getStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="执行模式" width="100">
          <template #default="scope">
            <el-tag size="small" effect="plain">
              {{ scope.row.mode }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="开始时间" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.startedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="结束时间" width="180">
          <template #default="scope">
            {{ formatDateTime(scope.row.stoppedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <el-button
              type="primary"
              :icon="View"
              size="small"
              @click="viewExecutionDetail(scope.row)"
            >
              查看
            </el-button>
            <el-button
              type="danger"
              :icon="Delete"
              size="small"
              @click="deleteExecutionRecord(scope.row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 执行工作流对话框 -->
    <el-dialog
      v-model="executeDialogVisible"
      title="确认执行"
      width="500px"
      :close-on-click-modal="false"
    >
      <div v-if="currentWorkflow" style="text-align: center; padding: 20px 0">
        <el-icon :size="60" color="#409eff" style="margin-bottom: 20px">
          <CaretRight />
        </el-icon>
        
        <h2 style="margin: 0 0 10px 0; font-size: 20px; color: #303133">
          {{ currentWorkflow.name }}
        </h2>
        
        <p style="color: #909399; margin-bottom: 20px">
          <el-icon style="vertical-align: middle"><Grid /></el-icon>
          包含 {{ currentWorkflow.nodes?.length || 0 }} 个节点
        </p>

        <el-divider style="margin: 20px 0" />

        <!-- 高级选项（折叠） -->
        <el-collapse v-model="showAdvancedOptions" style="text-align: left">
          <el-collapse-item name="advanced">
            <template #title>
              <span style="color: #909399; font-size: 14px">
                <el-icon style="vertical-align: middle"><Setting /></el-icon>
                高级选项（可选）
              </span>
            </template>
            <el-form label-position="top" size="small">
              <el-form-item label="自定义参数 (JSON)">
                <el-input
                  v-model="executeParams"
                  type="textarea"
                  :rows="6"
                  placeholder='{"key": "value"}'
                  style="font-family: monospace; font-size: 12px"
                />
                <div style="margin-top: 5px; color: #909399; font-size: 12px">
                  💡 如果工作流需要额外参数，可以在这里输入
                </div>
              </el-form-item>
            </el-form>
          </el-collapse-item>
        </el-collapse>

        <el-alert
          type="info"
          :closable="false"
          show-icon
          style="margin-top: 20px"
        >
          <template #title>
            <span style="font-size: 13px">
              工作流将使用 n8n 中配置的默认参数执行
            </span>
          </template>
        </el-alert>
      </div>

      <template #footer>
        <div style="display: flex; justify-content: center; gap: 10px">
          <el-button @click="executeDialogVisible = false" size="large">
            取消
          </el-button>
          <el-button
            type="primary"
            :icon="CaretRight"
            @click="executeWorkflowNow"
            :loading="executing"
            size="large"
          >
            {{ executing ? '执行中...' : '立即执行' }}
          </el-button>
        </div>
      </template>
    </el-dialog>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Promotion,
  Refresh,
  Connection,
  Search,
  DocumentCopy,
  Operation,
  Clock,
  Grid,
  CaretRight,
  View,
  More,
  VideoPause,
  VideoPlay,
  Edit,
  Delete,
  CircleCheck,
  CircleClose,
  Close,
  Document,
  Setting
} from '@element-plus/icons-vue'
import {
  getWorkflows,
  getWorkflow,
  executeWorkflow,
  getExecutions,
  getExecution,
  deleteExecution,
  activateWorkflow,
  deactivateWorkflow,
  checkHealth
} from './api'

// 数据
const workflows = ref([])
const executions = ref([])
const searchText = ref('')
const loading = ref(false)
const executionsLoading = ref(false)
const healthChecking = ref(false)
const healthStatus = ref('unknown')
const executeDialogVisible = ref(false)
const currentWorkflow = ref(null)
const executeParams = ref('')
const executing = ref(false)
const executionResult = ref(null)
const showAdvancedOptions = ref([]) // 高级选项折叠状态

// 计算属性
const filteredWorkflows = computed(() => {
  if (!searchText.value) return workflows.value
  const search = searchText.value.toLowerCase()
  return workflows.value.filter(w =>
    w.name.toLowerCase().includes(search) ||
    (w.tags && w.tags.some(t => t.name.toLowerCase().includes(search)))
  )
})

// 方法
const fetchWorkflows = async () => {
  loading.value = true
  try {
    const res = await getWorkflows()
    // 后端返回格式：{data: {data: [...], count: n}}
    const workflowList = res.data?.data || []
    console.log('获取到的工作流列表:', workflowList)
    console.log('工作流数量:', workflowList.length)
    
    // 过滤掉未激活的工作流（可选）
    workflows.value = workflowList
  } catch (error) {
    console.error('获取工作流失败:', error)
    ElMessage.error('获取工作流失败: ' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 强制刷新所有数据
const handleRefresh = async () => {
  await Promise.all([
    fetchWorkflows(),
    fetchExecutions()
  ])
  ElMessage.success('刷新成功')
}

const fetchExecutions = async () => {
  executionsLoading.value = true
  try {
    const res = await getExecutions(20)
    // 后端返回格式：{data: {data: [...], count: n}}
    const executionList = res.data?.data || []
    console.log('获取到的执行历史:', executionList)
    
    // 确保每个执行记录都有工作流名称
    executions.value = executionList.map(exec => ({
      ...exec,
      workflowName: exec.workflowName || exec.workflowId || '未知工作流'
    }))
  } catch (error) {
    ElMessage.error('获取执行历史失败: ' + (error.message || '未知错误'))
  } finally {
    executionsLoading.value = false
  }
}

const checkN8nHealth = async () => {
  healthChecking.value = true
  try {
    await checkHealth()
    healthStatus.value = 'healthy'
    ElMessage.success('n8n 服务运行正常')
  } catch (error) {
    healthStatus.value = 'unhealthy'
    ElMessage.error('n8n 服务不可用')
  } finally {
    healthChecking.value = false
  }
}

const showExecuteDialog = (workflow) => {
  currentWorkflow.value = workflow
  executeParams.value = ''
  showAdvancedOptions.value = [] // 默认折叠高级选项
  executeDialogVisible.value = true
}

const executeWorkflowNow = async () => {
  executing.value = true
  try {
    // 解析参数
    let params = {}
    if (executeParams.value.trim()) {
      try {
        params = JSON.parse(executeParams.value)
      } catch (e) {
        ElMessage.error('参数格式错误，请输入有效的 JSON')
        executing.value = false
        return
      }
    }

    console.log('准备执行工作流:', {
      id: currentWorkflow.value.id,
      name: currentWorkflow.value.name,
      params: params
    })

    const res = await executeWorkflow(currentWorkflow.value.id, params)
    
    // 调试：打印返回数据
    console.log('工作流执行返回数据:', res)
    
    // 检查响应结构
    if (res.code === 500 || res.code !== 0) {
      // 后端返回错误
      const errorMsg = res.message || res.msg || '执行失败'
      console.error('后端返回错误:', errorMsg)
      ElMessage.error(`执行失败: ${errorMsg}`)
      executing.value = false
      executeDialogVisible.value = false
      return
    }
    
    console.log('res.data:', res.data)
    
    // 后端返回格式：{code: 0, data: ExecutionData}
    if (!res.data) {
      console.error('返回数据为空')
      ElMessage.error('执行失败: 返回数据为空')
      executing.value = false
      executeDialogVisible.value = false
      return
    }
    
    executionResult.value = res.data
    executeDialogVisible.value = false
    
    console.log('executionResult 已设置:', executionResult.value)
    
    ElMessage.success('🎉 工作流执行成功！')
    
    // 刷新工作流列表和执行历史
    fetchWorkflows()
    fetchExecutions()
    
    // 滚动到结果区域
    setTimeout(() => {
      const resultCard = document.querySelector('.result-card')
      console.log('查找结果卡片:', resultCard)
      if (resultCard) {
        resultCard.scrollIntoView({ behavior: 'smooth', block: 'start' })
      } else {
        console.warn('未找到 .result-card 元素，executionResult:', executionResult.value)
      }
    }, 200)
  } catch (error) {
    console.error('执行异常:', error)
    ElMessage.error('执行失败: ' + (error.message || '未知错误'))
    executeDialogVisible.value = false
  } finally {
    executing.value = false
  }
}

const viewWorkflowDetail = async (workflow) => {
  try {
    const res = await getWorkflow(workflow.id)
    // 后端返回格式：{data: Workflow}
    ElMessageBox.alert(
      `<pre>${JSON.stringify(res.data, null, 2)}</pre>`,
      '工作流详情',
      {
        dangerouslyUseHTMLString: true,
        confirmButtonText: '关闭'
      }
    )
  } catch (error) {
    ElMessage.error('获取详情失败')
  }
}

const viewExecutionDetail = async (execution) => {
  try {
    const res = await getExecution(execution.id)
    // 后端返回格式：{data: ExecutionData}
    executionResult.value = res.data
    
    // 滚动到结果区域
    setTimeout(() => {
      const resultCard = document.querySelector('.result-card')
      if (resultCard) {
        resultCard.scrollIntoView({ behavior: 'smooth', block: 'start' })
      }
    }, 100)
  } catch (error) {
    ElMessage.error('获取执行详情失败')
  }
}

const deleteExecutionRecord = async (execution) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除执行记录"${execution.workflowName}"吗？`,
      '确认删除',
      {
        type: 'warning'
      }
    )
    
    await deleteExecution(execution.id)
    ElMessage.success('删除成功')
    fetchExecutions()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleWorkflowAction = async ({ action, workflow }) => {
  if (action === 'toggle') {
    try {
      if (workflow.active) {
        await deactivateWorkflow(workflow.id)
        ElMessage.success('工作流已停用')
      } else {
        await activateWorkflow(workflow.id)
        ElMessage.success('工作流已激活')
      }
      fetchWorkflows()
    } catch (error) {
      ElMessage.error('操作失败')
    }
  } else if (action === 'edit') {
    openN8nEditor(workflow.id)
  }
}

const openN8nEditor = (workflowId = null) => {
  const baseUrl = import.meta.env.VITE_N8N_URL || 'http://localhost:5678'
  const url = workflowId ? `${baseUrl}/workflow/${workflowId}` : baseUrl
  window.open(url, '_blank')
}

const viewInN8n = (execution) => {
  if (execution && execution.id) {
    const baseUrl = import.meta.env.VITE_N8N_URL || 'http://localhost:5678'
    window.open(`${baseUrl}/execution/${execution.id}`, '_blank')
  }
}

const formatTime = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  const now = new Date()
  const diff = now - date
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)
  
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes} 分钟前`
  if (hours < 24) return `${hours} 小时前`
  if (days < 7) return `${days} 天前`
  return date.toLocaleDateString()
}

const formatDateTime = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const getStatusType = (status) => {
  const map = {
    'success': 'success',
    'error': 'danger',
    'waiting': 'warning',
    'running': 'info'
  }
  return map[status] || 'info'
}

const getStatusText = (status) => {
  const map = {
    'success': '成功',
    'error': '失败',
    'waiting': '等待中',
    'running': '运行中'
  }
  return map[status] || status
}

// 格式化 JSON 数据
const formatJSON = (data) => {
  try {
    return JSON.stringify(data, null, 2)
  } catch (e) {
    return String(data)
  }
}

// 判断是否为 HTML 结果
const isHTMLResult = computed(() => {
  if (!executionResult.value?.data) return false
  const html = extractHTML(executionResult.value.data)
  return html && (html.includes('<html') || html.includes('<div') || html.includes('<table'))
})

// 从 n8n 返回数据中提取 HTML 内容
const extractHTML = (data) => {
  try {
    console.log('开始提取 HTML，数据结构:', data)
    
    // 处理不同的 n8n 返回格式
    if (typeof data === 'string') {
      return data
    }
    
    // 1. 检查 Webhook 返回格式：data.result
    if (data.result) {
      console.log('检测到 Webhook 返回格式: data.result')
      let result = data.result
      
      // 如果是字符串，尝试解析
      if (typeof result === 'string') {
        // 尝试解析为 JSON（n8n 可能返回 JSON 字符串）
        try {
          const parsed = JSON.parse(result)
          console.log('解析 JSON 字符串成功:', parsed)
          
          // 检查是否有 404 错误
          if (parsed.code === 404) {
            console.error('Webhook 返回 404 错误:', parsed.message)
            return `错误：${parsed.message}\n\n提示：请确保工作流已激活，并使用生产 Webhook URL。`
          }
          
          // 如果解析后的对象还有 result 字段
          if (parsed.result) {
            result = parsed.result
          }
        } catch (e) {
          // 不是 JSON，直接返回字符串
          console.log('不是 JSON 字符串，直接返回')
          return result
        }
      }
      
      return typeof result === 'string' ? result : JSON.stringify(result, null, 2)
    }
    
    // 2. 检查 resultData.runData 结构
    if (data.resultData?.runData) {
      const runData = data.resultData.runData
      // 遍历所有节点
      for (const nodeName in runData) {
        const nodeData = runData[nodeName]
        if (Array.isArray(nodeData)) {
          for (const execution of nodeData) {
            if (execution.data?.main?.[0]) {
              for (const item of execution.data.main[0]) {
                // 检查 json 字段中的 html 或 output
                if (item.json?.html) return item.json.html
                if (item.json?.output) return item.json.output
                if (item.json?.result) return item.json.result
                // 检查 binary 数据
                if (item.binary?.data?.data) {
                  try {
                    return atob(item.binary.data.data)
                  } catch (e) {
                    // ignore
                  }
                }
              }
            }
          }
        }
      }
    }
    
    // 如果都没找到，返回整个 JSON 字符串
    console.log('未找到 HTML，返回 JSON')
    return JSON.stringify(data, null, 2)
  } catch (e) {
    console.error('提取 HTML 失败:', e)
    return JSON.stringify(data, null, 2)
  }
}

// 复制结果到剪贴板
const copyResult = async () => {
  try {
    let text
    if (isHTMLResult.value) {
      text = extractHTML(executionResult.value.data)
    } else {
      text = formatJSON(executionResult.value.data)
    }
    await navigator.clipboard.writeText(text)
    ElMessage.success('✅ 结果已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

// 在新窗口中查看 HTML
const viewHTMLInNewTab = () => {
  const html = extractHTML(executionResult.value.data)
  const newWindow = window.open('', '_blank')
  if (newWindow) {
    newWindow.document.write(html)
    newWindow.document.close()
  }
}

// 生命周期
onMounted(() => {
  fetchWorkflows()
  fetchExecutions()
  checkN8nHealth()
})
</script>

<style scoped>
.workflows-container {
  padding: 20px;
}

.header-card {
  margin-bottom: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-section h2 {
  margin: 0 0 8px 0;
  font-size: 24px;
  color: #303133;
  display: flex;
  align-items: center;
}

.subtitle {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

.action-section {
  display: flex;
  gap: 10px;
  align-items: center;
}

.workflows-card,
.result-card,
.executions-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.empty-state {
  text-align: center;
  padding: 60px 0;
}

.workflows-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.workflow-card {
  transition: transform 0.2s;
}

.workflow-card:hover {
  transform: translateY(-2px);
}

.workflow-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.workflow-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.workflow-title h3 {
  margin: 0;
  font-size: 16px;
  color: #303133;
}

.workflow-meta {
  margin-bottom: 12px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #909399;
  font-size: 13px;
  margin-bottom: 6px;
}

.workflow-tags {
  margin-bottom: 16px;
  min-height: 24px;
}

.workflow-actions {
  display: flex;
  gap: 8px;
  padding-top: 12px;
  border-top: 1px solid #ebeef5;
}

.execution-result {
  padding: 10px 0;
}

.result-data {
  margin-top: 20px;
}

.result-data h4 {
  margin-bottom: 10px;
}

.result-data pre {
  background: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
  font-size: 13px;
  color: #303133;
  overflow-x: auto;
}

/* 执行结果卡片样式 */
.result-card {
  animation: slideDown 0.3s ease-out;
  border-left: 4px solid #67c23a;
  max-width: 100%;
  overflow: hidden;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.execution-result-content {
  padding: 10px 0;
}

.execution-result-content p {
  margin: 8px 0;
  font-size: 14px;
  color: #606266;
}

.execution-result-content strong {
  color: #303133;
  font-weight: 600;
}

.result-data-card {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 20px;
  margin-top: 10px;
}

.json-viewer {
  background: #ffffff;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 0;
}

.json-viewer pre {
  margin: 0;
  padding: 20px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 13px;
  line-height: 1.6;
  color: #2c3e50;
  overflow-x: auto;
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* 滚动条美化 */
.json-viewer :deep(.el-scrollbar__wrap) {
  overflow-x: auto;
}

/* HTML 结果查看器 */
.html-result-viewer {
  background: #ffffff;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
  max-width: 100%;
}

.html-iframe {
  width: 100%;
  min-height: 600px;
  border: none;
  background: white;
  display: block;
}

.html-content {
  padding: 20px;
  min-height: 100px;
  max-width: 100%;
  overflow-x: auto;
  word-wrap: break-word;
}

/* 让渲染的 HTML 内容样式正常显示 */
.html-content :deep(*) {
  max-width: 100%;
  box-sizing: border-box;
}

.html-content :deep(img) {
  max-width: 100%;
  height: auto;
}

.html-content :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 10px 0;
  table-layout: auto;
  overflow-x: auto;
  display: block;
}

.html-content :deep(.container) {
  max-width: 100% !important;
  margin: 0 !important;
}

.html-content :deep(table td),
.html-content :deep(table th) {
  border: 1px solid #ddd;
  padding: 8px;
}

.html-content :deep(table th) {
  background-color: #f5f7fa;
  font-weight: 600;
}

.html-content :deep(pre) {
  background: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  overflow-x: auto;
}

.html-content :deep(code) {
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'Monaco', 'Menlo', monospace;
}
</style>



