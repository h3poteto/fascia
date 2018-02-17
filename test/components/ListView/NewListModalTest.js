import ShallowRenderer from 'react-test-renderer/shallow'
import expect from 'expect'
import React from 'react'
import NewListModal from '../../../frontend/javascripts/components/ListView/NewListModal.jsx'
import * as NewListModalFixture from '../../fixtures/components/ListView/NewListModalFixture'

function setup(props) {
  let renderer = new ShallowRenderer()
  renderer.render(<NewListModal {...props} />)
  let output = renderer.getRenderOutput()

  return {
    props,
    output,
    renderer
  }
}

describe('ListView::NewListModal', () => {
  context('when new list modal close', () => {
    let state = NewListModalFixture.initState()
    it('should not render modal', () => {
      const { output } = setup(state)
      expect(output.props.isListModalOpen).toBe(false)
    })
  })
  context('when new list modal open', () => {
    let state = NewListModalFixture.openNewListModalState()
    it('should render modal', () => {
      const { output } = setup(state)
      expect(output.props.isListModalOpen).toBe(true)
    })
  })
})
