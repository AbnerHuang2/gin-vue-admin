<template>
  <div class="order-container">
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchForm" class="demo-form-inline" @keyup.enter="onSearch">
        <el-form-item label="订单ID">
          <el-input v-model="searchForm.orderId" placeholder="请输入订单ID" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item label="国家">
          <el-select v-model="searchForm.country" placeholder="请选择国家" clearable style="width: 120px">
            <el-option
              v-for="country in countryList"
              :key="country"
              :label="getCountryName(country)"
              :value="country"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="订单日期">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            style="width: 260px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button type="warning" icon="refresh" :loading="syncLoading" @click="handleSync">
            {{ syncLoading ? '同步中...' : '同步订单' }}
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <el-table
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="orderId"
        v-loading="loading"
      >
        <el-table-column align="center" label="序号" type="index" width="60" :index="indexMethod" />
        <el-table-column align="left" label="订单ID" prop="orderId" width="120" />
        <el-table-column align="left" label="订单日期" prop="orderDateLocal" width="160" />
        <el-table-column align="center" label="国家" width="80">
          <template #default="scope">
            <el-tag :type="getCountryTagType(scope.row.country)" size="small">
              {{ getCountryName(scope.row.country) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="center" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)" size="small">
              {{ getStatusName(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="right" label="订单金额" width="110">
          <template #default="scope">
            <span>{{ formatPrice(scope.row.price) }} {{ scope.row.currency }}</span>
          </template>
        </el-table-column>
        <el-table-column align="right" label="订单金额(CNY)" width="100">
          <template #default="scope">
            <span class="price-cny">¥{{ formatPrice(scope.row.priceCny) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="产品信息" min-width="380">
          <template #default="scope">
            <div v-if="scope.row.products && scope.row.products.length > 0" class="product-list">
              <div 
                v-for="(product, idx) in scope.row.products" 
                :key="product.productId" 
                class="product-item"
                :class="{ 'product-item-border': idx > 0 }"
              >
                <div class="product-info">
                  <span class="product-name" :title="product.productName">{{ product.productName || product.productId }}</span>
                  <span class="product-detail">
                    数量: <b>{{ product.quantity }}</b> | 
                    单价: {{ formatPrice(product.salePrice) }} | 
                    单价(CNY): <span class="price-cny">¥{{ formatPrice(product.salePriceCny) }}</span>
                  </span>
                </div>
                <el-button 
                  v-if="product.productUrl" 
                  type="primary" 
                  link 
                  size="small" 
                  @click="openProductUrl(product.productUrl)"
                >
                  查看
                </el-button>
              </div>
            </div>
            <span v-else class="no-product">无产品</span>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
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
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { getOrderList, syncOrders, getCountryList } from '@/api/emagOrder'

defineOptions({
  name: 'EmagOrder'
})

// 搜索表单
const searchForm = reactive({
  orderId: '',
  country: '',
  status: ''
})

// 日期范围
const dateRange = ref([])

// 分页
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 数据
const tableData = ref([])
const countryList = ref([])
const loading = ref(false)
const syncLoading = ref(false)

// 国家名称映射
const countryNameMap = {
  'RO': '罗马尼亚',
  'BG': '保加利亚',
  'HU': '匈牙利',
  'PL': '波兰'
}

// 国家标签颜色
const countryTagTypeMap = {
  'RO': 'success',
  'BG': 'info',
  'HU': 'warning',
  'PL': 'info'
}

// 状态名称映射
const statusNameMap = {
  'STATUS_FINALIZED': '已完成',
  'STATUS_NEW': '新订单',
  'STATUS_IN_PROGRESS': '处理中',
  'STATUS_CANCELLED': '已取消'
}

// 状态标签颜色
const statusTagTypeMap = {
  'STATUS_FINALIZED': 'success',
  'STATUS_NEW': 'warning',
  'STATUS_IN_PROGRESS': '',
  'STATUS_CANCELLED': 'danger'
}

// 获取国家名称
const getCountryName = (code) => {
  return countryNameMap[code] || code
}

// 获取国家标签类型
const getCountryTagType = (code) => {
  return countryTagTypeMap[code] || 'info'
}

// 获取状态名称
const getStatusName = (status) => {
  return statusNameMap[status] || status
}

// 获取状态标签类型
const getStatusTagType = (status) => {
  return statusTagTypeMap[status] || 'info'
}

// 格式化价格
const formatPrice = (price) => {
  if (price === null || price === undefined) return '-'
  return Number(price).toFixed(2)
}

// 序号方法
const indexMethod = (index) => {
  return (page.value - 1) * pageSize.value + index + 1
}

// 获取国家列表
const fetchCountryList = async () => {
  try {
    const res = await getCountryList()
    if (res.code === 0) {
      countryList.value = res.data || []
    }
  } catch (error) {
    console.error('获取国家列表失败', error)
  }
}

// 获取订单列表
const fetchOrderList = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value,
      ...searchForm
    }
    // 处理日期范围
    if (dateRange.value && dateRange.value.length === 2) {
      params.startDate = dateRange.value[0]
      params.endDate = dateRange.value[1]
    }

    const res = await getOrderList(params)
    if (res.code === 0) {
      tableData.value = res.data.list || []
      total.value = res.data.total
    } else {
      ElMessage.error(res.msg || '获取订单列表失败')
    }
  } catch (error) {
    ElMessage.error('获取订单列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const onSearch = () => {
  page.value = 1
  fetchOrderList()
}

// 重置
const onReset = () => {
  searchForm.orderId = ''
  searchForm.country = ''
  searchForm.status = ''
  dateRange.value = []
  page.value = 1
  fetchOrderList()
}

// 同步订单
const handleSync = async () => {
  try {
    syncLoading.value = true
    const res = await syncOrders({
      status: 'STATUS_FINALIZED',
      page: 0,
      limit: 100
    })
    if (res.code === 0) {
      ElMessage.success(res.msg || '同步任务已启动')
      // 3秒后刷新列表
      setTimeout(() => {
        fetchOrderList()
      }, 3000)
    } else {
      ElMessage.error(res.msg || '同步失败')
    }
  } catch (error) {
    ElMessage.error('同步失败')
  } finally {
    setTimeout(() => {
      syncLoading.value = false
    }, 3000)
  }
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  fetchOrderList()
}

const handleCurrentChange = (val) => {
  page.value = val
  fetchOrderList()
}

// 打开产品链接
const openProductUrl = (url) => {
  window.open(url, '_blank')
}

// 初始化
onMounted(() => {
  fetchCountryList()
  fetchOrderList()
})
</script>

<style scoped lang="scss">
.order-container {
  padding: 16px;
  min-height: calc(100vh - 100px);
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

.price-cny {
  color: #F56C6C;
  font-weight: 600;
}

.product-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.product-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 0;
}

.product-item-border {
  border-top: 1px dashed #e4e7ed;
  padding-top: 8px;
}

.product-info {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-width: 0;
}

.product-name {
  font-weight: 500;
  // color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 280px;
}

.product-detail {
  font-size: 12px;
  //color: #909399;
  margin-top: 2px;
}

.no-product {
  color: #909399;
  font-style: italic;
}
</style>

