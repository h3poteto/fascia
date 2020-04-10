import axios from 'axios'
import { Action, Dispatch } from 'redux'

import { Task, ServerTask, converter } from '@/entities/task'

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
      const data = converter(res.data)
      dispatch(receiveGetTask(data))
    })
  }
}

type Actions = ReturnType<typeof requestGetTask | typeof receiveGetTask>

export default Actions
