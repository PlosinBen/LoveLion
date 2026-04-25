# Release 流程

## 方式

Release 綁定在 Git tag 上，不使用 `gh release create`。

## 步驟

1. 確認所有變更已 commit
2. Push master 到 origin：`git push origin master`
3. **Fetch remote tags 並確認 commit 連續性**：
   ```bash
   git fetch --tags origin
   git tag --sort=-v:refname | head -5          # 確認最新版本號
   git log --oneline $(git tag --sort=-v:refname | head -1)..HEAD  # 確認新 commit 接在最新 tag 之後
   ```
   若 `git log` 為空表示沒有新 commit；若最新 tag 不在當前 branch 歷史中會報錯 — 兩者都應停下來排查。
4. 根據最新 tag 決定下一個版本號
5. 建立 tag：`git tag v<版本號>`
6. Push tag：`git push origin v<版本號>`

## 版本號規則

採用 [Semantic Versioning](https://semver.org/)：

- **MAJOR**（v2.0.0）：不相容的 API 變更
- **MINOR**（v1.1.0）：新增功能，向下相容
- **PATCH**（v1.0.1）：修復 bug，向下相容
