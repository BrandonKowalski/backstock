import { ref, computed } from 'vue'
import type { Location } from '../types'
import { api } from '../api/client'

export function useLocations() {
  const locations = ref<Location[]>([])

  async function fetchLocations() {
    locations.value = await api.listLocations()
  }

  // Flat list of all location names (parents + children) for stock dropdowns
  const flatLocationNames = computed(() => {
    const names: string[] = []
    function walk(locs: Location[]) {
      for (const loc of locs) {
        names.push(loc.name)
        if (loc.children?.length) walk(loc.children)
      }
    }
    walk(locations.value)
    return names
  })

  async function createLocation(name: string, parentId?: number | null, isFood = true) {
    await api.createLocation({ name, parent_id: parentId ?? null, is_food: isFood })
    await fetchLocations()
  }

  async function updateLocation(id: number, data: Partial<Location>) {
    await api.updateLocation(id, data)
    await fetchLocations()
  }

  async function deleteLocation(id: number) {
    await api.deleteLocation(id)
    await fetchLocations()
  }

  function locationNamesByFood(isFood: boolean): string[] {
    const names: string[] = []
    function walk(locs: Location[]) {
      for (const loc of locs) {
        if (loc.is_food === isFood) names.push(loc.name)
        if (loc.children?.length) walk(loc.children)
      }
    }
    walk(locations.value)
    return names
  }

  const foodLocationNames = computed(() => locationNamesByFood(true))
  const nonFoodLocationNames = computed(() => locationNamesByFood(false))

  // Map of location name -> is_food for quick lookup
  const locationFoodMap = computed(() => {
    const map = new Map<string, boolean>()
    function walk(locs: Location[]) {
      for (const loc of locs) {
        map.set(loc.name, loc.is_food)
        if (loc.children?.length) walk(loc.children)
      }
    }
    walk(locations.value)
    return map
  })

  // Set of location names that are excluded from default views
  const excludedLocationNames = computed(() => {
    const names = new Set<string>()
    function walk(locs: Location[], parentExcluded: boolean) {
      for (const loc of locs) {
        const excluded = loc.exclude_default || parentExcluded
        if (excluded) names.add(loc.name)
        if (loc.children?.length) walk(loc.children, excluded)
      }
    }
    walk(locations.value, false)
    return names
  })

  // Map of location name -> set of names (self + all descendants) for filtering
  const locationWithChildren = computed(() => {
    const map = new Map<string, Set<string>>()
    function collect(loc: Location): string[] {
      const names = [loc.name]
      if (loc.children?.length) {
        for (const child of loc.children) {
          names.push(...collect(child))
        }
      }
      return names
    }
    function walk(locs: Location[]) {
      for (const loc of locs) {
        map.set(loc.name, new Set(collect(loc)))
        if (loc.children?.length) walk(loc.children)
      }
    }
    walk(locations.value)
    return map
  })

  return { locations, flatLocationNames, foodLocationNames, nonFoodLocationNames, locationFoodMap, excludedLocationNames, locationWithChildren, fetchLocations, createLocation, updateLocation, deleteLocation }
}
