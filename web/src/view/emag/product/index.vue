<template>
  <div class="product-container">
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchForm" class="demo-form-inline" @keyup.enter="onSearch">
        <el-form-item label="产品编号Key">
          <el-input v-model="searchForm.pnk" placeholder="请输入产品编号Key" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item label="外部ID">
          <el-input v-model="searchForm.extId" placeholder="请输入外部ID" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item label="最小库存">
          <el-input-number v-model="searchForm.stockMin" :min="0" placeholder="最小库存" style="width: 140px" />
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
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable style="width: 140px">
            <el-option
              v-for="status in statusList"
              :key="status"
              :label="getStatusName(status)"
              :value="status"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
          <el-button type="warning" icon="refresh" :loading="syncLoading" @click="handleSync">
            {{ syncLoading ? '同步中...' : '同步产品' }}
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
        <el-table-column align="left" label="外部ID" prop="extId" width="150" show-overflow-tooltip />
        <el-table-column align="left" label="产品编号Key" prop="pnk" width="120" show-overflow-tooltip />
        <el-table-column align="center" label="国家" width="80">
          <template #default="scope">
            <el-tag :type="getCountryTagType(scope.row.country)" size="small">
              {{ getCountryName(scope.row.country) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="产品标题" prop="title" min-width="280" show-overflow-tooltip />
        <el-table-column align="center" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)" size="small">
              {{ getStatusName(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="right" label="销售价格" width="120">
          <template #default="scope">
            <span>{{ formatPrice(scope.row.salePrice) }} {{ scope.row.currency }}</span>
          </template>
        </el-table-column>
        <el-table-column align="right" label="含税价格" width="120">
          <template #default="scope">
            <span>{{ formatPrice(scope.row.afterTaxPrice) }} {{ scope.row.currency }}</span>
          </template>
        </el-table-column>
        <el-table-column align="right" label="销售价格(CNY)" width="120">
          <template #default="scope">
            <span class="price-cny">¥{{ formatPrice(scope.row.salePriceCn) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="right" label="成本价格(CNY)" width="120">
          <template #default="scope">
            <span v-if="scope.row.costPriceCn > 0" class="cost-price">¥{{ formatPrice(scope.row.costPriceCn) }}</span>
            <span v-else class="no-data">-</span>
          </template>
        </el-table-column>
        <el-table-column align="center" label="增值税率" prop="vat" width="100" />
        <el-table-column align="center" label="库存" prop="stock" width="80">
          <template #default="scope">
            <el-tag v-if="scope.row.stock <= 0" type="danger" size="small">{{ scope.row.stock }}</el-tag>
            <el-tag v-else-if="scope.row.stock < 10" type="warning" size="small">{{ scope.row.stock }}</el-tag>
            <span v-else>{{ scope.row.stock }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" label="购买按钮排名" width="120">
          <template #default="scope">
            <el-tag v-if="scope.row.buyButtonRank === 1" type="success" size="small">
              第{{ scope.row.buyButtonRank }}名
            </el-tag>
            <el-tag v-else-if="scope.row.buyButtonRank <= 3" type="warning" size="small">
              第{{ scope.row.buyButtonRank }}名
            </el-tag>
            <span v-else>第{{ scope.row.buyButtonRank }}名</span>
          </template>
        </el-table-column>
        <el-table-column align="center" label="竞争对手数" width="100">
          <template #default="scope">
            <span>{{ scope.row.buyButtonCnt }}</span>
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" width="100" fixed="right">
          <template #default="scope">
            <el-button 
              v-if="scope.row.url" 
              type="primary" 
              link 
              size="small" 
              @click="openProductUrl(scope.row.url)"
            >
              查看详情
            </el-button>
            <span v-else class="no-data">-</span>
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
import { ElMessage } from 'element-plus'
import { getProductList, syncProducts, getStatusList, getCountryList } from '@/api/emagProduct'

defineOptions({
  name: 'EmagProduct'
})

// 搜索表单
const searchForm = reactive({
  pnk: '',
  extId: '',
  stockMin: 0,
  country: '',
  status: ''
})

// 分页
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 数据
const tableData = ref([])
const statusList = ref([])
const countryList = ref([])
const loading = ref(false)
const syncLoading = ref(false)

// 国家名称映射
const countryNameMap = {
  'RO': '罗马尼亚',
  'BG': '保加利亚',
  'HU': '匈牙利'
}

// 国家标签颜色
const countryTagTypeMap = {
  'RO': 'success',
  'BG': 'info',
  'HU': 'warning'
}

// 状态名称映射
const statusNameMap = {
  'Active': '活跃',
  'Inactive': '未激活',
  'Blocked': '已屏蔽',
  'Auto Inactivated': '自动未激活'
}

// 状态标签颜色
const statusTagTypeMap = {
  'Active': 'success',
  'Inactive': 'info',
  'Blocked': 'danger',
  'Auto Inactivated': 'warning'
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

// 获取状态列表
const fetchStatusList = async () => {
  try {
    const res = await getStatusList()
    if (res.code === 0) {
      statusList.value = res.data || []
    }
  } catch (error) {
    console.error('获取状态列表失败', error)
  }
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

// 获取产品列表
const fetchProductList = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      pageSize: pageSize.value,
      ...searchForm
    }

    const res = await getProductList(params)
    if (res.code === 0) {
      tableData.value = res.data.list || []
      total.value = res.data.total
    } else {
      ElMessage.error(res.msg || '获取产品列表失败')
    }
  } catch (error) {
    ElMessage.error('获取产品列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const onSearch = () => {
  page.value = 1
  fetchProductList()
}

// 重置
const onReset = () => {
  searchForm.pnk = ''
  searchForm.extId = ''
  searchForm.stockMin = 0
  searchForm.country = ''
  searchForm.status = ''
  page.value = 1
  fetchProductList()
}

// 同步产品
const handleSync = async () => {
  try {
    syncLoading.value = true
    const res = await syncProducts({
      page: 0,
      limit: 100
    })
    if (res.code === 0) {
      ElMessage.success(res.msg || '同步任务已启动')
      // 3秒后刷新列表
      setTimeout(() => {
        fetchProductList()
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
  fetchProductList()
}

const handleCurrentChange = (val) => {
  page.value = val
  fetchProductList()
}

// 打开产品链接
const openProductUrl = (url) => {
  window.open(url, '_blank')
}

// 初始化
onMounted(() => {
  fetchStatusList()
  fetchCountryList()
  fetchProductList()
})
</script>

<style scoped lang="scss">
.product-container {
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

.cost-price {
  color: #E6A23C;
  font-weight: 600;
}

.no-data {
  color: #909399;
}
</style>
