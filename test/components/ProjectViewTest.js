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
      expect(props.projectActions.openNewProjectModal.calls.length).toBe(1)
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
})
