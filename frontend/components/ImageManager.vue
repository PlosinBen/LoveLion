<template>
  <div class="image-manager">
    <!-- Image Grid -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4" v-if="images.length > 0 || pendingUploads.length > 0">
      <div 
        v-for="(image, index) in [...images, ...pendingUploads.map(p => ({ id: 'pending-'+Math.random(), file_path: p.preview, isPending: true } as any))]" 
        :key="image.id"
        class="relative group aspect-square bg-gray-100 rounded-lg overflow-hidden border border-gray-200"
      >
        <img 
          :src="image.isPending ? image.file_path : getImageUrl(image.file_path)" 
          class="w-full h-full object-cover"
          alt="Uploaded image"
        />
        
        <!-- Overlay Actions -->
        <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
           <button 
             type="button"
             @click.stop="handleDelete(index)"
             class="p-2 bg-red-500 text-white rounded-full hover:bg-red-600 focus:outline-none"
             title="Delete"
           >
             <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
               <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
             </svg>
           </button>
           
           <div v-if="allowReorder && images.length > 1" class="flex gap-1">
             <button type="button" v-if="index > 0" @click.stop="move(index, -1)" class="p-1 bg-white text-gray-800 rounded shadow hover:bg-gray-100">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
             </button>
             <button type="button" v-if="index < images.length - 1" @click.stop="move(index, 1)" class="p-1 bg-white text-gray-800 rounded shadow hover:bg-gray-100">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                </svg>
             </button>
           </div>
        </div>
      </div>
    </div>

    <!-- Upload Area -->
    <div v-if="images.length < maxCount" class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center hover:border-blue-500 transition-colors cursor-pointer bg-gray-50" @click="triggerUpload">
      <input 
        type="file" 
        ref="fileInput" 
        class="hidden" 
        accept="image/png,image/jpeg,image/jpg"
        @change="handleFileSelect"
        :multiple="maxCount > 1"
      >
      <div v-if="uploading" class="text-blue-500 font-medium">Uploading...</div>
      <div v-else class="text-gray-500">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
        </svg>
        <span>Click to upload images</span>
        <div class="text-xs text-gray-400 mt-1">MAX 5MB (JPG, PNG)</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useImages, type Image } from '~/composables/useImages'

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
const fileInput = ref<HTMLInputElement | null>(null)
const uploading = ref(false)

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
    fileInput.value?.click()
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
