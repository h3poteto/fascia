import Request from 'superagent';

export const REQUEST_POSTS = 'REQUEST_POSTS';
function requestPosts() {
  return {
    type: REQUEST_POSTS
  };
}

export const RECEIVE_POSTS = 'RECEIVE_POSTS';
function receivePosts(projects) {
  return {
    type: RECEIVE_POSTS,
    projects: projects
  };
}

export function fetchProjects() {
  return dispatch => {
    dispatch(requestPosts());
    return Request
      .get('/projects/')
      .end((err, res)=> {
        if (res.body != null) {
          dispatch(receivePosts(res.body));
        }
      });
  };
}

export const NEW_PROJECT = 'NEW_PROJECT';
export function newProject() {
  return {
    type: NEW_PROJECT,
    isModalOpen: true
  };
}
