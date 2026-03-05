# 系統規格
- **Stack**: Nuxt 4 (Layer), Go (Gin/GORM), PostgreSQL.
- **ID 策略**:
  - URL 暴露 (Trip/Txn/Ledger/Store): `NanoID`
  - 內部關聯 (User/Member/Item): `UUID`
- **精確度**: 金額 `DECIMAL(10,2)`, 匯率 `DECIMAL(12,6)`.
- **必備欄位**: 所有表均需 `created_at`, `updated_at` (TIMESTAMPTZ).
