<template>
  <div class="flex flex-col gap-3">
    <label class="text-sm text-neutral-400">{{ label }}</label>

    <div class="relative" @keydown.enter.prevent="add">
      <BaseInput
        v-model="newItem"
        :placeholder="placeholder"
        class="pr-10"
      />
      <button
        type="button"
        @click="add"
        class="absolute right-1 top-1/2 -translate-y-1/2 p-2 bg-transparent border-0 cursor-pointer text-neutral-500 hover:text-neutral-300 transition-all"
        :class="{ 'opacity-0 scale-75 pointer-events-none': !newItem, 'opacity-100 scale-100': newItem }"
      >
        <Icon icon="mdi:plus" class="text-xl" />
      </button>
    </div>

    <div v-if="modelValue.length" class="flex flex-wrap gap-2">
      <div v-for="(item, index) in modelValue" :key="index" class="flex items-center gap-1 bg-neutral-700 px-3 py-1 rounded text-sm text-white">
        <span>{{ item }}</span>
        <button type="button" @click="remove(index)" class="bg-transparent border-0 cursor-pointer text-neutral-500 hover:text-neutral-300 p-0 transition-colors">
          <Icon icon="mdi:close" class="text-base" />
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
  }
  newItem.value = ''
}

const remove = (index: number) => {
  const newValue = [...props.modelValue]
  newValue.splice(index, 1)
  emit('update:modelValue', newValue)
}
</script>
