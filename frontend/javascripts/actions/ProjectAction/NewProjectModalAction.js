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
    dispatch(requestCreateProject())
    console.log(params)
    // TODO: repositoryを上手いこと作れるようにしておく
    // var repositoryID, repositoryOwner, repositoryName
    // if (repository != null) {
    //   repositoryID = repository.id
    //   repositoryOwner = repository.owner.login
    //   repositoryName = repository.name
    // }
    return Request
      .post('/projects')
      .type('form')
      .send(params)
      .end((err, res)=> {
        if (res.ok) {
          dispatch(receiveCreateProject(res.body))
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
