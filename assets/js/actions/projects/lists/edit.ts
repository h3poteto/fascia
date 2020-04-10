import axios from 'axios'
import { push } from 'connected-react-router'

import { getLists } from '@/actions/projects/lists'

export const RequestUpdateList = 'RequestUpdateList' as const
export const ReceiveUpdateList = 'ReceiveUpdateList' as const

export const requestUpdateList = () => ({
  type: RequestUpdateList
})

export const receiveUpdateList = () => ({
  type: ReceiveUpdateList
})

export const updateList = (projectID: number, listID: number, params: any) => {
  return (dispatch: Function) => {
    dispatch(requestUpdateList())
    return axios.patch<{}>(`/api/projects/${projectID}/lists/${listID}`, params).then(() => {
      dispatch(receiveUpdateList())
      dispatch(getLists(projectID))
      dispatch(push(`/projects/${projectID}`))
    })
  }
}

type Actions = ReturnType<typeof requestUpdateList | typeof receiveUpdateList>

export default Actions
