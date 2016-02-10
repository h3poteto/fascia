import Request from 'superagent';

export const UNAUTHORIZED = 'UNAUTHORIZED';
function unauthorized() {
  window.location.pathname = "/sign_in";
  return {
    type: UNAUTHORIZED
  };
}

export const NOT_FOUND = 'NOT_FOUND';
function notFound() {
  return {
    type: NOT_FOUND
  }
}

export const SERVER_ERROR = 'SERVER_ERROR';
function serverError() {
  return {
    type: SERVER_ERROR
  };
}

export const CLOSE_FLASH = "CLOSE_FLASH";
export function closeFlash() {
  return {
    type: CLOSE_FLASH
  };
}

export const OPEN_NEW_LIST = 'OPEN_NEW_LIST';
export function openNewListModal() {
  return {
    type: OPEN_NEW_LIST,
    isListModalOpen: true
  };
}

export const CLOSE_NEW_LIST = 'CLOSE_NEW_LIST';
export function closeNewListModal() {
  return {
    type: CLOSE_NEW_LIST,
    isListModalOpen: false
  };
}

export const OPEN_NEW_TASK = 'OPEN_NEW_TASK';
export function openNewTaskModal(list) {
  return {
    type: OPEN_NEW_TASK,
    isTaskModalOpen: true,
    list: list
  };
}

export const CLOSE_NEW_TASK = 'CLOSE_NEW_TASK';
export function closeNewTaskModal() {
  return {
    type: CLOSE_NEW_TASK,
    isTaskModalOpen: false
  };
}

export const OPEN_EDIT_LIST = 'OPEN_EDIT_LIST';
export function openEditListModal(list) {
  return {
    type: OPEN_EDIT_LIST,
    isListEditModalOpen: true,
    list: list
  };
}

export const CLOSE_EDIT_LIST = 'CLOSE_EDIT_LIST';
export function closeEditListModal() {
  return {
    type: CLOSE_EDIT_LIST,
    isListEditModalOpen: false
  };
}

export const UPDATE_NEW_LIST_TITLE = 'UPDATE_NEW_LIST_TITLE';
export function updateNewListTitle(ev) {
  return {
    type: UPDATE_NEW_LIST_TITLE,
    title: ev.target.value
  };
}

export const UPDATE_NEW_LIST_COLOR = 'UPDATE_NEW_LIST_COLOR';
export function updateNewListColor(ev) {
  return {
    type: UPDATE_NEW_LIST_COLOR,
    color: ev.target.value
  };
}

export const UPDATE_SELECTED_LIST_TITLE = 'UPDATE_SELECTED_LIST_TITLE';
export function updateSelectedListTitle(ev) {
  return {
    type: UPDATE_SELECTED_LIST_TITLE,
    title: ev.target.value
  };
}

export const UPDATE_SELECTED_LIST_COLOR = 'UPDATE_SELECTED_LIST_COLOR';
export function updateSelectedListColor(ev) {
  return {
    type: UPDATE_SELECTED_LIST_COLOR,
    color: ev.target.value
  };
}

export const UPDATE_NEW_TASK_TITLE = 'UPDATE_NEW_TASK_TITLE';
export function updateNewTaskTitle(ev) {
  return {
    type: UPDATE_NEW_TASK_TITLE,
    title: ev.target.value
  };
}

export const UPDATE_NEW_TASK_DESCRIPTION = 'UPDATE_NEW_TASK_DESCRIPTION'
export function updateNewTaskDescription(ev) {
  return {
    type: UPDATE_NEW_TASK_DESCRIPTION,
    description: ev.target.value
  }
}

export const REQUEST_LISTS = 'REQUEST_LISTS';
function requestLists() {
  return {
    type: REQUEST_LISTS
  };
}

export const RECEIVE_LISTS = 'RECEIVE_LISTS';
function receiveLists(lists) {
  return {
    type: RECEIVE_LISTS,
    lists: lists.Lists,
    noneList: lists.NoneList
  }
}

