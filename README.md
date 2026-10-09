# 拍档 (paidang)

给网球陪练用的经营工具：**预约 + 收款 + 学员管理** SaaS。

陪练 5 分钟生成个人预约主页，学员自助选时段下单，陪练微信确认，扫码收款，系统记流水、管学员。不做 C 端撮合，先做 B 端工具——**做陪练生意上的拍档**，把供给聚拢，再长出交易。

Slogan：好生意，需要好拍档。

## 技术栈

| 层 | 选型 |
|---|---|
| 小程序 | 原生微信小程序（WXML/WXSS/JS） |
| 后端 | Go 1.24 / Gin / GORM |
| 数据库 | PostgreSQL（直连现有实例，`paidang` 独立库） |
| 部署 | 二进制直跑（systemd），不用 Docker |

## 目录

```
server/       Go 后端（cmd/api 入口，internal 三层：handler/service/repository）
miniprogram/  微信小程序（陪练端 5 页 + 学员端 2 页）
docs/         PRD.md、schema.md、api.md
DESIGN.md     设计语言：极简、克制、艺术化，不要 AI 味
```

## 状态

MVP 开发中。第一版功能（8 个）：微信登录、陪练个人主页、在线预约、手动收款确认、学员管理、收入看板、今日课表、微信预约通知。详见 `docs/PRD.md`。

**合规红线**：永远不碰资金托管/二清——学员的钱直接到陪练微信，不经过平台账户。

## 本地开发

```bash
# 1. 在你的 PostgreSQL 里建库
createdb paidang

# 2. 配环境变量（别提交到仓库）
export DATABASE_URL="postgres://user:pass@host:5432/paidang?sslmode=disable"

# 3. 跑迁移 + 启动
cd server
go run ./cmd/api
```

## License

私有项目，保留所有权利。
