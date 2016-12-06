import expect from 'expect'
import nock from 'nock'
import mockStore from '../../support/MockStore'
import * as newProjectModalActions from '../../../frontend/javascripts/actions/ProjectAction/NewProjectModalAction'

describe('closeNewProjectModal', () => {
  it('should close new project modal', () => {
    const expectedAction = {
      type: newProjectModalActions.CLOSE_NEW_PROJECT,
    }
    expect(newProjectModalActions.closeNewProjectModal()).toEqual(expectedAction)
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
    }
    const postForm = `title=${title}&description=${description}&repository_id=${repository.id}`
    beforeEach(() => {
      nock('http://localhost')
        .post('/projects', postForm)
        .reply(201, {
          ok: true,
          ID: 1,
          UserID: 12,
          Title: title,
          Description: description
        })
    })

    it('call RECEIVE_CREATE_PROJECT and get project', (done) => {
      const expectedActions = [
        { type: newProjectModalActions.REQUEST_CREATE_PROJECT },
        { type: newProjectModalActions.RECEIVE_CREATE_PROJECT, project: {ID: 1, UserID: 12, Title: title, Description: description } }
      ]
      const store = mockStore({ project: null }, expectedActions, done)
      const params = {
        title: title,
        description: description,
        repository_id: repository.id,
      }
      store.dispatch(newProjectModalActions.fetchCreateProject(params))
    })
  })
})
