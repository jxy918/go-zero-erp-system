<template>
  <div class="product-unit">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>
            {{ currentProduct.id ? `${currentProduct.name} - 计量单位管理` : '产品计量单位管理' }}
          </span>
          <div>
            <el-button v-if="currentProduct.id" type="info" @click="handleClearFilter">查看全部产品</el-button>
            <el-button type="primary" @click="handleAdd" v-has-permission="'product_unit:create'">新增单位</el-button>
          </div>
        </div>
      </template>
      <div class="description-box">
        <span class="description-icon">ℹ️</span>
        <span class="description-text">产品计量单位管理：为产品配置多计量单位，支持主单位和辅助单位，设置换算比例实现自动单位换算。例如：1箱=24瓶。</span>
      </div>
      <div v-if="currentProduct.id" class="current-product-info">
        <span style="color: #1890ff; font-weight: bold;">📌 当前管理：{{ currentProduct.name }}</span>
      </div>
      <div v-else class="query-bar">
        <el-input
          v-model="searchQuery"
          placeholder="请输入产品名称"
          style="width: 300px"
          prefix-icon="el-icon-search"
        ></el-input>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
      </div>
      <el-table :data="unitList" border stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="60"></el-table-column>
        <el-table-column prop="product_name" label="产品名称" min-width="150"></el-table-column>
        <el-table-column prop="unit_name" label="单位名称" width="100"></el-table-column>
        <el-table-column prop="ratio" label="换算比例" width="120">
          <template #default="{ row }">
            1{{ row.unit_name }} = {{ row.ratio }}主单位
          </template>
        </el-table-column>
        <el-table-column prop="is_main" label="是否主单位" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_main === 1 ? 'success' : ''">
              {{ row.is_main === 1 ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160"></el-table-column>
        <el-table-column prop="updated_at" label="更新时间" width="160"></el-table-column>
        <el-table-column label="操作" width="180" align="center">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination-bar">
        <el-pagination
          :current-page="pagination.page"
          :page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        ></el-pagination>
      </div>
    </el-card>

    <el-dialog title="新增/编辑计量单位" :model-value="dialogVisible" @update:model-value="dialogVisible = $event" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="产品" prop="product_id">
          <el-select v-model="form.product_id" placeholder="请选择产品" :disabled="isEdit || currentProduct.id">
            <el-option v-for="product in products" :key="product.id" :label="product.name || product.product_name" :value="product.id"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="单位名称" prop="unit_name">
          <el-input v-model="form.unit_name" placeholder="请输入单位名称"></el-input>
        </el-form-item>
        <el-form-item label="换算比例" prop="ratio">
          <el-input-number v-model="form.ratio" :min="0.0001" :step="0.0001" placeholder="与主单位的换算比例"></el-input-number>
          <span style="margin-left: 10px; color: #999">例如：1箱=24瓶，则填24</span>
        </el-form-item>
        <el-form-item label="是否主单位">
          <el-switch v-model="form.is_main" :active-value="1" :inactive-value="0"></el-switch>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { productApi } from '@/api/product'
const route = useRoute();
const loading = ref(false);
const unitList = ref([]);
const searchQuery = ref('');
const currentProduct = ref({
  id: null,
  name: ''
});
const pagination = ref({
 page: 1,
 pageSize: 10,
 total: 0
});
const dialogVisible = ref(false);
const isEdit = ref(false);
const products = ref([]);
const form = reactive({
 id: 0,
 product_id: null, // 默认为 null，让 el-select 显示 placeholder
 unit_name: '',
 ratio: 1,
 is_main: 0
});
const rules = {
 product_id: [
 { required: true, message: '请选择产品', trigger: 'change' }
 ],
 unit_name: [
 { required: true, message: '请输入单位名称', trigger: 'blur' }
 ],
 ratio: [
 { required: true, message: '请输入换算比例', trigger: 'blur' },
 { validator: (rule, value, callback) => {
 if (value <= 0) {
 callback(new Error('换算比例必须大于0'));
 }
 else {
 callback();
 }
 }}
 ]
};
const formRef = ref();
const fetchUnitList = async () => {
 loading.value = true;
 try {
 const params = {
 page: pagination.value.page,
 page_size: pagination.value.pageSize
 };
 if (currentProduct.value.id) {
 params.product_id = currentProduct.value.id;
 } else if (searchQuery.value) {
 params.product_name = searchQuery.value;
 }
 const res = await productApi.listProductUnit(params);
 if (res.code === 0) {
 unitList.value = res.data.units;
 pagination.value.total = res.data.total;
 }
 }
 catch (error) {
 console.error('获取计量单位列表失败:', error);
 }
 finally {
 loading.value = false;
 }
};
const fetchProducts = async () => {
 try {
 const res = await productApi.getProductList({ page: 1, page_size: 100 });
 if (res.code === 0) {
 products.value = res.data.products;
 }
 }
 catch (error) {
 console.error('获取产品列表失败:', error);
 }
};
const handleSearch = () => {
 pagination.value.page = 1;
 fetchUnitList();
};
const handleSizeChange = (val) => {
 pagination.value.pageSize = val;
 fetchUnitList();
};
const handleCurrentChange = (val) => {
 pagination.value.page = val;
 fetchUnitList();
};
const handleAdd = async () => {
 isEdit.value = false;
 form.id = 0;
 form.product_id = currentProduct.value.id ? Number(currentProduct.value.id) : null;
 form.unit_name = '';
 form.ratio = 1;
 form.is_main = 0;
 // 如果有指定当前产品，确保产品列表中有这个产品
 if (currentProduct.value.id) {
 const exists = products.value.some(p => p.id === currentProduct.value.id);
 if (!exists) {
 products.value.unshift({
 id: currentProduct.value.id,
 name: currentProduct.value.name
 });
 }
 } else if (products.value.length === 0) {
 // 如果没有指定当前产品且产品列表为空，才获取产品列表
 await fetchProducts();
 }
 dialogVisible.value = true;
};

const handleClearFilter = () => {
 currentProduct.value = { id: null, name: '' };
 searchQuery.value = '';
 pagination.value.page = 1;
 fetchUnitList();
};
const handleEdit = async (row) => {
 isEdit.value = true;
 form.id = row.id;
 form.product_id = Number(row.product_id); // 确保是 Number 类型
 form.unit_name = row.unit_name;
 form.ratio = row.ratio;
 form.is_main = row.is_main;
 // 检查当前产品是否在列表中
 const exists = products.value.some(p => p.id === row.product_id);
 if (!exists) {
 if (currentProduct.value.id && currentProduct.value.id === row.product_id) {
 // 如果是从产品管理跳转过来的，直接使用当前产品信息
 products.value.unshift({
 id: currentProduct.value.id,
 name: currentProduct.value.name
 });
 } else if (!currentProduct.value.id) {
 // 如果没有指定当前产品，获取产品详情
 try {
 const res = await productApi.getProduct(row.product_id);
 if (res.code === 0 && res.data) {
 products.value.unshift(res.data);
 }
 } catch (error) {
 console.error('获取产品详情失败:', error);
 }
 }
 }
 dialogVisible.value = true;
};
const handleDelete = async (row) => {
 ElMessageBox.confirm('确定要删除这个计量单位吗？', '提示', {
 type: 'warning'
 }).then(async () => {
 try {
 const res = await productApi.deleteProductUnit({ id: row.id });
 if (res.code === 0) {
 ElMessage.success('删除成功');
 fetchUnitList();
 }
 }
 catch (error) {
 console.error('删除失败:', error);
 }
 }).catch(() => {
 });
};
const handleSubmit = async () => {
 if (!formRef.value)
 return;
 formRef.value.validate(async (valid) => {
 if (valid) {
 try {
 // 将驼峰格式转换为下划线格式，确保数据类型正确
 const submitData = {
 product_id: Number(form.product_id),
 unit_name: form.unit_name,
 ratio: Number(form.ratio),
 is_main: Number(form.is_main)
 };
 if (isEdit.value) {
 submitData.id = Number(form.id);
 }
 let res;
 if (isEdit.value) {
 res = await productApi.updateProductUnit(submitData);
 }
 else {
 res = await productApi.createProductUnit(submitData);
 }
 if (res.code === 0) {
 ElMessage.success(isEdit.value ? '更新成功' : '创建成功');
 dialogVisible.value = false;
 fetchUnitList();
 }
 }
 catch (error) {
 console.error('提交失败:', error);
 }
 }
 });
};
onMounted(() => {
 if (route.query.productId) {
 currentProduct.value.id = Number(route.query.productId);
 currentProduct.value.name = route.query.productName || '';
 }
 fetchUnitList();
 // 只在没有指定产品时才获取完整产品列表
 if (!currentProduct.value.id) {
 fetchProducts();
 }
});
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.description-box {
  display: flex;
  align-items: flex-start;
  padding: 12px 16px;
  background-color: #fff7f0;
  border-left: 4px solid #ff6b35;
  border-radius: 4px;
  margin-bottom: 20px;
}

.description-icon {
  font-size: 16px;
  margin-right: 10px;
  flex-shrink: 0;
}

.description-text {
  font-size: 13px;
  color: #d93026;
  line-height: 1.6;
}

.current-product-info {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  background-color: #e6f7ff;
  border-left: 4px solid #1890ff;
  border-radius: 4px;
  margin-bottom: 20px;
}

.query-bar {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
  flex-wrap: wrap;
  align-items: center;
}

.pagination-bar {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>