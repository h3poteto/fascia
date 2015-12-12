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

    it('create RECEIVE_LISTS', (done) => {
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

    it('create RECEIVE_PROJECT', (done) => {
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
    const postForm = "title=listTitle&color=ffffff"
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

    it('create RECEIVE_CREATE_LIST', (done) => {
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
    const postForm = `title=${title}`
    beforeEach(() => {
      nock('http://localhost')
        .post(`/projects/${projectId}/lists/${listId}/tasks`, postForm)
        .reply(200, {
          Id: 1,
          ListId: listId,
          Title: title
        })
    })

    it('create RECEIVE_CREATE_TASK', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_CREATE_TASK },
        { type: listActions.RECEIVE_CREATE_TASK, task: { Id: 1, ListId: listId, Title: title } }
      ]
      const store = mockStore({ task: null }, expectedActions, done)
      store.dispatch(listActions.fetchCreateTask(projectId, listId, title))
    })
  })
})

describe('fetchUpdateList', () => {
  afterEach(() => {
    nock.cleanAll()
  })
  context('when response is right', () => {
    const projectId = 1
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
    const postForm = `title=${list.Title.String}&color=${list.Color.String}`
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
    it('create RECEIVE_UPDATE_LIST', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_UPDATE_LIST },
        { type: listActions.RECEIVE_UPDATE_LIST, list: { Id: list.Id, ProjectId: list.ProjectId, Title: list.Title.String, Color: list.Color.String, ListTasks: [] } }
      ]
      const store = mockStore({ list: null }, expectedActions, done)
      store.dispatch(listActions.fetchUpdateList(projectId, list))
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
    it('create RECEIVE_FETCH_GITHUB', (done) => {
      const expectedActions = [
        { type: listActions.REQUEST_FETCH_GITHUB },
        { type: listActions.RECEIVE_FETCH_GITHUB, lists: { lists: ["list1", "list2"] } }
      ]
      const store = mockStore({ lists: [] }, expectedActions, done)
      store.dispatch(listActions.fetchProjectGithub(projectId))
    })
  })
})

// TODO: drag関連のテストは真面目に考えて書くこと
