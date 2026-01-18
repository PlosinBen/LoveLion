---
description: How to commit changes with AI author attribution
---

# Git Commit (AI-assisted)

When committing code created or modified by Antigravity AI assistant.

## Commit with Antigravity Author
// turbo
1. Stage and commit with AI author:
```bash
git add .
git commit --author="Antigravity <PlosinBen@gmail.com>" -m "your commit message"
```

## Commit Specific Files
2. Commit specific files with AI author:
```bash
git add <file1> <file2>
git commit --author="Antigravity <PlosinBen@gmail.com>" -m "your commit message"
```

## Amend Last Commit
3. Amend the last commit to change author:
```bash
git commit --amend --author="Antigravity <antigravity@ai.assistant>" --no-edit
```
