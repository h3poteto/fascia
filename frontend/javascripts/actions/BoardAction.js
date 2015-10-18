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
      .get('/projects')
      .end((err, res)=> {
        if (res.body != null) {
          dispatch(receivePosts(res.body));
        }
      });
  };
}

export const OPEN_NEW_PROJECT = 'OPEN_NEW_PROJECT';
export function openNewProjectModal() {
  return {
    type: OPEN_NEW_PROJECT,
    isModalOpen: true
  };
}

export const CLOSE_NEW_PROJECT = 'CLOSE_NEW_PROJECT';
export function closeNewProjectModal() {
  return {
    type: CLOSE_NEW_PROJECT,
    isModalOpen: false
  };
}

export const REQUEST_REPOSITORIES = 'REQUEST_REPOSITORIES';
function requestRepositories() {
  return {
    type: REQUEST_REPOSITORIES
  };
}

export const RECEIVE_REPOSITORIES = 'RECEIVE_REPOSITORIES';
function receiveRepositories(repositories) {
  return {
    type: RECEIVE_REPOSITORIES,
    repositories: repositories
  };
}

export function fetchRepositories() {
  return dispatch => {
    dispatch(requestRepositories());
    return Request
      .get('/github/repositories')
      .end((err, res)=> {
        if (res.body != null) {
          dispatch(receiveRepositories(res.body));
        }
      });
  };
}

export const CHANGE_SELECT_REPOSITORY = 'CHANGE_SELECT_REPOSITORY';
export function changeSelectedRepository(ev) {
  return {
    type: CHANGE_SELECT_REPOSITORY,
    selectEvent: ev.target
  };
}

export const REQUEST_CREATE_PROJECT = 'REQUEST_CREATE_PROJECT';
function requestCreateProject() {
  return {
    type: REQUEST_CREATE_PROJECT
  };
}

export const RECEIVE_CREATE_PROJECT = 'RECEIVE_CREATE_PROJECT';
function receiveCreateProject(id, userId, title) {
  return {
    type: RECEIVE_CREATE_PROJECT,
    project: {Id: id, UserId: userId, Title: title}
  };
}


export function fetchCreateProject(title, repository) {
  return dispatch => {
    dispatch(requestCreateProject());
    return Request
      .post('/projects')
      .type('form')
      .send({title: title, repository: repository})
      .end((err, res)=> {
        if (res.body != null) {
          dispatch(receiveCreateProject(res.body.Id, res.body.UserId, res.body.Title));
        }
      });
    };
}

export const UPDATE_NEW_PROJECT = 'UPDATE_NEW_PROJECT';
export function updateNewProject(ev) {
  return {
    type: UPDATE_NEW_PROJECT,
    title: ev.target.value
  };
}
