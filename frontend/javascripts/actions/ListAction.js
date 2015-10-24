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
