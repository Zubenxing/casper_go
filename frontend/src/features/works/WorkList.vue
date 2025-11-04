<template>
  <div class="work-list">
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
          <el-option label="待办" value="pending" />
          <el-option label="进行中" value="in_progress" />
          <el-option label="已完成" value="completed" />
          <el-option label="已取消" value="cancelled" />
        </el-select>

        <el-button type="primary" size="large" @click="handleAdd" class="add-btn">
          <el-icon><Plus /></el-icon>
          <span>新建任务</span>
        </el-button>
      </div>
    </div>

    <!-- 表格 -->
    <el-table
      v-loading="loading"
      :data="workList"
      stripe
      class="work-table"
      :header-cell-style="{ background: '#f5f7fa', color: '#606266', fontWeight: '600' }"
    >
      <el-table-column type="index" label="#" width="60" align="center" />

      <el-table-column label="任务" min-width="300">
        <template #default="{ row }">
          <div class="task-cell">
            <div class="task-title">{{ row.title }}</div>
            <div v-if="row.description" class="task-desc">{{ row.description }}</div>
            <div v-if="row.tags" class="task-tags">
              <el-tag
                v-for="(tag, index) in row.tags.split(',')"
                :key="index"
                size="small"
                class="task-tag"
              >
                {{ tag }}
              </el-tag>
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
                <el-dropdown-item command="pending">待办</el-dropdown-item>
                <el-dropdown-item command="in_progress">进行中</el-dropdown-item>
                <el-dropdown-item command="completed">已完成</el-dropdown-item>
                <el-dropdown-item command="cancelled">已取消</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>

      <el-table-column label="优先级" width="120" align="center" sortable prop="priority">
        <template #default="{ row }">
          <el-dropdown @command="(cmd) => handlePriorityChange(row, cmd)" trigger="click">
            <el-tag
              :type="getPriorityType(row.priority)"
              effect="plain"
              size="large"
              class="priority-tag clickable"
            >
              {{ getPriorityText(row.priority) }}
              <el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-tag>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="low">低</el-dropdown-item>
                <el-dropdown-item command="medium">中</el-dropdown-item>
                <el-dropdown-item command="high">高</el-dropdown-item>
                <el-dropdown-item command="urgent">紧急</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>

      <el-table-column label="截止日期" width="140" align="center" sortable prop="due_date">
        <template #default="{ row }">
          <span v-if="row.due_date" :class="{ 'overdue': isOverdue(row.due_date) }">
            {{ formatDate(row.due_date) }}
          </span>
          <span v-else class="text-muted">-</span>
        </template>
      </el-table-column>

      <el-table-column label="开始日期" width="140" align="center" sortable prop="start_date">
        <template #default="{ row }">
          <span v-if="row.start_date">{{ formatDate(row.start_date) }}</span>
          <span v-else class="text-muted">-</span>
        </template>
      </el-table-column>

      <el-table-column label="创建时间" width="140" align="center">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>

      <el-table-column label="操作" width="120" align="center" fixed="right">
        <template #default="{ row }">
          <el-button
            type="primary"
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
      class="work-dialog"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="90px"
      >
        <el-form-item label="任务标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入任务标题" />
        </el-form-item>

        <el-form-item label="任务描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="4"
            placeholder="请输入任务描述"
          />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select v-model="form.status" placeholder="请选择状态" style="width: 100%">
                <el-option label="待办" value="pending" />
                <el-option label="进行中" value="in_progress" />
                <el-option label="已完成" value="completed" />
                <el-option label="已取消" value="cancelled" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="优先级" prop="priority">
              <el-select v-model="form.priority" placeholder="请选择优先级" style="width: 100%">
                <el-option label="低" value="low" />
                <el-option label="中" value="medium" />
                <el-option label="高" value="high" />
                <el-option label="紧急" value="urgent" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="开始日期">
              <el-date-picker
                v-model="form.start_date"
                type="datetime"
                placeholder="选择开始日期"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="截止日期">
              <el-date-picker
                v-model="form.due_date"
                type="datetime"
                placeholder="选择截止日期"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="标签">
          <el-input v-model="form.tags" placeholder="多个标签用逗号分隔" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false" size="large">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting" size="large">
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
  Plus, Edit, Delete, Clock, Loading, CircleCheck, Tickets, ArrowDown
} from '@element-plus/icons-vue'
import {
  getWorkList, createWork, updateWork, deleteWork, getWorkStats
} from './api'

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('新建任务')
const filterStatus = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const workList = ref([])
const formRef = ref(null)
const currentEditId = ref(null)

const stats = ref({
  total: 0,
  pending: 0,
  in_progress: 0,
  completed: 0
})

