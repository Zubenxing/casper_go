export default {
  path: '/works',
  name: 'Works',
  component: () => import('@/features/works/Works.vue'),
  meta: {
    title: '工作记录',
    icon: 'Document',
    showInMenu: true,
    order: 4,
    // 顶部标签栏配置
    hasTabBar: true,
    tabs: [
      { name: 'works', label: '工作内容' },
      { name: 'issues', label: '问题记录' }
    ]
  }
}

