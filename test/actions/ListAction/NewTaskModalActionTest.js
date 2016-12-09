import * as newTaskModalActions from '../../../frontend/javascripts/actions/ListAction/NewTaskModalAction'
import expect from 'expect'
import nock from 'nock'
import mockStore from '../../support/MockStore'

describe('closeNewTaskModal', () => {
  it('should close new task modal', () => {
    const expectedAction = {
      type: newTaskModalActions.CLOSE_NEW_TASK,
      isTaskModalOpen: false
    }
    expect(newTaskModalActions.closeNewTaskModal()).toEqual(expectedAction)
  })
})

describe('fetchCreateTask', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    const projectID = 1
    const listID = 2
    const title = "taskTitle"
    const description = "taskDescription"
    const postForm = `title=${title}&description=${description}`
    beforeEach(() => {
      nock('http://localhost')
        .post(`/projects/${projectID}/lists/${listID}/tasks`, postForm)
        .reply(200, {
          Lists: [
            {
              ID: listID,
              Title: title,
              ListTasks: [
                {
                  ID: 1,
                  ListID: listID,
                  Title: title,
                  Description: description
                }
              ]
            }
          ],
          NoneList: []
        })
    })

    it('call RECEIVE_CREATE_TASK and get task', (done) => {
      const expectedActions = [
        { type: newTaskModalActions.REQUEST_CREATE_TASK },
        { type: newTaskModalActions.RECEIVE_CREATE_TASK,
          lists: [
            {
              ID: listID,
              Title: title,
              ListTasks: [
                {
                  ID: 1,
                  ListID: listID,
                  Title: title,
                  Description: description
                }
              ]
            }
          ],
          noneList: []
        }
      ]
      const store = mockStore({ task: null }, expectedActions, done)
      const params = {
        title: title,
        description: description,
      }
      store.dispatch(newTaskModalActions.fetchCreateTask(projectID, listID, params))
    })
  })
})
