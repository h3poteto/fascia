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
        }],
        isLoading: false
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

      let [ wholeLoading, flash, modal, items ] = output.props.children
      expect(items.props.className).toBe('items')

      let [ link, button ] = items.props.children
      expect(link[0].props.to).toBe(`/projects/${state.ProjectReducer.projects[0].Id}`)
      expect(button.type).toBe('button')
      button.props.onClick()
      expect(props.openNewProjectModal.calls.length).toBe(1)
    })
  })
  context('when one error, no repositories, and one project', () => {
    let state = {
      ProjectReducer: {
        isModalOpen: true,
        newProject: {
          title: "Title",
          description: ""
        },
        projects: [{
          Id: 1,
          Title: "project title",
          Description: "project description"
        }],
        isLoading: false
      },
      fetchProjects: expect.createSpy(),
      fetchRepositories: expect.createSpy(),
      closeFlash: expect.createSpy(),
      openNewProjectModal: expect.createSpy(),
      updateNewProjectTitle: expect.createSpy(),
      changeSelectedRepository: expect.createSpy(),
      fetchCreateProject: expect.createSpy()
    }
    it('should render page and modal', () => {
      const { output, props } = setup(state)

      expect(output.type).toBe('div')
      expect(output.props.id).toBe('projects')

      let [ wholeLoading, flash, modal, items ] = output.props.children
      expect(modal.props.isOpen).toBe(true)
      let formWrapper = modal.props.children
      let form = formWrapper.props.children
      let field = form.props.children
      let [ legend, titileLabel, titleInput, descriptionLabel, descriptionInput, repositoryLabel, repositorySelect, action ] = field.props.children
      expect(titleInput.props.value).toBe("Title")
      titleInput.props.onChange()
      expect(props.updateNewProjectTitle.calls.length).toBe(1)
      let button = action.props.children
      button.props.onClick()
      expect(props.fetchCreateProject.calls.length).toBe(1)
    })
  })
  context('when open project modal with repositories', () => {
    context('without selectedRepository', () => {
      let state = {
        ProjectReducer: {
          isModalOpen: true,
          newProject: {
            title: "Title",
            description: ""
          },
          projects: [{
            Id: 1,
            Title: "project title",
            Description: "project description"
          }],
          selectedRepository: null,
          repositories: [{
            id: 1,
            full_name: "repo1"
          }, {
            id: 2,
            full_name: "repo2"
          }],
          isLoading: false
        },
        fetchProjects: expect.createSpy(),
        fetchRepositories: expect.createSpy(),
        closeFlash: expect.createSpy(),
        openNewProjectModal: expect.createSpy(),
        updateNewProjectTitle: expect.createSpy(),
        changeSelectedRepository: expect.createSpy(),
        fetchCreateProject: expect.createSpy()
      }
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
      let state = {
        ProjectReducer: {
          isModalOpen: true,
          newProject: {
            title: "Title",
            description: ""
          },
          projects: [{
            Id: 1,
            Title: "project title",
            Description: "project description"
          }],
          selectedRepository: {
            id: 1
          },
          repositories: [{
            id: 1,
            full_name: "repo1"
          }, {
            id: 2,
            full_name: "repo2"
          }],
          isLoading: false
        },
        fetchProjects: expect.createSpy(),
        fetchRepositories: expect.createSpy(),
        closeFlash: expect.createSpy(),
        openNewProjectModal: expect.createSpy(),
        updateNewProjectTitle: expect.createSpy(),
        changeSelectedRepository: expect.createSpy(),
        fetchCreateProject: expect.createSpy()
      }
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
