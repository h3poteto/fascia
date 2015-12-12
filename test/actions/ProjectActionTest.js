import * as projectActions from '../../frontend/javascripts/actions/ProjectAction'
import expect from 'expect'
import nock from 'nock'
import mockStore from '../support/MockStore'

// normal tests

describe('closeFlash', () => {
  it('should close flash', () => {
    const expectedAction = {
      type: projectActions.CLOSE_FLASH
    }
    expect(projectActions.closeFlash()).toEqual(expectedAction)
  })
})

describe('openNewProjectModal', () => {
  it('should open new project modal', () => {
    const expectedAction = {
      type: projectActions.OPEN_NEW_PROJECT,
      isModalOpen: true
    }
    expect(projectActions.openNewProjectModal()).toEqual(expectedAction)
  })
})

describe('closeNewProjectModal', () => {
  it('should close new project modal', () => {
    const expectedAction = {
      type: projectActions.CLOSE_NEW_PROJECT,
      isModalOpen: false
    }
    expect(projectActions.closeNewProjectModal()).toEqual(expectedAction)
  })
})

describe('changeSelectedRepository', () => {
  it('should change selected repository', () => {
    const expectedAction = {
      type: projectActions.CHANGE_SELECT_REPOSITORY,
      selectEvent: "<element>"
    }
    const event = {
      target: "<element>"
    }
    expect(projectActions.changeSelectedRepository(event)).toEqual(expectedAction)
  })
})

describe('updateNewProjectTitle', () => {
  it('should update new project title', () => {
    const expectedAction = {
      type: projectActions.UPDATE_NEW_PROJECT_TITLE,
      title: "projectTitle"
    }
    const event = {
      target: {
        value: "projectTitle"
      }
    }
    expect(projectActions.updateNewProjectTitle(event)).toEqual(expectedAction)
  })
})

describe('updateNewProjectDescription', () => {
  it('should update new project description', () => {
    const expectedAction = {
      type: projectActions.UPDATE_NEW_PROJECT_DESCRIPTION,
      description: "projectDescription"
    }
    const event = {
      target: {
        value: "projectDescription"
      }
    }
    expect(projectActions.updateNewProjectDescription(event)).toEqual(expectedAction)
  })
})

// async tests
describe('fetchProjects', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    beforeEach(() => {
      nock('http://localhost')
        .get('/projects')
        .reply(200, { projects: ['project1', 'project2'] })
    })

    it('call RECEIVE_POSTS and get projects', (done) => {
      const expectedActions = [
        { type: projectActions.REQUEST_POSTS },
        { type: projectActions.RECEIVE_POSTS, projects: { projects: ['project1', 'project2']  } }
      ]
      const store = mockStore({ projects: [] }, expectedActions, done)
      store.dispatch(projectActions.fetchProjects())
    })
  })

  context('when response is invalid or server error', () => {
    beforeEach(() => {
      nock('http://localhost')
        .get('/projects')
        .reply(500, {})
    })

    it('called SERVER_ERROR', (done) => {
      const expectedActions = [
        { type: projectActions.REQUEST_POSTS },
        { type: projectActions.SERVER_ERROR }
      ]
      const store = mockStore({ projects: [] }, expectedActions, done)
      store.dispatch(projectActions.fetchProjects())
    })
  })
})

describe('fetchRepositories', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    beforeEach(() => {
      nock('http://localhost')
        .get('/github/repositories')
        .reply(200, { repositories: [ "repo1", "repo2" ] })
    })

    it('call RECEIVE_REPOSITORIES and get repositories', (done) => {
      const expectedActions = [
        { type: projectActions.REQUEST_REPOSITORIES },
        { type: projectActions.RECEIVE_REPOSITORIES, repositories: { repositories: [ "repo1", "repo2" ] } }
      ]
      const store = mockStore({ repositories: [] }, expectedActions, done)
      store.dispatch(projectActions.fetchRepositories())
    })
  })
})

describe('fetchCreateProject', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    const title = "projectTitle"
    const description = "projectDescription"
    const repository = {
      id: 1,
      name: "repo1",
      owner: {
        login: "ownerName"
      }
    }
    // TODO: できれば文字列じゃなくてハッシュでやりたいけど，とりあえずこれで通しておく
    const postForm = `title=${title}&description=${description}&repositoryId=${repository.id}&repositoryOwner=${repository.owner.login}&repositoryName=${repository.name}`
    beforeEach(() => {
      nock('http://localhost')
        .post('/projects', postForm)
        .reply(201, {
          ok: true,
          Id: 1,
          UserId: 12,
          Title: title,
          Description: description
        })
    })

    it('call RECEIVE_CREATE_PROJECT and get project', (done) => {
      const expectedActions = [
        { type: projectActions.REQUEST_CREATE_PROJECT },
        { type: projectActions.RECEIVE_CREATE_PROJECT, project: {Id: 1, UserId: 12, Title: title, Description: description } }
      ]
      const store = mockStore({ project: null }, expectedActions, done)
      store.dispatch(projectActions.fetchCreateProject(title, description, repository))
    })
  })
})

// TODO: update系のテストが足りてない
