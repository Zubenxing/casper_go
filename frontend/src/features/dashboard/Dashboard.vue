<template>
  <div class="dashboard">
    <!-- 顶部布局 -->
    <div class="dashboard-top">
      <!-- 左侧：实时数据 + 折线图 -->
      <div class="dashboard-left">
        <!-- 实时数据卡片 -->
        <div class="data-card">
          <h3 class="card-title">实时数据</h3>
          <div class="card-stats">
            <div class="stat-item">
              <div class="stat-label">URL监控总数</div>
              <div class="stat-value">{{ certStats.total }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">证书到期告警数</div>
              <div class="stat-value">{{ certStats.warning }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">URL监控失败数</div>
              <div class="stat-value danger">{{ certStats.expired }}</div>
            </div>
          </div>
        </div>

        <!-- 折线图 -->
        <div class="main-chart">
        <div class="chart-header">
          <h3 class="chart-title">近一周工作内容统计</h3>
          <div class="chart-legend">
            <span class="legend-item">
              <i class="legend-dot" style="background: #f59e0b;"></i>
              <span>累计未完成工作</span>
            </span>
            <span class="legend-item">
              <i class="legend-dot" style="background: #10b981;"></i>
              <span>累计未完成问题</span>
            </span>
          </div>
        </div>
          <div ref="trendChartRef" class="chart-body"></div>
        </div>
      </div>

      <!-- 右侧：两个圆环图 -->
      <div class="dashboard-right">
        <!-- 用户数据 - 证书状态分布 -->
        <div class="ring-card">
          <h3 class="card-title">证书数据</h3>
          <div class="card-divider"></div>
          <div class="ring-container">
            <div ref="certChartRef" class="ring-chart"></div>
            <div class="ring-center">
              <div class="ring-label">正常证书数</div>
              <div class="ring-value">{{ certStats.normal }}</div>
              <div class="ring-percent">占比 {{ getPercentage(certStats.normal, certStats.total) }}%</div>
            </div>
          </div>
        </div>

        <!-- 账号数据 - 密码分类分布 -->
        <div class="ring-card">
          <h3 class="card-title">账号数据</h3>
          <div class="card-divider"></div>
          <div class="ring-container">
            <div ref="passwordChartRef" class="ring-chart"></div>
            <div class="ring-center">
              <div class="ring-label">分类总数</div>
              <div class="ring-value">{{ passwordStats.categories.length }}</div>
              <div class="ring-percent">账号总数 {{ passwordStats.total }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部账号分类占比 -->
    <div class="bottom-section">
      <h3 class="section-title">账号分类占比</h3>
      <div class="type-legend">
        <span v-for="item in assetTypes" :key="item.name" class="type-item">
          <i class="type-dot" :style="{ background: item.color }"></i>
          <span>{{ item.name }}</span>
        </span>
      </div>
      <div class="type-bar">
        <div
          v-for="item in assetTypes"
          :key="item.name"
          class="bar-segment"
          :style="{ width: item.percent + '%', background: item.color }"
        ></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import * as echarts from 'echarts'
import request from '@/core/api/request'

// 图表实例
const trendChartRef = ref(null)
const certChartRef = ref(null)
const passwordChartRef = ref(null)
let trendChart = null
let certChart = null
let passwordChart = null

// 数据统计
const certStats = ref({
  total: 0,
  normal: 0,
  warning: 0,
  expired: 0
})

const workStats = ref({
  total: 0,
  pending: 0,
  in_progress: 0,
  completed: 0
})

const issueStats = ref({
  total: 0,
  open: 0,
  in_progress: 0,
  resolved: 0
})

const passwordStats = ref({
  total: 0,
  categories: []
})

// 资产类型分布（基于操作系统类型）
const assetTypes = computed(() => {
  const total = passwordStats.value.total
  const categories = passwordStats.value.categories || []
  
  return categories.map((cat, index) => {
    const colors = ['#16a34a', '#059669', '#0891b2']
    return {
      name: cat.name || '未分类',
      count: cat.count || 0,
      percent: total > 0 ? ((cat.count / total) * 100).toFixed(1) : 0,
      color: colors[index % colors.length]
    }
  })
})

// 计算百分比
const getPercentage = (value, total) => {
  if (total === 0) return 0
  return ((value / total) * 100).toFixed(1)
}

// 获取活跃密码数量（返回最大分类的数量）
const getActivePasswordCount = () => {
  const categories = passwordStats.value.categories || []
  if (categories.length === 0) return 0
  // 返回数量最多的分类的数量
  const maxCategory = categories.reduce((max, cat) => cat.count > max.count ? cat : max, categories[0])
  return maxCategory?.count || 0
}

// 获取证书统计
const fetchCertStats = async () => {
  try {
    // 获取所有证书数据（后端限制最大100，需要分页获取所有数据）
    let page = 1
    const pageSize = 100
    let allCerts = []
    let total = 0
    
    // 获取第一页
    const firstResponse = await request.get('/certificates', { params: { page: 1, page_size: pageSize } })
    if (firstResponse.code === 0) {
      const list = firstResponse.data.list || firstResponse.data.items || []
      allCerts = [...list]
      total = firstResponse.data.total || 0
      
      // 计算总页数
      const totalPages = Math.ceil(total / pageSize)
      
      // 获取剩余页
      for (let p = 2; p <= totalPages; p++) {
        const response = await request.get('/certificates', { params: { page: p, page_size: pageSize } })
        if (response.code === 0) {
          const list = response.data.list || response.data.items || []
          allCerts = [...allCerts, ...list]
        }
      }
      
      // 统计各状态数量
      certStats.value.total = total
      certStats.value.normal = allCerts.filter(c => c.status === 1).length
      certStats.value.warning = allCerts.filter(c => c.status === 2).length
      certStats.value.expired = allCerts.filter(c => c.status === 3).length
    }
  } catch (error) {
    console.error('获取证书统计失败', error)
  }
}

// 获取工作记录统计
const fetchWorkStats = async () => {
  try {
    const response = await request.get('/works/stats')
    if (response.code === 0) {
      workStats.value = response.data
    }
  } catch (error) {
    console.error('获取工作统计失败', error)
  }
}

// 获取问题记录统计
const fetchIssueStats = async () => {
  try {
    const response = await request.get('/work-issues/stats')
    if (response.code === 0) {
      issueStats.value = response.data
    }
  } catch (error) {
    console.error('获取问题统计失败', error)
  }
}

// 获取密码管理统计（按分类）
const fetchPasswordStats = async () => {
  try {
    const response = await request.get('/passwords/stats')
    if (response.code === 0) {
      passwordStats.value.total = response.data.total || 0
      
      // 获取所有密码并统计分类（分页获取）
      let allPasswords = []
      const pageSize = 100
      const total = response.data.total || 0
      const totalPages = Math.ceil(total / pageSize)
      
      for (let p = 1; p <= totalPages; p++) {
        const listResponse = await request.get('/passwords', { params: { page: p, page_size: pageSize } })
        if (listResponse.code === 0) {
          const list = listResponse.data.items || listResponse.data.list || []
          allPasswords = [...allPasswords, ...list]
        }
      }
      
      // 统计分类
      const categoryMap = {}
      allPasswords.forEach(item => {
        const category = item.category || '未分类'
        categoryMap[category] = (categoryMap[category] || 0) + 1
      })
      
      passwordStats.value.categories = Object.entries(categoryMap).map(([name, count]) => ({
        name,
        count
      }))
    }
  } catch (error) {
    console.error('获取密码统计失败', error)
  }
}

// 初始化趋势图表
const initTrendChart = async () => {
  if (!trendChartRef.value) return
  
  trendChart = echarts.init(trendChartRef.value)
  
  // 近7天的数据 - 获取每天实际的未完成任务数
  const dates = []
  const workData = []
  const issueData = []
  const now = new Date()
  
  try {
    // 获取所有工作记录和问题记录
    const [workResponse, issueResponse] = await Promise.all([
      request.get('/works', { params: { page: 1, page_size: 1000 } }),
      request.get('/work-issues', { params: { page: 1, page_size: 1000 } })
    ])
    
    const works = workResponse.data?.list || []
    const issues = issueResponse.data?.list || []
    
    // 统计近7天每天的累计未完成任务数
    for (let i = 6; i >= 0; i--) {
      const date = new Date(now)
      date.setDate(date.getDate() - i)
      date.setHours(23, 59, 59, 999)  // 设置为当天结束时间
      
      dates.push(`${date.getMonth() + 1}-${date.getDate()}`)
      
      // 统计截至当天，已创建且仍未完成的工作
      const dayWorks = works.filter(w => {
        const createdAt = new Date(w.created_at)
        // 在该天或之前创建，且状态为未完成
        return createdAt <= date && 
               (w.status === 'pending' || w.status === 'in_progress')
      }).length
      
      // 统计截至当天，已创建且仍未完成的问题
      const dayIssues = issues.filter(i => {
        const createdAt = new Date(i.created_at)
        // 在该天或之前创建，且状态为未完成
        return createdAt <= date && 
               (i.status === 'open' || i.status === 'in_progress')
      }).length
      
      workData.push(dayWorks)
      issueData.push(dayIssues)
    }
  } catch (error) {
    console.error('获取趋势数据失败', error)
    // 如果获取失败，使用默认数据
    for (let i = 6; i >= 0; i--) {
      const date = new Date(now)
      date.setDate(date.getDate() - i)
      dates.push(`${date.getMonth() + 1}-${date.getDate()}`)
      workData.push(0)
      issueData.push(0)
    }
  }
  
  const option = {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(255, 255, 255, 0.95)',
      borderColor: '#e4e7ed',
      borderWidth: 1,
      textStyle: {
        color: '#303133',
        fontSize: 13
      },
      padding: [8, 12],
      axisPointer: {
        type: 'cross',
        crossStyle: {
          color: '#909399'
        },
        lineStyle: {
          color: '#909399',
          type: 'dashed'
        }
      },
      formatter: function(params) {
        let result = params[0].axisValue + '<br/>'
        params.forEach(item => {
          result += `<span style="display:inline-block;margin-right:5px;border-radius:10px;width:10px;height:10px;background-color:${item.color}"></span>`
          result += `${item.seriesName}: ${item.value} 个<br/>`
        })
        return result
      }
    },
    legend: {
      show: false
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '10%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: dates,
      axisLine: {
        lineStyle: {
          color: '#e5e7eb'
        }
      },
      axisLabel: {
        color: '#6b7280'
      }
    },
    yAxis: {
      type: 'value',
      minInterval: 1,  // 强制Y轴只显示整数间隔
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: '#6b7280',
        formatter: '{value}'  // 只显示整数
      },
      splitLine: {
        lineStyle: {
          color: '#f3f4f6'
        }
      }
    },
    series: [
      {
        name: '累计未完成工作',
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        data: workData,
        lineStyle: {
          color: '#f59e0b',
          width: 3
        },
        itemStyle: {
          color: '#f59e0b'
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(245, 158, 11, 0.3)' },
            { offset: 1, color: 'rgba(245, 158, 11, 0.05)' }
          ])
        }
      },
      {
        name: '累计未完成问题',
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 8,
        data: issueData,
        lineStyle: {
          color: '#10b981',
          width: 3
        },
        itemStyle: {
          color: '#10b981'
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(16, 185, 129, 0.3)' },
            { offset: 1, color: 'rgba(16, 185, 129, 0.05)' }
          ])
        }
      }
    ]
  }
  
  trendChart.setOption(option)
}

