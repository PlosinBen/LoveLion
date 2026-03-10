export function useButtonStyle(variant: string = 'primary') {
  const base = 'inline-flex justify-center items-center px-4 py-2.5 text-sm rounded font-bold cursor-pointer transition-all active:scale-95 disabled:opacity-50 disabled:pointer-events-none border-0 no-underline shrink-0'

  const variants: Record<string, string> = {
    primary: 'bg-indigo-500 text-white hover:bg-indigo-600 shadow-lg',
    secondary: 'bg-neutral-800 text-neutral-400 hover:text-white hover:bg-neutral-700',
    danger: 'bg-red-500/10 text-red-500 hover:bg-red-500 hover:text-white border border-red-500/20',
    ghost: 'bg-transparent text-neutral-500 hover:text-neutral-300',
  }

  return `${base} ${variants[variant] || variants.primary}`
}
