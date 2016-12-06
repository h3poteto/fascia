import TestUtils from 'react-addons-test-utils'
import expect from 'expect'
import React from 'react'
import NewProjectModal from '../../../frontend/javascripts/components/ProjectView/NewProjectModal.jsx'
import * as NewProjectModalFixture from '../../fixtures/components/ProjectView/NewProjectModalFixture'

function setup(props) {
  let renderer = TestUtils.createRenderer()
  renderer.render(<NewProjectModal {...props} />)
  let output = renderer.getRenderOutput()

  return {
    props,
    output,
    renderer
  }
}

describe('ProjectView::NewProjectModal', () => {
  context('when close new project modal', () => {
    let state = NewProjectModalFixture.initState()
    it('should not render modal', () => {
      const { output } = setup(state)

      expect(output.props.isModalOpen).toBe(false)
    })
  })
})
