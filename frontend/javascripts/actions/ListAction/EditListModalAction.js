import axios from 'axios'
import { ErrorHandler, ServerError } from '../ErrorHandler'
import { startLoading, stopLoading } from '../Loading'

export const CLOSE_EDIT_LIST = 'CLOSE_EDIT_LIST'
export function closeEditListModal() {
  return {
    type: CLOSE_EDIT_LIST,
    isListEditModalOpen: false
  }
}

export const REQUEST_UPDATE_LIST = 'REQUEST_UPDATE_LIST'
function requestUpdateList() {
  return {
    type: REQUEST_UPDATE_LIST
  }
}

export const RECEIVE_UPDATE_LIST = 'RECEIVE_UPDATE_LIST'
function receiveUpdateList(list) {
  return {
    type: RECEIVE_UPDATE_LIST,
    list: {
      ID: list.ID,
      ProjectID: list.ProjectID,
      Title: list.Title,
      Color: list.Color,
      ListTasks: list.ListTasks,
      ListOptionID: list.ListOptionID
    }
  }
}

export function fetchUpdateList(params) {
  return (dispatch, getState) => {
    const {
      ListReducer: {
        project: { ID: projectID }
      }
    } = getState()
    const {
      ListReducer: {
        selectedList: { ID: listID }
      }
    } = getState()
    dispatch(startLoading())
    dispatch(requestUpdateList())
    const form = Object.assign({}, params, {
      option_id: params.option_id.toString(10)
    })
    return axios
      .patch(`/api/projects/${projectID}/lists/${listID}`, form)
      .then(res => {
        dispatch(stopLoading())
        dispatch(receiveUpdateList(res.data))
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
