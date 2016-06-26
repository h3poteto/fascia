import TestUtils from 'react-addons-test-utils'
import expect from 'expect'
import React from 'react'
import NewTaskModal from '../../../frontend/javascripts/components/ListView/NewTaskModal.jsx'

function setup(props) {
  let renderer = TestUtils.createRenderer()
  renderer.render(<NewTaskModal {...props} />)
  let output = renderer.getRenderOutput()

  return {
    props,
    output,
    renderer
  }
}

describe('ListView::NewTaskModal', () => {
  context('when task modal open', () => {
    let state = {
      isTaskModalOpen: true,
      newTask: {title: ""},
      selectedList: 1,
      projectID: 1,
      closeNewTaskModal: expect.createSpy(),
      updateNewTaskTitle: expect.createSpy(),
      updateNewTaskDescription: expect.createSpy(),
      fetchCreateTask: expect.createSpy()
    }
    it('should render task modal', () => {
      const { output } = setup(state)
      expect(output.props.isOpen).toBe(true)
    })
  })
})
