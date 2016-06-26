import TestUtils from 'react-addons-test-utils'
import expect from 'expect'
import React from 'react'
import NewListModal from '../../../frontend/javascripts/components/ListView/NewListModal.jsx'

function setup(props) {
  let renderer = TestUtils.createRenderer()
  renderer.render(<NewListModal {...props} />)
  let output = renderer.getRenderOutput()

  return {
    props,
    output,
    renderer
  }
}

describe('ListView::NewListModal', () => {
  context('when new list modal open', () => {
    let state = {
      isListModalOpen: true,
      newList: {title: "", color: "0effff"},
      projectID: 1,
      closeNewListModal: expect.createSpy(),
      updateNewListTitle: expect.createSpy(),
      updateNewListColor: expect.createSpy(),
      fetchCreateList: expect.createSpy
    }
    it('should render modal', () => {
      const { output } = setup(state)
      expect(output.props.isOpen).toBe(true)
    })
  })
})
