# Pending Optimizations

最後更新：2026-04-05

---

## LOW

### 1. 前端交易表單 add/edit 大量重複程式碼

`pages/spaces/[id]/transaction/add.vue` 和 `edit.vue` 有相似的表單邏輯。

**建議**：抽成共用 composable（如 `useTransactionForm`）。

---

### 2. 圖片上傳缺少大小限制

`backend/internal/handlers/image.go` 沒有檢查上傳檔案大小。

**建議**：在 handler 或 middleware 加入檔案大小檢查（建議上限 10MB）。

---

### 3. 缺少 CI/CD Pipeline

沒有 GitHub Actions workflow，所有測試靠手動 `bin/integration_test` 執行。

**建議**：建立基本 workflow 執行 `go test ./...` 和 `npm test`。

---

### 4. 後端缺少 HTTP request logging middleware

沒有統一的 request log，難以追蹤 production 問題。

**建議**：加入 middleware 記錄 method、path、status code、duration、user ID。

---

### 5. JWT token 過期時間寫死

`auth.go` 中 token expiry 固定 7 天，無法透過環境變數調整。

**建議**：改為可設定。

---
