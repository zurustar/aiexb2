-- database/init/01_create_extensions.sql
-- PostgreSQL拡張機能の有効化

-- UUID生成用
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 全文検索用（将来的に使用）
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- タイムゾーン処理
SET timezone = 'Asia/Tokyo';
