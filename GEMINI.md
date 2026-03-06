# LoveLion AI 核心規範 (Index)

## 🛠️ 環境執行標準 (MANDATORY)
- **換行規範**: **嚴禁使用 CRLF！** 所有檔案必須統一使用 **LF**。若偵測到 CRLF，視為嚴重違規，必須立即修正。
- **Git 提交**: 嚴禁在 `run_shell_command` 中使用 `&&`。必須拆分：
  1. `git add .`
  2. `bash bin/commit "標題" "描述"` (署名: Antigravity)
- **執行**: 宿主機無環境，指令必經 `docker compose exec backend/frontend`。

## 📚 知識索引
- **架構/ID**: [`.agent/architecture.md`](.agent/architecture.md) (NanoID/UUID/Money)
- **前端**: [`.agent/rules/frontend.md`](.agent/rules/frontend.md) (原生 Tailwind)
- **安全**: [`.agent/rules/security.md`](.agent/rules/security.md) (權限校驗/私隱)
- **工作流**: `.agent/workflows/` (DB/測試)

## ⚠️ 關鍵提醒
- 修改後必跑後端測試：`docker compose exec backend go test ./...`
- 變動模型後必驗 `cmd/seed` 編譯。
