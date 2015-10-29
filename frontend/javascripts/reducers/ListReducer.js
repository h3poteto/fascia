import * as listActions from '../actions/ListAction';

const initState = {
  isListModalOpen: false,
  isTaskModalOpen: false,
  newList: {title: ""},
  newTask: {title: ""},
  lists: [],
  selectedListId: null,
  project: null
};

export default function ListReducer(state = initState, action) {
  switch(action.type) {
  case listActions.OPEN_NEW_LIST:
  case listActions.CLOSE_NEW_LIST:
    return Object.assign({}, state, {
      isListModalOpen: action.isListModalOpen
    });
  case listActions.OPEN_NEW_TASK:
    return Object.assign({}, state, {
      isTaskModalOpen: action.isTaskModalOpen,
      selectedListId: action.listId
    });
  case listActions.CLOSE_NEW_TASK:
    return Object.assign({}, state, {
      isTaskModalOpen: action.isTaskModalOpen
    });
  case listActions.UPDATE_NEW_LIST_TITLE:
    var newList = state.newList;
    newList.title = action.title;
    return Object.assign({}, state, {
      newList: newList
    });
  case listActions.UPDATE_NEW_TASK_TITLE:
    var newTask = state.newTask;
    newTask.title = action.title;
    return Object.assign({}, state, {
      newTask: newTask
    });
  case listActions.RECEIVE_LISTS:
    var lists = action.lists.map(function(list, index) {
      if (list.ListTasks == null) {
        list.ListTasks = [];
        return list;
      } else {
        return list;
      }
    });
    return Object.assign({}, state, {
      lists: lists
    });
  case listActions.RECEIVE_PROJECT:
    return Object.assign({}, state, {
      project: action.project
    });
  case listActions.RECEIVE_CREATE_LIST:
    var createdList = action.list;
    if (createdList.ListTasks == null) {
      createdList.ListTasks = [];
    }
    var lists = state.lists.concat([createdList]);
    return Object.assign({}, state, {
      newList: {title: ""},
      lists: lists,
      isListModalOpen: false
    });
  case listActions.RECEIVE_CREATE_TASK:
    var lists = state.lists.map(function(l, index) {
      if (l.Id == action.task.ListId) {
        l.ListTasks = l.ListTasks.concat([action.task]);
        return l;
      } else {
        return l;
      }
    });
    return Object.assign({}, state, {
      newTask: {title: ""},
      lists: lists,
      isTaskModalOpen: false
    });
  default:
    return state;
  }
}
