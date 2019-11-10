import axios from 'axios'
import { ErrorHandler, ErrorHandlerWithoutSubmission, ServerError } from './ErrorHandler'
import { startLoading, stopLoading } from './Loading'

export const CLOSE_FLASH = 'CLOSE_FLASH'
export function closeFlash() {
  return {
    type: CLOSE_FLASH
  }
}

export const OPEN_NEW_LIST = 'OPEN_NEW_LIST'
export function openNewListModal() {
  return {
    type: OPEN_NEW_LIST,
    isListModalOpen: true
  }
}

export const OPEN_NEW_TASK = 'OPEN_NEW_TASK'
export function openNewTaskModal(list) {
  return {
    type: OPEN_NEW_TASK,
    isTaskModalOpen: true,
    list: list
  }
}

export const OPEN_EDIT_LIST = 'OPEN_EDIT_LIST'
export function openEditListModal(list) {
  return {
    type: OPEN_EDIT_LIST,
    isListEditModalOpen: true,
    list: list
  }
}

export const REQUEST_LISTS = 'REQUEST_LISTS'
function requestLists() {
  return {
    type: REQUEST_LISTS
  }
}

export const RECEIVE_LISTS = 'RECEIVE_LISTS'
function receiveLists(lists) {
  return {
    type: RECEIVE_LISTS,
    lists: lists.Lists,
    noneList: lists.NoneList
  }
}

export function fetchLists(projectID) {
  return dispatch => {
    dispatch(requestLists())
    return axios
      .get(`/api/projects/${projectID}/lists`)
      .then(res => {
        dispatch(receiveLists(res.data))
      })
      .catch(err => {
        ErrorHandler(err)
          .then()
          .catch(error => {
            dispatch(ServerError(error))
          })
      })
  }
}

export const REQUEST_PROJECT = 'REQUEST_PROJECT'
function requestProject() {
  return {
    type: REQUEST_PROJECT
  }
}

export const RECEIVE_PROJECT = 'RECEIVE_PROJECT'
function receiveProject(project) {
  return {
    type: RECEIVE_PROJECT,
    project: project
  }
}

export function fetchProject(projectID) {
  return dispatch => {
    dispatch(requestProject())
    return axios
      .get(`/api/projects/${projectID}/show`)
      .then(res => {
        dispatch(receiveProject(res.data))
      })
      .catch(err => {
        ErrorHandlerWithoutSubmission(err)
          .then()
          .catch(error => {
            dispatch(ServerError(error))
          })
      })
  }
}

export const TASK_DRAG_START = 'TASK_DRAG_START'
export function taskDragStart(ev) {
  ev.dataTransfer.effectAllowed = 'moved'
  ev.dataTransfer.setData('text/html', ev.currentTarget)
  return {
    type: TASK_DRAG_START,
    taskDragTarget: ev.currentTarget,
    taskDragFromList: ev.currentTarget.parentNode.parentNode
  }
}

export const TASK_DRAG_LEAVE = 'TASK_DRAG_LEAVE'
export const TASK_DRAG_LEAVE_IGNORE = 'TASK_DRAG_LEAVE_IGNORE'
export function taskDragLeave(ev) {
  // li.new-taskだけはdragleaveイベントが頻繁に発生するため，抑制する
  // こうしておいても，他の要素に移動した際には問題なくleave処理が為される
  if (ev.target.className == 'new-task' || ev.target.className == 'fa-plus') {
    return {
      type: TASK_DRAG_LEAVE_IGNORE
    }
  } else {
    return {
      type: TASK_DRAG_LEAVE
    }
  }
}

export const REQUEST_MOVE_TASK = 'REQUEST_MOVE_TASK'
function requestMoveTask() {
  return {
    type: REQUEST_MOVE_TASK
  }
}

export const RECEIVE_MOVE_TASK = 'RECEIVE_MOVE_TASK'
function receiveMoveTask(lists) {
  return {
    type: RECEIVE_MOVE_TASK,
    lists: lists.Lists,
    noneList: lists.NoneList
  }
}

