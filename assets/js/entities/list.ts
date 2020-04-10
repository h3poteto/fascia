import { ServerTask, Task, converter as taskConverter } from './task'

export type Lists = {
  Lists: Array<ServerList>
  NoneList: ServerList
}

export type ServerList = {
  ID: number
  UserID: number
  ProjectID: number
  Title: string
  Color: string
  ListOptionID: number
  IsHidden: boolean
  IsInitList: boolean
  ListTasks: Array<ServerTask>
}

export type List = {
  id: number
  user_id: number
  project_id: number
  title: string
  color: string
  list_option_id: number
  is_hidden: boolean
  is_init_list: boolean
  tasks: Array<Task>
}

export const converter = (l: ServerList): List => ({
  id: l.ID,
  user_id: l.UserID,
  project_id: l.ProjectID,
  title: l.Title,
  color: l.Color,
  list_option_id: l.ListOptionID,
  is_hidden: l.IsHidden,
  is_init_list: l.IsInitList,
  tasks: l.ListTasks.map(t => taskConverter(t))
})
