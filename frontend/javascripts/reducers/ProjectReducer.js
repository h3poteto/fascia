import * as projectActions from '../actions/ProjectAction'
import * as newProjectModalActions from '../actions/ProjectAction/NewProjectModalAction'

const initState = {
  isModalOpen: false,
  projects: [],
  repositories: [],
  isLoading: false,
  error: null
}

export default function ProjectReducer(state = initState, action) {
  switch(action.type) {
    //-----------------------------------
    // newProjectModalActions
    //-----------------------------------
  case newProjectModalActions.NOT_FOUND:
    return Object.assign({}, state, {
      isLoading: false,
      error: 'Error Not Found'
    })
  case newProjectModalActions.SERVER_ERROR:
    return Object.assign({}, state, {
      isLoading: false,
      error: 'Internal Server Error'
    })
  case newProjectModalActions.CLOSE_NEW_PROJECT:
    return Object.assign({}, state, {
      isModalOpen: false
    })
  case newProjectModalActions.REQUEST_CREATE_PROJECT:
    return Object.assign({}, state, {
      isLoading: true
    })
  case newProjectModalActions.RECEIVE_CREATE_PROJECT: {
    const projects = state.projects.concat([action.project])
    return Object.assign({}, state, {
      projects: projects,
      isLoading: false,
      isModalOpen: false
    })
  }

    //-----------------------------------
    // projectActions
    //-----------------------------------
  case projectActions.NOT_FOUND:
    return Object.assign({}, state, {
      isLoading: false,
      error: 'Error Not Found'
    })
  case projectActions.SERVER_ERROR:
    return Object.assign({}, state, {
      isLoading: false,
      error: 'Internal Server Error'
    })
  case projectActions.CLOSE_FLASH:
    return Object.assign({}, state, {
      error: null
    })
  case projectActions.RECEIVE_POSTS:
    var prj
    if (action.projects == null) {
      prj = []
    } else {
      prj = action.projects
    }
    return Object.assign({}, state, {
      projects: prj
    })
  case projectActions.OPEN_NEW_PROJECT:
    return Object.assign({}, state, {
      isModalOpen: action.isModalOpen
    })
  case projectActions.REQUEST_REPOSITORIES:
    return state
  case projectActions.RECEIVE_REPOSITORIES:
    var repo = action.repositories
    if (repo == null) {
      repo = []
    }
    return Object.assign({}, state, {
      repositories: repo,
    })
  default:
    return state
  }
}
