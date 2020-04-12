import { Dispatch, Action } from 'redux'
import axios from 'axios'
import { ServerUser, User, converter } from '@/entities/user'

export const RequestGetSession = 'RequestGetSession' as const
export const ReceiveGetSession = 'ReceiveGetSession' as const
export const RequestUpdatePassword = 'RequestUpdatePassword' as const
export const ReceiveUpdatePassword = 'ReceiveUpdatePassword' as const

export const requestGetSession = () => ({
  type: RequestGetSession
})

export const receiveGetSession = (user: User) => ({
  type: ReceiveGetSession,
  payload: user
})

export const getSession = () => {
  return (dispatch: Dispatch<Action>) => {
    dispatch(requestGetSession())
    axios.get<ServerUser>('/session').then((res) => {
      const data: User = converter(res.data)
      dispatch(receiveGetSession(data))
    })
  }
}

export const requestUpdatePassword = () => ({
  type: RequestUpdatePassword
})

export const receiveUpdatePassword = () => ({
  type: ReceiveUpdatePassword
})

export const updatePassword = (params: any) => {
  return (dispatch: Dispatch<Action>) => {
    dispatch(requestUpdatePassword())
    axios.patch<{}>('/api/settings/password', params).then(() => {
      dispatch(receiveUpdatePassword())
    })
  }
}

type Actions = ReturnType<typeof requestUpdatePassword | typeof receiveUpdatePassword | typeof requestGetSession | typeof receiveGetSession>

export default Actions
