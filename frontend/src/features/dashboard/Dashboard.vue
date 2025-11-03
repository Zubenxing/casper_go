<template>
  <div class="dashboard">
    <h2 class="page-title">数据概览</h2>

    <!-- 统计卡片 -->
    <el-row :gutter="20">
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: #409eff;">
            <el-icon :size="30"><Document /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.total }}</div>
            <div class="stat-label">监控总数</div>
          </div>
        </div>
      </el-col>

      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: #67c23a;">
            <el-icon :size="30"><CircleCheck /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.normal }}</div>
            <div class="stat-label">正常</div>
          </div>
        </div>
      </el-col>

      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: #e6a23c;">
            <el-icon :size="30"><Warning /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.warning }}</div>
            <div class="stat-label">警告</div>
          </div>
        </div>
      </el-col>

      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-icon" style="background: #f56c6c;">
            <el-icon :size="30"><CircleClose /></el-icon>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stats.expired }}</div>
            <div class="stat-label">过期</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 证书列表 -->
    <el-card class="cert-list-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>最近添加的证书</span>
          <el-button type="primary" size="small" @click="$router.push('/certificates')">
            查看全部
          </el-button>
        </div>
      </template>

      <el-table :data="recentCertificates" style="width: 100%" v-loading="loading">
        <el-table-column prop="url" label="URL" min-width="200" />
        <el-table-column prop="domain" label="域名" width="150" />
        <el-table-column prop="days_left" label="剩余天数" width="100">
          <template #default="{ row }">
            <el-tag :type="getDaysLeftType(row.days_left)">
              {{ row.days_left }} 天
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status_text" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ row.status_text }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_check_at" label="最后检查" width="180">
          <template #default="{ row }">
            {{ formatTime(row.last_check_at) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Document, CircleCheck, Warning, CircleClose } from '@element-plus/icons-vue'
import { formatTime, getDaysLeftType, getStatusType } from '@/core/utils/format'
import api from '@/core/api'

const loading = ref(false)
const recentCertificates = ref([])
const stats = ref({
  total: 0,
  normal: 0,
  warning: 0,
  expired: 0
})

const fetchData = async () => {
  loading.value = true
  try {
    const response = await api.certificates.getList(1, 5)
    if (response.code === 0) {
      recentCertificates.value = response.data.list
      stats.value.total = response.data.total
      
      // 计算统计数据
      stats.value.normal = recentCertificates.value.filter(c => c.status === 1).length
      stats.value.warning = recentCertificates.value.filter(c => c.status === 2).length
      stats.value.expired = recentCertificates.value.filter(c => c.status === 3).length
    }
  } catch (error) {
    console.error('获取数据失败', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.dashboard {
  width: 100%;
}

.page-title {
  margin-bottom: 20px;
  font-size: 24px;
  color: #303133;
}

.stat-card {
  background: white;
  border-radius: 8px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  transition: transform 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 5px;
}

.cert-list-card {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>

