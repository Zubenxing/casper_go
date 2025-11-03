<template>
  <div class="validate-cert-page">
    <!-- 输入表单 -->
    <el-card class="input-card">
      <template #header>
        <div class="card-header">
          <el-icon><Stamp /></el-icon>
          <span>KEY/CSR/SSL证书匹配工具</span>
        </div>
      </template>

      <el-alert
        title="验证说明"
        type="info"
        :closable="false"
        style="margin-bottom: 24px"
      >
        <div class="alert-content">
          <div>✓ 支持 PEM 格式的证书和私钥文件</div>
          <div>✓ 自动验证证书有效性、过期时间、证书链信息</div>
          <div>✓ 验证证书和私钥是否匹配配对（强烈推荐）</div>
          <div>✓ 支持拖拽上传或粘贴文件内容</div>
        </div>
      </el-alert>

      <!-- 双栏对称输入布局 -->
      <div class="dual-column-layout">
        <!-- 左侧：证书文件 -->
        <div class="input-column">
          <div class="column-header cert-header">
            <div class="header-left">
              <el-icon :size="20"><Document /></el-icon>
              <span>SSL证书文件</span>
            </div>
            <el-tag size="small" type="danger" effect="dark">必填</el-tag>
          </div>
          
          <div class="input-wrapper">
            <el-input
              v-model="certContent"
              type="textarea"
              :rows="14"
              placeholder="请粘贴SSL证书内容到此处，以 -----BEGIN CERTIFICATE----- 开头

或下方拖拽上传 .crt、.cer、.pem 文件"
              class="cert-input"
            />
          </div>

          <div class="upload-zone cert-upload">
            <el-upload
              :auto-upload="false"
              :show-file-list="false"
              accept=".crt,.cer,.pem,.txt"
              :on-change="handleCertFileChange"
              drag
            >
              <template #default>
                <el-icon class="upload-icon"><UploadFilled /></el-icon>
                <div class="upload-text">拖拽证书文件到此 或 <em>点击上传</em></div>
                <div class="upload-hint">支持 .crt, .cer, .pem 格式</div>
              </template>
            </el-upload>
          </div>
        </div>

        <!-- 右侧：私钥文件 -->
        <div class="input-column">
          <div class="column-header key-header">
            <div class="header-left">
              <el-icon :size="20"><Key /></el-icon>
              <span>私钥文件（KEY）</span>
            </div>
            <el-tag size="small" type="success" effect="plain">可选</el-tag>
          </div>
          
          <div class="input-wrapper">
            <el-input
              v-model="privateKeyContent"
              type="textarea"
              :rows="14"
              placeholder="请粘贴私钥内容到此处，以 -----BEGIN RSA PRIVATE KEY----- 开头

或下方拖拽上传 .key、.pem 文件

