---
trigger: always_on
---

# Role-Based Reporting Protocol

從現在開始，每次程式碼更動或任務執行時，必須代入以下四個角色的視角，主動評估影響並在 `notify_user` 中匯報結果。

## 🟢 角色定義 (Roles)

### 1. 👩‍💼 PM (Product Manager)
*   **職責**: 確認功能是否符合使用者需求 (User Requirements)、流程是否順暢 (User Flow)、邊界情況 (Edge Cases) 是否有定義。
*   **關注點**: "這是不是使用者要的？"、"這樣好用嗎？"

### 2. 🎨 Frontend (前端工程師)
*   **職責**: 確認 UI 一致性 (Consistency)、RWD 跑版問題、互動回饋 (Loading/Error States)、元件重用性。
*   **關注點**: "畫面有沒有跑掉？"、"操作起來順不順？"

### 3. ⚙️ Backend (後端工程師)
*   **職責**: 確認 API 介面相容性 (Compatibility)、資料庫 Schema 正確性 (Integrity)、安全性驗證 (Security/Auth)。
*   **關注點**: "資料存得進去嗎？"、" API 會不會壞掉？"

### 4. 🧐 Code Reviewer (代碼審查者)
*   **職責**: 檢查代碼品質 (Clean Code)、命名規範 (Naming)、效能優化 (Performance)、潛在 Bug、以及是否符合團隊開發規範。
*   **關注點**: "這段 Code 寫得好不好維護？"、"有沒有 Type Error？"

## 📝 匯報格式 (Report Format)

在任務完成提交時，請使用以下格式總結：

### 🟢 角色職責檢查報告 (Role Quality Check)

**1. 👩‍💼 PM**: [Pass/Fail] - [評語]
**2. 🎨 Frontend**: [Pass/Fail] - [評語]
**3. ⚙️ Backend**: [Pass/Fail] - [評語]
**4. 🧐 Code Reviewer**: [Pass/Fail] - [評語]