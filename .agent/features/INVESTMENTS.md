# 投資紀錄功能

## 概述

獨立於 Space 體系的個人投資損益追蹤功能。支援多種投資類型（期貨、股票等），每月結算各類型損益，按權重分配給各參與者。

## 存取控制

- 透過 `inv_members.user_id` + `is_owner` 判斷權限
- 登入用戶的 user_id 存在於 `inv_members`（active）→ 放���
- `is_owner = true` → 完整操作權限
- 非 owner → 僅能查看自己的損益紀錄
- 未在 inv_members 中的帳號完全看不到此功能（前端不渲染入口、後端拒絕請求）

## 前端入口

- 首頁 BottomNav 新增「投資」tab（僅 inv_members 可見）
- 路由：`/investments`，獨立頁面體系

## 核心概念

### 成員與權重

- 投資人可能包含本人及代操對象
- 成員可新增、可隱藏（`active` 狀態），不可刪除（保留歷史資料）
- 每人有投入本金，權重 = floor((上期 balance - 當期 withdrawal) / 5000)（最小 1）
- 上月結餘扣除當月出金後計算的權重 = 本月分配比例
- **加碼**：下個月生效
- **減碼**：當月立即生效，直接用減碼後的新權重算整個月損益

### 投資類型

- 可擴充設計，每種類型有自己的子表、表單、損益計算公式
- 主表（月結算、成員分配）固定不動，新增類型只需加子表和對應邏輯
- 目前支援：期貨、股票

### 月結算流程

1. 手動建立月結算紀錄（選擇年月），初始為 draft 狀態
2. 填入各投資類型的結算資料
3. 各類型依自己的公式算出月損益
4. 所有類型月損益加總
5. 按當月各人權重比例分配損益
6. 確認後標記為完成

### 損益分配計算

```
每單位權重損益 = floor(總損益 / 總權重)

其他人損益 = 每單位權重損益 × 個人權重
我的損益   = 總損益 - Σ(其他人損益)
```

- 其他成員：用每單位權重損益 × 個人權重計算
- 自己（is_owner）：用「總損益 - 已分配給其他人的損益」計算，作為尾差吸收者，避免小數誤差

#### 範例

```
總損益 +11,803  總權重 14

每單位權重損益 = floor(11803 / 14) = 843

A (權重 3)：843 × 3 = 2,529
B (權重 1)：843 × 1 = 843
已分配：2,529 + 843 = 3,372
我(權重 10)：11,803 - 3,372 = 8,431
```

### 成員入出金 vs 投資帳戶入出金

兩者完全獨立：

- **投資帳戶入出金**（期貨/股票 statement 的 deposit/withdrawal）：實際匯入匯出券商帳戶的錢，影響損益計算公式
- **成員入出金**（allocation 的 deposit/withdrawal）：成員之間的權益調整，只影響個人 balance 和權重，不一定對應實際帳戶操作

例如：owner 把自己的額度轉 10,000 給 A → owner withdrawal 10,000 + A deposit 10,000，但投資帳戶沒有任何動作。

### 前期資料

- 損益計算需要前期數據（期貨的前期實質權益、股票的前期總權益）
- 沒有前期資料時，前期值視為 0（入出金會抵銷，不影響損益正確性）
- 期貨從 2024-01 開始有歷史資料（初始值透過 seed data 建立）
- 股票從 2026-05 開始，無前期

## 投資類型：��貨

### 月結單欄位

| 欄位 | 說明 |
|------|------|
| 期末權益 | 帳戶期末總權益 |
| 浮動損益 | 未平倉部位的浮動損益 |
| 沖銷損益 | 已平倉的實現損益 |
| 總入金 | 當月入金總額 |
| 總出金 | 當月出金總額 |

### 損益計算

```
實質權益 = 期末權益 - 浮動損益
權益損益 = 實質權益 - 前期實質權益 - 入金 + 出金
期貨損益 = min(沖銷損益, 權益損益)
```

## 投資類型：股票

### 月結單欄位

| 欄位 | 說明 |
|------|------|
| 帳戶餘�� | 券商帳戶現金餘額 |
| 庫存現值 | 庫存現值總額 Σ(shares × closing_price) |
| 總入金 | 當月入金總額 |
| 總出金 | 當月出金總額 |
| 庫存結餘 | 各股：代號/股數/結算價 |

### 交易紀錄（輔助查閱，不參與損益主公式）

| 欄位 | 說明 |
|------|------|
| 日期 | 交易日期 |
| 代��� | 股票代號 |
| 股數 | 交易股數（正=買/負=賣） |
| 成交價 | 每股成交價 |
| 手續費 | 券商手續費 |
| 交易稅 | 證交稅 |

交易紀錄獨立於 settlement，可在 settlement 建立前新增。

### 損益計��

