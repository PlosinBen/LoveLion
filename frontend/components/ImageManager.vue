<template>
  <div class="image-manager">
    <!-- Image Grid -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4" v-if="images.length > 0 || pendingUploads.length > 0">
      <ImageGridItem
        v-for="(image, index) in allImages"
        :key="image.id"
        :image="image"
        :index="index"
        :allow-reorder="allowReorder && allImages.length > 1"
        :is-first="index === 0"
        :is-last="index === allImages.length - 1"
        @delete="handleDelete"
        @move="move"
      />
    </div>

    <!-- Upload Area -->
    <ImageUploadArea
      v-if="images.length < maxCount"
      :uploading="uploading"
      :max-count="maxCount"
      @trigger="triggerUpload"
      @file-select="handleFileSelect"
      ref="uploadAreaRef"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useImages, type Image } from '~/composables/useImages'
import ImageGridItem from '~/components/ImageManager/ImageGridItem.vue'
import ImageUploadArea from '~/components/ImageManager/ImageUploadArea.vue'

const props = defineProps({
  entityId: {
    type: String,
    required: true
  },
  entityType: {
    type: String,
    required: true
  },
  maxCount: {
    type: Number,
    default: Infinity
  },
  allowReorder: {
    type: Boolean,
    default: true
  },
  instantDelete: {
    type: Boolean,
    default: true
  },
  instantUpload: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['change'])
defineExpose({ commit })

const { getImages, uploadImage, deleteImage, reorderImages, getImageUrl } = useImages()

const images = ref<Image[]>([])
const pendingDeletes = ref<string[]>([])
const pendingUploads = ref<{ file: File, preview: string }[]>([])
const uploadAreaRef = ref(null)
const uploading = ref(false)

const allImages = computed(() => {
  return [
    ...images.value,
    ...pendingUploads.value.map(p => ({
      id: 'pending-' + Math.random(),
      file_path: p.preview,
      isPending: true
    } as any))
  ]
})

const loadImages = async () => {
    if (!props.entityId) return
    try {
        images.value = await getImages(props.entityId, props.entityType)
    } catch (e) {
        console.error("Failed to load images", e)
    }
}

// Watch for entityId changes (e.g. navigation)
watch(() => props.entityId, loadImages)

const triggerUpload = () => {
    if (uploadAreaRef.value && 'click' in uploadAreaRef.value) {
        (uploadAreaRef.value as any).click()
    }
}

const handleFileSelect = async (event: Event) => {
    const target = event.target as HTMLInputElement
    if (!target.files || target.files.length === 0) return

    uploading.value = true
    try {
        if (!target.files) return
        for (let i = 0; i < target.files.length; i++) {
            if (images.value.length + pendingUploads.value.length >= props.maxCount) break
            const file = target.files[i]
            if (file) {
                 if (props.instantUpload) {
                    await uploadImage(file, props.entityId, props.entityType)
                 } else {
                    // Deferred: Create preview and add to pending
                    const preview = URL.createObjectURL(file)
                    pendingUploads.value.push({ file, preview })
                 }
            }
        }
        if (props.instantUpload) {
            await loadImages()
        }
    } catch (e: any) {
        alert("Upload failed: " + (e.message || "Unknown error"))
    } finally {
        uploading.value = false
        target.value = '' 
    }
}

async function handleDelete(index: number) {
    if (!confirm("Confirm remove?")) return
    
    // Check if it's a pending upload
    // If index is within images array (existing images)
    if (index < images.value.length) {
        const image = images.value[index]
        if (!image) return

        if (props.instantDelete) {
            try {
                await deleteImage(image.id)
                images.value.splice(index, 1)
            } catch (e) {
                alert("Delete failed")
            }
        } else {
            images.value.splice(index, 1)
            pendingDeletes.value.push(image.id)
            emit('change', images.value)
        }
    } else {
        // It's a pending upload
        const pendingIndex = index - images.value.length
        const pending = pendingUploads.value[pendingIndex]
        if (pending) {
            URL.revokeObjectURL(pending.preview)
            pendingUploads.value.splice(pendingIndex, 1)
        }
    }
}

async function commit(overrideEntityId?: string) {
    const targetEntityId = overrideEntityId || props.entityId

    // Process pending deletes
    if (pendingDeletes.value.length > 0) {
        await Promise.all(pendingDeletes.value.map(id => deleteImage(id)))
        pendingDeletes.value = []
    }
    
    // Process pending uploads
    if (pendingUploads.value.length > 0) {
        await Promise.all(pendingUploads.value.map(p => uploadImage(p.file, targetEntityId, props.entityType)))
        // Clear pending
        pendingUploads.value.forEach(p => URL.revokeObjectURL(p.preview))
        pendingUploads.value = []
    }
}

const move = async (index: number, direction: number) => {
    if (index + direction < 0 || index + direction >= images.value.length) return

    const newImages = [...images.value]
    const deleted = newImages.splice(index, 1)
    if (deleted.length > 0) {
        newImages.splice(index + direction, 0, deleted[0])
        images.value = newImages 
        
        try {
            await reorderImages(newImages.map(img => img.id))
        } catch (e) {
            console.error(e)
            await loadImages() 
        }
    }
}

onMounted(() => {
    loadImages()
})
</script>
