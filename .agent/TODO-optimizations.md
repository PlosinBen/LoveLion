# Pending Optimizations

最後更新：2026-04-05

---

## FEATURE

### 1. 客製化報表

使用者可以自行組合公式，產出個人化的消費報表。

**目標**：
- 使用者可定義報表欄位（如：分類加總、成員消費比例、時間區間篩選）
- 支援自訂公式（如：`餐飲 + 交通` 合併計算、`總支出 - 付款` 淨額）
- 報表可儲存、重複使用
- 視覺化呈現（圖表 / 表格）

---

## LOW

### 2. 缺少 CI/CD Pipeline

沒有 GitHub Actions workflow，所有測試靠手動 `bin/integration_test` 執行。

**建議**：建立基本 workflow 執行 `go test ./...` 和 `npm test`。

---
