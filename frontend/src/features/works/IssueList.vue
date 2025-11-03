<template>
  <div class="issue-list">
    <!-- 顶部操作栏 -->
    <div class="action-bar">
      <div class="stats-cards">
        <div class="stat-card stat-total">
          <el-icon class="stat-icon"><DocumentChecked /></el-icon>
          <div class="stat-content">
            <div class="stat-value">{{ stats.total || 0 }}</div>
            <div class="stat-label">全部问题</div>
          </div>
        </div>
        <div class="stat-card stat-open">
          <el-icon class="stat-icon"><Warning /></el-icon>
          <div class="stat-content">
            <div class="stat-value">{{ stats.open || 0 }}</div>
            <div class="stat-label">待处理</div>
          </div>
        </div>
        <div class="stat-card stat-progress">
          <el-icon class="stat-icon"><Loading /></el-icon>
          <div class="stat-content">
            <div class="stat-value">{{ stats.in_progress || 0 }}</div>
            <div class="stat-label">处理中</div>
          </div>
        </div>
        <div class="stat-card stat-resolved">
          <el-icon class="stat-icon"><CircleCheck /></el-icon>
          <div class="stat-content">
            <div class="stat-value">{{ stats.resolved || 0 }}</div>
            <div class="stat-label">已解决</div>
          </div>
        </div>
      </div>

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

      <el-table-column label="问题" min-width="300">
        <template #default="{ row }">
          <div class="issue-cell">
            <div class="issue-title">{{ row.title }}</div>
            <div v-if="row.description" class="issue-desc">{{ row.description }}</div>
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
        </template>
      </el-table-column>

      <el-table-column label="状态" width="140" align="center">
        <template #default="{ row }">
          <el-tag
            :type="getStatusType(row.status)"
            effect="dark"
            size="large"
            class="status-tag"
          >
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="严重程度" width="120" align="center">
        <template #default="{ row }">
          <el-tag
            :type="getSeverityType(row.severity)"
            effect="plain"
            size="large"
          >
            {{ getSeverityText(row.severity) }}
          </el-tag>
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
          <el-button
            type="warning"
            link
            size="small"
            @click="handleEdit(row)"
          >
            <el-icon><Edit /></el-icon>
          </el-button>
          <el-button
            type="danger"
            link
            size="small"
            @click="handleDelete(row)"
          >
            <el-icon><Delete /></el-icon>
          </el-button>
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
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false" size="large">取消</el-button>
        <el-button type="warning" @click="handleSubmit" :loading="submitting" size="large">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, Edit, Delete, Warning, Loading, CircleCheck, 
  DocumentChecked, Checked
} from '@element-plus/icons-vue'
import {
  getWorkIssueList, createWorkIssue, updateWorkIssue, 
  deleteWorkIssue, getWorkIssueStats
} from './api'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('新建问题')
const filterStatus = ref('')
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
  tags: ''
})

const rules = {
  title: [{ required: true, message: '请输入问题标题', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }],
  severity: [{ required: true, message: '请选择严重程度', trigger: 'change' }]
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
    issueList.value = listRes.data.list || []
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
    medium: '',
    high: 'warning',
    critical: 'danger'
  }
  return map[severity] || ''
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

// 筛选变更
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
  Object.assign(form, {
    title: issue.title,
    description: issue.description,
    status: issue.status,
    severity: issue.severity,
    solution: issue.solution || '',
    tags: issue.tags || ''
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

    const data = { ...form }
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
    tags: ''
  })
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

.issue-cell {
  padding: 4px 0;
}

.issue-title {
  font-size: 14px;
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 4px;
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
