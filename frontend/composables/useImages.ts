import { useApi } from './useApi'

export interface Image {
    id: string
    entity_id: string
    entity_type: string
    file_path: string
    blur_hash?: string
    sort_order: number
    created_at: string
}

export function useImages() {
    const { get, upload, put, del } = useApi()

    const getImages = async (entityId: string, entityType: string) => {
        return await get<Image[]>(`/api/images?entity_id=${entityId}&entity_type=${entityType}`)
    }

    const getImagesBatch = async (entityIds: string[], entityType: string) => {
        if (entityIds.length === 0) return []
        const ids = entityIds.join(',')
        return await get<Image[]>(`/api/images?entity_ids=${ids}&entity_type=${entityType}`)
    }

    const uploadImage = async (file: File, entityId: string, entityType: string) => {
        let fileToUpload = file

        // Compression Logic
        try {
            const { default: imageCompression } = await import('browser-image-compression')

            let options = {
                maxSizeMB: 1, // Default safety limit
                maxWidthOrHeight: 1280,
                useWebWorker: true,
                initialQuality: 0.9
            }

            if (entityType === 'transaction') {
                // Receipts: Focus on clarity, resize only if huge, keep high quality
                options.initialQuality = 1.0
                // If file is already small (<1.5MB), maybe skip compression to avoid any specific artifacting?
                // But browser-image-compression with quality 1.0 simply resizes if needed.
                // Let's stick to the plan: MaxWidth 1280px, No Quality Compression (1.0).
            } else {
                // Covers/Others: Standard Mobile First
                options.initialQuality = 0.9
            }

            // Only compress if it's an image
            if (file.type.startsWith('image/')) {
                fileToUpload = await imageCompression(file, options)
            }
        } catch (error) {
            console.warn('Image compression failed, uploading original:', error)
        }

        const formData = new FormData()
        // Ensure filename is preserved, even if compressed
        formData.append('file', fileToUpload, file.name)
        formData.append('entity_id', entityId)
        formData.append('entity_type', entityType)
        return await upload<Image>('/api/images', formData)
    }

    const deleteImage = async (id: string) => {
        return await del<{ message: string }>(`/api/images/${id}`)
    }

    const reorderImages = async (ids: string[]) => {
        return await put<{ message: string }>('/api/images/order', { ids })
    }

    const getImageUrl = (path: string) => {
        return path
    }

    return {
        getImages,
        getImagesBatch,
        uploadImage,
        deleteImage,
        reorderImages,
        getImageUrl
    }
}
