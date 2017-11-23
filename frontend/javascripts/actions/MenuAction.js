import axios from 'axios'
import { ErrorHandler, ServerError } from './ErrorHandler'

const REQUEST_SIGN_OUT = 'REQUEST_SIGN_OUT'
function requestSignOut() {
  return {
    type: REQUEST_SIGN_OUT,
  }
}

export function signOut() {
  return dispatch => {
    dispatch(requestSignOut())
    return axios
      .delete('/sign_out')
      .then((_) => {
        window.location.pathname = '/sign_in'
      })
      .catch((err) => {
        ErrorHandler(err)
          .then()
          .catch((error) => {
            dispatch(ServerError(error))
          })
      })
  }
}
