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
        .reply(200, { projects: ['do something'] })
    })

    it('creates RECEIVE_POSTS when fetching projects has been done', (done) => {
      const expectedActions = [
        { type: projectActions.REQUEST_POSTS },
        { type: projectActions.RECEIVE_POSTS, projects: { projects: ['do something']  } }
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

    it('creates RECEIVE_REPOSITORIES', (done) => {
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

    beforeEach(() => {
      // TODO: できれば文字列じゃなくてハッシュでやりたいけど，とりあえずこれで通しておく
      const postForm = "title=projectTitle&description=projectDescription&repositoryId=1&repositoryOwner=ownerName&repositoryName=repo1"

      nock('http://localhost')
        .post('/projects', (body) => {
          return JSON.stringify(body) === JSON.stringify(postForm)
        })
        .reply(201, {
          ok: true,
          Id: 1,
          UserId: 12,
          Title: 'projectTitle',
          Description: 'projectDescription'
        })
    })

    it('creates RECEIVE_CREATE_PROJECT', (done) => {
      const expectedActions = [
        { type: projectActions.REQUEST_CREATE_PROJECT },
        { type: projectActions.RECEIVE_CREATE_PROJECT, project: {Id: 1, UserId: 12, Title: 'projectTitle', Description: 'projectDescription' } }
      ]
      const store = mockStore({ project: null }, expectedActions, done)
      const repository = {
        id: 1,
        name: "repo1",
        owner: {
          login: "ownerName"
        }
      }
      store.dispatch(projectActions.fetchCreateProject("projectTitle", "projectDescription", repository))
    })
  })
})

// TODO: update系のテストが足りてない
