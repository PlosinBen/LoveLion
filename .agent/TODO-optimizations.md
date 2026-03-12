# Pending Optimizations

Code review 於 2026-03-12 產出，已完成項目 1-4，以下為待處理項目。

---

## 5. 前端狀態管理不一致（中優先）

同時使用三種方式管理狀態：
- `useState`（SSR 安全）— `useSpace.ts`
- 全域 `ref`（非 SSR 安全）— `useAuth.ts`, `useLoading.ts`
- Pinia store — `spaceDetail.ts`

**建議**：統一策略，認證和全域狀態遷移到 Pinia store，避免 SSR hydration 問題。

---

## 6. 前端大量使用 `alert()` / `confirm()`（中優先）

約 42 處使用瀏覽器原生對話框，手機體驗差。

**建議**：封裝統一的 modal/toast 元件取代原生對話框。

---

## 7. 列表 API 沒有分頁（中優先）

所有 list endpoint（transactions, stores, products, members）一次回傳全部資料。

**影響**：交易量大時效能下降。

**建議**：加入 cursor 或 offset/limit 分頁，前端配合 infinite scroll 或分頁元件。

---

## 8. 認證檢查分散在各頁面（中優先）

每個頁面 `onMounted` 各自呼叫 `initAuth()` 並手動 `router.push('/login')`。

**建議**：改用 Nuxt route middleware 統一處理認證檢查和跳轉。

---

## 9. `NewShortID` 碰撞處理邏輯有問題（中優先）

`backend/internal/utils/id.go` — 重試 3 次都用同樣長度 5，只在迴圈結束後才加長到 6。

**建議**：碰撞後立即遞增長度再重試。

---

## 10. 文件與實際不一致（低優先）

- `.agent/skills/api-conventions/SKILL.md` 仍使用 `ledger`/`trip` 術語
- `bin/seed` 用的是舊 API 路徑（`/trips`, `/ledgers`），與現在的 `/spaces` 路由不符

**建議**：統一更新文件和 seed script 為 Space 術語。

---

## 11. 前端交易表單 add/edit 大量重複程式碼（低優先）

`pages/spaces/[id]/transaction/add.vue` 和 `edit.vue` 有相似的表單邏輯。

**建議**：抽成共用 composable（如 `useTransactionForm`）。

---

## 12. 圖片上傳缺少大小限制（低優先）

`backend/internal/handlers/image.go` 沒有檢查上傳檔案大小。

**建議**：在 handler 或 middleware 加入 `Content-Length` 或檔案大小檢查（建議上限 10MB）。
