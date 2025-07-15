# Golang Backend for Google Calendar Integration

Docker Compose対応のGolangバックエンドAPI

## 🐳 Docker Compose で起動

```bash
# バックエンドディレクトリに移動
cd backend

# Docker Composeでビルド・起動
docker-compose up --build

# バックグラウンドで起動
docker-compose up -d --build

# 停止
docker-compose down
```

## 📋 API エンドポイント

### ヘルスチェック
- **GET** `http://localhost:8080/health`
- レスポンス: `{"status": "ok"}`

### Googleカレンダー取得
- **GET** `http://localhost:8080/api/calendar`
- **Headers**: `Authorization: Bearer {access_token}`
- **機能**: Google Calendar APIから最新10件のイベントを取得
- **レスポンス**: 
```json
{
  "events": [
    {
      "id": "event_id",
      "summary": "イベント名",
      "start": {...},
      "end": {...}
    }
  ],
  "count": 10
}
```

## 🔧 環境構成

- **ポート**: 8080
- **CORS**: `http://localhost:3000` を許可
- **ヘルスチェック**: 30秒間隔で `/health` エンドポイントを監視

## 📝 使用方法

1. フロントエンド（Next.js）からGoogle OAuthでアクセストークンを取得
2. そのトークンを `Authorization: Bearer {token}` ヘッダーに設定
3. `/api/calendar` エンドポイントにリクエスト送信
4. Google Calendar APIから取得したデータが返される