# AGENTS.md

## 在做的事
**拍档 (paidang)** — 网球陪练经营工具：预约 + 收款 + 学员管理。仓库：`foxprince/paidang`（单仓）。

## 技术
- Go 1.22+ / Gin / GORM；PostgreSQL（现有实例，`paidang` 独立库，`DATABASE_URL` 配置）
- 原生微信小程序（WXML/WXSS/JS），不用跨端框架
- 不用 Docker；后端二进制直跑
- 密钥只走环境变量，永不进仓库

## 目录
- `server/`：`cmd/api` 入口；`internal/` 下 handler / service / repository / model / middleware / wechat
- `miniprogram/`：陪练端 5 页（今日课表、预约、学员、收入、我的）+ 学员端 2 页（主页、预约下单）
- `docs/`：PRD、schema、api

## 习惯
- 代码极简：够用、好改；不做过度分层，不过早抽象
- 表结构变更走 `server/migrations/` SQL，不靠 AutoMigrate
- 中文界面用简体；注释中文，短
- 设计语言见 `DESIGN.md`：极简、克制、艺术化，不要 AI 味

## 安全 P0
- `code2session` 只在服务端调用；JWT；限流；参数化查询；上传限制；HTTPS；手机号脱敏
- 平台不碰钱：学员付款直达陪练微信，不做托管和分账
