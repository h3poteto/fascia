import expect from 'expect'

export function initState() {
  return {
    isListEditModalOpen: false,
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
