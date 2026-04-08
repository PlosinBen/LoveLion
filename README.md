# LoveLion

共享記帳與比價應用程式。以「Space」作為統一概念，將交易、商店與成員歸屬於同一空間。

## 功能

### 空間管理
- 建立個人、旅行、群組空間
- 自訂幣別、分類、付款方式
- 空間置頂 / 封面圖片
- 成員邀請連結（一次性 / 多次使用 / 可設到期日）

### 消費記帳
- 新增 / 編輯 / 刪除消費紀錄
- 多幣別支援，自動換算匯率與手續費
- 項目明細（名稱、單價、數量、折扣）
- 分類、付款方式、Google Maps 地點連結
- 收據 / 照片上傳（Cloudflare R2 儲存）
- 消費模板：從現有交易儲存為模板，新增時一鍵套用

### 分帳與付款
- 多人分帳，支援自訂金額分配與均分
- 現付標記（當場結清）
- 付款紀錄追蹤還款進度

### 比價功能
- 建立商店（含 Google Maps 連結）
- 新增商品並記錄各店價格
- 同商品跨店比價一覽

### 統計分析
- 空間消費統計
- 分類支出分佈

### 使用者
- 帳號註冊 / 登入（JWT 認證）
- 個人資料管理

## 技術棧

- **後端**：Go 1.25 / Gin / GORM / PostgreSQL 18
- **前端**：Nuxt 4 (Vue 3) / Tailwind CSS / TypeScript
- **測試**：Go testing / Vitest / httpexpect（整合測試）
- **基礎設施**：Docker Compose / Cloudflare R2 / Cloudflare Tunnel

## 開發

所有指令透過 Docker 執行，宿主機不需安裝 Go / Node。

```bash
# 啟動開發環境
docker compose up -d

# 後端測試
docker compose exec backend go test ./...

# 前端測試
docker compose exec frontend npm test

# 完整整合測試（含 DB 重置、seed、API 測試）
./bin/integration_test

# 資料庫重置
./bin/refresh-database

# 執行遷移
./bin/migrate
```

測試帳號：`dev` / `dev123`

## 部署

使用 tag-based 自動部署，cron 每 5 分鐘檢查新 tag：

```bash
git tag v1.0.0
git push origin v1.0.0
```

Production 環境使用 `docker-compose.prod.yml`，包含：
- 獨立 migration container（啟動時自動遷移）
- Backend graceful shutdown（10 秒 timeout）
- 各 container 資源限制（適配 4GB RAM 機器）

## 專案結構

```
backend/
  main.go                  # 路由定義與伺服器啟動
  internal/
    handlers/              # HTTP 處理器
    models/                # GORM 模型
    services/              # 業務邏輯
    repositories/          # 資料存取
    middleware/            # 認證、權限、日誌、限流
  migrations/              # SQL 遷移檔
  cmd/seed/                # 種子資料

frontend/
  pages/                   # 檔案式路由
  components/              # 共用元件
  composables/             # 組合式函式
  stores/                  # Pinia 狀態管理
  types/                   # TypeScript 型別定義
```
