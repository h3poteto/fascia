import * as showTaskModalActions from '../../../frontend/javascripts/actions/ListAction/ShowTaskModalAction'
import expect from 'expect'
import nock from 'nock'
import mockStore from '../../support/MockStore'

describe('closeNewTaskModal', () => {
  it('should close new task modal', () => {
    const expectedAction = {
      type: showTaskModalActions.CLOSE_SHOW_TASK
    }
    expect(showTaskModalActions.closeShowTaskModal()).toEqual(expectedAction)
  })
})

describe('fetchUpdateTask', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    const projectID = 1
    const listID = 2
    const taskID = 3
    const title = "taskTitle"
    const description = "taskDescription"
    const postForm = `title=${title}&description=${description}`
    beforeEach(() => {
      nock('http://localhost')
        .post(`/projects/${projectID}/lists/${listID}/tasks/${taskID}`, postForm)
        .reply(200, {
          Lists: [
            {
              ID: listID,
              Title: title,
              ListTasks: [
                {
                  ID: taskID,
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
    it('call RECEIVE_UPDATE_TASK and get lists', (done) => {
      const expectedActions = [
        { type: showTaskModalActions.REQUEST_UPDATE_TASK },
        { type: showTaskModalActions.RECEIVE_UPDATE_TASK,
          lists: [
            {
              ID: listID,
              Title: title,
              ListTasks: [
                {
                  ID: taskID,
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
      const store = mockStore({}, expectedActions, done)
      const params = {
        title: title,
        description: description,
      }
      store.dispatch(showTaskModalActions.fetchUpdateTask(projectID, listID, taskID, params))
    })
  })
})
