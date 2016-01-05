import * as listActions from "../../frontend/javascripts/actions/ListAction"
import expect from 'expect'
import nock from 'nock'
import mockStore from '../support/MockStore'


// normal tests

describe('closeFlash', () => {
  it('should close flash', () => {
    const expectedAction = {
      type: listActions.CLOSE_FLASH
    }
    expect(listActions.closeFlash()).toEqual(expectedAction)
  })
})


describe('openNewListModal', () => {
  it('should open new list modal', () => {
    const expectedAction = {
      type: listActions.OPEN_NEW_LIST,
      isListModalOpen: true
    }
    expect(listActions.openNewListModal()).toEqual(expectedAction)
  })
})

describe('closeNewListModal', () => {
  it('should close new list modal', () => {
    const expectedAction = {
      type: listActions.CLOSE_NEW_LIST,
      isListModalOpen: false
    }
    expect(listActions.closeNewListModal()).toEqual(expectedAction)
  })
})

describe('openNewTaskModal', () => {
  it('should open new task modal', () => {
    const list = {
      Id: 1,
      Title: "listTitle"
    }
    const expectedAction = {
      type: listActions.OPEN_NEW_TASK,
      isTaskModalOpen: true,
      list: list
    }
    expect(listActions.openNewTaskModal(list)).toEqual(expectedAction)
  })
})

describe('closeNewTaskModal', () => {
  it('should close new task modal', () => {
    const expectedAction = {
      type: listActions.CLOSE_NEW_TASK,
      isTaskModalOpen: false
    }
    expect(listActions.closeNewTaskModal()).toEqual(expectedAction)
  })
})

describe('openEditListModal', () => {
  it('should open edit list modal', () => {
    const list = {
      Id: 1,
      Title: "listTitle"
    }
    const expectedAction = {
      type: listActions.OPEN_EDIT_LIST,
      isListEditModalOpen: true,
      list: list
    }
    expect(listActions.openEditListModal(list)).toEqual(expectedAction)
  })
})

describe('closeEditListModal', () => {
  it('should close edit list modal', () => {
    const expectedAction = {
      type: listActions.CLOSE_EDIT_LIST,
      isListEditModalOpen: false
    }
    expect(listActions.closeEditListModal()).toEqual(expectedAction)
  })
})

describe('updateNewListTitle', () => {
  it('should update new list title', () => {
    const title = "newTitle"
    const ev = {
      target: {
        value: title
      }
    }
    const expectedAction = {
      type: listActions.UPDATE_NEW_LIST_TITLE,
      title: title
    }
    expect(listActions.updateNewListTitle(ev)).toEqual(expectedAction)
  })
})

describe('updateNewListColor', () => {
  it('should update new list color', () => {
    const color = "ffffff"
    const ev = {
      target: {
        value: color
      }
    }
    const expectedAction = {
      type: listActions.UPDATE_NEW_LIST_COLOR,
      color: color
    }
    expect(listActions.updateNewListColor(ev)).toEqual(expectedAction)
  })
})


// async tests
describe('fetchLists', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    const projectId = 1
    beforeEach(() => {
      nock('http://localhost')
        .get(`/projects/${projectId}/lists`)
        .reply(200, { lists: ['list1', 'list2'] })
    })

    it('call RECEIVE_LISTS and get lists', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_LISTS },
        { type: listActions.RECEIVE_LISTS, lists: { lists: ['list1', 'list2'] } }
      ]
      const store = mockStore({ lists: [] }, expectedActions, done)
      store.dispatch(listActions.fetchLists(projectId))
    })
  })
})

describe('fetchProject', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    const projectId = 1
    beforeEach(() => {
      nock('http://localhost')
        .get(`/projects/${projectId}/show`)
        .reply(200, { Id: 1, Title: "projectTitle" } )
    })

    it('call RECEIVE_PROJECT and get project', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_PROJECT },
        { type: listActions.RECEIVE_PROJECT, project: { Id: 1, Title: "projectTitle" } }
      ]
      const store = mockStore({ project: [] }, expectedActions, done)
      store.dispatch(listActions.fetchProject(projectId))
    })
  })
})

describe('fetchCreateList', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    const projectId = 1
    const title = "listTitle"
    const color = "ffffff"
    const postForm = `title=${title}&color=${color}`
    beforeEach(() => {
      nock('http://localhost')
        .post(`/projects/${projectId}/lists`, postForm)
        .reply(200, {
          Id: 1,
          ProjectId: projectId,
          Title: title,
          Color: color,
          ListTasks: ["task1"]
        })
    })

    it('call RECEIVE_CREATE_LIST and get list', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_CREATE_LIST },
        { type: listActions.RECEIVE_CREATE_LIST, list: { Id: 1, ProjectId: projectId, Title: title, Color: color, ListTasks: ["task1"] } }
      ]
      const store = mockStore({ list: null }, expectedActions, done)
      store.dispatch(listActions.fetchCreateList(projectId, title, color))
    })
  })
})

