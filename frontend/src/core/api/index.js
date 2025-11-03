// API 统一导出
import authApi from '@/features/auth/api'
import certificatesApi from '@/features/certificates/api'

export default {
  auth: authApi,
  certificates: certificatesApi
}

