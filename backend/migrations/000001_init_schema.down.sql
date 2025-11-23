-- backend/migrations/000001_init_schema.down.sql
-- 初期スキーマのロールバック

DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS reservation_resources;
DROP TABLE IF EXISTS reservation_participants;
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS resources;
DROP TABLE IF EXISTS users;
