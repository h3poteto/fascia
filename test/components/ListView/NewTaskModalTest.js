import TestUtils from 'react-addons-test-utils'
import expect from 'expect'
import React from 'react'
import NewTaskModal from '../../../frontend/javascripts/components/ListView/NewTaskModal.jsx'
import * as NewTaskModalFixture from '../../fixtures/components/ListView/NewTaskModalFixture'

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
  context('when task modal close', () => {
    let state = NewTaskModalFixture.initState()
    it('should not render new task modal', () => {
      const { output } = setup(state)
      expect(output.props.isTaskModalOpen).toBe(false)
    })
  })
  context('when task modal open', () => {
    let state = NewTaskModalFixture.openNewTaskModalState()
    it('should render new task modal', () => {
      const { output } = setup(state)
      expect(output.props.isTaskModalOpen).toBe(true)
    })
  })
})
