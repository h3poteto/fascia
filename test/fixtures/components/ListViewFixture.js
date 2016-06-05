import expect from 'expect'

export function initState() {
  return {
    ListReducer: {
      isListModalOpen: false,
      isTaskModalOpen: false,
      isListEditModalOpen: false,
      isProjectEditModalOpen: false,
      isLoading: false,
      newList: {title: "", color: "0effff"},
      newTask: {title: ""},
      lists: [
        {
          ID: 1,
          Title: "list1",
          IsHidden: false,
          ListTasks: [
            {
              ID: 1,
              Title: "task1",
              PullRequest: true
            }, {
              ID: 2,
              Title: "task2",
              PullRequest: false
            }
          ]
        }, {
          ID: 2,
          Title: "list2",
          IsHidden: false,
          ListTasks: []
        }
      ],
      noneList: {
        ID: 3,
        ListTasks: [
          {
            ID: 3,
            Title: "task3",
            PullRequest: false
          }
        ]
      },
      listOptions: [
        {
          ID: 1,
          Action: "close"
        }, {
          ID: 2,
          Action: "open"
        }
      ],
      selectedListOption: null,
      selectedList: null,
      project: {
        Title: "testProject",
        Description: "description",
        RepositoryID: 0,
        ShowIssues: true,
        ShowPullRequests: true
      },
      selectedProject: {Title: "", Description: "", RepositoryID: 0, ShowIssues: true, ShowPullRequests: true},
      isTaskDraggingOver: false,
      taskDraggingFrom: null,
      taskDraggingTo: null,
      error: null
    },
    params: {
      projectID: 1
    },
    fetchLists: expect.createSpy(),
    fetchProject: expect.createSpy(),
    fetchListOptions: expect.createSpy(),
    fetchUpdateList: expect.createSpy(),
    closeFlash: expect.createSpy(),
    taskDrop: expect.createSpy(),
    openNewListModal: expect.createSpy(),
    taskDragStart: expect.createSpy(),
    openNewTaskModal: expect.createSpy()
  }
}


export function errorState() {
  let state = initState()
  state["ListReducer"]["error"] = "Server Error"
  return state
}

export function wholeLoadingState() {
  let state = initState()
  state["ListReducer"]["isLoading"] = true
  return state
}

export function listModalState() {
  let state = initState()
  state["ListReducer"]["isListModalOpen"] = true
  return state
}

export function taskModalState() {
  let state = initState()
  state["ListReducer"]["isTaskModalOpen"] = true
  return state
}

export function listEditModalState() {
  let state = initState()
  state["ListReducer"]["isListEditModalOpen"] = true
  state["ListReducer"]["selectedListOption"] = {
    ID: 1,
    Action: "close"
  }
  return state
}

export function noRepositoryProjectState(defaultState) {
  let state
  if (defaultState == null) {
    state = initState()
  } else {
    state = defaultState
  }
  state["ListReducer"]["project"] = {
    Title: "testProject",
    Description: "description",
    RepositoryID: 0,
    ShowIssues: true,
    ShowPullRequests: true
  }
  return state
}

export function repositoryProjectState(defaultState) {
  let state
  if (defaultState == null) {
    state = initState()
  } else {
    state = defaultState
  }
  state["ListReducer"]["project"] = {
    Title: "testProject",
    Description: "description",
    RepositoryID: 1,
    ShowIssues: true,
    ShowPullRequests: true
  }
  return state
}

export function hideIssueState() {
  let state = initState()
  state["ListReducer"]["project"] = {
    Title: "testProject",
    Description: "description",
    RepositoryID: 1,
    ShowIssues: false,
    ShowPullRequests: true
  }
  return state
}

export function showIssueState() {
  let state = initState()
  state["ListReducer"]["project"] = {
    Title: "testProject",
    Description: "description",
    RepositoryID: 1,
    ShowIssues: true,
    ShowPullRequests: false
  }
  return state
}

export function hiddenListState() {
  let state =initState()
  state["ListReducer"]["lists"][0]["IsHidden"] = true
  return state
}
