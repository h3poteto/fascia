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
  context('when list edit modal open', () => {
    context('when project does not have repository', () => {
      let state = EditListModalFixture.noRepositoryState(EditListModalFixture.openEditListModalState())
      it('should render list edit modal without action', () => {
        const { output, props } = setup(state)
        expect(output.props.isOpen).toBe(true)

        let listForm = output.props.children
        let form = listForm.props.children
        let fieldset = form.props.children
        let [ , , , , , , formAction ] = fieldset.props.children
        formAction.props.children.props.onClick()
        expect(props.fetchUpdateList.calls.length).toBe(1)
      })
    })
    context('when project has repository', () => {
      let state = EditListModalFixture.hasRepositoryState(EditListModalFixture.openEditListModalState())
      it('should render list edit modal with action', () => {
        const { output } = setup(state)
        expect(output.props.isOpen).toBe(true)

        let listForm = output.props.children
        let form = listForm.props.children
        let fieldset = form.props.children
        let [ , , , , , actionWrapper ] = fieldset.props.children
        let [ , actionSelect ] = actionWrapper.props.children
        expect(actionSelect.props.value).toBe(state.selectedListOption.ID)
      })
    })
  })
})
