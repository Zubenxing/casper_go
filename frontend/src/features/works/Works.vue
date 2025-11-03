<template>
  <div class="works-container">
    <!-- 页面标题和统计 -->
    <div class="page-header">
      <div class="header-left">
        <h1 class="page-title">
          <el-icon class="title-icon"><Document /></el-icon>
          工作记录管理
        </h1>
        <p class="page-subtitle">记录您的工作任务和遇到的问题</p>
      </div>
    </div>

    <!-- Tab 切换 -->
    <el-card class="tabs-card" shadow="hover">
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="工作记录" name="works">
          <template #label>
            <span class="tab-label">
              <el-icon><Notebook /></el-icon>
              <span>工作记录</span>
              <el-badge v-if="workStats.total > 0" :value="workStats.total" class="tab-badge" />
            </span>
          </template>
          <WorkList ref="workListRef" />
        </el-tab-pane>

        <el-tab-pane label="问题记录" name="issues">
          <template #label>
            <span class="tab-label">
              <el-icon><Warning /></el-icon>
              <span>问题记录</span>
              <el-badge v-if="issueStats.total > 0" :value="issueStats.total" class="tab-badge" type="warning" />
            </span>
          </template>
          <IssueList ref="issueListRef" />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Document, Notebook, Warning } from '@element-plus/icons-vue'
import WorkList from './WorkList.vue'
import IssueList from './IssueList.vue'
import { getWorkStats, getWorkIssueStats } from './api'

const activeTab = ref('works')
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

const handleTabChange = (tabName) => {
  console.log('切换到标签页:', tabName)
}

const loadStats = async () => {
  try {
    const [workRes, issueRes] = await Promise.all([
      getWorkStats(),
      getWorkIssueStats()
    ])
    workStats.value = workRes.data
    issueStats.value = issueRes.data
  } catch (error) {
    console.error('加载统计信息失败:', error)
  }
}

onMounted(() => {
  loadStats()
})

// 暴露刷新方法给子组件调用
defineExpose({
  refreshStats: loadStats
})
</script>

<style scoped>
.works-container {
  padding: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: calc(100vh - 60px);
}

/* 页面标题 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 24px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.25);
  transition: all 0.3s ease;
}

.page-header:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(102, 126, 234, 0.35);
}

.header-left {
  flex: 1;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #2d3748;
  margin: 0 0 8px 0;
  display: flex;
  align-items: center;
  gap: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.title-icon {
  font-size: 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-subtitle {
  color: #718096;
  font-size: 14px;
  margin: 0;
  padding-left: 44px;
}

/* Tabs 卡片 */
.tabs-card {
  border-radius: 16px;
  border: none;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.15);
  transition: all 0.3s ease;
}

.tabs-card:hover {
  box-shadow: 0 12px 40px rgba(102, 126, 234, 0.25);
}

.tabs-card :deep(.el-card__body) {
  padding: 0;
}

/* Tab 样式 */
.tabs-card :deep(.el-tabs__header) {
  margin: 0;
  padding: 20px 20px 0;
  background: linear-gradient(135deg, #f6f8fb 0%, #ffffff 100%);
}

.tabs-card :deep(.el-tabs__nav-wrap::after) {
  display: none;
}

.tabs-card :deep(.el-tabs__active-bar) {
  height: 3px;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  border-radius: 3px 3px 0 0;
}

.tabs-card :deep(.el-tabs__item) {
  font-size: 16px;
  font-weight: 500;
  color: #718096;
  padding: 0 24px;
  height: 48px;
  line-height: 48px;
  transition: all 0.3s ease;
}

.tabs-card :deep(.el-tabs__item:hover) {
  color: #667eea;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 8px 8px 0 0;
}

.tabs-card :deep(.el-tabs__item.is-active) {
  color: #667eea;
  font-weight: 600;
}

.tab-label {
  display: flex;
  align-items: center;
  gap: 8px;
  position: relative;
}

.tab-label .el-icon {
  font-size: 18px;
}

.tab-badge {
  margin-left: 4px;
}

.tab-badge :deep(.el-badge__content) {
  font-weight: 600;
  font-size: 11px;
}

/* Tab 内容区域 */
.tabs-card :deep(.el-tabs__content) {
  padding: 24px;
  background: #ffffff;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .works-container {
    padding: 16px;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .page-title {
    font-size: 24px;
  }

  .page-subtitle {
    padding-left: 38px;
  }

  .tabs-card :deep(.el-tabs__item) {
    padding: 0 16px;
    font-size: 14px;
  }

  .tab-label .el-icon {
    font-size: 16px;
  }
}
</style>

