<template>
  <div class="issue-list">
    <!-- 顶部操作栏 -->
    <div class="action-bar">
      <div class="action-buttons">
        <el-select
          v-model="filterStatus"
          placeholder="筛选状态"
          clearable
          @change="handleFilterChange"
          class="status-filter"
          size="large"
        >
          <el-option label="全部状态" value="" />
          <el-option label="待处理" value="open" />
          <el-option label="处理中" value="in_progress" />
          <el-option label="已解决" value="resolved" />
          <el-option label="已关闭" value="closed" />
        </el-select>

        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="-"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          @change="handleDateChange"
          size="large"
          class="date-filter"
          clearable
        />

        <el-button type="warning" size="large" @click="handleAdd" class="add-btn">
          <el-icon><Plus /></el-icon>
          <span>新建问题</span>
        </el-button>
      </div>
    </div>

    <!-- 表格 -->
    <el-table
      v-loading="loading"
      :data="issueList"
      stripe
      class="issue-table"
      :header-cell-style="{ background: '#fff5f5', color: '#c2410c', fontWeight: '600' }"
    >
      <el-table-column type="index" label="#" width="60" align="center" />

      <el-table-column label="问题" min-width="400">
        <template #default="{ row }">
          <div class="issue-cell-wrapper">
            <!-- 左侧：标题和描述区 -->
            <div class="issue-main-content">
              <div class="issue-cell">
                <div class="issue-title">{{ row.title }}</div>
                <div v-if="row.description" class="issue-desc">{{ row.description }}</div>
              </div>
              
              <!-- 解决方案和标签在下方 -->
              <div v-if="row.solution" class="issue-solution">
                <el-icon class="solution-icon"><Checked /></el-icon>
                <span>{{ row.solution }}</span>
              </div>
              <div v-if="row.tags" class="issue-tags">
                <el-tag
                  v-for="(tag, index) in row.tags.split(',')"
                  :key="index"
                  size="small"
                  class="issue-tag"
                >
                  {{ tag }}
                </el-tag>
              </div>
            </div>
            
            <!-- 右侧：图片区（覆盖标题+描述高度） -->
            <div v-if="row.images && parseImages(row.images).length > 0" class="issue-images-area">
              <el-image
                v-for="(img, index) in parseImages(row.images).slice(0, 3)"
                :key="index"
                :src="getImageUrl(img)"
                :preview-src-list="parseImages(row.images).map(i => getImageUrl(i))"
                :initial-index="index"
                :z-index="9999"
                :preview-teleported="true"
                :hide-on-click-modal="true"
                fit="cover"
                class="issue-image"
                lazy
              >
                <template #error>
                  <div class="image-error">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
            </div>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="状态" width="140" align="center" sortable prop="status">
        <template #default="{ row }">
          <el-dropdown @command="(cmd) => handleStatusChange(row, cmd)" trigger="click">
            <el-tag
              :type="getStatusType(row.status)"
              effect="dark"
              size="large"
              class="status-tag clickable"
            >
              {{ getStatusText(row.status) }}
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-tag>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="open">待处理</el-dropdown-item>
                <el-dropdown-item command="in_progress">处理中</el-dropdown-item>
                <el-dropdown-item command="resolved">已解决</el-dropdown-item>
                <el-dropdown-item command="closed">已关闭</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>

      <el-table-column label="严重程度" width="120" align="center" sortable prop="severity">
        <template #default="{ row }">
          <el-dropdown @command="(cmd) => handleSeverityChange(row, cmd)" trigger="click">
            <el-tag
              :type="getSeverityType(row.severity)"
              effect="plain"
              size="large"
              class="severity-tag clickable"
            >
              {{ getSeverityText(row.severity) }}
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-tag>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="low">低</el-dropdown-item>
                <el-dropdown-item command="medium">中</el-dropdown-item>
                <el-dropdown-item command="high">高</el-dropdown-item>
                <el-dropdown-item command="critical">严重</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>

      <el-table-column label="创建时间" width="140" align="center">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>

      <el-table-column label="解决时间" width="140" align="center">
        <template #default="{ row }">
          <span v-if="row.resolved_at" class="resolved-time">
            <el-icon><CircleCheck /></el-icon>
            {{ formatDate(row.resolved_at) }}
          </span>
          <span v-else class="text-muted">-</span>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="120" align="center" fixed="right">
        <template #default="{ row }">
          <el-tooltip content="查看详情" placement="top">
            <el-button
              type="primary"
              link
              size="small"
              @click="handleViewDetail(row)"
            >
              <el-icon><Document /></el-icon>
            </el-button>
          </el-tooltip>
          <el-tooltip content="编辑" placement="top">
            <el-button
              type="warning"
              link
              size="small"
              @click="handleEdit(row)"
            >
              <el-icon><Edit /></el-icon>
            </el-button>
          </el-tooltip>
          <el-tooltip content="删除" placement="top">
            <el-button
              type="danger"
              link
              size="small"
              @click="handleDelete(row)"
            >
              <el-icon><Delete /></el-icon>
            </el-button>
          </el-tooltip>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div v-if="total > 0" class="pagination">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchData"
        @current-change="fetchData"
      />
    </div>

    <!-- 新建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      class="issue-dialog"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="90px"
      >
        <el-form-item label="问题标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入问题标题" />
        </el-form-item>

        <el-form-item label="问题描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="4"
            placeholder="请详细描述遇到的问题"
          />
        </el-form-item>

        <el-form-item label="问题截图">
          <div class="image-upload-area">
            <el-upload
              v-model:file-list="imageFileList"
              :action="uploadUrl"
              :headers="uploadHeaders"
              :on-success="handleImageSuccess"
              :on-remove="handleImageRemove"
              :before-upload="beforeImageUpload"
              :limit="3"
              :on-exceed="handleImageExceed"
              list-type="picture-card"
              accept="image/*"
            >
              <el-icon><Plus /></el-icon>
            </el-upload>
            <div class="upload-tip">最多上传3张截图，单张不超过5MB</div>
          </div>
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select v-model="form.status" placeholder="请选择状态" style="width: 100%">
                <el-option label="待处理" value="open" />
                <el-option label="处理中" value="in_progress" />
                <el-option label="已解决" value="resolved" />
                <el-option label="已关闭" value="closed" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="严重程度" prop="severity">
              <el-select v-model="form.severity" placeholder="请选择严重程度" style="width: 100%">
                <el-option label="低" value="low" />
                <el-option label="中" value="medium" />
                <el-option label="高" value="high" />
                <el-option label="紧急" value="critical" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="解决方案">
          <el-input
            v-model="form.solution"
            type="textarea"
            :rows="4"
            placeholder="请输入解决方案"
          />
        </el-form-item>

        <el-form-item label="标签">
          <el-input v-model="form.tags" placeholder="多个标签用逗号分隔" />
        </el-form-item>

        <el-form-item label="详细内容">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="6"
            placeholder="请输入问题的详细内容、复现步骤、解决过程等（支持 Markdown 格式）"
          />
          <div class="content-tip">提示：可以在这里详细记录问题分析、解决过程、注意事项等</div>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false" size="large">取消</el-button>
        <el-button type="warning" @click="handleSubmit" :loading="submitting" size="large">
          确定
        </el-button>
      </template>
    </el-dialog>

    <!-- 问题详情编辑器 -->
    <IssueDetailEditor
      v-model="detailEditorVisible"
      :issue="currentIssue || {}"
      @edit="handleEditFromDetail"
      @refresh="fetchData"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, Edit, Delete, Warning, Loading, CircleCheck, 
  DocumentChecked, Checked, ArrowDown, Document, Picture
} from '@element-plus/icons-vue'
import {
  getWorkIssueList, createWorkIssue, updateWorkIssue, 
  deleteWorkIssue, getWorkIssueStats
} from './api'
import IssueDetailEditor from './IssueDetailEditor.vue'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('新建问题')
const filterStatus = ref('')
const dateRange = ref(null)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const issueList = ref([])
const formRef = ref(null)
const currentEditId = ref(null)

