<template>
  <div class="flex flex-col min-h-screen font-sans transition-colors duration-200" :class="[bgClass, textClass]">
    
    <!-- Optional Header Slot -->
    <div v-if="$slots.header" class="shrink-0 py-2" :class="headerClass">
      <slot name="header" />
    </div>

    <!-- Main Content -->
    <main class="flex-1 w-full mx-auto relative flex flex-col px-6">
      <slot />
    </main>

    <!-- Optional Footer Slot (BottomNav) -->
    <footer v-if="$slots.footer" class="shrink-0" :class="footerClass">
      <slot name="footer" />
    </footer>

  </div>
</template>

<script setup lang="ts">
interface Props {
  bgClass?: string
  textClass?: string
  headerClass?: string
  footerClass?: string
}

withDefaults(defineProps<Props>(), {
  bgClass: 'bg-neutral-900',
  textClass: 'text-neutral-50',
})
</script>
