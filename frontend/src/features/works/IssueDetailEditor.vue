<template>
  <el-dialog
    v-model="visible"
    :title="issue.title || '问题详情'"
    width="80%"
    top="5vh"
    align-center
    draggable
    class="issue-detail-dialog"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div class="editor-container">
      <!-- 头部信息栏 -->
      <div class="info-bar">
        <div class="info-item">
          <span class="label">状态：</span>
          <el-tag :type="getStatusType(issue.status)" effect="dark" size="small">
            {{ getStatusText(issue.status) }}
          </el-tag>
        </div>
        <div class="info-item">
          <span class="label">严重程度：</span>
          <el-tag :type="getSeverityType(issue.severity)" effect="plain" size="small">
            {{ getSeverityText(issue.severity) }}
          </el-tag>
        </div>
        <div class="info-item">
          <span class="label">创建时间：</span>
          <span class="value">{{ formatDate(issue.created_at) }}</span>
        </div>
        <div v-if="issue.resolved_at" class="info-item">
          <span class="label">解决时间：</span>
          <span class="value">{{ formatDate(issue.resolved_at) }}</span>
        </div>
      </div>

      <!-- 问题描述 -->
      <div class="description-section">
        <div class="section-title">问题描述</div>
        <div class="description-text">{{ issue.description || '无' }}</div>
      </div>

      <!-- 解决方案 -->
      <div v-if="issue.solution" class="solution-section">
        <div class="section-title">解决方案</div>
        <div class="solution-text">{{ issue.solution }}</div>
      </div>

      <!-- 截图展示 -->
      <div v-if="images.length > 0" class="images-section">
        <div class="section-title">问题截图</div>
        <div class="images-grid">
          <el-image
            v-for="(img, index) in images"
            :key="index"
            :src="img"
            :preview-src-list="images"
            :initial-index="index"
            fit="cover"
            class="screenshot-image"
          />
        </div>
      </div>

      <!-- 富文本内容编辑器 -->
      <div class="editor-section">
        <div class="section-header">
          <div class="section-title">详细内容</div>
          <div class="edit-actions">
            <el-button
              v-if="!isEditing"
              type="primary"
              size="small"
              @click="startEdit"
            >
              <el-icon><Edit /></el-icon>
              编辑模式
            </el-button>
            <template v-else>
              <el-button size="small" @click="cancelEdit">取消</el-button>
              <el-button type="primary" size="small" @click="saveContent">
                <el-icon><Check /></el-icon>
                保存
              </el-button>
            </template>
          </div>
        </div>

        <!-- 富文本编辑器 -->
        <div class="editor-wrapper">
          <Toolbar
            :editor="editorRef"
            :defaultConfig="toolbarConfig"
            :mode="mode"
            class="editor-toolbar"
            :class="{ 'disabled': !isEditing }"
          />
          <Editor
            v-model="contentHtml"
            :defaultConfig="editorConfig"
            :mode="mode"
            class="editor-content"
            :class="{ 'view-mode': !isEditing }"
            @onCreated="handleCreated"
          />
        </div>
      </div>

      <!-- 标签 -->
      <div v-if="issue.tags" class="tags-section">
        <div class="section-title">标签</div>
        <div class="tags-list">
          <el-tag
            v-for="(tag, index) in issue.tags.split(',')"
            :key="index"
            size="small"
            class="tag-item"
          >
            {{ tag }}
          </el-tag>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose" size="large">关闭</el-button>
        <el-button type="warning" @click="handleEditIssue" size="large">
          <el-icon><Edit /></el-icon>
          编辑问题
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, shallowRef, watch, onBeforeUnmount } from 'vue'
import { Edit, Check } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import '@wangeditor/editor/dist/css/style.css'
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import { updateWorkIssue } from './api'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  issue: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update:modelValue', 'edit', 'refresh'])

const visible = ref(false)
const isEditing = ref(false)
const contentHtml = ref('')
const originalContent = ref('')
const editorRef = shallowRef()
const mode = 'default'

