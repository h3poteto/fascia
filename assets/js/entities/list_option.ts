export type ServerListOption = {
  ID: number
  Action: string
}

export type ListOption = {
  id: number
  action: string
}

export const converter = (s: ServerListOption): ListOption => ({
  id: s.ID,
  action: s.Action
})
