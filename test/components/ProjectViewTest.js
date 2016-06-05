import TestUtils from 'react-addons-test-utils'
import expect from 'expect'
import React from 'react'
import ProjectView from '../../frontend/javascripts/components/ProjectView.jsx'
import * as ProjectViewFixture from '../fixtures/components/ProjectViewFixture'

function setup(props) {
  let renderer = TestUtils.createRenderer()
  renderer.render(<ProjectView {...props} />)
  let output = renderer.getRenderOutput()

  return {
    props,
    output,
    renderer
  }
}

describe('ProjectView', () => {
  context('when no error, no repositories, and one project', () => {
    let state = ProjectViewFixture.initState()
    it('should render correctly', () => {
      const { output, props } = setup(state)

      expect(output.type).toBe('div')
      expect(output.props.id).toBe('projects')

      let [ wholeLoading, flash, modal, items ] = output.props.children
      expect(items.props.className).toBe('items')

      let [ link, button ] = items.props.children
      expect(link[0].props.to).toBe(`/projects/${state.ProjectReducer.projects[0].ID}`)
      expect(button.type).toBe('button')
      button.props.onClick()
      expect(props.openNewProjectModal.calls.length).toBe(1)
    })
  })
  context('when one error, no repositories, and one project', () => {
    let state = ProjectViewFixture.errorState()
    it('should render error', () => {
      const { output } = setup(state)
      let [ wholeLoading, flash, modal, items ] = output.props.children
      expect(flash.props.children).toBe('Server Error')
    })
  })

  context('when open project modal with repositories', () => {
    context('without selectedRepository', () => {
      let state = ProjectViewFixture.openProjectModal()
      it('should render repositories', () => {
        const { output, props } = setup(state)

        expect(output.type).toBe('div')
        expect(output.props.id).toBe('projects')

        let [ wholeLoading, flash, modal, items ] = output.props.children
        expect(modal.props.isOpen).toBe(true)
        let formWrapper = modal.props.children
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
      let state = ProjectViewFixture.openProjectModalWithRepository()
      it('should render repositories', () => {
        const { output, props } = setup(state)

        expect(output.type).toBe('div')
        expect(output.props.id).toBe('projects')

        let [ wholeLoading, flash, modal, items ] = output.props.children
        expect(modal.props.isOpen).toBe(true)
        let formWrapper = modal.props.children
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
