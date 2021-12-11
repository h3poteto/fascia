import axios from 'axios'
import { Dispatch, Action } from 'redux'

import { List, ServerList, converter } from '@/entities/list'
import { ServerListOption, ListOption, converter as optionConverter } from '@/entities/list_option'

export const RequestGetList = 'RequestGetList' as const
export const ReceiveGetList = 'ReceiveGetList' as const
export const RequestGetListOptions = 'RequestGetListOptions' as const
export const ReceiveGetListOptions = 'ReceiveGetListOptions' as const

export const requestGetList = () => ({
  type: RequestGetList
})

export const receiveGetList = (list: List) => ({
  type: ReceiveGetList,
  payload: list
})

export const getList = (projectID: number, id: number) => {
  return (dispatch: Dispatch<Action>) => {
    dispatch(requestGetList())
    axios.get<ServerList>(`/api/projects/${projectID}/lists/${id}`).then((res) => {
      const data: List = converter(res.data)
      dispatch(receiveGetList(data))
    })
  }
}

export const requestGetListOptions = () => ({
  type: RequestGetListOptions
})

export const receiveGetListOptions = (listOptions: Array<ListOption>) => ({
  type: ReceiveGetListOptions,
  payload: listOptions
})

export const getListOptions = () => {
  return (dispatch: Dispatch<Action>) => {
    dispatch(requestGetListOptions())
    axios.get<Array<ServerListOption>>(`/api/list_options`).then((res) => {
      const data: Array<ListOption> = res.data.map((s) => optionConverter(s))
      dispatch(receiveGetListOptions(data))
    })
  }
}

type Actions = ReturnType<typeof requestGetList | typeof receiveGetList | typeof requestGetListOptions | typeof receiveGetListOptions>

export default Actions
