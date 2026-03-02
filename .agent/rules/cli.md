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
- **Commit Helper**: MUST use `bash bin/commit "message" ["detail"]` for all commits.
- **Auto-Config**: The script automatically handles `user.name "Antigravity"` and its cleanup.
- **Commit Style**: Concise messages. Primary goal: WHY, not WHAT.

## Command Execution
- Separate tool calls PREFERRED for tracking.
- If `&&` or `||` is used, `bash -c` is REQUIRED.