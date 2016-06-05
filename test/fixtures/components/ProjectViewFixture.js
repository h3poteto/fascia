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
    fetchProjects: expect.createSpy(),
    fetchRepositories: expect.createSpy(),
    fetchSession: expect.createSpy(),
    closeFlash: expect.createSpy(),
    openNewProjectModal: expect.createSpy(),
    updateNewProjectTitle: expect.createSpy(),
    changeSelectedRepository: expect.createSpy(),
    fetchCreateProject: expect.createSpy()
  }
}

export function errorState() {
  let state = initState()
  state["ProjectReducer"]["error"] = "Server Error"
  return state
}

export function openProjectModal() {
  let state = initState()
  state["ProjectReducer"]["isModalOpen"] = true
  state["ProjectReducer"]["newProject"] = {
    title: "Title",
    description: ""
  }
  return state
}

export function openProjectModalWithRepository() {
  let state = openProjectModal()
  state["ProjectReducer"]["selectedRepository"] = {
    id : 1
  }
  return state
}

