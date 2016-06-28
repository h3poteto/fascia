import expect from 'expect'

export function initState() {
  return {
    isModalOpen: false,
    newProject: {
      title: "title",
      description: ""
    },
    repositories: [{
      id: 1,
      full_name: "repo1"
    }, {
      id: 2,
      full_name: "repo2"
    }],
    selectedRepository: null,
    closeNewProjectModal: expect.createSpy(),
    updateNewProjectTitle: expect.createSpy(),
    updateNewProjectDescription: expect.createSpy(),
    changeSelectedRepository: expect.createSpy(),
    fetchCreateProject: expect.createSpy()
  }
}

export function openNewProjectModalState() {
  var state = initState()
  state["isModalOpen"] = true
  state["newProject"] = {
    title: "Title",
    description: ""
  }
  return state
}

export function openNewProjectModalWithRepositoryState() {
  var state = openNewProjectModalState()
  state["selectedRepository"] = {
    id: 1,
    full_name: "repo1"
  }
  return state
}
