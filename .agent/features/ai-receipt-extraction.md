# AI 發票辨識自動填入

> 透過 LLM Vision 自動辨識上傳的發票圖片，背景填入記帳的「日期」與「品項」，節省整理消費的時間。

## 目標與範圍

**目標**：使用者上傳發票照片後，不需等待 LLM 處理，立刻可以繼續其他操作；背景 worker 完成辨識後自動補上日期與品項。

**只擷取以下欄位**：
- `date` — 結帳日期（純日期，不含時分）
- `items[]` — 品項清單，每筆 `{ name, unit_price, quantity }`

**不處理**的欄位（仍由使用者手動填寫）：標題、分類、付款方式、總額、幣別、分帳、Google Maps 連結、note。

**整單折扣**：在 prompt 中要求 LLM 加入一筆 `name="折扣"`、`unit_price` 為負數的品項，不動 schema。

---

## 使用者流程

### 建立流程
1. 使用者進入「新增交易」頁
2. 填寫 title（必填）
3. 上傳一張發票圖片
4. 勾選「使用 AI 辨識品項與日期」checkbox
5. 勾選後：
   - 其他欄位（date、items、total_amount 等）全部 disable
   - 顯示提示：「將在背景處理，可能需要約 1 分鐘」
   - 必填驗證放寬，只剩 title + image
6. 按儲存 → 立刻回到帳本列表
7. 該筆 transaction 在列表上顯示「🤖 AI 處理中」標記
8. 過一段時間後（10 秒～1 分鐘），手動下拉刷新看到 date 與 items 已填好

### 取消流程
- 在 transaction 詳細頁顯示 banner「🤖 AI 辨識中...」+「取消辨識」按鈕
- 點取消 → `ai_status` 回到 `NULL`，整筆變成一般 transaction，可正常編輯
- 若 worker 已在處理中，取消會讓 worker 寫回時自動 discard 結果

### 失敗流程
- 列表項目顯示 ⚠️ icon，tooltip 顯示錯誤訊息
- 詳細頁顯示 banner「⚠️ AI 辨識失敗：{error}」（純資訊，**沒有任何按鈕**）
- 表單可正常編輯，且額外顯示 AI checkbox（預設未勾選）
- 使用者按儲存後，PUT 端點自動處理 status 轉換：
  - **未勾選 AI** → `ai_status=NULL, ai_error=NULL`，banner 與 icon 消失，變成一般 transaction
  - **勾選 AI** → `ai_status='pending', ai_error=NULL`，worker 下次輪詢會重新處理（沿用既有的 image）
- 換句話說：**failed 狀態只要編輯過就不會留著**，要不變 NULL（手動處理）、要不變 pending（重跑 AI）
- 不需要任何專門的「重試」「清除」按鈕，全部收斂到既有的編輯儲存流程

---

## 狀態機

```
                ┌─────────────────────────────────────────┐
                │                                         │
        [建立 with AI]                                    │
                ↓                                         │
            pending ──[使用者 cancel]──→ NULL（一般交易） │
                │                                         │
        [worker 認領]                                     │
                ↓                                         │
           processing ──[使用者 cancel]──→ NULL          │
                │             （worker 寫回時 discard）   │
        [worker LLM 完成]                                 │
                ↓                                         │
        completed / failed                                │
                │                                         │
        [使用者編輯儲存（PUT）on failed]                  │
        ├─ 未勾 AI ─→ NULL ──────────────────────────────┘
        └─ 勾 AI    ─→ pending（worker 重跑）─────────────┘
```

`ai_status` 可能值：
- `NULL` — 未使用 AI（一般 transaction）或已取消
- `pending` — 等待 worker 認領
- `processing` — worker 正在呼叫 LLM
- `completed` — 處理完成，date / items 已寫入
- `failed` — 處理失敗，`ai_error` 有錯誤訊息

---

## 資料庫變更

