import { ref } from 'vue'
import type { Category } from '../types'
import { api } from '../api/client'

export function useCategories() {
  const categories = ref<Category[]>([])

  async function fetchCategories() {
    categories.value = await api.listCategories()
  }

  async function createCategory(name: string, isFood = true) {
    await api.createCategory(name, isFood)
    await fetchCategories()
  }

  async function updateCategory(id: number, data: Partial<Category>) {
    await api.updateCategory(id, data)
    await fetchCategories()
  }

  async function deleteCategory(id: number) {
    await api.deleteCategory(id)
    await fetchCategories()
  }

  return { categories, fetchCategories, createCategory, updateCategory, deleteCategory }
}
