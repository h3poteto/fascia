import TestUtils from 'react-addons-test-utils'
import expect from 'expect'
import React from 'react'
import EditProjectModal from '../../../frontend/javascripts/components/ListView/EditProjectModal.jsx'
import * as EditProjectModalFixture from '../../fixtures/components/ListView/EditProjectModalFixture'

function setup(props) {
  let renderer = TestUtils.createRenderer()
  renderer.render(<EditProjectModal {...props} />)
  let output = renderer.getRenderOutput()

  return {
    props,
    output,
    renderer
  }
}

describe('ListView::EditProjectModal', () => {
  context('when project edit modal close', () => {
    let state = EditProjectModalFixture.initState()
    it('should not render list edit modal', () => {
      const { output } = setup(state)
      expect(output.props.isProjectEditModalOpen).toBe(false)
    })
  })
  context('when project edit modal open', () => {
    let state = EditProjectModalFixture.openEditProjectModalState()
    it('should render modal', () => {
      const { output } = setup(state)
      expect(output.props.isProjectEditModalOpen).toBe(true)
    })
  })
})
