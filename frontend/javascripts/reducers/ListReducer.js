import * as listActions from '../actions/ListAction';

const initState = {
  isListModalOpen: false,
  isTaskModalOpen: false,
  isListEditModalOpen: false,
  isProjectEditModalOpen: false,
  newList: {title: "", color: "008ed4"},
  newTask: {title: "", description: ""},
  lists: [],
  listOptions: [],
  noneList: {Id: 0, ListTasks: []},
  selectedList: null,
  selectedListOption: null,
  project: null,
  selectedProject: {Title: "", Description: "", RepositoryID: 0},
  isTaskDraggingOver: false,
  taskDraggingFrom: null,
  taskDraggingTo: null,
  isLoading: false,
  error: null
};

export default function ListReducer(state = initState, action) {
  switch(action.type) {
  case listActions.SERVER_ERROR:
    return Object.assign({}, state, {
      error: "Server Error",
      isLoading: false
    });
  case listActions.CLOSE_FLASH:
    return Object.assign({}, state, {
      error: null
    });
  case listActions.REQUEST_FETCH_GITHUB:
    return Object.assign({}, state, {
      isLoading: true
    })
  case listActions.OPEN_NEW_LIST:
  case listActions.CLOSE_NEW_LIST:
    return Object.assign({}, state, {
      isListModalOpen: action.isListModalOpen
    });
  case listActions.OPEN_NEW_TASK:
    return Object.assign({}, state, {
      isTaskModalOpen: action.isTaskModalOpen,
      selectedList: Object.assign({}, action.list)
    });
  case listActions.CLOSE_NEW_TASK:
    return Object.assign({}, state, {
      isTaskModalOpen: action.isTaskModalOpen,
      selectedList: null
    });
  case listActions.OPEN_EDIT_LIST:
    var selectedListOption = null
    if (action.list.ListOptionId != 0) {
      selectedListOption = {
        Id: action.list.ListOptionId
      }
    }
    return Object.assign({}, state, {
      isListEditModalOpen: action.isListEditModalOpen,
      selectedList: Object.assign({}, action.list),
      selectedListOption: selectedListOption
    });
  case listActions.CLOSE_EDIT_LIST:
    return Object.assign({}, state, {
      isListEditModalOpen: action.isListEditModalOpen,
      selectedList: null,
      selectedListOption: null
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
    list.Title = action.title;
    return Object.assign({}, state, {
      selectedList: list
    });
  case listActions.UPDATE_SELECTED_LIST_COLOR:
    var list = state.selectedList;
    list.Color = action.color;
    return Object.assign({}, state, {
      selectedList: list
    });
  case listActions.UPDATE_NEW_TASK_TITLE:
    var newTask = state.newTask;
    newTask.title = action.title;
    return Object.assign({}, state, {
      newTask: newTask
    });
  case listActions.UPDATE_NEW_TASK_DESCRIPTION:
    var newTask = state.newTask
    newTask.description = action.description
    return Object.assign({}, state, {
      newTask: newTask
    })
  case listActions.RECEIVE_LISTS:
  case listActions.RECEIVE_FETCH_GITHUB:
  case listActions.RECEIVE_MOVE_TASK:
    var lists
    if (action.lists == null) {
      lists = []
    } else {
      lists = action.lists.map(function(list, index) {
        if (list.ListTasks == null) {
          list.ListTasks = []
          return list
        } else {
          return list
        }
      })
    }
    var noneList = state.noneList
    if (action.noneList != null) {
      noneList = action.noneList
    }
    return Object.assign({}, state, {
      lists: lists,
      noneList: noneList,
      isLoading: false
    });
  case listActions.RECEIVE_PROJECT:
    return Object.assign({}, state, {
      project: action.project,
      selectedProject: Object.assign({}, action.project)
    })
  case listActions.REQUEST_CREATE_LIST:
    return Object.assign({}, state, {
      isLoading: true
    })
  case listActions.RECEIVE_CREATE_LIST:
    var createdList = action.list;
    if (createdList.ListTasks == null) {
      createdList.ListTasks = [];
    }
    var lists = state.lists.concat([createdList]);
    return Object.assign({}, state, {
      newList: {title: "", color: "008ed4"},
      lists: lists,
      isListModalOpen: false,
      isLoading: false
    });
  case listActions.REQUEST_CREATE_TASK:
    return Object.assign({}, state, {
      isLoading: true
    })
  case listActions.RECEIVE_CREATE_TASK:
    var lists = state.lists.map(function(l, index) {
      if (l.Id == action.task.ListId) {
        l.ListTasks = l.ListTasks.concat([action.task]);
        return l;
      } else {
        return l;
      }
    });
    var noneList = state.noneList
    if (action.task.ListId == noneList.Id) {
      noneList.ListTasks = noneList.ListTasks.concat([action.task])
    }
    return Object.assign({}, state, {
      newTask: {title: "", description: ""},
      lists: lists,
      isTaskModalOpen: false,
      isLoading: false,
      noneList: noneList
    });
  case listActions.REQUEST_UPDATE_LIST:
    return Object.assign({}, state, {
      isLoading: true
    })
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
      isListEditModalOpen: false,
      isLoading: false
    });
  case listActions.TASK_DRAG_START:
    var lists = state.lists;
    var taskDraggingFrom;
    state.lists.map(function(list, i) {
      if (list.Id == action.taskDragFromList.dataset.id) {
        list.ListTasks.map(function(task, j) {
          if (task.Id == action.taskDragTarget.dataset.id) {
            taskDraggingFrom = {fromList: list, fromTask: task}
          }
        })
      }
    })

    state.noneList.ListTasks.map(function(task, j) {
      if (task.Id == action.taskDragTarget.dataset.id) {
        taskDraggingFrom = {fromList: state.noneList, fromTask: task}
      }
    })

    return Object.assign({}, state, {
      taskDraggingFrom: taskDraggingFrom
    });
  case listActions.TASK_DRAG_LEAVE:
    // arrowを抜いて
    var lists = state.lists.map(function(list, i) {
      var taskIndex = null
      list.ListTasks.map(function(task, j) {
        if (task.draggedOn) {
          taskIndex = j
        }
      })
      if (taskIndex != null) {
        list.ListTasks.splice(taskIndex, 1)
      }
      return list
    })
    var noneList = state.noneList
    var taskIndex = null
    state.noneList.ListTasks.map(function(task, j) {
      if (task.draggedOn) {
        taskIndex = j
      }
    })
    if (taskIndex != null) {
      noneList.ListTasks.splice(taskIndex, 1)
    }
    return Object.assign({}, state, {
      isTaskDraggingOver: false,
      lists: lists,
      noneList: noneList,
      taskDraggingTo: null
    });
  case listActions.TASK_DROP:
    var lists = state.lists.map(function(list, i) {
      // arrowを抜く
      var taskIndex = null
      list.ListTasks.map(function(task, j) {
        if (task.draggedOn) {
          taskIndex = j
        }
      })
      if (taskIndex != null) {
        list.ListTasks.splice(taskIndex, 1)
      }
      return list
    })
    var noneList = state.noneList
    var taskIndex = null
    state.noneList.ListTasks.map(function(task, j) {
      if (task.draggedOn) {
        taskIndex = j
      }
    })
    if (taskIndex != null) {
      noneList.ListTasks.splice(taskIndex, 1)
    }
    return Object.assign({}, state, {
      isTaskDraggingOver: false,
      lists: lists,
      noneList: noneList,
      taskDraggingFrom: null,
      taskDraggingTo: null
    });
  case listActions.REQUEST_MOVE_TASK:
    var lists = state.lists.map(function(list, i) {
      // arrowを抜く
      var taskIndex = null
      list.ListTasks.map(function(task, j) {
        if (task.draggedOn) {
          taskIndex = j
        }
      })
      if (taskIndex != null) {
        list.ListTasks.splice(taskIndex, 1)
      }
      // loadingを表示する
      if (list.Id == state.taskDraggingFrom.fromList.Id || list.Id == state.taskDraggingTo.toList.Id) {
        list.isLoading = true
      }
      return list
    })
    var noneList = state.noneList
    var taskIndex = null
    state.noneList.ListTasks.map(function(task, j) {
      if (task.draggedOn) {
        taskIndex = j
      }
    })
    if (taskIndex != null) {
      noneList.ListTasks.splice(taskIndex, 1)
    }
    return Object.assign({}, state, {
      isTaskDraggingOver: false,
      lists: lists,
      noneList: noneList,
      taskDraggingFrom: null,
      taskDraggingTo: null
    });
  case listActions.TASK_DRAG_OVER:
    // arrowの操作のみ
    var toList = null
    var lists = state.lists
    var noneList = state.noneList
    var taskDraggingTo = state.taskDraggingTo
    if (!state.isTaskDraggingOver) {
      state.lists.map(function(list, i) {
        if (list.Id == action.taskDragToList.dataset.id) {
          toList = list;
        }
      });
      if (state.noneList.Id == action.taskDragToList.dataset.id) {
        toList = state.noneList
      }
      if (toList == null) {
        // こんな場合はありえないが
      } else if(action.taskDragToTask.className == "task") {
        // taskの直前に入れる
        lists = state.lists.map(function(list, i) {
          if (list.Id == toList.Id) {
            var taskIndex
            list.ListTasks.map(function(task, j) {
              if (task.Id == action.taskDragToTask.dataset.id) {
                taskIndex = j
                taskDraggingTo = {toList: list, prevToTask: task}
              }
            })
            list.ListTasks.splice(taskIndex, 0, {draggedOn: true})
            return list
          } else {
            return list
          }
        })
        var taskIndex
        if (noneList.Id == toList.Id) {
          state.noneList.ListTasks.map(function(task, j) {
            if (task.Id == action.taskDragToTask.dataset.id) {
              taskIndex = j
              taskDraggingTo = {toList: noneList, prevToTask: task}
            }
            return task
          })
          noneList.ListTasks.splice(taskIndex, 0, {draggedOn: true})
        }
      } else {
        // taskの末尾に入れる
        lists = state.lists.map(function(list, i) {
          if (list.Id == toList.Id) {
            list.ListTasks.push({draggedOn: true})
            taskDraggingTo = {toList: list, prevToTask: null}
            return list
          } else {
            return list
          }
        })
        if (noneList.Id == toList.Id) {
          noneList.ListTasks.push({draggedOn: true})
          taskDraggingTo = {toList: noneList, prevToTask: null}
        }
      }
    }
    return Object.assign({}, state, {
      isTaskDraggingOver: true,
      lists: lists,
      noneList: noneList,
      taskDraggingTo: taskDraggingTo
    });
  case listActions.RECEIVE_LIST_OPTIONS:
    return Object.assign({}, state, {
      listOptions: action.listOptions
    })
  case listActions.CHANGE_SELECTED_LIST_OPTION:
    var listOption = {
      Id: action.selectEvent.value
    }
    state.listOptions.map(function(option, index) {
      if (option.Id == action.selectEvent.value) {
        listOption = option
      }
    })
    return Object.assign({}, state, {
      selectedListOption: listOption
    })
  case listActions.OPEN_EDIT_PROJECT:
    return Object.assign({}, state, {
      isProjectEditModalOpen: true
    })
  case listActions.CLOSE_EDIT_PROJECT:
    return Object.assign({}, state, {
      isProjectEditModalOpen: false
    })
  case listActions.UPDATE_EDIT_PROJECT_TITLE:
    var selectedProject = state.selectedProject
    selectedProject.Title = action.title
    return Object.assign({}, state, {
      selectedProject: selectedProject
    })
  case listActions.UPDATE_EDIT_PROJECT_DESCRIPTION:
    var selectedProject = state.selectedProject
    selectedProject.Description = action.description
    return Object.assign({}, state, {
      selectedProject: selectedProject
    })
  case listActions.RECEIVE_UPDATE_PROJECT:
    return Object.assign({}, state, {
      project: action.project,
      isProjectEditModalOpen: false
    })
  default:
    return state
  }
}
