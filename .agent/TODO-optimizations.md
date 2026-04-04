# Pending Optimizations

Code review 於 2026-03-12 產出，Release readiness audit 於 2026-04-02 補充。
已完成項目 1-4、10，以下為待處理項目。

---

## P0 — CRITICAL（Release 前必須修）

### C1. R2 密鑰洩漏在 git 歷史

`.env-backend` 包含真實 Cloudflare R2 credentials 且被 git 追蹤。

**行動**：
1. 立即到 Cloudflare Dashboard 輪換所有 R2 access key
2. 把 `.env-backend` 加入 `.gitignore`
3. 用 `git filter-repo` 從 git 歷史中徹底移除該檔案
4. `.env-backend.example` 改為只放 placeholder 值

---

### C2. 缺少 CORS 設定

`backend/main.go` 沒有設定 CORS，production 下任何 origin 都能打 API。

**建議**：使用 `gin-contrib/cors`，從環境變數讀取 allowed origins。

---

### C3. 環境變數缺少 production 驗證

`backend/internal/config/config.go` 的 JWT_SECRET 預設值為 `"dev-secret-key"`，production 如果沒設環境變數會靜默使用弱密鑰。

**建議**：非 development 環境下，缺少 `JWT_SECRET`、`DATABASE_URL` 等必要變數時應直接 panic。

---

## P2 — MEDIUM

### M1. 缺少 SEO meta tags

`nuxt.config.ts` 缺少 description、og:title、og:image 等基本 meta。

---

### M2. Vite file polling 應限制在開發環境

`nuxt.config.ts` 的 `usePolling: true` 是 Docker 開發用的，production 不需要，會造成不必要的 filesystem overhead。

**建議**：加上 `process.env.NODE_ENV === 'development'` 條件判斷。

---

### M3. 後端缺少結構化 logging

使用 `fmt.Printf` / `log.Println` 做 log，缺少 request ID、log level 等結構化資訊。

**建議**：引入 `slog` 或 `logrus`，統一 log format。

---

### M4. API 缺少 Rate Limiting

所有 endpoint 沒有 rate limit，公開後容易被暴力攻擊。

**建議**：至少對 `/api/users/login`、`/api/users/register` 加上 rate limiting。

---

### M5. Database SSL mode 未強制

`config.go` 預設 `sslmode=disable`，production 應強制 `sslmode=require`。

---

### M6. 前端狀態管理不一致

同時使用三種方式管理狀態：
- `useState`（SSR 安全）— `useSpace.ts`
- 全域 `ref`（非 SSR 安全）— `useAuth.ts`, `useLoading.ts`
- Pinia store — `spaceDetail.ts`

**建議**：統一策略，認證和全域狀態遷移到 Pinia store，避免 SSR hydration 問題。

---

### M7. 前端大量使用 `alert()` / `confirm()`

約 42 處使用瀏覽器原生對話框，手機體驗差。

**建議**：封裝統一的 modal/toast 元件取代原生對話框。

---

### M8. 列表 API 沒有分頁

所有 list endpoint（transactions, stores, products, members）一次回傳全部資料。

**影響**：交易量大時效能下降。

**建議**：加入 cursor 或 offset/limit 分頁，前端配合 infinite scroll 或分頁元件。

---

### M9. 認證檢查分散在各頁面

每個頁面 `onMounted` 各自呼叫 `initAuth()` 並手動 `router.push('/login')`。

**建議**：改用 Nuxt route middleware 統一處理認證檢查和跳轉。

---

### M10. `NewShortID` 碰撞處理邏輯有問題

`backend/internal/utils/id.go` — 重試 3 次都用同樣長度 5，只在迴圈結束後才加長到 6。

**建議**：碰撞後立即遞增長度再重試。

---

## P3 — LOW

### L1. 前端交易表單 add/edit 大量重複程式碼

`pages/spaces/[id]/transaction/add.vue` 和 `edit.vue` 有相似的表單邏輯。

**建議**：抽成共用 composable（如 `useTransactionForm`）。

---

### L2. 圖片上傳缺少大小限制

`backend/internal/handlers/image.go` 沒有檢查上傳檔案大小。

**建議**：在 handler 或 middleware 加入檔案大小檢查（建議上限 10MB）。

---

### L3. 缺少 CI/CD Pipeline

沒有 GitHub Actions workflow，所有測試靠手動 `bin/integration_test` 執行。

**建議**：建立基本 workflow 執行 `go test ./...` 和 `npm test`。

---

### L4. 後端缺少 HTTP request logging middleware

沒有統一的 request log，難以追蹤 production 問題。

**建議**：加入 middleware 記錄 method、path、status code、duration、user ID。

---

### L5. JWT token 過期時間寫死

`auth.go` 中 token expiry 固定 7 天，無法透過環境變數調整。

**建議**：改為可設定。
