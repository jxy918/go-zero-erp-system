# 需求文档 - 库存管理功能分析与改进

## 1. 概述

本文档对当前库存管理功能进行深入分析，识别存在的问题，并提出具体的改进建议。

## 2. 现状分析

### 2.1 数据模型现状

#### 核心数据表

| 表名 | 说明 | 关联 |
|------|------|------|
| `products` | 产品表 | 主表 |
| `warehouses` | 仓库表 | 主表 |
| `warehouse_inventories` | 仓库库存表 | 实时库存 |
| `inventory_records` | 库存流水记录表 | 库存变动明细 |
| `inventory_changes` | 库存变动表 | 调整记录 |

#### 产品表 (products) 库存字段
- `stock` - 总库存数量（汇总字段，由流水计算得出）

#### 仓库库存表 (warehouse_inventories)
```go
type WarehouseInventory struct {
    ID          uint
    ProductID   uint  // 产品ID
    WarehouseID uint  // 仓库ID
    Quantity    int   // 当前库存数量
}
```

#### 库存流水表 (inventory_records)
```go
type InventoryRecord struct {
    ID          uint  // 主键
    ProductID   uint  // 产品ID
    WarehouseID uint  // 仓库ID
    Quantity    int   // 变动数量（正数入库，负数出库）
    Type        int   // 1:入库, 2:出库, 3:调整
    OrderID     uint  // 关联订单ID
    OrderType   int   // 1:采购订单, 2:销售订单
}
```

#### 库存变动表 (inventory_changes)
```go
type InventoryChange struct {
    ID              uint  // 主键
    ProductID       uint  // 产品ID
    WarehouseID     uint  // 仓库ID
    BeforeQuantity  int   // 变动前数量
    AfterQuantity   int   // 变动后数量
    Quantity        int   // 变动数量
    Type            int   // 1:入库, 2:出库, 3:调整
    Remark          string // 调整原因
}
```

### 2.2 业务流程现状

#### 入库流程
```
采购入库 → inventory_records (type=1) → warehouse_inventories 更新
```

#### 出库流程
```
销售出库 → inventory_records (type=2) → warehouse_inventories 更新
```

#### 库存调整流程
```
手动调整 → inventory_changes (type=3) → inventory_records (type=3)
```

### 2.3 API 接口现状

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 库存列表 | GET | /inventory/list | 分页查询库存流水 |
| 库存调整 | POST | /inventory/adjust | 手动调整库存 |
| 库存历史 | GET | /inventory/history | 查询某产品库存变动记录 |
| 库存预警 | GET | /erp/statistics/inventory-alert | ERP统计-库存预警 |

### 2.4 前端页面现状

| 功能 | 说明 | 入口 |
|------|------|------|
| 库存列表 | 显示库存流水记录 | ERP管理 → 库存管理 |
| 库存调整 | 手动增加/减少库存 | 库存列表 → 调整按钮 |
| 库存历史 | 查看某产品库存变动明细 | 库存列表 → 查看记录按钮 |

---

## 3. 存在问题

### 3.1 数据模型问题

#### 问题 1: 库存数据分散在两张表
- **现状**: `warehouse_inventories` 和 `inventory_records` 都存储库存数据
- **问题**:
  - 两表数据可能不一致
  - `warehouse_inventories.quantity` 是实时库存，但 `products.stock` 也声称是总库存
  - 缺乏唯一数据源

#### 问题 2: 库存计算逻辑混乱
- **现状**:
  - `InventoryModel.GetStockByProduct()` 从 `inventory_records` 计算
  - `adjustInventoryLogic.go` 直接操作 `warehouse_inventories`
  - 没有明确哪个是"真相来源"

#### 问题 3: 库存类型定义不一致
- **现状**: `inventory_records.type` 和 `inventory_changes.type` 都用 1/2/3 表示类型
- **问题**: 如果增加新的库存类型，两边都需要改

### 3.2 库存调整功能问题