### Migration

```sql
-- up
ALTER TABLE transactions
  ADD COLUMN ai_status varchar(20),
  ADD COLUMN ai_error text;

CREATE INDEX idx_transactions_ai_status
  ON transactions(ai_status)
  WHERE ai_status IS NOT NULL;

-- down
DROP INDEX IF EXISTS idx_transactions_ai_status;
ALTER TABLE transactions
  DROP COLUMN IF EXISTS ai_error,
  DROP COLUMN IF EXISTS ai_status;
```

### Model 變更
`internal/models/transaction.go` 的 `Transaction` 加：
```go
AIStatus *string `gorm:"type:varchar(20);column:ai_status" json:"ai_status,omitempty"`
AIError  string  `gorm:"type:text;column:ai_error" json:"ai_error,omitempty"`
```

---

## API 設計

### 修改既有端點：`POST /api/spaces/:spaceId/expenses`

把現有「先建 expense → 再用 `/api/images` 上傳」兩步驟合成單一 atomic 端點。

**Content-Type 行為**：
- `application/json`：保留既有純 JSON 行為（無 images、無 AI），向後相容
- `multipart/form-data`：新模式，可一次送出 expense 資料 + N 張圖片 + 選用 AI 旗標

**multipart fields**：
- `data` (required)：JSON 字串，內容跟既有 JSON payload 完全相同，可額外帶 `ai_extract: true`
- `images` (optional, repeatable)：圖片檔案，可有多張；jpg/jpeg/png/webp，每張 ≤ 5MB

**處理流程**（包在單一 `db.Transaction` 內）：
1. Parse multipart，解析 `data` JSON
2. 驗證：若 `ai_extract=true`，必須至少有一張 image；驗證放寬（total_amount 可 0、items 可空）
3. 建立 transaction、expense、items、debts（沿用既有 service 邏輯）
4. 依序上傳每張圖片到 R2 + 建立 image record（`entity_type='transaction'`, `entity_id=新 txn id`，`sort_order` 依上傳順序）
5. 若 `ai_extract=true`，將 transaction 的 `ai_status` 設為 `'pending'`
6. 任一步驟失敗 → rollback DB + 刪除已上傳的 R2 物件（沿用既有 `image.go:163` 的 cleanup pattern）

**回傳**：`201 Created` + Transaction 物件（含 `ai_status` 欄位）

**取代既有的 frontend buffered 流程**：原本 `add.vue` 透過 `ImageManager` 的 `commit(txn_id)` 二段上傳會被改寫成一次 multipart submit。`/api/images` 端點仍保留（給 edit 頁與其他 entity type 使用）。

### 修改既有端點：`PUT /api/spaces/:spaceId/expenses/:txnId`

不需要支援 multipart（edit 頁的圖片管理仍走既有 `/api/images`）。
只需在 JSON payload 加 optional `ai_extract: bool` 欄位，後端依下方「PUT 自動轉換」規則處理 ai_status。

### 新增端點

#### `POST /api/spaces/:spaceId/transactions/:txnId/ai-cancel`
取消處理中的 AI 辨識。**只處理 pending / processing**，failed 狀態交給 PUT 自動處理。
```sql
UPDATE transactions SET ai_status=NULL, ai_error=NULL
WHERE id=? AND space_id=? AND ai_status IN ('pending','processing')
```
- 用途：使用者在表單 disable 時主動取消（pending/processing 不能走 PUT，因為表單是 disable 的）
- RowsAffected=0 → 回 409
- RowsAffected=1 → 回 200
- 對處理中的 worker：worker 寫回時的 `WHERE ai_status='processing'` 條件會 RowsAffected=0，自動 discard 結果

### 其他既有端點調整

