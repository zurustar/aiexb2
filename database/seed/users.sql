-- database/seed/users.sql
-- ユーザーのシードデータ
-- 
-- 開発環境用のテストユーザーを作成
-- 本番環境では IdP から同期されるため、このファイルは使用しない

INSERT INTO users (
    id,
    sub,
    email,
    name,
    role,
    manager_id,
    penalty_score,
    penalty_score_expire_at,
    is_active
) VALUES
    (
        '11111111-1111-1111-1111-111111111111',
        'idp-admin-001',
        'admin@example.com',
        'System Admin',
        'ADMIN',
        NULL,
        0,
        NULL,
        true
    ),
    (
        '22222222-2222-2222-2222-222222222222',
        'idp-manager-001',
        'manager@example.com',
        'Manager Tanaka',
        'MANAGER',
        '11111111-1111-1111-1111-111111111111',
        0,
        NULL,
        true
    ),
    (
        '33333333-3333-3333-3333-333333333333',
        'idp-secretary-001',
        'secretary@example.com',
        'Secretary Sato',
        'SECRETARY',
        '22222222-2222-2222-2222-222222222222',
        0,
        NULL,
        true
    ),
    (
        '44444444-4444-4444-4444-444444444444',
        'idp-general-001',
        'general@example.com',
        'General Suzuki',
        'GENERAL',
        '22222222-2222-2222-2222-222222222222',
        1,
        NOW() + INTERVAL '30 days',
        true
    ),
    (
        '55555555-5555-5555-5555-555555555555',
        'idp-auditor-001',
        'auditor@example.com',
        'Auditor Kato',
        'AUDITOR',
        '11111111-1111-1111-1111-111111111111',
        0,
        NULL,
        true
    )
ON CONFLICT (sub) DO NOTHING;