#### 问题 4: 库存调整同时写两张表
```go
// adjustInventoryLogic.go:111-117
inv := model.InventoryRecord{...}
tx.Create(&inv)  // 写流水表

change := &model.InventoryChange{...}
l.svcCtx.InventoryChangeModel.Create(change)  // 写变动表（在事务外）
```
- **问题**: `inventory_changes` 在事务外创建，可能导致数据不一致

#### 问题 5: 库存调整不更新 products.stock
- **现状**: 库存调整只更新 `warehouse_inventories`，不更新 `products.stock`
- **问题**: `products.stock` 无法准确反映总库存

#### 问题 6: 库存调整缺少审核流程
- **现状**: 库存调整直接生效，没有审批流程
- **问题**:
  - 无法追溯谁调整的、为什么调整
  - 无法防止误操作

### 3.3 业务逻辑问题

#### 问题 7: 库存扣减没有乐观锁
- **现状**: 直接 UPDATE，可能出现超卖
```go
// 直接更新，没有版本控制
tx.Model(&Product{}).Where("id = ?", req.ProductID).Update("stock", gorm.Expr("stock + ?", quantity))
```

#### 问题 8: 库存查询语义不清晰
- **现状**: `/inventory/list` 查询的是 `inventory_records`（流水），不是实时库存
- **问题**: 用户期望看到的是当前库存，而不是流水记录

### 3.4 前端交互问题

#### 问题 9: 库存列表显示混乱
- **现状**: 列表显示流水记录，用户难以理解
- **问题**:
  - 一条入库记录后跟着一条出库记录，用户不知道当前剩多少
  - 需要计算才能得出当前库存

#### 问题 10: 缺少库存预警配置
- **现状**: 前端硬编码库存预警阈值（≤10 预警，≤50 提示）
- **问题**:
  - 预警阈值应该可配置
  - 应该支持不同产品不同阈值

---

## 4. 改进建议

### 4.1 数据模型改进

#### 建议 1: 明确库存数据源

**方案 A（推荐）: 以流水为准**
```
inventory_records 表作为唯一数据源
warehouse_inventories 作为缓存（可选）
products.stock 通过 SUM(quantity) 计算
```

**优点**:
- 数据一致性好
- 可以追溯任意时间点的库存
- 报表准确

**实现方式**:
```go
// 获取实时库存 = 从流水计算
func (m *InventoryModel) GetRealTimeStock(productID, warehouseID uint) (int, error) {
    var total int
    err := m.db.Model(&InventoryRecord{}).
        Select("COALESCE(SUM(quantity), 0)").
        Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).
        Scan(&total).Error
    return total, err
}
```

**方案 B: 以库存表为准**
```
warehouse_inventories 作为数据源
inventory_records 只记录流水
```

**适用场景**: 高并发写入场景

#### 建议 2: 统一库存类型枚举

```go
// 库存变动类型常量
const (
    InventoryTypeInbound   = 1  // 入库
    InventoryTypeOutbound = 2  // 出库
    InventoryTypeAdjust   = 3  // 调整
    InventoryTypeReturn   = 4  // 退货
    InventoryTypeLoss     = 5  // 报损
)
```

#### 建议 3: 添加库存调整审核表

```go
type InventoryAdjustRequest struct {
    ID            uint
    ProductID     uint
    WarehouseID   uint
    Quantity      int        // 调整数量（正数增加，负数减少）
    Reason        string     // 调整原因
    Status        int        // 1:待审核, 2:已审核, 3:已拒绝
    ApplicantID   uint       // 申请人
    ApproverID    uint       // 审核人
    ApproveTime   *time.Time // 审核时间
    Remark        string     // 审核备注
}
```

### 4.2 库存调整功能改进

#### 建议 4: 库存调整改为"申请-审核"模式

```
用户提交调整申请 → 管理员审核 → 执行调整
```

**流程**:
1. 用户提交库存调整申请（不直接生效）
2. 管理员审核申请
3. 审核通过后，执行调整并记录日志

**优点**:
- 防止误操作
- 保留审核记录
- 责任可追溯

#### 建议 5: 添加库存调整单据号

