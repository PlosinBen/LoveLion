# User Preferences

## Terminal Commands
- **Preferred Shell**: Bash (Git Bash or WSL)
- **Command Style**: Use bash syntax instead of PowerShell
- **Command Execution**: Always call commands as **separate tool calls** (not chained with `;` or `&&`) to ensure Allow List whitelist works correctly

## Tailwind CSS
- **Only use native Tailwind classes** - Do NOT use arbitrary/custom values
- ❌ Forbidden: `h-[32rem]`, `w-[200px]`, `text-[#1a1a1a]`, `p-(--my-spacing)`
- ✅ Allowed: `h-96`, `w-48`, `text-gray-900`, `p-4`
- **Use Standard Tailwind Colors**: Do NOT define custom colors in `tailwind.config.ts`. Map design tokens (e.g., primary) to the nearest standard Tailwind utility (e.g., `indigo-500`, `neutral-900`).
- If a native class doesn't exist, use the closest available option. Do not extend the theme unless absolutely necessary for non-color values.
