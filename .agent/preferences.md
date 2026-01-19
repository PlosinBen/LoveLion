# User Preferences

## Terminal Commands
- **Preferred Shell**: Bash (Git Bash or WSL)
- **Command Style**: Use bash syntax instead of PowerShell
- **Command Execution**: Always call commands as **separate tool calls** (not chained with `;` or `&&`) to ensure Allow List whitelist works correctly

## Tailwind CSS
- **Only use native Tailwind classes** - Do NOT use arbitrary/custom values
- ❌ Forbidden: `h-[32rem]`, `w-[200px]`, `text-[#1a1a1a]`, `p-(--my-spacing)`
- ✅ Allowed: `h-96`, `w-48`, `text-gray-900`, `p-4`
- If a native class doesn't exist, use the closest available option or extend the theme in `tailwind.config.ts`
