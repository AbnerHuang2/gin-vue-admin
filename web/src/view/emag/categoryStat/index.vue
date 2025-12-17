<template>
  <div class="category-stat-container">
    <!-- Tab 切换 -->
    <el-tabs v-model="activeTab" type="border-card" class="stat-tabs">
      <!-- Tab 1: 品类指标 Top20 -->
      <el-tab-pane label="品类指标 Top20" name="top20">
        <div class="gva-search-box">
          <el-form :inline="true" :model="top20Search" class="demo-form-inline" @keyup.enter="onTop20Submit">
            <el-form-item label="快照日期">
              <el-select v-model="top20Search.snapshotDate" placeholder="请选择快照日期" clearable style="width: 180px">
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
                v-model="top20Search.supperHotRate" 
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
                v-model="top20Search.oemSupperHotRate" 
                :min="0" 
                :max="100" 
                :precision="2"
                :step="0.1"
                placeholder="请输入"
                style="width: 140px"
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="search" @click="onTop20Submit">查询</el-button>
              <el-button icon="refresh" @click="onTop20Reset">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <div class="gva-table-box">
          <!-- 图表区域 -->
          <div class="chart-container" v-if="top20Data.length > 0">
            <h4 class="chart-title">品类指标 Top20 对比图</h4>
            <Chart :options="top20ChartOptions" height="400px" />
          </div>

          <!-- 表格区域 -->
          <el-table
            :data="top20Data"
            style="width: 100%"
            tooltip-effect="dark"
            row-key="id"
          >
            <el-table-column align="center" label="排名" type="index" width="70" />
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
          </el-table>
        </div>
      </el-tab-pane>

      <!-- Tab 2: 同比增长排名 -->
      <el-tab-pane label="同比增长排名" name="growth">
        <div class="gva-search-box">
          <el-form :inline="true" :model="growthSearch" class="demo-form-inline">
            <el-form-item label="显示数量">
              <el-input-number 
                v-model="growthSearch.limit" 
                :min="10" 
                :max="100" 
                :step="10"
                style="width: 140px"
              />
            </el-form-item>
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
            <h4 class="chart-title">总数同比增长率排名（横向条形图）</h4>
            <Chart :options="growthChartOptions" height="500px" />
          </div>

          <!-- 表格区域 -->
          <el-table
            :data="growthData"
            style="width: 100%"
            tooltip-effect="dark"
            row-key="categoryId"
          >
            <el-table-column align="center" label="排名" type="index" width="70" />
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
            <el-table-column align="right" label="超热销增长率(%)" width="140">
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
  getCategoryStatTop20,
  getCategoryStatGrowthRank
} from '@/api/emagCategoryStat'

defineOptions({
  name: 'EmagCategoryStat'
})

// Tab 状态
const activeTab = ref('top20')

// 快照日期列表
const snapshotDateList = ref([])

// Top20 搜索条件
const top20Search = reactive({
  snapshotDate: '',
  supperHotRate: 0,
  oemSupperHotRate: 0
})

// Top20 数据
const top20Data = ref([])

// 同比增长搜索条件
const growthSearch = reactive({
  limit: 50
})

// 同比增长数据
const growthData = ref([])
const growthDateInfo = reactive({
  currentDate: '',
  previousDate: ''
})

// 获取快照日期列表
const fetchSnapshotDateList = async () => {
  try {
    const res = await getSnapshotDateList()
    if (res.code === 0) {
      snapshotDateList.value = res.data || []
      // 默认选中最新日期
      if (snapshotDateList.value.length > 0) {
        top20Search.snapshotDate = snapshotDateList.value[0]
        // 自动查询
        fetchTop20Data()
      }
    }
  } catch (error) {
    ElMessage.error('获取快照日期列表失败')
  }
}

// 获取 Top20 数据
const fetchTop20Data = async () => {
  try {
    const res = await getCategoryStatTop20(top20Search)
    if (res.code === 0) {
      top20Data.value = res.data.list || []
    }
  } catch (error) {
    ElMessage.error('获取品类指标Top20失败')
  }
}

// 获取同比增长数据
const fetchGrowthData = async () => {
  try {
    const res = await getCategoryStatGrowthRank(growthSearch)
    if (res.code === 0) {
      growthData.value = res.data.list || []
      growthDateInfo.currentDate = res.data.currentDate || ''
      growthDateInfo.previousDate = res.data.previousDate || ''
    }
  } catch (error) {
    ElMessage.error('获取同比增长排名失败')
  }
}

// Top20 查询
const onTop20Submit = () => {
  fetchTop20Data()
}

// Top20 重置
const onTop20Reset = () => {
  top20Search.supperHotRate = 0
  top20Search.oemSupperHotRate = 0
  if (snapshotDateList.value.length > 0) {
    top20Search.snapshotDate = snapshotDateList.value[0]
  }
  fetchTop20Data()
}

// 同比增长查询
const onGrowthSubmit = () => {
  fetchGrowthData()
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

// Top20 图表配置
const top20ChartOptions = computed(() => {
  if (!top20Data.value || top20Data.value.length === 0) {
    return {}
  }
  const categories = top20Data.value.map(item => item.categoryId || '')
  const supperHotRates = top20Data.value.map(item => {
    const rate = item.supperHotRate
    return rate != null ? (rate * 100).toFixed(2) : 0
  })
  const oemSupperHotRates = top20Data.value.map(item => {
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
    const rate = item.totalGrowthRate
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
        return `${data.name}<br/>总数增长率: ${data.value}%`
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
        name: '总数增长率',
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
  background: #f5f7fa;
  min-height: calc(100vh - 100px);
}

.stat-tabs {
  background: #fff;
  border-radius: 4px;
}

.gva-search-box {
  padding: 16px;
  background: #fafafa;
  border-radius: 4px;
  margin-bottom: 16px;
}

.gva-table-box {
  padding: 16px;
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

