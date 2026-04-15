<template>
  <div class="selection-container">
    <!-- Python 选品服务不可用时提示 -->
    <el-alert
      v-if="pyBackendError"
      class="backend-alert"
      type="warning"
      :title="pyBackendError"
      show-icon
      :closable="false"
    />

    <el-tabs v-model="activeTab" class="selection-tabs">
      <!-- Tab1: 发起分析 -->
      <el-tab-pane label="发起分析" name="analyze">
        <div class="gva-search-box">
          <el-form :inline="true" :model="analyzeForm" @submit.prevent>
            <el-form-item label="商品链接">
              <el-input
                v-model="analyzeForm.url"
                placeholder="请输入 emag 商品链接"
                clearable
                style="width: 420px"
                :disabled="analyzeSubmitting"
              />
            </el-form-item>
            <el-form-item>
              <el-button
                type="primary"
                :loading="analyzeSubmitting"
                @click="handleStartAnalyze"
              >
                开始分析
              </el-button>
            </el-form-item>
          </el-form>
        </div>

        <!-- 分析进度 -->
        <div v-if="currentTaskId" class="progress-section gva-table-box">
          <div class="progress-header">
            <span class="label">任务ID</span>
            <el-text class="task-id" truncated>{{ currentTaskId }}</el-text>
            <el-tag v-if="taskStatus.status" :type="getStatusTagType(taskStatus.status)" size="small">
              {{ taskStatus.status }}
            </el-tag>
          </div>
          <el-progress
            v-if="taskStatus.status !== 'failed'"
            :percentage="getProgressPercent(taskStatus.status)"
            :status="taskStatus.status === 'done' ? 'success' : undefined"
          />
          <div v-if="taskStatus.status === 'failed'" class="fail-msg">
            <el-text type="danger">{{ taskStatus.error_msg || '任务失败' }}</el-text>
            <div v-if="taskStatus.failure_detail" class="fail-detail">
              {{ taskStatus.failure_detail }}
            </div>
          </div>
          <div v-if="taskStatus.critic_loops != null" class="meta-line">
            <el-text size="small" type="info">Critic 轮次: {{ taskStatus.critic_loops }}</el-text>
          </div>
        </div>

        <!-- 报告展示 -->
        <template v-if="reportData">
          <!-- 市场分析 -->
          <el-card v-if="reportData.emag_market" class="report-card" shadow="never">
            <template #header>市场分析</template>
            <el-descriptions :column="2" border size="small">
              <el-descriptions-item label="搜索词">
                {{ reportData.emag_market.search_keyword || '-' }}
              </el-descriptions-item>
              <el-descriptions-item label="竞争度">
                {{ reportData.emag_market.competition_level || '-' }}
              </el-descriptions-item>
              <el-descriptions-item label="商品数">
                {{ reportData.emag_market.total_products ?? '-' }}
              </el-descriptions-item>
              <el-descriptions-item label="卖家数">
                {{ reportData.emag_market.total_sellers ?? '-' }}
              </el-descriptions-item>
              <el-descriptions-item label="均价 (EUR)">
                {{ formatNum(reportData.emag_market.avg_price_eur) }}
              </el-descriptions-item>
              <el-descriptions-item label="均价 (RON)">
                {{ formatNum(reportData.emag_market.avg_price_ron) }}
              </el-descriptions-item>
              <el-descriptions-item label="价格区间 (RON)" :span="2">
                {{ formatNum(reportData.emag_market.min_price_ron) }} ~
                {{ formatNum(reportData.emag_market.max_price_ron) }}
              </el-descriptions-item>
              <el-descriptions-item label="平均评分">
                {{ formatNum(reportData.emag_market.avg_rating) }}
              </el-descriptions-item>
              <el-descriptions-item label="平均评论数">
                {{ reportData.emag_market.avg_review_count ?? '-' }}
              </el-descriptions-item>
              <el-descriptions-item label="OEM 卖家比例" :span="2">
                {{ formatPercentRatio(reportData.emag_market.oem_seller_ratio) }}
              </el-descriptions-item>
            </el-descriptions>

            <div v-if="reportData.emag_top5?.length" class="sub-table-wrap">
              <div class="sub-table-title">eMag Top5</div>
              <el-table :data="reportData.emag_top5" size="small" border stripe>
                <el-table-column prop="title" label="标题" min-width="160" show-overflow-tooltip />
                <el-table-column prop="price" label="价格" width="90" align="right">
                  <template #default="{ row }">{{ formatNum(row.price) }}</template>
                </el-table-column>
                <el-table-column prop="rating" label="评分" width="70" align="center" />
                <el-table-column prop="review_count" label="评论数" width="80" align="right" />
                <el-table-column prop="seller_name" label="卖家" width="100" show-overflow-tooltip />
                <el-table-column label="链接" width="80" align="center">
                  <template #default="{ row }">
                    <el-button v-if="row.product_url" type="primary" link size="small" @click="openUrl(row.product_url)">
                      打开
                    </el-button>
                    <span v-else>-</span>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-card>

          <!-- 利润评估 -->
          <el-card class="report-card" shadow="never">
            <template #header>利润评估</template>
            <el-descriptions :column="2" border size="small">
              <el-descriptions-item label="利润 (EUR)">
                {{ formatNum(reportData.profit_eur) }}
              </el-descriptions-item>
              <el-descriptions-item label="利润率">
                {{ formatPercentRatio(reportData.profit_margin) }}
              </el-descriptions-item>
            </el-descriptions>
            <template v-if="reportData.cost_breakdown && Object.keys(reportData.cost_breakdown).length">
              <div class="sub-table-title">成本明细</div>
              <el-descriptions :column="2" border size="small">
                <el-descriptions-item
                  v-for="(val, key) in reportData.cost_breakdown"
                  :key="key"
                  :label="formatCostKey(key)"
                >
                  {{ formatNum(val) }}
                </el-descriptions-item>
              </el-descriptions>
            </template>

            <div v-if="reportData.suppliers_top5?.length" class="sub-table-wrap">
              <div class="sub-table-title">1688 供应商 Top5</div>
              <el-table :data="reportData.suppliers_top5" size="small" border stripe>
                <el-table-column prop="supplier_name" label="供应商" width="110" show-overflow-tooltip />
                <el-table-column prop="product_title" label="产品" min-width="120" show-overflow-tooltip />
                <el-table-column prop="unit_cost_cny" label="单价(CNY)" width="95" align="right">
                  <template #default="{ row }">{{ formatNum(row.unit_cost_cny) }}</template>
                </el-table-column>
                <el-table-column prop="moq" label="MOQ" width="70" align="right" />
                <el-table-column prop="supplier_rating" label="店铺评分" width="85" align="center" />
                <el-table-column prop="reliability_score" label="可靠度" width="75" align="center" />
                <el-table-column label="链接" width="80" align="center">
                  <template #default="{ row }">
                    <el-button v-if="row.product_url" type="primary" link size="small" @click="openUrl(row.product_url)">
                      打开
                    </el-button>
                    <span v-else>-</span>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-card>

          <!-- 风险评估 -->
          <el-card class="report-card" shadow="never">
            <template #header>风险评估</template>
            <el-descriptions :column="2" border size="small">
              <el-descriptions-item label="品牌侵权风险">
                {{ reportData.risk_brand_infringement || '-' }}
              </el-descriptions-item>
              <el-descriptions-item label="物流难度">
                {{ reportData.risk_logistics_difficulty || '-' }}
              </el-descriptions-item>
              <el-descriptions-item label="是否易碎">
                {{ boolText(reportData.risk_fragile) }}
              </el-descriptions-item>
              <el-descriptions-item label="是否需认证">
                {{ boolText(reportData.risk_certification_needed) }}
              </el-descriptions-item>
              <el-descriptions-item label="是否季节性" :span="2">
                {{ boolText(reportData.risk_seasonal) }}
              </el-descriptions-item>
            </el-descriptions>
          </el-card>

          <!-- AI 综合评估 -->
          <el-card class="report-card" shadow="never">
            <template #header>AI 综合评估</template>
            <el-descriptions :column="2" border size="small" title="各维度评分">
              <el-descriptions-item label="市场需求">
                {{ formatNum(reportData.score_market_demand) }}
              </el-descriptions-item>
              <el-descriptions-item label="竞争度">
                {{ formatNum(reportData.score_competition) }}
              </el-descriptions-item>
              <el-descriptions-item label="利润">
                {{ formatNum(reportData.score_profit) }}
              </el-descriptions-item>
              <el-descriptions-item label="供应链">
                {{ formatNum(reportData.score_supply_chain) }}
              </el-descriptions-item>
              <el-descriptions-item label="风险">
                {{ formatNum(reportData.score_risk) }}
              </el-descriptions-item>
              <el-descriptions-item label="运营">
                {{ formatNum(reportData.score_operation) }}
              </el-descriptions-item>
              <el-descriptions-item v-if="reportData.critic_quality_score != null" label="Critic 质量分" :span="2">
                {{ formatNum(reportData.critic_quality_score) }}
              </el-descriptions-item>
            </el-descriptions>

            <div class="final-score-row">
              <div class="final-score-label">综合评分</div>
              <div class="final-score-value">{{ formatNum(reportData.final_score) }}</div>
              <el-tag v-if="reportData.is_recommended !== undefined && reportData.is_recommended !== null" size="large" :type="reportData.is_recommended ? 'success' : 'info'">
                {{ reportData.is_recommended ? '推荐' : '不推荐' }}
              </el-tag>
            </div>
            <div class="ai-text">
              <el-text>{{ reportData.ai_recommendation || '暂无建议文本' }}</el-text>
            </div>
          </el-card>
        </template>
      </el-tab-pane>

      <!-- Tab2: 历史报告 -->
      <el-tab-pane label="历史报告" name="history">
        <div class="gva-search-box">
          <el-form :inline="true" :model="historySearch" class="demo-form-inline" @keyup.enter="onHistorySearch">
            <el-form-item label="关键词">
              <el-input
                v-model="historySearch.keyword"
                placeholder="按任务ID或URL搜索"
                clearable
                style="width: 220px"
              />
            </el-form-item>
            <el-form-item label="状态">
              <el-select v-model="historySearch.status" placeholder="全部" clearable style="width: 140px">
                <el-option label="pending" value="pending" />
                <el-option label="running" value="running" />
                <el-option label="done" value="done" />
                <el-option label="failed" value="failed" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="search" @click="onHistorySearch">查询</el-button>
              <el-button icon="refresh" @click="onHistoryReset">重置</el-button>
            </el-form-item>
          </el-form>
        </div>

        <div class="gva-table-box">
          <el-table :data="historyTable" style="width: 100%" tooltip-effect="dark" v-loading="historyLoading">
            <el-table-column align="center" label="主图" width="80">
              <template #default="{ row }">
                <el-image
                  v-if="row.main_image_url"
                  :src="row.main_image_url"
                  :preview-src-list="[row.main_image_url]"
                  preview-teleported
                  fit="contain"
                  style="width: 50px; height: 50px"
                />
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column prop="task_id" label="任务ID" width="200" show-overflow-tooltip />
            <el-table-column prop="input_url" label="商品URL" min-width="250" show-overflow-tooltip />
            <el-table-column label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getStatusTagType(row.status)" size="small">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="利润EUR" width="100" align="right">
              <template #default="{ row }">{{ formatNum(row.profit_eur) }}</template>
            </el-table-column>
            <el-table-column label="利润率" width="100" align="right">
              <template #default="{ row }">{{ formatPercentRatio(row.profit_margin) }}</template>
            </el-table-column>
            <el-table-column prop="final_score" label="评分" width="80" align="center">
              <template #default="{ row }">{{ formatNum(row.final_score) }}</template>
            </el-table-column>
            <el-table-column label="推荐" width="80" align="center">
              <template #default="{ row }">
                {{ row.is_recommended === true ? '是' : row.is_recommended === false ? '否' : '-' }}
              </template>
            </el-table-column>
            <el-table-column label="创建时间" width="180" align="center">
              <template #default="{ row }">{{ formatDateTime(row.ctime) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="160" align="center" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link size="small" @click="handleViewReport(row)">查看报告</el-button>
                <el-button
                  v-if="row.status === 'pending' || row.status === 'failed'"
                  type="warning"
                  link
                  size="small"
                  :loading="retryingTaskId === row.task_id"
                  @click="handleRetryTask(row)"
                >
                  重试
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="gva-pagination">
            <el-pagination
              :current-page="historyPage"
              :page-size="historyPageSize"
              :page-sizes="[10, 25, 50, 100]"
              :total="historyTotal"
              layout="total, sizes, prev, pager, next, jumper"
              @current-change="handleHistoryPageChange"
              @size-change="handleHistorySizeChange"
            />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { submitAnalysis, getTaskStatus, getTaskReport, getTaskList, retryTask } from '@/api/selection'

