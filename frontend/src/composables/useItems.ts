import { ref, type InjectionKey, type Ref } from 'vue'
import type { Item, ItemForm } from '../types'
import { api } from '../api/client'

export interface ItemsContext {
  items: Ref<Item[]>
  loading: Ref<boolean>
  filters: Ref<{ location: string; category: string; sort: string; search: string }>
  fetchItems: () => Promise<void>
  createItem: (form: ItemForm) => Promise<Item>
  updateItem: (id: number, form: ItemForm) => Promise<Item>
  deleteItem: (id: number) => Promise<void>
  getItem: (id: number) => Promise<Item>
}

export const ItemsKey: InjectionKey<ItemsContext> = Symbol('items')

export function useItems(): ItemsContext {
  const items = ref<Item[]>([])
  const loading = ref(false)
  const filters = ref({ location: '', category: '', sort: 'name', search: '' })

  async function fetchItems() {
    loading.value = true
    try {
      items.value = await api.listItems({
        location: filters.value.location || undefined,
        category: filters.value.category || undefined,
        sort: filters.value.sort || undefined,
        search: filters.value.search || undefined,
      })
    } finally {
      loading.value = false
    }
  }

  async function createItem(form: ItemForm): Promise<Item> {
    const item = await api.createItem(form)
    await fetchItems()
    return item
  }

  async function updateItem(id: number, form: ItemForm): Promise<Item> {
    const item = await api.updateItem(id, form)
    await fetchItems()
    return item
  }

  async function deleteItem(id: number): Promise<void> {
    await api.deleteItem(id)
    await fetchItems()
  }

  async function getItem(id: number): Promise<Item> {
    return api.getItem(id)
  }

  return { items, loading, filters, fetchItems, createItem, updateItem, deleteItem, getItem }
}
