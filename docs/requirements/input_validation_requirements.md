# 前后端输入参数校验需求文档

## 1. 需求概述

### 1.1 需求背景

当前系统缺少统一的输入参数校验机制，导致：
- 前端输入未经验证直接提交，可能导致无效数据进入系统
- 后端缺乏参数校验，存在安全风险和数据一致性问题
- 错误提示不友好，用户体验差

### 1.2 目标

建立统一的前后端输入参数校验规范，确保：
- 前端输入数据格式正确、符合业务规则
- 后端参数安全、防止恶意请求和SQL注入
- 错误提示清晰、友好，帮助用户快速定位问题

### 1.3 范围

覆盖系统所有核心模块：
- 用户管理
- 角色管理
- 权限管理
- 菜单管理
- 供应商管理
- 客户管理
- 产品管理
- 仓库管理
- 采购管理
- 销售管理
- 库存管理

---

## 2. 前端校验规范

### 2.1 通用校验规则

| 字段类型 | 校验规则 | 错误提示 |
|---------|---------|---------|
| 字符串必填 | 非空且长度 > 0 | "请输入{字段名}" |
| 字符串长度 | 最小长度、最大长度限制 | "{字段名}长度应在{min}-{max}之间" |
| 数字必填 | 非空且为有效数字 | "请输入{字段名}" |
| 数字范围 | 最小值、最大值限制 | "{字段名}应在{min}-{max}之间" |
| 邮箱格式 | 符合邮箱正则 | "请输入正确的邮箱格式" |
| 手机号格式 | 符合手机号正则 | "请输入正确的手机号" |
| 日期格式 | 符合YYYY-MM-DD格式 | "请选择正确的日期格式" |
| 下拉选择 | 值不为空字符串或undefined | "请选择{字段名}" |

### 2.2 字段级别校验规则

#### 2.2.1 用户管理模块

| 字段名 | 类型 | 必填 | 校验规则 |
|-------|------|------|---------|
| username | string | 是 | 长度3-50，字母数字下划线 |
| password | string | 是 | 长度6-100 |
| email | string | 否 | 邮箱格式 |
| phone | string | 否 | 手机号格式 |
| status | int | 是 | 0或1 |

#### 2.2.2 角色管理模块

| 字段名 | 类型 | 必填 | 校验规则 |
|-------|------|------|---------|
| name | string | 是 | 长度1-100 |
| code | string | 是 | 长度1-50，唯一 |
| status | int | 是 | 0或1 |

#### 2.2.3 供应商管理模块

| 字段名 | 类型 | 必填 | 校验规则 |
|-------|------|------|---------|
| name | string | 是 | 长度1-100 |
| code | string | 是 | 长度1-50，唯一 |
| contact | string | 否 | 长度1-50 |
| phone | string | 否 | 手机号格式 |

#### 2.2.4 客户管理模块

| 字段名 | 类型 | 必填 | 校验规则 |
|-------|------|------|---------|
| name | string | 是 | 长度1-100 |
| code | string | 是 | 长度1-50，唯一 |
| contact | string | 否 | 长度1-50 |
| phone | string | 否 | 手机号格式 |

#### 2.2.5 产品管理模块

| 字段名 | 类型 | 必填 | 校验规则 |
|-------|------|------|---------|
| name | string | 是 | 长度1-100 |
| code | string | 是 | 长度1-50，唯一 |
| price | float | 是 | > 0 |
| costPrice | float | 否 | >= 0 |
| stock | int | 否 | >= 0 |
| categoryId | int | 是 | > 0 |

#### 2.2.6 仓库管理模块

| 字段名 | 类型 | 必填 | 校验规则 |
|-------|------|------|---------|
| name | string | 是 | 长度1-100 |
| code | string | 是 | 长度1-50，唯一 |
| contact | string | 否 | 长度1-50 |
| phone | string | 否 | 手机号格式 |

#### 2.2.7 采购订单模块

| 字段名 | 类型 | 必填 | 校验规则 |
|-------|------|------|---------|
| supplierId | int | 是 | > 0 |
| warehouseId | int | 是 | > 0 |
| items | array | 是 | 长度 >= 1 |
| items[].productId | int | 是 | > 0 |
| items[].quantity | int | 是 | >= 1 |
| items[].unitPrice | float | 是 | > 0 |

#### 2.2.8 销售订单模块

| 字段名 | 类型 | 必填 | 校验规则 |
|-------|------|------|---------|
| customerId | int | 是 | > 0 |
| warehouseId | int | 是 | > 0 |
| items | array | 是 | 长度 >= 1 |
| items[].productId | int | 是 | > 0 |
| items[].quantity | int | 是 | >= 1 |
| items[].unitPrice | float | 是 | > 0 |

---

## 3. 后端校验规范

### 3.1 通用校验规则

