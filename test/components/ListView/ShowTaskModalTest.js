import TestUtils from 'react-addons-test-utils'
import expect from 'expect'
import React from 'react'
import ShowTaskModal from '../../../frontend/javascripts/components/ListView/ShowTaskModal.jsx'
import * as ShowTaskModalFixture from '../../fixtures/components/ListView/ShowTaskModalFixture'

function setup(props) {
  let renderer = TestUtils.createRenderer()
  renderer.render(<ShowTaskModal {...props} />)
  let output = renderer.getRenderOutput()

  return {
    props,
    output,
    renderer
  }
}

describe('ListView::ShowTaskModal', () => {
  context('when task modal close', () => {
    let state = ShowTaskModalFixture.initState()
    it('should not render show task modal', () => {
      const { output } = setup(state)
      expect(output.props.isOpen).toBe(false)
    })
  })
  context('when task modal open', () => {
    let state = ShowTaskModalFixture.openShowTaskModalState()
    it('should render show task modal', () => {
      const { output } = setup(state)
      expect(output.props.isOpen).toBe(true)
    })
  })

  context('when edit task modal visible', () => {
    let state = ShowTaskModalFixture.visibleEditTaskModalState()
    it('should render edit task form', () => {
      const { output } = setup(state)
      expect(output.props.isOpen).toBe(true)

      let taskDetail = output.props.children
      let [ , body ]  = taskDetail.props.children
      expect(body.props.children.type).toBe('form')
    })
  })
})

