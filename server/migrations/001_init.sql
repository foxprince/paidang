-- 拍档 MVP 初始化表结构（PostgreSQL）
-- 执行顺序：按文件名前缀顺序

CREATE TABLE coaches (
    id              BIGSERIAL PRIMARY KEY,
    openid          VARCHAR(64)  NOT NULL UNIQUE,
    nickname        VARCHAR(64)  NOT NULL,
    avatar_url      VARCHAR(512),
    phone           VARCHAR(20),
    title           VARCHAR(64),
    bio             TEXT,
    home_court      VARCHAR(128),
    video_url       VARCHAR(512),
    price_per_hour  INTEGER NOT NULL DEFAULT 0,
    pay_qr_url      VARCHAR(512),
    schedule        JSONB NOT NULL DEFAULT '{}',
    verified        BOOLEAN NOT NULL DEFAULT FALSE,
    published       BOOLEAN NOT NULL DEFAULT FALSE,
    demerits        INTEGER NOT NULL DEFAULT 0,
    status          VARCHAR(16) NOT NULL DEFAULT 'active'
                      CHECK (status IN ('active','suspended','banned')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE students (
    id          BIGSERIAL PRIMARY KEY,
    coach_id    BIGINT NOT NULL REFERENCES coaches(id),
    name        VARCHAR(64) NOT NULL,
    phone       VARCHAR(20),
    wechat      VARCHAR(64),
    level       VARCHAR(32),
    remark      TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_students_coach ON students(coach_id);

CREATE TABLE bookings (
    id            BIGSERIAL PRIMARY KEY,
    coach_id      BIGINT NOT NULL REFERENCES coaches(id),
    student_id    BIGINT REFERENCES students(id),
    student_name  VARCHAR(64) NOT NULL,
    student_phone VARCHAR(20) NOT NULL,
    play_date     DATE NOT NULL,
    start_time    TIME NOT NULL,
    end_time      TIME NOT NULL,
    service_type  VARCHAR(16) NOT NULL DEFAULT 'single',
    price         INTEGER NOT NULL,
    status        VARCHAR(16) NOT NULL DEFAULT 'pending'
                    CHECK (status IN ('pending','confirmed','completed','cancelled')),
    pay_status    VARCHAR(16) NOT NULL DEFAULT 'unpaid'
                    CHECK (pay_status IN ('unpaid','paid','reconciling')),
    no_show_by    VARCHAR(16),
    cancel_reason VARCHAR(256),
    idem_key      VARCHAR(64) NOT NULL UNIQUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_bookings_coach_date ON bookings(coach_id, play_date);
CREATE INDEX idx_bookings_status ON bookings(status);

CREATE TABLE lesson_notes (
    id          BIGSERIAL PRIMARY KEY,
    booking_id  BIGINT NOT NULL REFERENCES bookings(id),
    coach_id    BIGINT NOT NULL REFERENCES coaches(id),
    student_id  BIGINT NOT NULL REFERENCES students(id),
    key_points  TEXT NOT NULL,
    next_plan   TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_notes_student ON lesson_notes(student_id, created_at DESC);

CREATE TABLE reviews (
    id          BIGSERIAL PRIMARY KEY,
    booking_id  BIGINT NOT NULL UNIQUE REFERENCES bookings(id),
    coach_id    BIGINT NOT NULL REFERENCES coaches(id),
    rating      SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment     TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_reviews_coach ON reviews(coach_id);

CREATE TABLE complaints (
    id          BIGSERIAL PRIMARY KEY,
    booking_id  BIGINT NOT NULL REFERENCES bookings(id),
    evidence    JSONB NOT NULL DEFAULT '[]',
    detail      TEXT,
    status      VARCHAR(16) NOT NULL DEFAULT 'open'
                  CHECK (status IN ('open','upheld','rejected')),
    handled_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE idempotency_keys (
    key         VARCHAR(64) PRIMARY KEY,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
