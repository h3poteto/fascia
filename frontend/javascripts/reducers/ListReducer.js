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
  taskDragTarget: null,
  taskDragFromList: null,
  isTaskDragging: false
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
    var taskDragTarget = null;
    var lists = state.lists;
    state.lists.map(function(list, i) {
      if (list.Id == action.taskDragFromList.dataset.id) {
        fromList = list;
        list.ListTasks.map(function(task, j) {
          if (task.Id == action.taskDragTarget.dataset.id) {
            taskDragTarget = task;
          }
        });
      }
    });
    return Object.assign({}, state, {
      taskDragTarget: taskDragTarget,
      isTaskDragging: true,
      lists: lists,
      taskDragFromList: fromList
    });
  case listActions.TASK_DRAG_END:
    // var lists = state.lists;
    // if (state.taskDragToList == null) {
    //   // これはどこのリストにも入れなかったやつ
    //   // なにもしない
    // } else if (state.taskDragToList.Id == state.taskDragFromList.Id) {
    //   // これは並び替え
    // } else {
    //   // targetListに追加する
    //   // TODO: 並び順も考慮して挿入できるようにしておく
    //   lists = state.lists.map(function(list, i) {
    //     if (list.Id == state.taskDragToList.Id) {
    //       list.ListTasks.push(state.taskDragTarget);
    //       return list;
    //     } else if (list.Id == state.taskDragFromList.Id) {
    //       var taskIndex;
    //       list.ListTasks.map(function(task, j) {
    //         if (task.Id == state.taskDragTarget.Id) {
    //           taskIndex = j;
    //         };
    //       });
    //       list.ListTasks.splice(taskIndex, 1);
    //       return list;
    //     } else {
    //       return list;
    //     }
    //   });
    // }
    // return Object.assign({}, state, {
    //   taskDragTarget: null,
    //   taskDragFromList: null,
    //   taskDragToList: null,
    //   isTaskDragging: false,
    //   lists: lists
    // });
    return state;
  case listActions.TASK_DROP:
    // TODO: あとでplaceholderか何かを作ってドラッグ中の仮想位置を表示する
    var toList = null;
    state.lists.map(function(list, i) {
      if (list.Id == action.taskDragToList.dataset.id) {
        toList = list;
      }
    });
    var lists = state.lists;
    if (toList == null) {
      // これはどこのリストにも入れなかったやつ
      // なにもしない
    } else if (toList.Id == state.taskDragFromList.Id) {
      // これは並び替え
    } else {
      // targetListに追加する
      // TODO: 並び順も考慮して挿入できるようにしておく
      lists = state.lists.map(function(list, i) {
        if (list.Id == toList.Id) {
          list.ListTasks.push(state.taskDragTarget);
          return list;
        } else if (list.Id == state.taskDragFromList.Id) {
          var taskIndex;
          list.ListTasks.map(function(task, j) {
            if (task.Id == state.taskDragTarget.Id) {
              taskIndex = j;
            };
          });
          list.ListTasks.splice(taskIndex, 1);
          return list;
        } else {
          return list;
        }
      });
    }
    return Object.assign({}, state, {
      taskDragTarget: null,
      taskDragFromList: null,
      isTaskDragging: false,
      lists: lists
    });
  default:
    return state;
  }
}
