# Pending Optimizations

最後更新：2026-04-05

---

## LOW

### 1. 缺少 CI/CD Pipeline

沒有 GitHub Actions workflow，所有測試靠手動 `bin/integration_test` 執行。

**建議**：建立基本 workflow 執行 `go test ./...` 和 `npm test`。

---
