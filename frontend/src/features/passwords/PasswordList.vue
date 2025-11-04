<template>
  <div class="password-list-page">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon" :size="40" color="#409eff">
              <Lock />
            </el-icon>
            <div class="stat-text">
              <div class="stat-value">{{ stats.total || 0 }}</div>
              <div class="stat-label">总账户数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon" :size="40" color="#f56c6c">
              <Star />
            </el-icon>
            <div class="stat-text">
              <div class="stat-value">{{ stats.fav_count || 0 }}</div>
              <div class="stat-label">收藏账户</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon" :size="40" color="#67c23a">
              <FolderOpened />
            </el-icon>
            <div class="stat-text">
              <div class="stat-value">{{ Object.keys(stats.category_counts || {}).length }}</div>
              <div class="stat-label">分类数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <el-icon class="stat-icon" :size="40" color="#e6a23c">
              <Setting />
            </el-icon>
            <div class="stat-text">
              <div class="stat-value">AES-256</div>
              <div class="stat-label">加密算法</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 操作栏 -->
    <el-card class="operation-card" shadow="never">
      <el-row :gutter="12" align="middle">
        <el-col :span="8">
          <el-input
            v-model="keyword"
            placeholder="搜索标题、URL、用户名或备注..."
            :prefix-icon="Search"
            clearable
            @clear="fetchData"
            @keyup.enter="fetchData"
          />
        </el-col>
        <el-col :span="5">
          <el-select
            v-model="selectedCategory"
            placeholder="选择分类"
            clearable
            @change="fetchData"
            style="width: 100%"
          >
            <el-option
              v-for="cat in categories"
              :key="cat"
              :label="cat"
              :value="cat"
            />
          </el-select>
        </el-col>
        <el-col :span="3">
          <el-button type="primary" :icon="Search" @click="fetchData">
            搜索
          </el-button>
        </el-col>
        <el-col :span="8" style="text-align: right;">
          <el-button type="primary" :icon="Plus" @click="showAddDialog = true">
            添加账户
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 数据表格 -->
    <el-card class="table-card" shadow="never">
      <el-table
        :data="tableData"
        style="width: 100%"
        v-loading="loading"
        empty-text="暂无数据"
      >
        <el-table-column label="序号" width="80" type="index" :index="indexMethod" />
        
        <el-table-column label="收藏" width="80" align="center">
          <template #default="{ row }">
            <el-icon 
              :size="20" 
              :color="row.is_fav ? '#f56c6c' : '#dcdfe6'"
              style="cursor: pointer"
              @click="toggleFavorite(row)"
            >
              <StarFilled v-if="row.is_fav" />
              <Star v-else />
            </el-icon>
          </template>
        </el-table-column>

        <el-table-column prop="title" label="标题" width="200">
          <template #default="{ row }">
            <div class="title-cell">
              <el-icon :size="16" color="#409eff">
                <Link />
              </el-icon>
              <span>{{ row.title }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="url" label="网站URL" min-width="250">
          <template #default="{ row }">
            <a 
              v-if="row.url" 
              :href="row.url" 
              target="_blank" 
              class="url-link"
            >
              {{ row.url }}
              <el-icon :size="12"><TopRight /></el-icon>
            </a>
            <span v-else style="color: #909399">-</span>
          </template>
        </el-table-column>

        <el-table-column prop="username" label="用户名" width="220">
          <template #default="{ row }">
            <div v-if="row.username" class="username-cell">
              <span class="username-text">{{ row.username }}</span>
              <el-icon 
                :size="14" 
                class="copy-icon"
                @click="copyToClipboard(row.username, '用户名')"
              >
                <CopyDocument />
              </el-icon>
            </div>
            <span v-else style="color: #909399">-</span>
          </template>
        </el-table-column>

        <el-table-column prop="category" label="分类" width="140">
          <template #default="{ row }">
            <el-tag v-if="row.category" size="small" type="info">
              {{ row.category }}
            </el-tag>
            <span v-else style="color: #909399">-</span>
          </template>
        </el-table-column>

        <el-table-column label="密码" width="150" align="center">
          <template #default="{ row }">
            <el-button
              size="small"
              type="primary"
              :icon="View"
              @click="viewPassword(row)"
              :loading="row.passwordLoading"
            >
              查看密码
            </el-button>
          </template>
        </el-table-column>

        <el-table-column prop="updated_at" label="更新时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.updated_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button
              size="small"
              type="primary"
              :icon="Edit"
              @click="handleEdit(row)"
              link
            >
              编辑
            </el-button>
            <el-button
              size="small"
              type="danger"
              :icon="Delete"
              @click="handleDelete(row)"
              link
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="showAddDialog"
      :title="editingId ? '编辑账户' : '添加账户'"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="标题" prop="title">
          <el-input
            v-model="form.title"
            placeholder="例如：Google账户、GitHub等"
            :prefix-icon="Document"
          />
        </el-form-item>

        <el-form-item label="网站URL" prop="url">
          <el-input
            v-model="form.url"
            placeholder="https://example.com"
            :prefix-icon="Link"
          />
        </el-form-item>

        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="form.username"
            placeholder="用户名或邮箱"
            :prefix-icon="User"
          />
        </el-form-item>

        <el-form-item label="密码" prop="password">
          <el-input
            v-model="form.password"
            :type="showPassword ? 'text' : 'password'"
            placeholder="请输入密码"
            :prefix-icon="Lock"
          >
            <template #suffix>
              <el-icon
                style="cursor: pointer"
                @click="showPassword = !showPassword"
              >
                <View v-if="showPassword" />
                <Hide v-else />
              </el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item label="分类" prop="category">
          <el-select
            v-model="form.category"
            placeholder="选择或输入分类"
            allow-create
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="cat in categories"
              :key="cat"
              :label="cat"
              :value="cat"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="备注" prop="notes">
          <el-input
            v-model="form.notes"
            type="textarea"
            :rows="3"
            placeholder="备注信息"
          />
        </el-form-item>

        <el-form-item label="收藏" prop="is_fav">
          <el-switch v-model="form.is_fav" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">
          保存
        </el-button>
      </template>
    </el-dialog>

    <!-- 查看密码对话框 -->
    <el-dialog
      v-model="showPasswordDialog"
      title="查看密码"
      width="500px"
    >
      <div class="password-view">
        <el-alert
          title="安全提示"
          type="warning"
          :closable="false"
          show-icon
          style="margin-bottom: 20px"
        >
          请确保周围环境安全，不要在公共场合查看密码
        </el-alert>

        <div class="password-info">
          <div class="info-item">
            <span class="label">标题：</span>
            <span class="value">{{ currentPassword.title }}</span>
          </div>
          <div class="info-item">
            <span class="label">用户名：</span>
            <span class="value">{{ currentPassword.username || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="label">密码：</span>
            <div class="password-display">
              <el-input
                v-model="currentPassword.password"
                readonly
                :type="showCurrentPassword ? 'text' : 'password'"
              >
                <template #suffix>
                  <el-icon
                    style="cursor: pointer; margin-right: 5px"
                    @click="showCurrentPassword = !showCurrentPassword"
                  >
                    <View v-if="showCurrentPassword" />
                    <Hide v-else />
                  </el-icon>
                  <el-icon
                    style="cursor: pointer"
                    @click="copyToClipboard(currentPassword.password, '密码')"
                  >
                    <CopyDocument />
                  </el-icon>
                </template>
              </el-input>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <el-button @click="showPasswordDialog = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus, Search, Edit, Delete, View, Hide, Lock, Star, StarFilled,
  Link, User, Document, InfoFilled, TopRight, CopyDocument, Setting,
  FolderOpened
} from '@element-plus/icons-vue'
import api from './api'

// 数据状态
const loading = ref(false)
const saving = ref(false)
const tableData = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')
const selectedCategory = ref('')
const categories = ref([])
const stats = ref({})

// 表单状态
const showAddDialog = ref(false)
const showPassword = ref(false)
const editingId = ref(null)
const formRef = ref(null)
const form = reactive({
  title: '',
  url: '',
  username: '',
  password: '',
  category: '',
  notes: '',
  is_fav: false
})

// 密码查看状态
const showPasswordDialog = ref(false)
const showCurrentPassword = ref(false)
const currentPassword = ref({})

// 表单验证规则
const formRules = {
  title: [
    { required: true, message: '请输入标题', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    
    if (keyword.value) {
      params.keyword = keyword.value
    }
    
    if (selectedCategory.value) {
      params.category = selectedCategory.value
    }

    const response = await api.getList(params)
    tableData.value = response.data.items || []
    total.value = response.data.total || 0
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '获取列表失败')
  } finally {
    loading.value = false
  }
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    const response = await api.getCategories()
    categories.value = response.data.categories || []
  } catch (error) {
    console.error('获取分类失败:', error)
  }
}

// 获取统计信息
const fetchStats = async () => {
  try {
    const response = await api.getStats()
    stats.value = response.data || {}
  } catch (error) {
    console.error('获取统计失败:', error)
  }
}

// 查看密码
const viewPassword = async (row) => {
  row.passwordLoading = true
  try {
    const response = await api.getPassword(row.id)
    currentPassword.value = {
      title: row.title,
      username: row.username,
      password: response.data.password
    }
    showPasswordDialog.value = true
    showCurrentPassword.value = false
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '获取密码失败')
  } finally {
    row.passwordLoading = false
  }
}

