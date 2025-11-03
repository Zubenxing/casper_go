<template>
  <div class="profile">
    <h2 class="page-title">个人中心</h2>

    <el-card shadow="never">
      <el-descriptions :column="1" border v-if="authStore.user">
        <el-descriptions-item label="用户ID">
          {{ authStore.user.id }}
        </el-descriptions-item>
        <el-descriptions-item label="用户名">
          {{ authStore.user.username }}
        </el-descriptions-item>
        <el-descriptions-item label="昵称">
          {{ authStore.user.nickname }}
        </el-descriptions-item>
        <el-descriptions-item label="邮箱">
          {{ authStore.user.email }}
        </el-descriptions-item>
        <el-descriptions-item label="角色">
          <el-tag :type="authStore.user.role === 'admin' ? 'danger' : 'primary'">
            {{ authStore.user.role === 'admin' ? '管理员' : '普通用户' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="authStore.user.status === 1 ? 'success' : 'danger'">
            {{ authStore.user.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="注册时间">
          {{ formatTime(authStore.user.created_at) }}
        </el-descriptions-item>
      </el-descriptions>

      <div class="profile-actions">
        <el-button type="primary" :icon="Refresh" @click="handleRefresh">
          刷新信息
        </el-button>
        <el-button type="danger" :icon="SwitchButton" @click="handleLogout">
          退出登录
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, SwitchButton } from '@element-plus/icons-vue'
import { formatTime } from '@/core/utils/format'
import { useAuthStore } from './store'

const router = useRouter()
const authStore = useAuthStore()

const handleRefresh = async () => {
  try {
    await authStore.fetchUserInfo()
    ElMessage.success('信息已刷新')
  } catch (error) {
    ElMessage.error('刷新失败')
  }
}

const handleLogout = () => {
  ElMessageBox.confirm('确定要退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    await authStore.logout()
    ElMessage.success('已退出登录')
    router.push('/login')
  }).catch(() => {})
}
</script>

<style scoped>
.profile {
  width: 100%;
  max-width: 800px;
}

.page-title {
  margin-bottom: 20px;
  font-size: 24px;
  color: #303133;
}

.profile-actions {
  margin-top: 30px;
  display: flex;
  gap: 10px;
}
</style>

