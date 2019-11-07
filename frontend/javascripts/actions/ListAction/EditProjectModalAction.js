import axios from 'axios'
import { ErrorHandler, ServerError } from '../ErrorHandler'
import { startLoading, stopLoading } from '../Loading'

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
export function createWebhook() {
  return (dispatch, getState) => {
    const {
      ListReducer: {
        project: { ID: projectID }
      }
    } = getState()
    dispatch(startLoading())
    dispatch(requestCreateWebhook())
    return axios
      .post(`/api/projects/${projectID}/webhook`)
      .then(_res => {
        dispatch(stopLoading())
        dispatch(receiveCreateWebhook())
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

export const CLOSE_EDIT_PROJECT = 'CLOSE_EDIT_PROJECT'
export function closeEditProjectModal() {
  return {
    type: CLOSE_EDIT_PROJECT
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
export function fetchUpdateProject(params) {
  return (dispatch, getState) => {
    const {
      ListReducer: {
        project: { ID: projectID }
      }
    } = getState()
    dispatch(startLoading())
    dispatch(requestUpdateProject())
    return axios
      .patch(`/api/projects/${projectID}`, params)
      .then(res => {
        dispatch(stopLoading())
        dispatch(receiveUpdateProject(res.data))
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
