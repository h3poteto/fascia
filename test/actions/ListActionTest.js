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


describe('openNewTaskModal', () => {
  it('should open new task modal', () => {
    const list = {
      ID: 1,
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

describe('openEditListModal', () => {
  it('should open edit list modal', () => {
    const list = {
      ID: 1,
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


// async tests
describe('fetchLists', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    const projectID = 1
    beforeEach(() => {
      nock('http://localhost')
        .get(`/projects/${projectID}/lists`)
        .reply(200, { Lists: ['list1', 'list2'], NoneList: "noneList" })
    })

    it('call RECEIVE_LISTS and get lists', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_LISTS },
        { type: listActions.RECEIVE_LISTS, lists: ['list1', 'list2'], noneList: "noneList" }
      ]
      const store = mockStore({ lists: [] }, expectedActions, done)
      store.dispatch(listActions.fetchLists(projectID))
    })
  })
})

describe('fetchProject', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    const projectID = 1
    beforeEach(() => {
      nock('http://localhost')
        .get(`/projects/${projectID}/show`)
        .reply(200, { ID: 1, Title: "projectTitle" } )
    })

    it('call RECEIVE_PROJECT and get project', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_PROJECT },
        { type: listActions.RECEIVE_PROJECT, project: { ID: 1, Title: "projectTitle" } }
      ]
      const store = mockStore({ project: [] }, expectedActions, done)
      store.dispatch(listActions.fetchProject(projectID))
    })
  })
})


describe('fetchPorjectGithub', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    const projectID = 1
    beforeEach(() => {
      nock('http://localhost')
        .post(`/projects/${projectID}/fetch_github`)
        .reply(200, {
          Lists: ["list1", "list2"],
          NoneList: "noneList"
        })
    })
    it('call RECEIVE_FETCH_GITHUB and get lists', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_FETCH_GITHUB },
        { type: listActions.RECEIVE_FETCH_GITHUB, lists: ["list1", "list2"], noneList: "noneList" }
      ]
      const store = mockStore({ lists: [] }, expectedActions, done)
      store.dispatch(listActions.fetchProjectGithub(projectID))
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
    const event = {
      target: {
        className: "task"
      }
    }
    expect(listActions.taskDragLeave(event)).toEqual(expectedAction)
  })
})

describe('taskDrop', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when add task to list at the last', () => {
    const projectID = 1
    const taskDraggingFrom = {
      fromList: {
        ID: 1
      },
      fromTask: {
        ID: 5
      }
    }
    const taskDraggingTo = {
      toList: {
        ID: 2
      },
      prevToTask: null
    }
    beforeEach(() => {
      nock('http://localhost')
        .post(`/projects/${projectID}/lists/${taskDraggingFrom.fromList.ID}/tasks/${taskDraggingFrom.fromTask.ID}/move_task`)
        .reply(200, {
          Lists: ["list1", "list2"],
          NoneList: "noneList"
        })
    })

    it('call RECEIVE_MOVE_TASK and get lists', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_MOVE_TASK },
        { type: listActions.RECEIVE_MOVE_TASK, lists: ["list1", "list2"], noneList: "noneList" }
      ]
      const store = mockStore({ lists: [] }, expectedActions, done)
      store.dispatch(listActions.taskDrop(projectID, taskDraggingFrom, taskDraggingTo))
    })
  })
  context('when add task to list at halfway', () => {
    const projectID = 1
    const taskDraggingFrom = {
      fromList: {
        ID: 1
      },
      fromTask: {
        ID: 5
      }
    }
    const taskDraggingTo = {
      toList: {
        ID: 2
      },
      prevToTask: {
        ID: 6
      }
    }
    beforeEach(() => {
      nock('http://localhost')
        .post(`/projects/${projectID}/lists/${taskDraggingFrom.fromList.ID}/tasks/${taskDraggingFrom.fromTask.ID}/move_task`)
        .reply(200, {
          Lists: ["list1", "list2"],
          NoneList: "noneList"
        })
    })

    it('call RECEIVE_MOVE_TASK and get lists', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_MOVE_TASK },
        { type: listActions.RECEIVE_MOVE_TASK, lists: ["list1", "list2"], noneList: "noneList" }
      ]
      const store = mockStore({ lists: [] }, expectedActions, done)
      store.dispatch(listActions.taskDrop(projectID, taskDraggingFrom, taskDraggingTo))
    })
  })
  context('when dragg target is undefined', () => {
    it('call TASK_DROP and do nothing', () => {
      const projectID = 1
      const taskDraggingFrom = {
        fromList: {
          ID: 1
        },
        fromTask: {
          ID: 5
        }
      }
      const taskDraggingTo = null
      const expectedAction = {
        type: listActions.TASK_DROP
      }
      expect(listActions.taskDrop(projectID, taskDraggingFrom, taskDraggingTo)).toEqual(expectedAction)
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
            ID: 1
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
              ID: 1
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
                ID: 2
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
