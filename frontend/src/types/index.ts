export interface Unit {
  id: number
  name: string
  abbreviation: string
  is_food: boolean | null
}

export interface Category {
  id: number
  name: string
  is_food: boolean
}

export interface Stock {
  id: number
  item_id: number
  location: string
  quantity: number
  date_added: string
  updated_at: string
}

export interface Item {
  id: number
  name: string
  is_food: boolean
  unit_id: number | null
  package_size: number | null
  expiration_date: string | null
  best_by_date: string | null
  low_quantity_threshold: number | null
  created_at: string
  updated_at: string
  categories: Category[]
  stock: Stock[]
  total_quantity: number
  unit: Unit | null
}

export interface ItemForm {
  name: string
  is_food: boolean
  unit_id: number | null
  package_size: number | null
  expiration_date: string | null
  best_by_date: string | null
  low_quantity_threshold: number | null
  category_ids: number[]
}

export interface StockForm {
  location: string
  quantity: number
}

export interface StockMoveRequest {
  to_location: string
  quantity: number
}

export interface AuditEntry {
  id: number
  item_name: string
  quantity: number
  created_at: string
}

export interface Location {
  id: number
  name: string
  parent_id: number | null
  is_food: boolean
  exclude_default: boolean
  children?: Location[]
}