const route = useRoute()

defineOptions({
  name: 'EmagSelection',
})

const activeTab = ref('analyze')
const pyBackendError = ref('')

const analyzeForm = reactive({
  url: '',
})
const analyzeSubmitting = ref(false)

const currentTaskId = ref('')
const taskStatus = reactive({
  status: '',
  error_msg: null,
  failure_detail: null,
  critic_loops: 0,
})

const reportData = ref(null)

let statusPollTimer = null

const historySearch = reactive({
  keyword: '',
  status: '',
})
const historyPage = ref(1)
const historyPageSize = ref(10)
const historyTotal = ref(0)
const historyTable = ref([])
const historyLoading = ref(false)

/** 清理轮询 */
const clearStatusPoll = () => {
  if (statusPollTimer) {
    clearInterval(statusPollTimer)
    statusPollTimer = null
  }
}

/** 状态 -> el-tag 类型 */
const getStatusTagType = (status) => {
  const map = {
    pending: 'warning',
    running: 'primary',
    done: 'success',
    failed: 'danger',
  }
  return map[status] || 'info'
}

/** 状态 -> 进度条百分比（无状态视为刚提交，按 pending 展示） */
const getProgressPercent = (status) => {
  if (!status || status === 'pending') return 10
  if (status === 'running') return 50
  if (status === 'done') return 100
  return 10
}

