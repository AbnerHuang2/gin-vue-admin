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

    <el-tabs v-model="activeTab" type="border-card">
      <!-- ====== Tab 1: 候选商品 ====== -->
      <el-tab-pane label="候选商品" name="candidates">
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
              <el-button type="success" icon="plus" :loading="discoveryLoading" :disabled="discoveryLoading" @click="showDiscoveryDialog">
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
            <el-table-column align="center" label="主图" width="80">
              <template #default="scope">
                <el-image
                  v-if="scope.row.main_image_url"
                  :src="scope.row.main_image_url"
                  :preview-src-list="[scope.row.main_image_url]"
                  preview-teleported
                  fit="contain"
                  style="width: 50px; height: 50px"
                />
                <span v-else class="no-data">-</span>
              </template>
            </el-table-column>
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
            <el-table-column align="center" label="操作" width="280" fixed="right">
              <template #default="scope">
                <el-button
                  v-if="scope.row.status === 'new' || scope.row.status === 'pending'"
                  type="primary"
                  size="small"
                  :loading="submittingId === scope.row.id"
                  @click="handleSubmitCandidate(scope.row)"
                >
                  提交选品
                </el-button>
                <el-button
                  v-if="scope.row.status === 'new' || scope.row.status === 'pending'"
                  type="warning"
                  size="small"
                  :loading="analyzingId === scope.row.id"
                  @click="handleAnalyzeCandidate(scope.row)"
                >
                  加入选品分析
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
      </el-tab-pane>

      <!-- ====== Tab 2: 品类关键词 ====== -->
      <el-tab-pane label="品类关键词" name="keywords">
        <div class="gva-search-box">
          <el-form :inline="true" :model="kwSearchForm">
            <el-form-item label="类目">
              <el-select v-model="kwSearchForm.category" placeholder="全部类目" clearable style="width: 140px">
                <el-option
                  v-for="item in categoryOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="search" @click="fetchKeywordList">查询</el-button>
              <el-button type="success" icon="plus" @click="showKwAddDialog">新增关键词</el-button>
            </el-form-item>
          </el-form>
        </div>

        <div class="gva-table-box">
          <el-table :data="kwTableData" style="width: 100%" v-loading="kwLoading">
            <el-table-column align="center" label="序号" type="index" width="60" />
            <el-table-column align="center" label="类目" width="120">
              <template #default="scope">
                <el-tag type="info" size="small">{{ getCategoryLabel(scope.row.category) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column align="left" label="关键词" prop="keyword" min-width="250" />
            <el-table-column align="center" label="状态" width="100">
              <template #default="scope">
                <el-switch
                  v-model="scope.row.is_active"
                  @change="handleToggleKwActive(scope.row)"
                />
              </template>
            </el-table-column>
            <el-table-column align="center" label="创建时间" width="180">
              <template #default="scope">
                <span>{{ formatDateTime(scope.row.ctime) }}</span>
              </template>
            </el-table-column>
            <el-table-column align="center" label="操作" width="160" fixed="right">
              <template #default="scope">
                <el-button type="primary" link size="small" @click="showKwEditDialog(scope.row)">编辑</el-button>
                <el-button type="danger" link size="small" @click="handleDeleteKw(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- 发现候选品 - 类目选择弹窗 -->
    <el-dialog v-model="discoveryDialogVisible" title="选择抓取类目" width="420px" :close-on-click-modal="false">
      <div style="margin-bottom: 12px">
        <el-checkbox v-model="selectAllCategories" @change="handleSelectAll">全选</el-checkbox>
      </div>
      <el-checkbox-group v-model="selectedCategories">
        <el-checkbox
          v-for="item in categoryOptions"
          :key="item.value"
          :label="item.value"
          :value="item.value"
          style="display: block; margin-bottom: 8px"
        >
          {{ item.label }}
        </el-checkbox>
      </el-checkbox-group>
      <template #footer>
        <el-button @click="discoveryDialogVisible = false">取消</el-button>
        <el-button type="primary" :disabled="selectedCategories.length === 0" @click="handleConfirmDiscovery">
          确认
        </el-button>
      </template>
    </el-dialog>

    <!-- 新增/编辑关键词弹窗 -->
    <el-dialog
      v-model="kwDialogVisible"
      :title="kwDialogMode === 'add' ? '新增关键词' : '编辑关键词'"
      width="460px"
      :close-on-click-modal="false"
    >
      <el-form :model="kwForm" label-width="80px">
        <el-form-item label="类目" required>
          <el-select v-model="kwForm.category" placeholder="请选择类目" style="width: 100%">
            <el-option
              v-for="item in categoryOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词" required>
          <el-input v-model="kwForm.keyword" placeholder="请输入搜索关键词（罗语）" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="kwForm.is_active" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="kwDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="kwSubmitting" @click="handleKwSubmit">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getCandidateList, triggerDiscovery, submitCandidate, analyzeCandidate,
  getKeywordList, createKeyword, updateKeyword, deleteKeyword,
} from '@/api/discovery'

const router = useRouter()

defineOptions({
  name: 'EmagDiscovery'
})

const activeTab = ref('candidates')

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

// ── 候选商品相关 ──

const searchForm = reactive({
  category: '',
  status: '',
  crawl_date: ''
})

const page = ref(1)
const pageSize = ref(50)
const total = ref(0)

const tableData = ref([])
const loading = ref(false)
const discoveryLoading = ref(false)
const serviceUnavailable = ref(false)
const submittingId = ref(null)
const analyzingId = ref(null)

const discoveryDialogVisible = ref(false)
const selectedCategories = ref([])
const selectAllCategories = ref(false)

watch(selectedCategories, (val) => {
  selectAllCategories.value = val.length === categoryOptions.length
})

const handleSelectAll = (val) => {
  selectedCategories.value = val ? categoryOptions.map((i) => i.value) : []
}

const showDiscoveryDialog = () => {
  selectedCategories.value = []
  selectAllCategories.value = false
  discoveryDialogVisible.value = true
}

const categoryLabelMap = Object.fromEntries(categoryOptions.map((i) => [i.value, i.label]))

const getCategoryLabel = (code) => {
  return categoryLabelMap[code] || code || '-'
}

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

const formatPrice = (price) => {
  if (price === null || price === undefined) return '-'
  return Number(price).toFixed(2)
}

const formatScore = (score) => {
  if (score === null || score === undefined) return '-'
  return Number(score).toFixed(1)
}

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

const formatDateTime = (dt) => {
  if (!dt) return '-'
  return String(dt).replace('T', ' ').slice(0, 19)
}

const indexMethod = (index) => {
  return (page.value - 1) * pageSize.value + index + 1
}

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

const onSearch = () => {
  page.value = 1
  fetchCandidateList()
}

const onReset = () => {
  searchForm.category = ''
  searchForm.status = ''
  searchForm.crawl_date = ''
  page.value = 1
  fetchCandidateList()
}

const handleConfirmDiscovery = async () => {
  discoveryDialogVisible.value = false
  discoveryLoading.value = true
  const minWait = new Promise((r) => setTimeout(r, 3000))
  try {
    await triggerDiscovery(selectedCategories.value)
    ElMessage.success('发现任务已触发')
  } catch (error) {
    console.error('触发发现失败', error)
    ElMessage.error('触发发现失败，请检查 Python 服务')
  } finally {
    await minWait
    discoveryLoading.value = false
    await fetchCandidateList()
  }
}

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

const handleAnalyzeCandidate = async (row) => {
  if (!row?.id) return
  analyzingId.value = row.id
  try {
    const res = await analyzeCandidate(row.id)
    ElMessage.success('选品分析已启动')
    await fetchCandidateList()
    if (res.task_id) {
      router.push({ name: 'emagSelection', query: { task_id: res.task_id } })
    }
  } catch (error) {
    console.error('加入选品分析失败', error)
    ElMessage.error(error?.response?.data?.detail || '加入选品分析失败')
  } finally {
    analyzingId.value = null
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

// ── 品类关键词相关 ──

const kwSearchForm = reactive({ category: '' })
const kwTableData = ref([])
const kwLoading = ref(false)

const kwDialogVisible = ref(false)
const kwDialogMode = ref('add')
const kwSubmitting = ref(false)
const kwEditingId = ref(null)
const kwForm = reactive({
  category: '',
  keyword: '',
  is_active: true,
})

const fetchKeywordList = async () => {
  kwLoading.value = true
  try {
    const params = {}
    if (kwSearchForm.category) params.category = kwSearchForm.category
    const res = await getKeywordList(params)
    kwTableData.value = res.items || []
    serviceUnavailable.value = false
  } catch (error) {
    console.error('获取关键词列表失败', error)
    kwTableData.value = []
    serviceUnavailable.value = true
  } finally {
    kwLoading.value = false
  }
}

const showKwAddDialog = () => {
  kwDialogMode.value = 'add'
  kwEditingId.value = null
  kwForm.category = ''
  kwForm.keyword = ''
  kwForm.is_active = true
  kwDialogVisible.value = true
}

const showKwEditDialog = (row) => {
  kwDialogMode.value = 'edit'
  kwEditingId.value = row.id
  kwForm.category = row.category
  kwForm.keyword = row.keyword
  kwForm.is_active = row.is_active
  kwDialogVisible.value = true
}

const handleKwSubmit = async () => {
  if (!kwForm.category || !kwForm.keyword) {
    ElMessage.warning('请填写类目和关键词')
    return
  }
  kwSubmitting.value = true
  try {
    if (kwDialogMode.value === 'add') {
      await createKeyword({
        category: kwForm.category,
        keyword: kwForm.keyword.trim(),
        is_active: kwForm.is_active,
      })
      ElMessage.success('新增成功')
    } else {
      await updateKeyword(kwEditingId.value, {
        category: kwForm.category,
        keyword: kwForm.keyword.trim(),
        is_active: kwForm.is_active,
      })
      ElMessage.success('更新成功')
    }
    kwDialogVisible.value = false
    await fetchKeywordList()
  } catch (error) {
    const msg = error?.response?.data?.detail || (kwDialogMode.value === 'add' ? '新增失败' : '更新失败')
    ElMessage.error(msg)
  } finally {
    kwSubmitting.value = false
  }
}

const handleToggleKwActive = async (row) => {
  try {
    await updateKeyword(row.id, { is_active: row.is_active })
    ElMessage.success(row.is_active ? '已启用' : '已禁用')
  } catch (error) {
    row.is_active = !row.is_active
    ElMessage.error('状态更新失败')
  }
}

const handleDeleteKw = async (row) => {
  try {
    await ElMessageBox.confirm(`确定删除关键词「${row.keyword}」？`, '确认删除', {
      type: 'warning',
      confirmButtonText: '删除',
      cancelButtonText: '取消',
    })
    await deleteKeyword(row.id)
    ElMessage.success('已删除')
    await fetchKeywordList()
  } catch {
    // 用户取消
  }
}

// 切换 tab 时懒加载关键词数据
watch(activeTab, (tab) => {
  if (tab === 'keywords' && kwTableData.value.length === 0) {
    fetchKeywordList()
  }
})

onMounted(async () => {
  try {
    await fetchCandidateList()
  } catch (e) {
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
