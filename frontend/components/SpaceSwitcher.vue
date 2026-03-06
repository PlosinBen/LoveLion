<template>
  <div class="relative group cursor-pointer select-none" v-click-outside="() => showSwitcher = false">
    <!-- Trigger Header -->
    <div @click="showSwitcher = !showSwitcher" class="flex items-center gap-2 text-indigo-400 font-black py-1 px-3 -ml-3 rounded-2xl hover:bg-indigo-500/5 transition-all active:scale-95">
       <span class="truncate max-w-[120px] sm:max-w-[200px] tracking-tight">{{ currentSpace?.name || '正在載入...' }}</span>
       <Icon 
         icon="mdi:chevron-down" 
         class="text-xl flex-shrink-0 transition-transform duration-300" 
         :class="{ 'rotate-180': showSwitcher }" 
       />
    </div>
    
    <!-- Dropdown List -->
    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="transform scale-95 opacity-0 -translate-y-4"
      enter-to-class="transform scale-100 opacity-100 translate-y-0"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="transform scale-100 opacity-100 translate-y-0"
      leave-to-class="transform scale-95 opacity-0 -translate-y-4"
    >
      <div v-if="showSwitcher" class="absolute top-full left-0 mt-3 w-72 bg-neutral-900 border border-neutral-800 rounded-3xl shadow-2xl z-[100] overflow-hidden py-2 backdrop-blur-xl">
          <div v-for="space in allSpaces" :key="space.id" 
            class="px-5 py-4 hover:bg-neutral-800/80 transition-all flex items-center justify-between group/item border-0 cursor-pointer"
            :class="{ 'bg-indigo-500/10': space.id === currentSpace?.id }"
            @click.stop="onSwitch(space.id)"
          >
            <div class="flex flex-col min-w-0 flex-1">
              <span class="font-bold truncate tracking-tight" :class="space.id === currentSpace?.id ? 'text-indigo-400' : 'text-neutral-200'">{{ space.name }}</span>
              <span class="text-[10px] text-neutral-500 font-black uppercase tracking-widest mt-0.5">@{{ space.creator_name || '系統' }}</span>
            </div>

            <div class="flex items-center gap-3 shrink-0 ml-4">
              <!-- Settings Link -->
              <button 
                v-if="space.user_id === user?.id"
                @click.stop="router.push(`/spaces/${space.id}/settings`); showSwitcher = false"
                class="w-10 h-10 rounded-xl flex items-center justify-center text-neutral-500 hover:bg-neutral-700 hover:text-white transition-all border-0 cursor-pointer bg-transparent"
              >
                <Icon icon="mdi:cog-outline" class="text-xl" />
              </button>
              
              <Icon v-if="space.id === currentSpace?.id" icon="mdi:check-circle" class="text-xl text-indigo-500 shrink-0" />
            </div>
          </div>
          
          <div class="h-px bg-neutral-800 mx-5 my-2"></div>
          
          <NuxtLink to="/spaces/add-new" class="mx-2 px-4 py-4 hover:bg-neutral-800/50 rounded-2xl text-xs font-black text-indigo-400 flex items-center gap-3 no-underline transition-colors active:bg-neutral-800" @click="showSwitcher = false">
            <div class="w-8 h-8 rounded-lg bg-indigo-500/10 flex items-center justify-center">
              <Icon icon="mdi:plus-circle-outline" class="text-lg" />
            </div>
            新增空間
          </NuxtLink>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Icon } from '@iconify/vue'
import { useSpace } from '~/composables/useSpace'
import { useAuth } from '~/composables/useAuth'

const router = useRouter()
const { user } = useAuth()
const { allSpaces, currentSpace, selectSpace } = useSpace()
const showSwitcher = ref(false)

const onSwitch = (id: string) => {
  selectSpace(id)
  showSwitcher.value = false
  router.push(`/spaces/${id}`)
}
</script>
