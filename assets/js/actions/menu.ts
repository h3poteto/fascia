import { Dispatch, Action } from 'redux'
import axios from 'axios'

export const RequestLogout = 'RequestLogout' as const

export const requestLogout = () => ({
  type: RequestLogout
})

export const logout = () => {
  return async (dispatch: Dispatch<Action>) => {
    dispatch(requestLogout())
    return axios.delete<{}>(`/sign_out`).then(() => {
      window.location.pathname = '/sign_in'
    })
  }
}

type Actions = ReturnType<typeof requestLogout>

export default Actions
