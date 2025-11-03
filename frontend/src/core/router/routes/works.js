export default {
  path: '/works',
  name: 'Works',
  component: () => import('@/features/works/Works.vue'),
  meta: {
    title: '工作记录',
    icon: 'Document',
    showInMenu: true,
    order: 4
  }
}

