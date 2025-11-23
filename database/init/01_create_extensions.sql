-- database/init/01_create_extensions.sql
-- PostgreSQL拡張機能の有効化
-- 
-- このスクリプトはデータベース初期化時に実行され、
-- ESMSで使用する各種PostgreSQL拡張機能を有効化します。

-- UUID生成機能（主キー用）
-- UUIDv4を生成するための拡張機能
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- pgcrypto（パスワードハッシュ、署名ハッシュ用）
-- 暗号化関数を提供（監査ログの署名ハッシュ生成に使用）
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- pg_trgm（全文検索用）
-- トライグラムインデックスによる高速な部分一致検索を提供
-- 会議室名、予定タイトルの検索に使用
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- btree_gist（時間範囲検索用）
-- GiSTインデックスでB-tree演算子をサポート
-- 予約の時間範囲重複チェックに使用
CREATE EXTENSION IF NOT EXISTS "btree_gist";

-- pgvector（将来のAI機能用 - Phase 2）
-- ベクトル類似度検索を提供
-- 会議準備サジェスト機能で使用予定
-- Note: Phase 1では使用しないが、事前にインストール
CREATE EXTENSION IF NOT EXISTS "vector";

-- タイムゾーン設定
-- デフォルトタイムゾーンを日本時間に設定
SET timezone = 'Asia/Tokyo';

-- 拡張機能の確認
-- インストールされた拡張機能を表示
SELECT extname, extversion 
FROM pg_extension 
WHERE extname IN ('uuid-ossp', 'pgcrypto', 'pg_trgm', 'btree_gist', 'vector')
ORDER BY extname;
