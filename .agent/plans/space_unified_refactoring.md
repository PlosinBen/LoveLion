# LoveLion "大一統空間" 重構計畫 (Unified Space Refactoring)

## 1. 核心理念 (Core Concept)
*   **空間 (Space)** 是 App 的最高級主體。
*   所有的功能（記帳、比價、統計、成員管理）都作為「模組」掛載在空間內。
*   使用者透過 **「釘選 (Pin)」** 來切換目前的心理帳戶（例如：平時釘選「日常開銷」，出國時釘選「日本旅遊」）。

## 2. 功能模組 (Modules per Space)
每個空間進入後，透過分頁標籤 (Tabs) 切換：
*   **帳務 (Ledger)**：交易紀錄、多幣別換算、成員拆帳。
*   **比價 (Comparison)**：商店清單、產品價格紀錄（全域通用）。
*   **統計 (Stats)**：視覺化分析、每人支出比例。
*   **設定 (Settings)**：空間基礎資訊、成員管理、邀請連結。

## 3. UI/UX 變更點 (UX Design)
*   **首頁 (Dashboard)**：
    *   **釘選區 (Pinned)**：頂部大卡片或橫向滑動，快速進入目前焦點空間。
    *   **全部空間 (All Spaces)**：下方列表顯示歷史或非活躍空間。
*   **空間入口**：點擊空間卡片進入「帳務」模組。
*   **氛圍引擎**：根據 `Space.Type` (如 `personal`, `trip`, `group`) 調整 UI 配色與顯示欄位。

---

## 🚀 重構 Todo List (Todo List)

### 第一階段：後端模型大統 (Backend Refactoring)
- [ ] **模型合併與欄位擴充**：
    - [ ] 擴充 `Ledger` 模型，新增 `Type`, `StartDate`, `EndDate`, `Description`, `CoverImage`, `IsPinned` 等欄位。
    - [ ] 在代碼語意中將 `Ledger` 視為 `Space`。
- [ ] **資料遷移 (Data Migration)**：
    - [ ] 撰寫 Migration 將 `trips` 表資料搬遷至 `ledgers` 表。
    - [ ] 遷移 `trip_members` 到 `ledger_members`。
    - [ ] 更新 `comparison_stores` 的關聯鍵，從 `trip_id` 改為 `ledger_id`。
- [ ] **API 路由統一**：
    - [ ] 建立 `/api/spaces` 路由，整合舊有的 `/api/trips` 與 `/api/ledgers` 邏輯。

### 第二階段：前端核心重構 (Frontend Core)
- [ ] **Composables 升級**：
    - [ ] 將 `useLedger.ts` 升級為 `useSpace.ts`，實作全域空間管理與釘選邏輯。
- [ ] **首頁 Dashboard 重寫**：
    - [ ] 實作「釘選區」與「所有空間」列表介面。
- [ ] **通用空間佈局 (Unified Space View)**：
    - [ ] 建立 `/pages/spaces/[id]` 結構，內含 `Ledger`, `Comparison`, `Stats`, `Settings` 分頁。

### 第三階段：功能模組遷移與清理 (Modules & Cleanup)
- [ ] **功能移植**：
    - [ ] 將原有的比價、統計功能搬移至新的空間路徑。
    - [ ] 確保分帳功能在所有空間類型下皆正常運作。
- [ ] **代碼清理**：
    - [ ] 移除所有舊有的 `trips` 相關路由、組件與 Controller。
    - [ ] 根據空間類型優化卡片呈現視覺。

---

## 📈 最終目標 (Goal)
LoveLion 將不再只是一個記帳或旅遊工具，而是一個以「專案/空間」為導向的**全方位消費管理引擎**。