export const TASK_DROP = 'TASK_DROP'
export function taskDrop(projectID, taskDraggingFrom, taskDraggingTo) {
  if (taskDraggingTo != undefined && taskDraggingTo != null) {
    var prevToTaskID
    if (taskDraggingTo.prevToTask == null) {
      prevToTaskID = null
    } else {
      prevToTaskID = taskDraggingTo.prevToTask.ID
    }
    return dispatch => {
      dispatch(requestMoveTask())
      return axios
        .post(`/api/projects/${projectID}/lists/${taskDraggingFrom.fromList.ID}/tasks/${taskDraggingFrom.fromTask.ID}/move_task`, {
          to_list_id: taskDraggingTo.toList.ID,
          prev_to_task_id: prevToTaskID
        })
        .then(res => {
          dispatch(receiveMoveTask(res.data))
        })
        .catch(err => {
          // TODO: ここはドラッグしたviewを元に戻す必要がある
          ErrorHandlerWithoutSubmission(err)
            .then()
            .catch(error => {
              dispatch(ServerError(error))
            })
        })
    }
  } else {
    return {
      type: TASK_DROP
    }
  }
}

export const TASK_DRAG_OVER = 'TASK_DRAG_OVER'
export function taskDragOver(ev) {
  ev.preventDefault()
  var targetList
  switch (ev.target.dataset.droppedDepth) {
    case '0':
      targetList = ev.target
      break
    case '1':
      targetList = ev.target.parentNode
      break
    case '2':
      targetList = ev.target.parentNode.parentNode
      break
    case '3':
      targetList = ev.target.parentNode.parentNode.parentNode
      break
    default:
      targetList = ev.target.parentNode.parentNode
      break
  }
  return {
    type: TASK_DRAG_OVER,
    taskDragToTask: ev.target,
    taskDragToList: targetList
  }
}

export const REQUEST_FETCH_GITHUB = 'REQUEST_FETCH_GITHUB'
function requestFetchGithub() {
  return {
    type: REQUEST_FETCH_GITHUB
  }
}

export const RECEIVE_FETCH_GITHUB = 'RECEIVE_FETCH_GITHUB'
function receiveFetchGithub(lists) {
  return {
    type: RECEIVE_FETCH_GITHUB,
    lists: lists.Lists,
    noneList: lists.NoneList
  }
}

export const FETCH_PROJECT_GITHUB = 'FETCH_PROJECT_GITHUB'
export function fetchProjectGithub(projectID) {
  return dispatch => {
    dispatch(startLoading())
    dispatch(requestFetchGithub())
    return axios
      .post(`/api/projects/${projectID}/fetch_github`)
      .then(res => {
        dispatch(stopLoading())
        dispatch(receiveFetchGithub(res.data))
      })
      .catch(err => {
        dispatch(stopLoading())
        ErrorHandlerWithoutSubmission(err)
          .then()
          .catch(error => {
            dispatch(ServerError(error))
          })
      })
  }
}

export const REQUEST_LIST_OPTIONS = 'REQUEST_LIST_OPTIONS'
function requestListOptions() {
  return {
    type: REQUEST_LIST_OPTIONS
  }
}

export const RECEIVE_LIST_OPTIONS = 'RECEIVE_LIST_OPTIONS'
function receiveListOptions(listOptions) {
  return {
    type: RECEIVE_LIST_OPTIONS,
    listOptions: listOptions
  }
}

export const FETCH_LIST_OPTIONS = 'FETCH_LIST_OPTIONS'
export function fetchListOptions() {
  return dispatch => {
    dispatch(requestListOptions())
    return axios
      .get('/api/list_options')
      .then(res => {
        dispatch(receiveListOptions(res.data))
      })
      .catch(err => {
        ErrorHandlerWithoutSubmission(err)
          .then()
          .catch(error => {
            dispatch(ServerError(error))
          })
      })
  }
}

