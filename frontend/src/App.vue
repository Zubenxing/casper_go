<template>
  <router-view />
</template>

<script setup>
// App 入口组件
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 
    'Microsoft YaHei', 'Helvetica Neue', Helvetica, Arial, sans-serif;
  background-color: #f0f2f5;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  width: 100%;
  min-height: 100vh;
}

/* 全局滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-thumb {
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background-color: rgba(0, 0, 0, 0.3);
}

::-webkit-scrollbar-track {
  background-color: transparent;
}

/* Element Plus 全局样式覆盖 */
.el-card {
  border-radius: 8px;
  border: none;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.el-table {
  font-size: 14px;
}

.el-table th {
  background-color: #fafafa;
  color: #606266;
  font-weight: 600;
}

.el-table__row:hover {
  background-color: #f5f7fa;
}

.el-button {
  font-weight: 500;
}

.el-tag {
  font-weight: 500;
  border: none;
}

/* 状态标签通用样式 */
.status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.status-badge.success {
  background-color: #f0f9ff;
  color: #0369a1;
}

.status-badge.warning {
  background-color: #fffbeb;
  color: #d97706;
}

.status-badge.danger {
  background-color: #fef2f2;
  color: #dc2626;
}

.status-badge.info {
  background-color: #f0f9ff;
  color: #2563eb;
}

.status-badge.primary {
  background-color: #eff6ff;
  color: #3b82f6;
}
</style>