// 初始化证书状态饼图
const initCertChart = () => {
  if (!certChartRef.value) return
  
  certChart = echarts.init(certChartRef.value)
  
  const chartData = [
    { value: certStats.value.normal, name: '正常' },
    { value: certStats.value.warning, name: '告警' },
    { value: certStats.value.expired, name: '过期' }
  ]
  
  const option = {
    color: ['#10b981', '#f59e0b', '#ef4444'],
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} 个 ({d}%)',
      backgroundColor: 'rgba(255, 255, 255, 0.95)',
      borderColor: '#e4e7ed',
      borderWidth: 1,
      textStyle: {
        color: '#303133',
        fontSize: 13
      },
      padding: [8, 12]
    },
    series: [
      {
        name: '证书状态',
        type: 'pie',
        radius: ['70%', '90%'],
        avoidLabelOverlap: false,
        label: {
          show: false
        },
        labelLine: {
          show: false
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.3)'
          }
        },
        data: chartData
      }
    ]
  }
  
  certChart.setOption(option)
}

// 初始化密码分类饼图
const initPasswordChart = () => {
  if (!passwordChartRef.value) return
  
  passwordChart = echarts.init(passwordChartRef.value)
  
  const categories = passwordStats.value.categories || []
  
  // 确保有数据
  const chartData = categories.length > 0 
    ? categories.map((cat, index) => {
        const colors = ['#16a34a', '#059669', '#0891b2', '#0284c7', '#2563eb']
        return {
          value: cat.count,
          name: cat.name,
          itemStyle: {
            color: colors[index % colors.length]
          }
        }
      })
    : [{ value: 1, name: '暂无数据', itemStyle: { color: '#e5e7eb' } }]
  
  const option = {
    color: ['#16a34a', '#059669', '#0891b2', '#0284c7', '#2563eb'],
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} 个账号 ({d}%)',
      backgroundColor: 'rgba(255, 255, 255, 0.95)',
      borderColor: '#e4e7ed',
      borderWidth: 1,
      textStyle: {
        color: '#303133',
        fontSize: 13
      },
      padding: [8, 12]
    },
    series: [
      {
        name: '密码分类',
        type: 'pie',
        radius: ['70%', '90%'],
        avoidLabelOverlap: false,
        label: {
          show: false
        },
        labelLine: {
          show: false
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.3)'
          }
        },
        data: chartData
      }
    ]
  }
  
  passwordChart.setOption(option)
}

