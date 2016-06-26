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

export const REQUEST_CREATE_WEBHOOK = 'REQUEST_CREATE_WEBHOOK'
export function requestCreateWebhook() {
  return {
    type: REQUEST_CREATE_WEBHOOK
  }
}

export const RECEIVE_CREATE_WEBHOOK = 'RECEIVE_CREATE_WEBHOOK'
function receiveCreateWebhook() {
  return {
    type: RECEIVE_CREATE_WEBHOOK
  }
}

export const CREATE_WEBHOOK = 'CREATE_WEBHOOK'
export function createWebhook(projectID) {
  return dispatch => {
    dispatch(requestCreateWebhook())
    return Request
      .post(`/projects/${projectID}/webhook`)
      .end((err, res) => {
        if (res.ok) {
          dispatch(receiveCreateWebhook())
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

export const CLOSE_EDIT_PROJECT = 'CLOSE_EDIT_PROJECT'
export function closeEditProjectModal() {
  return {
    type: CLOSE_EDIT_PROJECT
  }
}


export const UPDATE_EDIT_PROJECT_TITLE = 'UPDATE_EDIT_PROJECT_TITLE'
export function updateEditProjectTitle(ev) {
  return {
    type: UPDATE_EDIT_PROJECT_TITLE,
    title: ev.target.value
  }
}

export const UPDATE_EDIT_PROJECT_DESCRIPTION = 'UPDATE_EDIT_PROJECT_DESCRIPTION'
export function updateEditProjectDescription(ev) {
  return {
    type: UPDATE_EDIT_PROJECT_DESCRIPTION,
    description: ev.target.value
  }
}

export const REQUEST_UPDATE_PROJECT = 'REQUEST_UPDATE_PROJECT'
function requestUpdateProject() {
  return {
    type: REQUEST_UPDATE_PROJECT
  }
}

export const RECEIVE_UPDATE_PROJECT = 'RECEIVE_UPDATE_PROJECT'
function receiveUpdateProject(project) {
  return {
    type: RECEIVE_UPDATE_PROJECT,
    project: project
  }
}

export const FETCH_UPDATE_PROJECT = 'FETCH_UPDATE_PROJECT'
export function fetchUpdateProject(projectID, project) {
  return dispatch => {
    dispatch(requestUpdateProject())
    return Request
      .post(`/projects/${projectID}`)
      .type('form')
      .send({title: project.Title, description: project.Description})
      .end((err, res) => {
        if (res.ok) {
          dispatch(receiveUpdateProject(res.body))
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
