import * as editProjectModalActions from '../../../frontend/javascripts/actions/ListAction/EditProjectModalAction'
import expect from 'expect'
import nock from 'nock'
import mockStore from '../../support/MockStore'

describe('closeEditProjectModal', () => {
  it('should close edit project modal', () => {
    const expectedAction = {
      type: editProjectModalActions.CLOSE_EDIT_PROJECT
    }
    expect(editProjectModalActions.closeEditProjectModal()).toEqual(expectedAction)
  })
})

describe('fetchUpdateProject', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    const project = {
      ID: 1,
      Title: "title",
      Description: "description"
    }
    const postForm = `title=${project.Title}&description=${project.Description}`
    beforeEach(() => {
      nock('http://localhost')
        .post(`/projects/${project.ID}`, postForm)
        .reply(200, {
          ID: project.ID,
          Title: project.Title,
          Description: project.Description
        })
    })
    it('call RECEIVE_UPDATE_PROJECT and get project', (done) => {
      const expectedActions = [
        { type: editProjectModalActions.REQUEST_UPDATE_PROJECT },
        { type: editProjectModalActions.RECEIVE_UPDATE_PROJECT, project: project }
      ]
      const store = mockStore({ project: null }, expectedActions, done)
      store.dispatch(editProjectModalActions.fetchUpdateProject(project.ID, project))
    })
  })
})

describe('createWebhook', () => {
  afterEach(() => {
    nock.cleanAll
  })
  context('when response is right', () => {
    const projectID = 1
    beforeEach(() => {
      nock('http://localhost')
        .post(`/projects/${projectID}/webhook`)
        .reply(200, {
        })
    })
    it('call RECEIVE_CREATE_WEBHOOK', (done) => {
      const expectedActions = [
        { type: editProjectModalActions.REQUEST_CREATE_WEBHOOK },
        { type: editProjectModalActions.RECEIVE_CREATE_WEBHOOK }
      ]
      const store = mockStore({}, expectedActions, done)
      store.dispatch(editProjectModalActions.createWebhook(projectID))
    })
  })
})