const stats = ref({
  total: 0,
  open: 0,
  in_progress: 0,
  resolved: 0
})

const form = reactive({
  title: '',
  description: '',
  status: 'open',
  severity: 'medium',
  solution: '',
  tags: '',
  images: [],
  content: ''
})

const rules = {
  title: [{ required: true, message: '请输入问题标题', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }],
  severity: [{ required: true, message: '请选择严重程度', trigger: 'change' }]
}

// 图片上传相关
const imageFileList = ref([])
const uploadUrl = 'http://localhost:8080/api/work-issues/upload'
// 动态获取 token，避免过期问题
const uploadHeaders = computed(() => ({
  'Authorization': `Bearer ${localStorage.getItem('token')}`
}))

// 详情编辑器
const detailEditorVisible = ref(false)
const currentIssue = ref(null)

// 获取统计数据
const fetchStats = async () => {
  try {
    const statsRes = await getWorkIssueStats()
    stats.value = statsRes.data
  } catch (error) {
    console.error('获取统计失败:', error)
  }
}

// 排序规则
const sortIssueList = (list) => {
  // 状态优先级：处理中 > 待处理 > 已解决 > 已关闭
  const statusOrder = {
    'in_progress': 1,
    'open': 2,
    'resolved': 3,
    'closed': 4
  }
  
  // 严重程度排序：严重 > 高 > 中 > 低
  const severityOrder = {
    'critical': 1,
    'high': 2,
    'medium': 3,
    'low': 4
  }
  
  return list.sort((a, b) => {
    // 首先按状态排序
    const statusDiff = (statusOrder[a.status] || 999) - (statusOrder[b.status] || 999)
    if (statusDiff !== 0) return statusDiff
    
    // 状态相同时按严重程度排序
    const severityDiff = (severityOrder[a.severity] || 999) - (severityOrder[b.severity] || 999)
    if (severityDiff !== 0) return severityDiff
    
    // 都相同时按创建时间倒序（新的在前）
    return new Date(b.created_at) - new Date(a.created_at)
  })
}

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const [listRes, statsRes] = await Promise.all([
      getWorkIssueList({
        page: page.value,
        page_size: pageSize.value,
        status: filterStatus.value
      }),
      getWorkIssueStats()
    ])
    issueList.value = sortIssueList(listRes.data.list || [])
    total.value = listRes.data.total || 0
    stats.value = statsRes.data
  } catch (error) {
    ElMessage.error('获取数据失败')
  } finally {
    loading.value = false
  }
}