describe('fetchCreateTask', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    const projectId = 1
    const listId = 2
    const title = "taskTitle"
    const description = "taskDescription"
    const postForm = `title=${title}&description=${description}`
    beforeEach(() => {
      nock('http://localhost')
        .post(`/projects/${projectId}/lists/${listId}/tasks`, postForm)
        .reply(200, {
          Id: 1,
          ListId: listId,
          Title: title,
          Description: description
        })
    })

    it('call RECEIVE_CREATE_TASK and get task', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_CREATE_TASK },
        { type: listActions.RECEIVE_CREATE_TASK, task: { Id: 1, ListId: listId, Title: title, Description: description } }
      ]
      const store = mockStore({ task: null }, expectedActions, done)
      store.dispatch(listActions.fetchCreateTask(projectId, listId, title, description))
    })
  })
})

describe('fetchUpdateList', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    const projectId = 1
    const option = {
      Id: 1,
      Action: "close"
    }
    const list = {
      Id: 2,
      Title: {
        String: "listTitle"
      },
      Color: {
        String: "ffffff"
      },
      ProjectId: projectId
    }
    const postForm = `title=${list.Title.String}&color=${list.Color.String}&action=${option.Action}`
    beforeEach(() => {
      nock('http://localhost')
        .post(`/projects/${projectId}/lists/${list.Id}`, postForm)
        .reply(200, {
          Id: list.Id,
          ProjectId: list.ProjectId,
          Title: list.Title.String,
          Color: list.Color.String,
          ListTasks: []
        })
    })
    it('call RECEIVE_UPDATE_LIST and get list', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_UPDATE_LIST },
        { type: listActions.RECEIVE_UPDATE_LIST, list: { Id: list.Id, ProjectId: list.ProjectId, Title: list.Title.String, Color: list.Color.String, ListTasks: [] } }
      ]
      const store = mockStore({ list: null }, expectedActions, done)
      store.dispatch(listActions.fetchUpdateList(projectId, list, option))
    })
  })
})

describe('fetchPorjectGithub', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    const projectId = 1
    beforeEach(() => {
      nock('http://localhost')
        .post(`/projects/${projectId}/fetch_github`)
        .reply(200, {
          lists: ["list1", "list2"]
        })
    })
    it('call RECEIVE_FETCH_GITHUB and get lists', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_FETCH_GITHUB },
        { type: listActions.RECEIVE_FETCH_GITHUB, lists: { lists: ["list1", "list2"] } }
      ]
      const store = mockStore({ lists: [] }, expectedActions, done)
      store.dispatch(listActions.fetchProjectGithub(projectId))
    })
  })
})


describe('fetchListOptions', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    beforeEach(() => {
      nock('http://localhost')
        .get('/list_options')
        .reply(200, {
          listOptions: ["option1", "option2"]
        })
    })
    it('call RECEIVE_LIST_OPTIONS and get listOptions', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_LIST_OPTIONS },
        { type: listActions.RECEIVE_LIST_OPTIONS, listOptions: { listOptions: ["option1", "option2"] } }
      ]
      const store = mockStore({ listOptions: [] }, expectedActions, done)
      store.dispatch(listActions.fetchListOptions())
    })
  })
})

// drag

describe('taskDragStart', () => {
  it('should set drag target', () => {
    const event = {
      dataTransfer: {
        effectAllowed: null,
        setData: (format, target) => {
          return format, target
        }
      },
      currentTarget: {
        parentNode: {
          parentNode: {
            name: "parent"
          }
        }
      }
    }
    const expectedAction = {
      type: listActions.TASK_DRAG_START,
      taskDragTarget: event.currentTarget,
      taskDragFromList: event.currentTarget.parentNode.parentNode
    }
    expect(listActions.taskDragStart(event)).toEqual(expectedAction)
  })
})

describe('taskDragLeave', () => {
  it('should free drag target', () => {
    const expectedAction = {
      type: listActions.TASK_DRAG_LEAVE,
    }
    expect(listActions.taskDragLeave()).toEqual(expectedAction)
  })
})

