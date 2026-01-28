---
trigger: always_on
---

# Terminal & CLI Rules

## Usage Preferences
- **Preferred Shell**: Bash (Git Bash or WSL)
- **Command Style**: Use bash syntax instead of PowerShell
- **Command Execution**: Always call commands as **separate tool calls** (not chained with `;` or `&&`) to ensure Allow List whitelist works correctly

## Docker Environment
- **Go Commands**: MUST be run inside the backend container.
  - Usage: `docker compose exec -T backend go ...`
- **NPM/Node Commands**: MUST be run inside the frontend container.
  - Usage: `docker compose exec -T frontend npm ...`
- **Host System**: Do NOT assume go, node, or npm are installed on the host Windows system.
