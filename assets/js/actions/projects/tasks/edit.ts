import axios from 'axios'
import { push } from 'connected-react-router'

import { getLists } from '@/actions/projects/lists'
import { receiveCreateTask } from './new'

export const RequestUpdateTask = 'RequestUpdateTask' as const
export const ReceiveUpdateTask = 'ReceiveUpdateTask' as const

export const requestUpdateTask = () => ({
  type: RequestUpdateTask
})

export const receiveUpdateTask = () => ({
  type: ReceiveUpdateTask
})

export const updateTask = (projectID: number, listID: number, taskID: number, params: any) => {
  return async (dispatch: Function) => {
    dispatch(requestUpdateTask())
    return axios.patch<{}>(`/api/projects/${projectID}/lists/${listID}/tasks/${taskID}`, params).then(() => {
      dispatch(receiveCreateTask())
      dispatch(getLists(projectID))
      dispatch(push(`/projects/${projectID}/lists/${listID}/tasks/${taskID}`))
    })
  }
}

type Actions = ReturnType<typeof requestUpdateTask | typeof receiveUpdateTask>

export default Actions
