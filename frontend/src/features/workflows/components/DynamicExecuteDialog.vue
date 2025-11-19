<template>
  <el-dialog
    v-model="visible"
    :title="`执行工作流：${config?.workflow?.name || ''}`"
    width="600px"
    @close="handleClose"
  >
    <div v-if="config && config.workflow" class="dialog-content">
      <el-alert
        :title="config.description || '准备执行工作流'"
        type="info"
        :closable="false"
        style="margin-bottom: 20px"
      />

      <!-- 简单执行类型：无需输入 -->
      <div v-if="config.type === 'simple'" class="simple-execute">
        <el-result icon="success" title="准备就绪">
          <template #sub-title>
            点击下方按钮即可执行工作流
          </template>
        </el-result>
      </div>

      <!-- 文件上传类型 -->
      <div v-else-if="config.type === 'upload_files'" class="file-upload">
        <el-form :model="formData" label-width="100px">
          <el-form-item
            v-for="field in config.fields"
            :key="field.name"
            :label="field.label"
            :required="field.required"
          >
            <el-upload
              :ref="field.name"
              class="upload-demo"
              drag
              :auto-upload="false"
              :limit="1"
              :accept="field.accept"
              :on-change="(file) => handleFileChange(field.name, file)"
              :on-exceed="handleExceed"
            >
              <el-icon class="el-icon--upload"><upload-filled /></el-icon>
              <div class="el-upload__text">
                拖拽文件到这里 或 <em>点击上传</em>
              </div>
              <template #tip>
                <div class="el-upload__tip">
                  {{ field.description || `支持 ${field.accept} 格式` }}
                </div>
              </template>
            </el-upload>
            <div v-if="formData[field.name]" class="file-info">
              <el-tag type="success">{{ formData[field.name].name }}</el-tag>
            </div>
          </el-form-item>
        </el-form>
      </div>

      <!-- 表单输入类型 -->
      <div v-else-if="config.type === 'form_input'" class="form-input">
        <el-form :model="formData" label-width="100px">
          <el-form-item
            v-for="field in config.fields"
            :key="field.name"
            :label="field.label"
            :required="field.required"
          >
            <!-- 日期选择器 -->
            <el-date-picker
              v-if="field.type === 'date'"
              v-model="formData[field.name]"
              type="date"
              :placeholder="field.placeholder || '请选择日期'"
              style="width: 100%"
            />
            
            <!-- 文本输入 -->
            <el-input
              v-else-if="field.type === 'text'"
              v-model="formData[field.name]"
              :placeholder="field.placeholder || `请输入${field.label}`"
            />
            
            <!-- 文本域 -->
            <el-input
              v-else-if="field.type === 'textarea'"
              v-model="formData[field.name]"
              type="textarea"
              :rows="4"
              :placeholder="field.placeholder || `请输入${field.label}`"
            />
            
            <!-- 数字输入 -->
            <el-input-number
              v-else-if="field.type === 'number'"
              v-model="formData[field.name]"
              :placeholder="field.placeholder"
              style="width: 100%"
            />
            
            <!-- 下拉选择 -->
            <el-select
              v-else-if="field.type === 'select'"
              v-model="formData[field.name]"
              :placeholder="field.placeholder || `请选择${field.label}`"
              style="width: 100%"
            >
              <el-option
                v-for="option in field.options"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              />
            </el-select>
          </el-form-item>
        </el-form>
      </div>

      <!-- Webhook 类型 -->
      <div v-else-if="config.type === 'webhook'" class="webhook-execute">
        <el-result icon="info" title="Webhook 调用">
          <template #sub-title>
            <p>此工作流将通过 Webhook 直接调用</p>
            <p v-if="config.webhookUrl" class="webhook-url">
              <el-tag size="small">{{ config.webhookUrl }}</el-tag>
            </p>
          </template>
        </el-result>
        <el-form v-if="config.fields && config.fields.length > 0" :model="formData" label-width="100px">
          <el-form-item
            v-for="field in config.fields"
            :key="field.name"
            :label="field.label"
            :required="field.required"
          >
            <el-input
              v-model="formData[field.name]"
              :placeholder="field.placeholder || `请输入${field.label}`"
            />
          </el-form-item>
        </el-form>
      </div>

      <!-- 批量处理类型 -->
      <div v-else-if="config.type === 'batch'" class="batch-execute">
        <el-result icon="warning" title="批量处理">
          <template #sub-title>
            <p>此工作流将批量处理数据</p>
            <p v-if="config.batchSize" class="batch-size">
              每批处理 {{ config.batchSize }} 条数据
            </p>
          </template>
        </el-result>
        <el-form :model="formData" label-width="100px">
          <el-form-item label="数据文件" required>
            <el-upload
              ref="batchFileUpload"
              class="upload-demo"
              drag
              :auto-upload="false"
              :limit="1"
              accept=".json,.csv"
              :on-change="handleBatchFileChange"
            >
              <el-icon class="el-icon--upload"><upload-filled /></el-icon>
              <div class="el-upload__text">
                拖拽 JSON/CSV 文件到这里 或 <em>点击上传</em>
              </div>
              <template #tip>
                <div class="el-upload__tip">
                  支持 JSON 和 CSV 格式，每行一个数据记录
                </div>
              </template>
            </el-upload>
          </el-form-item>
        </el-form>
      </div>
    </div>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button
        type="primary"
        :loading="executing"
        :disabled="!canExecute"
        @click="handleExecute"
      >
        {{ config.executeLabel || '执行' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: Boolean,
  config: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['update:modelValue', 'execute'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const formData = ref({})
const executing = ref(false)

// 是否可以执行
const canExecute = computed(() => {
  if (props.config.type === 'simple') {
    return true
  }
  
  if (props.config.type === 'batch') {
    return formData.value.batchData && formData.value.batchData.length > 0
  }
  
  if (props.config.type === 'webhook') {
    // Webhook 类型如果没有配置字段，则总是可执行
    if (!props.config.fields || props.config.fields.length === 0) {
      return true
    }
    // 如果有字段，检查必填字段
    return props.config.fields
      .filter(f => f.required)
      .every(f => formData.value[f.name])
  }
  
  if (props.config.fields) {
    // 检查所有必填字段是否已填写
    return props.config.fields
      .filter(f => f.required)
      .every(f => formData.value[f.name])
  }
  
  return true
})

// 监听对话框打开，重置表单
watch(visible, (newVal) => {
  if (newVal) {
    formData.value = {}
  }
})

const handleFileChange = (fieldName, file) => {
  formData.value[fieldName] = file.raw
}

const handleBatchFileChange = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      let batchData = []
      const content = e.target.result
      
      if (file.name.endsWith('.json')) {
        // 解析 JSON 文件
        batchData = JSON.parse(content)
        if (!Array.isArray(batchData)) {
          throw new Error('JSON 文件必须包含数组')
        }
      } else if (file.name.endsWith('.csv')) {
        // 解析 CSV 文件
        const lines = content.split('\n').filter(line => line.trim())
        const headers = lines[0].split(',').map(h => h.trim())
        
        batchData = lines.slice(1).map(line => {
          const values = line.split(',').map(v => v.trim())
          const item = {}
          headers.forEach((header, index) => {
            item[header] = values[index] || ''
          })
          return item
        })
      }
      
      formData.value.batchData = batchData
      formData.value.batchFileName = file.name
      ElMessage.success(`成功加载 ${batchData.length} 条数据`)
    } catch (error) {
      ElMessage.error(`文件解析失败: ${error.message}`)
      formData.value.batchData = []
    }
  }
  reader.readAsText(file.raw)
}

const handleExceed = () => {
  ElMessage.warning('只能上传一个文件')
}

const handleExecute = () => {
  if (!canExecute.value) {
    ElMessage.warning('请完成所有必填项')
    return
  }
  
  executing.value = true
  emit('execute', formData.value, () => {
    executing.value = false
  })
}

const handleClose = () => {
  if (!executing.value) {
    visible.value = false
  }
}
</script>

<style scoped>
.dialog-content {
  padding: 10px 0;
}

.simple-execute {
  text-align: center;
  padding: 20px 0;
}

.webhook-execute {
  text-align: center;
  padding: 10px 0;
}

.webhook-url {
  margin-top: 10px;
}

.batch-execute {
  text-align: center;
  padding: 10px 0;
}

.batch-size {
  margin-top: 10px;
  color: #e6a23c;
}

.file-info {
  margin-top: 10px;
}

.upload-demo {
  width: 100%;
}

:deep(.el-result__title) {
  margin-top: 15px;
}

:deep(.el-result__subtitle) {
  margin-top: 10px;
  line-height: 1.6;
}
</style>
