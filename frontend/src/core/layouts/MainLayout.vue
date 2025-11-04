<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside :width="isCollapse ? '64px' : '240px'" class="sidebar">
      <div class="logo">
        <div class="logo-icon">
          <el-icon :size="28">
            <Monitor />
          </el-icon>
        </div>
        <span v-if="!isCollapse" class="logo-text">Casper Platform</span>
      </div>

      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :router="true"
        class="sidebar-menu"
      >
        <el-menu-item
          v-for="menuItem in menuItems"
          :key="menuItem.path"
          :index="menuItem.path"
          class="menu-item"
        >
          <el-icon :size="18">
            <component :is="menuItem.meta.icon" />
          </el-icon>
          <template #title>{{ menuItem.meta.title }}</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- 主内容区 -->
    <el-container>
      <!-- 顶部栏 -->
      <el-header class="header" :class="{ 'has-tabs': hasTabBarPage }">
        <div class="header-left">
          <el-button 
            class="collapse-btn" 
            :icon="isCollapse ? Expand : Fold" 
            @click="toggleCollapse"
            text
          />

          <!-- 非标签栏页面显示页面标题 -->
          <div v-if="!hasTabBarPage" class="page-title">
            <h2>{{ currentTitle }}</h2>
          </div>

          <!-- 动态标签栏（根据路由配置） -->
          <div v-if="currentTabBarConfig.length > 0" class="header-tabs">
            <div 
              v-for="tab in currentTabBarConfig"
              :key="tab.name"
              class="tab-item"
              :class="{ active: certificateTab === tab.name }"
              @click="certificateTab = tab.name"
            >
              {{ tab.label }}
            </div>
          </div>
        </div>

        <div class="header-right">
          <!-- 搜索框 -->
          <el-input
            v-model="searchKeyword"
            class="search-input"
            placeholder="搜索..."
            :prefix-icon="Search"
            clearable
          />
          
          <!-- 证书监控页面显示刷新信息 -->
          <el-tooltip v-if="isCertificatesPage" content="证书自动检查时间" placement="bottom">
            <el-tag type="info" size="small" effect="plain">
              <el-icon><Clock /></el-icon>
              每天 00:05
            </el-tag>
          </el-tooltip>

          <!-- 通知按钮 -->
          <el-badge :value="0" :hidden="true" class="notification-badge">
            <el-button :icon="Bell" circle text />
          </el-badge>
          
          <!-- 用户信息 -->
          <el-dropdown @command="handleCommand" class="user-dropdown">
            <div class="user-info">
              <el-avatar :size="36" :icon="UserFilled" class="user-avatar" />
              <span class="username">{{ authStore.user?.nickname || authStore.user?.username }}</span>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人中心
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 内容区 -->
      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { ref, computed, onMounted, provide, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import * as ElementPlusIcons from '@element-plus/icons-vue'
import { useAuthStore } from '@/features/auth/store'
import { menuRoutes } from '@/core/config/menu'

const { 
  HomeFilled, 
  Document, 
  User, 
  Monitor,
  Fold,
  Expand,
  UserFilled,
  SwitchButton,
  Clock,
  Search,
  Bell
} = ElementPlusIcons

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const isCollapse = ref(false)
const searchKeyword = ref('')
const activeMenu = computed(() => route.path)
const currentTitle = computed(() => route.meta.title || '首页')

// 菜单项（从配置读取）
const menuItems = computed(() => menuRoutes)

// 判断当前页面是否有标签栏
const hasTabBarPage = computed(() => {
  const matched = router.currentRoute.value.matched
  return matched.some(route => route.meta?.hasTabBar)
})

// 判断是否是证书监控页面（用于显示自动刷新标签）
const isCertificatesPage = computed(() => route.path.startsWith('/certificates'))

// 获取当前页面的标签栏配置
const currentTabBarConfig = computed(() => {
  const matched = router.currentRoute.value.matched
  const routeWithTabs = matched.find(route => route.meta?.hasTabBar)
  return routeWithTabs?.meta?.tabs || []
})

// 动态 tab 管理
const certificateTab = ref('list')

// 根据当前路由初始化默认 tab
const initializeTab = () => {
  if (currentTabBarConfig.value.length > 0) {
    certificateTab.value = currentTabBarConfig.value[0].name
  }
}

// 监听路由变化，重置 tab
watch(() => route.path, () => {
  initializeTab()
}, { immediate: true })

// 提供给子组件使用
provide('certificateTab', certificateTab)

const toggleCollapse = () => {
  isCollapse.value = !isCollapse.value
}

const handleCommand = async (command) => {
  if (command === 'logout') {
    ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      await authStore.logout()
      ElMessage.success('已退出登录')
      router.push('/login')
    }).catch(() => {})
  } else if (command === 'profile') {
    router.push('/profile')
  }
}