建议上传私钥以验证证书-私钥配对关系"
              class="key-input"
            />
          </div>

          <div class="upload-zone key-upload">
            <el-upload
              :auto-upload="false"
              :show-file-list="false"
              accept=".key,.pem,.txt"
              :on-change="handleKeyFileChange"
              drag
            >
              <template #default>
                <el-icon class="upload-icon"><UploadFilled /></el-icon>
                <div class="upload-text">拖拽私钥文件到此 或 <em>点击上传</em></div>
                <div class="upload-hint">支持 .key, .pem 格式</div>
              </template>
            </el-upload>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <el-button 
          type="primary" 
          size="large"
          @click="handleValidate" 
          :loading="validating"
          :disabled="!certContent.trim()"
        >
          <el-icon v-if="!validating"><Check /></el-icon>
          <span v-if="validating">验证中...</span>
          <span v-else-if="privateKeyContent.trim()">验证证书并校验配对</span>
          <span v-else>验证证书</span>
        </el-button>
        <el-button size="large" @click="handleClear">
          <el-icon><Delete /></el-icon>
          清空全部
        </el-button>
      </div>
    </el-card>

    <!-- 验证结果 -->
    <el-card v-if="result" class="result-card">
      <template #header>
        <div class="card-header">
          <el-icon v-if="result.valid" style="color: #67c23a" :size="24">
            <SuccessFilled />
          </el-icon>
          <el-icon v-else style="color: #f56c6c" :size="24">
            <CircleCloseFilled />
          </el-icon>
          <span>{{ result.valid ? '✓ 证书验证成功' : '✗ 验证失败' }}</span>
        </div>
      </template>

      <el-alert
        v-if="!result.valid"
        :title="result.error_message"
        type="error"
        :closable="false"
      />

      <!-- 证书-私钥配对验证结果 -->
      <el-alert
        v-if="result.valid && result.key_pair_checked"
        :type="result.key_pair_matched ? 'success' : 'error'"
        :closable="false"
        style="margin-bottom: 20px"
      >
        <template #title>
          <div style="display: flex; align-items: center; gap: 8px; font-size: 16px; font-weight: 600;">
            <el-icon :size="20">
              <component :is="result.key_pair_matched ? 'SuccessFilled' : 'CircleCloseFilled'" />
            </el-icon>
            <span>{{ result.key_pair_matched ? '证书和私钥配对成功' : '证书和私钥不匹配' }}</span>
          </div>
        </template>
        <div style="margin-top: 8px;">
          <span v-if="result.key_pair_matched">
            证书的公钥和私钥完美匹配，证书可以正常使用！
          </span>
          <span v-else>
            证书和私钥不配对，无法一起使用。请检查是否上传了正确的私钥文件。
          </span>
        </div>
      </el-alert>

      <div v-if="result.valid">
        <!-- 证书状态 -->
        <div class="status-section">
          <el-tag 
            :type="getStatusType(result)" 
            size="large" 
            effect="dark"
            style="font-size: 15px; padding: 12px 20px;"
          >
            {{ getStatusText(result) }}
          </el-tag>
        </div>

        <!-- 证书详细信息 -->
        <div class="info-section">
          <h3 class="section-title">证书详细信息</h3>
          
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="info-item">
                <span class="label">通用名称 (CN)</span>
                <span class="value">{{ result.common_name }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="info-item">
                <span class="label">组织</span>
                <span class="value">{{ result.organization || '-' }}</span>
              </div>
            </el-col>
          </el-row>

          <el-row :gutter="20">
            <el-col :span="12">
              <div class="info-item">
                <span class="label">国家</span>
                <span class="value">{{ result.country || '-' }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="info-item">
                <span class="label">证书版本</span>
                <span class="value">v{{ result.version }}</span>
              </div>
            </el-col>
          </el-row>

          <div class="info-item">
            <span class="label">序列号</span>
            <span class="value">{{ result.serial_number }}</span>
          </div>

          <div class="info-item">
            <span class="label">颁发者</span>
            <span class="value">{{ result.issuer }}</span>
          </div>

          <div class="info-item">
            <span class="label">主体</span>
            <span class="value">{{ result.subject }}</span>
          </div>
        </div>

        <!-- 有效期 -->
        <div class="info-section">
          <h3 class="section-title">有效期</h3>
          
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="info-item">
                <span class="label">生效时间</span>
                <span class="value">{{ formatDate(result.not_before) }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="info-item">
                <span class="label">过期时间</span>
                <span class="value">{{ formatDate(result.not_after) }}</span>
              </div>
            </el-col>
          </el-row>

          <div class="info-item">
            <span class="label">剩余天数</span>
            <span class="value" :style="{ color: result.days_left <= 30 ? '#f56c6c' : '#67c23a' }">
              {{ result.days_left }} 天
            </span>
          </div>
        </div>

        <!-- 域名和IP -->
        <div class="info-section">
          <h3 class="section-title">备用名称 (SAN)</h3>
          
          <div v-if="result.dns_names && result.dns_names.length > 0" class="info-item">
            <span class="label">DNS 名称</span>
            <div class="tag-list">
              <el-tag v-for="(name, index) in result.dns_names" :key="index" style="margin: 4px;">
                {{ name }}
              </el-tag>
            </div>
          </div>

          <div v-if="result.ip_addresses && result.ip_addresses.length > 0" class="info-item">
            <span class="label">IP 地址</span>
            <div class="tag-list">
              <el-tag v-for="(ip, index) in result.ip_addresses" :key="index" type="success" style="margin: 4px;">
                {{ ip }}
              </el-tag>
            </div>
          </div>

          <div v-if="result.email_addresses && result.email_addresses.length > 0" class="info-item">
            <span class="label">电子邮件</span>
            <div class="tag-list">
              <el-tag v-for="(email, index) in result.email_addresses" :key="index" type="info" style="margin: 4px;">
                {{ email }}
              </el-tag>
            </div>
          </div>
        </div>

        <!-- 签名和密钥 -->
        <div class="info-section">
          <h3 class="section-title">签名和密钥信息</h3>
          
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="info-item">
                <span class="label">签名算法</span>
                <span class="value">{{ result.signature_algorithm }}</span>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="info-item">
                <span class="label">公钥算法</span>
                <span class="value">{{ result.public_key_algorithm }}</span>
              </div>
            </el-col>
          </el-row>

          <div class="info-item">
            <span class="label">公钥长度</span>
            <span class="value">{{ result.public_key_size }} 位</span>
          </div>

          <div v-if="result.is_ca !== undefined" class="info-item">
            <span class="label">CA 证书</span>
            <span class="value">
              <el-tag :type="result.is_ca ? 'warning' : 'info'" size="small">
                {{ result.is_ca ? '是' : '否' }}
              </el-tag>
            </span>
          </div>
        </div>

        <!-- 密钥用途 -->
        <div v-if="result.key_usage && result.key_usage.length > 0" class="info-section">
          <h3 class="section-title">密钥用途</h3>
          
          <div class="tag-list">
            <el-tag v-for="(usage, index) in result.key_usage" :key="index" type="warning" style="margin: 4px;">
              {{ usage }}
            </el-tag>
          </div>
        </div>

        <!-- 扩展密钥用途 -->
        <div v-if="result.ext_key_usage && result.ext_key_usage.length > 0" class="info-section">
          <h3 class="section-title">扩展密钥用途</h3>
          
          <div class="tag-list">
            <el-tag v-for="(usage, index) in result.ext_key_usage" :key="index" type="success" style="margin: 4px;">
              {{ usage }}
            </el-tag>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { 
  Stamp, 
  Check, 
  SuccessFilled, 
  CircleCloseFilled,
  UploadFilled,
  Document,
  Key,
  Delete
} from '@element-plus/icons-vue'
import api from './api'

const certContent = ref('')
const privateKeyContent = ref('')
const validating = ref(false)
const result = ref(null)

const handleCertFileChange = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    certContent.value = e.target.result
    ElMessage.success('证书文件读取成功')
  }
  reader.onerror = () => {
    ElMessage.error('证书文件读取失败')
  }
  reader.readAsText(file.raw)
}

const handleKeyFileChange = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    privateKeyContent.value = e.target.result
    ElMessage.success('私钥文件读取成功')
  }
  reader.onerror = () => {
    ElMessage.error('私钥文件读取失败')
  }
  reader.readAsText(file.raw)
}

