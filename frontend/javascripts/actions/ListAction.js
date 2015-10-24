import Request from 'superagent';

export const OPEN_NEW_LIST = 'PEN_NEW_LIST';
export function openNewListModal() {
  return {
    type: OPEN_NEW_LIST,
    isModalOpen: true
  };
}

export const CLOSE_NEW_LIST = 'CLOSE_NEW_LIST';
export function closeNewListModal() {
  return {
    type: CLOSE_NEW_LIST,
    isModalOpen: false
  };
}

export const UPDATE_NEW_LIST_TITLE = 'UPDATE_NEW_LIST_TITLE';
export function updateNewListTitle(ev) {
  return {
    type: UPDATE_NEW_LIST_TITLE,
    title: ev.target.value
  };
}

export const REQUEST_POSTS = 'REQUEST_POSTS';
function requestPosts() {
  return {
    type: REQUEST_POSTS
  };
}

export const RECEIVE_POSTS = 'RECEIVE_POSTS';
function receivePosts(lists) {
  return {
    type: RECEIVE_POSTS,
    lists: lists
  };
}

export function fetchLists(projectId) {
  return dispatch => {
    dispatch(requestPosts());
    return Request
      .get(`/projects/${projectId}/lists`)
      .end((err, res)=> {
        if (res.body != null) {
          dispatch(receivePosts(res.body));
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
    list: {Id: list.Id, ProjectId: list.ProjectId, Title: list.Title}
  };
}

export function fetchCreateList(projectId, title) {
  return dispatch => {
    dispatch(requestCreateList());
    return Request
      .post(`/projects/${projectId}/lists`)
      .type('form')
      .send({title: title})
      .end((err, res)=> {
        if(res.body != null) {
          dispatch(receiveCreateList(res.body));
        }
      });
  };
}
