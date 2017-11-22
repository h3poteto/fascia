import * as loadingActions from '../actions/Loading.js'
import * as errorHandler from '../actions/ErrorHandler.js'
import * as listActions from '../actions/ListAction.js'
import * as newListModalActions from '../actions/ListAction/NewListModalAction.js'
import * as editListModalActions from '../actions/ListAction/EditListModalAction.js'
import * as newTaskModalActions from '../actions/ListAction/NewTaskModalAction.js'
import * as editProjectModalActions from '../actions/ListAction/EditProjectModalAction.js'
import * as showTaskModalActions from '../actions/ListAction/ShowTaskModalAction.js'
import * as deleteProjectModalActions from '../actions/ListAction/DeleteProjectModalAction.js'

const initState = {
  isListModalOpen: false,
  isTaskModalOpen: false,
  isListEditModalOpen: false,
  isProjectEditModalOpen: false,
  isTaskShowModalOpen: false,
  isEditTaskModalVisible: false,
  isDeleteProjectModalOpen: false,
  lists: [],
  listOptions: [],
  noneList: {ID: 0, ListTasks: []},
  selectedList: {},
  project: {Title: '', Description: '', RepositoryID: 0, ShowIssues: true, ShowPullRequests: true},
  selectedTask: {Title: '', IssueNumber: 0, Description: 'description'},
  isTaskDraggingOver: false,
  taskDraggingFrom: null,
  taskDraggingTo: null,
  color: '008ed4',
  isLoading: false,
  error: null
}

