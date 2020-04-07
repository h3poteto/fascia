import axios from 'axios'

import { getLists } from '../lists'
import { List, ServerList, converter } from '@/entities/list'

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
      const data = converter(res.data)
      dispatch(receiveCreateList(data))
      dispatch(getLists(projectID))
    })
  }
}

type Actions = ReturnType<typeof requestCreateList | typeof receiveCreateList>

export default Actions
