import Request from 'superagent'

export const UNAUTHORIZED = 'UNAUTHORIZED'
function unauthorized() {
  window.location.pathname = "/sign_in"
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

export const CLOSE_SHOW_TASK = 'CLOSE_SHOW_TASK'
export function closeShowTaskModal() {
  return {
    type: CLOSE_SHOW_TASK
  }
}

export const CHANGE_EDIT_MODE = 'CHANGE_EDIT_MODE'
export function changeEditMode(task) {
  return {
    type: CHANGE_EDIT_MODE,
    task: task
  }
}

export const REQUEST_UPDATE_TASK = 'REQUEST_UPDATE_TASK'
function requestUpdateTask() {
  return {
    type: REQUEST_UPDATE_TASK
  }
}

export const RECEIVE_UPDATE_TASK = 'RECEIVE_UPDATE_TASK'
function receiveUpdateTask(lists) {
  return {
    type: RECEIVE_UPDATE_TASK,
    lists: lists.Lists,
    noneList: lists.NoneList
  }
}

export function fetchUpdateTask(projectID, listID, taskID, params) {
  return dispatch => {
    dispatch(requestUpdateTask())
    return Request
      .post(`/projects/${projectID}/lists/${listID}/tasks/${taskID}`)
      .type('form')
      .send(params)
      .end((err, res)=> {
        if(res.ok) {
          dispatch(receiveUpdateTask(res.body))
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

export const REQUEST_DELETE_TASK = 'REQUEST_DELETE_TASK'
function requestDeleteTask() {
  return {
    type: REQUEST_DELETE_TASK
  }
}

export const RECEIVE_DELETE_TASK = 'RECEIVE_DELETE_TASK'
function receiveDeleteTask(lists) {
  return {
    type: RECEIVE_DELETE_TASK,
    lists: lists.Lists,
    noneList: lists.NoneList
  }
}

export function fetchDeleteTask(projectID, listID, taskID) {
  return dispatch => {
    dispatch(requestDeleteTask())
    return Request
      .del(`/projects/${projectID}/lists/${listID}/tasks/${taskID}`)
      .end((err, res)=> {
        if (res.ok) {
          dispatch(receiveDeleteTask(res.body))
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
