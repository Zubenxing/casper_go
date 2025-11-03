<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside :width="isCollapse ? '64px' : '200px'" class="sidebar">
      <div class="logo">
        <el-icon v-if="!isCollapse" :size="30" color="#409eff">
          <Monitor />
        </el-icon>
        <span v-if="!isCollapse" class="logo-text">Casper</span>
      </div>

      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :router="true"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <el-menu-item
          v-for="menuItem in menuItems"
          :key="menuItem.path"
          :index="menuItem.path"
        >
          <el-icon>
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
          <el-icon class="collapse-icon" @click="toggleCollapse">
            <Fold v-if="!isCollapse" />
            <Expand v-else />
          </el-icon>

          <!-- 非标签栏页面显示面包屑 -->
          <el-breadcrumb v-if="!hasTabBarPage" separator="/">
            <el-breadcrumb-item>{{ currentTitle }}</el-breadcrumb-item>
          </el-breadcrumb>

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
          <!-- 证书监控页面显示刷新信息 -->
          <el-tag v-if="isCertificatesPage" type="info" size="small" effect="plain">
            <el-icon><Clock /></el-icon>
            自动刷新：每天00:05
          </el-tag>
          
          <span class="username">{{ authStore.user?.nickname || authStore.user?.username }}</span>
          <el-dropdown @command="handleCommand">
            <el-avatar :size="35" :icon="UserFilled" />
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
  Clock
} = ElementPlusIcons

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const isCollapse = ref(false)
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
.layout-container {
  height: 100vh;
}

.sidebar {
  background-color: #304156;
  transition: width 0.3s;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  background-color: #1f2d3d;
  color: white;
  font-size: 20px;
  font-weight: bold;
}

.logo-text {
  color: #409eff;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
  padding: 0 20px;
  height: 60px;
}

.header.has-tabs {
  border-bottom: 1px solid #e4e7ed;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
  flex: 1;
}

/* 标签栏 */
.header-tabs {
  display: flex;
  gap: 0;
  margin-left: 10px;
}

.tab-item {
  padding: 0 24px;
  height: 60px;
  line-height: 60px;
  cursor: pointer;
  font-size: 14px;
  color: #606266;
  background: transparent;
  border-bottom: 3px solid transparent;
  transition: all 0.3s ease;
  user-select: none;
  font-weight: 500;
  position: relative;
  margin-bottom: -1px;
}

.tab-item:hover {
  color: #409eff;
  background: #ecf5ff;
}

.tab-item.active {
  color: #409eff;
  border-bottom-color: #409eff;
  background: #ecf5ff;
}

.collapse-icon {
  font-size: 20px;
  cursor: pointer;
  transition: color 0.3s;
}

.collapse-icon:hover {
  color: #409eff;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 15px;
}

.username {
  color: #606266;
  font-size: 14px;
}

.main-content {
  background-color: #f5f7fa;
  padding: 20px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
