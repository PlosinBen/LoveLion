# DECISIONS

重要的架構決策紀錄。每條包含：決策內容、原因、不要改動的理由。

---

## D1: Space 統一概念取代 Ledger/Trip 分離架構

**決策**：以 Space（空間）作為唯一的頂層容器，透過 `type` 欄位（personal/trip/group）區分用途，不再分開維護 Ledger 和 Trip 兩套模型。

**原因**：Ledger 和 Trip 共享 90% 以上的邏輯（交易、成員、邀請），分開維護導致大量重複程式碼與不一致的 bug。統一後路由、權限、UI 全部收斂為一套。

**不要改動**：前端仍保留 `useLedger()` 向後相容別名（`composables/useSpace.ts`），若要移除需全面搜尋確認無殘留引用。

---

## D2: 雙 ID 策略 — NanoID（URL）+ UUID（內部）

**決策**：URL 暴露的實體（Transaction、Store、Invite）使用自製短 ID（`utils.NewShortID`，a-zA-Z0-9，起始 5 字元），內部實體（User、SpaceMember、Item）使用 UUID。

**原因**：短 ID 讓 URL 更簡潔易讀（`/transactions/aB3kZ` vs UUID），同時避免暴露自增 ID 的可猜測性。UUID 用於不出現在 URL 的內部關聯，提供唯一性保證。

**不要改動**：短 ID 的碰撞處理邏輯（`id.go` 中碰撞時自動加長 +1 字元，最多重試 3 次）是刻意設計。Transaction.ID 欄位為 VARCHAR(21) 而非 UUID，改動需同步遷移與所有 Repository。

---

## D3: JSONB 儲存空間設定（categories / currencies / paymentMethods / splitMembers）

**決策**：Space 的 categories、currencies、payment_methods、split_members 使用 PostgreSQL JSONB 陣列儲存，而非獨立關聯表。

**原因**：這些設定是每個空間獨立的短列表（通常 < 20 項），不需要跨空間查詢或外鍵約束。JSONB 避免了大量 join 和多對多表的複雜性，讀寫都在單一 row 完成。

**不要改動**：前端的預設值（`useTransactionForm.ts` 中的 DEFAULT_CATEGORIES、DEFAULT_CURRENCIES）必須與後端 seed 資料保持一致。

---

## D4: 全 Docker 開發環境，宿主機不裝 Go/Node

**決策**：開發環境完全透過 Docker Compose 運行，所有指令都經由 `docker compose exec` 執行。

**原因**：確保所有開發者（含 AI agent）的環境一致，避免 Go/Node 版本差異問題。backend container 啟動時自動 `go mod tidy && migrate && run`，frontend 自動 `npm install && dev`。

**不要改動**：所有 bin/ 腳本都假設 Docker 環境。不要在指令文件中加入本機 go/npm 指令。

---

## D5: 多態圖片表（polymorphic Image）

**決策**：使用單一 images 表，透過 `entity_type`（space/transaction）+ `entity_id` 區分歸屬，而非每種實體各建圖片表。

**原因**：圖片的上傳/刪除/排序邏輯完全相同，多態設計讓前後端只需維護一套 CRUD。ImageManager 元件透過 props 切換 entity_type 即可複用。

**不要改動**：SpaceAccess 中介層已在 preload 時加上 `entity_type='space'` 過濾條件（`middleware/space.go`），新增 entity_type 時需同步更新。

---

## D6: AI Worker 採用 Polling（非 Queue）

**決策**：AI 收據辨識使用單一 goroutine 每 10 秒 polling `ai_status='pending'` 的 rows，而非 message queue（如 Redis/RabbitMQ）。

**原因**：目前用戶規模不大，polling 足以應付。避免引入額外基礎設施依賴（Redis/RabbitMQ），降低部署複雜度。PostgreSQL 本身就是可靠的狀態儲存。

**不要改動**：Worker 使用 conditional UPDATE（`WHERE ai_status='pending'`）作為樂觀鎖，確保多 worker 不會重複處理同一筆。即使目前只有單一 worker，此設計為未來水平擴展預留了基礎。

