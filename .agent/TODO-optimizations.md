# Pending Optimizations

最後更新：2026-04-04

---

## MEDIUM

### 1. 後端缺少結構化 logging

使用 `fmt.Printf` / `log.Println` 做 log，缺少 request ID、log level 等結構化資訊。

**建議**：引入 `slog`，統一 log format。

---

### 2. 前端狀態管理不一致

同時使用三種方式管理狀態：
- `useState`（SSR 安全）— `useSpace.ts`
- 全域 `ref`（非 SSR 安全）— `useAuth.ts`, `useLoading.ts`
- Pinia store — `spaceDetail.ts`

**建議**：統一策略，認證和全域狀態遷移到 Pinia store，避免 SSR hydration 問題。

---

### 3. 前端大量使用 `alert()` / `confirm()`

約 42 處使用瀏覽器原生對話框，手機體驗差。

**建議**：封裝統一的 modal/toast 元件取代原生對話框。

---

### 4. 列表 API 沒有分頁

所有 list endpoint（transactions, stores, products, members）一次回傳全部資料。

**影響**：交易量大時效能下降。

**建議**：加入 cursor 或 offset/limit 分頁，前端配合 infinite scroll 或分頁元件。

---

## LOW

### 5. 前端交易表單 add/edit 大量重複程式碼

`pages/spaces/[id]/transaction/add.vue` 和 `edit.vue` 有相似的表單邏輯。

**建議**：抽成共用 composable（如 `useTransactionForm`）。

---

### 6. 圖片上傳缺少大小限制

`backend/internal/handlers/image.go` 沒有檢查上傳檔案大小。

**建議**：在 handler 或 middleware 加入檔案大小檢查（建議上限 10MB）。

---

### 7. 缺少 CI/CD Pipeline

沒有 GitHub Actions workflow，所有測試靠手動 `bin/integration_test` 執行。

**建議**：建立基本 workflow 執行 `go test ./...` 和 `npm test`。

---

### 8. 後端缺少 HTTP request logging middleware

沒有統一的 request log，難以追蹤 production 問題。

**建議**：加入 middleware 記錄 method、path、status code、duration、user ID。

---

### 9. JWT token 過期時間寫死

`auth.go` 中 token expiry 固定 7 天，無法透過環境變數調整。

**建議**：改為可設定。
