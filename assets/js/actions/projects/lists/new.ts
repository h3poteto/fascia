import axios from 'axios'

import { List, ServerList, getLists } from '../lists'

export const RequestCreateList = 'RequestCreateList' as const
export const ReceiveCreateList = 'ReceiveCreateList' as const

export const requestCreateList = () => ({
  type: RequestCreateList
})

export const receiveCreateList = (list: List) => ({
  type: ReceiveCreateList,
  payload: list
})

export const createList = (projectID: number, params: any) => {
  return async (dispatch: Function) => {
    dispatch(requestCreateList())
    return axios.post<ServerList>(`/api/projects/${projectID}/lists`, params).then(res => {
      const data: List = {
        id: res.data.ID,
        user_id: res.data.UserID,
        project_id: res.data.ProjectID,
        title: res.data.Title,
        color: res.data.Color,
        list_option_id: res.data.ListOptionID,
        is_hidden: res.data.IsHidden,
        is_init_list: res.data.IsInitList,
        tasks: res.data.ListTasks.map(t => ({
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
      dispatch(receiveCreateList(data))
      dispatch(getLists(projectID))
    })
  }
}

type Actions = ReturnType<typeof requestCreateList | typeof receiveCreateList>

export default Actions
