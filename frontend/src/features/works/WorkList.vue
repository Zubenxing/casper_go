<template>
  <div class="work-list">
    <!-- 顶部操作栏 -->
    <div class="action-bar">
      <div class="stats-cards">
        <div class="stat-card stat-total">
          <el-icon class="stat-icon"><Tickets /></el-icon>
          <div class="stat-content">
            <div class="stat-value">{{ stats.total || 0 }}</div>
            <div class="stat-label">全部任务</div>
          </div>
        </div>
        <div class="stat-card stat-pending">
          <el-icon class="stat-icon"><Clock /></el-icon>
          <div class="stat-content">
            <div class="stat-value">{{ stats.pending || 0 }}</div>
            <div class="stat-label">待办</div>
          </div>
        </div>
        <div class="stat-card stat-progress">
          <el-icon class="stat-icon"><Loading /></el-icon>
          <div class="stat-content">
            <div class="stat-value">{{ stats.in_progress || 0 }}</div>
            <div class="stat-label">进行中</div>
          </div>
        </div>
        <div class="stat-card stat-completed">
          <el-icon class="stat-icon"><CircleCheck /></el-icon>
          <div class="stat-content">
            <div class="stat-value">{{ stats.completed || 0 }}</div>
            <div class="stat-label">已完成</div>
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
        >
          <el-option label="全部状态" value="" />
          <el-option label="待办" value="pending" />
          <el-option label="进行中" value="in_progress" />
          <el-option label="已完成" value="completed" />
          <el-option label="已取消" value="cancelled" />
        </el-select>

        <el-button type="primary" @click="handleAdd" class="add-btn">
          <el-icon><Plus /></el-icon>
          <span>新建任务</span>
        </el-button>
      </div>
    </div>

    <!-- 工作列表 -->
    <div v-loading="loading" class="work-cards">
      <el-empty v-if="!loading && workList.length === 0" description="暂无工作记录" />

      <transition-group name="list">
        <div
          v-for="work in workList"
          :key="work.id"
          class="work-card"
          :class="`status-${work.status}`"
        >
          <!-- 卡片头部 -->
          <div class="card-header">
            <div class="header-left">
              <el-tag
                :type="getStatusType(work.status)"
                size="small"
                effect="dark"
                class="status-tag"
              >
                {{ getStatusText(work.status) }}
              </el-tag>
              <el-tag
                :type="getPriorityType(work.priority)"
                size="small"
                effect="plain"
                class="priority-tag"
              >
                {{ getPriorityText(work.priority) }}
              </el-tag>
            </div>
            <div class="header-right">
              <el-button
                type="primary"
                size="small"
                link
                @click="handleEdit(work)"
              >
                <el-icon><Edit /></el-icon>
              </el-button>
              <el-button
                type="danger"
                size="small"
                link
                @click="handleDelete(work)"
              >
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
          </div>

          <!-- 卡片内容 -->
          <div class="card-body">
            <h3 class="work-title">{{ work.title }}</h3>
            <p v-if="work.description" class="work-desc">{{ work.description }}</p>

            <div v-if="work.tags" class="work-tags">
              <el-tag
                v-for="(tag, index) in work.tags.split(',')"
                :key="index"
                size="small"
                class="work-tag"
              >
                {{ tag }}
              </el-tag>
            </div>
          </div>

          <!-- 卡片底部 -->
          <div class="card-footer">
            <div class="time-info">
              <el-icon><Calendar /></el-icon>
              <span>创建于 {{ formatDate(work.created_at) }}</span>
            </div>
            <div v-if="work.due_date" class="due-date" :class="{ overdue: isOverdue(work.due_date) }">
              <el-icon><Timer /></el-icon>
              <span>{{ formatDate(work.due_date) }}</span>
            </div>
          </div>
        </div>
      </transition-group>
    </div>

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
        label-width="80px"
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

        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status" placeholder="请选择状态">
            <el-option label="待办" value="pending" />
            <el-option label="进行中" value="in_progress" />
            <el-option label="已完成" value="completed" />
            <el-option label="已取消" value="cancelled" />
          </el-select>
        </el-form-item>

        <el-form-item label="优先级" prop="priority">
          <el-select v-model="form.priority" placeholder="请选择优先级">
            <el-option label="低" value="low" />
            <el-option label="中" value="medium" />
            <el-option label="高" value="high" />
            <el-option label="紧急" value="urgent" />
          </el-select>
        </el-form-item>

        <el-form-item label="开始日期">
          <el-date-picker
            v-model="form.start_date"
            type="datetime"
            placeholder="选择开始日期"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="截止日期">
          <el-date-picker
            v-model="form.due_date"
            type="datetime"
            placeholder="选择截止日期"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="标签">
          <el-input v-model="form.tags" placeholder="多个标签用逗号分隔" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
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
  Plus, Edit, Delete, Calendar, Timer, Clock, Loading, CircleCheck, Tickets
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
    medium: '',
    high: 'warning',
    urgent: 'danger'
  }
  return map[priority] || ''
}

const getPriorityText = (priority) => {
  const map = {
    low: '低优先级',
    medium: '中优先级',
    high: '高优先级',
    urgent: '紧急'
  }
  return map[priority] || priority
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
  width: 140px;
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

/* 工作卡片 */
.work-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.work-card {
  background: white;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  overflow: hidden;
  transition: all 0.3s ease;
  cursor: pointer;
}

.work-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
  border-color: #667eea;
}

.work-card.status-pending {
  border-left: 4px solid #fbbf24;
}

.work-card.status-in_progress {
  border-left: 4px solid #3b82f6;
}

.work-card.status-completed {
  border-left: 4px solid #10b981;
}

.work-card.status-cancelled {
  border-left: 4px solid #ef4444;
  opacity: 0.7;
}

/* 卡片头部 */
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: linear-gradient(135deg, #f6f8fb 0%, #ffffff 100%);
  border-bottom: 1px solid #e2e8f0;
}

.header-left {
  display: flex;
  gap: 8px;
}

.status-tag {
  font-weight: 600;
  border-radius: 6px;
}

.priority-tag {
  border-radius: 6px;
}

.header-right {
  display: flex;
  gap: 4px;
}

/* 卡片内容 */
.card-body {
  padding: 16px;
}

.work-title {
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
  margin: 0 0 8px 0;
  line-height: 1.4;
}

.work-desc {
  font-size: 14px;
  color: #718096;
  margin: 0 0 12px 0;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.work-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.work-tag {
  background: #f7fafc;
  border: 1px solid #e2e8f0;
  color: #4a5568;
  border-radius: 4px;
}

/* 卡片底部 */
.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f7fafc;
  border-top: 1px solid #e2e8f0;
  font-size: 12px;
  color: #718096;
}

.time-info, .due-date {
  display: flex;
  align-items: center;
  gap: 4px;
}

.due-date.overdue {
  color: #ef4444;
  font-weight: 600;
}

/* 分页 */
.pagination {
  display: flex;
  justify-content: center;
  padding: 24px 0;
}

/* 列表动画 */
.list-enter-active, .list-leave-active {
  transition: all 0.3s ease;
}

.list-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.list-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

/* 对话框 */
.work-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.work-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
}

.work-dialog :deep(.el-dialog__headerbtn .el-dialog__close) {
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

  .work-cards {
    grid-template-columns: 1fr;
  }
}
</style>