// 编辑器配置
const toolbarConfig = {
  toolbarKeys: [
    'headerSelect',
    'bold',
    'italic',
    'underline',
    'through',
    '|',
    'fontSize',
    'fontFamily',
    'lineHeight',
    '|',
    'color',
    'bgColor',
    '|',
    'bulletedList',
    'numberedList',
    'todo',
    '|',
    'emotion',
    'insertLink',
    'insertImage',
    'insertTable',
    'codeBlock',
    'divider',
    '|',
    'undo',
    'redo',
    '|',
    'fullScreen'
  ]
}

const editorConfig = {
  placeholder: '请输入问题的详细内容...',
  MENU_CONF: {
    uploadImage: {
      server: 'http://localhost:8080/api/work-issues/upload',
      fieldName: 'file',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      maxFileSize: 5 * 1024 * 1024, // 5MB
      allowedFileTypes: ['image/*'],
      customInsert(res, insertFn) {
        console.log('[富文本编辑器] 图片上传响应:', res)
        // 后端 response.Success 返回的 code 是 0，不是 200
        if (res.code === 0 || res.code == 0) {
          // 后端返回文件名，前端拼接完整URL
          const filename = res.data.filename
          const imageUrl = `http://localhost:8080/api/files/work-issues/${filename}`
          console.log('[富文本编辑器] 插入图片URL:', imageUrl)
          insertFn(imageUrl, res.data.originalName || '', imageUrl)
          ElMessage.success('图片上传成功')
        } else {
          console.error('[富文本编辑器] 图片上传失败，code:', res.code)
          ElMessage.error(res.message || '图片上传失败')
        }
      }
    }
  }
}

// 解析图片 - 将数据库中的文件名转换为完整URL
const images = ref([])
const parseImages = (imagesStr) => {
  if (!imagesStr) return []
  try {
    const parsed = JSON.parse(imagesStr)
    return Array.isArray(parsed) ? parsed.map(img => {
      // 如果已经是完整URL，直接返回
      if (img.startsWith('http')) return img
      // 如果是旧格式的相对路径（/api/files/xxx），加上域名
      if (img.startsWith('/api/')) return `http://localhost:8080${img}`
      // 如果只是文件名（1001_xxx.png），拼接完整路径
      return `http://localhost:8080/api/files/work-issues/${img}`
    }) : []
  } catch {
    return []
  }
}

// 状态和严重程度映射
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
    critical: '严重'
  }
  return map[severity] || severity
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 编辑器实例创建
const handleCreated = (editor) => {
  editorRef.value = editor
  if (!isEditing.value) {
    editor.disable()
  }
}

// 开始编辑
const startEdit = () => {
  isEditing.value = true
  originalContent.value = contentHtml.value
  if (editorRef.value) {
    editorRef.value.enable()
  }
}

// 取消编辑
const cancelEdit = () => {
  isEditing.value = false
  contentHtml.value = originalContent.value
  if (editorRef.value) {
    editorRef.value.disable()
  }
}

// 保存内容
const saveContent = async () => {
  try {
    const updateData = {
      title: props.issue.title,
      description: props.issue.description,
      status: props.issue.status,
      severity: props.issue.severity,
      solution: props.issue.solution || '',
      tags: props.issue.tags || '',
      images: props.issue.images || '',
      content: contentHtml.value
    }
    await updateWorkIssue(props.issue.id, updateData)
    ElMessage.success('保存成功')
    isEditing.value = false
    if (editorRef.value) {
      editorRef.value.disable()
    }
    emit('refresh')
  } catch (error) {
    ElMessage.error(error.message || '保存失败')
  }
}

// 编辑问题
const handleEditIssue = () => {
  emit('edit', props.issue)
  handleClose()
}

// 关闭对话框
const handleClose = () => {
  if (isEditing.value) {
    cancelEdit()
  }
  visible.value = false
  emit('update:modelValue', false)
}

