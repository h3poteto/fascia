import * as boardActions from '../actions/BoardAction';

const initState = {
  isModalOpen: false,
  newProject: "",
  projects: [],
  repositories: [],
  selectedRepository: null
};

export default function BoardReducer(state = initState, action) {
  switch(action.type) {
  case boardActions.REQUEST_POSTS:
    return state;
  case boardActions.RECEIVE_POSTS:
    return Object.assign({}, state, {
      projects: action.projects
    });
  case boardActions.OPEN_NEW_PROJECT:
    return Object.assign({}, state, {
      isModalOpen: action.isModalOpen
    });
  case boardActions.CLOSE_NEW_PROJECT:
    return Object.assign({}, state, {
      isModalOpen: action.isModalOpen
    });
  case boardActions.REQUEST_REPOSITORIES:
    return state;
  case boardActions.RECEIVE_REPOSITORIES:
    return Object.assign({}, state, {
      repositories: action.repositories
    });
  case boardActions.CHANGE_SELECT_REPOSITORY:
    return Object.assign({}, state, {
      selectedRepository: action.selectEvent.value,
      newProject: action.selectEvent.options[action.selectEvent.selectedIndex].text
    });
  case boardActions.RECEIVE_CREATE_PROJECT:
    const projects = state.projects.concat([action.project]);
    return Object.assign({}, state, {
      newProject: "",
      projects: projects,
      isModalOpen: false
    });
  case boardActions.UPDATE_NEW_PROJECT:
    return Object.assign({}, state, {
      newProject: action.title
    });
  default:
    return state;
  }
}
