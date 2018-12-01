import expect from 'expect'

export function initState() {
  return {
    ListReducer: {
      isListModalOpen: false,
      isTaskModalOpen: false,
      isListEditModalOpen: false,
      isProjectEditModalOpen: false,
      isTaskShowModalOpen: false,
      isDeleteProjectModalOpen: false,
      isEditTaskModalVisible: false,
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
      selectedList: {},
      project: {
        Title: "testProject",
        Description: "description",
        RepositoryID: 0,
        ShowIssues: true,
        ShowPullRequests: true
      },
      selectedProject: {Title: "", Description: "", RepositoryID: 0, ShowIssues: true, ShowPullRequests: true},
      selectedTask: {Title: "", Description: "", IssueNumber: 0},
      isTaskDraggingOver: false,
      taskDraggingFrom: null,
      taskDraggingTo: null,
      error: null
    },
    match: {
      params: {
        projectID: '1'
      }
    },
    listActions: {
      fetchLists: expect.createSpy(),
      fetchProject: expect.createSpy(),
      fetchListOptions: expect.createSpy(),
      closeFlash: expect.createSpy(),
      showPullRequests: expect.createSpy(),
      showIssues: expect.createSpy(),
      fetchProjectGithub: expect.createSpy(),
      openEditListModal: expect.createSpy(),
      taskDragOver: expect.createSpy(),
      taskDrop: expect.createSpy(),
      taskDragLeave: expect.createSpy(),
      displayList: expect.createSpy(),
      hideList: expect.createSpy(),
      taskDragStart: expect.createSpy(),
      openShowTaskModal: expect.createSpy(),
      openNewTaskModal: expect.createSpy(),
      openEditProjectModal: expect.createSpy(),
      openNewListModal: expect.createSpy(),
      openDeleteProjectModal: expect.createSpy(),
    },
    newListModalActions: {
      closeNewListModal: expect.createSpy(),
      fetchCreateList: expect.createSpy()
    },
    newTaskModalActions: {
      closeNewTaskModal: expect.createSpy(),
      fetchCreateTask: expect.createSpy()
    },
    editListModalActions: {
      closeEditListModal: expect.createSpy(),
      fetchUpdateList: expect.createSpy()
    },
    editProjectModalActions: {
      closeEditProjectModal: expect.createSpy(),
      fetchUpdateProject: expect.createSpy(),
      createWebhook: expect.createSpy()
    },
    showTaskModalActions: {
      closeShowTaskModal: expect.createSpy(),
      changeEditMode: expect.createSpy(),
      fetchUpdateTask: expect.createSpy(),
      fetchDeleteTask: expect.createSpy()
    },
    deleteProjectModalActions: {
      closeDeleteProjectModal: expect.createSpy(),
      fetchDeleteProject: expect.createSpy()
    },
  }
}


export function errorState() {
  let state = initState()
  state["ListReducer"]["error"] = "Server Error"
  return state
}


export function noRepositoryProjectState() {
  let state = initState()
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
  let state = initState()
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
