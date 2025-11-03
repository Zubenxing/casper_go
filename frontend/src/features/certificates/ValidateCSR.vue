<template>
  <div class="validate-csr-container">
    <el-card class="input-card">
      <template #header>
        <div class="card-header">
          <el-icon><DocumentChecked /></el-icon>
          <span>验证证书签名请求 (CSR)</span>
        </div>
      </template>

      <el-alert
        title="支持的格式"
        type="info"
        :closable="false"
        style="margin-bottom: 20px"
      >
        请粘贴 PEM 格式的 CSR 内容，以 <code>-----BEGIN CERTIFICATE REQUEST-----</code> 开头
      </el-alert>

      <el-form>
        <el-form-item>
          <el-input
            v-model="csrContent"
            type="textarea"
            :rows="12"
            placeholder="粘贴 CSR 内容到此处...

-----BEGIN CERTIFICATE REQUEST-----
MIIDCDCCAfACAQAwdjELMAkGA1UEBhMCQ04...
-----END CERTIFICATE REQUEST-----"
          />
        </el-form-item>

        <el-form-item>
          <div class="upload-section">
            <el-upload
              :auto-upload="false"
              :show-file-list="false"
              :on-change="handleFileChange"
              accept=".csr,.pem,.txt"
              drag
            >
              <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
              <div class="el-upload__text">
                拖拽文件到此处，或 <em>点击选择文件</em>
              </div>
              <template #tip>
                <div class="el-upload__tip">
                  支持 .csr, .pem, .txt 格式文件
                </div>
              </template>
            </el-upload>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button 
            type="primary" 
            @click="handleValidate" 
            :loading="validating"
            :disabled="!csrContent.trim()"
            :icon="Check"
          >
            验证 CSR
          </el-button>
          <el-button @click="handleClear">清空</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 验证结果 -->
    <el-card v-if="result" class="result-card">
      <template #header>
        <div class="card-header">
          <el-icon v-if="result.valid" style="color: #67c23a">
            <SuccessFilled />
          </el-icon>
          <el-icon v-else style="color: #f56c6c">
            <CircleCloseFilled />
          </el-icon>
          <span>{{ result.valid ? '验证成功' : '验证失败' }}</span>
        </div>
      </template>

      <el-alert
        v-if="!result.valid"
        :title="result.error_message"
        type="error"
        :closable="false"
      />

      <div v-if="result.valid" class="info-grid">
        <div class="info-item">
          <label>通用名 (CN)</label>
          <div class="info-value">{{ result.common_name || '-' }}</div>
        </div>

        <div class="info-item">
          <label>国家 (C)</label>
          <div class="info-value">{{ result.country || '-' }}</div>
        </div>

        <div class="info-item">
          <label>省份/州 (ST)</label>
          <div class="info-value">{{ result.province || '-' }}</div>
        </div>

        <div class="info-item">
          <label>城市 (L)</label>
          <div class="info-value">{{ result.locality || '-' }}</div>
        </div>

        <div class="info-item">
          <label>组织 (O)</label>
          <div class="info-value">{{ result.organization || '-' }}</div>
        </div>

        <div class="info-item">
          <label>部门 (OU)</label>
          <div class="info-value">{{ result.organizational_unit || '-' }}</div>
        </div>

        <div class="info-item full-width" v-if="result.dns_names && result.dns_names.length > 0">
          <label>DNS 名称 (SAN)</label>
          <div class="info-value">
            <el-tag 
              v-for="dns in result.dns_names" 
              :key="dns" 
              size="small" 
              type="info"
              style="margin-right: 8px; margin-bottom: 8px"
            >
              {{ dns }}
            </el-tag>
          </div>
        </div>

        <div class="info-item full-width" v-if="result.email_addresses && result.email_addresses.length > 0">
          <label>邮箱地址</label>
          <div class="info-value">
            <el-tag 
              v-for="email in result.email_addresses" 
              :key="email" 
              size="small"
              style="margin-right: 8px; margin-bottom: 8px"
            >
              {{ email }}
            </el-tag>
          </div>
        </div>

        <div class="info-item">
          <label>公钥算法</label>
          <div class="info-value">
            <el-tag type="success">{{ result.public_key_algorithm }}</el-tag>
          </div>
        </div>

        <div class="info-item">
          <label>公钥长度</label>
          <div class="info-value">{{ result.public_key_size }} 位</div>
        </div>

        <div class="info-item">
          <label>签名算法</label>
          <div class="info-value">{{ result.signature_algorithm }}</div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  DocumentChecked, 
  Check, 
  SuccessFilled, 
  CircleCloseFilled,
  UploadFilled 
} from '@element-plus/icons-vue'
import api from './api'

const csrContent = ref('')
const validating = ref(false)
const result = ref(null)

const handleFileChange = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    csrContent.value = e.target.result
    ElMessage.success('文件加载成功')
  }
  reader.onerror = () => {
    ElMessage.error('文件读取失败')
  }
  reader.readAsText(file.raw)
}

const handleValidate = async () => {
  if (!csrContent.value.trim()) {
    ElMessage.warning('请输入 CSR 内容')
    return
  }

  validating.value = true
  result.value = null

  try {
    const response = await api.validateCSR({
      csr_content: csrContent.value.trim()
    })

    result.value = response.data
    
    if (response.data.valid) {
      ElMessage.success('CSR 验证成功')
    } else {
      ElMessage.error('CSR 验证失败: ' + response.data.error_message)
    }

    // 滚动到结果区域
    setTimeout(() => {
      document.querySelector('.result-card')?.scrollIntoView({ behavior: 'smooth' })
    }, 100)
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '验证请求失败')
  } finally {
    validating.value = false
  }
}

const handleClear = () => {
  csrContent.value = ''
  result.value = null
}
</script>

<style scoped>
.validate-csr-container {
  max-width: 1200px;
  margin: 0 auto;
}

.input-card,
.result-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
}

.upload-section {
  width: 100%;
  margin-top: 16px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
  margin-top: 20px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-item.full-width {
  grid-column: 1 / -1;
}

.info-item label {
  font-size: 13px;
  color: #909399;
  font-weight: 500;
}

.info-value {
  font-size: 14px;
  color: #303133;
  word-break: break-all;
}

code {
  background: #e6e8eb;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

:deep(.el-upload-dragger) {
  padding: 30px;
}

:deep(.el-icon--upload) {
  font-size: 40px;
  color: #409eff;
  margin-bottom: 16px;
}
</style>

