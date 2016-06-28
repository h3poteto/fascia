import expect from 'expect'

export function initState() {
  return {
    isListEditModalOpen: false,
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
}

export function openEditListModalState() {
  let state = initState()
  state["isListEditModalOpen"] = true
  return state
}

export function noRepositoryState(state) {
  state["project"] = {
    Title: "testProject",
    Description: "description",
    RepositoryID: 0,
    ShowIssues: true,
    ShowPullRequests: true
  }
  return state
}

export function hasRepositoryState(state) {
  state["project"] = {
    Title: "testProject",
    Description: "description",
    RepositoryID: 1,
    ShowIssues: true,
    ShowPullRequests: true
  }
  return state
}
