<template>
  <div class="inventory-alert">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>库存预警</span>
          <div>
            <el-button type="primary" @click="handleCheckAlert">手动检查</el-button>
          </div>
        </div>
      </template>
      <div class="description-box">
        <span class="description-icon">ℹ️</span>
        <span class="description-text">库存预警：系统自动监测库存数量，当库存低于安全库存、低于最低库存或高于最高库存时触发预警，帮助及时发现库存异常。</span>
      </div>
      <div class="query-bar">
        <el-select v-model="alertTypeFilter" placeholder="预警类型" clearable style="width: 150px">
          <el-option :value="0" label="全部"></el-option>
          <el-option :value="1" label="低于安全库存"></el-option>
          <el-option :value="2" label="低于最低库存"></el-option>
          <el-option :value="3" label="高于最高库存"></el-option>
        </el-select>
        <el-button type="primary" @click="handleSearch">搜索</el-button>
      </div>
      <el-table :data="alertList" border stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="60"></el-table-column>
        <el-table-column prop="product_name" label="产品名称" min-width="150"></el-table-column>
        <el-table-column prop="product_code" label="产品编码" min-width="120"></el-table-column>
        <el-table-column prop="warehouse" label="仓库" width="100"></el-table-column>
        <el-table-column prop="quantity" label="当前库存" width="100">
          <template #default="{ row }">
            <span :class="getQuantityClass(row)">{{ row.quantity }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="min_stock" label="最低库存" width="100"></el-table-column>
        <el-table-column prop="safety_stock" label="安全库存" width="100"></el-table-column>
        <el-table-column prop="max_stock" label="最高库存" width="100"></el-table-column>
        <el-table-column prop="alert_type" label="预警类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getAlertTypeTagType(row.alert_type)">{{ getAlertTypeText(row.alert_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="alert_level" label="预警级别" width="100">
          <template #default="{ row }">
            <el-tag :type="getAlertLevelTagType(row.alert_level)">{{ getAlertLevelText(row.alert_level) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="预警时间" width="160"></el-table-column>
        <el-table-column label="操作" width="120" align="center">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="handleViewProduct(row)">查看产品</el-button>
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { inventoryApi } from '@/api/inventory';

const router = useRouter();
const loading = ref(false);
const alertList = ref([]);
const alertTypeFilter = ref(0);
const pagination = ref({
  page: 1,
  pageSize: 10,
  total: 0
});
const fetchAlertList = async () => {
  loading.value = true;
  try {
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    };
    if (alertTypeFilter.value > 0) {
      params.alert_type = alertTypeFilter.value;
    }
    const res = await inventoryApi.listInventoryAlert(params);
    if (res.code === 0) {
      alertList.value = res.data.alerts;
      pagination.value.total = res.data.total;
    }
  }
  catch (error) {
    console.error('获取库存预警列表失败:', error);
  }
  finally {
    loading.value = false;
  }
};
const handleSearch = () => {
  pagination.value.page = 1;
  fetchAlertList();
};
const handleSizeChange = (val) => {
  pagination.value.pageSize = val;
  fetchAlertList();
};
const handleCurrentChange = (val) => {
  pagination.value.page = val;
  fetchAlertList();
};
const handleCheckAlert = async () => {
  loading.value = true;
  try {
    const res = await inventoryApi.checkInventoryAlert();
    if (res.code === 0) {
      alertList.value = res.data.alerts;
      pagination.value.total = res.data.total;
      ElMessage.success('库存检查完成');
    }
  }
  catch (error) {
    console.error('检查库存预警失败:', error);
  }
  finally {
    loading.value = false;
  }
};
const handleViewProduct = (row) => {
  // 跳转到产品管理页面，并搜索该产品
  router.push({
    path: '/product',
    query: { search: row.product_name || row.product_code }
  });
};
const getAlertTypeText = (type) => {
  const map = {
    1: '低于安全库存',
    2: '低于最低库存',
    3: '高于最高库存'
  };
  return map[type] || '未知';
};
const getAlertTypeTagType = (type) => {
  if (type === null || type === undefined) {
    return 'info';
  }
  const map = {
    1: 'warning',
    2: 'danger',
    3: 'info'
  };
  return map[type] || 'info';
};
const getAlertLevelText = (level) => {
  const map = {
    1: '低',
    2: '中',
    3: '高'
  };
  return map[level] || '未知';
};
const getAlertLevelTagType = (level) => {
  if (level === null || level === undefined) {
    return 'warning';
  }
  const map = {
    1: 'success',
    2: 'warning',
    3: 'danger'
  };
  return map[level] || 'warning';
};
const getQuantityClass = (row) => {
  if (row.quantity <= row.min_stock) {
    return 'text-danger';
  }
  else if (row.quantity <= row.safety_stock) {
    return 'text-warning';
  }
  else if (row.quantity >= row.max_stock) {
    return 'text-info';
  }
  return '';
};
onMounted(() => {
  fetchAlertList();
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

.text-danger {
  color: #f56c6c;
  font-weight: bold;
}

.text-warning {
  color: #e6a23c;
  font-weight: bold;
}

.text-info {
  color: #67c23a;
  font-weight: bold;
}
</style>
