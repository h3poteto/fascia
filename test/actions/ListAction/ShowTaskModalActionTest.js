import * as showTaskModalActions from '../../../frontend/javascripts/actions/ListAction/ShowTaskModalAction'
import expect from 'expect'

describe('closeNewTaskModal', () => {
  it('should close new task modal', () => {
    const expectedAction = {
      type: showTaskModalActions.CLOSE_SHOW_TASK,
    }
    expect(showTaskModalActions.closeShowTaskModal()).toEqual(expectedAction)
  })
})
