/**
 * 证书监控路由配置
 */
export default {
  path: '/certificates',
  name: 'Certificates',
  component: () => import('@/features/certificates/Certificates.vue'),
  meta: {
    title: '证书监控',
    icon: 'Document',
    showInMenu: true,
    order: 2,
    // 特殊配置：显示标签栏
    hasTabBar: true,
    tabs: [
      { name: 'list', label: '证书列表' },
      { name: 'update', label: '更新证书' },
      { name: 'generate-csr', label: '生成 CSR' },
      { name: 'validate-csr', label: '验证 CSR' },
      { name: 'validate-cert', label: '验证证书' }
    ]
  }
}

