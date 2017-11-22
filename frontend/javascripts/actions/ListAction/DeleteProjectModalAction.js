import axios from 'axios'
import { browserHistory } from 'react-router'
import { ErrorHandler, ServerError } from '../ErrorHandler'
import { startLoading, stopLoading } from '../Loading'

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
function receiveDeleteProject() {
  return {
    type: RECEIVE_DELETE_PROJECT,
  }
}

export function fetchDeleteProject() {
  return (dispatch, getState) => {
    const { ListReducer: { project: { ID: projectID }}} = getState()
    dispatch(startLoading())
    dispatch(requestDeleteProject())
    return axios
      .delete(`/projects/${projectID}`)
      .then((res) => {
        dispatch(stopLoading())
        dispatch(receiveDeleteProject(res.body))
        browserHistory.push('/')
      })
      .catch((err) => {
        dispatch(stopLoading())
        ErrorHandler(err)
          .then()
          .catch((error) => {
            dispatch(ServerError(error))
          })
      })
  }
}
