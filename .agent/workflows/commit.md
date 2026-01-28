---
description: How to commit changes with AI author attribution
---

// turbo-all

# Git Commit (AI-assisted)

When committing code created or modified by Antigravity AI assistant.

> **IMPORTANT**: Always call commands as **separate tool calls** (not chained with `;` or `&&`). This ensures the Allow List whitelist works correctly.

## Commit with Antigravity Author

1. Stage all changes:
```bash
git add .
```

2. Set local user.name to Antigravity:
```bash
git config user.name "Antigravity"
```

3. Commit with AI author:
```bash
git commit -m "your commit message"
```

4. Restore user.name (unset to use global):
```bash
git config --unset user.name
```

## Commit Specific Files

1. Stage specific files:
```bash
git add <file1> <file2>
```

2. Set local user.name to Antigravity:
```bash
git config user.name "Antigravity"
```

3. Commit with AI author:
```bash
git commit -m "your commit message"
```

4. Restore user.name (unset to use global):
```bash
git config --unset user.name
```

## Amend Last Commit

1. Set local user.name to Antigravity:
```bash
git config user.name "Antigravity"
```

2. Amend the last commit to change author:
```bash
git commit --amend --reset-author --no-edit
```

3. Restore user.name (unset to use global):
```bash
git config --unset user.name
```
