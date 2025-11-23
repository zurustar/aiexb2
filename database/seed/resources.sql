-- database/seed/resources.sql
-- リソース（会議室・備品）のシードデータ
--
-- テスト・開発環境用のサンプルデータを投入します。
-- 本番環境では実際のリソース情報に置き換えてください。

-- ============================================================================
-- 会議室データ
-- ============================================================================

-- 小会議室（4-6名）
INSERT INTO resources (id, name, type, capacity, location, equipment, required_role, is_active) VALUES
(
    '10000000-0000-0000-0000-000000000001'::uuid,
    'A会議室',
    'MEETING_ROOM',
    4,
    '本社ビル 3階',
    '{"projector": true, "whiteboard": true, "tv": false, "video_conference": false}'::jsonb,
    NULL,  -- 全ロールが予約可能
    true
),
(
    '10000000-0000-0000-0000-000000000002'::uuid,
    'B会議室',
    'MEETING_ROOM',
    6,
    '本社ビル 3階',
    '{"projector": true, "whiteboard": true, "tv": true, "video_conference": true}'::jsonb,
    NULL,
    true
);

-- 中会議室（8-12名）
INSERT INTO resources (id, name, type, capacity, location, equipment, required_role, is_active) VALUES
(
    '10000000-0000-0000-0000-000000000003'::uuid,
    'C会議室',
    'MEETING_ROOM',
    8,
    '本社ビル 4階',
    '{"projector": true, "whiteboard": true, "tv": true, "video_conference": true, "microphone": true}'::jsonb,
    NULL,
    true
),
(
    '10000000-0000-0000-0000-000000000004'::uuid,
    'D会議室',
    'MEETING_ROOM',
    10,
    '本社ビル 4階',
    '{"projector": true, "whiteboard": true, "tv": true, "video_conference": true, "microphone": true}'::jsonb,
    NULL,
    true
),
(
    '10000000-0000-0000-0000-000000000005'::uuid,
    'E会議室',
    'MEETING_ROOM',
    12,
    '本社ビル 5階',
    '{"projector": true, "whiteboard": true, "tv": true, "video_conference": true, "microphone": true, "recording": true}'::jsonb,
    NULL,
    true
);

-- 大会議室（20-30名）
INSERT INTO resources (id, name, type, capacity, location, equipment, required_role, is_active) VALUES
(
    '10000000-0000-0000-0000-000000000006'::uuid,
    'F会議室（大）',
    'MEETING_ROOM',
    20,
    '本社ビル 6階',
    '{"projector": true, "whiteboard": true, "tv": true, "video_conference": true, "microphone": true, "recording": true, "sound_system": true}'::jsonb,
    NULL,
    true
),
(
    '10000000-0000-0000-0000-000000000007'::uuid,
    'G会議室（大）',
    'MEETING_ROOM',
    30,
    '本社ビル 6階',
    '{"projector": true, "whiteboard": true, "tv": true, "video_conference": true, "microphone": true, "recording": true, "sound_system": true}'::jsonb,
    NULL,
    true
);

-- 役員会議室（アクセス制限あり）
INSERT INTO resources (id, name, type, capacity, location, equipment, required_role, is_active) VALUES
(
    '10000000-0000-0000-0000-000000000008'::uuid,
    '役員会議室',
    'MEETING_ROOM',
    15,
    '本社ビル 7階',
    '{"projector": true, "whiteboard": true, "tv": true, "video_conference": true, "microphone": true, "recording": true, "sound_system": true, "secure": true}'::jsonb,
    'MANAGER',  -- マネージャー以上のみ予約可能
    true
),
(
    '10000000-0000-0000-0000-000000000009'::uuid,
    '社長室会議スペース',
    'MEETING_ROOM',
    8,
    '本社ビル 7階',
    '{"projector": true, "whiteboard": true, "tv": true, "video_conference": true, "microphone": true, "recording": true, "secure": true}'::jsonb,
    'ADMIN',  -- 管理者のみ予約可能
    true
);

-- ============================================================================
-- 備品データ
-- ============================================================================

-- プロジェクター（持ち運び可能）
INSERT INTO resources (id, name, type, capacity, location, equipment, required_role, is_active) VALUES
(
    '20000000-0000-0000-0000-000000000001'::uuid,
    'モバイルプロジェクター #1',
    'EQUIPMENT',
    NULL,
    '本社ビル 3階 備品室',
    '{"portable": true, "hdmi": true, "wireless": true, "brightness": "3000lm"}'::jsonb,
    NULL,
    true
),
(
    '20000000-0000-0000-0000-000000000002'::uuid,
    'モバイルプロジェクター #2',
    'EQUIPMENT',
    NULL,
    '本社ビル 4階 備品室',
    '{"portable": true, "hdmi": true, "wireless": true, "brightness": "3000lm"}'::jsonb,
    NULL,
    true
);

