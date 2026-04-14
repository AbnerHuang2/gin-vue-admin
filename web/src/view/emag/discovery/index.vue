<template>
  <div class="discovery-container">
    <!-- Python 服务不可用时提示 -->
    <el-alert
      v-if="serviceUnavailable"
      class="service-alert"
      type="warning"
      :closable="false"
      show-icon
      title="Python 分析服务未连接"
    />

    <div class="gva-search-box">
      <el-form :inline="true" :model="searchForm" class="demo-form-inline" @keyup.enter="onSearch">
        <el-form-item label="类目">
          <el-select v-model="searchForm.category" placeholder="请选择类目" clearable style="width: 140px">
            <el-option
              v-for="item in categoryOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable style="width: 140px">
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="抓取日期">
          <el-date-picker
            v-model="searchForm.crawl_date"
            type="date"
            placeholder="选择日期"
            clearable
            value-format="YYYY-MM-DD"
            style="width: 160px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button type="success" icon="plus" :loading="discoveryLoading" @click="handleTriggerDiscovery">
            {{ discoveryLoading ? '发现中...' : '发现候选品' }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <el-table
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        v-loading="loading"
      >
        <el-table-column align="center" label="序号" type="index" width="60" :index="indexMethod" />
        <el-table-column align="left" label="商品名称" prop="title" min-width="280" show-overflow-tooltip />
        <el-table-column align="center" label="类目" width="100">
          <template #default="scope">
            <el-tag type="info" size="small">{{ getCategoryLabel(scope.row.category) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="关键词" prop="keyword" width="120" show-overflow-tooltip />
        <el-table-column align="right" label="价格" width="100">
          <template #default="scope">
            <span>{{ formatPrice(scope.row.price) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" label="评分" prop="rating" width="80" />
        <el-table-column align="center" label="评论数" prop="review_count" width="100" />
        <el-table-column align="left" label="卖家" prop="seller_name" width="120" show-overflow-tooltip />
        <el-table-column align="center" label="排名" prop="rank_position" width="80" />
        <el-table-column align="center" label="排名变化" width="100">
          <template #default="scope">
            <span :class="getRankChangeClass(scope.row.rank_change)">
              {{ formatRankChange(scope.row.rank_change) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column align="right" label="发现分数" width="100">
          <template #default="scope">
            <span>{{ formatScore(scope.row.discovery_score) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getCandidateStatusTagType(scope.row.status)" size="small">
              {{ getCandidateStatusLabel(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" width="220" fixed="right">
          <template #default="scope">
            <el-button
              v-if="scope.row.status === 'new'"
              type="primary"
              size="small"
              :loading="submittingId === scope.row.id"
              @click="handleSubmitCandidate(scope.row)"
            >
              提交选品
            </el-button>
            <el-button
              v-if="scope.row.status === 'submitted' && scope.row.selection_task_id"
              type="success"
              link
              size="small"
              @click="goToReport(scope.row.selection_task_id)"
            >
              查看报告
            </el-button>
            <el-button
              v-if="scope.row.product_url"
              type="primary"
              link
              size="small"
              @click="openProductUrl(scope.row.product_url)"
            >
              查看详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 25, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getCandidateList, triggerDiscovery, submitCandidate } from '@/api/discovery'

const router = useRouter()

defineOptions({
  name: 'EmagDiscovery'
})

// 类目下拉选项
const categoryOptions = [
  { label: '日用', value: 'daily' },
  { label: '玩具', value: 'toys' },
  { label: '宠物', value: 'pet' },
  { label: '电子', value: 'electronics' },
  { label: '运动', value: 'sports' },
  { label: '家居', value: 'home' }
]

// 状态下拉选项
const statusOptions = [
  { label: '新建', value: 'new' },
  { label: '已提交', value: 'submitted' },
  { label: '已忽略', value: 'ignored' }
]

// 搜索表单
const searchForm = reactive({
  category: '',
  status: '',
  crawl_date: ''
})

// 分页（与后端 offset/limit 对齐）
const page = ref(1)
const pageSize = ref(50)
const total = ref(0)

const tableData = ref([])
const loading = ref(false)
const discoveryLoading = ref(false)
// Python 服务不可用时展示告警
const serviceUnavailable = ref(false)
// 当前正在提交的行 id，用于按钮 loading
const submittingId = ref(null)

// 类目展示文案
const categoryLabelMap = Object.fromEntries(categoryOptions.map((i) => [i.value, i.label]))

const getCategoryLabel = (code) => {
  return categoryLabelMap[code] || code || '-'
}

// 候选状态展示与标签颜色（new=primary, submitted=success, ignored=info）
const candidateStatusLabelMap = {
  new: '新建',
  submitted: '已提交',
  ignored: '已忽略'
}

const candidateStatusTagTypeMap = {
  new: 'primary',
  submitted: 'success',
  ignored: 'info'
}

const getCandidateStatusLabel = (status) => {
  return candidateStatusLabelMap[status] || status || '-'
}

const getCandidateStatusTagType = (status) => {
  return candidateStatusTagTypeMap[status] || 'info'
}

// 价格保留两位小数
const formatPrice = (price) => {
  if (price === null || price === undefined) return '-'
  return Number(price).toFixed(2)
}

// 发现分数保留一位小数
const formatScore = (score) => {
  if (score === null || score === undefined) return '-'
  return Number(score).toFixed(1)
}

// 排名变化展示：正数绿色↑，负数红色↓，0 灰色 -
const formatRankChange = (val) => {
  if (val === null || val === undefined || val === '') return '-'
  const n = Number(val)
  if (Number.isNaN(n)) return '-'
  if (n > 0) return `↑${n}`
  if (n < 0) return `↓${Math.abs(n)}`
  return '-'
}

const getRankChangeClass = (val) => {
  if (val === null || val === undefined || val === '') return 'rank-change-zero'
  const n = Number(val)
  if (Number.isNaN(n)) return 'rank-change-zero'
  if (n > 0) return 'rank-change-up'
  if (n < 0) return 'rank-change-down'
  return 'rank-change-zero'
}

// 表格序号（跨页连续）
const indexMethod = (index) => {
  return (page.value - 1) * pageSize.value + index + 1
}

// 组装列表查询参数
const buildListParams = () => {
  const params = {
    limit: pageSize.value,
    offset: (page.value - 1) * pageSize.value
  }
  if (searchForm.category) params.category = searchForm.category
  if (searchForm.status) params.status = searchForm.status
  if (searchForm.crawl_date) params.crawl_date = searchForm.crawl_date
  return params
}

// 拉取候选品列表
const fetchCandidateList = async () => {
  loading.value = true
  try {
    const res = await getCandidateList(buildListParams())
    tableData.value = res.items || []
    total.value = typeof res.total === 'number' ? res.total : 0
    serviceUnavailable.value = false
  } catch (error) {
    console.error('获取候选品列表失败', error)
    tableData.value = []
    total.value = 0
    serviceUnavailable.value = true
  } finally {
    loading.value = false
  }
}

// 查询
const onSearch = () => {
  page.value = 1
  fetchCandidateList()
}

// 重置筛选并刷新
const onReset = () => {
  searchForm.category = ''
  searchForm.status = ''
  searchForm.crawl_date = ''
  page.value = 1
  fetchCandidateList()
}

// 触发发现任务
const handleTriggerDiscovery = async () => {
  discoveryLoading.value = true
  try {
    await triggerDiscovery()
    ElMessage.success('发现任务已触发')
    await fetchCandidateList()
  } catch (error) {
    console.error('触发发现失败', error)
    ElMessage.error('触发发现失败，请检查 Python 服务')
  } finally {
    discoveryLoading.value = false
  }
}

// 提交单个候选为选品任务
const handleSubmitCandidate = async (row) => {
  if (!row?.id) return
  submittingId.value = row.id
  try {
    const res = await submitCandidate(row.id)
    ElMessage.success(`已提交，任务 ID：${res.task_id || ''}`)
    await fetchCandidateList()
  } catch (error) {
    console.error('提交选品失败', error)
    ElMessage.error('提交选品失败')
  } finally {
    submittingId.value = null
  }
}

const goToReport = (taskId) => {
  router.push({ name: 'emagSelection', query: { task_id: taskId } })
}

const openProductUrl = (url) => {
  if (!url) return
  window.open(url, '_blank')
}

const handleSizeChange = (val) => {
  pageSize.value = val
  page.value = 1
  fetchCandidateList()
}

const handleCurrentChange = (val) => {
  page.value = val
  fetchCandidateList()
}

// 首屏加载：失败不抛错，空表 + 告警
onMounted(async () => {
  try {
    await fetchCandidateList()
  } catch (e) {
    // fetchCandidateList 内部已吞掉异常，此处双保险
    tableData.value = []
    serviceUnavailable.value = true
  }
})
</script>

<style scoped lang="scss">
.discovery-container {
  padding: 16px;
  min-height: calc(100vh - 100px);
}

.service-alert {
  margin-bottom: 16px;
}

.gva-search-box {
  padding: 16px;
  border-radius: 4px;
  margin-bottom: 16px;
  background: var(--el-bg-color);
}

.gva-table-box {
  padding: 16px;
  background: var(--el-bg-color);
  border-radius: 4px;
}

.gva-pagination {
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
}

.rank-change-up {
  color: #67c23a;
  font-weight: 600;
}

.rank-change-down {
  color: #f56c6c;
  font-weight: 600;
}

.rank-change-zero {
  color: #909399;
}
</style>
