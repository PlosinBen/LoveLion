export interface ComparisonStore {
  id: string
  space_id: string
  name: string
  google_map_url: string
  location: string
  created_at: string
  updated_at: string
  products?: ComparisonProduct[]
}

export interface ComparisonProduct {
  id: string
  store_id: string
  name: string
  price: string
  currency: string
  unit: string
  note: string
  created_at: string
  updated_at: string
  store?: ComparisonStore
}
