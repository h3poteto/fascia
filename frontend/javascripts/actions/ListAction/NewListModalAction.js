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
    list: {ID: list.ID, ProjectID: list.ProjectID, Title: list.Title, Color: list.Color, ListTasks: list.ListTasks}
  }
}

export function fetchCreateList(projectID, params) {
  return dispatch => {
    dispatch(requestCreateList())
    return Request
      .post(`/projects/${projectID}/lists`)
      .type('form')
      .send(params)
      .end((err, res)=> {
        if(res.ok) {
          dispatch(receiveCreateList(res.body))
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
