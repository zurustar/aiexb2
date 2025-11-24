# 監査ログ・バックアップ保持ポリシー

## 概要

本ドキュメントは、Enterprise Schedule Management System (ESMS) における監査ログとデータベースバックアップの保持ポリシーおよび運用手順を定義します。

## 監査ログ保持ポリシー

### 保持期間

- **オンライン検索用**: 1年間
  - CloudWatch Logs / Datadog等のログ管理サービスに保存
  - 即座に検索・分析可能な状態で保持
- **アーカイブ保存**: 追加2年間（合計3年間）
  - S3 Glacier等のコールドストレージへ移行
  - 必要時にリストア可能

### ログ形式

- **フォーマット**: JSON構造化ログ
- **必須フィールド**:
  - `timestamp`: ISO8601形式のタイムスタンプ
  - `user_id`: 操作実行ユーザーのID
  - `action`: 実行されたアクション（CREATE, UPDATE, DELETE等）
  - `resource_type`: 対象リソースタイプ（Reservation, Resource等）
  - `resource_id`: 対象リソースのID
  - `ip_address`: クライアントIPアドレス
  - `user_agent`: クライアントUser-Agent
  - `signature_hash`: 改ざん検知用HMAC-SHA256ハッシュ

### 改ざん防止

- **署名ハッシュ**: 各ログレコードにHMAC-SHA256署名を付与
- **WORM (Write Once Read Many)**: S3 Object Lockを使用して上書き・削除を防止
- **SIEM連携**: リアルタイムで外部SIEMへ転送し、独立した監査証跡を保持

### S3 Lifecycle設定例

```json
{
  "Rules": [
    {
      "Id": "AuditLogRetentionPolicy",
      "Status": "Enabled",
      "Filter": {
        "Prefix": "audit-logs/"
      },
      "Transitions": [
        {
          "Days": 365,
          "StorageClass": "GLACIER"
        }
      ],
      "Expiration": {
        "Days": 1095
      }
    }
  ]
}
```

## データベースバックアップポリシー

### RPO/RTO目標

- **RPO (Recovery Point Objective)**: 15分
  - 最大15分間のデータ損失を許容
- **RTO (Recovery Time Objective)**: 1時間
  - 障害発生から1時間以内にサービス復旧

### バックアップ戦略

#### 1. 自動スナップショット

- **頻度**: 毎日1回（深夜2:00 JST）
- **保持期間**: 7日間
- **対象**: RDS全体（全データベース）

#### 2. Point-in-Time Recovery (PITR)

- **有効化**: 必須
- **保持期間**: 7日間
- **粒度**: 5分間隔
- **用途**: 特定時点へのリカバリ

#### 3. クロスリージョンレプリケーション

- **プライマリリージョン**: ap-northeast-1 (東京)
- **セカンダリリージョン**: ap-northeast-3 (大阪)
- **レプリケーション遅延**: 5分以内
- **用途**: リージョン障害時のDR

### バックアップ手順

#### RDS自動バックアップ設定

```bash
# RDS自動バックアップの有効化
aws rds modify-db-instance \
  --db-instance-identifier esms-prod \
  --backup-retention-period 7 \
  --preferred-backup-window "17:00-18:00" \
  --apply-immediately

# PITR有効化確認
aws rds describe-db-instances \
  --db-instance-identifier esms-prod \
  --query 'DBInstances[0].BackupRetentionPeriod'
```

#### 手動スナップショット作成

```bash
# 重要な変更前に手動スナップショット作成
aws rds create-db-snapshot \
  --db-instance-identifier esms-prod \
  --db-snapshot-identifier esms-prod-manual-$(date +%Y%m%d-%H%M%S)
```

### リストア手順

#### PITRによるリストア

```bash
# 特定時点へのリストア
aws rds restore-db-instance-to-point-in-time \
  --source-db-instance-identifier esms-prod \
  --target-db-instance-identifier esms-prod-restored \
  --restore-time "2025-11-24T10:00:00Z"
```

#### スナップショットからのリストア

```bash
# スナップショットからのリストア
aws rds restore-db-instance-from-db-snapshot \
  --db-instance-identifier esms-prod-restored \
  --db-snapshot-identifier esms-prod-snapshot-20251124
```

### 復旧演習

- **頻度**: 四半期ごと（3ヶ月に1回）
- **内容**:
  1. 最新スナップショットからのリストア
  2. アプリケーション接続確認
  3. データ整合性検証
  4. 復旧時間計測（RTO達成確認）
- **記録**: 演習結果を `docs/ops/dr-drill-reports/` に保存

## PII（個人情報）の取り扱い

### ログ出力時のマスキング

- **対象フィールド**: email, password, token, secret
- **マスキング方法**: `***MASKED***` に置換
- **実装**: `backend/internal/util/logger.go` の `maskPII()` 関数

### データベース暗号化

- **保存時暗号化**: RDS暗号化を有効化（AWS KMS使用）
- **転送時暗号化**: TLS 1.2以上を強制

## モニタリング・アラート

### 監査ログ

- **欠損検知**: 1時間以上ログが記録されない場合にアラート
- **異常パターン**: 短時間での大量削除操作を検知

### バックアップ

- **失敗検知**: 自動バックアップ失敗時に即座にアラート
- **容量監視**: バックアップストレージ使用率が80%を超えた場合にアラート

## 責任者

- **監査ログ管理**: セキュリティチーム
- **バックアップ管理**: インフラチーム
- **復旧演習実施**: SREチーム

## 参考資料

- [AWS RDS Backup and Restore](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_CommonTasks.BackupRestore.html)
- [S3 Object Lock](https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-lock.html)
- [CloudWatch Logs Retention](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/Working-with-log-groups-and-streams.html)