| 校验类型 | 说明 | 处理方式 |
|---------|------|---------|
| 必填校验 | 非空检查 | 返回400错误，提示字段必填 |
| 类型校验 | 数据类型匹配 | 返回400错误，提示类型错误 |
| 范围校验 | 数值范围检查 | 返回400错误，提示范围错误 |
| 格式校验 | 邮箱、手机号等格式 | 返回400错误，提示格式错误 |
| 长度校验 | 字符串长度限制 | 返回400错误，提示长度错误 |
| 唯一校验 | 唯一字段检查 | 返回400错误，提示已存在 |
| SQL注入防护 | 防止恶意SQL | 使用参数化查询，过滤特殊字符 |
| 安全过滤 | XSS防护 | 过滤HTML标签和脚本 |

### 3.2 后端API校验要求

所有API接口必须实现：

1. **参数绑定校验**：使用Go结构体tag进行自动校验
2. **业务规则校验**：在Logic层实现业务规则校验
3. **统一错误响应**：返回统一格式的错误响应

### 3.3 错误响应格式

```json
{
  "code": 400,
  "data": null,
  "message": "参数校验失败：username不能为空"
}
```

---

## 4. 前端校验实现方案

### 4.1 使用Element Plus表单校验

```vue
<el-form :model="form" :rules="rules" ref="formRef">
  <el-form-item label="用户名" prop="username">
    <el-input v-model="form.username"></el-input>
  </el-form-item>
</el-form>
```

### 4.2 校验规则定义

```javascript
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度应在3-50之间', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 100, message: '密码长度应在6-100之间', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
}
```

---

## 5. 后端校验实现方案

### 5.1 使用Go结构体tag校验

```go
type LoginRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Password string `json:"password" validate:"required,min=6,max=100"`
}
```

### 5.2 在Logic层实现校验

```go
func (l *CreateUserLogic) CreateUser(req *types.CreateUserRequest) (resp *types.CreateUserResponse, err error) {
    // 参数校验
    if err := l.svcCtx.Validator.Validate(req); err != nil {
        return nil, errors.New("参数校验失败：" + err.Error())
    }
    
    // 业务规则校验
    if req.Price <= 0 {
        return nil, errors.New("单价必须大于0")
    }
    
    // 唯一校验
    exists, err := l.svcCtx.UserModel.ExistsByUsername(req.Username)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("用户名已存在")
    }
    
    // 业务逻辑...
}
```

---

## 6. 测试方案

### 6.1 前端测试用例

| 测试场景 | 输入 | 预期结果 |
|---------|------|---------|
| 空用户名登录 | username="", password="123456" | 提示"请输入用户名" |
| 短用户名 | username="ab", password="123456" | 提示"用户名长度应在3-50之间" |
| 无效邮箱 | email="invalid" | 提示"请输入正确的邮箱格式" |
| 无效手机号 | phone="123456789" | 提示"请输入正确的手机号" |
| 负数单价 | price=-1 | 提示"单价必须大于0" |

### 6.2 后端测试用例

| 测试场景 | 请求参数 | 预期结果 |
|---------|---------|---------|
| 空参数 | {} | 返回code=400，提示必填字段 |
| 无效类型 | {"username": 123} | 返回code=400，提示类型错误 |
| 重复用户名 | {"username": "admin", ...} | 返回code=400，提示已存在 |
| SQL注入尝试 | {"username": "' OR 1=1 --"} | 返回code=400或正常过滤 |

---

## 7. 实施计划

| 阶段 | 任务 | 负责人 | 时间 |
|------|------|-------|------|
| 阶段1 | 需求分析与文档编写 | AgentA | 1天 |
| 阶段2 | 前端校验实现 | AgentB | 2天 |
| 阶段3 | 后端校验实现 | AgentC | 2天 |
| 阶段4 | 测试文档编写与测试 | AgentD | 1天 |
| 阶段5 | Bug修复与优化 | 对应Agent | 1天 |

---

## 8. 交付物

| 交付物 | 描述 | 位置 |
|-------|------|------|
| 需求文档 | 输入参数校验需求规范 | docs/requirements/input_validation_requirements.md |
| 前端代码 | 包含校验规则的Vue组件 | frontend/src/views/**/*.vue |
| 后端代码 | 包含校验逻辑的Go文件 | admin/internal/logic/**/*.go |
| 测试文档 | 测试用例与测试报告 | test/ |

---

## 9. 注意事项

1. 前后端校验规则需保持一致
2. 错误提示需友好、清晰
3. 避免过度校验影响用户体验
4. 敏感信息（如密码）校验后不应返回详细错误信息
5. 性能敏感接口需考虑校验开销

---

**文档版本**: v1.0  
**创建时间**: 2026-04-29  
**最后更新**: 2026-04-29