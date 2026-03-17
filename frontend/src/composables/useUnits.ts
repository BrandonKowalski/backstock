import { ref } from 'vue'
import type { Unit } from '../types'
import { api } from '../api/client'

export function useUnits() {
  const units = ref<Unit[]>([])

  async function fetchUnits() {
    units.value = await api.listUnits()
  }

  async function createUnit(name: string, abbreviation: string) {
    await api.createUnit(name, abbreviation)
    await fetchUnits()
  }

  async function updateUnit(id: number, data: Partial<Unit>) {
    await api.updateUnit(id, data)
    await fetchUnits()
  }

  async function deleteUnit(id: number) {
    await api.deleteUnit(id)
    await fetchUnits()
  }

  return { units, fetchUnits, createUnit, updateUnit, deleteUnit }
}
