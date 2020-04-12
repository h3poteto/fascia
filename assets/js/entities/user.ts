export type ServerUser = {
  ID: number
  Email: string
  UserName: string
  Avatar: string
}

export type User = {
  id: number
  email: string
  user_name: string
  avatar: string
}

export const converter = (s: ServerUser): User => ({
  id: s.ID,
  email: s.Email,
  user_name: s.UserName,
  avatar: s.Avatar
})