const formatNum = (val) => {
  if (val === null || val === undefined || val === '') return '-'
  const n = Number(val)
  if (Number.isNaN(n)) return String(val)
  return n.toFixed(2)
}

/** 0-1 小数转百分比展示 */
const formatPercentRatio = (ratio) => {
  if (ratio === null || ratio === undefined || ratio === '') return '-'
  const n = Number(ratio)
  if (Number.isNaN(n)) return '-'
  return `${(n * 100).toFixed(2)}%`
}

/** OEM 等已是 0-1 比例时同上 */
const formatDateTime = (str) => {
  if (!str) return '-'
  try {
    const d = new Date(str)
    if (!Number.isNaN(d.getTime())) return d.toLocaleString()
  } catch {
    /* ignore */
  }
  return str.length > 19 ? str.substring(0, 19).replace('T', ' ') : str
}

const boolText = (v) => {
  if (v === true) return '是'
  if (v === false) return '否'
  return '-'
}

const costKeyLabels = {
  product_cost_cny: '商品成本 (CNY)',
  logistics_cny: '物流 (CNY)',
  total_cost_eur: '总成本 (EUR)',
}

const formatCostKey = (key) => costKeyLabels[key] || key

const openUrl = (url) => {
  if (url) window.open(url, '_blank')
}

