<template>
  <div class="relative group cursor-pointer select-none" v-click-outside="() => showSwitcher = false">
    <!-- Trigger Header -->
    <div @click="showSwitcher = !showSwitcher" class="flex items-center gap-1.5 text-indigo-400 font-bold py-1 px-2 -ml-2 rounded-lg hover:bg-white/5 transition-colors">
       <span class="truncate max-w-xs">{{ currentSpace?.name || '正在載入...' }}</span>
       <Icon 
         icon="mdi:chevron-down" 
         class="text-lg flex-shrink-0 transition-transform duration-300" 
         :class="{ 'rotate-180': showSwitcher }" 
       />
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
      <div v-if="showSwitcher" class="absolute top-full left-0 mt-2 w-72 bg-neutral-900 border border-neutral-800 rounded-2xl shadow-2xl z-40 overflow-hidden py-1">
          <div v-for="space in allSpaces" :key="space.id" 
            class="px-4 py-3 hover:bg-neutral-800 transition-colors flex items-center justify-between group/item border-0 cursor-pointer"
            :class="{ 'bg-indigo-500/10 text-indigo-400': space.id === currentSpace?.id }"
            @click.stop="onSwitch(space.id)"
          >
            <div class="flex flex-col min-w-0 flex-1">
              <span class="font-medium truncate">{{ space.name }}</span>
              <span class="text-xs text-neutral-500">{{ space.creator_name || '系統' }} 的空間</span>
            </div>

            <div class="flex items-center gap-2 shrink-0 ml-3">
              <button 
                v-if="space.user_id === user?.id"
                @click.stop="router.push(`/spaces/${space.id}/settings`); showSwitcher = false"
                class="w-8 h-8 rounded-full flex items-center justify-center text-neutral-500 hover:bg-neutral-700 hover:text-white transition-colors border-0 cursor-pointer bg-transparent"
              >
                <Icon icon="mdi:cog-outline" class="text-lg" />
              </button>
              
              <Icon v-if="space.id === currentSpace?.id" icon="mdi:check" class="text-lg shrink-0" />
            </div>
          </div>
          
          <div class="border-t border-neutral-800 my-1"></div>
          
          <NuxtLink to="/spaces/add-new" class="w-full text-left px-4 py-3 hover:bg-neutral-800 text-sm text-neutral-400 flex items-center gap-2 no-underline" @click="showSwitcher = false">
            <Icon icon="mdi:plus-circle-outline" /> 新增空間
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
