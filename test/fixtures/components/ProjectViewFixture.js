import expect from 'expect'

export function initState() {
  return {
    ProjectReducer: {
      isModalOpen: false,
      newProject: {
        title: "",
        description: ""
      },
      projects: [{
        ID: 1,
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
    projectActions: {
      fetchProjects: expect.createSpy(),
      fetchRepositories: expect.createSpy(),
      fetchSession: expect.createSpy(),
      closeFlash: expect.createSpy(),
      openNewProjectModal: expect.createSpy()
    },
    newProjectModalActions: {
      closeNewProjectModal: expect.createSpy(),
      updateNewProjectTitle: expect.createSpy(),
      updateNewProjectDescription: expect.createSpy(),
      changeSelectedRepository: expect.createSpy(),
      fetchCreateProject: expect.createSpy()
    }
  }
}

export function errorState() {
  let state = initState()
  state["ProjectReducer"]["error"] = "Server Error"
  return state
}

