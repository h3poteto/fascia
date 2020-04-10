import axios from 'axios'
import { push } from 'connected-react-router'

import { getLists } from '@/actions/projects/lists'

export const RequestCreateTask = 'RequestCreateTask' as const
export const ReceiveCreateTask = 'ReceiveCreateTask' as const

export const requestCreateTask = () => ({
  type: RequestCreateTask
})

export const receiveCreateTask = () => ({
  type: ReceiveCreateTask
})

export const createTask = (projectID: number, listID: number, params: any) => {
  return async (dispatch: Function) => {
    dispatch(requestCreateTask())
    return axios.post<{}>(`/api/projects/${projectID}/lists/${listID}/tasks`, params).then(() => {
      dispatch(receiveCreateTask())
      dispatch(getLists(projectID))
      dispatch(push(`/projects/${projectID}`))
    })
  }
}

type Actions = ReturnType<typeof requestCreateTask | typeof receiveCreateTask>

export default Actions
