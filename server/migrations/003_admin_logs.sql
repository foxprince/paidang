-- 003: 管理操作审计日志
CREATE TABLE admin_logs (
    id           BIGSERIAL PRIMARY KEY,
    admin_openid VARCHAR(64) NOT NULL,
    action       VARCHAR(32) NOT NULL,
    target_type  VARCHAR(32) NOT NULL,
    target_id    BIGINT NOT NULL,
    detail       TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_admin_logs_target ON admin_logs(target_type, target_id);
