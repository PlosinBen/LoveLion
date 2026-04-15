# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 專案簡介

LoveLion 是一個共享記帳與比價應用程式。以「Space」作為統一概念（取代舊有的 Ledger/Trip 分離架構），將交易、商店與成員歸屬於同一空間。

## 技術棧

- **後端**: Go 1.25, Gin, GORM, PostgreSQL 18
- **前端**: Nuxt 4 (Vue 3), Tailwind CSS, TypeScript, Vitest
- **基礎設施**: 全部透過 Docker Compose 運行（宿主機不需安裝 Go/Node）

## 開發指令

所有指令都透過 Docker 執行，宿主機無 Go/Node 環境。

```bash
# 啟動開發環境
docker compose up -d

# 查看日誌
docker compose logs -f backend
docker compose logs -f frontend

# 執行後端測試（任何後端修改後必須執行）
docker compose exec backend go test ./...

# 執行單一測試
docker compose exec backend go test ./internal/handlers -run TestFunctionName

# 執行前端測試
docker compose exec frontend npm test

# 資料庫完整重置（清除 + 遷移 + 種子資料）
./bin/refresh-database
# 重置後測試帳號：dev / dev123

# 僅執行遷移
./bin/migrate

# 建立新的遷移檔
docker compose exec backend migrate create -ext sql -dir /app/migrations -seq <名稱>

# 模型變更後驗證 seed 可編譯
docker compose exec backend go build ./cmd/seed/

# 連線 PostgreSQL
docker compose exec postgres psql -U postgres -d lovelion
```

## 架構

### 後端結構 (`backend/`)
- `main.go` — 路由定義與伺服器啟動
- `internal/handlers/` — HTTP 處理器（auth, space, space_sharing, transaction, expense, payment, expense_template, comparison, image, announcement）
- `internal/models/` — GORM 模型（user, space, transaction, transaction_debt, comparison, image, announcement, expense_template）
- `internal/services/` — 業務邏輯（AI worker, AI extract, AI announcement, invite, transaction）
- `internal/repositories/` — 資料存取（transaction, expense, debt, member, invite）
- `internal/middleware/` — 認證（JWT）、空間權限、管理員權限、速率限制、請求日誌
- `internal/config/` — 環境設定載入
- `internal/database/` — 資料庫連線
- `internal/storage/` — Cloudflare R2 儲存
- `internal/utils/` — ID 生成工具
- `cmd/migrate/` — 遷移執行器
- `cmd/seed/` — 種子資料產生器
- `cmd/cleanup-r2/` — R2 清理工具
- `migrations/` — SQL 遷移檔（golang-migrate，依序編號）

### 前端結構 (`frontend/`)
- `pages/` — 檔案式路由。空間頁面在 `spaces/[id]/`，交易頁面在 `spaces/[id]/transaction/`
- `composables/` — `useApi`, `useAuth`, `useSpace`, `useImages`, `useTransactionForm`, `useExpenseTemplates`, `useLoading`, `useToast`, `useConfirm`, `usePrompt`, `useButtonStyle`
- `components/` — 共用元件（ImageManager, ExpenseForm, PaymentForm, DebtEditor, BroadcastBar, AnnouncementForm, BaseCard, BaseModal 等）
- `stores/` — Pinia 狀態管理（auth, spaceDetail, toast, confirm, loading, prompt）
- `types/` — TypeScript 型別定義
- API 代理：Nuxt 將 `/api/**` 代理至 `http://backend:8080/api/**`

### ID 策略
- **NanoID**（`utils.NewShortID`）：用於 URL 可見的實體（spaces, transactions, stores, invites, announcements）
- **UUID**：用於內部實體（users, members, items）

### 金額精度
- 金額：`DECIMAL(10,2)`
- 匯率：`DECIMAL(12,6)`

### API 慣例
- RESTful 風格，複數名詞，kebab-case（`/spaces`, `/transactions`）
- 巢狀資源：`/spaces/:id/transactions`, `/spaces/:id/stores`
- 成功回傳資源本身（無包裝）。錯誤回傳 `{"error": "訊息"}`
- POST 回傳 201，其餘回傳 200
- 所有變更端點需經 `SpaceAccess` 中介層；僅擁有者操作使用 `SpaceOwnerOnly`

### 安全規則
- PUT/POST/DELETE 必須驗證 `SpaceMember` 權限
- 僅 `owner` 可修改設定、管理成員、撤銷邀請
- 禁止在 API 回應中暴露 `PasswordHash`
- 多表更新必須使用資料庫交易（DB Transaction）

## 前端規則
- Tailwind：僅使用原生類名，禁止任意值（`-[...]`），禁止在 tailwind.config 自定義顏色
- Mobile-First 佈局，底部導航列

## 檔案編碼
- 所有檔案必須使用 **LF** 換行（嚴禁 CRLF）

## 補充文件
更多細節請參考 `.agent/` 目錄 — 包含檔案索引（PROJECT_MAP）、架構決策（DECISIONS）、踩坑紀錄（GOTCHAS）、發布流程（RELEASE）與功能規格（features/）。
