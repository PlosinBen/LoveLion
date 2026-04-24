<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="visible" class="fixed inset-0 z-30 flex items-end justify-center">
        <div class="absolute inset-0 bg-black/60" @click="$emit('close')"></div>
        <div class="relative w-full max-w-lg bg-neutral-900 border-t border-neutral-800 rounded-t-2xl shadow-2xl max-h-3/4 flex flex-col">
          <!-- Header -->
          <div class="flex justify-between items-center p-4 border-b border-neutral-800 shrink-0">
            <span class="text-sm font-bold text-white">套用模板</span>
            <button
              @click="$emit('close')"
              class="text-neutral-500 hover:text-white bg-transparent border-0 cursor-pointer"
            >
              <Icon icon="mdi:close" class="text-xl" />
            </button>
          </div>

          <!-- Content -->
          <div class="flex-1 overflow-y-auto p-4">
            <div v-if="loading" class="text-center text-neutral-400 py-8">載入中...</div>
            <div v-else-if="templates.length === 0" class="text-center text-neutral-500 py-8">尚無模板</div>
            <div v-else class="flex flex-col gap-2">
              <div
                v-for="template in templates"
                :key="template.id"
                class="flex items-center justify-between bg-neutral-800 rounded-xl p-4 cursor-pointer hover:bg-neutral-700 transition-colors"
                @click="$emit('select', template)"
              >
                <div class="flex-1 min-w-0">
                  <div class="text-sm font-bold text-white truncate">{{ template.name }}</div>
                  <div class="text-xs text-neutral-400 mt-1 flex gap-2">
                    <span v-if="template.data.category">{{ template.data.category }}</span>
                    <span v-if="template.data.total_amount !== '0'">{{ template.data.currency }} {{ formatAmount(template.data.total_amount) }}</span>
                  </div>
                </div>
                <button
                  @click.stop="$emit('delete', template.id)"
                  class="text-neutral-600 hover:text-red-400 bg-transparent border-0 cursor-pointer p-2 transition-colors"
                >
                  <Icon icon="mdi:trash-can-outline" class="text-lg" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import type { ExpenseTemplate } from '~/types'

defineProps<{
  visible: boolean
  templates: ExpenseTemplate[]
  loading: boolean
}>()

defineEmits<{
  close: []
  select: [template: ExpenseTemplate]
  delete: [templateId: string]
}>()

const formatAmount = (amount: string) => {
  const num = parseFloat(amount)
  return isNaN(num) ? '0' : num.toLocaleString('zh-TW')
}
</script>
