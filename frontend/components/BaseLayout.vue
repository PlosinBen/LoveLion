<template>
  <div class="flex min-h-screen min-w-0 w-full flex-col overflow-x-hidden font-sans transition-colors duration-200" :class="[bgClass, textClass]">
    
    <!-- Optional Header Slot -->
    <div v-if="$slots.header" class="shrink-0 py-2" :class="headerClass">
      <slot name="header" />
    </div>

    <!-- Main Content -->
    <!-- Auto-apply bottom padding if footer (BottomNav) exists to prevent overlap -->
    <main 
      class="relative mx-auto flex w-full min-w-0 max-w-full flex-1 flex-col px-4"
      :class="[{ 'pb-32': $slots.footer }, mainClass]"
    >
      <slot />
    </main>

    <!-- Optional Footer Slot (BottomNav) -->
    <footer v-if="$slots.footer" class="fixed inset-x-0 bottom-0 z-30 pointer-events-none" :class="footerClass">
      <div class="pointer-events-auto">
        <slot name="footer" />
      </div>
    </footer>

    <!-- Global UI overlays (client-only to avoid SSR hydration mismatch from Teleport) -->
    <ClientOnly>
      <AppToast />
      <AppConfirm />
      <AppPrompt />
    </ClientOnly>
  </div>
</template>

<script setup lang="ts">
interface Props {
  bgClass?: string
  textClass?: string
  headerClass?: string
  footerClass?: string
  mainClass?: string
}

withDefaults(defineProps<Props>(), {
  bgClass: 'bg-neutral-900',
  textClass: 'text-neutral-50',
})
</script>
