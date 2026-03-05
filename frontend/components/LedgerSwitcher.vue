<template>
  <div class="relative group cursor-pointer" v-click-outside="() => showSwitcher = false">
    <div @click="showSwitcher = !showSwitcher" class="flex items-center gap-1.5 text-indigo-400 font-medium py-1 px-2 -ml-2 rounded-lg hover:bg-white/5 transition-colors">
       <span class="truncate max-w-[160px]">{{ currentLedger?.name || '正在載入...' }}</span>
       <Icon icon="mdi:chevron-down" class="text-lg flex-shrink-0 transition-transform duration-300" :class="{ 'rotate-180': showSwitcher }" />
    </div>
    
    <!-- Dropdown List -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="transform scale-95 opacity-0 -translate-y-2"
      enter-to-class="transform scale-100 opacity-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="transform scale-100 opacity-100 translate-y-0"
      leave-to-class="transform scale-95 opacity-0 -translate-y-2"
    >
      <div v-if="showSwitcher" class="absolute top-full left-0 mt-2 w-72 bg-neutral-900 border border-neutral-800 rounded-2xl shadow-2xl z-50 overflow-hidden py-1">
          <div v-for="l in allLedgers" :key="l.id" 
            class="px-4 py-3 hover:bg-neutral-800 transition-colors flex items-center justify-between group/item"
            :class="{ 'bg-indigo-500/10 text-indigo-400': l.id === currentLedger?.id }"
            @click.stop="onSwitch(l.id)"
          >
            <div class="flex flex-col min-w-0 flex-1">
              <span class="font-medium truncate">{{ l.name }}</span>
              <span class="text-[10px] text-neutral-500">{{ l.user?.display_name || '系統' }} 的帳本</span>
            </div>

            <div class="flex items-center gap-2 shrink-0 ml-3">
              <!-- Settings Button (Visible only if owner) -->
              <button 
                v-if="l.user_id === user?.id"
                @click.stop="router.push(`/ledger/${l.id}/settings`); showSwitcher = false"
                class="w-8 h-8 rounded-full flex items-center justify-center text-neutral-500 hover:bg-neutral-700 hover:text-white transition-colors border-0 cursor-pointer bg-transparent"
              >
                <Icon icon="mdi:cog-outline" class="text-lg" />
              </button>
              
              <Icon v-if="l.id === currentLedger?.id" icon="mdi:check" class="text-lg shrink-0" />
            </div>
          </div>
          
          <div class="border-t border-neutral-800 my-1"></div>
          
          <NuxtLink to="/ledger/add-new" class="w-full text-left px-4 py-3 hover:bg-neutral-800 text-sm text-neutral-400 flex items-center gap-2 no-underline" @click="showSwitcher = false">
            <Icon icon="mdi:plus-circle-outline" /> 新增帳本
          </NuxtLink>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Icon } from '@iconify/vue'
import { useLedger } from '~/composables/useLedger'
import { useAuth } from '~/composables/useAuth'

const router = useRouter()
const { user } = useAuth()
const { allLedgers, currentLedger, selectLedger } = useLedger()
const showSwitcher = ref(false)

const onSwitch = (id: string) => {
  selectLedger(id)
  showSwitcher.value = false
}
</script>
