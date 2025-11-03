<template>
  <div class="generate-csr-container">
    <el-card class="form-card">
      <template #header>
        <div class="card-header">
          <el-icon><Document /></el-icon>
          <span>生成证书签名请求 (CSR)</span>
        </div>
      </template>

      <el-form :model="form" :rules="rules" ref="formRef" label-width="140px">
        <el-divider content-position="left">基本信息</el-divider>
        
        <el-form-item label="通用名 (CN)" prop="common_name">
          <el-input 
            v-model="form.common_name" 
            placeholder="例如: example.com"
            clearable
          />
          <span class="form-tip">域名或服务器名称（必填）</span>
        </el-form-item>

        <el-form-item label="国家代码 (C)" prop="country">
          <el-input 
            v-model="form.country" 
            placeholder="例如: CN, US"
            maxlength="2"
            clearable
          />
          <span class="form-tip">两位国家代码，如 CN（中国）、US（美国）</span>
        </el-form-item>

        <el-form-item label="省份/州 (ST)" prop="province">
          <el-input 
            v-model="form.province" 
            placeholder="例如: Beijing, California"
            clearable
          />
        </el-form-item>

        <el-form-item label="城市 (L)" prop="locality">
          <el-input 
            v-model="form.locality" 
            placeholder="例如: Beijing, San Francisco"
            clearable
          />
        </el-form-item>

        <el-form-item label="组织名称 (O)" prop="organization">
          <el-input 
            v-model="form.organization" 
            placeholder="例如: Example Corp"
            clearable
          />
        </el-form-item>

        <el-form-item label="部门 (OU)" prop="organizational_unit">
          <el-input 
            v-model="form.organizational_unit" 
            placeholder="例如: IT Department"
            clearable
          />
        </el-form-item>

        <el-form-item label="邮箱地址" prop="email_address">
          <el-input 
            v-model="form.email_address" 
            placeholder="例如: admin@example.com"
            type="email"
            clearable
          />
        </el-form-item>

        <el-divider content-position="left">Subject Alternative Names (SAN)</el-divider>

        <el-form-item label="DNS 名称" prop="dns_names">
          <el-select
            v-model="form.dns_names"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="输入域名后按回车添加"
            style="width: 100%"
          >
          </el-select>
          <span class="form-tip">多个域名，例如: example.com, www.example.com</span>
        </el-form-item>

        <el-divider content-position="left">密钥配置</el-divider>

        <el-form-item label="密钥算法" prop="key_algorithm">
          <el-radio-group v-model="form.key_algorithm">
            <el-radio label="RSA">RSA</el-radio>
            <el-radio label="ECDSA">ECDSA</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="密钥长度" prop="key_size">
          <el-select v-model="form.key_size" placeholder="选择密钥长度">
            <template v-if="form.key_algorithm === 'RSA'">
              <el-option label="2048 位（推荐）" :value="2048" />
              <el-option label="4096 位（更安全）" :value="4096" />
            </template>
            <template v-else>
              <el-option label="P-256（推荐）" :value="256" />
              <el-option label="P-384（更安全）" :value="384" />
            </template>
          </el-select>
          <span class="form-tip">
            {{ form.key_algorithm === 'RSA' ? 'RSA 推荐 2048 位或更高' : 'ECDSA 推荐 P-256 或 P-384' }}
          </span>
        </el-form-item>

        <el-form-item>
          <el-button 
            type="primary" 
            @click="handleGenerate" 
            :loading="generating"
            :icon="Check"
          >
            生成 CSR
          </el-button>
          <el-button @click="handleReset">重置表单</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 生成结果 -->
    <el-card v-if="result" class="result-card">
      <template #header>
        <div class="card-header">
          <el-icon><SuccessFilled /></el-icon>
          <span>生成成功</span>
        </div>
      </template>

      <el-alert
        title="重要提示"
        type="warning"
        :closable="false"
        style="margin-bottom: 20px"
      >
        <template #default>
          私钥仅显示一次，请务必妥善保存！丢失私钥将无法使用对应的证书。
        </template>
      </el-alert>

      <div class="result-section">
        <div class="section-header">
          <h4>私钥 (Private Key)</h4>
          <el-button 
            size="small" 
            type="primary"
            @click="downloadFile(result.private_key, 'private.key')"
            :icon="Download"
          >
            下载私钥
          </el-button>
        </div>
        <el-input
          v-model="result.private_key"
          type="textarea"
          :rows="10"
          readonly
        />
      </div>

      <div class="result-section">
        <div class="section-header">
          <h4>证书签名请求 (CSR)</h4>
          <el-button 
            size="small" 
            type="primary"
            @click="downloadFile(result.csr, 'request.csr')"
            :icon="Download"
          >
            下载 CSR
          </el-button>
        </div>
        <el-input
          v-model="result.csr"
          type="textarea"
          :rows="10"
          readonly
        />
      </div>

      <el-divider />

      <div class="next-steps">
        <h4>📝 下一步操作：</h4>
        <ol>
          <li>下载并妥善保存私钥文件（<code>private.key</code>）</li>
          <li>下载 CSR 文件（<code>request.csr</code>）</li>
          <li>将 CSR 文件提交给证书颁发机构（CA）申请 SSL 证书</li>
          <li>收到证书后，使用私钥和证书配置您的服务器</li>
        </ol>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Document, Check, SuccessFilled, Download } from '@element-plus/icons-vue'