describe('taskDrop', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when add task to list at the last', () => {
    const projectId = 1
    const taskDraggingFrom = {
      fromList: {
        Id: 1
      },
      fromTask: {
        Id: 5
      }
    }
    const taskDraggingTo = {
      toList: {
        Id: 2
      },
      prevToTask: null
    }
    beforeEach(() => {
      nock('http://localhost')
        .post(`/projects/${projectId}/lists/${taskDraggingFrom.fromList.Id}/tasks/${taskDraggingFrom.fromTask.Id}/move_task`)
        .reply(200, {
          lists: ["list1", "list2"]
        })
    })

    it('call RECEIVE_MOVE_TASK and get lists', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_MOVE_TASK },
        { type: listActions.RECEIVE_MOVE_TASK, lists: { lists: ["list1", "list2"] } }
      ]
      const store = mockStore({ lists: [] }, expectedActions, done)
      store.dispatch(listActions.taskDrop(projectId, taskDraggingFrom, taskDraggingTo))
    })
  })
  context('when add task to list at halfway', () => {
    const projectId = 1
    const taskDraggingFrom = {
      fromList: {
        Id: 1
      },
      fromTask: {
        Id: 5
      }
    }
    const taskDraggingTo = {
      toList: {
        Id: 2
      },
      prevToTask: {
        Id: 6
      }
    }
    beforeEach(() => {
      nock('http://localhost')
        .post(`/projects/${projectId}/lists/${taskDraggingFrom.fromList.Id}/tasks/${taskDraggingFrom.fromTask.Id}/move_task`)
        .reply(200, {
          lists: ["list1", "list2"]
        })
    })

    it('call RECEIVE_MOVE_TASK and get lists', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_MOVE_TASK },
        { type: listActions.RECEIVE_MOVE_TASK, lists: { lists: ["list1", "list2"] } }
      ]
      const store = mockStore({ lists: [] }, expectedActions, done)
      store.dispatch(listActions.taskDrop(projectId, taskDraggingFrom, taskDraggingTo))
    })
  })
  context('when dragg target is undefined', () => {
    it('call TASK_DROP and do nothing', () => {
      const projectId = 1
      const taskDraggingFrom = {
        fromList: {
          Id: 1
        },
        fromTask: {
          Id: 5
        }
      }
      const taskDraggingTo = null
      const expectedAction = {
        type: listActions.TASK_DROP
      }
      expect(listActions.taskDrop(projectId, taskDraggingFrom, taskDraggingTo)).toEqual(expectedAction)
    })
  })
})

describe('taskDragOver', () => {
  context('when drag over list element', () => {
    it('should get target task and list', () => {
      const event = {
        preventDefault: () => {
          return true
        },
        target: {
          dataset: {
            droppedDepth: "0"
          }
        }
      }
      const expectedAction = {
        type: listActions.TASK_DRAG_OVER,
        taskDragToTask: event.target,
        taskDragToList: event.target
      }
      expect(listActions.taskDragOver(event)).toEqual(expectedAction)
    })
  })
  context('when drag over list title element', () => {
    it('should get target task and list', () => {
      const event = {
        preventDefault: () => {
          return true
        },
        target: {
          dataset: {
            droppedDepth: "1"
          },
          parentNode: {
            Id: 1
          }
        }
      }
      const expectedAction = {
        type: listActions.TASK_DRAG_OVER,
        taskDragToTask: event.target,
        taskDragToList: event.target.parentNode
      }
      expect(listActions.taskDragOver(event)).toEqual(expectedAction)
    })
  })
  context('when drag over li element', () => {
    it('should get target task and list', () => {
      const event = {
        preventDefault: () => {
          return true
        },
        target: {
          dataset: {
            droppedDepth: "2"
          },
          parentNode: {
            parentNode: {
              Id: 1
            }
          }
        }
      }
      const expectedAction = {
        type: listActions.TASK_DRAG_OVER,
        taskDragToTask: event.target,
        taskDragToList: event.target.parentNode.parentNode
      }
      expect(listActions.taskDragOver(event)).toEqual(expectedAction)
    })
  })
  context('when drag over icon element', () => {
    it('should get target task and list', () => {
      const event = {
        preventDefault: () => {
          return true
        },
        target: {
          dataset: {
            droppedDepth: "3"
          },
          parentNode: {
            parentNode: {
              parentNode: {
                Id: 2
              }
            }
          }
        }
      }
      const expectedAction = {
        type: listActions.TASK_DRAG_OVER,
        taskDragToTask: event.target,
        taskDragToList: event.target.parentNode.parentNode.parentNode
      }
      expect(listActions.taskDragOver(event)).toEqual(expectedAction)
    })
  })
})