onMounted(() => {
  authStore.init()
  if (authStore.isLoggedIn && !authStore.user) {
    authStore.fetchUserInfo()
  }
  initializeTab()
})
</script>

<style scoped>
/* ========== 整体容器 ========== */
.layout-container {
  height: 100vh;
  background-color: #f0f2f5;
}

/* ========== 侧边栏样式 ========== */
.sidebar {
  background: linear-gradient(180deg, #1a1f36 0%, #14182b 100%);
  transition: width 0.28s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.15);
  position: relative;
  z-index: 1000;
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 0 16px;
  background: rgba(0, 0, 0, 0.2);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.logo-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  color: white;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.logo-text {
  color: #ffffff;
  font-size: 18px;
  font-weight: 600;
  letter-spacing: 0.5px;
  white-space: nowrap;
}

/* 侧边栏菜单 */
.sidebar-menu {
  border: none;
  background: transparent !important;
  padding: 8px;
}

.sidebar-menu :deep(.el-menu-item) {
  height: 48px;
  line-height: 48px;
  margin: 4px 0;
  border-radius: 8px;
  color: rgba(255, 255, 255, 0.75);
  background: transparent;
  transition: all 0.3s ease;
  font-size: 14px;
  font-weight: 500;
}

.sidebar-menu :deep(.el-menu-item:hover) {
  background: rgba(255, 255, 255, 0.1);
  color: #ffffff;
}

.sidebar-menu :deep(.el-menu-item.is-active) {
  background: linear-gradient(90deg, rgba(102, 126, 234, 0.2) 0%, rgba(118, 75, 162, 0.1) 100%);
  color: #ffffff;
  border-left: 3px solid #667eea;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.2);
}

.sidebar-menu :deep(.el-menu-item .el-icon) {
  color: inherit;
  margin-right: 12px;
}

/* ========== 顶部栏样式 ========== */
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #ffffff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  padding: 0 24px;
  height: 64px;
  border-bottom: 1px solid #f0f0f0;
  position: relative;
  z-index: 999;
}

.header.has-tabs {
  border-bottom: none;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
}

.collapse-btn {
  font-size: 20px;
  color: #5a5e66;
  width: 40px;
  height: 40px;
}

.collapse-btn:hover {
  background-color: #f5f7fa;
  color: #409eff;
}

.page-title h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

/* 标签栏 */
.header-tabs {
  display: flex;
  gap: 0;
  margin-left: 24px;
  height: 64px;
}

.tab-item {
  padding: 0 28px;
  height: 64px;
  line-height: 64px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  color: #606266;
  background: transparent;
  border-bottom: 2px solid transparent;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  user-select: none;
  position: relative;
}

.tab-item:hover {
  color: #409eff;
  background: #f5f7fa;
}

.tab-item.active {
  color: #409eff;
  border-bottom-color: #409eff;
  background: #f0f7ff;
  font-weight: 600;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

/* 搜索框 */
.search-input {
  width: 200px;
}

.search-input :deep(.el-input__wrapper) {
  background-color: #f5f7fa;
  box-shadow: none;
  border-radius: 20px;
  transition: all 0.3s;
}

.search-input :deep(.el-input__wrapper:hover),
.search-input :deep(.el-input__wrapper.is-focus) {
  background-color: #ffffff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

/* 通知徽章 */
.notification-badge {
  cursor: pointer;
}

/* 用户下拉菜单 */
.user-dropdown {
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 4px 12px 4px 4px;
  border-radius: 20px;
  transition: all 0.3s;
}

.user-info:hover {
  background-color: #f5f7fa;
}

.user-avatar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.username {
  color: #303133;
  font-size: 14px;
  font-weight: 500;
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* ========== 主内容区 ========== */
.main-content {
  background-color: #f0f2f5;
  padding: 24px;
  overflow-y: auto;
}

/* ========== 过渡动画 ========== */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