const form = reactive({
  title: '',
  description: '',
  status: 'pending',
  priority: 'medium',
  start_date: null,
  due_date: null,
  tags: ''
})

const rules = {
  title: [{ required: true, message: '请输入任务标题', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }],
  priority: [{ required: true, message: '请选择优先级', trigger: 'change' }]
}

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const [listRes, statsRes] = await Promise.all([
      getWorkList({
        page: page.value,
        page_size: pageSize.value,
        status: filterStatus.value
      }),
      getWorkStats()
    ])
    workList.value = listRes.data.list || []
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
    pending: 'info',
    in_progress: 'warning',
    completed: 'success',
    cancelled: 'danger'
  }
  return map[status] || 'info'
}

const getStatusText = (status) => {
  const map = {
    pending: '待办',
    in_progress: '进行中',
    completed: '已完成',
    cancelled: '已取消'
  }
  return map[status] || status
}

const getPriorityType = (priority) => {
  const map = {
    low: 'info',
    medium: 'warning',
    high: 'warning',
    urgent: 'danger'
  }
  return map[priority] || 'info'
}

const getPriorityText = (priority) => {
  const map = {
    low: '低',
    medium: '中',
    high: '高',
    urgent: '紧急'
  }
  return map[priority] || priority
}

// 快速更改状态
const handleStatusChange = async (row, newStatus) => {
  if (row.status === newStatus) return
  
  try {
    await updateWork(row.id, { status: newStatus })
    ElMessage.success('状态已更新')
    row.status = newStatus
    // 刷新统计数据
    await fetchStats()
  } catch (error) {
    ElMessage.error(error.message || '状态更新失败')
  }
}

// 快速更改优先级
const handlePriorityChange = async (row, newPriority) => {
  if (row.priority === newPriority) return
  
  try {
    await updateWork(row.id, { priority: newPriority })
    ElMessage.success('优先级已更新')
    row.priority = newPriority
  } catch (error) {
    ElMessage.error(error.message || '优先级更新失败')
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

const isOverdue = (dueDate) => {
  if (!dueDate) return false
  return new Date(dueDate) < new Date()
}

// 筛选变更
const handleFilterChange = () => {
  page.value = 1
  fetchData()
}

// 新建任务
const handleAdd = () => {
  dialogTitle.value = '新建任务'
  currentEditId.value = null
  resetForm()
  dialogVisible.value = true
}

// 编辑任务
const handleEdit = (work) => {
  dialogTitle.value = '编辑任务'
  currentEditId.value = work.id
  Object.assign(form, {
    title: work.title,
    description: work.description,
    status: work.status,
    priority: work.priority,
    start_date: work.start_date ? new Date(work.start_date) : null,
    due_date: work.due_date ? new Date(work.due_date) : null,
    tags: work.tags || ''
  })
  dialogVisible.value = true
}

// 删除任务
const handleDelete = async (work) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除任务"${work.title}"吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await deleteWork(work.id)
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
      await updateWork(currentEditId.value, data)
      ElMessage.success('更新成功')
    } else {
      await createWork(data)
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
    status: 'pending',
    priority: 'medium',
    start_date: null,
    due_date: null,
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
.work-list {
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
  background: linear-gradient(135deg, #f6f8fb 0%, #ffffff 100%);
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  transition: all 0.3s ease;
  cursor: pointer;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  font-size: 32px;
  padding: 8px;
  border-radius: 10px;
}

.stat-total .stat-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.stat-pending .stat-icon {
  background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%);
  color: white;
}

.stat-progress .stat-icon {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
}

.stat-completed .stat-icon {
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  font-weight: 600;
}

.add-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(102, 126, 234, 0.3);
}

/* 表格样式 */
.work-table {
  width: 100%;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.work-table :deep(.el-table__row) {
  transition: all 0.2s ease;
}

.work-table :deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}

.task-cell {
  padding: 4px 0;
}

.task-title {
  font-size: 14px;
  font-weight: 600;
  color: #2d3748;
  margin-bottom: 4px;
}

.task-desc {
  font-size: 13px;
  color: #718096;
  margin-bottom: 6px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.task-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 6px;
}

.task-tag {
  background: #f7fafc;
  border: 1px solid #e2e8f0;
  color: #4a5568;
  font-size: 12px;
}

.status-tag {
  font-weight: 600;
  padding: 6px 12px;
}

.priority-tag {
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

.overdue {
  color: #ef4444;
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
.work-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.work-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
  font-size: 18px;
}

.work-dialog :deep(.el-dialog__headerbtn .el-dialog__close) {
  color: white;
  font-size: 20px;
}

.work-dialog :deep(.el-dialog__headerbtn:hover .el-dialog__close) {
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
