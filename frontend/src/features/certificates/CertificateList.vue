<template>
  <div class="certificate-list">
    <el-card shadow="never" class="operation-card">
      <el-row :gutter="10">
        <el-col :xs="24" :sm="12" :md="10">
          <el-button type="primary" :icon="Plus" @click="showAddDialog = true">
            添加监控
          </el-button>
          <el-button :icon="Refresh" @click="handleRefresh" :loading="refreshing">
            刷新列表
          </el-button>
          <el-button :icon="Refresh" @click="handleCheckAll" :loading="checkingAll">
            检查所有
          </el-button>
          
          <el-dropdown @command="handleExport">
            <el-button :icon="Download">
              导出 <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="json">
                  <el-icon><Document /></el-icon> JSON 格式
                </el-dropdown-item>
                <el-dropdown-item command="csv">
                  <el-icon><Document /></el-icon> CSV 格式
                </el-dropdown-item>
                <el-dropdown-item command="excel">
                  <el-icon><Document /></el-icon> Excel 格式
                </el-dropdown-item>
                <el-dropdown-item command="txt">
                  <el-icon><Document /></el-icon> TXT 文本
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          
          <el-button :icon="Upload" @click="showImportDialog = true">
            导入
          </el-button>
        </el-col>
        <el-col :xs="24" :sm="12" :md="8" :offset="6" class="text-right">
          <el-input
            v-model="searchText"
            placeholder="搜索 URL 或域名"
            :prefix-icon="Search"
            clearable
            @input="handleSearch"
          />
        </el-col>
      </el-row>
    </el-card>

    <!-- 证书列表 -->
    <el-card shadow="never" class="table-card">
      <el-table
        :data="filteredCertificates"
        v-loading="loading"
        style="width: 100%"
        :default-sort="{ prop: 'created_at', order: 'descending' }"
      >
        <el-table-column type="index" label="#" width="60" />
        
        <el-table-column prop="url" label="URL" min-width="200" show-overflow-tooltip />
        
        <el-table-column prop="customer_name" label="客户名" width="150">
          <template #default="{ row }">
            <el-tag v-if="row.customer_name" type="info">{{ row.customer_name }}</el-tag>
            <span v-else style="color: #909399;">未设置</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="issuer" label="颁发者" width="200" show-overflow-tooltip />
        
        <el-table-column prop="days_left" label="剩余天数" width="120" sortable>
          <template #default="{ row }">
            <el-tag :type="getDaysLeftType(row.days_left)" effect="dark">
              {{ row.days_left }} 天
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ row.status_text }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="last_check_at" label="最后检查" width="180">
          <template #default="{ row }">
            {{ formatTime(row.last_check_at) }}
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" :icon="View" @click="handleView(row)" link>
              详情
            </el-button>
            <el-button size="small" :icon="Edit" @click="handleEdit(row)" link type="primary">
              编辑
            </el-button>
            <el-button size="small" :icon="Refresh" @click="handleUpdate(row.id)" link>
              更新
            </el-button>
            <el-button 
              size="small" 
              :icon="Delete" 
              type="danger" 
              @click="handleDelete(row.id)" 
              link
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <!-- 添加对话框 -->
    <el-dialog v-model="showAddDialog" title="添加证书监控" width="500px">
      <el-alert 
        title="支持自动补全" 
        type="success" 
        :closable="false"
        style="margin-bottom: 15px;"
      >
        <p>• 输入 www.baidu.com 会自动补充为 https://www.baidu.com</p>
        <p>• 输入 http:// 开头的会自动转为 https://</p>
      </el-alert>
      
      <el-form :model="addForm" :rules="addRules" ref="addFormRef">
        <el-form-item label="URL" prop="url">
          <el-input 
            v-model="addForm.url" 
            placeholder="www.baidu.com 或 https://www.baidu.com"
            :prefix-icon="Link"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAdd" :loading="adding">
          添加
        </el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="showDetailDialog" title="证书详情" width="600px">
      <el-descriptions :column="1" border v-if="currentCert">
        <el-descriptions-item label="URL">{{ currentCert.url }}</el-descriptions-item>
        <el-descriptions-item label="域名">{{ currentCert.domain }}</el-descriptions-item>
        <el-descriptions-item label="证书主体">{{ currentCert.subject }}</el-descriptions-item>
        <el-descriptions-item label="组织">{{ currentCert.organization }}</el-descriptions-item>
        <el-descriptions-item label="颁发者">{{ currentCert.issuer }}</el-descriptions-item>
        <el-descriptions-item label="生效时间">
          {{ formatTime(currentCert.not_before) }}
        </el-descriptions-item>
        <el-descriptions-item label="过期时间">
          {{ formatTime(currentCert.not_after) }}
        </el-descriptions-item>
        <el-descriptions-item label="剩余天数">
          <el-tag :type="getDaysLeftType(currentCert.days_left)" effect="dark">
            {{ currentCert.days_left }} 天
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentCert.status)">
            {{ currentCert.status_text }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="是否有效">
          <el-tag :type="currentCert.is_valid ? 'success' : 'danger'">
            {{ currentCert.is_valid ? '有效' : '无效' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="最后检查">
          {{ formatTime(currentCert.last_check_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">
          {{ formatTime(currentCert.created_at) }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <!-- 批量导入对话框 -->
    <el-dialog v-model="showImportDialog" title="批量导入证书" width="700px">
      <el-alert 
        title="导入说明" 
        type="info" 
        :closable="false"
        style="margin-bottom: 15px;"
      >
        <p>• <strong>格式：</strong>URL 客户名（空格分隔，客户名可选）</p>
        <p>• <strong>示例：</strong>www.baidu.com 百度公司</p>
        <p>• <strong>仅URL：</strong>www.baidu.com（不带客户名也可以）</p>
        <p>• <strong>CSV文件：</strong>第一列URL，第二列客户名</p>
      </el-alert>

      <el-tabs v-model="importTab" class="import-tabs">
        <!-- 文件上传 -->
        <el-tab-pane label="文件上传" name="file">
          <el-upload
            ref="uploadRef"
            class="upload-area"
            drag
            :auto-upload="false"
            :on-change="handleFileChange"
            :limit="1"
            accept=".json,.csv,.txt"
            :file-list="fileList"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              拖拽文件到此处 或 <em>点击选择文件</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                支持 JSON、CSV、TXT 格式，文件大小不超过 2MB
              </div>
            </template>
          </el-upload>
          
          <div v-if="parsedUrls.length > 0" class="preview-area">
            <el-divider content-position="left">
              解析结果：共 {{ parsedUrls.length }} 个 URL
            </el-divider>
            <div class="url-preview-list">
              <div 
                v-for="(item, index) in parsedUrls.slice(0, 10)" 
                :key="index"
                class="url-preview-item"
              >
                <el-tag style="margin-right: 10px;">{{ item.url }}</el-tag>
                <el-tag v-if="item.customer_name" type="success">{{ item.customer_name }}</el-tag>
                <el-button 
                  size="small" 
                  :icon="Delete" 
                  circle
                  @click="removeParsedUrl(index)"
                  style="margin-left: 10px;"
                />
              </div>
            </div>
            <div v-if="parsedUrls.length > 10" style="margin-top: 10px; color: #909399;">
              还有 {{ parsedUrls.length - 10 }} 个...
            </div>
          </div>
        </el-tab-pane>

        <!-- 手动输入 -->
        <el-tab-pane label="手动输入" name="manual">
          <el-form :model="importForm" ref="importFormRef">
            <el-form-item label="URL 列表">
              <el-input 
                v-model="importForm.urls" 
                type="textarea"
                :rows="12"
                placeholder="每行一个，格式：URL 客户名（空格分隔）&#10;&#10;www.baidu.com 百度公司&#10;www.taobao.com 淘宝公司&#10;www.jd.com 京东&#10;github.com"
              />
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
      
      <template #footer>
        <el-button @click="closeImportDialog">取消</el-button>
        <el-button type="primary" @click="handleImport" :loading="importing">
          <el-icon><Upload /></el-icon>
          开始导入 {{ importTab === 'file' && parsedUrls.length > 0 ? `(${parsedUrls.length} 个)` : '' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 编辑对话框 -->
    <el-dialog v-model="showEditDialog" title="编辑证书信息" width="500px">
      <el-form :model="editForm" ref="editFormRef" label-width="100px">
        <el-form-item label="客户名">
          <el-input 
            v-model="editForm.customer_name" 
            placeholder="输入客户名称"
            clearable
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input 
            v-model="editForm.remark" 
            type="textarea"
            :rows="3"
            placeholder="输入备注信息"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleEditSave">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onBeforeUnmount } from 'vue'
import { ElMessage, ElMessageBox, ElLoading, ElNotification } from 'element-plus'
import { 
  Plus, 
  Refresh, 
  Search, 
  View, 
  Delete, 
  Link,
  Document,
  CircleCheck,
  Warning,
  CircleClose,
  Download,
  Upload,
  ArrowDown,
  UploadFilled,
  Edit
} from '@element-plus/icons-vue'
import { formatTime, getDaysLeftType, getStatusType } from '@/core/utils/format'
import api from '@/core/api'
import * as XLSX from 'xlsx'

const loading = ref(false)
const adding = ref(false)
const checkingAll = ref(false)
const importing = ref(false)
const refreshing = ref(false)
const showAddDialog = ref(false)
const showDetailDialog = ref(false)
const showImportDialog = ref(false)
const showEditDialog = ref(false)
const addFormRef = ref(null)
const editFormRef = ref(null)
const importFormRef = ref(null)
const uploadRef = ref(null)
const searchText = ref('')
const importTab = ref('file')
const fileList = ref([])
const parsedUrls = ref([])

const certificates = ref([])
const currentCert = ref(null)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const addForm = reactive({
  url: ''
})

const importForm = reactive({
  urls: ''
})

const editForm = reactive({
  id: null,
  customer_name: '',
  remark: ''
})

const addRules = {
  url: [
    { required: true, message: '请输入 URL', trigger: 'blur' }
  ]
}

// 自动补充 https:// 前缀
const normalizeUrl = (url) => {
  url = url.trim()
  
  // 如果没有协议，自动添加 https://
  if (!url.startsWith('http://') && !url.startsWith('https://')) {
    url = 'https://' + url
  }
  
  // 如果是 http://，转换为 https://
  if (url.startsWith('http://')) {
    url = url.replace('http://', 'https://')
  }
  
  return url
}

// 搜索过滤
const filteredCertificates = computed(() => {
  if (!searchText.value) return certificates.value
  const search = searchText.value.toLowerCase()
  return certificates.value.filter(cert =>
    cert.url.toLowerCase().includes(search) ||
    cert.domain.toLowerCase().includes(search)
  )
})

// 统计数据
const stats = computed(() => {
  return {
    total: certificates.value.length,
    normal: certificates.value.filter(c => c.status === 1).length,
    warning: certificates.value.filter(c => c.status === 2).length,
    expired: certificates.value.filter(c => c.status === 3).length
  }
})

const fetchData = async () => {
  loading.value = true
  try {
    const response = await api.certificates.getList(pagination.page, pagination.pageSize)
    if (response.code === 0) {
      certificates.value = response.data.list || []
      pagination.total = response.data.total
    }
  } catch (error) {
    console.error('获取证书列表失败', error)
  } finally {
    loading.value = false
  }
}

const handleAdd = async () => {
  if (!addFormRef.value) return

  await addFormRef.value.validate(async (valid) => {
    if (!valid) return

    adding.value = true
    try {
      // 自动规范化 URL
      const normalizedUrl = normalizeUrl(addForm.url)
      
      const response = await api.certificates.add(normalizedUrl)
      if (response.code === 0) {
        ElMessage.success('添加成功')
        showAddDialog.value = false
        addForm.url = ''
        fetchData()
      }
    } catch (error) {
      console.error('添加失败', error)
    } finally {
      adding.value = false
    }
  })
}

const handleView = (row) => {
  currentCert.value = row
  showDetailDialog.value = true
}

const handleEdit = (row) => {
  editForm.id = row.id
  editForm.customer_name = row.customer_name || ''
  editForm.remark = row.remark || ''
  showEditDialog.value = true
}

const handleEditSave = async () => {
  if (!editFormRef.value) return

  try {
    // 调用更新接口（需要后端支持）
    const response = await api.certificates.updateInfo(editForm.id, {
      customer_name: editForm.customer_name,
      remark: editForm.remark
    })
    
    if (response.code === 0) {
      ElMessage.success('保存成功')
      showEditDialog.value = false
      fetchData()
    }
  } catch (error) {
    console.error('保存失败', error)
  }
}

const handleUpdate = async (id) => {
  try {
    const response = await api.certificates.update(id)
    if (response.code === 0) {
      ElMessage.success('更新成功')
      fetchData()
    }
  } catch (error) {
    console.error('更新失败', error)
  }
}

const handleDelete = async (id) => {
  ElMessageBox.confirm('确定要删除这个证书监控吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const response = await api.certificates.delete(id)
      if (response.code === 0) {
        ElMessage.success('删除成功')
        fetchData()
      }
    } catch (error) {
      console.error('删除失败', error)
    }
  }).catch(() => {})
}

const handleCheckAll = async () => {
  // 防止重复点击
  if (checkingAll.value) {
    ElMessage.warning('正在检查中，请稍候...')
    return
  }

  checkingAll.value = true
  
  // 显示加载提示
  const loadingInstance = ElLoading.service({
    lock: true,
    text: '正在并发检查所有证书，请稍候...',
    background: 'rgba(0, 0, 0, 0.7)',
  })

  try {
    const response = await api.certificates.checkAll()
    if (response.code === 0) {
      const result = response.data
      
      // 显示详细的检查结果
      if (result.failed > 0) {
        ElNotification({
          title: '检查完成（部分失败）',
          message: `总数：${result.total}，成功：${result.success}，失败：${result.failed}，耗时：${result.duration}`,
          type: 'warning',
          duration: 5000
        })
      } else {
        ElNotification({
          title: '检查完成',
          message: `成功检查了 ${result.success} 个证书，耗时：${result.duration}`,
          type: 'success',
          duration: 3000
        })
      }
      
      // 刷新列表
      await fetchData()
    }
  } catch (error) {
    console.error('检查失败', error)
    ElMessage.error(error.response?.data?.message || '检查失败，请稍后重试')
  } finally {
    checkingAll.value = false
    loadingInstance.close()
  }
}

const handleSearch = () => {
  // 搜索由 computed 自动处理
}

// 导出证书列表
const handleExport = (format) => {
  if (certificates.value.length === 0) {
    ElMessage.warning('没有可导出的数据')
    return
  }

  const timestamp = new Date().toISOString().slice(0, 10)
  
  switch (format) {
    case 'json':
      exportJSON(timestamp)
      break
    case 'csv':
      exportCSV(timestamp)
      break
    case 'excel':
      exportExcel(timestamp)
      break
    case 'txt':
      exportTXT(timestamp)
      break
  }
}

// 导出为 JSON
const exportJSON = (timestamp) => {
  const exportData = certificates.value.map(cert => ({
    URL: cert.url,
    客户名: cert.customer_name || '',
    域名: cert.domain,
    颁发者: cert.issuer,
    剩余天数: cert.days_left,
    状态: cert.status_text,
    是否有效: cert.is_valid ? '是' : '否',
    备注: cert.remark || '',
    最后检查时间: formatTime(cert.last_check_at)
  }))

  const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' })
  downloadFile(blob, `certificates_${timestamp}.json`)
  ElMessage.success('JSON 导出成功')
}

// 导出为 CSV
const exportCSV = (timestamp) => {
  const headers = ['URL', '客户名', '域名', '颁发者', '剩余天数', '状态', '是否有效', '备注', '最后检查时间']
  const rows = certificates.value.map(cert => [
    cert.url,
    cert.customer_name || '',
    cert.domain,
    cert.issuer,
    cert.days_left,
    cert.status_text,
    cert.is_valid ? '是' : '否',
    cert.remark || '',
    formatTime(cert.last_check_at)
  ])

  let csvContent = headers.join(',') + '\n'
  rows.forEach(row => {
    csvContent += row.map(cell => `"${cell}"`).join(',') + '\n'
  })

  const blob = new Blob(['\ufeff' + csvContent], { type: 'text/csv;charset=utf-8;' })
  downloadFile(blob, `certificates_${timestamp}.csv`)
  ElMessage.success('CSV 导出成功')
}

// 导出为 Excel
const exportExcel = (timestamp) => {
  const data = [
    ['URL', '客户名', '域名', '颁发者', '证书主体', '生效时间', '过期时间', '剩余天数', '状态', '是否有效', '备注', '最后检查时间']
  ]

  certificates.value.forEach(cert => {
    data.push([
      cert.url,
      cert.customer_name || '',
      cert.domain,
      cert.issuer,
      cert.subject,
      formatTime(cert.not_before),
      formatTime(cert.not_after),
      cert.days_left,
      cert.status_text,
      cert.is_valid ? '是' : '否',
      cert.remark || '',
      formatTime(cert.last_check_at)
    ])
  })

  const ws = XLSX.utils.aoa_to_sheet(data)
  
  // 设置列宽
  ws['!cols'] = [
    { wch: 35 }, // URL
    { wch: 20 }, // 客户名
    { wch: 20 }, // 域名
    { wch: 30 }, // 颁发者
    { wch: 20 }, // 证书主体
    { wch: 20 }, // 生效时间
    { wch: 20 }, // 过期时间
    { wch: 10 }, // 剩余天数
    { wch: 10 }, // 状态
    { wch: 10 }, // 是否有效
    { wch: 30 }, // 备注
    { wch: 20 }  // 最后检查时间
  ]

  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, '证书监控列表')
  XLSX.writeFile(wb, `certificates_${timestamp}.xlsx`)
  
  ElMessage.success('Excel 导出成功')
}

// 导出为 TXT
const exportTXT = (timestamp) => {
  let txtContent = '证书监控列表\n'
  txtContent += '='.repeat(80) + '\n\n'

  certificates.value.forEach((cert, index) => {
    txtContent += `${index + 1}. ${cert.url}\n`
    if (cert.customer_name) {
      txtContent += `   客户名: ${cert.customer_name}\n`
    }
    txtContent += `   域名: ${cert.domain}\n`
    txtContent += `   颁发者: ${cert.issuer}\n`
    txtContent += `   剩余天数: ${cert.days_left} 天\n`
    txtContent += `   状态: ${cert.status_text}\n`
    txtContent += `   是否有效: ${cert.is_valid ? '是' : '否'}\n`
    if (cert.remark) {
      txtContent += `   备注: ${cert.remark}\n`
    }
    txtContent += `   最后检查: ${formatTime(cert.last_check_at)}\n`
    txtContent += '\n' + '-'.repeat(80) + '\n\n'
  })

  txtContent += `\n导出时间: ${formatTime(new Date())}\n`
  txtContent += `总计: ${certificates.value.length} 条记录\n`

  const blob = new Blob([txtContent], { type: 'text/plain;charset=utf-8' })
  downloadFile(blob, `certificates_${timestamp}.txt`)
  ElMessage.success('TXT 导出成功')
}

// 通用下载函数
const downloadFile = (blob, filename) => {
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  link.click()
  URL.revokeObjectURL(url)
}

// 处理文件变化
const handleFileChange = (file) => {
  const reader = new FileReader()
  
  reader.onload = (e) => {
    const content = e.target.result
    const fileName = file.name.toLowerCase()
    
    try {
      if (fileName.endsWith('.json')) {
        parseJSON(content)
      } else if (fileName.endsWith('.csv')) {
        parseCSV(content)
      } else if (fileName.endsWith('.txt')) {
        parseTXT(content)
      } else {
        ElMessage.error('不支持的文件格式')
        fileList.value = []
      }
    } catch (error) {
      ElMessage.error('文件解析失败：' + error.message)
      fileList.value = []
    }
  }
  
  reader.readAsText(file.raw)
  fileList.value = [file]
}

// 解析 JSON 文件
const parseJSON = (content) => {
  const data = JSON.parse(content)
  const items = []
  
  if (Array.isArray(data)) {
    data.forEach(item => {
      if (typeof item === 'string') {
        items.push({ url: item, customer_name: '' })
      } else if (item.url || item.URL) {
        items.push({
          url: item.url || item.URL,
          customer_name: item.customer_name || item.CustomerName || item.客户名 || ''
        })
      }
    })
  }
  
  parsedUrls.value = items
  ElMessage.success(`解析成功：找到 ${items.length} 个 URL`)
}

// 解析 CSV 文件
const parseCSV = (content) => {
  const lines = content.split('\n').map(l => l.trim()).filter(l => l)
  const items = []
  
  lines.forEach((line, index) => {
    // 跳过表头
    if (index === 0 && (line.includes('URL') || line.includes('url') || line.includes('域名') || line.includes('客户'))) {
      return
    }
    
    // 解析 CSV 行（第一列URL，第二列客户名）
    const columns = line.split(',').map(c => c.replace(/^"|"$/g, '').trim())
    
    if (columns.length >= 1 && columns[0] && columns[0].includes('.')) {
      items.push({
        url: columns[0],
        customer_name: columns[1] || ''
      })
    }
  })
  
  parsedUrls.value = items
  ElMessage.success(`解析成功：找到 ${items.length} 个 URL`)
}

// 解析 TXT 文件
const parseTXT = (content) => {
  const lines = content.split('\n').map(l => l.trim()).filter(l => l)
  const items = []
  
  lines.forEach(line => {
    // 跳过注释和空行
    if (!line || line.startsWith('#')) return
    
    // 格式：URL 客户名（空格或逗号分隔）
    let url = ''
    let customerName = ''
    
    // 优先按空格分隔
    if (line.includes(' ')) {
      const parts = line.split(/\s+/) // 按空格分隔
      url = parts[0]
      customerName = parts.slice(1).join(' ') // 后面的都是客户名
    } else if (line.includes(',')) {
      const parts = line.split(',').map(p => p.trim())
      url = parts[0]
      customerName = parts[1] || ''
    } else {
      url = line
    }
    
    if (url && url.includes('.')) {
      items.push({ url, customer_name: customerName })
    }
  })
  
  parsedUrls.value = items
  ElMessage.success(`解析成功：找到 ${items.length} 个 URL`)
}

// 移除解析的 URL
const removeParsedUrl = (index) => {
  parsedUrls.value.splice(index, 1)
}

// 关闭导入对话框
const closeImportDialog = () => {
  showImportDialog.value = false
  importForm.urls = ''
  parsedUrls.value = []
  fileList.value = []
  importTab.value = 'file'
}

// 批量导入证书
const handleImport = async () => {
  let items = []
  
  if (importTab.value === 'file') {
    // 文件上传模式
    if (parsedUrls.value.length === 0) {
      ElMessage.warning('请先上传文件')
      return
    }
    items = parsedUrls.value
  } else {
    // 手动输入模式
    const lines = importForm.urls.split('\n').map(u => u.trim()).filter(u => u)
    if (lines.length === 0) {
      ElMessage.warning('请输入要导入的 URL')
      return
    }
    
    // 解析每行：格式 URL 客户名（空格或逗号分隔）
    items = lines.map(line => {
      let url = ''
      let customerName = ''
      
      // 优先按空格分隔
      if (line.includes(' ')) {
        const parts = line.split(/\s+/) // 按空格分隔
        url = parts[0]
        customerName = parts.slice(1).join(' ') // 后面的都是客户名
      } else if (line.includes(',')) {
        const parts = line.split(',').map(p => p.trim())
        url = parts[0]
        customerName = parts[1] || ''
      } else {
        url = line
      }
      
      return { url, customer_name: customerName }
    })
  }

  importing.value = true
  let successCount = 0
  let failCount = 0
  const failedItems = []
  const errorMessages = []

  console.log('开始导入，共', items.length, '个 URL')

  try {
    for (const item of items) {
      try {
        const normalizedUrl = normalizeUrl(item.url)
        const displayName = item.customer_name ? `${item.customer_name}(${item.url})` : item.url
        console.log(`导入中: ${displayName} -> ${normalizedUrl}`)
        
        // 先添加证书
        const response = await api.certificates.add(normalizedUrl)
        
        if (response.code === 0) {
          // 如果有客户名，再更新客户名
          if (item.customer_name) {
            await api.certificates.updateInfo(response.data.id, {
              customer_name: item.customer_name,
              remark: ''
            })
          }
          successCount++
          console.log(`✓ 成功: ${displayName}`)
        } else {
          failCount++
          failedItems.push(displayName)
          errorMessages.push(`${displayName}: ${response.message}`)
          console.error(`✗ 失败: ${displayName} - ${response.message}`)
        }
      } catch (error) {
        failCount++
        const displayName = item.customer_name ? `${item.customer_name}(${item.url})` : item.url
        failedItems.push(displayName)
        const errMsg = error.response?.data?.message || error.message || '未知错误'
        errorMessages.push(`${displayName}: ${errMsg}`)
        console.error(`✗ 失败: ${displayName}`, error)
      }
      
      // 添加小延迟，避免请求过快
      await new Promise(resolve => setTimeout(resolve, 200))
    }

    console.log(`导入完成 - 成功: ${successCount}, 失败: ${failCount}`)

    if (failCount === 0) {
      ElMessage.success(`导入完成！成功导入 ${successCount} 个证书`)
    } else {
      ElMessageBox.alert(
        `成功: ${successCount} 个\n失败: ${failCount} 个\n\n失败详情:\n${errorMessages.slice(0, 10).join('\n')}${errorMessages.length > 10 ? '\n...' : ''}`,
        '导入结果',
        {
          confirmButtonText: '确定',
          type: failCount > successCount ? 'warning' : 'success',
          dangerouslyUseHTMLString: false
        }
      )
    }
    
    closeImportDialog()
    fetchData()
  } finally {
    importing.value = false
  }
}

// 手动刷新
const handleRefresh = async () => {
  refreshing.value = true
  try {
    await fetchData()
    ElMessage.success('刷新成功')
  } catch (error) {
    console.error('刷新失败', error)
  } finally {
    refreshing.value = false
  }
}

// 自动刷新（每天00:05）
let autoRefreshTimer = null

const startAutoRefresh = () => {
  const scheduleNextRefresh = () => {
    const now = new Date()
    const tomorrow = new Date(now)
    tomorrow.setDate(tomorrow.getDate() + 1)
    tomorrow.setHours(0, 5, 0, 0) // 设置为明天00:05
    
    const timeUntilRefresh = tomorrow - now
    
    console.log(`下次自动刷新时间: ${tomorrow.toLocaleString()}`)
    
    autoRefreshTimer = setTimeout(() => {
      console.log('执行每日自动刷新...')
      fetchData()
      scheduleNextRefresh() // 安排下一次刷新
    }, timeUntilRefresh)
  }
  
  scheduleNextRefresh()
}

const stopAutoRefresh = () => {
  if (autoRefreshTimer) {
    clearTimeout(autoRefreshTimer)
    autoRefreshTimer = null
  }
}

onMounted(() => {
  fetchData()
  startAutoRefresh() // 启动自动刷新
})

onBeforeUnmount(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
.certificate-list {
  width: 100%;
}

.operation-card {
  margin-bottom: 20px;
}

.table-card {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.text-right {
  text-align: right;
}

.import-tabs {
  margin-top: 10px;
}

.upload-area {
  width: 100%;
}

:deep(.el-upload-dragger) {
  padding: 40px;
  width: 100%;
}

.preview-area {
  margin-top: 20px;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 4px;
  max-height: 300px;
  overflow-y: auto;
}

.url-preview-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.url-preview-item {
  display: flex;
  align-items: center;
  padding: 8px;
  background: white;
  border-radius: 4px;
}

@media (max-width: 768px) {
  .text-right {
    text-align: left;
    margin-top: 10px;
  }
}
</style>