```
庫存現值 = Σ(各股結算價 × 股數)
本期總權益 = 帳戶餘額 + 庫存現值
股票損益 = 本期總權益 - 前期總權益 - 入金 + 出金
```

## UI

### 路由結構 (所有子頁面使用 overlay)

```
/investments                                    ← 損益列表 (parent, 預設頁)
/investments/settlements                        ← 結算列表 (parent)
/investments/settlements/:ym                    ← 結算總覽 (overlay, 唯讀展示所有類型欄位 + 分配預覽)
/investments/settlements/:ym/futures            ← 期貨月結單編輯 (overlay)
/investments/settlements/:ym/stocks             ← 股票月結單編輯 (overlay, 含庫存)
/investments/trades                             ← 股票交易清單 (parent)
/investments/trades/add                         ← 新增交易 (overlay)
/investments/trades/:id/edit                    ← 編輯交易 (overlay)
/investments/transactions                       ← 成員入出金明細 (parent)
/investments/transactions/add                   ← 新增 (overlay)
/investments/transactions/:id/edit              ← 編輯 (overlay)
/investments/settings                           ← 成員管理 (overlay, 右上角齒輪進入)
```

### BottomNav

- Owner：`損益 | 結算 | 交易 | 異動`，右上角設定 icon → 成員管理
- 非 Owner：無 BottomNav，僅顯示損益列表（只有自己的資料）

### 損益列表頁 `/investments`

水平滾動表格，Y 軸（月份）凍結在左側：

```
         |    我                    |    A                    |    B              |
         | 入金  出金  損益  結餘    | 入金  出金  損益  結餘    | 入金  出金  損益  結餘 |
2026-05  | 0    0    +8,431 58,431 | 0    0    +2,529 17,529 | ...              |
2026-04  | ...                     | ...                     | ...              |
...（預設 12 個月）
```

### 結算列表頁 `/investments/settlements`

緊密卡片模擬表格，點擊 row 進入結算總覽：

```
┌──────────┬────────────┬──────┐
│ 2026-05  │ +11,803    │ 🟡   │  ← draft (黃色)
├──────────┼────────────┼──────┤
│ 2026-04  │ +13,500    │ 🟢   │  ← completed (綠色)
└──────────┴────────────┴──────┘

FAB → 新增月結算
```

### 結算總覽頁 `/investments/settlements/:ym`

唯讀展示所有類型月結單欄位 + 分配預覽。每個類型區塊有「編輯」按鈕進入對應編輯頁。

```
2026-05 (draft)

── 期貨 ──────────────────── [編輯]
期末權益        120,000
浮動損益          5,000
沖銷損益         15,000
入金                  0
出金                  0
損益            +15,000

── 股票 ──────────────────── [編輯]
帳戶餘額         80,000
庫存現值        150,000
入金                  0
出金                  0
損益             -3,197

── 成員入出金 ──────────── [→ 異動頁面]  (唯讀，顯示該月加總)
我    入金 0     出金 0
A     入金 10,000  出金 0
B     入金 0     出金 0

── 分配預覽 ──
總損益 +11,803
我     +8,431 (權重 10)
A      +2,529 (權重 3)
B        +843 (權重 1)

[完成結算]  /  [重新開啟] (已完成時)
```

### 交易清單頁 `/investments/trades`

預設顯示近 12 個月，可調整篩選範圍。FAB 新增，點擊 row 進入編輯。

### 成員管理 `/investments/settings`

成員列表，支援新增��切換 active、調整排序。右上角齒輪 icon 進入。

## DB Schema

金額欄位皆為 `integer`（台幣無小數），僅股價為 `decimal(10,2)`。

