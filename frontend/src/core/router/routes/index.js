/**
 * 路由模块聚合
 */
import dashboard from './dashboard'
import certificates from './certificates'
import passwords from './passwords'
import works from './works'
import profile from './profile'

// 导出所有子路由（按 order 排序）
export const moduleRoutes = [
  dashboard,
  certificates,
  passwords,
  works,
  profile
].sort((a, b) => (a.meta?.order || 999) - (b.meta?.order || 999))

// 导出菜单配置（仅显示在菜单中的路由）
export const menuRoutes = moduleRoutes.filter(route => route.meta?.showInMenu)

