<template>
  <div class="category-stat-container">
    <!-- Tab 切换 -->
    <el-tabs v-model="activeTab" type="border-card" class="stat-tabs">
      <!-- Tab 1: 品类指标列表 -->
      <el-tab-pane label="品类指标列表" name="list">
        <div class="gva-search-box">
          <el-form :inline="true" :model="listSearch" class="demo-form-inline" @keyup.enter="onListSubmit">
            <el-form-item label="快照日期">
              <el-select v-model="listSearch.snapshotDate" placeholder="请选择快照日期" clearable style="width: 180px">
                <el-option
                  v-for="date in snapshotDateList"
                  :key="date"
                  :label="date"
                  :value="date"
                />
              </el-select>
            </el-form-item>
            <el-form-item label="超热销率 ≥">
              <el-input-number 
                v-model="listSearch.supperHotRate" 
                :min="0" 
                :max="100" 
                :precision="2"
                :step="0.1"
                placeholder="请输入"
                style="width: 140px"
              />
            </el-form-item>
            <el-form-item label="OEM超热销率 ≥">
              <el-input-number 
                v-model="listSearch.oemSupperHotRate" 
                :min="0" 
                :max="100" 
                :precision="2"
                :step="0.1"
                placeholder="请输入"
                style="width: 140px"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="search" @click="onListSubmit">查询</el-button>
              <el-button icon="refresh" @click="onListReset">重置</el-button>
              <el-button type="warning" icon="refresh" :loading="updateTaskLoading" @click="handleTriggerUpdate">
                {{ updateTaskLoading ? '更新中...' : '手动更新数据' }}
              </el-button>
            </el-form-item>
          </el-form>
        </div>

        <div class="gva-table-box">
          <!-- 图表区域 -->
          <div class="chart-container" v-if="listData.length > 0">
            <h4 class="chart-title">品类指标对比图（当前页）</h4>
            <Chart :options="listChartOptions" height="400px" />
          </div>

          <!-- 表格区域 -->
          <el-table
            :data="listData"
            style="width: 100%"
            tooltip-effect="dark"
            row-key="id"
          >
            <el-table-column align="center" label="序号" type="index" width="70" :index="listIndexMethod" />
            <el-table-column align="left" label="品类ID" prop="categoryId" width="150" />
            <el-table-column align="left" label="品类名称" prop="categoryName" min-width="150" show-overflow-tooltip />
            <el-table-column align="left" label="子品类名称" prop="subcategoryName" min-width="150" show-overflow-tooltip />
            <el-table-column align="right" label="总数" prop="total" width="100" />
            <el-table-column align="right" label="超热销总数" prop="supperHotTotal" width="120" />
            <el-table-column align="right" label="OEM总数" prop="oemTotal" width="100" />
            <el-table-column align="right" label="OEM超热销总数" prop="oemSupperHotTotal" width="140" />
            <el-table-column align="right" label="超热销率(%)" width="120">
              <template #default="scope">
                <span :class="getRateClass(scope.row.supperHotRate)">
                  {{ formatRate(scope.row.supperHotRate) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column align="right" label="OEM超热销率(%)" width="140">
              <template #default="scope">
                <span :class="getRateClass(scope.row.oemSupperHotRate)">
                  {{ formatRate(scope.row.oemSupperHotRate) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column align="left" label="快照日期" width="120">
              <template #default="scope">
                <span>{{ formatSnapshotDate(scope.row.snapshotDate) }}</span>
              </template>
            </el-table-column>
            <el-table-column align="left" label="标签" prop="tags" min-width="100" />
            <el-table-column align="center" label="操作" width="100" fixed="right">
              <template #default="scope">
                <el-popconfirm
                  title="确定要标记为不关注吗？标记后将不再显示此品类"
                  confirm-button-text="确定"
                  cancel-button-text="取消"
                  @confirm="handleMarkAsNotCare(scope.row)"
                >
                  <template #reference>
                    <el-button type="danger" link size="small">不关注</el-button>
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
              :current-page="listPage"
              :page-size="listPageSize"
              :page-sizes="[10, 30, 50, 100]"
              :total="listTotal"
              layout="total, sizes, prev, pager, next, jumper"
              @current-change="handleListCurrentChange"
              @size-change="handleListSizeChange"
            />
          </div>
        </div>
      </el-tab-pane>

      <!-- Tab 2: 热销增长率排名 -->
      <el-tab-pane label="热销增长率排名" name="growth">
        <div class="gva-search-box">
          <el-form :inline="true" class="demo-form-inline">
            <el-form-item>
              <el-button type="primary" icon="search" @click="onGrowthSubmit">查询</el-button>
            </el-form-item>
          </el-form>
          <div class="date-info" v-if="growthDateInfo.currentDate">
            <el-tag type="info">
              对比日期：{{ growthDateInfo.previousDate }} → {{ growthDateInfo.currentDate }}
            </el-tag>
          </div>
        </div>

        <div class="gva-table-box">
          <!-- 图表区域 -->
          <div class="chart-container" v-if="growthData.length > 0">
            <h4 class="chart-title">热销增长率排名（当前页前20条）</h4>
            <Chart :options="growthChartOptions" height="500px" />
          </div>

          <!-- 表格区域 -->
          <el-table
            :data="growthData"
            style="width: 100%"
            tooltip-effect="dark"
            row-key="categoryId"
          >
            <el-table-column align="center" label="序号" type="index" width="70" :index="growthIndexMethod" />
            <el-table-column align="left" label="品类ID" prop="categoryId" width="150" />
            <el-table-column align="left" label="品类名称" prop="categoryName" min-width="150" show-overflow-tooltip />
            <el-table-column align="left" label="子品类名称" prop="subcategoryName" min-width="150" show-overflow-tooltip />
            <el-table-column align="right" label="当前总数" prop="currentTotal" width="100" />
            <el-table-column align="right" label="上期总数" prop="previousTotal" width="100" />
            <el-table-column align="right" label="总数增长率(%)" width="130">
              <template #default="scope">
                <span :class="getGrowthClass(scope.row.totalGrowthRate)">
                  {{ formatGrowthRate(scope.row.totalGrowthRate) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column align="right" label="当前超热销" prop="currentSupperHotTotal" width="110" />
            <el-table-column align="right" label="上期超热销" prop="previousSupperHotTotal" width="110" />
            <el-table-column align="right" label="热销增长率(%)" width="140">
              <template #default="scope">
                <span :class="getGrowthClass(scope.row.supperHotGrowthRate)">
                  {{ formatGrowthRate(scope.row.supperHotGrowthRate) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column align="right" label="当前超热销率(%)" width="140">
              <template #default="scope">
                {{ formatRate(scope.row.currentSupperHotRate) }}
              </template>
            </el-table-column>
            <el-table-column align="right" label="上期超热销率(%)" width="140">
              <template #default="scope">
                {{ formatRate(scope.row.previousSupperHotRate) }}
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination
              :current-page="growthPage"
              :page-size="growthPageSize"
              :page-sizes="[10, 30, 50, 100]"
              :total="growthTotal"
              layout="total, sizes, prev, pager, next, jumper"
              @current-change="handleGrowthCurrentChange"
              @size-change="handleGrowthSizeChange"
            />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import Chart from '@/components/charts/index.vue'
import {
  getSnapshotDateList,
  getCategoryStatList,
  getCategoryStatGrowthRank,
  triggerUpdateTask,
  markAsNotCare
} from '@/api/emagCategoryStat'

defineOptions({
  name: 'EmagCategoryStat'
})

// Tab 状态
const activeTab = ref('list')

// 快照日期列表
const snapshotDateList = ref([])

// 品类指标列表搜索条件
const listSearch = reactive({
  snapshotDate: '',
  supperHotRate: 0.10,
  oemSupperHotRate: 0.05
})

// 品类指标列表分页
const listPage = ref(1)
const listPageSize = ref(10)
const listTotal = ref(0)

// 品类指标列表数据
const listData = ref([])

// 同比增长分页
const growthPage = ref(1)
const growthPageSize = ref(10)
const growthTotal = ref(0)

// 同比增长数据
const growthData = ref([])
const growthDateInfo = reactive({
  currentDate: '',
  previousDate: ''
})

// 更新任务状态
const updateTaskLoading = ref(false)

// 获取快照日期列表
const fetchSnapshotDateList = async () => {
  try {
    const res = await getSnapshotDateList()
    if (res.code === 0) {
      snapshotDateList.value = res.data || []
      // 默认选中最新日期
      if (snapshotDateList.value.length > 0) {
        listSearch.snapshotDate = snapshotDateList.value[0]
        // 自动查询
        fetchListData()
      }
    }
  } catch (error) {
    ElMessage.error('获取快照日期列表失败')
  }
}

// 获取品类指标列表数据
const fetchListData = async () => {
  try {
    const res = await getCategoryStatList({
      page: listPage.value,
      pageSize: listPageSize.value,
      ...listSearch
    })
    if (res.code === 0) {
      listData.value = res.data.list || []
      listTotal.value = res.data.total
    }
  } catch (error) {
    ElMessage.error('获取品类指标列表失败')
  }
}

// 获取同比增长数据
const fetchGrowthData = async () => {
  try {
    const res = await getCategoryStatGrowthRank({
      page: growthPage.value,
      pageSize: growthPageSize.value
    })
    if (res.code === 0) {
      growthData.value = res.data.list || []
      growthTotal.value = res.data.total || 0
      growthDateInfo.currentDate = res.data.currentDate || ''
      growthDateInfo.previousDate = res.data.previousDate || ''
    }
  } catch (error) {
    ElMessage.error('获取热销增长率排名失败')
  }
}

// 品类指标列表查询
const onListSubmit = () => {
  listPage.value = 1
  fetchListData()
}

// 品类指标列表重置
const onListReset = () => {
  listSearch.supperHotRate = 0
  listSearch.oemSupperHotRate = 0
  if (snapshotDateList.value.length > 0) {
    listSearch.snapshotDate = snapshotDateList.value[0]
  }
  listPage.value = 1
  fetchListData()
}

// 品类指标列表分页
const handleListSizeChange = (val) => {
  listPageSize.value = val
  fetchListData()
}

const handleListCurrentChange = (val) => {
  listPage.value = val
  fetchListData()
}

// 品类指标列表序号方法
const listIndexMethod = (index) => {
  return (listPage.value - 1) * listPageSize.value + index + 1
}

// 手动触发更新任务
const handleTriggerUpdate = async () => {
  try {
    updateTaskLoading.value = true
    const res = await triggerUpdateTask()
    if (res.code === 0) {
      ElMessage.success(res.msg || '任务已触发，请稍后刷新页面查看结果')
    } else {
      ElMessage.error(res.msg || '触发任务失败')
    }
  } catch (error) {
    ElMessage.error('触发任务失败')
  } finally {
    // 3秒后恢复按钮状态
    setTimeout(() => {
      updateTaskLoading.value = false
    }, 3000)
  }
}

// 标记品类为不关注
const handleMarkAsNotCare = async (row) => {
  try {
    const res = await markAsNotCare(row.categoryId)
    if (res.code === 0) {
      ElMessage.success('标记成功')
      // 刷新列表
      fetchListData()
    } else {
      ElMessage.error(res.msg || '标记失败')
    }
  } catch (error) {
    ElMessage.error('标记失败')
  }
}

// 同比增长查询
const onGrowthSubmit = () => {
  growthPage.value = 1
  fetchGrowthData()
}

// 同比增长分页
const handleGrowthSizeChange = (val) => {
  growthPageSize.value = val
  fetchGrowthData()
}

const handleGrowthCurrentChange = (val) => {
  growthPage.value = val
  fetchGrowthData()
}

// 同比增长序号方法
const growthIndexMethod = (index) => {
  return (growthPage.value - 1) * growthPageSize.value + index + 1
}

// 格式化比率
const formatRate = (rate) => {
  if (rate === null || rate === undefined) return '-'
  return (rate * 100).toFixed(2)
}

// 格式化增长率
const formatGrowthRate = (rate) => {
  if (rate === null || rate === undefined) return '-'
  const prefix = rate > 0 ? '+' : ''
  return prefix + rate.toFixed(2)
}

// 格式化快照日期
const formatSnapshotDate = (date) => {
  if (!date) return '-'
  return date.split('T')[0]
}

// 获取比率样式类
const getRateClass = (rate) => {
  if (rate >= 0.3) return 'rate-high'
  if (rate >= 0.1) return 'rate-medium'
  return 'rate-low'
}

// 获取增长率样式类
const getGrowthClass = (rate) => {
  if (rate > 0) return 'growth-positive'
  if (rate < 0) return 'growth-negative'
  return ''
}

// 品类指标列表图表配置
const listChartOptions = computed(() => {
  if (!listData.value || listData.value.length === 0) {
    return {}
  }
  const categories = listData.value.map(item => item.categoryId || '')
  const supperHotRates = listData.value.map(item => {
    const rate = item.supperHotRate
    return rate != null ? (rate * 100).toFixed(2) : 0
  })
  const oemSupperHotRates = listData.value.map(item => {
    const rate = item.oemSupperHotRate
    return rate != null ? (rate * 100).toFixed(2) : 0
  })

  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      data: ['超热销率(%)', 'OEM超热销率(%)']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: categories,
      axisLabel: {
        rotate: 45,
        interval: 0
      }
    },
    yAxis: {
      type: 'value',
      name: '比率(%)',
      axisLabel: {
        formatter: '{value}%'
      }
    },
    series: [
      {
        name: '超热销率(%)',
        type: 'bar',
        data: supperHotRates,
        itemStyle: {
          color: '#409EFF'
        }
      },
      {
        name: 'OEM超热销率(%)',
        type: 'bar',
        data: oemSupperHotRates,
        itemStyle: {
          color: '#67C23A'
        }
      }
    ]
  }
})

// 同比增长图表配置 - 横向条形图
const growthChartOptions = computed(() => {
  if (!growthData.value || growthData.value.length === 0) {
    return {}
  }
  // 取前20条用于图表展示
  const displayData = growthData.value.slice(0, 20)
  const categories = displayData.map(item => item.categoryId || '')
  const growthRates = displayData.map(item => {
    const rate = item.supperHotGrowthRate
    return rate != null ? Number(rate.toFixed(2)) : 0
  })

  // 复制数组再 reverse，避免修改原数组
  const reversedCategories = [...categories].reverse()
  const reversedGrowthRates = [...growthRates].reverse()

  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      },
      formatter: (params) => {
        const data = params[0]
        return `${data.name}<br/>热销增长率: ${data.value}%`
      }
    },
    grid: {
      left: '15%',
      right: '10%',
      bottom: '5%',
      top: '5%',
      containLabel: true
    },
    xAxis: {
      type: 'value',
      name: '增长率(%)',
      axisLabel: {
        formatter: '{value}%'
      }
    },
    yAxis: {
      type: 'category',
      data: reversedCategories,
      axisLabel: {
        interval: 0
      }
    },
    series: [
      {
        name: '热销增长率',
        type: 'bar',
        data: reversedGrowthRates,
        itemStyle: {
          color: (params) => {
            return params.value >= 0 ? '#67C23A' : '#F56C6C'
          }
        },
        label: {
          show: true,
          position: 'right',
          formatter: '{c}%'
        }
      }
    ]
  }
})

// 初始化
onMounted(() => {
  fetchSnapshotDateList()
  fetchGrowthData()
})
</script>

<style scoped lang="scss">
.category-stat-container {
  padding: 16px;
  //background: #f5f7fa;
  min-height: calc(100vh - 100px);
}

.stat-tabs {
  //background: #fff;
  border-radius: 4px;
}

.gva-search-box {
  padding: 16px;
  //background: #fafafa;
  border-radius: 4px;
  margin-bottom: 16px;
}

.gva-table-box {
  padding: 16px;
}

.gva-pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.chart-container {
  margin-bottom: 24px;
  padding: 16px;
  background: #fafafa;
  border-radius: 8px;
}

.chart-title {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.date-info {
  display: inline-block;
  margin-left: 16px;
}

.rate-high {
  color: #67C23A;
  font-weight: 600;
}

.rate-medium {
  color: #E6A23C;
  font-weight: 500;
}

.rate-low {
  color: #909399;
}

.growth-positive {
  color: #67C23A;
  font-weight: 600;
}

.growth-negative {
  color: #F56C6C;
  font-weight: 600;
}
</style>