- `GET /api/spaces/:id/transactions` — 回傳新增 `ai_status`、`ai_error` 欄位
- `GET /api/spaces/:id/transactions/:txnId` — 同上
- `POST /api/spaces/:id/payments` — 不變
- `PUT /api/spaces/:id/payments/:txnId` — 不變
- `POST /api/images` — 不變（edit 頁與其他 entity type 仍會用）

#### PUT 端點的 ai_status 自動轉換

PUT 的 payload 新增 optional 欄位 `ai_extract: bool`。後端在 update 既有欄位的同一個 `db.Transaction` 內，依下列規則寫 ai_status：

| 當前 ai_status | payload `ai_extract` | 寫回 ai_status | 寫回 ai_error | items 處理 |
|---|---|---|---|---|
| `failed` | `false` / 未提供 | `NULL` | `NULL` | 依 payload 中的 items 正常更新 |
| `failed` | `true` | `pending` | `NULL` | **清空 items**（worker 會重新填） |
| `pending` / `processing` | 任意 | **拒絕，回 409** | — | — |
| `NULL` / `completed` | 任意 | 不動 | 不動 | 依 payload 中的 items 正常更新 |

要點：
- `pending` / `processing` 狀態不允許 PUT（表單是 disable 的，前端不該送）
- 只有 `failed` 狀態會被 PUT 自動轉換，這是 V1 的範圍
- `NULL` 與 `completed` 的 transaction 沒有 AI checkbox，PUT 不會碰 ai_status
- 當 `ai_extract=true` 重跑時，前端應該 disable 並不送出 date / items（跟 add 流程一致），後端清空 items 表並把 date 留給 worker 寫回

---

## Worker 設計

### 結構
位置：`internal/services/ai_worker.go`

- 單一 long-running goroutine（不做併發）
- 在 `main.go` 啟動時 spawn，傳入 root context
- 輪詢間隔：**10 秒**
- 訂閱 root context，收到 cancel 時立刻終止 LLM 呼叫並退出

### 啟動 recovery
Worker 啟動時，先執行：
```sql
UPDATE transactions SET ai_status='pending'
WHERE ai_status='processing'
```
把上次換版時卡住的 row 重置回 pending。因為是 single instance worker，啟動時不可能有其他人在處理。

### 主迴圈

```go
for {
    select {
    case <-ctx.Done():
        return
    case <-time.After(10 * time.Second):
    }

    txns := pickPending(ctx, limit=5)
    for _, txn := range txns {
        if ctx.Err() != nil { return }
        processOne(ctx, txn)
    }
}
```

### 三階段處理（不持鎖跨 LLM）

#### 階段 1：認領（短 UPDATE，無鎖）
```sql
UPDATE transactions SET ai_status='processing'
WHERE id=? AND ai_status='pending'
```
- `RowsAffected=0` → 已被取消，跳過
- `RowsAffected=1` → 認領成功，繼續

#### 階段 2：呼叫 LLM（DB 完全無關）
- 從 transaction 取第一張圖（`ORDER BY sort_order ASC LIMIT 1`，`entity_type='transaction'`）
- 若 transaction 沒有任何 image → 直接進入「失敗寫回」，`ai_error="無圖片"`
- 從 R2 下載圖片內容
- `context.WithTimeout(ctx, 30*time.Second)` 包住 HTTP 呼叫
- 呼叫 Gemini API，帶 `responseSchema` 強制結構化輸出
- 收到 shutdown 信號 → context cancel → HTTP 立即中斷 → 函式返回，**不**寫回 DB
- 失敗（網路、API error、解析失敗）→ 進入「失敗寫回」

#### 階段 3：寫回（短 db.Transaction）
沿用既有 service pattern（參考 `transaction_service.go:269` 的 `CreateExpense`）：

