import * as projectActions from '../actions/ProjectAction';

const initState = {
  isModalOpen: false,
  newProject: {title: "", description: ""},
  projects: [],
  repositories: [],
  selectedRepository: null
};

export default function ProjectReducer(state = initState, action) {
  switch(action.type) {
  case projectActions.REQUEST_POSTS:
    return state;
  case projectActions.RECEIVE_POSTS:
    return Object.assign({}, state, {
      projects: action.projects
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
    return Object.assign({}, state, {
      selectedRepository: action.selectEvent.value,
      newProject: action.selectEvent.options[action.selectEvent.selectedIndex].text
    });
  case projectActions.RECEIVE_CREATE_PROJECT:
    const projects = state.projects.concat([action.project]);
    return Object.assign({}, state, {
      newProject: {title: "", description: ""},
      projects: projects,
      isModalOpen: false
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
