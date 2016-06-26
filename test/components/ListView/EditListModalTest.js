import TestUtils from 'react-addons-test-utils'
import expect from 'expect'
import React from 'react'
import EditListModal from '../../../frontend/javascripts/components/ListView/EditListModal.jsx'

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
      let state = {
        isListEditModalOpen: true,
        selectedListOption: {
          ID: 1,
          Action: "close"
        },
        selectedList: 1,
        project: {
          Title: "testProject",
          Description: "description",
          RepositoryID: 0,
          ShowIssues: true,
          ShowPullRequests: true
        },
        projectID: 1,
        listOptions: [
          {
            ID: 1,
            Action: "close"
          }, {
            ID: 2,
            Action: "open"
          }
        ],
        closeEditListModal: expect.createSpy(),
        updateSelectedListTitle: expect.createSpy(),
        updateSelectedListColor: expect.createSpy(),
        changeSelectedListOption: expect.createSpy(),
        fetchUpdateList: expect.createSpy()
      }
      it('should render list edit modal without action', () => {
        const { output, props } = setup(state)
        expect(output.props.isOpen).toBe(true)

        let listForm = output.props.children
        let form = listForm.props.children
        let fieldset = form.props.children
        let [ legend, titleLabel, titleInput, colorLabel, colorInput, nil, formAction ] = fieldset.props.children
        formAction.props.children.props.onClick()
        expect(props.fetchUpdateList.calls.length).toBe(1)
      })
    })
    context('when project has repository', () => {
      let state = {
        isListEditModalOpen: true,
        selectedListOption: {
          ID: 1,
          Action: "close"
        },
        selectedList: 1,
        project: {
          Title: "testProject",
          Description: "description",
          RepositoryID: 1,
          ShowIssues: true,
          ShowPullRequests: true
        },
        projectID: 1,
        listOptions: [
          {
            ID: 1,
            Action: "close"
          }, {
            ID: 2,
            Action: "open"
          }
        ],
        closeEditListModal: expect.createSpy(),
        updateSelectedListTitle: expect.createSpy(),
        updateSelectedListColor: expect.createSpy(),
        changeSelectedListOption: expect.createSpy(),
        fetchUpdateList: expect.createSpy()
      }
      it('should render list edit modal with action', () => {
        const { output, props } = setup(state)
        expect(output.props.isOpen).toBe(true)

        let listForm = output.props.children
        let form = listForm.props.children
        let fieldset = form.props.children
        let [ legend, titleLabel, titleInput, colorLabel, colorInput, actionWrapper, formAction ] = fieldset.props.children
        let [ actionLabel, actionSelect ] = actionWrapper.props.children
        expect(actionSelect.props.value).toBe(state.selectedListOption.ID)
      })
    })
  })
})
