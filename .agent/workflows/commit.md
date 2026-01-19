---
description: How to commit changes with AI author attribution
---

# Git Commit (AI-assisted)

When committing code created or modified by Antigravity AI assistant.

> **IMPORTANT**: Always call commands as **separate tool calls** (not chained with `;` or `&&`). This ensures the Allow List whitelist works correctly.

## Commit with Antigravity Author
// turbo
1. Stage all changes:
```bash
git add .
```

// turbo
2. Commit with AI author:
```bash
git commit --author="Antigravity <PlosinBen@gmail.com>" -m "your commit message"
```

## Commit Specific Files
// turbo
3. Stage specific files:
```bash
git add <file1> <file2>
```

// turbo
4. Commit with AI author:
```bash
git commit --author="Antigravity <PlosinBen@gmail.com>" -m "your commit message"
```

## Amend Last Commit
// turbo
5. Amend the last commit to change author:
```bash
git commit --amend --author="Antigravity <PlosinBen@gmail.com>" --no-edit
```
