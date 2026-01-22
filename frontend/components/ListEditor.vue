<template>
  <div class="flex flex-col gap-2">
    <label class="text-sm text-neutral-400">{{ label }}</label>
    <div class="flex flex-wrap gap-2">
      <div v-for="(item, index) in modelValue" :key="index" class="flex items-center gap-1 bg-neutral-700 px-3 py-1 rounded-full text-sm">
        <span>{{ item }}</span>
        <button type="button" @click="remove(index)" class="text-neutral-400 hover:text-white">
          <Icon icon="mdi:close" class="text-base" />
        </button>
      </div>
      <div class="relative flex items-center">
        <input 
          v-model="newItem" 
          type="text" 
          class="bg-transparent border border-neutral-700 rounded-full px-3 py-1 text-sm text-white focus:outline-none focus:border-indigo-500 w-32 placeholder-neutral-600"
          :placeholder="placeholder"
          @keydown.enter.prevent="add"
          @blur="add"
        />
        <button type="button" @click="add" class="absolute right-2 text-neutral-400 hover:text-white" :class="{ 'hidden': !newItem }">
            <Icon icon="mdi:plus" class="text-base" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Icon } from '@iconify/vue'

const props = defineProps<{
  modelValue: string[]
  label: string
  placeholder?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string[]): void
}>()

const newItem = ref('')

const add = () => {
  const val = newItem.value.trim()
  if (val && !props.modelValue.includes(val)) {
    emit('update:modelValue', [...props.modelValue, val])
    newItem.value = ''
  } else {
      newItem.value = ''
  }
}

const remove = (index: number) => {
  const newValue = [...props.modelValue]
  newValue.splice(index, 1)
  emit('update:modelValue', newValue)
}
</script>