export const OPEN_EDIT_PROJECT = 'OPEN_EDIT_PROJECT'
export function openEditProjectModal(project) {
  return {
    type: OPEN_EDIT_PROJECT,
    project: project
  }
}

export const REQUEST_SETTINGS_PROJECT = 'REQUEST_SETTINGS_PROJECT'
export function requestSettingsProject() {
  return {
    type: REQUEST_SETTINGS_PROJECT
  }
}

export const SHOW_ISSUES = 'SHOW_ISSUES'
export function showIssues(projectID, showIssues, showPullRequests) {
  return dispatch => {
    dispatch(requestSettingsProject())
    return axios
      .patch(`/api/projects/${projectID}/settings`, {
        show_issues: !showIssues,
        show_pull_requests: showPullRequests
      })
      .then(res => {
        dispatch(receiveProject(res.data))
      })
      .catch(err => {
        ErrorHandlerWithoutSubmission(err)
          .then()
          .catch(error => {
            dispatch(ServerError(error))
          })
      })
  }
}

export const SHOW_PULL_REQUESTS = 'SHOW_PULL_REQUESTS'
export function showPullRequests(projectID, showIssues, showPullRequests) {
  return dispatch => {
    dispatch(requestSettingsProject())
    return axios
      .patch(`/api/projects/${projectID}/settings`, {
        show_issues: showIssues,
        show_pull_requests: !showPullRequests
      })
      .then(res => {
        dispatch(receiveProject(res.data))
      })
      .catch(err => {
        ErrorHandlerWithoutSubmission(err)
          .then()
          .catch(error => {
            dispatch(ServerError(error))
          })
      })
  }
}

export const REQUEST_HIDE_LIST = 'REQUEST_HIDE_LIST'
function requestHideList() {
  return {
    type: REQUEST_HIDE_LIST
  }
}

export const RECEIVE_HIDE_LIST = 'RECEIVE_HIDE_LIST'
function receiveHideList(lists) {
  return {
    type: RECEIVE_HIDE_LIST,
    lists: lists.Lists,
    noneList: lists.NoneList
  }
}

export const HIDE_LIST = 'HIDE_LIST'
export function hideList(projectID, listID) {
  return dispatch => {
    dispatch(requestHideList())
    return axios
      .patch(`/api/projects/${projectID}/lists/${listID}/hide`)
      .then(res => {
        dispatch(receiveHideList(res.data))
      })
      .catch(err => {
        ErrorHandlerWithoutSubmission(err)
          .then()
          .catch(error => {
            dispatch(ServerError(error))
          })
      })
  }
}

export const REQUEST_DISPLAY_LIST = 'REQUEST_DISPLAY_LIST'
function requestDisplayList() {
  return {
    type: REQUEST_DISPLAY_LIST
  }
}

export const RECEIVE_DISPLAY_LIST = 'RECEIVE_DISPLAY_LIST'
function receiveDisplayList(lists) {
  return {
    type: RECEIVE_DISPLAY_LIST,
    lists: lists.Lists,
    noneList: lists.NoneList
  }
}

export const DISPLAY_LIST = 'DISPLAY_LIST'
export function displayList(projectID, listID) {
  return dispatch => {
    dispatch(requestDisplayList())
    return axios
      .patch(`/api/projects/${projectID}/lists/${listID}/display`)
      .then(res => {
        dispatch(receiveDisplayList(res.data))
      })
      .catch(err => {
        ErrorHandlerWithoutSubmission(err)
          .then()
          .catch(error => {
            dispatch(ServerError(error))
          })
      })
  }
}

export const OPEN_SHOW_TASK = 'OPEN_SHOW_TASK'
export function openShowTaskModal(task) {
  return {
    type: OPEN_SHOW_TASK,
    task: task
  }
}

export const OPEN_DELETE_PROJECT = 'OPEN_DELETE_PROJECT'
export function openDeleteProjectModal() {
  return {
    type: OPEN_DELETE_PROJECT
  }
}
