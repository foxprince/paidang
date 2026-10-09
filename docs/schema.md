# 拍档 · 数据表设计（PostgreSQL）

7 张表。状态用 varchar + CHECK 约束，不用 enum（enum 改起来疼）。

```sql
-- 陪练
CREATE TABLE coaches (
    id              BIGSERIAL PRIMARY KEY,
    openid          VARCHAR(64)  NOT NULL UNIQUE,          -- 微信 openid
    nickname        VARCHAR(64)  NOT NULL,
    avatar_url      VARCHAR(512),
    phone           VARCHAR(20),                            -- 展示时脱敏
    title           VARCHAR(64),                            -- 头衔，如"前省队"
    bio             TEXT,
    home_court      VARCHAR(128),                           -- 常出没球场
    video_url       VARCHAR(512),                           -- 30秒对打视频，发布必填
    price_per_hour  INTEGER NOT NULL DEFAULT 0,             -- 分
    pay_qr_url      VARCHAR(512),                           -- 微信收款码
    schedule        JSONB NOT NULL DEFAULT '{}',            -- 可约时段模板
    verified        BOOLEAN NOT NULL DEFAULT FALSE,         -- 平台验证徽章
    published       BOOLEAN NOT NULL DEFAULT FALSE,         -- 主页是否发布
    demerits        INTEGER NOT NULL DEFAULT 0,             -- 履约污点计数
    status          VARCHAR(16) NOT NULL DEFAULT 'active'
                      CHECK (status IN ('active','suspended','banned')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 学员（归属某个陪练）
CREATE TABLE students (
    id          BIGSERIAL PRIMARY KEY,
    coach_id    BIGINT NOT NULL REFERENCES coaches(id),
    name        VARCHAR(64) NOT NULL,
    phone       VARCHAR(20),                                -- 展示时脱敏
    wechat      VARCHAR(64),
    level       VARCHAR(32),                                -- 水平标签：新手/2.5/3.0…
    remark      TEXT,                                       -- 学员整体备注
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_students_coach ON students(coach_id);

-- 预约订单
CREATE TABLE bookings (
    id            BIGSERIAL PRIMARY KEY,
    coach_id      BIGINT NOT NULL REFERENCES coaches(id),
    student_id    BIGINT REFERENCES students(id),            -- 老学员关联；新学员下单时自动建档
    student_name  VARCHAR(64) NOT NULL,                      -- 下单时填写，快照
    student_phone VARCHAR(20) NOT NULL,
    play_date     DATE NOT NULL,
    start_time    TIME NOT NULL,
    end_time      TIME NOT NULL,
    service_type  VARCHAR(16) NOT NULL DEFAULT 'single',     -- v2 预留：拼单/团课
    price         INTEGER NOT NULL,                           -- 分，下单时快照
    status        VARCHAR(16) NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending','confirmed','completed','cancelled')),
    pay_status    VARCHAR(16) NOT NULL DEFAULT 'unpaid'
                    CHECK (pay_status IN ('unpaid','paid','reconciling')),
    no_show_by    VARCHAR(16),                                -- 爽约方：coach/student
    cancel_reason VARCHAR(256),
    idem_key      VARCHAR(64) NOT NULL UNIQUE,                -- 幂等键，防重复提交
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_bookings_coach_date ON bookings(coach_id, play_date);
CREATE INDEX idx_bookings_status ON bookings(status);

-- 按次课备注
CREATE TABLE lesson_notes (
    id          BIGSERIAL PRIMARY KEY,
    booking_id  BIGINT NOT NULL REFERENCES bookings(id),
    coach_id    BIGINT NOT NULL REFERENCES coaches(id),
    student_id  BIGINT NOT NULL REFERENCES students(id),
    key_points  TEXT NOT NULL,                                -- 本节要点
    next_plan   TEXT,                                        -- 下节计划
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_notes_student ON lesson_notes(student_id, created_at DESC);

-- 评价（仅已完成订单）
CREATE TABLE reviews (
    id          BIGSERIAL PRIMARY KEY,
    booking_id  BIGINT NOT NULL UNIQUE REFERENCES bookings(id),
    coach_id    BIGINT NOT NULL REFERENCES coaches(id),
    rating      SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment     TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_reviews_coach ON reviews(coach_id);

-- 投诉
CREATE TABLE complaints (
    id          BIGSERIAL PRIMARY KEY,
    booking_id  BIGINT NOT NULL REFERENCES bookings(id),
    evidence    JSONB NOT NULL DEFAULT '[]',                 -- 付款截图/聊天记录 URL
    detail      TEXT,
    status      VARCHAR(16) NOT NULL DEFAULT 'open'
                  CHECK (status IN ('open','upheld','rejected')),
    handled_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 幂等键（独立表，定期清理）
CREATE TABLE idempotency_keys (
    key         VARCHAR(64) PRIMARY KEY,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

## 约定
- 金额全部用"分"存整数，不用浮点
- 手机号只存原文，展示层统一脱敏（138****8000）
- `updated_at` 由应用层更新，不用触发器
- 迁移脚本放 `server/migrations/`，按 `001_xxx.sql` 顺序执行，不靠 AutoMigrate
