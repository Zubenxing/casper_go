<template>
  <div class="cert-update">
    <h2 class="page-title">更新证书</h2>

    <el-card shadow="never">
      <el-alert
        title="批量更新证书信息"
        type="info"
        :closable="false"
        style="margin-bottom: 20px;"
      >
        <p>此功能用于批量更新证书的最新状态信息</p>
        <p>• 支持选择特定证书或全部证书</p>
        <p>• 支持按状态筛选（即将过期、已过期等）</p>
        <p>• 可设置更新频率和通知方式</p>
      </el-alert>

      <el-row :gutter="20">
        <el-col :span="24">
          <el-card shadow="never" class="function-card">
            <template #header>
              <div class="card-header">
                <span>快速操作</span>
              </div>
            </template>
            
            <el-space wrap size="large">
              <el-button type="primary" size="large" :icon="Refresh" @click="updateAll" :loading="updating">
                更新所有证书
              </el-button>
              
              <el-button type="warning" size="large" :icon="Warning" @click="updateWarning" :loading="updatingWarning">
                更新即将过期证书
              </el-button>
              
              <el-button type="danger" size="large" :icon="CircleClose" @click="updateExpired" :loading="updatingExpired">
                更新已过期证书
              </el-button>
            </el-space>
          </el-card>
        </el-col>
      </el-row>

      <!-- 更新进度 -->
      <el-card v-if="showProgress" shadow="never" style="margin-top: 20px;">
        <template #header>
          <div class="card-header">
            <span>更新进度</span>
          </div>
        </template>
        
        <el-progress 
          :percentage="progress" 
          :status="progressStatus"
          :stroke-width="20"
        />
        
        <div style="margin-top: 15px; color: #606266;">
          <p>已处理：{{ processedCount }} / {{ totalCount }}</p>
          <p v-if="currentUrl">当前：{{ currentUrl }}</p>
        </div>
      </el-card>

      <!-- 更新结果 -->
      <el-card v-if="updateResult" shadow="never" style="margin-top: 20px;">
        <template #header>
          <div class="card-header">
            <span>更新结果</span>
          </div>
        </template>
        
        <el-result
          :icon="updateResult.success ? 'success' : 'warning'"
          :title="updateResult.title"
          :sub-title="updateResult.message"
        >
          <template #extra>
            <el-button type="primary" @click="closeResult">知道了</el-button>
          </template>
        </el-result>
      </el-card>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh, Warning, CircleClose } from '@element-plus/icons-vue'
import api from '@/core/api'

const updating = ref(false)
const updatingWarning = ref(false)
const updatingExpired = ref(false)
const showProgress = ref(false)
const progress = ref(0)
const progressStatus = ref('')
const processedCount = ref(0)
const totalCount = ref(0)
const currentUrl = ref('')
const updateResult = ref(null)

// 更新所有证书
const updateAll = async () => {
  updating.value = true
  showProgress.value = true
  progress.value = 0
  
  try {
    const response = await api.certificates.checkAll()
    if (response.code === 0) {
      progress.value = 100
      progressStatus.value = 'success'
      
      updateResult.value = {
        success: true,
        title: '更新完成',
        message: '所有证书已成功更新'
      }
      
      ElMessage.success('所有证书更新完成')
    }
  } catch (error) {
    progressStatus.value = 'exception'
    updateResult.value = {
      success: false,
      title: '更新失败',
      message: error.response?.data?.message || error.message
    }
  } finally {
    updating.value = false
  }
}

// 更新即将过期证书
const updateWarning = () => {
  ElMessage.info('此功能开发中...')
  // TODO: 后端添加按状态更新的API
}

// 更新已过期证书
const updateExpired = () => {
  ElMessage.info('此功能开发中...')
  // TODO: 后端添加按状态更新的API
}

// 关闭结果
const closeResult = () => {
  updateResult.value = null
  showProgress.value = false
}
</script>

<style scoped>
.cert-update {
  width: 100%;
}

.page-title {
  margin-bottom: 20px;
  font-size: 24px;
  color: #303133;
}

.function-card {
  background: #f5f7fa;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}
</style>

