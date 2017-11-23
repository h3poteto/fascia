import axios from 'axios'
import { ErrorHandler, ServerError } from '../ErrorHandler'
import { startLoading, stopLoading } from '../Loading'

export const CLOSE_NEW_PROJECT = 'CLOSE_NEW_PROJECT'
export function closeNewProjectModal() {
  return {
    type: CLOSE_NEW_PROJECT
  }
}

export const REQUEST_CREATE_PROJECT = 'REQUEST_CREATE_PROJECT'
function requestCreateProject() {
  return {
    type: REQUEST_CREATE_PROJECT
  }
}

export const RECEIVE_CREATE_PROJECT = 'RECEIVE_CREATE_PROJECT'
function receiveCreateProject(body) {
  return {
    type: RECEIVE_CREATE_PROJECT,
    project: {ID: body.ID, UserID: body.UserID, Title: body.Title, Description: body.Description}
  }
}


export function fetchCreateProject(params) {
  return dispatch => {
    dispatch(startLoading())
    dispatch(requestCreateProject())
    return axios
      .post('/projects', params)
      .then((res) => {
        dispatch(stopLoading())
        dispatch(receiveCreateProject(res.data))
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
