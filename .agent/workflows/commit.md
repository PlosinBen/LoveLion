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
2. Set local user.name to Antigravity:
```bash
git config user.name "Antigravity"
```

// turbo
3. Commit with AI author:
```bash
git commit -m "your commit message"
```

// turbo
4. Restore user.name (unset to use global):
```bash
git config --unset user.name
```

## Commit Specific Files
// turbo
1. Stage specific files:
```bash
git add <file1> <file2>
```

// turbo
2. Set local user.name to Antigravity:
```bash
git config user.name "Antigravity"
```

// turbo
3. Commit with AI author:
```bash
git commit -m "your commit message"
```

// turbo
4. Restore user.name (unset to use global):
```bash
git config --unset user.name
```

## Amend Last Commit
// turbo
1. Set local user.name to Antigravity:
```bash
git config user.name "Antigravity"
```

// turbo
2. Amend the last commit to change author:
```bash
git commit --amend --reset-author --no-edit
```

// turbo
3. Restore user.name (unset to use global):
```bash
git config --unset user.name
```