export function fetchLists(projectId) {
  return dispatch => {
    dispatch(requestLists())
    return Request
      .get(`/projects/${projectId}/lists`)
      .end((err, res)=> {
        if (res.ok) {
          dispatch(receiveLists(res.body))
        } else if (res.unauthorized) {
          dispatch(unauthorized())
        } else if (res.notFound) {
          dispatch(notFound())
        } else {
          dispatch(serverError())
        }
      })
  }
}

export const REQUEST_PROJECT = 'REQUEST_PROJECT';
function requestProject() {
  return {
    type: REQUEST_PROJECT
  };
}

export const RECEIVE_PROJECT = 'RECEIVE_PROJECT';
function receiveProject(project) {
  return {
    type: RECEIVE_PROJECT,
    project: project
  };
}

export function fetchProject(projectId) {
  return dispatch => {
    dispatch(requestProject());
    return Request
      .get(`/projects/${projectId}/show`)
      .end((err, res)=> {
        if (res.ok) {
          dispatch(receiveProject(res.body));
        } else if (res.unauthorized) {
          dispatch(unauthorized());
        } else if (res.notFound) {
          dispatch(notFound())
        } else {
          dispatch(serverError());
        }
      });
  };
}


export const REQUEST_CREATE_LIST = 'REQUEST_CREATE_LIST';
function requestCreateList() {
  return {
    type: REQUEST_CREATE_LIST
  };
}

export const RECEIVE_CREATE_LIST = 'RECEIVE_CREATE_LIST';
function receiveCreateList(list) {
  return {
    type: RECEIVE_CREATE_LIST,
    list: {Id: list.Id, ProjectId: list.ProjectId, Title: list.Title, Color: list.Color, ListTasks: list.ListTasks}
  }
}

export function fetchCreateList(projectId, title, color) {
  return dispatch => {
    dispatch(requestCreateList());
    return Request
      .post(`/projects/${projectId}/lists`)
      .type('form')
      .send({title: title, color: color})
      .end((err, res)=> {
        if(res.ok) {
          dispatch(receiveCreateList(res.body));
        } else if (res.unauthorized) {
          dispatch(unauthorized());
        } else if (res.notFound) {
          dispatch(notFound())
        } else {
          dispatch(serverError());
        }
      });
  };
}

export const REQUEST_CREATE_TASK = 'REQUEST_CREATE_TASK';
function requestCreateTask() {
  return {
    type: REQUEST_CREATE_TASK
  };
}

export const RECEIVE_CREATE_TASK = 'RECEIVE_CREATE_TASK';
function receiveCreateTask(task) {
  return {
    type: RECEIVE_CREATE_TASK,
    task: task
  };
}

export function fetchCreateTask(projectId, listId, title, description) {
  return dispatch => {
    dispatch(requestCreateTask());
    return Request
      .post(`/projects/${projectId}/lists/${listId}/tasks`)
      .type('form')
      .send({title: title, description: description})
      .end((err, res)=> {
        if(res.ok) {
          dispatch(receiveCreateTask(res.body));
        } else if (res.unauthorized) {
          dispatch(unauthorized());
        } else if (res.notFound) {
          dispatch(notFound())
        } else {
          dispatch(serverError());
        }
      });
  };
}

export const REQUEST_UPDATE_LIST = 'REQUEST_UPDATE_LIST';
function requestUpdateList() {
  return {
    type: REQUEST_UPDATE_LIST
  };
}

export const RECEIVE_UPDATE_LIST = 'RECEIVE_UPDATE_LIST';
function receiveUpdateList(list) {
  return {
    type: RECEIVE_UPDATE_LIST,
    list: {Id: list.Id, ProjectId: list.ProjectId, Title: list.Title, Color: list.Color, ListTasks: list.ListTasks, ListOptionId: list.ListOptionId}
  };
}