```
── inv_members (投資成員)
   id              varchar(21) PK       ← NanoID
   name            varchar(50)
   user_id         UUID FK → users (nullable)
   is_owner        boolean              ← 尾差吸收者，僅一人 true
   active          boolean
   sort_order      integer
   net_investment  integer              ← 累計投入淨額 (總入金 - 總出金)
   created_at

── inv_settlements (月結算主表)
   year_month              varchar(7) PK    ← "2026-05"
   status                  varchar(10)      ← draft / completed
   total_profit_loss       integer
   total_weight            integer
   profit_loss_per_weight  integer          ← floor(total_profit_loss / total_weight)
   created_at
   updated_at

── inv_member_transactions (成員入出金/損益明細)
   id          UUID PK
   member_id   varchar(21) FK → inv_members
   date        date
   type        varchar(20)          ← deposit / withdrawal / profit_loss
   amount      integer
   note        text

── inv_settlement_allocations (月結算分配，completed 時從 member_transactions 加總寫入)
   year_month   varchar(7) PK, FK → inv_settlements
   member_id    varchar(21) PK, FK → inv_members
   weight       integer              ← 本期分配權重 = floor((上期 balance - 當期 withdrawal) / 5000), 最小 1
   amount       integer              ← 本期分配損益 (從 member_transactions 加總)
   deposit      integer              ← 本期入金 (從 member_transactions 加總)
   withdrawal   integer              ← 本期出金 (從 member_transactions 加總)
   balance      integer              ← 本期結餘 (上期結餘 + 入金 - 出金 + 損益)

── inv_futures_statements (期貨月結單)
   year_month              varchar(7) PK, FK → inv_settlements
   ending_equity           integer
   floating_profit_loss    integer
   realized_profit_loss    integer
   deposit                 integer
   withdrawal              integer
   profit_loss             integer

── inv_stock_statements (股票月結單)
   year_month      varchar(7) PK, FK → inv_settlements
   account_balance integer
   market_value    integer          ← 庫存現值總額 Σ(shares × closing_price)
   deposit         integer
   withdrawal      integer
   profit_loss     integer

── inv_stock_holdings (股票庫存結餘)
   id              UUID PK
   year_month      varchar(7) FK → inv_stock_statements
   symbol          varchar(20)
   shares          integer
   closing_price   decimal(10,2)

── inv_stock_trades (股票交易紀錄，獨立於 settlement)
   id              UUID PK
   trade_date      date             ← 交易日期，查詢時用日期範圍篩選
   symbol          varchar(20)
   shares          integer          ← 正=買/負=賣
   price           decimal(10,2)
   fee             integer
   tax             integer
   note            text
```

## API

### Middleware: InvestmentAccess

- user_id 存在於 `inv_members` (active) → 放行
- `is_owner = true` → 完整操作權限
- 非 owner → 僅能 `GET /api/investments/allocations`（只回自己的資料）

### 成員管理 (owner only)

```
GET    /api/investments/members                        ← 列表
POST   /api/investments/members                        ← 新增
PUT    /api/investments/members/:id                    ← 更新 (name, active, sort_order)
```

### 月結算 (owner only)

```
GET    /api/investments/settlements                     ← 列表 (倒序)
POST   /api/investments/settlements                     ← 新增 draft
GET    /api/investments/settlements/:ym                 ← 詳情 (含各類型明細 + 試算分配)
PUT    /api/investments/settlements/:ym/complete        ← draft → completed (寫入 allocations)
PUT    /api/investments/settlements/:ym/reopen          ← completed → draft (清除 allocations)
DELETE /api/investments/settlements/:ym                 ← 刪除 (僅 draft)
```

完成驗證：所有已啟用投資類型皆有月結單才能 complete。

### 期貨月結單 (owner only)

```
PUT    /api/investments/settlements/:ym/futures         ← upsert
```

### 股票月結單 (owner only)

```
PUT    /api/investments/settlements/:ym/stocks          ← upsert (含 holdings)
```

### 成員入出金 (owner only)

```
GET    /api/investments/members/transactions?from=&to=  ← 列表 (日期 range 篩選, optional)
POST   /api/investments/members/transactions            ← 新增
PUT    /api/investments/members/transactions/:id        ← 更新
DELETE /api/investments/members/transactions/:id        ← 刪除
```

### 股票交易紀錄 (owner only)

```
GET    /api/investments/stocks/trades?from=&to=         ← 列表 (日期 range 篩選, optional)
POST   /api/investments/stocks/trades                   ← 新增
PUT    /api/investments/stocks/trades/:id               ← 更新
DELETE /api/investments/stocks/trades/:id               ← 刪除
```

### 個人損益紀錄 (所有成員)

```
GET    /api/investments/allocations?from=&to=           ← owner 回全部人，非 owner 只回自己
```

## 已決議

- 減碼當月重算：直接用減碼後的新權重算整個月，不按天數拆分
- 權重取整：floor（無條件���去）
- 擴充機制：主表固定，每次新增投資類型需開發子表 + 表單 + 損益公式
- 成員管理：可新增、可隱藏，不可刪除
- 月結算建立：手動建立（未來可加 cron 自動建 draft）
- 金額：台幣限定，integer 無小數
- 成員 ID：NanoID (varchar(21))
- 存取控制：靠 inv_members.user_id + is_owner 判斷
- completed ↔ draft 可雙向切換
- complete 時驗證所有已啟用投資類型皆有月結單
- 前期資料不存在時視為 0
- 成員入出金與投資帳戶入出金完全獨立
- 成員入出金/損益透過 inv_member_transactions 記錄明細，allocation 為加總快照
- 股票交易紀錄獨立於 settlement，不帶 FK
- 權重計算：floor((上期 balance - 當期 withdrawal) / 5000)，最小 1
- profit_loss 在 complete 時寫入 member_transactions（date 為該月最後一天），reopen 後重新 complete 時覆蓋更新
