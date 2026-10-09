# 拍档 · API 接口

Base: `https://api.paidang.example/api/v1`（示例）

## 鉴权
陪练端接口需 `Authorization: Bearer <JWT>`。学员端公开接口免登录，下单时用短信验证码验手机号。

## 统一格式
```json
// 成功
{ "code": 0, "msg": "ok", "data": {} }
// 失败
{ "code": 40001, "msg": "预约时段已被占用", "data": null }
```
列表接口统一：`?page=1&page_size=20`，返回 `{ "list": [], "total": 0 }`。

## 端点

| 方法 | 路径 | 鉴权 | 说明 |
|------|------|------|------|
| POST | /auth/wechat-login | 无 | `{code}` → `{token, is_new}` |
| GET | /coach/me | 陪练 | 我的档案 |
| PUT | /coach/profile | 陪练 | 编辑主页（头像/简介/视频/价格/时段/收款码） |
| POST | /coach/publish | 陪练 | 发布/下架主页（视频必填校验） |
| GET | /coach/today | 陪练 | 今日课表 |
| GET | /bookings | 陪练 | 预约列表 `?status=&date=` |
| PATCH | /bookings/:id/confirm | 陪练 | 确认预约 |
| PATCH | /bookings/:id/reject | 陪练 | 拒绝 `{reason}` |
| PATCH | /bookings/:id/complete | 陪练 | 标记完成 |
| PATCH | /bookings/:id/cancel | 陪练 | 取消（记爽约方） |
| PATCH | /bookings/:id/mark-paid | 陪练 | 确认收款 → 流水入账 |
| GET | /students | 陪练 | 学员列表 |
| POST | /students | 陪练 | 手动建档 |
| GET | /students/:id | 陪练 | 学员详情 + 按次课备注时间线 |
| POST | /students/:id/notes | 陪练 | 写按次课备注 `{booking_id, key_points, next_plan}` |
| GET | /income/summary | 陪练 | `?month=` 本月/上月实收、订单数、课时数 |
| POST | /complaints | 陪练/学员 | 投诉 `{booking_id, evidence[], detail}` |
| GET | /public/coach/:id | 无 | 学员端：陪练主页展示 |
| GET | /public/coach/:id/slots | 无 | 学员端：`?date=` 可约时段 |
| POST | /public/bookings | 无 | 学员下单 `{coach_id, date, start, end, name, phone, sms_code, idem_key}` |
| POST | /public/sms-code | 无 | 发送短信验证码（限流 1/分钟） |
| POST | /public/reviews | 无 | 学员评价 `{booking_id, rating, comment}`（仅已完成订单） |

## 后台管理（需管理员白名单）
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /admin/ping | 是否管理员 |
| GET | /admin/coaches?verified=&status= | 陪练列表 |
| GET | /admin/coaches/:id | 陪练详情 |
| PATCH | /admin/coaches/:id/verify | 审核 `{verified, note}` |
| PATCH | /admin/coaches/:id/status | 改状态 `{status, note}` |
| GET | /admin/complaints?status= | 投诉列表 |
| GET | /admin/complaints/:id | 投诉详情 |
| PATCH | /admin/complaints/:id | 裁决 `{result: upheld/rejected, note}` |
| GET | /admin/stats | 数据看板 |
| POST | /admin/coaches | 管理员手动建档 `{name, phone, title, price_yuan, bio}` → 返回认领码 |

## 陪练端
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /coach/claim | 输入认领码绑定档案 `{claim_code}` → 返回新 token |

## 错误码
- `40001` 参数错误 / `40101` 未登录或 token 过期 / `40301` 无权限
- `40401` 不存在 / `40901` 时段冲突 / `42901` 触发限流

## 限流
公开接口按 IP：短信 1/分钟，下单 10/分钟；登录后接口按用户 120/分钟。
