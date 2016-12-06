import * as editListModalActions from '../../../frontend/javascripts/actions/ListAction/EditListModalAction'
import expect from 'expect'
import nock from 'nock'
import mockStore from '../../support/MockStore'

describe('closeEditListModal', () => {
  it('should close edit list modal', () => {
    const expectedAction = {
      type: editListModalActions.CLOSE_EDIT_LIST,
      isListEditModalOpen: false
    }
    expect(editListModalActions.closeEditListModal()).toEqual(expectedAction)
  })
})

describe('fetchUpdateList', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    const projectID = 1
    const option = {
      ID: 1,
      Action: "close"
    }
    const list = {
      ID: 2,
      Title: "listTitle",
      Color: "ffffff",
      ProjectID: projectID,
      ListOptionID: option.ID
    }
    const postForm = `title=${list.Title}&color=${list.Color}&option_id=${option.ID}`
    beforeEach(() => {
      nock('http://localhost')
        .post(`/projects/${projectID}/lists/${list.ID}`, postForm)
        .reply(200, {
          ID: list.ID,
          ProjectID: list.ProjectID,
          Title: list.Title,
          Color: list.Color,
          ListTasks: [],
          ListOptionID: option.ID
        })
    })
    it('call RECEIVE_UPDATE_LIST and get list', (done) => {
      const expectedActions = [
        { type: editListModalActions.REQUEST_UPDATE_LIST },
        { type: editListModalActions.RECEIVE_UPDATE_LIST, list: { ID: list.ID, ProjectID: list.ProjectID, Title: list.Title, Color: list.Color, ListTasks: [], ListOptionID: option.ID } }
      ]
      const store = mockStore({ list: null }, expectedActions, done)
      const params = {
        title: list.Title,
        color: list.Color,
        option_id: option.ID,
      }
      store.dispatch(editListModalActions.fetchUpdateList(projectID, list.ID, params))
    })
  })
})