**成功寫回**：
```go
db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
    result := tx.Model(&Transaction{}).
        Where("id = ? AND ai_status = ?", txnID, "processing").
        Updates(map[string]interface{}{
            "ai_status": "completed",
            "date":      parsedDate,
        })
    if result.Error != nil { return result.Error }
    if result.RowsAffected == 0 { return nil } // 被取消，整筆 rollback

    items := buildExpenseItems(expenseID, parsedItems) // 沿用既有 helper
    if len(items) > 0 {
        if err := tx.Create(&items).Error; err != nil { return err }
    }
    // 順便更新 total_amount
    return tx.Model(&Transaction{}).Where("id=?", txnID).
        Update("total_amount", calcTotal(items)).Error
})
```

**失敗寫回**：
```go
db.Model(&Transaction{}).
    Where("id = ? AND ai_status = ?", txnID, "processing").
    Updates(map[string]interface{}{
        "ai_status": "failed",
        "ai_error":  errMsg,
    })
```

### 為什麼這樣設計

| 設計選擇 | 理由 |
|---|---|
| Single goroutine | 避免並發大量 LLM 呼叫被 API rate limit / 收費 |
| 輪詢而非直接 dispatch | 容器重啟後自動恢復，不需特別 startup 邏輯 |
| 無 row lock，純 WHERE 條件 | LLM 呼叫期間 DB 完全無連線占用，換版安全 |
| 啟動時 recovery | 補上 LLM 呼叫期間被砍掉留下的孤兒 |
| Context cancel | Graceful shutdown 時立刻中斷 LLM 呼叫，不用等 timeout |

---

## LLM Provider

**選定**：Google Gemini 2.5 Flash

理由：
- 有持續性免費 tier（AI Studio 約 15 RPM、1500 RPD）
- Vision 品質對中文發票良好
- 支援 `responseSchema` 強制 JSON 結構化輸出
- 對個人專案規模綽綽有餘

### 抽象介面
`internal/services/ai_extract.go` 定義 interface，方便日後切換 provider：

```go
type ReceiptData struct {
    Date  *time.Time
    Items []ReceiptItem
}

type ReceiptItem struct {
    Name      string
    UnitPrice decimal.Decimal
    Quantity  decimal.Decimal
}

type ReceiptExtractor interface {
    Extract(ctx context.Context, image []byte, mimeType string) (*ReceiptData, error)
}
```

實作：`GeminiReceiptExtractor`

### Prompt（在 service 內固定）

System / User instruction：
```
你是發票辨識助手。請從圖片擷取消費資訊，並以指定的 JSON Schema 回傳。

規則：
- date 取發票上的消費日期，格式 YYYY-MM-DD。若無法判讀填 null。
- items 依發票上的順序列出。
- unit_price 是單價（不是小計）。quantity 是數量。
- 若有整單折扣，加一筆 name="折扣" 的品項，unit_price 為負數，quantity=1。
- 不要輸出 total，呼叫端會自行計算。
```

Response Schema：
```json
{
  "type": "object",
  "properties": {
    "date": { "type": "string", "nullable": true },
    "items": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "name": { "type": "string" },
          "unit_price": { "type": "number" },
          "quantity": { "type": "number" }
        },
        "required": ["name", "unit_price", "quantity"]
      }
    }
  },
  "required": ["items"]
}
```

### 環境變數
新增到 `.env-backend.example`：
```
GEMINI_API_KEY=
GEMINI_MODEL=gemini-2.5-flash
RECEIPT_EXTRACT_ENABLED=true
RECEIPT_EXTRACT_RATE_LIMIT_PER_DAY=20
```

`internal/config/config.go` 加對應欄位。

### 串接細節（不引入 SDK）

**HTTP direct，不使用 google.golang.org/genai SDK**：
- 端點：`POST https://generativelanguage.googleapis.com/v1beta/models/{model}:generateContent?key={api_key}`
- 直接用 `net/http`，手刻 request / response struct，整個 extractor 約 150 行
- 理由：避免引入大型相依、httptest mock 容易、行為可控

**圖片以 inline base64 傳送（不用 File API、不用 presigned URL）**：

