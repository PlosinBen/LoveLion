# GOTCHAS

實作過程中踩過的坑或非直覺的行為。每條包含：現象描述、根本原因、解法/注意事項。

---

## G1: AI Worker 重啟後重試狀態歸零

**現象**：Worker 重啟後，之前已經重試 2 次的 pending 交易會從 attempt 1 重新開始，指數退避也被重置。

**根本原因**：重試計數（`retries` map）和退避時間（`retryAt` map）存在記憶體中，不持久化到資料庫。`recoverStuck()` 只把 processing 重設為 pending，不恢復計數。

**注意事項**：這是可接受的取捨——worst case 是多重試幾次。若需要精確計數，需在 transactions 表加 `ai_retry_count` 欄位。

---

## G2: AI 取消與 Worker 的競態條件

**現象**：使用者按下「取消辨識」時，如果 Worker 正好在 processing，兩邊都在寫 ai_status。

**根本原因**：取消操作設定 `ai_status = NULL`，Worker 的 writeSuccess/writeFailure 使用 `WHERE ai_status = 'processing'` 條件。

**解法**：這已被正確處理。Worker 的 conditional UPDATE 在 ai_status 已不是 processing 時會 RowsAffected=0，等同 no-op。取消方先到先贏，Worker 的寫入靜默失敗。不需要額外處理。

---

## G3: 401 時前端使用 window.location.href 而非 Vue Router

**現象**：API 收到 401 後，頁面會整頁重新載入到 /login，而非 SPA 內跳轉。

**根本原因**：`useApi.ts` 中 401 處理使用 `window.location.href = '/login'`，這是刻意的硬跳轉。

**注意事項**：這確保所有 Pinia store 和 composable 狀態完全重置，避免殘留的認證狀態導致後續 API 呼叫異常。不要改成 `navigateTo('/login')`，那樣會保留舊狀態。

---

## G4: useApi 偵測 FormData 時跳過 Content-Type header

**現象**：手動設定 `Content-Type: multipart/form-data` 上傳圖片時，後端解析失敗。

**根本原因**：FormData 需要瀏覽器自動設定 boundary 參數在 Content-Type 中。手動設定會遺漏 boundary。

**解法**：`useApi.ts` 在偵測到 body 是 FormData 時，不設定 Content-Type header（`lines 40-42`），讓瀏覽器自動帶上正確的 `multipart/form-data; boundary=...`。不要手動加 Content-Type。

---

## G5: 收據圖片壓縮品質刻意設為 1.0

**現象**：收據圖片上傳後檔案很大，壓縮似乎沒有生效。

**根本原因**：`useImages.ts` 中 transaction 類型的圖片使用 `initialQuality: 1.0`（不壓縮品質），只限制尺寸到 1280px / 1MB。

**注意事項**：這是刻意的——收據需要保留文字清晰度供 Gemini OCR 辨識。空間封面圖才用 0.9 品質壓縮。不要統一調低品質。

---

## G6: ImageManager 在新增交易時使用緩衝模式

**現象**：在新增交易頁選了圖片後，圖片沒有立刻上傳到伺服器。

**根本原因**：新增交易頁的 ImageManager 設定 `instantUpload=false`，檔案緩衝在前端的 `pendingUploads` 陣列中，等到表單送出時才與交易資料一起透過 multipart 送出。

**注意事項**：這確保交易建立和圖片上傳是原子操作（在同一個 DB Transaction 中）。編輯頁面則可能使用 `instantUpload=true`。取得緩衝檔案用 `getBufferedFiles()`。

---

## G7: Space 切換時 spaceDetail store 全部清空

**現象**：切換空間後，原空間的交易資料消失，需要重新載入。

**根本原因**：`spaceDetail.ts` 的 `setSpaceId()` 會重置所有 state（transactions、stores、members 等）和 fetched flags。

**注意事項**：這是刻意的，防止顯示上一個空間的資料。代價是切換空間時必定觸發重新載入。如果需要快取多空間資料，需要改為 Map 結構（目前不需要）。

---

## G8: settled_amount 使用向上取整（Ceiling）

**現象**：外幣消費分帳時，各人分攤金額加總可能微幅超過帳單金額。

**根本原因**：`transaction_service.go` 的 `calcSettledAmount()` 使用 `math.Ceil()` 計算每人的 settled_amount：`ceil(amount / totalAmount * billingAmount)`。

**注意事項**：這是刻意的保護收款方設計——確保收款方不會因為捨入而少收錢。若要改為四捨五入或向下取整，需要同步更新前端的統計計算。

---

## G9: useLedger() 是 useSpace() 的向後相容別名