export default function ListReducer(state = initState, action) {
  switch(action.type) {
      //-----------------------------------
      // LoadingActions
      //-----------------------------------
    case errorHandler.SERVER_ERROR:
      return Object.assign({}, state, {
        error: action.message
      })
      //-----------------------------------
      // LoadingActions
      //-----------------------------------
    case loadingActions.START_LOADING:
      return Object.assign({}, state, {
        isLoading: true,
      })
    case loadingActions.STOP_LOADING:
      return Object.assign({}, state, {
        isLoading: false,
      })
      //-----------------------------------
      // newListModalActions
      //-----------------------------------
    case newListModalActions.CLOSE_NEW_LIST:
      return Object.assign({}, state, {
        isListModalOpen: false,
      })
    case newListModalActions.RECEIVE_CREATE_LIST:
      var createdList = action.list
      if (createdList.ListTasks == null) {
        createdList.ListTasks = []
      }
      var lists = state.lists.concat([createdList])
      return Object.assign({}, state, {
        lists: lists,
        isListModalOpen: false,
      })
    case newListModalActions.CHANGE_COLOR:
      return Object.assign({}, state, {
        color: action.color,
      })


      //------------------------------------
      // editListModalActions
      //------------------------------------
    case editListModalActions.CLOSE_EDIT_LIST:
      return Object.assign({}, state, {
        isListEditModalOpen: action.isListEditModalOpen,
        selectedList: {},
      })
    case editListModalActions.RECEIVE_UPDATE_LIST: {
      let updatedList = action.list
      if (updatedList.ListTasks == null) {
        updatedList.ListTasks = []
      }
      let lists = state.lists.map(function(list, _) {
        if (list.ID == updatedList.ID) {
          return updatedList
        } else {
          return list
        }
      })
      return Object.assign({}, state, {
        lists: lists,
        isListEditModalOpen: false,
      })
    }
    case editListModalActions.CHANGE_COLOR:
      return Object.assign({}, state, {
        color: action.color,
      })
      //------------------------------------
      // newTaskModalActions
      //------------------------------------
    case newTaskModalActions.CLOSE_NEW_TASK:
      return Object.assign({}, state, {
        isTaskModalOpen: action.isTaskModalOpen,
        selectedList: {},
      })
    case newTaskModalActions.RECEIVE_CREATE_TASK:
      let lists
      if (action.lists == null) {
        lists = []
      } else {
        lists = action.lists.map(function(list, _) {
          if (list.ListTasks == null) {
            list.ListTasks = []
            return list
          } else {
            return list
          }
        })
      }
      let noneList = state.noneList
      if (action.noneList != null) {
        noneList = action.noneList
      }
      return Object.assign({}, state, {
        lists: lists,
        isTaskModalOpen: false,
        noneList: noneList
      })

      //------------------------------------
      // editProjectModalActions
      //------------------------------------
    case editProjectModalActions.NOT_FOUND:
      return Object.assign({}, state, {
        error: 'Error Not Found',
        isLoading: false
      })
    case editProjectModalActions.SERVER_ERROR:
      return Object.assign({}, state, {
        error: 'Internal Server Error',
        isLoading: false
      })
    case editProjectModalActions.REQUEST_CREATE_WEBHOOK:
      return Object.assign({}, state, {
        isProjectEditModalOpen: false
      })
    case editProjectModalActions.CLOSE_EDIT_PROJECT:
      return Object.assign({}, state, {
        isProjectEditModalOpen: false
      })
    case editProjectModalActions.RECEIVE_UPDATE_PROJECT:
      return Object.assign({}, state, {
        project: action.project,
        isProjectEditModalOpen: false,
        isLoading: false
      })

      //------------------------------------
      // deleteProjectModalActions
      //------------------------------------
    case deleteProjectModalActions.NOT_FOUND:
      return Object.assign({}, state, {
        error: 'Error Not Found',
        isLoading: false
      })
    case deleteProjectModalActions.SERVER_ERROR:
      return Object.assign({}, state, {
        error: 'Internal Server Error',
        isLoading: false
      })
    case deleteProjectModalActions.CLOSE_DELETE_PROJECT:
      return Object.assign({}, state, {
        isDeleteProjectModalOpen: false
      })
    case deleteProjectModalActions.REQUEST_DELETE_PROJECT:
      return Object.assign({}, state, {
        isLoading: true
      })
    case deleteProjectModalActions.RECEIVE_DELETE_PROJECT:
      return Object.assign({}, state, {
        isLoading: false,
        isDeleteProjectModalOpen: false
      })


      //------------------------------------
      // showTaskModalActions
      //------------------------------------
    case showTaskModalActions.CLOSE_SHOW_TASK:
      return Object.assign({}, state, {
        isTaskShowModalOpen: false,
        isEditTaskModalVisible: false
      })
    case showTaskModalActions.CHANGE_EDIT_MODE:
      return Object.assign({}, state, {
        isEditTaskModalVisible: !state.isEditTaskModalVisible,
      })
    case showTaskModalActions.RECEIVE_UPDATE_TASK: {
      let lists
      if (action.lists == null) {
        lists = []
      } else {
        lists = action.lists.map(function(list, _) {
          if (list.ListTasks == null) {
            list.ListTasks = []
            return list
          } else {
            return list
          }
        })
      }
      let noneList = state.noneList
      if (action.noneList != null) {
        noneList = action.noneList
      }
      return Object.assign({}, state, {
        lists: lists,
        noneList: noneList,
        isEditTaskModalVisible: false,
        isTaskShowModalOpen: false
      })
    }
    case showTaskModalActions.RECEIVE_DELETE_TASK: {
      let lists
      if (action.lists == null) {
        lists = []
      } else {
        lists = action.lists.map(function(list, _) {
          if (list.ListTasks == null) {
            list.ListTasks = []
            return list
          } else {
            return list
          }
        })
      }
      let noneList = state.noneList
      if (action.noneList != null) {
        noneList = action.noneList
      }
      return Object.assign({}, state, {
        lists: lists,
        noneList: noneList,
        isEditTaskModalVisible: false,
        isTaskShowModalOpen: false
      })
    }

      //------------------------------------
      // listActions
      //------------------------------------
    case listActions.NOT_FOUND:
      return Object.assign({}, state, {
        error: 'Error Not Found',
        isLoading: false
      })
    case listActions.SERVER_ERROR:
      return Object.assign({}, state, {
        error: 'Internal Server Error',
        isLoading: false
      })
    case listActions.CLOSE_FLASH:
      return Object.assign({}, state, {
        error: null
      })
    case listActions.REQUEST_FETCH_GITHUB:
      return Object.assign({}, state, {
        isLoading: true
      })
    case listActions.OPEN_NEW_LIST:
      return Object.assign({}, state, {
        isListModalOpen: action.isListModalOpen
      })
    case listActions.OPEN_NEW_TASK:
      return Object.assign({}, state, {
        isTaskModalOpen: action.isTaskModalOpen,
        selectedList: Object.assign({}, action.list)
      })
    case listActions.OPEN_EDIT_LIST:
      return Object.assign({}, state, {
        isListEditModalOpen: action.isListEditModalOpen,
        selectedList: Object.assign({}, action.list),
      })
    case listActions.OPEN_SHOW_TASK:
      return Object.assign({}, state, {
        isTaskShowModalOpen: true,
        selectedTask: Object.assign({}, action.task)
      })
    case listActions.RECEIVE_LISTS:
    case listActions.RECEIVE_FETCH_GITHUB:
    case listActions.RECEIVE_MOVE_TASK:
    case listActions.RECEIVE_HIDE_LIST:
    case listActions.RECEIVE_DISPLAY_LIST: {
      let lists
      if (action.lists == null) {
        lists = []
      } else {
        lists = action.lists.map(function(list, _) {
          if (list.ListTasks == null) {
            list.ListTasks = []
            return list
          } else {
            return list
          }
        })
      }
      let noneList = state.noneList
      if (action.noneList != null) {
        noneList = action.noneList
      }
      return Object.assign({}, state, {
        lists: lists,
        noneList: noneList,
        isLoading: false
      })
    }
    case listActions.RECEIVE_PROJECT:
      return Object.assign({}, state, {
        project: action.project
      })
    case listActions.TASK_DRAG_START: {
      let taskDraggingFrom
      state.lists.map(function(list, _) {
        if (list.ID == action.taskDragFromList.dataset.id) {
          list.ListTasks.map(function(task, _) {
            if (task.ID == action.taskDragTarget.dataset.id) {
              taskDraggingFrom = {fromList: list, fromTask: task}
            }
          })
        }
      })

      state.noneList.ListTasks.map(function(task, _) {
        if (task.ID == action.taskDragTarget.dataset.id) {
          taskDraggingFrom = {fromList: state.noneList, fromTask: task}
        }
      })

      return Object.assign({}, state, {
        taskDraggingFrom: taskDraggingFrom
      })
    }
    case listActions.TASK_DRAG_LEAVE: {
      // arrowを抜いて
      let lists = state.lists.map(function(list, _) {
        let taskIndex = null
        list.ListTasks.map(function(task, j) {
          if (task.draggedOn) {
            taskIndex = j
          }
        })
        if (taskIndex != null) {
          list.ListTasks.splice(taskIndex, 1)
          list.isDraggingOver = false
        }
        return list
      })
      let noneList = state.noneList
      let taskIndex = null
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
      })
    }
    case listActions.TASK_DROP: {
      let lists = state.lists.map(function(list, _) {
        // arrowを抜く
        let taskIndex = null
        list.ListTasks.map(function(task, j) {
          if (task.draggedOn) {
            taskIndex = j
          }
        })
        if (taskIndex != null) {
          list.ListTasks.splice(taskIndex, 1)
          list.isDraggingOver = false
        }
        return list
      })
      let noneList = state.noneList
      let taskIndex = null
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
      })
    }
    case listActions.REQUEST_MOVE_TASK: {
      let lists = state.lists.map(function(list, _) {
        // arrowを抜く
        let taskIndex = null
        list.ListTasks.map(function(task, j) {
          if (task.draggedOn) {
            taskIndex = j
          }
        })
        if (taskIndex != null) {
          list.ListTasks.splice(taskIndex, 1)
          list.isDraggingOver = false
        }
        // loadingを表示する
        if (list.ID == state.taskDraggingFrom.fromList.ID || list.ID == state.taskDraggingTo.toList.ID) {
          list.isLoading = true
        }
        return list
      })
      let noneList = state.noneList
      let taskIndex = null
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
      })
    }
    case listActions.TASK_DRAG_OVER: {
      // arrowの操作のみ
      let toList = null
      let lists = state.lists
      let noneList = state.noneList
      let taskDraggingTo = state.taskDraggingTo
      if (!state.isTaskDraggingOver) {
        state.lists.map(function(list, _) {
          if (list.ID == action.taskDragToList.dataset.id) {
            toList = list
          }
        })
        if (state.noneList.ID == action.taskDragToList.dataset.id) {
          toList = state.noneList
        }
        if (toList == null) {
          // こんな場合はありえないが
        } else if(action.taskDragToTask.className == 'task') {
          // taskの直前に入れる
          lists = state.lists.map(function(list, _) {
            if (list.ID == toList.ID) {
              let taskIndex
              list.ListTasks.map(function(task, j) {
                if (task.ID == action.taskDragToTask.dataset.id) {
                  taskIndex = j
                  taskDraggingTo = {toList: list, prevToTask: task}
                }
              })
              list.ListTasks.splice(taskIndex, 0, {draggedOn: true})
              list.isDraggingOver = true
              return list
            } else {
              return list
            }
          })
          let taskIndex
          if (noneList.ID == toList.ID) {
            state.noneList.ListTasks.map(function(task, j) {
              if (task.ID == action.taskDragToTask.dataset.id) {
                taskIndex = j
                taskDraggingTo = {toList: noneList, prevToTask: task}
              }
              return task
            })
            noneList.ListTasks.splice(taskIndex, 0, {draggedOn: true})
          }
        } else {
          // taskの末尾に入れる
          lists = state.lists.map(function(list, _) {
            if (list.ID == toList.ID) {
              list.ListTasks.push({draggedOn: true})
              list.isDraggingOver = true
              taskDraggingTo = {toList: list, prevToTask: null}
              return list
            } else {
              return list
            }
          })
          if (noneList.ID == toList.ID) {
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
      })
    }
    case listActions.RECEIVE_LIST_OPTIONS:
      return Object.assign({}, state, {
        listOptions: action.listOptions
      })
    case listActions.OPEN_EDIT_PROJECT:
      return Object.assign({}, state, {
        isProjectEditModalOpen: true
      })
    case listActions.OPEN_DELETE_PROJECT:
      return Object.assign({}, state, {
        isDeleteProjectModalOpen: true
      })

    default:
      return state
  }
}
