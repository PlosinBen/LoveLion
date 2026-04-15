# Release 流程

## 方式

Release 綁定在 Git tag 上，不使用 `gh release create`。

## 步驟

1. 確認所有變更已 commit
2. Push master 到 origin：`git push origin master`
3. 建立 tag：`git tag v<版本號>`
4. Push tag：`git push origin v<版本號>`

## 版本號規則

採用 [Semantic Versioning](https://semver.org/)：

- **MAJOR**（v2.0.0）：不相容的 API 變更
- **MINOR**（v1.1.0）：新增功能，向下相容
- **PATCH**（v1.0.1）：修復 bug，向下相容
