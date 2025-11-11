export default {
  path: '/workflows',
  name: 'Workflows',
  component: () => import('@/features/workflows/Workflows.vue'),
  meta: {
    title: 'AI 工作流',
    requiresAuth: true,
    icon: 'Promotion',
    showInMenu: true,
    order: 4
  }
}

