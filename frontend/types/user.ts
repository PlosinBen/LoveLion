export interface User {
  id: string
  username: string
  display_name: string
  role: string
  inv_access?: boolean
  inv_is_owner?: boolean
  created_at: string
  updated_at: string
}

export interface Announcement {
  id: string
  title: string
  content: string
  status: string
  broadcast_start: string | null
  broadcast_end: string | null
  created_at: string
  updated_at: string
}