import api from './api'

const formRef = ref(null)
const generating = ref(false)
const result = ref(null)

const form = reactive({
  common_name: '',
  country: '',
  province: '',
  locality: '',
  organization: '',
  organizational_unit: '',
  email_address: '',
  dns_names: [],
  key_algorithm: 'RSA',
  key_size: 2048
})

const rules = {
  common_name: [
    { required: true, message: '请输入通用名（域名）', trigger: 'blur' }
  ],
  key_algorithm: [
    { required: true, message: '请选择密钥算法', trigger: 'change' }
  ],
  key_size: [
    { required: true, message: '请选择密钥长度', trigger: 'change' }
  ]
}

// 监听算法变化，自动调整密钥长度
watch(() => form.key_algorithm, (newVal) => {
  if (newVal === 'RSA') {
    form.key_size = 2048
  } else {
    form.key_size = 256
  }
})

const handleGenerate = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    generating.value = true
    result.value = null

    try {
      const response = await api.generateCSR({
        common_name: form.common_name,
        country: form.country || undefined,
        province: form.province || undefined,
        locality: form.locality || undefined,
        organization: form.organization || undefined,
        organizational_unit: form.organizational_unit || undefined,
        email_address: form.email_address || undefined,
        dns_names: form.dns_names.length > 0 ? form.dns_names : undefined,
        key_algorithm: form.key_algorithm,
        key_size: form.key_size
      })

      result.value = response.data
      ElMessage.success('CSR 生成成功！请妥善保存私钥。')
      
      // 滚动到结果区域
      setTimeout(() => {
        document.querySelector('.result-card')?.scrollIntoView({ behavior: 'smooth' })
      }, 100)
    } catch (error) {
      ElMessage.error(error.response?.data?.message || '生成失败，请检查输入信息')
    } finally {
      generating.value = false
    }
  })
}

const handleReset = () => {
  formRef.value?.resetFields()
  result.value = null
}

const downloadFile = (content, filename) => {
  const blob = new Blob([content], { type: 'text/plain' })
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  window.URL.revokeObjectURL(url)
  ElMessage.success(`已下载 ${filename}`)
}
</script>

<style scoped>
.generate-csr-container {
  max-width: 1200px;
  margin: 0 auto;
}

.form-card,
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

.form-tip {
  font-size: 12px;
  color: #909399;
  display: block;
  margin-top: 4px;
}

.result-section {
  margin-bottom: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.section-header h4 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
}

.next-steps {
  background: #f5f7fa;
  padding: 16px;
  border-radius: 4px;
}

.next-steps h4 {
  margin: 0 0 12px 0;
  font-size: 14px;
  color: #303133;
}

.next-steps ol {
  margin: 0;
  padding-left: 24px;
  color: #606266;
}

.next-steps li {
  margin-bottom: 8px;
  line-height: 1.6;
}

.next-steps code {
  background: #e6e8eb;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
}

:deep(.el-divider__text) {
  font-weight: 600;
  color: #303133;
}
</style>

