// frontend/src/app/layout.tsx
// ルートレイアウトコンポーネント
//
// 責務:
// - アプリケーション全体のレイアウト定義
// - グローバルスタイルの適用
// - メタデータの設定
// - プロバイダーのラップ

import './globals.css'

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  )
}
