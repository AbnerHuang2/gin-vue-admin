<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="品类ID">
          <el-input v-model="searchInfo.categoryId" placeholder="请输入品类ID" clearable />
        </el-form-item>
        <el-form-item label="品类名称">
          <el-input v-model="searchInfo.categoryName" placeholder="请输入品类名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog('add')">新增</el-button>
        <el-button
          :disabled="!multipleSelection.length"
          icon="delete"
          @click="onDelete"
        >删除</el-button>
      </div>
      <el-table
        ref="multipleTableRef"
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="id"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="ID" prop="id" width="80" />
        <el-table-column align="left" label="品类ID" prop="categoryId" width="180" />
        <el-table-column align="left" label="品类名称" prop="categoryName" min-width="150" />
        <el-table-column align="left" label="子品类名称" prop="subcategoryName" min-width="150" />
        <el-table-column align="left" label="创建时间" width="180">
          <template #default="scope">
            <span>{{ formatDate(scope.row.createdAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="更新时间" width="180">
          <template #default="scope">
            <span>{{ formatDate(scope.row.updatedAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" width="200">
          <template #default="scope">
            <el-button type="primary" link icon="edit" @click="openDialog('edit', scope.row)">编辑</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogFormVisible"
      :title="dialogTitle"
      width="500px"
    >
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="品类ID" prop="categoryId">
          <el-input v-model="formData.categoryId" placeholder="请输入品类ID" :disabled="dialogType === 'edit'" />
        </el-form-item>
        <el-form-item label="品类名称" prop="categoryName">
          <el-input v-model="formData.categoryName" placeholder="请输入品类名称" />
        </el-form-item>
        <el-form-item label="子品类名称" prop="subcategoryName">
          <el-input v-model="formData.subcategoryName" placeholder="请输入子品类名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取消</el-button>
          <el-button type="primary" @click="submitForm">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createEmagCategory,
  deleteEmagCategory,
  deleteEmagCategoryByIds,
  updateEmagCategory,
  getEmagCategoryList
} from '@/api/emagCategory'
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { formatDate } from '@/utils/format'

defineOptions({
  name: 'EmagCategory'
})

// 搜索条件
const searchInfo = reactive({
  categoryId: '',
  categoryName: ''
})

// 分页
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 表格数据
const tableData = ref([])
const multipleSelection = ref([])
const multipleTableRef = ref(null)

// 对话框
const dialogFormVisible = ref(false)
const dialogType = ref('add')
const dialogTitle = ref('新增品类')
const formRef = ref(null)
const formData = reactive({
  categoryId: '',
  categoryName: '',
  subcategoryName: ''
})

// 表单验证规则
const rules = {
  categoryId: [
    { required: true, message: '请输入品类ID', trigger: 'blur' }
  ],
  categoryName: [
    { required: true, message: '请输入品类名称', trigger: 'blur' }
  ]
}

// 获取表格数据
const getTableData = async () => {
  const res = await getEmagCategoryList({
    page: page.value,
    pageSize: pageSize.value,
    ...searchInfo
  })
  if (res.code === 0) {
    tableData.value = res.data.list || []
    total.value = res.data.total
  }
}

// 初始化加载
getTableData()

// 分页事件
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 多选
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 搜索
const onSubmit = () => {
  page.value = 1
  getTableData()
}

// 重置
const onReset = () => {
  searchInfo.categoryId = ''
  searchInfo.categoryName = ''
  page.value = 1
  getTableData()
}

// 打开对话框
const openDialog = (type, row = null) => {
  dialogType.value = type
  dialogTitle.value = type === 'add' ? '新增品类' : '编辑品类'
  if (type === 'edit' && row) {
    formData.categoryId = row.categoryId
    formData.categoryName = row.categoryName
    formData.subcategoryName = row.subcategoryName
  } else {
    formData.categoryId = ''
    formData.categoryName = ''
    formData.subcategoryName = ''
  }
  dialogFormVisible.value = true
}

// 关闭对话框
const closeDialog = () => {
  dialogFormVisible.value = false
  formRef.value?.resetFields()
}

// 提交表单
const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      let res
      if (dialogType.value === 'add') {
        res = await createEmagCategory(formData)
      } else {
        res = await updateEmagCategory(formData)
      }
      if (res.code === 0) {
        ElMessage.success(dialogType.value === 'add' ? '新增成功' : '更新成功')
        closeDialog()
        getTableData()
      }
    }
  })
}

// 删除单条
const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除该品类吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await deleteEmagCategory({ categoryId: row.categoryId })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getTableData()
    }
  })
}

// 批量删除
const onDelete = () => {
  ElMessageBox.confirm('确定要删除选中的品类吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const ids = multipleSelection.value.map(item => item.id)
    const res = await deleteEmagCategoryByIds({ ids })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      getTableData()
    }
  })
}
</script>

<style scoped>
</style>

