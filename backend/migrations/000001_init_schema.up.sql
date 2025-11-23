-- backend/migrations/000001_init_schema.up.sql
-- 初期スキーマ作成
--
-- このマイグレーションは以下のテーブルを作成します:
-- - users: ユーザー情報
-- - resources: 会議室・備品マスターデータ
-- - reservations: 予約基本情報（パーティション親テーブル）
-- - reservation_instances: 予約インスタンス展開テーブル
-- - reservation_participants: 予約参加者
-- - reservation_resources: 予約リソース（多対多）
-- - audit_logs: 監査ログ

-- ============================================================================
-- Users テーブル
-- ============================================================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sub VARCHAR(255) NOT NULL UNIQUE,  -- IdPから取得したユーザー識別子（不変）
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'GENERAL',  -- GENERAL, SECRETARY, MANAGER, ADMIN, AUDITOR
    manager_id UUID REFERENCES users(id),  -- 上長のユーザーID（承認フロー用）
    penalty_score INT NOT NULL DEFAULT 0,  -- キャンセルペナルティスコア
    penalty_score_expire_at TIMESTAMPTZ,  -- ペナルティスコア有効期限
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ  -- 論理削除
);

COMMENT ON TABLE users IS 'ユーザー情報テーブル';
COMMENT ON COLUMN users.sub IS 'IdP側の不変ユーザー識別子';
COMMENT ON COLUMN users.role IS 'ロール: GENERAL, SECRETARY, MANAGER, ADMIN, AUDITOR';
COMMENT ON COLUMN users.penalty_score IS 'キャンセルペナルティスコア（90日ローテーション）';

