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
export const RequestSyncGithub = 'RequestSyncGithub' as const
export const RequestMoveTask = 'RequestMoveTask' as const
export const RequestHideList = 'RequestHideList' as const
export const ReceiveHideList = 'ReceiveHideList' as const
export const RequestDisplayList = 'RequestDisplayList' as const
export const ReceiveDisplayList = 'ReceiveDisplayList' as const

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
    axios.get<Lists>(`/api/projects/${projectID}/lists`).then((res) => {
      const data: Array<List> = res.data.Lists.map((l) => converter(l))
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

export const requestSyncGithub = () => ({
  type: RequestSyncGithub
})

export const syncGithub = (projectID: number) => {
  return (dispatch: Dispatch<Action>) => {
    dispatch(requestSyncGithub())
    axios.post<Lists>(`/api/projects/${projectID}/fetch_github`).then((res) => {
      const data: Array<List> = res.data.Lists.map((l) => converter(l))
      dispatch(receiveGetLists(data))
      const none = converter(res.data.NoneList)
      dispatch(receiveNoneList(none))
    })
  }
}

export const requestMoveTask = () => ({
  type: RequestMoveTask
})

export const moveTask = (projectID: number, fromListID: number, toListID: number, taskID: number, prevToTaskID: number | null) => {
  return (dispatch: Dispatch<Action>) => {
    dispatch(requestMoveTask())
    axios
      .post<Lists>(`/api/projects/${projectID}/lists/${fromListID}/tasks/${taskID}/move_task`, {
        to_list_id: toListID,
        prev_to_task_id: prevToTaskID
      })
      .then((res) => {
        const data: Array<List> = res.data.Lists.map((l) => converter(l))
        dispatch(receiveGetLists(data))
        const none = converter(res.data.NoneList)
        dispatch(receiveNoneList(none))
      })
  }
}

export const requestHideList = () => ({
  type: RequestHideList
})

export const receiveHideList = () => ({
  type: ReceiveHideList
})

export const hideList = (projectID: number, listID: number) => {
  return (dispatch: Dispatch<Action>) => {
    dispatch(requestHideList())
    axios.patch<Lists>(`/api/projects/${projectID}/lists/${listID}/hide`).then((res) => {
      dispatch(receiveHideList())
      const data: Array<List> = res.data.Lists.map((l) => converter(l))
      dispatch(receiveGetLists(data))
      const none = converter(res.data.NoneList)
      dispatch(receiveNoneList(none))
    })
  }
}

export const requestDisplayList = () => ({
  type: RequestDisplayList
})

export const receiveDisplayList = () => ({
  type: ReceiveDisplayList
})

export const displayList = (projectID: number, listID: number) => {
  return (dispatch: Dispatch<Action>) => {
    dispatch(requestHideList())
    axios.patch<Lists>(`/api/projects/${projectID}/lists/${listID}/display`).then((res) => {
      dispatch(receiveDisplayList())
      const data: Array<List> = res.data.Lists.map((l) => converter(l))
      dispatch(receiveGetLists(data))
      const none = converter(res.data.NoneList)
      dispatch(receiveNoneList(none))
    })
  }
}

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
  | typeof requestSyncGithub
  | typeof requestMoveTask
  | typeof requestHideList
  | typeof receiveHideList
  | typeof requestDisplayList
  | typeof requestDisplayList
>

export default Actions
