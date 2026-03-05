# 安全規範
- **校驗**: PUT/POST/DELETE 必驗 `LedgerMember` 權限。
- **歸屬**: 僅 `owner` 可改設定/增減成員/撤銷邀請。
- **洩漏**: 禁止返回 `PasswordHash` 或 Secrets。
- **事務**: 多表更新必用 DB Transaction。