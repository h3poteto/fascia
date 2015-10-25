import * as listActions from '../actions/ListAction';

const initState = {
  isModalOpen: false,
  newList: {title: ""},
  lists: [],
  project: null
};

export default function ListReducer(state = initState, action) {
  switch(action.type) {
  case listActions.OPEN_NEW_LIST:
    return Object.assign({}, state, {
      isModalOpen: action.isModalOpen
    });
  case listActions.CLOSE_NEW_LIST:
    return Object.assign({}, state, {
      isModalOpen: action.isModalOpen
    });
  case listActions.UPDATE_NEW_LIST_TITLE:
    var newList = state.newList;
    newList.title = action.title;
    return Object.assign({}, state, {
      newList: newList
    });
  case listActions.RECEIVE_LISTS:
    return Object.assign({}, state, {
      lists: action.lists
    });
  case listActions.RECEIVE_PROJECT:
    return Object.assign({}, state, {
      project: action.project
    });
  case listActions.RECEIVE_CREATE_LIST:
    const lists = state.lists.concat([action.list]);
    return Object.assign({}, state, {
      newList: {title: ""},
      lists: lists,
      isModalOpen: false
    });
  default:
    return state;
  }
}
