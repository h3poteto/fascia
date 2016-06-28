import expect from 'expect'
import nock from 'nock'
import mockStore from '../../support/MockStore'
import * as newProjectModalActions from '../../../frontend/javascripts/actions/ProjectAction/NewProjectModalAction'

describe('closeNewProjectModal', () => {
  it('should close new project modal', () => {
    const expectedAction = {
      type: newProjectModalActions.CLOSE_NEW_PROJECT,
      isModalOpen: false
    }
    expect(newProjectModalActions.closeNewProjectModal()).toEqual(expectedAction)
  })
})

describe('changeSelectedRepository', () => {
  it('should change selected repository', () => {
    const expectedAction = {
      type: newProjectModalActions.CHANGE_SELECT_REPOSITORY,
      selectEvent: "<element>"
    }
    const event = {
      target: "<element>"
    }
    expect(newProjectModalActions.changeSelectedRepository(event)).toEqual(expectedAction)
  })
})

describe('updateNewProjectTitle', () => {
  it('should update new project title', () => {
    const expectedAction = {
      type: newProjectModalActions.UPDATE_NEW_PROJECT_TITLE,
      title: "projectTitle"
    }
    const event = {
      target: {
        value: "projectTitle"
      }
    }
    expect(newProjectModalActions.updateNewProjectTitle(event)).toEqual(expectedAction)
  })
})

describe('updateNewProjectDescription', () => {
  it('should update new project description', () => {
    const expectedAction = {
      type: newProjectModalActions.UPDATE_NEW_PROJECT_DESCRIPTION,
      description: "projectDescription"
    }
    const event = {
      target: {
        value: "projectDescription"
      }
    }
    expect(newProjectModalActions.updateNewProjectDescription(event)).toEqual(expectedAction)
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
    const postForm = `title=${title}&description=${description}&repository_id=${repository.id}&repository_owner=${repository.owner.login}&repository_name=${repository.name}`
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
      store.dispatch(newProjectModalActions.fetchCreateProject(title, description, repository))
    })
  })
})