```go
type InventoryAdjustRequest struct {
    RequestNo    string  // 单据号: ADJ-20240101-001
    ...
}
```

**格式**: `ADJ-{日期}-{序号}`

#### 建议 6: 库存调整事务保证

```go
func (l *AdjustInventoryLogic) AdjustInventory(req *types.AdjustInventoryRequest) error {
    return model.DB.Transaction(func(tx *gorm.DB) error {
        // 1. 锁定库存行（悲观锁）
        var inventory WarehouseInventory
        err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
            Where("product_id = ? AND warehouse_id = ?", req.ProductID, req.WarehouseID).
            First(&inventory).Error

        // 2. 检查库存是否足够
        if inventory.Quantity + req.Quantity < 0 {
            return errors.New("库存不足")
        }

        // 3. 更新库存
        inventory.Quantity += req.Quantity
        if err := tx.Save(&inventory).Error; err != nil {
            return err
        }

        // 4. 记录流水（在同一个事务中）
        record := InventoryRecord{...}
        if err := tx.Create(&record).Error; err != nil {
            return err
        }

        return nil
    })
}
```

### 4.3 业务逻辑改进

#### 建议 7: 添加库存预警功能

```go
type ProductAlert struct {
    ID          uint
    ProductID   uint
    WarehouseID uint  // 0 表示所有仓库
    AlertType   int   // 1:低于阈值, 2:高于阈值
    Threshold   int   // 阈值
    Enabled     bool
}
```

**触发条件**:
- 库存 ≤ 预警阈值 → 发送预警通知

#### 建议 8: 区分"库存查询"和"流水查询"

| 接口 | 说明 | 数据来源 |
|------|------|----------|
| `/inventory/current` | 当前库存查询 | warehouse_inventories 或计算 |
| `/inventory/list` | 库存流水查询 | inventory_records |
| `/inventory/changes` | 库存变动查询 | inventory_changes |

### 4.4 前端交互改进

#### 建议 9: 库存列表显示当前库存

```
┌─────────────────────────────────────────────────────────────┐
│  产品名称    │ 当前库存 │ 仓库    │ 最近变动 │ 状态        │
├─────────────────────────────────────────────────────────────┤
│  产品A       │ 100      │ 仓库1   │ 2024-01  │ ⚠️ 库存预警 │
│  产品B       │ 50       │ 仓库2   │ 2024-01  │ ✓ 正常      │
└─────────────────────────────────────────────────────────────┘
```

**改进点**:
- 默认显示当前库存，而不是流水
- 增加"查看流水"按钮
- 库存预警状态可视化

#### 建议 10: 库存调整页面改进

**改进点**:
- 调整类型选择：入库/出库/调整
- 调整数量：输入框（支持负数或正数）
- 调整原因：必填，下拉选择 + 文本输入
- 预估结果：调整后库存预览

---

## 5. 改进优先级

### P0（必须修复）
1. 库存调整事务完整性（inventory_changes 在事务外）
2. 库存计算逻辑统一（明确数据源）

### P1（重要）
3. 添加库存审核流程
4. 库存扣减乐观锁/悲观锁
5. products.stock 同步更新

### P2（优化）
6. 库存预警配置化
7. 前端库存列表显示改进
8. 区分库存查询和流水查询接口

---

## 6. 附录

### 6.1 术语说明

| 术语 | 说明 |
|------|------|
| 实时库存 | 某仓库某产品的当前库存数量 |
| 库存流水 | 库存变动的历史记录 |
| 库存变动 | 调整单导致的库存变化 |
| 预警阈值 | 触发库存预警的库存数量 |

### 6.2 相关文件

| 文件 | 说明 |
|------|------|
| `admin/internal/model/inventoryModel.go` | 库存模型 |
| `admin/internal/logic/adjustInventoryLogic.go` | 库存调整逻辑 |
| `admin/internal/logic/listInventoryLogic.go` | 库存列表逻辑 |
| `frontend/src/views/inventory/Inventory.vue` | 库存管理页面 |

---

**文档版本**: v1.0
**创建日期**: 2026-04-30
**状态**: 待评审
