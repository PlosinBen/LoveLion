# E2E Testing — Playwright

## 目標

為 LoveLion 建立瀏覽器自動化測試，覆蓋前端互動 → API → 頁面結果的完整鏈路。

## 架構決策

### Playwright 獨立 Container

Playwright 跑在獨立的 Docker service（`profiles: [e2e]`），不安裝在 frontend image 中。

**原因**：Playwright + Chromium 約 500MB，放進 frontend image 會拖慢日常開發的 build 與啟動。獨立 container 只在需要時啟動，保持開發環境輕量。

### 測試環境網路

```
playwright container
  └─ browser → http://frontend:3000（頁面）
                  └─ Nuxt proxy /api/** → http://backend:8080（API）
```

瀏覽器載入頁面後，前端的 API 呼叫走相對路徑 `/api/**`，由 Nuxt proxy 轉發到 backend。不需要修改現有架構。

### DB 重置

複用現有的 `bin/refresh-database` 流程（drop → create → migrate → seed），確保每次 E2E 測試都從乾淨狀態開始。

## 檔案清單

### 新建

| 檔案 | 說明 |
|------|------|
| `frontend/playwright.config.ts` | Playwright 設定：baseURL `http://frontend:3000`、headless、timeout、reporter |
| `frontend/e2e/fixtures/auth.ts` | 共用 fixture：登入後取得已認證的 page，避免每個測試重複登入 |
| `frontend/e2e/01-auth.spec.ts` | 認證測試：登入、註冊新帳號、登出 |
| `frontend/e2e/02-core-flow.spec.ts` | 核心流程：建立空間 → 新增消費 → 驗證帳本列表出現該筆消費 |
| `frontend/e2e/03-announcements.spec.ts` | 公告測試：公告列表頁、公告詳情頁（需先透過 API 建立測試資料） |
| `bin/e2e_test` | 一鍵執行腳本：DB 重置 → 確認服務就緒 → 啟動 playwright container → 輸出結果 |

### 修改

| 檔案 | 變更 |
|------|------|
| `docker-compose.yml` | 新增 `playwright` service，使用 `mcr.microsoft.com/playwright` image，`profiles: [e2e]`，掛載 `frontend/` |
| `frontend/package.json` | devDependencies 加入 `@playwright/test`，scripts 加入 `test:e2e` |
| `frontend/.gitignore` | 加入 `playwright-report/`、`test-results/`、`blob-report/` |
| `.github/workflows/ci.yml` | 新增 `e2e` job（DB 重置 → 啟動服務 → 跑 Playwright） |

## 測試案例設計

### 01-auth.spec.ts

1. **登入成功**：輸入 dev/dev123 → 導向首頁（空間列表）
2. **登入失敗**：錯誤密碼 → 顯示錯誤訊息
3. **註冊新帳號**：填入新帳號 → 自動登入 → 導向首頁
4. **登出**：設定頁點登出 → 導回登入頁

### 02-core-flow.spec.ts

前置：以 dev/dev123 登入

1. **建立空間**：首頁點新增 → 填入名稱/類型 → 建立成功 → 進入空間
2. **新增消費**：帳本頁 → 新增消費 → 填寫金額/備註 → 儲存
3. **驗證帳本**：回到帳本列表 → 確認剛才的消費出現在列表中

### 03-announcements.spec.ts

前置：透過 API 以 admin 身份建立一則 published 公告

1. **公告列表**：設定頁 → 點公告連結 → 列表頁顯示該公告
2. **公告詳情**：點擊公告 → 詳情頁顯示標題與 Markdown 內容

## SPA 注意事項

- `ssr: false` 代表頁面載入後需等 Vue mount 完成，使用 `page.waitForSelector()` 或 Playwright 的 auto-wait 定位可見元素
- 導航後用 `page.waitForURL()` 確認路由切換完成
- LoadingOverlay 出現時需等其消失再做 assertion
- DatePicker 元件複雜，測試中使用預設日期（今天），不操作 DatePicker

## 執行方式

```bash
# 本地一鍵執行
./bin/e2e_test

# 手動執行（開發/除錯用）
docker compose --profile e2e run --rm playwright npx playwright test

# 只跑特定測試
docker compose --profile e2e run --rm playwright npx playwright test e2e/01-auth.spec.ts

# 查看報告
docker compose --profile e2e run --rm playwright npx playwright show-report
```
