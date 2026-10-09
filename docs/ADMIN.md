# 拍档 · 后台管理设计

## 形态
不做独立 Web 后台。**小程序内「管理」区**，`ADMIN_OPENIDS` 白名单（环境变量，逗号分隔）。

理由：不新增技术栈、不新增部署；管理员也是微信用户，手机上就能审。等单量起来再考虑 Web 版。

## 权限
- 服务端 `AdminOnly` 中间件：JWT → openid 在白名单里，否则 403
- 小程序端：`GET /admin/ping` 返回 `{is_admin}`，「我的」页据此显示管理入口
- 管理员自己的 openid：登录后调 `GET /coach/me`，返回里有 `openid`

## 模块

### A1 陪练审核
- 列表：待审核（已发布但未验证），按申请时间倒序
- 详情：对打视频、简介、价格、球场
- 操作：**通过**（verified=true）/ **下架**（published=false）
- 审核意见：通过/下架都要留一句话，记 admin_logs

### A0 手动建档（冷启动）
陪练自助注册是主路径（打开小程序自动建档），但招募期需要管理员代建：
- 管理员填姓名/手机/头衔/价格 → 生成 6 位认领码
- 微信把认领码发给陪练 → 他在「我的」页输入认领码 → 档案绑定到他的微信
- 认领后自动建档的空档案会被删掉，认领码作废
- 未认领的档案 openid 为空，不出现在学员端

### A2 投诉处理
- 列表：待处理（status=open），按创建时间；**超过 48 小时未处理的标红**
- 详情：订单信息、学员/陪练、证据图片、投诉说明
- 操作：**成立** / **不成立**
  - 成立 → 陪练 demerits+1，自动执行：
    - 累计 2 次 → 下架主页（published=false, status=suspended）
    - 累计 3 次 → 封号（status=banned）
  - 不成立 → 直接关闭，记 handled_at
- 规则写死在 service 层，不靠人记

### A3 履约管理
- 搜索陪练：看 demerits、status、主页状态
- 操作：手动下架 / 解禁 / 封号（二次确认）
- 所有操作记 admin_logs：谁、什么时候、对谁、做了什么

### A4 数据看板
- 陪练：总数、已发布、待审核
- 订单：总数、今日新增、各状态分布
- 收入：GMV（已收款）、本月 GMV
- 投诉：待处理数（>0 标红）

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /admin/ping | 是否管理员 `{is_admin}` |
| GET | /admin/coaches?verified=&status= | 陪练列表 |
| GET | /admin/coaches/:id | 陪练详情（含 demerits、订单数） |
| PATCH | /admin/coaches/:id/verify | 审核 `{verified, note}` |
| PATCH | /admin/coaches/:id/status | 改状态 `{status, note}` |
| GET | /admin/complaints?status= | 投诉列表 |
| GET | /admin/complaints/:id | 投诉详情（含订单、证据） |
| PATCH | /admin/complaints/:id | 裁决 `{result: upheld/rejected, note}` |
| GET | /admin/stats | 数据看板 |

## 数据表
只加一张：

```sql
CREATE TABLE admin_logs (
    id          BIGSERIAL PRIMARY KEY,
    admin_openid VARCHAR(64) NOT NULL,
    action      VARCHAR(32) NOT NULL,   -- verify_coach / set_status / judge_complaint
    target_type VARCHAR(32) NOT NULL,   -- coach / complaint
    target_id   BIGINT NOT NULL,
    detail      TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_admin_logs_target ON admin_logs(target_type, target_id);
```

其余复用现有表：complaints.status/handled_at、coaches.verified/published/status/demerits。

## 小程序页面
- pages/admin/ — 管理首页：四个入口 + 待办数字（待审核 n、待处理投诉 n）
- pages/admin-coaches/ — 审核列表（待审核 tab + 全部 tab）
- pages/admin-coach-detail/ — 审核详情 + 状态操作
- pages/admin-complaints/ — 投诉列表（待处理标红超时）
- pages/admin-complaint-detail/ — 投诉详情 + 裁决
- pages/admin-stats/ — 数据看板

「我的」页加一个管理入口（is_admin 才显示）。

## 不做
- 复杂的筛选和导出（MVP 不需要）
- 多角色权限（只有一种管理员）
- Web 版（单量起来再说）
