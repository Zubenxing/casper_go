/**
 * 首页路由配置
 */
export default {
  path: '/dashboard',
  name: 'Dashboard',
  component: () => import('@/features/dashboard/Dashboard.vue'),
  meta: {
    title: '首页',
    icon: 'HomeFilled',
    showInMenu: true,
    order: 1
  }
}

