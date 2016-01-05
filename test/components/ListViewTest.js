import TestUtils from 'react-addons-test-utils'
import expect from 'expect'
import React from 'react'
import ListView from '../../frontend/javascripts/components/ListView.jsx'

function setup(props) {
  let renderer = TestUtils.createRenderer()
  renderer.render(<ListView {...props} />)
  let output = renderer.getRenderOutput()

  return {
    props,
    output,
    renderer
  }
}

describe('ListView', () => {
  context('when no error', () => {
    let state = {
      ListReducer: {
        isListModalOpen: false,
        isTaskModalOpen: false,
        isListEditModalOpen: false,
        newList: {title: "", color: "0effff"},
        newTask: {title: ""},
        lists: [
          {
            Id: 1,
            Title: "list1",
            ListTasks: [
              {
                Id: 1,
                Title: "task1"
              }, {
                Id: 2,
                Title: "task2"
              }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }
        ],
        listOptions: [],
        selectedListOption: null,
        selectedList: null,
        project: {
          Title: "testProject"
        },
        isTaskDraggingOver: false,
        taskDraggingFrom: null,
        taskDraggingTo: null,
        error: null
      },
      params: {
        projectId: 1
      },
      fetchLists: expect.createSpy(),
      fetchProject: expect.createSpy(),
      fetchListOptions: expect.createSpy(),
      closeFlash: expect.createSpy(),
      taskDrop: expect.createSpy(),
      openNewListModal: expect.createSpy(),
      taskDragStart: expect.createSpy(),
      openNewTaskModal: expect.createSpy()
    }
    it('should render correctly', () => {
      const { output, props } = setup(state)

      expect(output.type).toBe('div')
      expect(output.props.id).toBe('lists')

      let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectTitleWrapper, items ] = output.props.children
      let [ icon, projectTitle ] = projectTitleWrapper.props.children
      expect(projectTitle.props.children).toBe('testProject')

      expect(items.props.className).toBe('items')

      let [ list, button ] = items.props.children
      expect(list[0].props['data-id']).toBe(1)
      expect(list[1].props['data-id']).toBe(2)

      list[0].props.onDrop()
      expect(props.taskDrop.calls.length).toBe(1)

      button.props.onClick()
      expect(props.openNewListModal.calls.length).toBe(1)

      // list which have tasks
      let [ firstListMenu, firstList1Title, firstTasks ] = list[0].props.children
      let [ firstListTasks, firstNewTask ] = firstTasks.props.children
      let [ task1, task2 ] = firstListTasks
      expect(task1.props['data-id']).toBe(1)
      expect(task1.props.children).toBe('task1')
      expect(task2.props['data-id']).toBe(2)
      expect(task2.props.children).toBe('task2')

      task2.props.onDragStart()
      expect(props.taskDragStart.calls.length).toBe(1)

      // list which do not have tasks
      let [ secondListMenu, secondListTitle, secondTasks ] = list[1].props.children
      let [ secondListTasks, secondNewTask ] = secondTasks.props.children
      expect(secondNewTask.props.className).toBe('new-task')

      secondNewTask.props.onClick()
      expect(props.openNewTaskModal.calls.length).toBe(1)
    })
  })

  context('when one error, not modal', () => {
    let state = {
      ListReducer: {
        isListModalOpen: false,
        isTaskModalOpen: false,
        isListEditModalOpen: false,
        newList: {title: "", color: "0effff"},
        newTask: {title: ""},
        lists: [
          {
            Id: 1,
            Title: "list1",
            ListTasks: [
              {
                Id: 1,
                Title: "task1"
              }, {
                Id: 2,
                Title: "task2"
              }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }
        ],
        listOptions: [],
        selectedListOption: null,
        selectedList: null,
        project: {
          Title: "testProject"
        },
        isTaskDraggingOver: false,
        taskDraggingFrom: null,
        taskDraggingTo: null,
        error: "Server Error"
      },
      params: {
        projectId: 1
      },
      fetchLists: expect.createSpy(),
      fetchProject: expect.createSpy(),
      fetchListOptions: expect.createSpy(),
      closeFlash: expect.createSpy(),
      taskDrop: expect.createSpy(),
      openNewListModal: expect.createSpy(),
      taskDragStart: expect.createSpy(),
      openNewTaskModal: expect.createSpy()
    }
    it('should render error', () => {
      const { output } = setup(state)
      let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectTitleWrapper, items ] = output.props.children
      expect(flash.props.children).toBe('Server Error')
    })
  })
  context('when whole loading is open', () => {
    let state = {
      ListReducer: {
        isListModalOpen: false,
        isTaskModalOpen: false,
        isListEditModalOpen: false,
        isLoading: true,
        newList: {title: "", color: "0effff"},
        newTask: {title: ""},
        lists: [
          {
            Id: 1,
            Title: "list1",
            ListTasks: [
              {
                Id: 1,
                Title: "task1"
              }, {
                Id: 2,
                Title: "task2"
              }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }
        ],
        listOptions: [],
        selectedListOption: null,
        selectedList: null,
        project: {
          Title: "testProject"
        },
        isTaskDraggingOver: false,
        taskDraggingFrom: null,
        taskDraggingTo: null,
        error: null
      },
      params: {
        projectId: 1
      },
      fetchLists: expect.createSpy(),
      fetchProject: expect.createSpy(),
      fetchListOptions: expect.createSpy(),
      closeFlash: expect.createSpy(),
      taskDrop: expect.createSpy(),
      openNewListModal: expect.createSpy(),
      taskDragStart: expect.createSpy(),
      openNewTaskModal: expect.createSpy()
    }
    it('should render loading window', () => {
      const { output } = setup(state)
      let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectTitleWrapper, items ] = output.props.children
      expect(wholeLoading.type).toBe('div')
      expect(wholeLoading.props.className).toBe('whole-loading')
    })
  })

  context('when list modal open', () => {
    let state = {
      ListReducer: {
        isListModalOpen: true,
        isTaskModalOpen: false,
        isListEditModalOpen: false,
        newList: {title: "", color: "0effff"},
        newTask: {title: ""},
        lists: [
          {
            Id: 1,
            Title: "list1",
            ListTasks: [
              {
                Id: 1,
                Title: "task1"
              }, {
                Id: 2,
                Title: "task2"
              }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }
        ],
        listOptions: [],
        selectedListOption: null,
        selectedList: null,
        project: {
          Title: "testProject"
        },
        isTaskDraggingOver: false,
        taskDraggingFrom: null,
        taskDraggingTo: null,
        error: null
      },
      params: {
        projectId: 1
      },
      fetchLists: expect.createSpy(),
      fetchProject: expect.createSpy(),
      fetchListOptions: expect.createSpy(),
      closeFlash: expect.createSpy(),
      taskDrop: expect.createSpy(),
      openNewListModal: expect.createSpy(),
      taskDragStart: expect.createSpy(),
      openNewTaskModal: expect.createSpy()
    }
    it('should render modal', () => {
      const { output, props } = setup(state)
      let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectTitleWrapper, items ] = output.props.children
      expect(listModal.props.isOpen).toBe(true)
    })
  })
  context('when task modal open', () => {
    let state = {
      ListReducer: {
        isListModalOpen: false,
        isTaskModalOpen: true,
        isListEditModalOpen: false,
        newList: {title: "", color: "0effff"},
        newTask: {title: ""},
        lists: [
          {
            Id: 1,
            Title: "list1",
            ListTasks: [
              {
                Id: 1,
                Title: "task1"
              }, {
                Id: 2,
                Title: "task2"
              }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }
        ],
        listOptions: [],
        selectedListOption: null,
        selectedList: null,
        project: {
          Title: "testProject"
        },
        isTaskDraggingOver: false,
        taskDraggingFrom: null,
        taskDraggingTo: null,
        error: null
      },
      params: {
        projectId: 1
      },
      fetchLists: expect.createSpy(),
      fetchProject: expect.createSpy(),
      fetchListOptions: expect.createSpy(),
      closeFlash: expect.createSpy(),
      taskDrop: expect.createSpy(),
      openNewListModal: expect.createSpy(),
      taskDragStart: expect.createSpy(),
      openNewTaskModal: expect.createSpy()
    }
    it('should render task modal', () => {
      const { output } = setup(state)
      let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectTitleWrapper, items ] = output.props.children
      expect(taskModal.props.isOpen).toBe(true)
    })
  })
  context('when list edit modal open', () => {
    let state = {
      ListReducer: {
        isListModalOpen: false,
        isTaskModalOpen: false,
        isListEditModalOpen: true,
        newList: {title: "", color: "0effff"},
        newTask: {title: ""},
        lists: [
          {
            Id: 1,
            Title: "list1",
            ListTasks: [
              {
                Id: 1,
                Title: "task1"
              }, {
                Id: 2,
                Title: "task2"
              }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }
        ],
        listOptions: [],
        selectedListOption: null,
        selectedList: null,
        project: {
          Title: "testProject"
        },
        isTaskDraggingOver: false,
        taskDraggingFrom: null,
        taskDraggingTo: null,
        error: null
      },
      params: {
        projectId: 1
      },
      fetchLists: expect.createSpy(),
      fetchProject: expect.createSpy(),
      fetchListOptions: expect.createSpy(),
      closeFlash: expect.createSpy(),
      taskDrop: expect.createSpy(),
      openNewListModal: expect.createSpy(),
      taskDragStart: expect.createSpy(),
      openNewTaskModal: expect.createSpy()
    }
    it('should render list edit modal', () => {
      const { output, props } = setup(state)
      let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectTitleWrapper, items ] = output.props.children
      expect(listEditModal.props.isOpen).toBe(true)
    })
  })
})
