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

export const CLOSE_DELETE_PROJECT = 'CLOSE_DELETE_PROJECT'
export function closeDeleteProjectModal() {
  return {
    type: CLOSE_DELETE_PROJECT,
  }
}

export function fetchDeleteProject() {
}
