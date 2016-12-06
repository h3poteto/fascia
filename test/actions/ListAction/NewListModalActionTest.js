import * as newListModalActions from '../../../frontend/javascripts/actions/ListAction/NewListModalAction'
import expect from 'expect'
import nock from 'nock'
import mockStore from '../../support/MockStore'

describe('closeNewListModal', () => {
  it('should close new list modal', () => {
    const expectedAction = {
      type: newListModalActions.CLOSE_NEW_LIST,
    }
    expect(newListModalActions.closeNewListModal()).toEqual(expectedAction)
  })
})

describe('fetchCreateList', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    const projectID = 1
    const title = "listTitle"
    const color = "ffffff"
    const postForm = `title=${title}&color=${color}`
    beforeEach(() => {
      nock('http://localhost')
        .post(`/projects/${projectID}/lists`, postForm)
        .reply(200, {
          ID: 1,
          ProjectID: projectID,
          Title: title,
          Color: color,
          ListTasks: ["task1"]
        })
    })

    it('call RECEIVE_CREATE_LIST and get list', (done) => {
      const expectedActions = [
        { type: newListModalActions.REQUEST_CREATE_LIST },
        { type: newListModalActions.RECEIVE_CREATE_LIST, list: { ID: 1, ProjectID: projectID, Title: title, Color: color, ListTasks: ["task1"] } }
      ]
      const store = mockStore({ list: null }, expectedActions, done)
      const params = {
        title: title,
        color: color,
      }
      store.dispatch(newListModalActions.fetchCreateList(projectID, params))
    })
  })
})
