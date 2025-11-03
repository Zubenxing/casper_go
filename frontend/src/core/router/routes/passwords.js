import Passwords from '@/features/passwords/Passwords.vue'
import { Lock } from '@element-plus/icons-vue'

export default {
  path: '/passwords',
  name: 'Passwords',
  component: Passwords,
  meta: {
    title: '密码管理',
    icon: 'Lock',
    showInMenu: true,
    order: 3
  }
}

