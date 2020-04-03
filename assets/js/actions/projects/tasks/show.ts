import axios from 'axios'
import { Action, Dispatch } from 'redux'

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

export const RequestGetTask = 'RequestGetTask' as const
export const ReceiveGetTask = 'ReceiveGetTask' as const

export const requestGetTask = () => ({
  type: RequestGetTask
})

export const receiveGetTask = (task: Task) => ({
  type: ReceiveGetTask,
  payload: task
})

export const getTask = (projectID: number, listID: number, taskID: number) => {
  return (dispatch: Dispatch<Action>) => {
    dispatch(requestGetTask())
    axios.get<ServerTask>(`/api/projects/${projectID}/lists/${listID}/tasks/${taskID}`).then(res => {
      const data: Task = {
        id: res.data.ID,
        list_id: res.data.ListID,
        user_id: res.data.UserID,
        issue_number: res.data.IssueNumber,
        title: res.data.Title,
        description: res.data.Description,
        html_url: res.data.HTMLURL,
        pull_request: res.data.PullRequest
      }
      dispatch(receiveGetTask(data))
    })
  }
}

type Actions = ReturnType<typeof requestGetTask | typeof receiveGetTask>

export default Actions
