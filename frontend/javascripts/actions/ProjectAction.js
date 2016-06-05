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

export const CLOSE_FLASH = "CLOSE_FLASH"
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
    return Request
      .post('/session')
      .end((err, res) => {
        if (res.ok) {
          dispatch(receiveSession())
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
    return Request
      .get('/projects')
      .end((err, res)=> {
        if (res.ok) {
          dispatch(receivePosts(res.body))
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

export const OPEN_NEW_PROJECT = 'OPEN_NEW_PROJECT'
export function openNewProjectModal() {
  return {
    type: OPEN_NEW_PROJECT,
    isModalOpen: true
  }
}

export const CLOSE_NEW_PROJECT = 'CLOSE_NEW_PROJECT'
export function closeNewProjectModal() {
  return {
    type: CLOSE_NEW_PROJECT,
    isModalOpen: false
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
    return Request
      .get('/github/repositories')
      .end((err, res)=> {
        if (res.ok) {
          dispatch(receiveRepositories(res.body))
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

export const CHANGE_SELECT_REPOSITORY = 'CHANGE_SELECT_REPOSITORY'
export function changeSelectedRepository(ev) {
  return {
    type: CHANGE_SELECT_REPOSITORY,
    selectEvent: ev.target
  }
}

export const REQUEST_CREATE_PROJECT = 'REQUEST_CREATE_PROJECT'
function requestCreateProject() {
  return {
    type: REQUEST_CREATE_PROJECT
  }
}

export const RECEIVE_CREATE_PROJECT = 'RECEIVE_CREATE_PROJECT'
function receiveCreateProject(id, userID, title, description) {
  return {
    type: RECEIVE_CREATE_PROJECT,
    project: {ID: id, UserID: userID, Title: title, Description: description}
  }
}


export function fetchCreateProject(title, description, repository) {
  return dispatch => {
    dispatch(requestCreateProject())
    var repositoryID, repositoryOwner, repositoryName
    if (repository != null) {
      repositoryID = repository.id
      repositoryOwner = repository.owner.login
      repositoryName = repository.name
    }
    return Request
      .post('/projects')
      .type('form')
      .send({title: title, description: description, repositoryID: repositoryID, repositoryOwner: repositoryOwner, repositoryName: repositoryName})
      .end((err, res)=> {
        if (res.ok) {
          dispatch(receiveCreateProject(res.body.ID, res.body.UserID, res.body.Title, res.body.Description))
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

export const UPDATE_NEW_PROJECT_TITLE = 'UPDATE_NEW_PROJECT_TITLE'
export function updateNewProjectTitle(ev) {
  return {
    type: UPDATE_NEW_PROJECT_TITLE,
    title: ev.target.value
  }
}

export const UPDATE_NEW_PROJECT_DESCRIPTION = 'UPDATE_NEW_PROJECT_DESCRIPTION'
export function updateNewProjectDescription(ev) {
  return {
    type: UPDATE_NEW_PROJECT_DESCRIPTION,
    description: ev.target.value
  }
}
