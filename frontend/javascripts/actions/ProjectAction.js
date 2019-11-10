import axios from 'axios'
import { ErrorHandler, ServerError } from './ErrorHandler'

export const CLOSE_FLASH = 'CLOSE_FLASH'
export function closeFlash() {
  return {
    type: CLOSE_FLASH
  }
}

export const REQUEST_SESSION = 'REQUEST_SESSION'
function requestSession() {
  return {
    type: REQUEST_SESSION
  }
}

export const RECEIVE_SESSION = 'RECEIVE_SESSION'
function receiveSession() {
  return {
    type: RECEIVE_SESSION
  }
}

export function fetchSession() {
  return dispatch => {
    dispatch(requestSession())
    return axios
      .patch('/session')
      .then(_ => {
        dispatch(receiveSession())
      })
      .catch(err => {
        ErrorHandler(err)
          .then()
          .catch(error => {
            dispatch(ServerError(error))
          })
      })
  }
}

export const REQUEST_POSTS = 'REQUEST_POSTS'
function requestPosts() {
  return {
    type: REQUEST_POSTS
  }
}

export const RECEIVE_POSTS = 'RECEIVE_POSTS'
function receivePosts(projects) {
  return {
    type: RECEIVE_POSTS,
    projects: projects
  }
}

export function fetchProjects() {
  return dispatch => {
    dispatch(requestPosts())
    return axios
      .get('/api/projects')
      .then(res => {
        dispatch(receivePosts(res.data))
      })
      .catch(err => {
        ErrorHandler(err)
          .then()
          .catch(error => {
            dispatch(ServerError(error))
          })
      })
  }
}

export const OPEN_NEW_PROJECT = 'OPEN_NEW_PROJECT'
export function openNewProjectModal() {
  return {
    type: OPEN_NEW_PROJECT,
    isModalOpen: true
  }
}

export const REQUEST_REPOSITORIES = 'REQUEST_REPOSITORIES'
function requestRepositories() {
  return {
    type: REQUEST_REPOSITORIES
  }
}

export const RECEIVE_REPOSITORIES = 'RECEIVE_REPOSITORIES'
function receiveRepositories(repositories) {
  return {
    type: RECEIVE_REPOSITORIES,
    repositories: repositories
  }
}

export function fetchRepositories() {
  return dispatch => {
    dispatch(requestRepositories())
    return axios
      .get('/api/github/repositories')
      .then(res => {
        dispatch(receiveRepositories(res.data))
      })
      .catch(err => {
        ErrorHandler(err)
          .then()
          .catch(error => {
            dispatch(ServerError(error))
          })
      })
  }
}
