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
function receiveCreateProject(id, userId, title, description) {
  return {
    type: RECEIVE_CREATE_PROJECT,
    project: {Id: id, UserId: userId, Title: title, Description: description}
  };
}


export function fetchCreateProject(title, description, repository) {
  return dispatch => {
    dispatch(requestCreateProject());
    return Request
      .post('/projects')
      .type('form')
      .send({title: title, description: description, repository: repository})
      .end((err, res)=> {
        if (res.body != null) {
          dispatch(receiveCreateProject(res.body.Id, res.body.UserId, res.body.Title, res.body.Description));
        }
      });
    };
}

export const UPDATE_NEW_PROJECT_TITLE = 'UPDATE_NEW_PROJECT_TITLE';
export function updateNewProjectTitle(ev) {
  return {
    type: UPDATE_NEW_PROJECT_TITLE,
    title: ev.target.value
  };
}

export const UPDATE_NEW_PROJECT_DESCRIPTION = 'UPDATE_NEW_PROJECT_DESCRIPTION';
export function updateNewProjectDescription(ev) {
  return {
    type: UPDATE_NEW_PROJECT_DESCRIPTION,
    description: ev.target.value
  };
}