// 状态相关
const getStatusType = (status) => {
  const map = {
    open: 'danger',
    in_progress: 'warning',
    resolved: 'success',
    closed: 'info'
  }
  return map[status] || 'info'
}

const getStatusText = (status) => {
  const map = {
    open: '待处理',
    in_progress: '处理中',
    resolved: '已解决',
    closed: '已关闭'
  }
  return map[status] || status
}

const getSeverityType = (severity) => {
  const map = {
    low: 'info',
    medium: 'warning',
    high: 'warning',
    critical: 'danger'
  }
  return map[severity] || 'info'
}

const getSeverityText = (severity) => {
  const map = {
    low: '低',
    medium: '中',
    high: '高',
    critical: '紧急'
  }
  return map[severity] || severity
}

// 快速更改状态
const handleStatusChange = async (row, newStatus) => {
  if (row.status === newStatus) return
  
  try {
    await updateWorkIssue(row.id, { status: newStatus })
    ElMessage.success('状态已更新')
    row.status = newStatus
    // 刷新统计数据并重新排序
    await fetchStats()
    issueList.value = sortIssueList([...issueList.value])
  } catch (error) {
    ElMessage.error(error.message || '状态更新失败')
  }
}

// 快速更改严重程度
const handleSeverityChange = async (row, newSeverity) => {
  if (row.severity === newSeverity) return
  
  try {
    await updateWorkIssue(row.id, { severity: newSeverity })
    ElMessage.success('严重程度已更新')
    row.severity = newSeverity
    // 重新排序列表
    issueList.value = sortIssueList([...issueList.value])
  } catch (error) {
    ElMessage.error(error.message || '严重程度更新失败')
  }
}

// 日期格式化
const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return d.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

// 解析图片JSON
const parseImages = (images) => {
  if (!images) return []
  try {
    return typeof images === 'string' ? JSON.parse(images) : images
  } catch {
    return []
  }
}

// 获取图片URL - 将数据库中的文件名转换为完整URL
const getImageUrl = (path) => {
  // 如果路径已经是完整URL，直接返回
  if (path.startsWith('http')) return path
  // 如果是旧格式的相对路径（/api/files/xxx），加上域名
  if (path.startsWith('/api/')) return `http://localhost:8080${path}`
  // 如果只是文件名（1001_xxx.png），拼接完整路径
  return `http://localhost:8080/api/files/work-issues/${path}`
}

