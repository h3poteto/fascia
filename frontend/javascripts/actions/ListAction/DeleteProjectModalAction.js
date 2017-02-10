import Request from 'superagent'
import { browserHistory } from 'react-router'

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

export const CLOSE_DELETE_PROJECT = 'CLOSE_DELETE_PROJECT'
export function closeDeleteProjectModal() {
  return {
    type: CLOSE_DELETE_PROJECT,
  }
}

export const REQUEST_DELETE_PROJECT = 'REQUEST_DELETE_PROJECT'
function requestDeleteProject() {
  return {
    type: REQUEST_DELETE_PROJECT,
  }
}

export const RECEIVE_DELETE_PROJECT = 'RECEIVE_DELETE_PROJECT'
function receiveDeleteProject(body) {
  return {
    type: RECEIVE_DELETE_PROJECT,
  }
}

export function fetchDeleteProject(projectID) {
  return dispatch => {
    dispatch(requestDeleteProject())
    return Request
      .del(`/projects/${projectID}`)
      .end((err, res)=> {
        if (res.ok) {
          dispatch(receiveDeleteProject(res.body))
          browserHistory.push('/')
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