**現象**：程式碼中同時出現 `useLedger()` 和 `useSpace()`，看起來像兩套系統。

**根本原因**：專案從 Ledger 概念遷移到 Space 概念，`useLedger()` 只是一個屬性映射別名（allLedgers → allSpaces 等）。

**注意事項**：新程式碼一律使用 `useSpace()`。`useLedger()` 存在是為了已有頁面的漸進遷移，最終會移除。

---

## G10: 後端啟動順序 — migrate 在 server 前

**現象**：首次啟動或 schema 變更後，直接連 API 會遇到表不存在的錯誤。

**根本原因**：`docker-compose.yml` 中 backend 的 command 是 `go mod tidy && go run cmd/migrate/main.go && go run main.go`，遷移跑完才啟動伺服器。

**注意事項**：如果遷移失敗，伺服器不會啟動。檢查 `docker compose logs backend` 確認遷移狀態。`bin/refresh-database` 會額外在 seed 前再跑一次遷移確保一致。

---

## G11: CORS 環境變數為空時不啟用 CORS 中介層

**現象**：開發環境跨域請求正常（因為走 Nuxt 代理），但直連後端時 CORS 被擋。

**根本原因**：`main.go` 中只在 `CORS_ORIGINS` 非空時才加上 CORS 中介層。開發環境不需要 CORS 因為前端走 Nuxt proxy。

**注意事項**：生產環境必須設定 `CORS_ORIGINS`。如果前端直連後端（不走代理），開發環境也需要設定。

---

## G12: AI Rate Limit 設為 0 或負數時完全停用

**現象**：設定 `RECEIPT_EXTRACT_RATE_LIMIT_PER_DAY=0` 後，AI 辨識沒有任何限制。

**根本原因**：rate limiter 邏輯中，0 或負數視為「不檢查」而非「完全禁止」。

**注意事項**：要禁用 AI 功能應設定 `RECEIPT_EXTRACT_ENABLED=false`，不要用 rate limit 為 0 來達成。

---

## G13: Transaction ID 是 String 而非 UUID

**現象**：嘗試用 UUID 格式查詢交易時找不到。

**根本原因**：Transaction.ID 是 VARCHAR(21) 的短字串（由 `utils.NewShortID` 生成），不是 UUID。這是為了 URL 可讀性。

**注意事項**：Repository 和 Handler 中 transaction ID 的參數型別是 `string`，不要用 `uuid.UUID` 解析。新增遷移時 foreign key 需對應 VARCHAR(21)。

---

## G14: SpaceAccess 中介層會 preload Space.Images

**現象**：Handler 中存取 `space.Images` 時已經有資料，不需要額外查詢。

**根本原因**：`middleware/space.go` 在驗證成員權限時同時 preload Images（`WHERE entity_type='space'`）並呼叫 `PopulateCoverImage()`。

**注意事項**：這意味著 space 物件在進入 handler 前就帶有 Images。不要在 handler 中重複查詢 Images。新增 entity_type 時要檢查此 preload 條件。

---

## G15: 前端 useSpace 的 fetchSpaces 有快取邏輯

**現象**：呼叫 `fetchSpaces()` 多次但只看到一次 API 請求。

**根本原因**：`useSpace.ts` 的 `fetchSpaces()` 檢查 `allSpaces.value.length > 0` 就跳過請求，除非傳入 `force=true`。

**注意事項**：新增/刪除空間後必須呼叫 `fetchSpaces(true)` 強制重新載入。同理，`spaceDetail.ts` 的各 fetch 方法也有 `fetched[key]` 快取，mutation 後要呼叫 `invalidate('transactions')` 等。

---

## G16: Payment 的 query param 預填

**現象**：從統計頁點「建議付款」跳轉到新增頁時，表單自動填入了付款人/收款人/金額。

**根本原因**：`transaction/add.vue` 會讀取 URL query params `?type=payment&payer=...&payee=...&amount=...`，自動切換到付款模式並填入欄位。

**注意事項**：這個預填邏輯在頁面 mount 時執行。如果要新增其他預填欄位，需要同步更新統計頁的連結生成和 add.vue 的解析。

---

## G17: Vite 在 Docker 中使用 polling 偵測檔案變更

**現象**：在 Docker 環境下修改前端檔案後，HMR 有時延遲或不觸發。

**根本原因**：`nuxt.config.ts` 中非 production 環境下啟用 `vite.server.watch.usePolling = true`，因為 Docker volume mount 的 filesystem events 在某些平台不可靠。

**注意事項**：這會增加 CPU 使用。如果 HMR 有問題，先確認 polling 已開啟。macOS + Docker Desktop 通常需要 polling。
