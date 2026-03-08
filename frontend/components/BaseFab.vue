<template>
  <BaseButton
    @click="$emit('click')"
    variant="primary"
    size="icon"
    class="fixed right-6 bottom-24 z-20 !rounded-full !w-14 !h-14"
  >
    <Icon :icon="icon" class="text-3xl" />
  </BaseButton>
</template>

<script setup lang="ts">
import { Icon } from '@iconify/vue'
import BaseButton from '~/components/BaseButton.vue'

interface Props {
  icon?: string
}

const props = withDefaults(defineProps<Props>(), {
  icon: 'mdi:plus',
})

defineEmits(['click'])
</script>
