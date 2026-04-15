# PROJECT MAP

意圖導向的檔案索引。用功能意圖快速定位檔案。

## 進入點與啟動

| 意圖 | 路徑 | 說明 |
|------|------|------|
| 看所有 API 路由 | `backend/main.go` | 路由定義、中介層串接、依賴注入、AI Worker 啟動 |
| 看 Nuxt 設定 | `frontend/nuxt.config.ts` | SSR 開關、API 代理規則、Vite/Tailwind 設定 |
| 看 Docker 服務編排 | `docker-compose.yml` | postgres / backend / frontend 三服務定義 |
| 看 E2E 測試服務編排 | `docker-compose.test.yml` | 隔離測試環境：ephemeral postgres、seed、promote-admin、playwright |
| 看生產部署設定 | `docker-compose.prod.yml` | 獨立 migration container、資源限制、graceful shutdown |
| 看 PWA 設定 | `frontend/public/manifest.json` | Web App Manifest：standalone 模式、icon、theme color |

## 認證與授權

| 意圖 | 路徑 | 說明 |
|------|------|------|
| 改 JWT 驗證邏輯 | `backend/internal/middleware/auth.go` | Bearer token 解析、claims 驗證、userID 注入 context |
| 改空間存取權限 | `backend/internal/middleware/space.go` | SpaceAccess（成員驗證）、SpaceOwnerOnly（擁有者限制） |
| 改管理員權限 | `backend/internal/middleware/admin.go` | AdminOnly 中介層，檢查 user.Role == "admin" |
| 改速率限制 | `backend/internal/middleware/ratelimit.go` | 通用 RateLimit |
| 改 AI 速率限制 | `backend/internal/middleware/ai_ratelimit.go` | AI 專用 AIRateLimiter |
| 改請求日誌 | `backend/internal/middleware/request_logger.go` | 結構化 slog 請求日誌 |
| 改前端路由守衛 | `frontend/middleware/auth.global.ts` | 全局 auth middleware，公開路由：/login、/join/* |
| 改登入/註冊/登出 | `frontend/stores/auth.ts` | Pinia store：token 持久化、localStorage、自動初始化 |
| 改 HTTP 客戶端 | `frontend/composables/useApi.ts` | Bearer 注入、401 自動登出、FormData 偵測 |

## 空間管理

| 意圖 | 路徑 | 說明 |
|------|------|------|
| 改空間 CRUD API | `backend/internal/handlers/space.go` | 建立/讀取/更新/刪除空間、離開空間 |
| 改空間資料模型 | `backend/internal/models/space.go` | JSONB 欄位（currencies/categories/paymentMethods/splitMembers） |
| 改成員/邀請 API | `backend/internal/handlers/space_sharing.go` | 成員列表/別名/移除、邀請建立/列出/撤銷/加入 |
| 改成員資料模型 | `backend/internal/models/space.go` | SpaceMember：role (member/owner)、alias（與 Space 同檔） |
| 改邀請業務邏輯 | `backend/internal/services/invite_service.go` | 邀請驗證、使用次數、到期檢查 |
| 改空間列表頁 | `frontend/pages/index.vue` | 首頁：所有空間列表、置頂/分類 |
| 改空間設定頁 | `frontend/pages/spaces/[id]/settings.vue` | 空間名稱/幣別/分類/成員/邀請管理 |
| 改空間 composable | `frontend/composables/useSpace.ts` | useState 全域空間列表、置頂、離開；含 useLedger() 向後相容 |
| 改空間詳情 store | `frontend/stores/spaceDetail.ts` | 當前空間資料、懶載入、invalidation 機制 |

## 交易記帳

| 意圖 | 路徑 | 說明 |
|------|------|------|
| 改消費/付款 API | `backend/internal/handlers/transaction.go` | 列表/詳情/刪除；路由分流 expense vs payment |
| 改消費 CRUD API | `backend/internal/handlers/expense.go` | 建立/更新消費（含 multipart 圖片上傳） |
| 改付款 CRUD API | `backend/internal/handlers/payment.go` | 建立/更新付款 |
| 改交易業務邏輯 | `backend/internal/services/transaction_service.go` | DB Transaction 包裝多表寫入、R2 上傳、AI 狀態轉換 |
| 改交易資料模型 | `backend/internal/models/transaction.go` | Transaction：short ID、AIStatus/AIError、關聯 Expense/Debts/Images |
| 改消費明細模型 | `backend/internal/models/transaction.go` | TransactionExpense + TransactionExpenseItem（與 Transaction 同檔） |
| 改分帳模型 | `backend/internal/models/transaction_debt.go` | TransactionDebt：payer/payee、settled_amount、is_spot_paid |
| 改交易 Repository | `backend/internal/repositories/transaction_repo.go` | 分頁查詢、搜尋/分類/日期範圍篩選 |
| 改帳本頁面 | `frontend/pages/spaces/[id]/ledger.vue` | 交易列表、搜尋、篩選 |
| 改新增交易頁 | `frontend/pages/spaces/[id]/transaction/add.vue` | 快速掃描收據入口、消費/付款表單、AI 辨識開關、模板套用 |
| 改編輯交易頁 | `frontend/pages/spaces/[id]/transaction/[txnId]/edit.vue` | 編輯消費/付款、AI 重試 |
| 改交易詳情頁 | `frontend/pages/spaces/[id]/transaction/[txnId]/index.vue` | 交易明細、圖片、分帳資訊 |
| 改交易表單邏輯 | `frontend/composables/useTransactionForm.ts` | 表單狀態、驗證、multipart payload 組裝、幣別換算 |
| 改消費表單元件 | `frontend/components/ExpenseForm.vue` | 消費表單 UI（品項/分類/幣別/付款方式） |
| 改付款表單元件 | `frontend/components/PaymentForm.vue` | 付款表單 UI（付款人/收款人/金額） |
| 改分帳編輯元件 | `frontend/components/DebtEditor.vue` | 分帳人員與金額編輯 |

## AI 收據辨識

| 意圖 | 路徑 | 說明 |
|------|------|------|
| 改 AI Worker | `backend/internal/services/ai_worker.go` | 背景 polling（10s）、狀態機、重試與錯誤恢復 |
| 改 Gemini 呼叫 | `backend/internal/services/ai_extract.go` | Gemini Vision API 呼叫、結構化輸出解析、超時 35s |
| 改 AI 取消 API | `backend/internal/handlers/transaction.go` | POST /transactions/:id/ai-cancel |
| 改 AI 速率限制 | `backend/internal/middleware/ai_ratelimit.go` | 每用戶每日上限（預設 20 次） |
| 看 AI 功能規格 | `.agent/features/ai-receipt-extraction.md` | 完整功能規格文件 |

## 比價功能

| 意圖 | 路徑 | 說明 |
|------|------|------|
| 改商店/商品 API | `backend/internal/handlers/comparison.go` | 商店 CRUD、商品 CRUD、跨店比價 |
| 改比價資料模型 | `backend/internal/models/comparison.go` | Store、Product 模型 |
| 改商店列表頁 | `frontend/pages/spaces/[id]/stores.vue` | 商店列表、Google Maps 連結 |
| 改商品列表頁 | `frontend/pages/spaces/[id]/products.vue` | 商品與價格比較 |

## 圖片管理

| 意圖 | 路徑 | 說明 |
|------|------|------|
| 改圖片 API | `backend/internal/handlers/image.go` | 上傳/列表/排序/刪除 |
| 改圖片資料模型 | `backend/internal/models/image.go` | 多態關聯（entity_type + entity_id）、blur_hash |
| 改 R2 儲存 | `backend/internal/storage/r2.go` | Cloudflare R2 (S3 相容) 上傳/下載/刪除 |
| 改圖片元件 | `frontend/components/ImageManager.vue` | 緩衝上傳模式、拖曳排序、即時/延遲刪除 |
| 改圖片壓縮 | `frontend/composables/useImages.ts` | browser-image-compression、收據保留原始品質 |

## 公告系統

| 意圖 | 路徑 | 說明 |
|------|------|------|
| 改公告 API | `backend/internal/handlers/announcement.go` | 公開讀取 + Admin CRUD + AI 生成 |
| 改公告資料模型 | `backend/internal/models/announcement.go` | status (draft/published)、broadcast 時間窗 |
| 改公告管理頁 | `frontend/pages/admin/announcements/` | Admin 公告管理介面 |
| 改公告列表/詳情 | `frontend/pages/announcements/` | 使用者端公告檢視（Markdown 渲染） |
| 改廣播列 | `frontend/components/BroadcastBar.vue` | Header 廣播列、localStorage 記憶關閉狀態 |
| 看公告功能規格 | `.agent/features/announcement-system.md` | 完整功能規格文件 |

## 模板功能

| 意圖 | 路徑 | 說明 |
|------|------|------|
| 改模板 API | `backend/internal/handlers/expense_template.go` | 消費模板 CRUD |
| 改模板資料模型 | `backend/internal/models/expense_template.go` | ExpenseTemplate 模型 |
| 改模板選擇元件 | `frontend/components/TemplatePickerModal.vue` | 模板選擇彈窗 |
| 改模板 composable | `frontend/composables/useExpenseTemplates.ts` | 模板資料 CRUD |

## 統計分析

| 意圖 | 路徑 | 說明 |
|------|------|------|
| 改統計頁面 | `frontend/pages/spaces/[id]/stats.vue` | 消費統計、分類分佈、建議付款連結 |
| 改統計元件 | `frontend/components/SpaceStats.vue` | 統計圖表與資料計算 |

## 使用者設定

| 意圖 | 路徑 | 說明 |
|------|------|------|
| 改使用者 API | `backend/internal/handlers/auth.go` | 註冊/登入/取得個人資料/更新個人資料 |
| 改使用者模型 | `backend/internal/models/user.go` | bcrypt 密碼、IsAdmin()、角色 (user/admin) |
| 改登入頁 | `frontend/pages/login.vue` | 登入/註冊表單 |
| 改設定頁 | `frontend/pages/settings.vue` | 個人資料、密碼變更 |
| 改關於頁 | `frontend/pages/about.vue` | 版本號、產品故事、GitHub 連結 |

## 共用 UI 元件

| 意圖 | 路徑 | 說明 |
|------|------|------|
| 改基礎按鈕 | `frontend/components/BaseButton.vue` | 統一按鈕樣式 |
| 改基礎輸入 | `frontend/components/BaseInput.vue` | 統一輸入框 |
| 改基礎 Modal | `frontend/components/BaseModal.vue` | 統一彈窗 |
| 改 Header | `frontend/components/Header.vue` | 頂部導航列 |
| 改底部導航 | `frontend/components/BottomNav.vue` | 空間內三分頁導航 |
| 改沉浸式 Header | `frontend/components/ImmersiveHeader.vue` | 帶封面圖的 Header |
| 改 Toast 通知 | `frontend/components/AppToast.vue` + `frontend/stores/toast.ts` | 全域提示訊息 |
| 改確認對話框 | `frontend/components/AppConfirm.vue` + `frontend/stores/confirm.ts` | Promise-based 確認 |
| 改載入覆蓋層 | `frontend/components/LoadingOverlay.vue` + `frontend/stores/loading.ts` | 全域 loading |
| 改佈局 | `frontend/layouts/default.vue` | Header + BroadcastBar + slot + BottomNav |
| 改空佈局 | `frontend/layouts/empty.vue` | 無 Header/Nav 的乾淨佈局 |
| 改根元件 | `frontend/app.vue` | NuxtLayout + NuxtPage + LoadingOverlay |

## ID 生成與工具

| 意圖 | 路徑 | 說明 |
|------|------|------|
| 改 Short ID 生成 | `backend/internal/utils/id.go` | 自製短 ID（a-zA-Z0-9），碰撞自動加長，crypto/rand |
| 改環境設定 | `backend/internal/config/config.go` | 所有環境變數定義與預設值 |
| 改資料庫連線 | `backend/internal/database/database.go` | GORM PostgreSQL 連線設定 |

## 資料庫遷移

| 意圖 | 路徑 | 說明 |
|------|------|------|
| 看/加遷移檔 | `backend/migrations/` | golang-migrate SQL 檔，依序 000001-000010 |
| 執行遷移 | `backend/cmd/migrate/main.go` | 遷移執行器 |
| 看種子資料 | `backend/cmd/seed/main.go` | dev/ming/mei 三用戶 + 示範空間/交易/比價 |
| 重置資料庫 | `bin/refresh-database` | 停服務→刪庫→重建→遷移→seed |
| 執行遷移腳本 | `bin/migrate` | docker compose exec 包裝 |

## 測試

| 意圖 | 路徑 | 說明 |
|------|------|------|
| 後端整合測試 | `backend/integration/*_test.go` | httpexpect 整合測試，01-07 號測試檔 |
| 後端單元測試 | `backend/internal/handlers/*_test.go` | httptest + testutil 單元測試 |
| 後端測試工具 | `backend/internal/testutil/` | TestDB、CreateTestUser、AuthContext 等 helper |
| 前端單元測試 | `frontend/tests/composables/` | Vitest 測試 useAuth/useSpace/useTransactionForm 等 |
| 前端單元測試設定 | `frontend/vitest.config.ts` | Vitest 設定檔 |
| E2E 測試 | `frontend/e2e/` | Playwright E2E 測試：auth、core-flow、announcements |
| E2E 測試 fixtures | `frontend/e2e/fixtures/auth.ts` | 登入 helper 與 authedPage fixture |
| E2E 測試設定 | `frontend/playwright.config.ts` | Playwright 設定：單 worker、Chromium、HTML reporter |
| 執行整合測試 | `bin/integration_test` | DB 重置 + seed + 完整 API 測試 |
| 執行 E2E 測試 | `bin/e2e_test` | 隔離 Docker 環境啟動 → seed → playwright → cleanup |
| 看 E2E 功能規格 | `.agent/features/e2e-testing.md` | E2E 測試架構與設計文件 |

## TypeScript 型別

| 意圖 | 路徑 | 說明 |
|------|------|------|
| 改使用者/公告型別 | `frontend/types/user.ts` | User、Announcement interface |
| 改空間相關型別 | `frontend/types/space.ts` | Space、SpaceMember、SpaceInvite interface |
| 改交易相關型別 | `frontend/types/transaction.ts` | Transaction、Expense、Debt 等 interface |
| 改比價型別 | `frontend/types/comparison.ts` | Store、Product interface |
| 改圖片型別 | `frontend/types/image.ts` | Image interface |
