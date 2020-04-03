import { Action, Dispatch } from 'redux'
import axios from 'axios'

type Lists = {
  Lists: Array<ServerList>
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

type ServerList = {
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

export const RequestGetLists = 'RequestGetLists' as const
export const ReceiveGetLists = 'ReceiveGetLists' as const

export const requestGetLists = () => ({
  type: RequestGetLists
})

export const receiveGetLists = (lists: Array<List>) => ({
  type: ReceiveGetLists,
  payload: lists
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
    })
  }
}

type Actions = ReturnType<typeof requestGetLists | typeof receiveGetLists>

export default Actions
