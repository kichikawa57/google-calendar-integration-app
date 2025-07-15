# Google Calendar Integration App

A full-stack application that integrates with Google Calendar API using Next.js frontend and Golang backend.

## Features

- Google OAuth authentication
- Access token management
- Google Calendar API integration
- Recent calendar events display

## Setup

### Slack
- アプリ一覧ページ
  - https://api.slack.com/apps
- ngrokを使ってhttpsで起動できるようにする
- Slackのアプリを作成する
- Incoming Webhooksをonにする
- Manage DistributionをPublic Distributionにする
- Redirect URLsも設定とOAuth Tokensをonにする
- Client IDとClient Secretを取得する

### 1. Google Cloud Console Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google Calendar API
4. Create OAuth 2.0 credentials:
   - APIとサービスの認証情報にアクセス
   - Application type: Web application
   - Authorized redirect URIs: `http://localhost:3000/api/auth/callback/google`
   - Note down the Client ID and Client Secret

### 2. Frontend Setup

```bash
cd frontend
npm install
```

Create `.env.local` file:
```
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-secret-key-here
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
```

### 3. Backend Setup

```bash
cd backend
go mod tidy
```

### 4. Running the Application

Terminal 1 (Frontend):
```bash
cd frontend
npm run dev
```

Terminal 2 (Backend):
```bash
cd backend
go run main.go
```

## Usage

1. Open http://localhost:3000
2. Click "Sign in with Google"
3. Authorize the application
4. Click "Fetch Calendar Data" to retrieve recent calendar events

## API Endpoints

- `GET /api/calendar` - Fetch calendar events (requires Bearer token)

## Tech Stack

- **Frontend**: Next.js, NextAuth.js, TypeScript, Tailwind CSS
- **Backend**: Go, Gin, Google Calendar API
- **Authentication**: Google OAuth 2.0