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


describe('fetchSession', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    beforeEach(() => {
      nock('http://localhost')
        .post('/session')
        .reply(200, {
          ok: true
        })
    })

    it('call RECEIVE_SESSION', (done) => {
      const expectedActions = [
        { type: projectActions.REQUEST_SESSION },
        { type: projectActions.RECEIVE_SESSION }
      ]
      const store = mockStore({}, expectedActions, done)
      store.dispatch(projectActions.fetchSession())
    })
  })
})
