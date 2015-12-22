import TestUtils from 'react-addons-test-utils'
import expect from 'expect'
import React from 'react'
import ProjectView from '../../frontend/javascripts/components/ProjectView.jsx'

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
    let state = {
      ProjectReducer: {
        isModalOpen: false,
        newProject: {
          title: "",
          description: ""
        },
        projects: [{
          Id: 1,
          Title: "project title",
          Description: "project description"
        }]
      },
      fetchProjects: expect.createSpy(),
      fetchRepositories: expect.createSpy(),
      closeFlash: expect.createSpy(),
      openNewProjectModal: expect.createSpy()
    }
    it('should render correctly', () => {
      const { output, props } = setup(state)

      expect(output.type).toBe('div')
      expect(output.props.id).toBe('projects')

      let [ flash, modal, items ] = output.props.children
      expect(items.props.className).toBe('items')

      let [ link, button ] = items.props.children
      expect(link[0].props.to).toBe(`/projects/${state.ProjectReducer.projects[0].Id}`)
      expect(button.type).toBe('button')
      button.props.onClick()
      expect(props.openNewProjectModal.calls.length).toBe(1)
    })
  })
  // TODO: write
  context('when one error, no repositories, and one project', () => {
  })
  context('when open project modal with repositories', () => {
  })
})