-- ホワイトボード（持ち運び可能）
INSERT INTO resources (id, name, type, capacity, location, equipment, required_role, is_active) VALUES
(
    '20000000-0000-0000-0000-000000000003'::uuid,
    'モバイルホワイトボード #1',
    'EQUIPMENT',
    NULL,
    '本社ビル 3階 備品室',
    '{"portable": true, "size": "large", "magnetic": true}'::jsonb,
    NULL,
    true
),
(
    '20000000-0000-0000-0000-000000000004'::uuid,
    'モバイルホワイトボード #2',
    'EQUIPMENT',
    NULL,
    '本社ビル 4階 備品室',
    '{"portable": true, "size": "large", "magnetic": true}'::jsonb,
    NULL,
    true
);

-- ビデオカメラ
INSERT INTO resources (id, name, type, capacity, location, equipment, required_role, is_active) VALUES
(
    '20000000-0000-0000-0000-000000000005'::uuid,
    'ビデオカメラ #1',
    'EQUIPMENT',
    NULL,
    '本社ビル 5階 AV機器室',
    '{"4k": true, "tripod": true, "microphone": true, "storage": "256GB"}'::jsonb,
    NULL,
    true
),
(
    '20000000-0000-0000-0000-000000000006'::uuid,
    'ビデオカメラ #2',
    'EQUIPMENT',
    NULL,
    '本社ビル 5階 AV機器室',
    '{"4k": true, "tripod": true, "microphone": true, "storage": "256GB"}'::jsonb,
    NULL,
    true
);

-- マイク・スピーカーシステム
INSERT INTO resources (id, name, type, capacity, location, equipment, required_role, is_active) VALUES
(
    '20000000-0000-0000-0000-000000000007'::uuid,
    'ワイヤレスマイクセット #1',
    'EQUIPMENT',
    NULL,
    '本社ビル 5階 AV機器室',
    '{"wireless": true, "microphones": 4, "receiver": true, "battery": "rechargeable"}'::jsonb,
    NULL,
    true
),
(
    '20000000-0000-0000-0000-000000000008'::uuid,
    'ポータブルスピーカー #1',
    'EQUIPMENT',
    NULL,
    '本社ビル 5階 AV機器室',
    '{"bluetooth": true, "battery": "10hours", "power": "50W"}'::jsonb,
    NULL,
    true
);

-- ノートPC（貸出用）
INSERT INTO resources (id, name, type, capacity, location, equipment, required_role, is_active) VALUES
(
    '20000000-0000-0000-0000-000000000009'::uuid,
    '貸出用ノートPC #1',
    'EQUIPMENT',
    NULL,
    '本社ビル IT管理室',
    '{"os": "Windows 11", "cpu": "Core i7", "ram": "16GB", "storage": "512GB SSD"}'::jsonb,
    NULL,
    true
),
(
    '20000000-0000-0000-0000-000000000010'::uuid,
    '貸出用ノートPC #2',
    'EQUIPMENT',
    NULL,
    '本社ビル IT管理室',
    '{"os": "Windows 11", "cpu": "Core i7", "ram": "16GB", "storage": "512GB SSD"}'::jsonb,
    NULL,
    true
);

-- 高額備品（アクセス制限あり）
INSERT INTO resources (id, name, type, capacity, location, equipment, required_role, is_active) VALUES
(
    '20000000-0000-0000-0000-000000000011'::uuid,
    '高性能ビデオ会議システム',
    'EQUIPMENT',
    NULL,
    '本社ビル 6階 AV機器室',
    '{"4k_camera": true, "ai_tracking": true, "noise_cancellation": true, "multi_display": true, "price": "high"}'::jsonb,
    'MANAGER',  -- マネージャー以上のみ予約可能
    true
);

-- ============================================================================
-- 完了メッセージ
-- ============================================================================
DO $$
DECLARE
    meeting_room_count INT;
    equipment_count INT;
BEGIN
    SELECT COUNT(*) INTO meeting_room_count FROM resources WHERE type = 'MEETING_ROOM';
    SELECT COUNT(*) INTO equipment_count FROM resources WHERE type = 'EQUIPMENT';
    
    RAISE NOTICE 'リソースシードデータ投入完了';
    RAISE NOTICE '会議室: % 件', meeting_room_count;
    RAISE NOTICE '備品: % 件', equipment_count;
    RAISE NOTICE '合計: % 件', meeting_room_count + equipment_count;
END $$;
