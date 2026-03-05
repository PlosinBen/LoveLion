<template>
  <ImmersiveHeader
    :fallback-icon="icon"
    height="h-48"
    class="rounded-3xl shadow-xl mb-6 mx-4"
  >
    <template #top-left>
      <button 
        v-if="showBack" 
        @click="router.back()" 
        class="w-10 h-10 rounded-full bg-black/30 backdrop-blur-md text-white flex items-center justify-center hover:bg-black/50 transition-colors border-0 cursor-pointer"
      >
        <Icon icon="mdi:arrow-left" class="text-xl" />
      </button>
    </template>

    <template #top-right>
      <button 
        v-if="showSettings && isOwner" 
        @click="router.push(`/ledger/${currentLedger?.id}/settings`)" 
        class="w-10 h-10 rounded-full bg-black/30 backdrop-blur-md text-white flex items-center justify-center hover:bg-black/50 transition-colors border-0 cursor-pointer"
      >
        <Icon icon="mdi:share-variant" class="text-xl" />
      </button>
    </template>

    <template #bottom>
      <div class="flex flex-col">
        <h1 class="text-2xl font-bold text-white shadow-sm">{{ title }}</h1>
        <LedgerSwitcher />
      </div>
    </template>
  </ImmersiveHeader>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Icon } from '@iconify/vue'
import { useLedger } from '~/composables/useLedger'
import { useAuth } from '~/composables/useAuth'
import ImmersiveHeader from './ImmersiveHeader.vue'
import LedgerSwitcher from './LedgerSwitcher.vue'

const props = defineProps({
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
  },
  showSettings: {
    type: Boolean,
    default: false
  }
})

const router = useRouter()
const { user } = useAuth()
const { currentLedger } = useLedger()

const isOwner = computed(() => {
  return currentLedger.value && currentLedger.value.user_id === user.value?.id
})
</script>
