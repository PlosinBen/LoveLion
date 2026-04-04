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