// 切换收藏状态
const toggleFavorite = async (row) => {
  try {
    await api.update(row.id, { is_fav: !row.is_fav })
    row.is_fav = !row.is_fav
    ElMessage.success(row.is_fav ? '已添加到收藏' : '已取消收藏')
    fetchStats()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

// 编辑
const handleEdit = async (row) => {
  editingId.value = row.id
  Object.assign(form, {
    title: row.title,
    url: row.url,
    username: row.username,
    password: '', // 不显示原密码
    category: row.category,
    notes: row.notes,
    is_fav: row.is_fav
  })
  showAddDialog.value = true
}

// 保存
const handleSave = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    saving.value = true
    try {
      if (editingId.value) {
        // 编辑时，如果密码为空则不更新密码
        const updateData = { ...form }
        if (!updateData.password) {
          delete updateData.password
        }
        await api.update(editingId.value, updateData)
        ElMessage.success('更新成功')
      } else {
        await api.create(form)
        ElMessage.success('添加成功')
      }
      
      showAddDialog.value = false
      resetForm()
      fetchData()
      fetchCategories()
      fetchStats()
    } catch (error) {
      ElMessage.error(error.response?.data?.message || '保存失败')
    } finally {
      saving.value = false
    }
  })
}

// 删除
const handleDelete = async (row) => {
  ElMessageBox.confirm(
    `确定要删除账户 "${row.title}" 吗？此操作不可恢复！`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await api.delete(row.id)
      ElMessage.success('删除成功')
      fetchData()
      fetchStats()
    } catch (error) {
      ElMessage.error(error.response?.data?.message || '删除失败')
    }
  }).catch(() => {})
}

