<template>
  <div class="works-container">
    <!-- 统计卡片 -->
    <div class="stats-overview">
      <el-row :gutter="16">
        <el-col :span="6">
          <el-card class="stat-card stat-total" shadow="hover">
            <div class="stat-content">
              <el-icon class="stat-icon" :size="36"><Tickets /></el-icon>
              <div class="stat-info">
                <div class="stat-value">{{ currentStats.total || 0 }}</div>
                <div class="stat-label">全部{{ activeTab === 'works' ? '任务' : '问题' }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card stat-pending" shadow="hover">
            <div class="stat-content">
              <el-icon class="stat-icon" :size="36"><Clock /></el-icon>
              <div class="stat-info">
                <div class="stat-value">{{ currentStats.pending || currentStats.open || 0 }}</div>
                <div class="stat-label">{{ activeTab === 'works' ? '待办' : '待处理' }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card stat-progress" shadow="hover">
            <div class="stat-content">
              <el-icon class="stat-icon" :size="36"><Loading /></el-icon>
              <div class="stat-info">
                <div class="stat-value">{{ currentStats.in_progress || 0 }}</div>
                <div class="stat-label">{{ activeTab === 'works' ? '进行中' : '处理中' }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card class="stat-card stat-completed" shadow="hover">
            <div class="stat-content">
              <el-icon class="stat-icon" :size="36"><CircleCheck /></el-icon>
              <div class="stat-info">
                <div class="stat-value">{{ currentStats.completed || currentStats.resolved || 0 }}</div>
                <div class="stat-label">{{ activeTab === 'works' ? '已完成' : '已解决' }}</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 内容区域 -->
    <div class="content-area">
      <WorkList v-if="activeTab === 'works'" ref="workListRef" :hide-stats="true" @stats-update="handleWorkStatsUpdate" />
      <IssueList v-else-if="activeTab === 'issues'" ref="issueListRef" :hide-stats="true" @stats-update="handleIssueStatsUpdate" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, inject, watch, onMounted } from 'vue'
import { Tickets, Clock, Loading, CircleCheck } from '@element-plus/icons-vue'
import WorkList from './WorkList.vue'
import IssueList from './IssueList.vue'
import { getWorkStats, getWorkIssueStats } from './api'

// 从 MainLayout 注入的 tab 状态
const activeTab = inject('certificateTab', ref('works'))

const workListRef = ref(null)
const issueListRef = ref(null)

const workStats = ref({
  total: 0,
  pending: 0,
  in_progress: 0,
  completed: 0
})

const issueStats = ref({
  total: 0,
  open: 0,
  in_progress: 0,
  resolved: 0
})

// 根据当前 tab 返回对应的统计数据
const currentStats = computed(() => {
  return activeTab.value === 'works' ? workStats.value : issueStats.value
})

const loadWorkStats = async () => {
  try {
    const workRes = await getWorkStats()
    workStats.value = workRes.data || {}
  } catch (error) {
    console.error('加载工作统计信息失败:', error)
  }
}

const loadIssueStats = async () => {
  try {
    const issueRes = await getWorkIssueStats()
    issueStats.value = issueRes.data || {}
  } catch (error) {
    console.error('加载问题统计信息失败:', error)
  }
}

// 处理工作统计更新
const handleWorkStatsUpdate = (stats) => {
  workStats.value = stats
}

// 处理问题统计更新
const handleIssueStatsUpdate = (stats) => {
  issueStats.value = stats
}

// 监听 tab 变化，加载对应的统计数据
watch(activeTab, (newTab) => {
  if (newTab === 'works') {
    loadWorkStats()
  } else if (newTab === 'issues') {
    loadIssueStats()
  }
}, { immediate: true })

onMounted(() => {
  console.log('工作记录页面加载完成')
})

// 暴露刷新方法给子组件调用
defineExpose({
  refreshStats: () => {
    loadWorkStats()
    loadIssueStats()
  }
})
</script>

<style scoped>
.works-container {
  padding: 0;
  background: #ffffff;
}

/* 统计卡片 */
.stats-overview {
  margin-bottom: 24px;
}

.stat-card {
  border: none;
  border-radius: 8px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.12);
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
}

.stat-total .stat-icon {
  background: linear-gradient(135deg, rgba(64, 158, 255, 0.1) 0%, rgba(64, 158, 255, 0.05) 100%);
  color: #409eff;
}

.stat-pending .stat-icon {
  background: linear-gradient(135deg, rgba(230, 162, 60, 0.1) 0%, rgba(230, 162, 60, 0.05) 100%);
  color: #e6a23c;
}

.stat-progress .stat-icon {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.05) 100%);
  color: #667eea;
}

.stat-completed .stat-icon {
  background: linear-gradient(135deg, rgba(103, 194, 58, 0.1) 0%, rgba(103, 194, 58, 0.05) 100%);
  color: #67c23a;
}

.stat-info {
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

/* 内容区域 */
.content-area {
  background: #ffffff;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .works-container {
    padding: 0;
  }

  .tabs-card :deep(.el-tabs__item) {
    padding: 0 12px;
    font-size: 14px;
  }

  .tab-label .el-icon {
    font-size: 14px;
  }
  
  .stat-value {
    font-size: 24px;
  }
  
  .stat-icon {
    width: 48px;
    height: 48px;
  }
}
</style>

