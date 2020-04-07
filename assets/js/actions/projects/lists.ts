import { Action, Dispatch } from 'redux'
import axios from 'axios'

import { List, Lists, converter } from '@/entities/list'

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
      const data: Array<List> = res.data.Lists.map(l => converter(l))
      dispatch(receiveGetLists(data))
      const none = converter(res.data.NoneList)
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
