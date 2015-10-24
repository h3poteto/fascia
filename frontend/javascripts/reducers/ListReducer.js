import * as listActions from '../actions/ListAction';

const initState = {
  isModalOpen: false,
  newList: {title: ""},
  lists: []
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
  default:
    return state;
  }
}