// 图片上传成功（el-upload 的 on-success 回调：(response, uploadFile, uploadFiles) => void）
const handleImageSuccess = (response, uploadFile, uploadFiles) => {
  console.log('==================== [图片上传] ====================')
  console.log('[图片上传] 完整响应对象:', response)
  console.log('[图片上传] response.code 值:', response.code)
  console.log('[图片上传] response.code 类型:', typeof response.code)
  console.log('[图片上传] response.code === 0:', response.code === 0)
  console.log('[图片上传] response.code == 0:', response.code == 0)
  console.log('[图片上传] 当前 form.images:', JSON.stringify(form.images))
  
  // 后端 response.Success 返回的 code 是 0，使用宽松比较
  if (response.code == 0 || response.code === 0) {
    // 后端返回文件名，前端存储文件名到数组
    const filename = response.data.filename
    console.log('[图片上传] ✅ 提取的文件名:', filename)
    
    // 添加到 form.images 数组
    form.images.push(filename)
    console.log('[图片上传] ✅ 添加后的 form.images:', JSON.stringify(form.images))
    
    ElMessage.success('图片上传成功')
  } else {
    console.error('[图片上传] ❌ 条件判断失败！')
    console.error('[图片上传] response.code:', response.code, '(type:', typeof response.code, ')')
    ElMessage.error(response.message || '图片上传失败')
  }
  console.log('====================================================')
}

// 图片移除（el-upload 的 on-remove 回调：(uploadFile, uploadFiles) => void）
const handleImageRemove = (uploadFile, uploadFiles) => {
  console.log('==================== [图片移除] ====================')
  console.log('[图片移除] 移除的文件:', uploadFile)
  console.log('[图片移除] 剩余文件:', uploadFiles)
  console.log('[图片移除] 移除前 form.images:', JSON.stringify(form.images))
  
  // 找到要移除的文件在 imageFileList 中的索引
  const index = imageFileList.value.findIndex(item => item.uid === uploadFile.uid)
  console.log('[图片移除] 找到的索引:', index)
  
  if (index !== -1 && form.images[index]) {
    // 移除 form.images 中对应的文件名
    form.images.splice(index, 1)
    console.log('[图片移除] 移除后 form.images:', JSON.stringify(form.images))
  }
  console.log('====================================================')
}

// 图片上传前检查
const beforeImageUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt5M = file.size / 1024 / 1024 < 5

  if (!isImage) {
    ElMessage.error('只能上传图片文件!')
    return false
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB!')
    return false
  }
  return true
}

// 超出上传数量限制
const handleImageExceed = () => {
  ElMessage.warning('最多只能上传 3 张图片!')
}

// 查看详情
const handleViewDetail = (row) => {
  currentIssue.value = { ...row }
  detailEditorVisible.value = true
}

// 从详情编辑器编辑问题
const handleEditFromDetail = (issue) => {
  handleEdit(issue)
}

// 日期筛选变更
const handleDateChange = () => {
  page.value = 1
  fetchData()
}

// 状态筛选变更
const handleFilterChange = () => {
  page.value = 1
  fetchData()
}

// 新建问题
const handleAdd = () => {
  dialogTitle.value = '新建问题'
  currentEditId.value = null
  resetForm()
  dialogVisible.value = true
}

// 编辑问题
const handleEdit = (issue) => {
  dialogTitle.value = '编辑问题'
  currentEditId.value = issue.id
  
  // 解析图片数据
  const images = parseImages(issue.images)
  form.images = images
  
  // 构建图片文件列表用于显示
  imageFileList.value = images.map((url, index) => ({
    uid: Date.now() + index,
    name: `image-${index + 1}`,
    url: getImageUrl(url)
  }))
  
  Object.assign(form, {
    title: issue.title,
    description: issue.description,
    status: issue.status,
    severity: issue.severity,
    solution: issue.solution || '',
    tags: issue.tags || '',
    content: issue.content || ''
  })
  dialogVisible.value = true
}

// 删除问题
const handleDelete = async (issue) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除问题"${issue.title}"吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await deleteWorkIssue(issue.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    console.log('[提交] form.images:', form.images)
    console.log('[提交] imageFileList:', imageFileList.value)

    const data = {
      ...form,
      // 将图片数组转换为 JSON 字符串
      images: form.images.length > 0 ? JSON.stringify(form.images) : '[]'
    }
    
    console.log('[提交] 最终数据:', data)
    
    if (currentEditId.value) {
      await updateWorkIssue(currentEditId.value, data)
      ElMessage.success('更新成功')
    } else {
      await createWorkIssue(data)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    fetchData()
  } catch (error) {
    if (error !== false) {
      ElMessage.error(currentEditId.value ? '更新失败' : '创建失败')
    }
  } finally {
    submitting.value = false
  }
}