export function fetchUpdateList(projectId, list, option) {
  var action
  if (option != undefined && option != null) {
    action = option.Action
  }
  return dispatch => {
    dispatch(requestUpdateList());
    return Request
      .post(`/projects/${projectId}/lists/${list.Id}`)
      .type('form')
      .send({title: list.Title, color: list.Color, action: action})
      .end((err, res)=> {
        if(res.ok) {
          dispatch(receiveUpdateList(res.body));
        } else if (res.unauthorized) {
          dispatch(unauthorized());
        } else if (res.notFound) {
          dispatch(notFound())
        } else {
          dispatch(serverError());
        }
      });
  };
}


export const TASK_DRAG_START = "TASK_DRAG_START";
export function taskDragStart(ev) {
  ev.dataTransfer.effectAllowed = "moved";
  ev.dataTransfer.setData("text/html", ev.currentTarget);
  return {
    type: TASK_DRAG_START,
    taskDragTarget: ev.currentTarget,
    taskDragFromList: ev.currentTarget.parentNode.parentNode
  };
}

export const TASK_DRAG_LEAVE = "TASK_DRAG_LEAVE";
export const TASK_DRAG_LEAVE_IGNORE = "TASK_DRAG_LEAVE_IGNORE"
export function taskDragLeave(ev) {
  // li.new-taskだけはdragleaveイベントが頻繁に発生するため，抑制する
  // こうしておいても，他の要素に移動した際には問題なくleave処理が為される
  if (ev.target.className == "new-task" || ev.target.className == "fa-plus") {
    return {
      type: TASK_DRAG_LEAVE_IGNORE
    }
  } else {
    return {
      type: TASK_DRAG_LEAVE
    }
  }
}

export const REQUEST_MOVE_TASK = 'REQUEST_MOVE_TASK';
function requestMoveTask() {
  return {
    type: REQUEST_MOVE_TASK
  };
}

export const RECEIVE_MOVE_TASK = 'RECEIVE_MOVE_TASK';
function receiveMoveTask(lists) {
  return {
    type: RECEIVE_MOVE_TASK,
    lists: lists.Lists,
    noneList: lists.NoneList
  }
}

export const TASK_DROP = 'TASK_DROP';
export function taskDrop(projectId, taskDraggingFrom, taskDraggingTo) {
  if (taskDraggingTo != undefined && taskDraggingTo != null) {
    var prevToTaskId;
    if (taskDraggingTo.prevToTask == null) {
      prevToTaskId = null;
    } else {
      prevToTaskId = taskDraggingTo.prevToTask.Id;
    }
    return dispatch => {
      dispatch(requestMoveTask());
      return Request
        .post(`/projects/${projectId}/lists/${taskDraggingFrom.fromList.Id}/tasks/${taskDraggingFrom.fromTask.Id}/move_task`)
        .type('form')
        .send({to_list_id: taskDraggingTo.toList.Id, prev_to_task_id: prevToTaskId})
        .end((err, res) => {
          if(res.ok) {
            dispatch(receiveMoveTask(res.body));
          } else if(res.unauthorized) {
            dispatch(unauthorized());
          } else if (res.notFound) {
            dispatch(notFound())
          } else {
            // TODO: ここはドラッグしたviewを元に戻す必要がある
            dispatch(serverError());
          }
        });
    };
  } else {
    return {
      type: TASK_DROP
    };
  }
}

export const TASK_DRAG_OVER = "TASK_DRAG_OVER";
export function taskDragOver(ev) {
  ev.preventDefault();
  var targetList;
  switch(ev.target.dataset.droppedDepth) {
  case "0":
    targetList = ev.target;
    break;
  case "1":
    targetList = ev.target.parentNode;
    break;
  case "2":
    targetList = ev.target.parentNode.parentNode;
    break;
  case "3":
    targetList = ev.target.parentNode.parentNode.parentNode;
    break;
  default:
    targetList = ev.target.parentNode.parentNode;
    break;
  }
  return {
    type: TASK_DRAG_OVER,
    taskDragToTask: ev.target,
    taskDragToList: targetList
  };
}