---

## D7: AI 狀態機設計（NULL → pending → processing → completed/failed）

**決策**：ai_status 使用 NULLABLE 欄位，NULL 表示「從未啟用 AI」，與 failed 明確區分。

**原因**：NULL 語意上表示「不適用」，而非「失敗」。前端可以根據 NULL vs failed 顯示不同 UI（無 badge vs 失敗重試 badge）。partial index `WHERE ai_status IS NOT NULL` 確保 polling 只掃描有意義的 rows。

**不要改動**：Update 時的狀態轉換規則（`transaction_service.go`）：failed + AIExtract=true → pending（重試）、failed + AIExtract=false → NULL（手動編輯覆蓋）。這是刻意的雙向轉換。

---

## D8: R2 上傳在 DB Transaction 內執行

**決策**：圖片上傳到 Cloudflare R2 的操作放在資料庫 transaction 內部，rollback 時用 background context 手動清理 R2 物件。

**原因**：確保 DB record 和 R2 object 的一致性。若 DB 寫入失敗，已上傳的 R2 物件會被清理。使用 background context 是因為 request context 可能已經取消。

**不要改動**：`uploadedKeys` 追蹤機制（`transaction_service.go`）是 rollback cleanup 的關鍵。cleanup 使用 `context.Background()` 而非 request context，是刻意設計。

---

## D9: 前端 SSR 關閉（SPA 模式）

**決策**：`nuxt.config.ts` 設定 `ssr: false`，全部客戶端渲染。

**原因**：此 App 是登入後才能使用的工具型應用，不需要 SEO。SPA 模式簡化狀態管理（不需處理 hydration mismatch），也降低伺服器負載。

**不要改動**：若開啟 SSR，需重新審視所有 localStorage 存取（auth store、useApi）和 window 物件使用，影響範圍大。

---

## D10: 金額使用 shopspring/decimal，不用 float

**決策**：所有金額相關欄位使用 `shopspring/decimal.Decimal`，資料庫為 `DECIMAL(10,2)`，匯率為 `DECIMAL(12,6)`。

**原因**：浮點數的精度問題在金融計算中會造成分帳不準確。decimal 提供精確的十進位運算，避免 0.1 + 0.2 != 0.3 的問題。

**不要改動**：前端傳入的金額字串會由 GORM 自動解析為 decimal。settled_amount 的 ceiling 捨入（`math.Ceil`）是刻意向上取整，保護收款方。

---

## D11: Nuxt API 代理而非直連後端

**決策**：前端透過 Nuxt 的 `routeRules` 將 `/api/**` 代理至 `http://backend:8080/api/**`，而非讓瀏覽器直連後端。

**原因**：避免 CORS 問題、隱藏後端實際位址、統一前端的 API base path。`apiBase` 預設為空字串（相對路徑），觸發 Nuxt 代理。

**不要改動**：生產環境可透過 `NUXT_PUBLIC_API_BASE` 環境變數指向實際後端 URL，但開發環境必須保持空字串走代理。

---

## D12: Tag-based 自動部署

**決策**：使用 git tag（如 `v1.0.0`）觸發部署，cron 每 5 分鐘檢查新 tag。

**原因**：簡單且可靠。不依賴 CI/CD 平台，適合小團隊。tag 同時作為版本號和部署觸發器。

**不要改動**：`docker-compose.prod.yml` 中 migration container 會在啟動時自動執行遷移，無需手動介入。

---

## D13: 測試架構 — httptest 整合測試為主

**決策**：後端測試以 httptest 整合測試為主（`handlers/*_test.go`），透過 `testutil.TestDB(t)` 建立獨立測試資料庫，每個測試函式完整走 HTTP → Handler → Service → DB。

**原因**：整合測試最接近實際行為，能捕捉到中介層串接、JSON 序列化、資料庫 constraint 等 mock 無法覆蓋的問題。

**不要改動**：testutil 中的 `TestRouter()` 已正確串接所有中介層。測試使用真實 DB 而非 mock，這是刻意的設計選擇。
