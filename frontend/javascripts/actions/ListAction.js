import Request from 'superagent';

export const OPEN_NEW_LIST = 'PEN_NEW_LIST';
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
    lists: lists
  };
}

export function fetchLists(projectId) {
  return dispatch => {
    dispatch(requestLists());
    return Request
      .get(`/projects/${projectId}/lists`)
      .end((err, res)=> {
        if (res.body != null) {
          dispatch(receiveLists(res.body));
        }
      });
  };
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
        if (res.body != null) {
          dispatch(receiveProject(res.body));
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
  };
}

export function fetchCreateList(projectId, title, color) {
  return dispatch => {
    dispatch(requestCreateList());
    return Request
      .post(`/projects/${projectId}/lists`)
      .type('form')
      .send({title: title, color: color})
      .end((err, res)=> {
        if(res.body != null) {
          dispatch(receiveCreateList(res.body));
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
    task: {Id: task.Id, ListId: task.ListId, Title: task.Title }
  };
}

export function fetchCreateTask(projectId, listId, title) {
  return dispatch => {
    dispatch(requestCreateTask());
    return Request
      .post(`/projects/${projectId}/lists/${listId}/tasks`)
      .type('form')
      .send({title: title})
      .end((err, res)=> {
        if(res.body != null) {
          dispatch(receiveCreateTask(res.body));
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
    list: {Id: list.Id, ProjectId: list.ProjectId, Title: list.Title, Color: list.Color, ListTasks: list.ListTasks}
  };
}

export function fetchUpdateList(projectId, list) {
  return dispatch => {
    console.log(list);
    dispatch(requestUpdateList());
    return Request
      .post(`/projects/${projectId}/lists/${list.Id}`)
      .type('form')
      .send({title: list.Title.String, color: list.Color.String})
      .end((err, res)=> {
        if(res.body != null) {
          dispatch(receiveUpdateList(res.body));
        }
      });
  };
}
