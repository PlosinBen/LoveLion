<template>
  <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-blue-500 transition-colors cursor-pointer bg-gray-50" @click="$emit('trigger')">
    <input 
      type="file" 
      ref="fileInputRef" 
      class="hidden" 
      accept="image/png,image/jpeg,image/jpg"
      @change="handleFileSelect"
      :multiple="maxCount > 1"
    >
    <div v-if="uploading" class="text-blue-500 font-medium">Uploading...</div>
    <div v-else class="text-gray-500">
      <Icon icon="mdi:cloud-upload" class="h-8 w-8 mx-auto mb-2" />
      <span>Click to upload images</span>
      <div class="text-xs text-gray-400 mt-1">MAX 5MB (JPG, PNG)</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { Icon } from '@iconify/vue'

const props = defineProps<{
  uploading: boolean
  maxCount: number
}>()

const emit = defineEmits(['trigger', 'fileSelect'])

const fileInputRef = ref<HTMLInputElement | null>(null)

defineExpose({
  click: () => fileInputRef.value?.click()
})

const handleFileSelect = (event: Event) => {
  emit('fileSelect', event)
}
</script>