const setPyError = (err) => {
  const msg = err?.response?.data?.detail || err?.message || 'Python 选品服务不可用或请求失败'
  pyBackendError.value = String(msg)
}

const clearPyError = () => {
  pyBackendError.value = ''
}

/** 拉取任务状态并更新 UI */
const pollTaskStatusOnce = async () => {
  if (!currentTaskId.value) return
  try {
    clearPyError()
    const res = await getTaskStatus(currentTaskId.value)
    taskStatus.status = res.status || ''
    taskStatus.error_msg = res.error_msg
    taskStatus.failure_detail = res.failure_detail
    taskStatus.critic_loops = res.critic_loops ?? 0

    if (res.status === 'done') {
      clearStatusPoll()
      await fetchReport()
    }
    if (res.status === 'failed') {
      clearStatusPoll()
    }
  } catch (err) {
    setPyError(err)
    clearStatusPoll()
  }
}

/** 拉取分析报告 */
const fetchReport = async () => {
  if (!currentTaskId.value) return
  try {
    clearPyError()
    const res = await getTaskReport(currentTaskId.value)
    reportData.value = res.report || null
    if (!reportData.value) {
      ElMessage.warning('暂无报告内容')
    }
  } catch (err) {
    setPyError(err)
    reportData.value = null
  }
}

/** 开始轮询（3s） */
const startStatusPoll = () => {
  clearStatusPoll()
  statusPollTimer = setInterval(() => {
    pollTaskStatusOnce()
  }, 3000)
}

/** 发起分析 */
const handleStartAnalyze = async () => {
  const url = (analyzeForm.url || '').trim()
  if (!url) {
    ElMessage.warning('请输入商品链接')
    return
  }
  analyzeSubmitting.value = true
  clearStatusPoll()
  reportData.value = null
  currentTaskId.value = ''
  taskStatus.status = ''
  taskStatus.error_msg = null
  taskStatus.failure_detail = null
  taskStatus.critic_loops = 0
  try {
    clearPyError()
    const res = await submitAnalysis({ url })
    currentTaskId.value = res.task_id || ''
    taskStatus.status = res.status || 'pending'
    if (!currentTaskId.value) {
      ElMessage.error('未返回任务ID')
      return
    }
    ElMessage.success('任务已提交')
    await pollTaskStatusOnce()
    if (taskStatus.status !== 'done' && taskStatus.status !== 'failed') {
      startStatusPoll()
    }
  } catch (err) {
    setPyError(err)
  } finally {
    analyzeSubmitting.value = false
  }
}

