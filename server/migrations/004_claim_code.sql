-- 004: 陪练支持管理员手动建档 + 认领码
-- openid 允许空（未认领），空值不参与唯一约束
ALTER TABLE coaches ALTER COLUMN openid DROP NOT NULL;
ALTER TABLE coaches DROP CONSTRAINT IF EXISTS coaches_openid_key;
CREATE UNIQUE INDEX IF NOT EXISTS idx_coaches_openid ON coaches(openid) WHERE openid <> '';

ALTER TABLE coaches ADD COLUMN claim_code VARCHAR(6);
CREATE UNIQUE INDEX IF NOT EXISTS idx_coaches_claim_code ON coaches(claim_code) WHERE claim_code IS NOT NULL;
