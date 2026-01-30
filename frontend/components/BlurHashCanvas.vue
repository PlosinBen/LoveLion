<template>
  <canvas ref="canvasRef" :width="width" :height="height" class="w-full h-full object-cover"></canvas>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { decode } from 'blurhash'

const props = defineProps<{
  hash: string
  width?: number
  height?: number
  punch?: number
}>()

const canvasRef = ref<HTMLCanvasElement | null>(null)

const draw = () => {
  const canvas = canvasRef.value
  if (!canvas || !props.hash) return

  // Default resolution if not provided
  // We can render small and let CSS scale it up for better performance + blur effect
  const w = props.width || 32
  const h = props.height || 32

  try {
    const pixels = decode(props.hash, w, h, props.punch || 1)
    const ctx = canvas.getContext('2d')
    if (!ctx) return

    const imageData = ctx.createImageData(w, h)
    imageData.data.set(pixels)
    ctx.putImageData(imageData, 0, 0)
  } catch (e) {
    // console.warn("Failed to decode blurhash", e)
  }
}

onMounted(draw)
watch(() => props.hash, draw)
</script>
