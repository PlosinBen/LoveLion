import type { User } from './user'
import type { Image } from './image'

export interface Space {
  id: string
  user_id: string
  name: string
  description: string
  type: 'personal' | 'trip' | 'group'
  base_currency: string
  currencies: string[]
  split_members: string[]
  categories: string[]
  payment_methods: string[]
  start_date: string | null
  end_date: string | null
  cover_image: string
  is_pinned: boolean
  created_at: string
  updated_at: string
  my_role?: 'owner' | 'member'
  member_count?: number
  user?: User
  members?: Member[]
  invites?: Invite[]
  images?: Image[]
}

export interface Member {
  id: string
  space_id: string
  user_id: string
  role: 'owner' | 'member'
  alias: string
  weight?: number
  created_at: string
  updated_at: string
  user?: User
}

export interface InviteInfo {
  space_name: string
  creator_name: string
  is_one_time: boolean
}

export interface Invite {
  id: string
  space_id: string
  token: string
  is_one_time: boolean
  max_uses: number
  use_count: number
  expires_at: string | null
  created_by: string
  created_at: string
  updated_at: string
  creator?: User
}