// 初始化编辑器内容
const initializeContent = (issue) => {
  let content = issue.content || ''
  images.value = parseImages(issue.images)
  
  // 如果 content 为空或只是默认内容，但有上传的图片，则自动插入图片
  if ((!content || content === '<p>暂无详细内容</p>') && images.value.length > 0) {
    // 生成包含所有图片的 HTML
    const imagesHtml = images.value.map(img => `<p><img src="${img}" alt="问题截图" style="max-width: 100%;"/></p>`).join('')
    content = `<p>问题截图：</p>${imagesHtml}<p><br></p>`
  }
  
  contentHtml.value = content || '<p>暂无详细内容</p>'
  originalContent.value = contentHtml.value
}

// 监听 props 变化
watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val) {
    initializeContent(props.issue)
    isEditing.value = false
  }
})

watch(() => props.issue, (newVal) => {
  if (newVal && visible.value) {
    initializeContent(newVal)
  }
}, { deep: true })

// 组件销毁时清理编辑器
onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor) {
    editor.destroy()
  }
})
</script>

<style scoped>
/* 对话框样式 */
.issue-detail-dialog :deep(.el-dialog) {
  margin: 0 auto !important;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 12px 48px rgba(0, 0, 0, 0.2);
}

.issue-detail-dialog :deep(.el-dialog__header) {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px 24px;
  cursor: move;
  user-select: none;
}

.issue-detail-dialog :deep(.el-dialog__title) {
  color: white;
  font-weight: 600;
  font-size: 20px;
}

.issue-detail-dialog :deep(.el-dialog__headerbtn) {
  top: 16px;
  right: 16px;
}

.issue-detail-dialog :deep(.el-dialog__headerbtn .el-dialog__close) {
  color: white;
  font-size: 20px;
  font-weight: bold;
  transition: all 0.3s;
}

.issue-detail-dialog :deep(.el-dialog__headerbtn:hover .el-dialog__close) {
  color: #fbbf24;
  transform: scale(1.1);
}

.issue-detail-dialog :deep(.el-dialog__body) {
  padding: 0;
  max-height: 75vh;
  overflow-y: auto;
  background: white;
}

.editor-container {
  padding: 24px;
}

/* 信息栏 */
.info-bar {
  display: flex;
  gap: 24px;
  padding: 16px;
  background: #f8f9fa;
  border-radius: 8px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.info-item .label {
  color: #64748b;
  font-size: 14px;
}

.info-item .value {
  color: #1e293b;
  font-weight: 500;
  font-size: 14px;
}

/* 各个部分 */
.description-section,
.solution-section,
.images-section,
.editor-section,
.tags-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 2px solid #e2e8f0;
}

.description-text,
.solution-text {
  padding: 12px;
  background: #f8f9fa;
  border-radius: 6px;
  color: #475569;
  line-height: 1.6;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.solution-text {
  background: #f0fdf4;
  color: #15803d;
  border-left: 3px solid #22c55e;
}

/* 截图网格 */
.images-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
}

.screenshot-image {
  width: 100%;
  height: 200px;
  border-radius: 8px;
  cursor: pointer;
  transition: transform 0.2s;
  border: 2px solid #e2e8f0;
}

.screenshot-image:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 编辑器部分 */
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.edit-actions {
  display: flex;
  gap: 8px;
}

.editor-wrapper {
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  overflow: hidden;
}

.editor-toolbar {
  border-bottom: 1px solid #e2e8f0;
  background: #fafafa;
}

.editor-toolbar.disabled {
  opacity: 0.6;
  pointer-events: none;
  user-select: none;
}

.editor-content {
  min-height: 400px;
  max-height: 600px;
  overflow-y: auto;
}

.editor-content.view-mode {
  background: #f8f9fa;
  min-height: 200px;
}

.editor-content :deep(.w-e-text-container) {
  background: white;
  min-height: 400px;
  padding: 16px;
}

.editor-content.view-mode :deep(.w-e-text-container) {
  background: #f8f9fa;
  min-height: 200px;
  padding: 16px;
}

.editor-content :deep(.w-e-text-container p) {
  margin: 8px 0;
  line-height: 1.8;
}

.editor-content :deep(.w-e-text-container img) {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
  margin: 12px 0;
}

/* 标签 */
.tags-list {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.tag-item {
  font-size: 13px;
}

/* 底部按钮 */
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>