// 重置表单
const resetForm = () => {
  Object.assign(form, {
    title: '',
    description: '',
    status: 'open',
    severity: 'medium',
    solution: '',
    tags: '',
    images: [],
    content: ''
  })
  imageFileList.value = []
  formRef.value?.clearValidate()
}

// 对话框关闭
const handleDialogClose = () => {
  resetForm()
}

onMounted(() => {
  fetchData()
})

defineExpose({
  refresh: fetchData
})
</script>

<style scoped>
.issue-list {
  width: 100%;
}

/* 操作栏 */
.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  gap: 24px;
  flex-wrap: wrap;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 16px;
  flex: 1;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: linear-gradient(135deg, #fff5f5 0%, #ffffff 100%);
  border-radius: 12px;
  border: 1px solid #ffe0e0;
  transition: all 0.3s ease;
  cursor: pointer;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(239, 68, 68, 0.15);
}

.stat-icon {
  font-size: 32px;
  padding: 8px;
  border-radius: 10px;
}

.stat-total .stat-icon {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
}

.stat-open .stat-icon {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  color: white;
}

.stat-progress .stat-icon {
  background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
  color: white;
}

.stat-resolved .stat-icon {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #2d3748;
  line-height: 1;
}

.stat-label {
  font-size: 12px;
  color: #718096;
  margin-top: 4px;
}

.action-buttons {
  display: flex;
  gap: 12px;
  align-items: center;
}

.status-filter {
  width: 160px;
}

.date-filter {
  width: 280px;
}

.add-btn {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  border: none;
  color: white;
  font-weight: 600;
}

.add-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(245, 158, 11, 0.3);
}

/* 表格样式 */
.issue-table {
  width: 100%;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.issue-table :deep(.el-table__row) {
  transition: all 0.2s ease;
}

.issue-table :deep(.el-table__row:hover) {
  background-color: #fff5f5;
}

/* 左右分栏容器 */
.issue-cell-wrapper {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 8px 0;
}

/* 左侧主内容区 */
.issue-main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* 标题和描述容器 */
.issue-cell {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.issue-title {
  font-size: 14px;
  font-weight: 600;
  color: #2d3748;
  line-height: 1.4;
}

.issue-desc {
  font-size: 13px;
  color: #718096;
  margin-bottom: 6px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.issue-solution {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 10px;
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  border-radius: 6px;
  margin: 6px 0;
  font-size: 12px;
  color: #065f46;
}

.solution-icon {
  color: #059669;
}

.issue-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 6px;
}

.issue-tag {
  background: #fff7ed;
  border: 1px solid #fed7aa;
  color: #c2410c;
  font-size: 12px;
}

.status-tag {
  font-weight: 600;
  padding: 6px 12px;
}

.severity-tag {
  font-weight: 600;
  padding: 6px 12px;
}

.clickable {
  cursor: pointer;
  user-select: none;
  transition: all 0.3s ease;
}

.clickable:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.text-muted {
  color: #a0aec0;
}

.resolved-time {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  color: #10b981;
  font-weight: 600;
}

/* 右侧图片区（覆盖标题+描述高度） */
.issue-images-area {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
  align-self: flex-start;
}

.issue-image {
  width: 55px;
  height: 55px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  border: 2px solid #e2e8f0;
  flex-shrink: 0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.issue-image:hover {
  transform: scale(1.05);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.image-error {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: #f5f5f5;
  color: #999;
  font-size: 20px;
}

/* 分页 */
.pagination {
  display: flex;
  justify-content: center;
  padding: 24px 0;
  background: white;
  border-radius: 0 0 8px 8px;
}

/* 对话框 */
.issue-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  padding: 20px;
}

.issue-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
  font-size: 18px;
}

.issue-dialog :deep(.el-dialog__headerbtn .el-dialog__close) {
  color: white;
  font-size: 20px;
}

.issue-dialog :deep(.el-dialog__headerbtn:hover .el-dialog__close) {
  color: white;
}

/* 图片上传区域 */
.image-upload-area {
  width: 100%;
}

.upload-tip {
  margin-top: 8px;
  font-size: 12px;
  color: #909399;
}

.content-tip {
  margin-top: 4px;
  font-size: 12px;
  color: #909399;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .action-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .stats-cards {
    grid-template-columns: repeat(2, 1fr);
  }

  .action-buttons {
    justify-content: space-between;
  }

  .status-filter {
    flex: 1;
  }
}
</style>
