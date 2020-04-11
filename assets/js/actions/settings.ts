import { Dispatch, Action } from 'redux'
import axios from 'axios'

export const RequestUpdatePassword = 'RequestUpdatePassword' as const
export const ReceiveUpdatePassword = 'ReceiveUpdatePassword' as const

export const requestUpdatePassword = () => ({
  type: RequestUpdatePassword
})

export const receiveUpdatePassword = () => ({
  type: ReceiveUpdatePassword
})

export const updatePassword = (params: any) => {
  return (dispatch: Dispatch<Action>) => {
    dispatch(requestUpdatePassword())
    axios.patch<{}>('/settings/password', params).then(() => {
      dispatch(receiveUpdatePassword())
    })
  }
}

type Actions = ReturnType<typeof requestUpdatePassword | typeof receiveUpdatePassword>

export default Actions
