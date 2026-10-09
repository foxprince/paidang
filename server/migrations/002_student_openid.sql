-- 002: 预约表加下单人 openid（用于微信订阅消息）
ALTER TABLE bookings ADD COLUMN student_openid VARCHAR(64);
CREATE INDEX idx_bookings_student_openid ON bookings(student_openid);