-- ============================================================================
-- Resources テーブル
-- ============================================================================
CREATE TABLE resources (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,  -- MEETING_ROOM, EQUIPMENT
    capacity INT,  -- 収容人数（会議室の場合）
    location VARCHAR(255),  -- 場所
    equipment JSONB,  -- 設備情報（プロジェクター、ホワイトボード等）
    required_role VARCHAR(50),  -- アクセス制御用の必要ロール
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE resources IS '会議室・備品マスターデータ';
COMMENT ON COLUMN resources.type IS 'リソース種別: MEETING_ROOM, EQUIPMENT';
COMMENT ON COLUMN resources.equipment IS '設備情報（JSON形式）';
COMMENT ON COLUMN resources.required_role IS '予約に必要な最低ロール';

-- ============================================================================
-- Reservations テーブル（パーティション親テーブル）
-- ============================================================================
CREATE TABLE reservations (
    id UUID NOT NULL DEFAULT uuid_generate_v4(),
    organizer_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    rrule VARCHAR(255),  -- 繰り返しルール (iCalendar RFC 5545形式)
    is_private BOOLEAN NOT NULL DEFAULT false,
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Tokyo',
    approval_status VARCHAR(50) NOT NULL DEFAULT 'CONFIRMED',  -- PENDING, CONFIRMED, REJECTED
    updated_by UUID REFERENCES users(id),
    version INT NOT NULL DEFAULT 1,  -- 楽観的ロック用バージョン
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,  -- 論理削除
    PRIMARY KEY (id, start_at)  -- パーティションキーを含む複合主キー
) PARTITION BY RANGE (start_at);

COMMENT ON TABLE reservations IS '予約基本情報（パーティション親テーブル）';
COMMENT ON COLUMN reservations.rrule IS '繰り返しルール（iCalendar RFC 5545形式）';
COMMENT ON COLUMN reservations.approval_status IS '承認状態: PENDING, CONFIRMED, REJECTED';
COMMENT ON COLUMN reservations.version IS '楽観的ロック用バージョン番号';

-- Note: パーティションテーブルのため、id単独のUNIQUE制約は作成できない
-- 外部キー参照には (id, start_at) の複合キーを使用する

-- 年別パーティション作成（2025年〜2027年）
CREATE TABLE reservations_2025 PARTITION OF reservations
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');

CREATE TABLE reservations_2026 PARTITION OF reservations
    FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');

CREATE TABLE reservations_2027 PARTITION OF reservations
    FOR VALUES FROM ('2027-01-01') TO ('2028-01-01');

-- ============================================================================
-- ReservationInstances テーブル
-- ============================================================================
CREATE TABLE reservation_instances (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reservation_id UUID NOT NULL,
    reservation_start_at TIMESTAMPTZ NOT NULL,  -- 親テーブルのパーティションキー（外部キー参照用）
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    original_start_at TIMESTAMPTZ,  -- 繰り返し例外時の元の開始日時
    status VARCHAR(20) NOT NULL DEFAULT 'CONFIRMED',  -- CONFIRMED, CANCELLED, CHECKED_IN, COMPLETED, NO_SHOW
    checked_in_at TIMESTAMPTZ,  -- チェックイン日時
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- 外部キー制約（パーティションテーブルへの参照には複合キーが必要）
    CONSTRAINT fk_reservation_instances_reservation
        FOREIGN KEY (reservation_id, reservation_start_at) 
        REFERENCES reservations(id, start_at) 
        ON DELETE CASCADE
);

COMMENT ON TABLE reservation_instances IS '予約インスタンス展開テーブル（繰り返し予定の個別インスタンス）';
COMMENT ON COLUMN reservation_instances.reservation_start_at IS '親予約の開始日時（パーティションキー、外部キー用）';
COMMENT ON COLUMN reservation_instances.original_start_at IS '繰り返し例外時の元の開始日時';
COMMENT ON COLUMN reservation_instances.status IS 'ステータス: CONFIRMED, CANCELLED, CHECKED_IN, COMPLETED, NO_SHOW';
COMMENT ON CONSTRAINT fk_reservation_instances_reservation ON reservation_instances IS 'reservations テーブルへの複合外部キー';

-- ============================================================================
-- ReservationParticipants テーブル
-- ============================================================================
CREATE TABLE reservation_participants (
    reservation_instance_id UUID NOT NULL REFERENCES reservation_instances(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    role VARCHAR(50) NOT NULL DEFAULT 'ATTENDEE',  -- ORGANIZER, ATTENDEE, APPROVER
    status VARCHAR(20) NOT NULL DEFAULT 'NEEDS_ACTION',  -- NEEDS_ACTION, ACCEPTED, DECLINED
    response_at TIMESTAMPTZ,  -- 回答日時
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (reservation_instance_id, user_id)
);

COMMENT ON TABLE reservation_participants IS '予約参加者テーブル';
COMMENT ON COLUMN reservation_participants.role IS '役割: ORGANIZER, ATTENDEE, APPROVER';
COMMENT ON COLUMN reservation_participants.status IS '参加ステータス: NEEDS_ACTION, ACCEPTED, DECLINED';

-- ============================================================================
-- ReservationResources テーブル（排他制御の対象）
-- ============================================================================
CREATE TABLE reservation_resources (
    reservation_instance_id UUID NOT NULL REFERENCES reservation_instances(id) ON DELETE CASCADE,
    resource_id UUID NOT NULL REFERENCES resources(id),
    is_required BOOLEAN NOT NULL DEFAULT true,  -- 必須リソースかどうか
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (reservation_instance_id, resource_id),
    UNIQUE (reservation_instance_id, resource_id)  -- 重複防止
);

COMMENT ON TABLE reservation_resources IS '予約リソーステーブル（多対多関係、排他制御対象）';
COMMENT ON COLUMN reservation_resources.is_required IS '必須リソースフラグ';

-- ============================================================================
-- AuditLogs テーブル
-- ============================================================================
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL,  -- CREATE_RESERVATION, UPDATE_RESERVATION, CANCEL_RESERVATION, etc.
    target_type VARCHAR(50) NOT NULL,  -- RESERVATION, RESOURCE, USER
    target_id UUID NOT NULL,
    details JSONB,  -- 詳細情報（JSON形式）
    ip_address INET,  -- IPアドレス
    user_agent TEXT,  -- ユーザーエージェント
    signature_hash VARCHAR(64),  -- 改ざん検知用署名ハッシュ（HMAC-SHA256）
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE audit_logs IS '監査ログテーブル（WORM/SIEM転送用）';
COMMENT ON COLUMN audit_logs.action IS '操作種別: CREATE_RESERVATION, UPDATE_RESERVATION, CANCEL_RESERVATION等';
COMMENT ON COLUMN audit_logs.signature_hash IS '改ざん検知用署名ハッシュ（HMAC-SHA256）';

-- ============================================================================
-- インデックス作成（Phase 1: カレンダー表示最適化）
-- ============================================================================

-- Users テーブル
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_role ON users(role) WHERE is_active = true;

-- Resources テーブル
CREATE INDEX idx_resources_type ON resources(type) WHERE is_active = true;
CREATE INDEX idx_resources_name_trgm ON resources USING gin(name gin_trgm_ops);  -- 全文検索用

-- Reservations テーブル（各パーティションに自動作成）
CREATE INDEX idx_reservations_organizer ON reservations(organizer_id, start_at DESC);
CREATE INDEX idx_reservations_time_range ON reservations USING gist(tstzrange(start_at, end_at));
CREATE INDEX idx_reservations_approval_status ON reservations(approval_status) WHERE deleted_at IS NULL;

-- ReservationInstances テーブル
CREATE INDEX idx_instances_reservation ON reservation_instances(reservation_id);
CREATE INDEX idx_instances_time_range ON reservation_instances USING gist(tstzrange(start_at, end_at));
CREATE INDEX idx_instances_status ON reservation_instances(status);

-- ReservationParticipants テーブル
CREATE INDEX idx_participants_user ON reservation_participants(user_id, status);

-- ReservationResources テーブル（排他制御用）
CREATE INDEX idx_resources_resource ON reservation_resources(resource_id);

-- AuditLogs テーブル
CREATE INDEX idx_audit_user ON audit_logs(user_id, created_at DESC);
CREATE INDEX idx_audit_action ON audit_logs(action, created_at DESC);
CREATE INDEX idx_audit_target ON audit_logs(target_type, target_id);

-- ============================================================================
-- トリガー関数: updated_at自動更新
-- ============================================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Users テーブルのトリガー
CREATE TRIGGER trigger_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Resources テーブルのトリガー
CREATE TRIGGER trigger_resources_updated_at
    BEFORE UPDATE ON resources
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Reservations テーブルのトリガー
CREATE TRIGGER trigger_reservations_updated_at
    BEFORE UPDATE ON reservations
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ReservationInstances テーブルのトリガー
CREATE TRIGGER trigger_instances_updated_at
    BEFORE UPDATE ON reservation_instances
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- 制約チェック関数
-- ============================================================================

-- 予約時間の妥当性チェック（start_at < end_at）
CREATE OR REPLACE FUNCTION check_reservation_time()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.start_at >= NEW.end_at THEN
        RAISE EXCEPTION '開始時刻は終了時刻より前である必要があります';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_check_reservation_time
    BEFORE INSERT OR UPDATE ON reservations
    FOR EACH ROW
    EXECUTE FUNCTION check_reservation_time();

CREATE TRIGGER trigger_check_instance_time
    BEFORE INSERT OR UPDATE ON reservation_instances
    FOR EACH ROW
    EXECUTE FUNCTION check_reservation_time();

-- ============================================================================
-- 完了メッセージ
-- ============================================================================
DO $$
BEGIN
    RAISE NOTICE '初期スキーマ作成完了';
    RAISE NOTICE 'テーブル数: 7';
    RAISE NOTICE 'パーティション: reservations_2025, reservations_2026, reservations_2027';
    RAISE NOTICE 'インデックス数: 14';
END $$;
