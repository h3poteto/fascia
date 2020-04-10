import { Action, Dispatch } from 'redux'
import axios from 'axios'

export const RequestUpdateProject = 'RequestUpdateProject' as const
export const ReceiveUpdateProject = 'ReceiveUpdateProject' as const

export const requestUpdateProject = () => ({
  type: RequestUpdateProject
})

export const receiveUpdateProject = () => ({
  type: ReceiveUpdateProject
})

export const updateProject = (id: number, params: any) => {
  return async (dispatch: Dispatch<Action>) => {
    dispatch(requestUpdateProject())
    return axios.patch<{}>(`/api/projects/${id}`, params).then(() => {
      dispatch(receiveUpdateProject())
    })
  }
}

type Actions = ReturnType<typeof requestUpdateProject | typeof receiveUpdateProject>

export default Actions
