import * as newListModalActions from '../../../frontend/javascripts/actions/ListAction/NewListModalAction'
import expect from 'expect'
import nock from 'nock'
import mockStore from '../../support/MockStore'

describe('closeNewListModal', () => {
  it('should close new list modal', () => {
    const expectedAction = {
      type: newListModalActions.CLOSE_NEW_LIST,
      isListModalOpen: false
    }
    expect(newListModalActions.closeNewListModal()).toEqual(expectedAction)
  })
})

describe('updateNewListTitle', () => {
  it('should update new list title', () => {
    const title = "newTitle"
    const ev = {
      target: {
        value: title
      }
    }
    const expectedAction = {
      type: newListModalActions.UPDATE_NEW_LIST_TITLE,
      title: title
    }
    expect(newListModalActions.updateNewListTitle(ev)).toEqual(expectedAction)
  })
})

describe('updateNewListColor', () => {
  it('should update new list color', () => {
    const color = "ffffff"
    const ev = {
      target: {
        value: color
      }
    }
    const expectedAction = {
      type: newListModalActions.UPDATE_NEW_LIST_COLOR,
      color: color
    }
    expect(newListModalActions.updateNewListColor(ev)).toEqual(expectedAction)
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
      store.dispatch(newListModalActions.fetchCreateList(projectID, title, color))
    })
  })
})
