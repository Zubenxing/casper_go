/**
 * 个人中心路由配置
 */
export default {
  path: '/profile',
  name: 'Profile',
  component: () => import('@/features/auth/Profile.vue'),
  meta: {
    title: '个人中心',
    icon: 'User',
    showInMenu: true,
    order: 4
  }
}

