-- backend/migrations/000001_init_schema.down.sql
-- 初期スキーマのロールバック
--
-- このマイグレーションは000001_init_schema.up.sqlで作成した
-- 全てのテーブル、トリガー、関数を削除します。
-- 依存関係の逆順で削除を実行します。

-- ============================================================================
-- トリガーの削除
-- ============================================================================
DROP TRIGGER IF EXISTS trigger_check_instance_time ON reservation_instances;
DROP TRIGGER IF EXISTS trigger_check_reservation_time ON reservations;
DROP TRIGGER IF EXISTS trigger_instances_updated_at ON reservation_instances;
DROP TRIGGER IF EXISTS trigger_reservations_updated_at ON reservations;
DROP TRIGGER IF EXISTS trigger_resources_updated_at ON resources;
DROP TRIGGER IF EXISTS trigger_users_updated_at ON users;

-- ============================================================================
-- トリガー関数の削除
-- ============================================================================
DROP FUNCTION IF EXISTS check_reservation_time();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- ============================================================================
-- テーブルの削除（依存関係の逆順）
-- ============================================================================
-- 注: CASCADE指定により、以下も自動的に削除されます:
--   - 外部キー制約 (fk_reservation_instances_reservation等)
--   - その他の制約とインデックス

-- 監査ログテーブル（他テーブルへの依存なし）
DROP TABLE IF EXISTS audit_logs CASCADE;

-- 予約関連の中間テーブル（外部キー依存あり）
DROP TABLE IF EXISTS reservation_resources CASCADE;
DROP TABLE IF EXISTS reservation_participants CASCADE;

-- 予約インスタンステーブル（reservationsへの外部キー制約あり）
DROP TABLE IF EXISTS reservation_instances CASCADE;

-- 予約テーブル（パーティション親テーブル）
-- パーティション子テーブルも自動的に削除される
DROP TABLE IF EXISTS reservations CASCADE;

-- リソーステーブル
DROP TABLE IF EXISTS resources CASCADE;

-- ユーザーテーブル
DROP TABLE IF EXISTS users CASCADE;

-- ============================================================================
-- 完了メッセージ
-- ============================================================================
DO $$
BEGIN
    RAISE NOTICE '初期スキーマのロールバック完了';
    RAISE NOTICE '削除されたテーブル: users, resources, reservations (パーティション含む), reservation_instances, reservation_participants, reservation_resources, audit_logs';
    RAISE NOTICE '削除されたトリガー: 6個';
    RAISE NOTICE '削除された関数: 2個';
    RAISE NOTICE '削除された制約: 外部キー制約、UNIQUE制約等（CASCADE により自動削除）';
END $$;
