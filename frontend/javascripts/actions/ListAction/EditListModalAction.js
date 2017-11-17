import Request from 'superagent'

export const UNAUTHORIZED = 'UNAUTHORIZED'
function unauthorized() {
  window.location.pathname = '/sign_in'
  return {
    type: UNAUTHORIZED
  }
}

export const NOT_FOUND = 'NOT_FOUND'
function notFound() {
  return {
    type: NOT_FOUND
  }
}

export const SERVER_ERROR = 'SERVER_ERROR'
function serverError() {
  return {
    type: SERVER_ERROR
  }
}

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
    list: {ID: list.ID, ProjectID: list.ProjectID, Title: list.Title, Color: list.Color, ListTasks: list.ListTasks, ListOptionID: list.ListOptionID}
  }
}

export function fetchUpdateList(projectID, listID, params) {
  return dispatch => {
    dispatch(requestUpdateList())
    return Request
      .post(`/projects/${projectID}/lists/${listID}`)
      .type('form')
      .send(params)
      .end((err, res)=> {
        if(res.ok) {
          dispatch(receiveUpdateList(res.body))
        } else if (res.unauthorized) {
          dispatch(unauthorized())
        } else if (res.notFound) {
          dispatch(notFound())
        } else {
          dispatch(serverError())
        }
      })
  }
}


export const CHANGE_COLOR = 'CHANGE_COLOR'
export function changeColor(color) {
  return {
    type: CHANGE_COLOR,
    color: color,
  }
}
