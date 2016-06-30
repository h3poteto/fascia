import expect from 'expect'

export function initState() {
  return {
    isProjectEditModalOpen: false,
    project: null,
    selectedProject: {Title: "", Description: "", RepositoryID: 0, ShowIssues: true, ShowPullRequests: true},
    closeEditProjectModal: expect.createSpy(),
    updateEditProjectTitle: expect.createSpy(),
    updateEditProjectDescription: expect.createSpy(),
    fetchUpdateProject: expect.createSpy()
  }
}

export function openEditProjectModalState() {
  let state = initState()
  state["isProjectEditModalOpen"] = true
  return state
}