Gemini `generateContent` 接受圖片的方式只有兩種：
1. `inlineData` — base64 字串塞在 request body
2. `fileData` — 引用 Google File API 上傳的 fileUri（**不接受任意 R2/S3 URL**）

R2 presigned URL 對 Gemini 沒用，Google 不會主動 fetch。File API 適合「同一張圖多次 call LLM」可省重複上傳，但我們是一張收據 call 一次，反而多一次 round trip + 等待 ACTIVE 狀態。所以選 inline base64：

```
R2 GetObject (io.ReadCloser)
   → base64.NewEncoder(stdEncoding, &requestBodyBuffer)
   → 寫入 generateContent request body 的 inlineData.data
```

收據圖片限制 ≤ 5MB，base64 後約 6.7MB，單一 worker 單次 request 完全 OK。可以用 streaming encoder 不必整檔載入記憶體。

**Request 結構**：
```go
type GenerateContentRequest struct {
    Contents         []Content        `json:"contents"`
    GenerationConfig GenerationConfig `json:"generationConfig"`
}
type Content struct {
    Parts []Part `json:"parts"`
}
type Part struct {
    Text       string      `json:"text,omitempty"`
    InlineData *InlineData `json:"inlineData,omitempty"`
}
type InlineData struct {
    MimeType string `json:"mimeType"`
    Data     string `json:"data"` // base64
}
type GenerationConfig struct {
    ResponseMimeType string          `json:"responseMimeType"` // "application/json"
    ResponseSchema   json.RawMessage `json:"responseSchema"`
}
```

`Extract` 函式流程：
1. `context.WithTimeout(ctx, 30*time.Second)`
2. 組 request（system instruction + image part）
3. POST，檢查 status code
4. parse `candidates[0].content.parts[0].text` 為已定義的 ReceiptData JSON
5. 回傳

**Mock 測試**：用 `httptest.NewServer` 攔截，把 base URL 注入 extractor（constructor 接收 `baseURL string` 參數，預設 `https://generativelanguage.googleapis.com`）。

---

## Rate Limit

新增 `internal/middleware/per_user_ratelimit.go`：
- key 用 JWT 拿到的 user_id
- 預設：每 user 每天 20 次
- 套用方式較麻煩：因為 `POST /expenses` 既處理一般 create 又處理 AI create，rate limit 不能套整個端點
  - 解法：在 handler 內，當 parse 出 `ai_extract=true` 時手動檢查 rate limit（呼叫一個 reusable 的 `CheckAIRateLimit(userID)` 函式）
- 同樣套在 `PUT /expenses/:txnId` 當 `ai_extract=true` 時
- 不套在 `ai-cancel`（不會打 LLM）

---

## 前端設計

### 型別定義
`frontend/types/transaction.ts` 加：
```ts
ai_status?: 'pending' | 'processing' | 'completed' | 'failed' | null
ai_error?: string
```

### 元件：`components/AiStatusBadge.vue`
共用 badge，根據 `ai_status` 顯示對應 icon：
- `pending` → 🤖 灰色
- `processing` → 🤖 藍色 + 旋轉動畫
- `failed` → ⚠️ 紅色（tooltip 顯示 ai_error）
- `completed` / `null` → 不渲染

### 修改 `pages/spaces/[id]/transaction/add.vue`
- 上傳圖片區下方加 checkbox：「☐ 使用 AI 辨識品項與日期」
- 勾選後：
  - 整個表單除了 title 與圖片區外全部 disable
  - 顯示說明文字：「將在背景處理，可能需要約 1 分鐘」
  - 必須至少有一張圖片才能提交
