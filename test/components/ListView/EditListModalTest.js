import TestUtils from 'react-addons-test-utils'
import expect from 'expect'
import React from 'react'
import EditListModal from '../../../frontend/javascripts/components/ListView/EditListModal.jsx'
import * as EditListModalFixture from '../../fixtures/components/ListView/EditListModalFixture'

function setup(props) {
  let renderer = TestUtils.createRenderer()
  renderer.render(<EditListModal {...props} />)
  let output = renderer.getRenderOutput()

  return {
    props,
    output,
    renderer
  }
}

describe('ListView::EditListModal', () => {
  context('when list edit modal close', () => {
    let state = EditListModalFixture.initState()
    it('should not render list edit modal', () => {
      const { output } = setup(state)
      expect(output.props.isListEditModalOpen).toBe(false)
    })
  })

  context('when list edit modal open', () => {
    context('when project does not have repository', () => {
      let state = EditListModalFixture.noRepositoryState(EditListModalFixture.openEditListModalState())
      it('should render list edit modal without action', () => {
        const { output, props } = setup(state)
        expect(output.props.isListEditModalOpen).toBe(true)
      })
    })
  })
})
