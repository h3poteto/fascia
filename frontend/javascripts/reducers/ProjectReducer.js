import * as projectActions from '../actions/ProjectAction';

const initState = {
  isModalOpen: false,
  newProject: {title: "", description: ""},
  projects: [],
  repositories: [],
  selectedRepository: null,
  isLoading: false,
  error: null
};

export default function ProjectReducer(state = initState, action) {
  switch(action.type) {
  case projectActions.SERVER_ERROR:
    return Object.assign({}, state, {
      isLoading: false,
      error: "Server Error"
    });
  case projectActions.CLOSE_FLASH:
    return Object.assign({}, state, {
      error: null
    });
  case projectActions.RECEIVE_POSTS:
    var prj;
    if (action.projects == null) {
      prj = [];
    } else {
      prj = action.projects;
    }
    return Object.assign({}, state, {
      projects: prj
    });
  case projectActions.OPEN_NEW_PROJECT:
    return Object.assign({}, state, {
      isModalOpen: action.isModalOpen
    });
  case projectActions.CLOSE_NEW_PROJECT:
    return Object.assign({}, state, {
      isModalOpen: action.isModalOpen
    });
  case projectActions.REQUEST_REPOSITORIES:
    return state;
  case projectActions.RECEIVE_REPOSITORIES:
    return Object.assign({}, state, {
      repositories: action.repositories
    });
  case projectActions.CHANGE_SELECT_REPOSITORY:
    var newProject = state.newProject;
    newProject.title = action.selectEvent.options[action.selectEvent.selectedIndex].text;
    // repositoryはオブジェクトを渡したい
    var repository;
    state.repositories.map(function(repo, index) {
      if (repo.id == action.selectEvent.value) {
        repository = repo;
      }
    });
    return Object.assign({}, state, {
      selectedRepository: repository,
      newProject: newProject
    });
  case projectActions.REQUEST_CREATE_PROJECT:
    return Object.assign({}, state, {
      isModalOpen: false,
      isLoading: true
    })
  case projectActions.RECEIVE_CREATE_PROJECT:
    const projects = state.projects.concat([action.project]);
    return Object.assign({}, state, {
      newProject: {title: "", description: ""},
      projects: projects,
      isLoading: false
    });
  case projectActions.UPDATE_NEW_PROJECT_TITLE:
    var newProject = state.newProject;
    newProject.title = action.title;
    return Object.assign({}, state, {
      newProject: newProject
    });
  case projectActions.UPDATE_NEW_PROJECT_DESCRIPTION:
    var newProject = state.newProject;
    newProject.description = action.description;
    return Object.assign({}, state, {
      newProject: newProject
    });
  default:
    return state;
  }
}
