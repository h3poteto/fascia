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

      expect(output.props.isOpen).toBe(false)
    })
  })

  context('when open project modal with repositories', () => {
    context('without selectedRepository', () => {
      let state = NewProjectModalFixture.openNewProjectModalState()
      it('should render repositories', () => {
        const { output } = setup(state)

        expect(output.props.isOpen).toBe(true)
        let formWrapper = output.props.children
        let form = formWrapper.props.children
        let field = form.props.children
        let [ legend, titileLabel, titleInput, descriptionLabel, descriptionInput, repositoryLabel, repositorySelect, action ] = field.props.children
        let [ option1, repos ] = repositorySelect.props.children
        let [ repo1, repo2 ] = repos
        expect(repo1.props.children).toBe('repo1')
        expect(repo2.props.children).toBe('repo2')
      })
    })
    context('with selectedRepository', () => {
      let state = NewProjectModalFixture.openNewProjectModalWithRepositoryState()
      it('should render repositories', () => {
        const { output } = setup(state)

        expect(output.props.isOpen).toBe(true)
        let formWrapper = output.props.children
        let form = formWrapper.props.children
        let field = form.props.children
        let [ legend, titileLabel, titleInput, descriptionLabel, descriptionInput, repositoryLabel, repositorySelect, action ] = field.props.children
        let [ option1, repos ] = repositorySelect.props.children
        let [ repo1, repo2 ] = repos
        expect(repo1.props.children).toBe('repo1')
        expect(repo1.props.selected).toBe(true)
        expect(repo2.props.children).toBe('repo2')
      })
    })
  })
})
