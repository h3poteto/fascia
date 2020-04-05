import { Action, Dispatch } from 'redux'
import axios from 'axios'

type Lists = {
  Lists: Array<ServerList>
  NoneList: ServerList
}

type ServerTask = {
  ID: number
  ListID: number
  UserID: number
  IssueNumber: number
  Title: string
  Description: string
  HTMLURL: string
  PullRequest: boolean
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

export type Task = {
  id: number
  list_id: number
  user_id: number
  issue_number: number
  title: string
  description: string
  html_url: string
  pull_request: boolean
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

export type ServerProject = {
  ID: number
  UserID: number
  Title: string
  Description: string
  RepositoryID: number | null
  ShowIssues: boolean
  ShowPullRequests: boolean
}

export type Project = {
  id: number
  userID: number
  title: string
  description: string
  repositoryID: number | null
  showIssues: boolean
  showPullRequests: boolean
}

export const RequestGetLists = 'RequestGetLists' as const
export const ReceiveGetLists = 'ReceiveGetLists' as const
export const ReceiveNoneList = 'ReceiveNoneList' as const
export const OpenDelete = 'OpenDelete' as const
export const CloseDelete = 'CloseDelete' as const
export const OpenNewList = 'OpenNewList' as const
export const CloseNewList = 'CloseNewList' as const
export const OpenEditProject = 'OpenEditProject' as const
export const CloseEditProject = 'CloseEditProject' as const

export const requestGetLists = () => ({
  type: RequestGetLists
})

export const receiveGetLists = (lists: Array<List>) => ({
  type: ReceiveGetLists,
  payload: lists
})

export const receiveNoneList = (list: List) => ({
  type: ReceiveNoneList,
  payload: list
})

export const getLists = (projectID: number) => {
  return (dispatch: Dispatch<Action>) => {
    dispatch(requestGetLists())
    axios.get<Lists>(`/api/projects/${projectID}/lists`).then(res => {
      const data: Array<List> = res.data.Lists.map(l => ({
        id: l.ID,
        user_id: l.UserID,
        project_id: l.ProjectID,
        title: l.Title,
        color: l.Color,
        list_option_id: l.ListOptionID,
        is_hidden: l.IsHidden,
        is_init_list: l.IsInitList,
        tasks: l.ListTasks.map(t => ({
          id: t.ID,
          list_id: t.ListID,
          user_id: t.UserID,
          issue_number: t.IssueNumber,
          title: t.Title,
          description: t.Description,
          html_url: t.HTMLURL,
          pull_request: t.PullRequest
        }))
      }))
      dispatch(receiveGetLists(data))
      const none: List = {
        id: res.data.NoneList.ID,
        user_id: res.data.NoneList.UserID,
        project_id: res.data.NoneList.ProjectID,
        title: res.data.NoneList.Title,
        color: res.data.NoneList.Color,
        list_option_id: res.data.NoneList.ListOptionID,
        is_hidden: res.data.NoneList.IsHidden,
        is_init_list: res.data.NoneList.IsInitList,
        tasks: res.data.NoneList.ListTasks.map(t => ({
          id: t.ID,
          list_id: t.ListID,
          user_id: t.UserID,
          issue_number: t.IssueNumber,
          title: t.Title,
          description: t.Description,
          html_url: t.HTMLURL,
          pull_request: t.PullRequest
        }))
      }
      dispatch(receiveNoneList(none))
    })
  }
}

export const openDelete = () => ({
  type: OpenDelete
})

export const closeDelete = () => ({
  type: CloseDelete
})

export const openNewList = () => ({
  type: OpenNewList
})

export const closeNewList = () => ({
  type: CloseNewList
})

export const openEditProject = () => ({
  type: OpenEditProject
})

export const closeEditProject = () => ({
  type: CloseEditProject
})

type Actions = ReturnType<
  | typeof requestGetLists
  | typeof receiveGetLists
  | typeof receiveNoneList
  | typeof openDelete
  | typeof closeDelete
  | typeof openNewList
  | typeof closeNewList
  | typeof openEditProject
  | typeof closeEditProject
>

export default Actions
