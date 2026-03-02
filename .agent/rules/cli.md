---
trigger: always_on
---

# CLI & Git Mandatory Rules

## Execution Guidelines
- **Shell**: MANDATORY use of Git Bash (`sh/bash`).
- **Chained Commands**: MUST be wrapped in `bash -c '...'` (e.g., `bash -c "git add . && git commit -m '...'"`).
- **Environment**: 
  - Backend: `docker compose exec -T backend go ...`
  - Frontend: `docker compose exec -T frontend npm ...`
- **Host System**: NEVER assume Go/Node presence on host.

## AI Contribution Attribution
- **Author Identity**: BEFORE any commit, MUST set `git config user.name "Antigravity"`.
- **Cleanup**: AFTER commit, MUST run `git config --unset user.name`.
- **Commit Style**: Concise messages. Primary goal: WHY, not WHAT.

## Command Execution
- Separate tool calls PREFERRED for tracking.
- If `&&` or `||` is used, `bash -c` is REQUIRED.