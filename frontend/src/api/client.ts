import type { Item, ItemForm, Stock, StockForm, StockMoveRequest, Category, Unit, AuditEntry, Location } from '../types'
import { useToast } from '../composables/useToast'

const BASE = '/api'

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    const text = await res.text()
    const msg = text || res.statusText
    const { addToast } = useToast()
    addToast(msg)
    throw new Error(msg)
  }
  if (res.status === 204) return undefined as T
  return res.json()
}

export const api = {
  // Items
  listItems(params?: { location?: string; category?: string; sort?: string; search?: string }): Promise<Item[]> {
    const q = new URLSearchParams()
    if (params?.location) q.set('location', params.location)
    if (params?.category) q.set('category', params.category)
    if (params?.sort) q.set('sort', params.sort)
    if (params?.search) q.set('search', params.search)
    const qs = q.toString()
    return request<Item[]>(`/items${qs ? '?' + qs : ''}`)
  },
  getItem(id: number): Promise<Item> {
    return request<Item>(`/items/${id}`)
  },
  createItem(form: ItemForm): Promise<Item> {
    return request<Item>('/items', { method: 'POST', body: JSON.stringify(form) })
  },
  updateItem(id: number, form: ItemForm): Promise<Item> {
    return request<Item>(`/items/${id}`, { method: 'PUT', body: JSON.stringify(form) })
  },
  deleteItem(id: number): Promise<void> {
    return request<void>(`/items/${id}`, { method: 'DELETE' })
  },

  // Stock
  listStock(itemId: number): Promise<Stock[]> {
    return request<Stock[]>(`/items/${itemId}/stock`)
  },
  addStock(itemId: number, form: StockForm): Promise<Stock> {
    return request<Stock>(`/items/${itemId}/stock`, { method: 'POST', body: JSON.stringify(form) })
  },
  updateStock(stockId: number, quantity: number): Promise<Stock | void> {
    return request<Stock | void>(`/stock/${stockId}`, { method: 'PUT', body: JSON.stringify({ quantity }) })
  },
  deleteStock(stockId: number): Promise<void> {
    return request<void>(`/stock/${stockId}`, { method: 'DELETE' })
  },
  moveStock(stockId: number, req: StockMoveRequest): Promise<Stock | void> {
    return request<Stock | void>(`/stock/${stockId}/move`, { method: 'POST', body: JSON.stringify(req) })
  },

  // Categories
  listCategories(): Promise<Category[]> {
    return request<Category[]>('/categories')
  },
  createCategory(name: string, isFood = true): Promise<Category> {
    return request<Category>('/categories', { method: 'POST', body: JSON.stringify({ name, is_food: isFood }) })
  },
  updateCategory(id: number, data: Partial<Category>): Promise<Category> {
    return request<Category>(`/categories/${id}`, { method: 'PUT', body: JSON.stringify(data) })
  },
  deleteCategory(id: number): Promise<void> {
    return request<void>(`/categories/${id}`, { method: 'DELETE' })
  },

  // Units
  listUnits(): Promise<Unit[]> {
    return request<Unit[]>('/units')
  },
  createUnit(name: string, abbreviation: string): Promise<Unit> {
    return request<Unit>('/units', { method: 'POST', body: JSON.stringify({ name, abbreviation }) })
  },
  updateUnit(id: number, data: Partial<Unit>): Promise<Unit> {
    return request<Unit>(`/units/${id}`, { method: 'PUT', body: JSON.stringify(data) })
  },
  deleteUnit(id: number): Promise<void> {
    return request<void>(`/units/${id}`, { method: 'DELETE' })
  },

  // Locations
  listLocations(): Promise<Location[]> {
    return request<Location[]>('/locations')
  },
  createLocation(loc: { name: string; parent_id?: number | null; is_food?: boolean }): Promise<Location> {
    return request<Location>('/locations', { method: 'POST', body: JSON.stringify({ is_food: true, ...loc }) })
  },
  updateLocation(id: number, loc: Partial<Location>): Promise<Location> {
    return request<Location>(`/locations/${id}`, { method: 'PUT', body: JSON.stringify(loc) })
  },
  deleteLocation(id: number): Promise<void> {
    return request<void>(`/locations/${id}`, { method: 'DELETE' })
  },

  // Audit
  listAudit(): Promise<AuditEntry[]> {
    return request<AuditEntry[]>('/audit')
  },
}
