<template>
  <div class="media-view-page bg-black fixed inset-0 flex flex-col items-center justify-center">
    <div class="absolute top-6 left-6 z-10">
      <BaseButton @click="router.back()" variant="ghost" class="!p-0 w-10 h-10 rounded-full backdrop-blur-md">
        <Icon icon="mdi:close" class="text-xl" />
      </BaseButton>
    </div>

    <div v-if="loading" class="text-neutral-500 flex flex-col items-center gap-3">
        <Icon icon="mdi:loading" class="text-4xl animate-spin" />
        <span class="text-xs font-bold uppercase tracking-widest">載入相片...</span>
    </div>
    <div v-else class="w-full h-full flex items-center justify-center">
      <img :src="imageUrl" class="max-w-full max-h-full object-contain shadow-2xl" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useImages } from '~/composables/useImages'
import BaseButton from '~/components/BaseButton.vue'

definePageMeta({
  layout: 'empty'
})

const route = useRoute()
const router = useRouter()
const { getImageUrl } = useImages()

const imageUrl = ref('')
const loading = ref(true)

onMounted(() => {
  const path = route.query.path as string
  if (path) {
    imageUrl.value = getImageUrl(path)
  }
  loading.value = false
})
</script>