- **提交流程改寫**（不只是 AI 流程，是整個 add.vue 的儲存流程）：
  - 原本：`api.post('/api/spaces/.../expenses', json)` → `imageManagerRef.commit(txn_id)` 兩段
  - 新版：組成 `multipart/form-data`，`data` 帶 JSON payload，`images` 帶從 ImageManager 拿到的 File 物件們，一次 POST
  - `ImageManager` 仍維持 buffered 模式（`instant-upload="false"`），但 commit 的目標從「呼叫 /api/images」改為「把 File 物件回傳給 add.vue 組 multipart」
  - 提交成功後 navigate 回帳本列表
- 不勾 AI：行為相同（也走 multipart，只是不帶 ai_extract）

### 修改 `composables/useTransactionForm.ts`
- 加 `aiExtract: ref(false)`
- `buildExpensePayload` 加上 `ai_extract` 欄位（當 aiExtract=true 時帶上）
- 加 `buildExpenseMultipart(files: File[])`：把 `buildExpensePayload()` 的結果序列化成 `data` 欄位，配 `images` 檔案組成 FormData
- 既有 validation 條件式放寬（aiExtract 為 true 時不檢查 total_amount / items，但要求至少一張圖）

### 修改 `components/ImageManager.vue`
- 既有 buffered 模式 commit 行為改變：
  - 原本：`commit(entityId)` → 對每個 buffered file 呼叫 `/api/images`
  - 新版：加一個 `getBufferedFiles(): File[]` 方法給外層拿原始 File 物件
  - 既有 `commit(entityId)` 流程仍保留，給 edit 頁用（edit 頁仍走 `/api/images`）
- 換句話說 ImageManager 同時支援兩種 commit 模式

### 列表頁
在每筆 transaction row 顯示 `<AiStatusBadge :status="txn.ai_status" :error="txn.ai_error" />`。

### 詳細頁（編輯頁）
- 若 `ai_status` 為 `pending` / `processing`：
  - 頂部顯示 banner：「🤖 AI 辨識中...」+「取消辨識」按鈕
  - 整個表單 disable
  - 「取消辨識」呼叫 `POST /ai-cancel`
- 若 `ai_status` 為 `failed`：
  - 頂部顯示 banner：「⚠️ AI 辨識失敗：{ai_error}」(**純資訊，無按鈕**)
  - 表單可正常編輯
  - 額外顯示「使用 AI 重新辨識」checkbox（預設未勾選）
  - 勾選 → date / items 區塊變 disable（同 add 頁）
  - 儲存時送出 PUT，payload 帶 `ai_extract` 對應 checkbox 狀態，後端依規則自動轉換 ai_status
  - 不論勾或不勾，存完後都不會再是 failed
- 若 `ai_status` 為 `NULL` / `completed`：
  - 不顯示 banner、不顯示 AI checkbox
  - 表單行為與 AI 功能加入前完全相同

### 不做的事
- ❌ Auto polling — 使用者手動下拉刷新
- ❌ 通知中心 / toast 提醒 — 失敗只靠列表 ⚠️ icon

---

## 安全與部署

### 隱私揭露
在 add.vue 的 checkbox 旁邊顯示一行小字：

> 上傳的發票會傳送至 Google Gemini 進行辨識，不會用於模型訓練。

### 圖片驗證
- 大小 ≤ 5MB
- Type 白名單：jpeg / png / webp
- 既有 `image.go:82` 已有副檔名檢查，沿用即可

### 日誌
- 不寫圖片內容到 log，只記錄 size / user_id / 處理時間 / 結果狀態
- LLM 原始錯誤訊息只記在 server log，不回傳給前端（前端只看到友善訊息）

### Graceful shutdown
`main.go` 啟動 worker 時：
```go
ctx, cancel := context.WithCancel(context.Background())
go aiWorker.Run(ctx)

// shutdown 流程
<-quit
cancel()              // worker 收到信號，立刻終止 LLM 呼叫
srv.Shutdown(ctx2)    // 既有 HTTP server graceful shutdown
```

換版時 worker 中段被砍 → row 留在 `processing` → 下次 startup 的 recovery 自動撿回。

---