export const REQUEST_FETCH_GITHUB = "REQUEST_FETCH_GITHUB";
function requestFetchGithub() {
  return {
    type: REQUEST_FETCH_GITHUB
  };
}

export const RECEIVE_FETCH_GITHUB = "RECEIVE_FETCH_GITHUB";
function receiveFetchGithub(lists) {
  return {
    type: RECEIVE_FETCH_GITHUB,
    lists: lists.Lists,
    noneList: lists.NoneList
  }
}

export const FETCH_PROJECT_GITHUB = "FETCH_PROJECT_GITHUB";
export function fetchProjectGithub(projectId) {
  return dispatch => {
    dispatch(requestFetchGithub());
    return Request
      .post(`/projects/${projectId}/fetch_github`)
      .end((err, res) => {
        if (res.ok) {
          dispatch(receiveFetchGithub(res.body));
        } else if(res.unauthorized) {
          dispatch(unauthorized());
        } else if (res.notFound) {
          dispatch(notFound())
        } else {
          dispatch(serverError());
        }
      });
  };
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
    return Request
      .get('/list_options')
      .end((err, res) => {
        if (res.ok) {
          dispatch(receiveListOptions(res.body))
        } else if (res.unauthorized) {
          dispatch(unauthorized())
        } else if (res.notFound) {
          dispatch(notFound())
        } else {
          dispatch(serverError())
        }
      })
  }
}

export const CHANGE_SELECTED_LIST_OPTION = 'CHANGE_SELECTED_LIST_OPTION'
export function changeSelectedListOption(ev) {
  return {
    type: CHANGE_SELECTED_LIST_OPTION,
    selectEvent: ev.target
  }
}


export const OPEN_EDIT_PROJECT = 'OPEN_EDIT_PROJECT';
export function openEditProjectModal(project) {
  return {
    type: OPEN_EDIT_PROJECT,
    project: project
  }
}

export const CLOSE_EDIT_PROJECT = 'CLOSE_EDIT_PROJECT';
export function closeEditProjectModal() {
  return {
    type: CLOSE_EDIT_PROJECT
  }
}


export const UPDATE_EDIT_PROJECT_TITLE = 'UPDATE_EDIT_PROJECT_TITLE'
export function updateEditProjectTitle(ev) {
  return {
    type: UPDATE_EDIT_PROJECT_TITLE,
    title: ev.target.value
  }
}

export const UPDATE_EDIT_PROJECT_DESCRIPTION = 'UPDATE_EDIT_PROJECT_DESCRIPTION'
export function updateEditProjectDescription(ev) {
  return {
    type: UPDATE_EDIT_PROJECT_DESCRIPTION,
    description: ev.target.value
  }
}

export const REQUEST_UPDATE_PROJECT = 'REQUEST_UPDATE_PROJECT'
function requestUpdateProject() {
  return {
    type: REQUEST_UPDATE_PROJECT
  }
}

export const RECEIVE_UPDATE_PROJECT = 'RECEIVE_UPDATE_PROJECT'
function receiveUpdateProject(project) {
  return {
    type: RECEIVE_UPDATE_PROJECT,
    project: project
  }
}

export const FETCH_UPDATE_PROJECT = 'FETCH_UPDATE_PROJECT'
export function fetchUpdateProject(projectID, project) {
  return dispatch => {
    dispatch(requestUpdateProject())
    return Request
      .post(`/projects/${projectID}`)
      .type('form')
      .send({title: project.Title, description: project.Description})
      .end((err, res) => {
        if (res.ok) {
          dispatch(receiveUpdateProject(res.body))
        } else if (res.unauthorized) {
          dispatch(unauthorized())
        } else if (res.notFound) {
          dispatch(notFound())
        } else {
          dispatch(serverError())
        }
      })
  }
}
