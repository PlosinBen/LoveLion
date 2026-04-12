# 公告系統設計

## 資料面

- `users` 表加 `role` 欄位（`user`/`admin`，default `user`）
- `announcements` 表：
  - `id` — NanoID（URL 可見）
  - `title` — 標題
  - `content` — 內容（Markdown）
  - `status` — `draft`（草稿，僅 admin 可見）/ `published`（已發布，公開可見）
  - `broadcast_start` — 廣播開始時間，NULL = 不廣播
  - `broadcast_end` — 廣播結束時間，NULL = 不廣播
  - `created_at`
  - `updated_at`
- `AdminOnly` middleware 檢查 role
- Admin 身份：第一個 admin 手動 SQL 設定，後續可考慮 UI 管理
- 前端判斷 admin：`GET /api/users/me` 回傳的 User 已含 `role` 欄位

## API

### 公開
- `GET /api/announcements` — 列表（status = published，依 created_at 排序，暫不分頁）
- `GET /api/announcements/:id` — 單則詳情（published only）
- `GET /api/announcements/broadcast` — 當前廣播中的最新一則（published AND now BETWEEN broadcast_start AND broadcast_end）

### Admin only
- `GET /api/admin/announcements` — 列表（全部，含 draft）
- `POST /api/admin/announcements` — 新增
- `PUT /api/admin/announcements/:id` — 編輯
- `DELETE /api/admin/announcements/:id` — 刪除
- `POST /api/admin/announcements/generate` — AI 生成公告（同步，即時回傳 title + content）

## 前端

### Header 下方公告條
- 所有頁面都顯示
- 顯示當前廣播中的最新一則公告
- 點擊內容進入 `/announcements/:id`
- 右邊有 X 可關閉，localStorage 記住已關閉的公告 ID，不再顯示同一則
- 無廣播中公告時隱藏

### `/announcements` — 公告列表頁
- 顯示所有已發布的公告

### `/announcements/:id` — 公告詳情頁
- 麵包屑：首頁 > 公告 > {公告標題}
- 內容以 Markdown 渲染
- 底部「查看所有公告」連結回列表頁

### 入口
- Settings 頁面內連結 → `/announcements`
- Header 公告條點擊 → `/announcements/:id`
- 公告詳情頁麵包屑 + 底部連結 → `/announcements`

### Admin 管理頁面（`/admin/announcements`）
- 只有 admin 可見，入口放 settings 頁
- 公告列表（含 draft，顯示狀態標籤）
- 新增/編輯公告：
  - title、content（textarea，Markdown）
  - status 切換（draft / published）
  - 廣播期間設定（可選）
- AI 生成：輸入一句描述，呼叫 Gemini 同步生成 title + content（繁體中文），填入表單供 admin 修改後儲存。若 Gemini 未設定（無 API key）則隱藏按鈕
- 刪除公告（硬刪；不需要刪除時 admin 可改 status 為 draft 隱藏）
- Markdown 渲染使用 `marked`

## About 頁面

- 獨立 `/about`，從 settings 連過去
- 顯示版本號（env 注入 `NUXT_PUBLIC_APP_VERSION`，deploy script 設定）
- Built by Ben © 2026
- GitHub 連結：https://github.com/PlosinBen/LoveLion

### 產品介紹文案

**LoveLion** 最初是因為和朋友出國時的分帳問題而誕生的。

多人旅行中，代付、外幣找零、統一刷卡的匯率與手續費，每一筆都讓結算變得更複雜。以前我都是先用記帳 app 記錄，回國後再手動整理到 Google Sheet，逐筆填入匯率和手續費，最後再一個一個算出每個人該付多少。這個過程既分散又容易出錯。

另一個困擾是比價。旅途中看到想買的東西，會沿路比較不同店家的價格，但每次都是拍照記錄，等到要確認的時候又得翻相簿翻很久，很難快速對照。

所以我做了 LoveLion — 把記帳、比價、匯率換算、分帳和結算收在同一個地方。旅途中就能即時記錄，每個人隨時都能看到消費明細和自己還需要付多少錢，比價紀錄也不用再靠照片大海撈針。

後來發現，日常生活中的記帳需求其實也差不多，就把它一起收進來了。

這是一個為了解決自己問題而做的小工具，希望也能幫到有同樣困擾的人。如果你有任何功能建議或想法，歡迎到 [GitHub](https://github.com/PlosinBen/LoveLion) 上跟我聊聊。
