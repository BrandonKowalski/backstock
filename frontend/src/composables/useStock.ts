import { api } from '../api/client'
import type { Stock, StockForm, StockMoveRequest } from '../types'

export function useStock() {
  async function listStock(itemId: number): Promise<Stock[]> {
    return api.listStock(itemId)
  }

  async function addStock(itemId: number, form: StockForm): Promise<Stock> {
    return api.addStock(itemId, form)
  }

  async function updateStock(stockId: number, quantity: number): Promise<Stock | void> {
    return api.updateStock(stockId, quantity)
  }

  async function deleteStock(stockId: number): Promise<void> {
    return api.deleteStock(stockId)
  }

  async function moveStock(stockId: number, req: StockMoveRequest): Promise<Stock | void> {
    return api.moveStock(stockId, req)
  }

  return { listStock, addStock, updateStock, deleteStock, moveStock }
}