## 範圍外（不做）

明確排除，避免功能蔓延：

- ❌ BYOK（使用者自帶 API key）→ 留待 Phase 2
- ❌ 處理時分（只取日期）
- ❌ 處理分類 / 付款方式 / 分帳
- ❌ 多張發票同時辨識
- ❌ Auto polling / WebSocket / SSE
- ❌ 通知中心
- ❌ 預算追蹤、發票歷史等延伸功能
- ❌ Pending images 表（async 流程下不需要）
- ❌ Worker 啟動 dispatch（純輪詢，不做事件觸發）

---

## 實作工作清單

### Backend
1. Migration：`add_ai_status_to_transactions`（含 index）
2. `internal/models/transaction.go` — 加 `AIStatus`、`AIError` 欄位
3. `internal/config/config.go` + `.env-backend.example` — Gemini 設定
4. `internal/services/ai_extract.go` — `ReceiptExtractor` interface + `GeminiReceiptExtractor` 實作
5. `internal/services/ai_extract_test.go` — 用 httptest mock Gemini API
6. `internal/services/ai_worker.go` — 輪詢 worker + 三階段處理 + recovery
7. `internal/services/ai_worker_test.go` — 用 fake extractor 注入測試
8. `internal/handlers/expense.go` — **改寫 `Create` handler 接受 multipart，內聯處理圖片上傳到 R2 + image record 建立**；改寫 `Update` handler 接受 `ai_extract` 並依規則轉換 ai_status
9. `internal/handlers/transaction.go` — 新增 `AICancel` handler
10. `internal/handlers/expense_test.go` / `transaction_test.go` — 對應測試（含 multipart 路徑）
11. `internal/middleware/per_user_ratelimit.go` — 新增 per-user rate limit helper（不是 middleware，是給 handler 內呼叫的函式）
12. `internal/services/transaction_service.go` — `CreateExpense` 簽章可能需要擴充：接受 image files + ai_extract，在同一個 db.Transaction 內處理圖片寫入；或新增 `CreateExpenseWithImages` 平行函式
13. `main.go` — 啟動 worker goroutine、graceful shutdown 整合

### Frontend
14. `frontend/types/transaction.ts` — 加 `ai_status`、`ai_error`
15. `frontend/components/AiStatusBadge.vue` — 新增
16. `frontend/components/ImageManager.vue` — 加 `getBufferedFiles()` 方法給外層拿原始 File；既有 `commit()` 行為保留給 edit 頁
17. `frontend/composables/useTransactionForm.ts` — 加 `aiExtract` 狀態、multipart builder、條件式驗證
18. `frontend/composables/useApi.ts`（如有需要）— 確認 multipart upload helper 可用
19. `frontend/pages/spaces/[id]/transaction/add.vue` — checkbox + 條件式 disable + 改寫提交流程為單一 multipart POST
20. `frontend/pages/spaces/[id]/transaction/[txnId]/edit.vue` — failed 狀態時顯示 banner + AI checkbox；pending/processing 顯示 banner + cancel 按鈕（form disable）；PUT payload 加 ai_extract
21. 帳本列表頁面 — 顯示 `AiStatusBadge`
22. `frontend/tests/AiStatusBadge.test.ts` — 新增
23. `frontend/tests/useTransactionForm.test.ts` — 補 ai 流程的測試

### 文件
24. `.env-backend.example` 加註解
25. README 不需動（feature 屬於空間管理大類，README 已涵蓋）

---

## 開放問題

實作前的最後確認：

1. **Gemini API key 由誰提供** — 你會自己申請並填到 `.env-backend`，我會在 example 加佔位即可
2. **Rate limit 數字** — 預設 20 次/日/user，可在 env 調整
3. **Worker 輪詢間隔** — 10 秒
4. **詳細頁路徑** — 需確認 `pages/spaces/[id]/transaction/[txnId]/` 的具體檔名