// 响应式调整图表大小
const handleResize = () => {
  trendChart?.resize()
  certChart?.resize()
  passwordChart?.resize()
}

// 获取所有数据并初始化图表
const initDashboard = async () => {
  await Promise.all([
    fetchCertStats(),
    fetchWorkStats(),
    fetchIssueStats(),
    fetchPasswordStats()
  ])
  
  await nextTick()
  
  await initTrendChart()  // 等待趋势图初始化（需要获取数据）
  initCertChart()
  initPasswordChart()
  
  window.addEventListener('resize', handleResize)
}

onMounted(() => {
  initDashboard()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  trendChart?.dispose()
  certChart?.dispose()
  passwordChart?.dispose()
})
</script>

<style scoped>
.dashboard {
  width: 100%;
  max-width: 100%;
  margin: 0;
  padding: 0;
}

/* 顶部布局 */
.dashboard-top {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 16px;
  margin-bottom: 16px;
}

/* 左侧：实时数据 + 折线图 */
.dashboard-left {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* 右侧：两个圆环图 */
.dashboard-right {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* 数据卡片 */
.data-card {
  background: white;
  border-radius: 4px;
  padding: 24px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
  transition: box-shadow 0.3s ease;
}

.data-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
}

.card-title {
  font-size: 14px;
  font-weight: 500;
  color: #606266;
  margin: 0 0 16px 0;
  padding: 0;
}

/* 卡片分割线 */
.card-divider {
  height: 1px;
  background: #e4e7ed;
  margin-bottom: 20px;
}

/* 实时数据卡片样式 */
.card-stats {
  display: flex;
  justify-content: space-between;
  gap: 16px;
}

.stat-item {
  flex: 1;
  text-align: center;
}

.stat-label {
  font-size: 12px;
  color: #909399;
  margin-bottom: 12px;
  font-weight: 400;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
}

.stat-value.danger {
  color: #f56c6c;
}

/* 折线图 */
.main-chart {
  background: white;
  border-radius: 4px;
  padding: 24px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
  flex: 1;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.chart-title {
  font-size: 14px;
  font-weight: 500;
  color: #606266;
  margin: 0;
}

.chart-legend {
  display: flex;
  gap: 24px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #606266;
}

.legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}

.chart-body {
  width: 100%;
  height: 300px;
}

/* 圆环图卡片 */
.ring-card {
  background: white;
  border-radius: 4px;
  padding: 20px 24px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
  flex: 1;
  display: flex;
  flex-direction: column;
}

.ring-container {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  flex: 1;
}

.ring-chart {
  width: 100%;
  height: 180px;
}

.ring-center {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  pointer-events: none;
}

.ring-label {
  font-size: 12px;
  color: #909399;
  margin-bottom: 8px;
  font-weight: 400;
}

.ring-value {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
}

.ring-percent {
  font-size: 12px;
  color: #909399;
}

/* 底部账号分类占比 */
.bottom-section {
  background: white;
  border-radius: 4px;
  padding: 24px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
}

.section-title {
  font-size: 14px;
  font-weight: 500;
  color: #606266;
  margin: 0 0 20px 0;
}

.type-legend {
  display: flex;
  gap: 24px;
  margin-bottom: 16px;
  font-size: 12px;
  color: #606266;
  flex-wrap: wrap;
}

.type-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.type-dot {
  width: 8px;
  height: 8px;
  border-radius: 2px;
  display: inline-block;
}

.type-bar {
  display: flex;
  width: 100%;
  height: 24px;
  border-radius: 4px;
  overflow: hidden;
  background: #f5f7fa;
}

.bar-segment {
  height: 100%;
  transition: all 0.3s ease;
}

.bar-segment:hover {
  opacity: 0.8;
}

/* 响应式设计 */
@media (max-width: 1600px) {
  .dashboard-top {
    grid-template-columns: 1.5fr 1fr;
  }
}

@media (max-width: 1200px) {
  .dashboard-top {
    grid-template-columns: 1fr;
  }
  
  .dashboard-right {
    flex-direction: row;
    gap: 16px;
  }
  
  .chart-body {
    height: 280px;
  }
  
  .ring-chart {
    height: 160px;
  }
}

@media (max-width: 900px) {
  .dashboard-right {
    flex-direction: column;
  }
  
  .stat-value {
    font-size: 24px;
  }
  
  .card-stats {
    flex-direction: column;
    gap: 20px;
  }
}

@media (max-width: 600px) {
  .data-card,
  .main-chart,
  .ring-card,
  .bottom-section {
    padding: 16px;
  }
  
  .chart-body {
    height: 250px;
  }
  
  .ring-chart {
    height: 140px;
  }
  
  .stat-value {
    font-size: 20px;
  }
}
</style>

