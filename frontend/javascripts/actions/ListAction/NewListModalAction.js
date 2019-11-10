import axios from 'axios'
import { ErrorHandler, ServerError } from '../ErrorHandler'
import { startLoading, stopLoading } from '../Loading'

export const CLOSE_NEW_LIST = 'CLOSE_NEW_LIST'
export function closeNewListModal() {
  return {
    type: CLOSE_NEW_LIST
  }
}

export const REQUEST_CREATE_LIST = 'REQUEST_CREATE_LIST'
function requestCreateList() {
  return {
    type: REQUEST_CREATE_LIST
  }
}

export const RECEIVE_CREATE_LIST = 'RECEIVE_CREATE_LIST'
function receiveCreateList(list) {
  return {
    type: RECEIVE_CREATE_LIST,
    list: {
      ID: list.ID,
      ProjectID: list.ProjectID,
      Title: list.Title,
      Color: list.Color,
      ListTasks: list.ListTasks
    }
  }
}

export function fetchCreateList(params) {
  return (dispatch, getState) => {
    const {
      ListReducer: {
        project: { ID }
      }
    } = getState()
    dispatch(startLoading())
    dispatch(requestCreateList())
    return axios
      .post(`/api/projects/${ID}/lists`, params)
      .then(res => {
        dispatch(stopLoading())
        dispatch(receiveCreateList(res.data))
      })
      .catch(err => {
        dispatch(stopLoading())
        ErrorHandler(err)
          .then()
          .catch(error => {
            dispatch(ServerError(error))
          })
      })
  }
}

export const CHANGE_COLOR = 'CHANGE_COLOR'
export function changeColor(color) {
  return {
    type: CHANGE_COLOR,
    color: color
  }
}
