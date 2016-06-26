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

export const CLOSE_NEW_TASK = 'CLOSE_NEW_TASK'
export function closeNewTaskModal() {
  return {
    type: CLOSE_NEW_TASK,
    isTaskModalOpen: false
  }
}

export const UPDATE_NEW_TASK_TITLE = 'UPDATE_NEW_TASK_TITLE'
export function updateNewTaskTitle(ev) {
  return {
    type: UPDATE_NEW_TASK_TITLE,
    title: ev.target.value
  }
}

export const UPDATE_NEW_TASK_DESCRIPTION = 'UPDATE_NEW_TASK_DESCRIPTION'
export function updateNewTaskDescription(ev) {
  return {
    type: UPDATE_NEW_TASK_DESCRIPTION,
    description: ev.target.value
  }
}

export const REQUEST_CREATE_TASK = 'REQUEST_CREATE_TASK'
function requestCreateTask() {
  return {
    type: REQUEST_CREATE_TASK
  }
}

export const RECEIVE_CREATE_TASK = 'RECEIVE_CREATE_TASK'
function receiveCreateTask(task) {
  return {
    type: RECEIVE_CREATE_TASK,
    task: task
  }
}

export function fetchCreateTask(projectID, listID, title, description) {
  return dispatch => {
    dispatch(requestCreateTask())
    return Request
      .post(`/projects/${projectID}/lists/${listID}/tasks`)
      .type('form')
      .send({title: title, description: description})
      .end((err, res)=> {
        if(res.ok) {
          dispatch(receiveCreateTask(res.body))
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
