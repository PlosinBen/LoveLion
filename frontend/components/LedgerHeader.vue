<template>
  <div class="px-2 pt-0 pb-4 flex items-center gap-3">
    <!-- Back Button -->
    <button 
      v-if="showBack" 
      @click="router.back()" 
      class="w-10 h-10 rounded-full bg-neutral-800 text-white flex items-center justify-center hover:bg-neutral-700 transition-colors border-0 cursor-pointer shrink-0"
    >
      <Icon icon="mdi:arrow-left" class="text-xl" />
    </button>

    <!-- Title and Switcher -->
    <div class="flex-1 min-w-0">
      <h1 class="text-xl font-bold text-white tracking-tight truncate">{{ title }}</h1>
      <LedgerSwitcher />
    </div>
  </div>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import { useLedger } from '~/composables/useLedger'
import LedgerSwitcher from './LedgerSwitcher.vue'

defineProps({
  title: {
    type: String,
    required: true
  },
  icon: {
    type: String,
    default: 'mdi:wallet-outline'
  },
  showBack: {
    type: Boolean,
    default: false
  }
})

const router = useRouter()
const { currentLedger } = useLedger()
</script>
