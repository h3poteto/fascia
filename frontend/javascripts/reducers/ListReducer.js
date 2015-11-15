import * as listActions from '../actions/ListAction';

const initState = {
  isListModalOpen: false,
  isTaskModalOpen: false,
  isListEditModalOpen: false,
  newList: {title: "", color: "0effff"},
  newTask: {title: ""},
  lists: [],
  selectedList: null,
  project: null,
  isTaskDraggingOver: false,
  taskDraggingFrom: null,
  taskDraggingTo: null
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
      selectedList: action.list
    });
  case listActions.CLOSE_NEW_TASK:
    return Object.assign({}, state, {
      isTaskModalOpen: action.isTaskModalOpen,
      selectedList: null
    });
  case listActions.OPEN_EDIT_LIST:
    return Object.assign({}, state, {
      isListEditModalOpen: action.isListEditModalOpen,
      selectedList: action.list
    });
  case listActions.CLOSE_EDIT_LIST:
    return Object.assign({}, state, {
      isListEditModalOpen: action.isListEditModalOpen,
      selectedList: null
    });
  case listActions.UPDATE_NEW_LIST_TITLE:
    var newList = state.newList;
    newList.title = action.title;
    return Object.assign({}, state, {
      newList: newList
    });
  case listActions.UPDATE_NEW_LIST_COLOR:
    var newList = state.newList;
    newList.color = action.color;
    return Object.assign({}, state, {
      newList: newList
    });
  case listActions.UPDATE_SELECTED_LIST_TITLE:
    var list = state.selectedList;
    list.Title.String = action.title;
    return Object.assign({}, state, {
      selectedList: list
    });
  case listActions.UPDATE_SELECTED_LIST_COLOR:
    var list = state.selectedList;
    list.Color.String = action.color;
    return Object.assign({}, state, {
      selectedList: list
    });
  case listActions.UPDATE_NEW_TASK_TITLE:
    var newTask = state.newTask;
    newTask.title = action.title;
    return Object.assign({}, state, {
      newTask: newTask
    });
  case listActions.RECEIVE_LISTS:
  case listActions.RECEIVE_MOVE_TASK:
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
      newList: {title: "", color: "0effff"},
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
      newTask: {title: "", color: "0effff"},
      lists: lists,
      isTaskModalOpen: false
    });
  case listActions.RECEIVE_UPDATE_LIST:
    var updatedList = action.list;
    if (updatedList.ListTasks == null) {
      updatedList.ListTasks = [];
    }
    var lists = state.lists.map(function(list, index) {
      if (list.Id == updatedList.Id) {
        return updatedList;
      } else {
        return list;
      }
    });
    return Object.assign({}, state, {
      lists: lists,
      isListEditModalOpen: false
    });
  case listActions.TASK_DRAG_START:
    var fromList = null;
    var lists = state.lists;
    var taskDraggingFrom;
    state.lists.map(function(list, i) {
      if (list.Id == action.taskDragFromList.dataset.id) {
        list.ListTasks.map(function(task, j) {
          if (task.Id == action.taskDragTarget.dataset.id) {
            taskDraggingFrom = {fromList: list, fromTask: task};
          }
        });
      }
    });

    return Object.assign({}, state, {
      lists: lists,
      taskDraggingFrom: taskDraggingFrom
    });
  case listActions.TASK_DRAG_LEAVE:
    // arrowを抜いて
    var lists = state.lists.map(function(list, i) {
      var taskIndex = null;
      list.ListTasks.map(function(task, j) {
        if (task.draggedOn) {
          taskIndex = j;
        }
      });
      if (taskIndex != null) {
        list.ListTasks.splice(taskIndex, 1);
      }
      return list;
    });
    return Object.assign({}, state, {
      isTaskDraggingOver: false,
      lists: lists,
      taskDraggingTo: null
    });
  case listActions.TASK_DROP:
  case listActions.REQUEST_MOVE_TASK:
    var lists = state.lists.map(function(list, i) {
      // arrowを抜く
      var taskIndex = null;
      list.ListTasks.map(function(task, j) {
        if (task.draggedOn) {
          taskIndex = j;
        }
        if (taskIndex != null) {
          list.ListTasks.splice(taskIndex, 1);
        }
      });
      return list;
    });
    return Object.assign({}, state, {
      isTaskDraggingOver: false,
      lists: lists,
      taskDraggingFrom: null,
      taskDraggingTo: null
    });
  case listActions.TASK_DRAG_OVER:
    // arrowの操作のみ
    var toList = null;
    var lists = state.lists;
    var taskDraggingTo = state.taskDraggingTo;
    if (!state.isTaskDraggingOver) {
      state.lists.map(function(list, i) {
        if (list.Id == action.taskDragToList.dataset.id) {
          toList = list;
        }
      });
      if (toList == null) {
        // こんな場合はありえないが
      } else if(action.taskDragToTask.className == "task") {
        // taskの直前に入れる
        lists = state.lists.map(function(list, i) {
          if (list.Id == toList.Id) {
            var taskIndex;
            list.ListTasks.map(function(task, j) {
              if (task.Id == action.taskDragToTask.dataset.id) {
                taskIndex = j;
                taskDraggingTo = {toList: list, prevToTask: task};
              }
            });
            list.ListTasks.splice(taskIndex, 0, {draggedOn: true});
            return list;
          } else {
            return list;
          }
        });
      } else {
        // taskの末尾に入れる
        lists = state.lists.map(function(list, i) {
          if (list.Id == toList.Id) {
            list.ListTasks.push({draggedOn: true});
            taskDraggingTo = {toList: list, prevToTaks: null};
            return list;
          } else {
            return list;
          }
        });
      }
    }
    return Object.assign({}, state, {
      isTaskDraggingOver: true,
      lists: lists,
      taskDraggingTo: taskDraggingTo
    });
  default:
    return state;
  }
}
