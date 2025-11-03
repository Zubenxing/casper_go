/**
 * 菜单配置
 * 从路由配置中自动生成菜单
 */
import { menuRoutes } from '@/core/router/routes'

export { menuRoutes }

/**
 * 根据路由 meta 判断是否显示标签栏
 */
export function hasTabBar(route) {
  return route.meta?.hasTabBar || false
}

/**
 * 获取标签栏配置
 */
export function getTabBarConfig(route) {
  return route.meta?.tabs || []
}

/**
 * 判断路由路径是否匹配
 */
export function isRouteMatch(currentPath, routePath) {
  return currentPath.startsWith(routePath)
}

