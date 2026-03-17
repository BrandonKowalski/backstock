export function formatLocation(loc: string): string {
  return loc.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}

export function formatDate(date: string | null): string {
  if (!date) return ''
  const d = new Date(date + 'T00:00:00')
  return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}

export function formatQuantity(qty: number, abbreviation?: string): string {
  const formatted = Number.isInteger(qty) ? qty.toString() : qty.toFixed(1)
  return abbreviation ? `${formatted} ${abbreviation}` : formatted
}

export function isExpiringSoon(date: string | null, days = 7): boolean {
  if (!date) return false
  const exp = new Date(date + 'T00:00:00')
  const now = new Date()
  const diff = (exp.getTime() - now.getTime()) / (1000 * 60 * 60 * 24)
  return diff >= 0 && diff <= days
}

export function isExpired(date: string | null): boolean {
  if (!date) return false
  const exp = new Date(date + 'T00:00:00')
  return exp < new Date()
}

export function isPastBestBy(date: string | null): boolean {
  if (!date) return false
  const bb = new Date(date + 'T00:00:00')
  return bb < new Date()
}

export function isStaleStock(dateAdded: string | null, years = 1): boolean {
  if (!dateAdded) return false
  const added = new Date(dateAdded)
  const cutoff = new Date()
  cutoff.setFullYear(cutoff.getFullYear() - years)
  return added < cutoff
}
