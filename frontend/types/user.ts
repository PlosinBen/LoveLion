export interface User {
  id: string
  username: string
  display_name: string
  role: string
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