/** 历史列表 */
const fetchHistoryList = async () => {
  historyLoading.value = true
  try {
    clearPyError()
    const offset = (historyPage.value - 1) * historyPageSize.value
    const params = {
      limit: historyPageSize.value,
      offset,
    }
    if (historySearch.keyword) params.keyword = historySearch.keyword
    if (historySearch.status) params.status = historySearch.status
    const res = await getTaskList(params)
    historyTable.value = res.items || []
    historyTotal.value = res.total ?? 0
  } catch (err) {
    setPyError(err)
    historyTable.value = []
    historyTotal.value = 0
  } finally {
    historyLoading.value = false
  }
}

const onHistorySearch = () => {
  historyPage.value = 1
  fetchHistoryList()
}

const onHistoryReset = () => {
  historySearch.keyword = ''
  historySearch.status = ''
  historyPage.value = 1
  fetchHistoryList()
}

const handleHistorySizeChange = (val) => {
  historyPageSize.value = val
  fetchHistoryList()
}

const handleHistoryPageChange = (val) => {
  historyPage.value = val
  fetchHistoryList()
}

const retryingTaskId = ref('')

/** 重试 pending/failed 任务 */
const handleRetryTask = async (row) => {
  if (!row?.task_id) return
  retryingTaskId.value = row.task_id
  try {
    await retryTask(row.task_id)
    ElMessage.success('已重新提交分析')
    await fetchHistoryList()
  } catch (err) {
    ElMessage.error(err?.response?.data?.detail || '重试失败')
  } finally {
    retryingTaskId.value = ''
  }
}

/** 从历史查看报告：切到分析 Tab 并加载，进行中的任务自动轮询 */
const handleViewReport = async (row) => {
  if (!row?.task_id) return
  clearStatusPoll()
  activeTab.value = 'analyze'
  currentTaskId.value = row.task_id
  taskStatus.status = row.status || ''
  taskStatus.error_msg = null
  taskStatus.failure_detail = null
  reportData.value = null
  try {
    clearPyError()
    const st = await getTaskStatus(row.task_id)
    taskStatus.status = st.status || row.status
    taskStatus.error_msg = st.error_msg
    taskStatus.failure_detail = st.failure_detail
    taskStatus.critic_loops = st.critic_loops ?? 0
  } catch (err) {
    setPyError(err)
  }

  if (taskStatus.status === 'done') {
    await fetchReport()
  } else if (taskStatus.status === 'pending' || taskStatus.status === 'running') {
    startStatusPoll()
  }
}

onMounted(() => {
  fetchHistoryList()

  const queryTaskId = route.query.task_id
  if (queryTaskId) {
    handleViewReport({ task_id: queryTaskId, status: '' })
  }
})

onUnmounted(() => {
  clearStatusPoll()
})
</script>

<style scoped lang="scss">
.selection-container {
  padding: 16px;
  min-height: calc(100vh - 100px);
}

.backend-alert {
  margin-bottom: 12px;
}

.selection-tabs {
  background: transparent;
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

.progress-section {
  margin-bottom: 16px;
}

.progress-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  flex-wrap: wrap;

  .label {
    color: var(--el-text-color-secondary);
    font-size: 13px;
  }

  .task-id {
    max-width: 360px;
    font-family: monospace;
  }
}

.fail-msg {
  margin-top: 8px;
}

.fail-detail {
  margin-top: 6px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.meta-line {
  margin-top: 8px;
}

.report-card {
  margin-bottom: 16px;
}

.sub-table-wrap {
  margin-top: 12px;
}

.sub-table-title {
  font-size: 13px;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--el-text-color-primary);
}

.final-score-row {
  display: flex;
  align-items: center;
  gap: 16px;
  margin: 16px 0 12px;
  flex-wrap: wrap;
}

.final-score-label {
  font-size: 14px;
  color: var(--el-text-color-secondary);
}

.final-score-value {
  font-size: 36px;
  font-weight: 700;
  color: var(--el-color-primary);
  line-height: 1;
}

.ai-text {
  margin-top: 8px;
  line-height: 1.6;
}
</style>