// 重置表单
const resetForm = () => {
  editingId.value = null
  Object.assign(form, {
    title: '',
    url: '',
    username: '',
    password: '',
    category: '',
    notes: '',
    is_fav: false
  })
  formRef.value?.resetFields()
}

// 复制到剪贴板
const copyToClipboard = async (text, label) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success(`${label}已复制到剪贴板`)
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

// 序号方法（根据分页计算）
const indexMethod = (index) => {
  return (page.value - 1) * pageSize.value + index + 1
}

// 格式化日期
const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 初始化
onMounted(() => {
  fetchData()
  fetchCategories()
  fetchStats()
})
</script>

<style scoped>
.password-list-page {
  padding: 0;
}

/* ========== 统计卡片样式 ========== */
.stats-row {
  margin-bottom: 24px;
}

.stat-card {
  cursor: pointer;
  border: none;
  border-radius: 12px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  overflow: hidden;
  position: relative;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  opacity: 0;
  transition: opacity 0.3s;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.12);
}

.stat-card:hover::before {
  opacity: 1;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 4px;
}

.stat-icon {
  flex-shrink: 0;
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.05) 100%);
}

.stat-text {
  flex: 1;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: #303133;
  line-height: 1.2;
  letter-spacing: -0.5px;
}

.stat-label {
  font-size: 13px;
  color: #909399;
  margin-top: 4px;
  font-weight: 500;
}

/* ========== 操作卡片样式 ========== */
.operation-card {
  margin-bottom: 24px;
  border-radius: 8px;
  border: none;
}

.operation-card :deep(.el-card__body) {
  padding: 20px 24px;
}

/* ========== 表格卡片样式 ========== */
.table-card {
  margin-bottom: 24px;
  border-radius: 8px;
  border: none;
  overflow: hidden;
}

.table-card :deep(.el-card__body) {
  padding: 0;
}

.table-card :deep(.el-table) {
  font-size: 14px;
}

.table-card :deep(.el-table th) {
  background-color: #fafafa !important;
  color: #303133;
  font-weight: 600;
  font-size: 13px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  padding: 16px 0;
}

.table-card :deep(.el-table td) {
  padding: 16px 0;
  border-bottom: 1px solid #f0f0f0;
}

.table-card :deep(.el-table__row:hover) {
  background-color: #f5f7fa !important;
}

/* ========== 表格单元格样式 ========== */
.title-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  color: #303133;
}

.url-link {
  color: #409eff;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  transition: all 0.3s;
  font-weight: 500;
}

.url-link:hover {
  color: #66b1ff;
  text-decoration: underline;
}

.username-cell {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
}

.username-text {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: 500;
  color: #606266;
}

.copy-icon {
  flex-shrink: 0;
  cursor: pointer;
  color: #409eff;
  transition: all 0.3s;
  padding: 4px;
  border-radius: 4px;
}

.copy-icon:hover {
  color: #66b1ff;
  background-color: #ecf5ff;
}

/* ========== 分页容器 ========== */
.pagination-container {
  margin-top: 20px;
  padding: 16px 24px;
  display: flex;
  justify-content: flex-end;
  background: white;
  border-top: 1px solid #f0f0f0;
}

.password-view {
  padding: 10px 0;
}

.password-info {
  margin-top: 20px;
}

.info-item {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.info-item .label {
  font-weight: 600;
  color: #606266;
  min-width: 80px;
}

.info-item .value {
  color: #303133;
}

.password-display {
  flex: 1;
}
</style>

