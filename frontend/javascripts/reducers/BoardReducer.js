import * as boardActions from '../actions/BoardAction';

const initState = {
  isModalOpen: false,
  newProject: "",
  projects: [],
  repositories: [],
  selectedRepository: []
};

export default function posts(state = initState, action) {
  switch(action.type) {
  case boardActions.REQUEST_POSTS:
    return state;
  case boardActions.RECEIVE_POSTS:
    return Object.assign({}, state, {
      projects: action.projects
    });
  case boardActions.NEW_PROJECT:
    return Object.assign({}, state, {
      isModalOpen: action.isModalOpen
    });
  default:
    return state;
  }
}
