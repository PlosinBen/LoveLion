import { ref } from 'vue'
import { useApi } from '~/composables/useApi'
import type { ExpenseTemplate, ExpenseTemplateData } from '~/types'

export function useExpenseTemplates(spaceId: string) {
  const api = useApi()
  const templates = ref<ExpenseTemplate[]>([])
  const loading = ref(false)

  const fetchTemplates = async () => {
    loading.value = true
    try {
      templates.value = await api.get<ExpenseTemplate[]>(`/api/spaces/${spaceId}/expense-templates`)
    } catch (e) {
      console.error('Failed to fetch templates', e)
    } finally {
      loading.value = false
    }
  }

  const createTemplate = async (name: string, data: ExpenseTemplateData) => {
    return await api.post<ExpenseTemplate>(`/api/spaces/${spaceId}/expense-templates`, { name, data })
  }

  const deleteTemplate = async (templateId: string) => {
    await api.del(`/api/spaces/${spaceId}/expense-templates/${templateId}`)
    templates.value = templates.value.filter(t => t.id !== templateId)
  }

  return { templates, loading, fetchTemplates, createTemplate, deleteTemplate }
}