const handleValidate = async () => {
  if (!certContent.value.trim()) {
    ElMessage.warning('请输入证书内容')
    return
  }

  validating.value = true
  result.value = null

  try {
    const requestData = {
      cert_content: certContent.value.trim()
    }

    // 如果提供了私钥，添加到请求中
    if (privateKeyContent.value.trim()) {
      requestData.private_key_content = privateKeyContent.value.trim()
    }

    const response = await api.validateCert(requestData)

    result.value = response.data
    
    if (response.data.valid) {
      if (response.data.key_pair_checked) {
        if (response.data.key_pair_matched) {
          ElMessage.success('证书验证成功，证书和私钥配对正确！')
        } else {
          ElMessage.error('证书验证成功，但证书和私钥不匹配')
        }
      } else if (response.data.is_expired) {
        ElMessage.warning('证书验证成功，但证书已过期')
      } else if (response.data.is_not_yet_valid) {
        ElMessage.warning('证书验证成功，但证书尚未生效')
      } else {
        ElMessage.success('证书验证成功')
      }
    } else {
      ElMessage.error('证书验证失败: ' + response.data.error_message)
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
  certContent.value = ''
  privateKeyContent.value = ''
  result.value = null
  ElMessage.info('已清空所有内容')
}

const getStatusType = (result) => {
  if (result.is_expired) return 'danger'
  if (result.is_not_yet_valid) return 'warning'
  if (result.days_left <= 30) return 'warning'
  return 'success'
}

const getStatusText = (result) => {
  if (result.is_expired) return '证书已过期'
  if (result.is_not_yet_valid) return '证书尚未生效'
  if (result.days_left <= 30) return `即将过期 (${result.days_left} 天)`
  return '证书有效'
}

const formatDate = (dateStr) => {
  return new Date(dateStr).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}
</script>

<style scoped>
.validate-cert-page {
  padding: 20px;
  max-width: 1600px;
  margin: 0 auto;
}

.input-card {
  margin-bottom: 24px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.alert-content {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 14px;
  line-height: 1.6;
}

/* 双栏对称布局 */
.dual-column-layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  margin-bottom: 24px;
}

.input-column {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.column-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-radius: 8px;
  font-size: 15px;
  font-weight: 600;
}

.cert-header {
  background: linear-gradient(135deg, #409eff 0%, #53a8ff 100%);
  color: white;
}

.key-header {
  background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%);
  color: white;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.input-wrapper {
  flex: 1;
}

.cert-input :deep(.el-textarea__inner),
.key-input :deep(.el-textarea__inner) {
  font-family: 'Courier New', Consolas, monospace;
  font-size: 13px;
  border-radius: 8px;
  border: 2px solid #dcdfe6;
  transition: all 0.3s;
  line-height: 1.6;
}

.cert-input :deep(.el-textarea__inner) {
  background: #f0f9ff;
}

.cert-input :deep(.el-textarea__inner):focus {
  border-color: #409eff;
  box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.1);
  background: white;
}

.key-input :deep(.el-textarea__inner) {
  background: #f0f9f0;
}

.key-input :deep(.el-textarea__inner):focus {
  border-color: #67c23a;
  box-shadow: 0 0 0 3px rgba(103, 194, 58, 0.1);
  background: white;
}

.upload-zone {
  height: 100px;
}

.upload-zone :deep(.el-upload) {
  width: 100%;
}

.upload-zone :deep(.el-upload-dragger) {
  padding: 16px;
  height: 100%;
  width: 100%;
  border-radius: 8px;
  border: 2px dashed #dcdfe6;
  transition: all 0.3s;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.cert-upload :deep(.el-upload-dragger):hover {
  border-color: #409eff;
  background-color: #f0f9ff;
}

.key-upload :deep(.el-upload-dragger):hover {
  border-color: #67c23a;
  background-color: #f0f9f0;
}

.upload-icon {
  font-size: 28px;
  color: #409eff;
  margin-bottom: 4px;
}

.upload-text {
  font-size: 14px;
  color: #606266;
  margin-bottom: 4px;
}

.upload-text em {
  color: #409eff;
  font-style: normal;
  font-weight: 500;
}

.upload-hint {
  font-size: 12px;
  color: #909399;
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  justify-content: center;
  gap: 16px;
  padding-top: 8px;
}

.action-buttons .el-button {
  min-width: 200px;
  font-size: 15px;
  height: 48px;
  font-weight: 500;
}

/* 结果卡片 */
.result-card {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.status-section {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 0;
  margin-bottom: 24px;
  border-bottom: 1px solid #e4e7ed;
}

.info-section {
  margin-bottom: 24px;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 8px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin: 0 0 16px 0;
  padding-bottom: 8px;
  border-bottom: 2px solid #409eff;
}

.info-item {
  display: flex;
  align-items: flex-start;
  padding: 12px 0;
  border-bottom: 1px solid #ebeef5;
}

.info-item:last-child {
  border-bottom: none;
}

.info-item .label {
  flex-shrink: 0;
  width: 140px;
  font-weight: 600;
  color: #606266;
  font-size: 14px;
}

.info-item .value {
  flex: 1;
  color: #303133;
  word-break: break-all;
  font-size: 14px;
}

.tag-list {
  flex: 1;
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .dual-column-layout {
    grid-template-columns: 1fr;
  }
}
</style>
