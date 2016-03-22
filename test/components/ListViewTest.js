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
        isProjectEditModalOpen: false,
        newList: {title: "", color: "0effff"},
        newTask: {title: ""},
        lists: [
          {
            ID: 1,
            Title: "list1",
            ListTasks: [
              {
                ID: 1,
                Title: "task1"
              }, {
                ID: 2,
                Title: "task2"
              }
            ]
          }, {
            ID: 2,
            Title: "list2",
            ListTasks: []
          }
        ],
        noneList: {
          ID: 3,
          ListTasks: [
            {
              ID: 3,
              Title: "task3"
            }
          ]
        },
        listOptions: [],
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

      let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectEditModal, projectTitleWrapper, items, noneList ] = output.props.children
      let [ icon, projectTitle ] = projectTitleWrapper.props.children
      let [ title, editButton ] = projectTitle.props.children
      expect(title).toBe('testProject')

      expect(items.props.className).toBe('items')

      let [ list, button ] = items.props.children
      expect(list[0].props['data-id']).toBe(1)
      expect(list[1].props['data-id']).toBe(2)

      list[0].props.onDrop()
      expect(props.taskDrop.calls.length).toBe(1)

      button.props.onClick()
      expect(props.openNewListModal.calls.length).toBe(1)

      // list which have tasks
      let [ firstListMenu, firstListTitle, firstTasks ] = list[0].props.children
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

      let [ tasks, newTask ] = noneList.props.children.props.children
      expect(tasks[0].props['data-id']).toBe(3)
      expect(newTask.props.className).toBe('new-task pure-button button-blue')

      newTask.props.onClick()
      expect(props.openNewTaskModal.calls.length).toBe(2)
    })
  })

  context('when one error, not modal', () => {
    let state = {
      ListReducer: {
        isListModalOpen: false,
        isTaskModalOpen: false,
        isListEditModalOpen: false,
        isProjectEditModalOpen: false,
        newList: {title: "", color: "0effff"},
        newTask: {title: ""},
        lists: [
          {
            ID: 1,
            Title: "list1",
            ListTasks: [
              {
                ID: 1,
                Title: "task1"
              }, {
                ID: 2,
                Title: "task2"
              }
            ]
          }, {
            ID: 2,
            Title: "list2",
            ListTasks: []
          }
        ],
        noneList: {ID: 0, ListTasks: [] },
        listOptions: [],
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
        error: "Server Error"
      },
      params: {
        projectID: 1
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
      let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectEditModal, projectTitleWrapper, items, noneList ] = output.props.children
      expect(flash.props.children).toBe('Server Error')
    })
  })
  context('when whole loading is open', () => {
    let state = {
      ListReducer: {
        isListModalOpen: false,
        isTaskModalOpen: false,
        isListEditModalOpen: false,
        isProjectEditModalOpen: false,
        isLoading: true,
        newList: {title: "", color: "0effff"},
        newTask: {title: ""},
        lists: [
          {
            ID: 1,
            Title: "list1",
            ListTasks: [
              {
                ID: 1,
                Title: "task1"
              }, {
                ID: 2,
                Title: "task2"
              }
            ]
          }, {
            ID: 2,
            Title: "list2",
            ListTasks: []
          }
        ],
        noneList: {ID: 0, ListTasks: [] },
        listOptions: [],
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
      closeFlash: expect.createSpy(),
      taskDrop: expect.createSpy(),
      openNewListModal: expect.createSpy(),
      taskDragStart: expect.createSpy(),
      openNewTaskModal: expect.createSpy()
    }
    it('should render loading window', () => {
      const { output } = setup(state)
      let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectEditModal, projectTitleWrapper, items, noneList ] = output.props.children
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
        isProjectEditModalOpen: false,
        newList: {title: "", color: "0effff"},
        newTask: {title: ""},
        lists: [
          {
            ID: 1,
            Title: "list1",
            ListTasks: [
              {
                ID: 1,
                Title: "task1"
              }, {
                ID: 2,
                Title: "task2"
              }
            ]
          }, {
            ID: 2,
            Title: "list2",
            ListTasks: []
          }
        ],
        noneList: {ID: 0, ListTasks: [] },
        listOptions: [],
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
      closeFlash: expect.createSpy(),
      taskDrop: expect.createSpy(),
      openNewListModal: expect.createSpy(),
      taskDragStart: expect.createSpy(),
      openNewTaskModal: expect.createSpy()
    }
    it('should render modal', () => {
      const { output, props } = setup(state)
      let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectEditModal, projectTitleWrapper, items, noneList ] = output.props.children
      expect(listModal.props.isOpen).toBe(true)
    })
  })
  context('when task modal open', () => {
    let state = {
      ListReducer: {
        isListModalOpen: false,
        isTaskModalOpen: true,
        isListEditModalOpen: false,
        isProjectEditModalOpen: false,
        newList: {title: "", color: "0effff"},
        newTask: {title: ""},
        lists: [
          {
            ID: 1,
            Title: "list1",
            ListTasks: [
              {
                ID: 1,
                Title: "task1"
              }, {
                ID: 2,
                Title: "task2"
              }
            ]
          }, {
            ID: 2,
            Title: "list2",
            ListTasks: []
          }
        ],
        noneList: {ID: 0, ListTasks: [] },
        listOptions: [],
        selectedListOption: null,
        selectedList: null,
        project: {
          Title: "testProject"
        },
        selectedProject: {Title: "", Description: "", RepositoryID: 0},
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
      closeFlash: expect.createSpy(),
      taskDrop: expect.createSpy(),
      openNewListModal: expect.createSpy(),
      taskDragStart: expect.createSpy(),
      openNewTaskModal: expect.createSpy()
    }
    it('should render task modal', () => {
      const { output } = setup(state)
      let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectEditModal, projectTitleWrapper, items, noneList ] = output.props.children
      expect(taskModal.props.isOpen).toBe(true)
    })
  })
  context('when list edit modal open', () => {
    context('when project does not have repository', () => {
      let state = {
        ListReducer: {
          isListModalOpen: false,
          isTaskModalOpen: false,
          isListEditModalOpen: true,
          isProjectEditModalOpen: false,
          newList: {title: "", color: "0effff"},
          newTask: {title: ""},
          lists: [
            {
              ID: 1,
              Title: "list1",
              ListTasks: [
                {
                  ID: 1,
                  Title: "task1"
                }, {
                  ID: 2,
                  Title: "task2"
                }
              ]
            }, {
              ID: 2,
              Title: "list2",
              ListTasks: []
            }
          ],
          noneList: {ID: 0, ListTasks: [] },
          listOptions: [
            {
              ID: 1,
              Action: "close"
            }, {
              ID: 2,
              Action: "open"
            }
          ],
          selectedListOption: {
            ID: 1,
            Action: "close"
          },
          selectedList: null,
          project: {
            Title: "testProject",
            RepositoryID: 0
          },
          selectedProject: {Title: "", Description: "", RepositoryID: 0},
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
      it('should render list edit modal without action', () => {
        const { output, props } = setup(state)
        let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectEditModal, projectTitleWrapper, items, noneList ] = output.props.children
        expect(listEditModal.props.isOpen).toBe(true)

        let listForm = listEditModal.props.children
        let form = listForm.props.children
        let fieldset = form.props.children
        let [ legend, titleLabel, titleInput, colorLabel, colorInput, nil, formAction ] = fieldset.props.children
        formAction.props.children.props.onClick()
        expect(props.fetchUpdateList.calls.length).toBe(1)
      })
    })
    context('when project has repository', () => {
      let state = {
        ListReducer: {
          isListModalOpen: false,
          isTaskModalOpen: false,
          isListEditModalOpen: true,
          isProjectEditModalOpen: false,
          newList: {title: "", color: "0effff"},
          newTask: {title: ""},
          lists: [
            {
              ID: 1,
              Title: "list1",
              ListTasks: [
                {
                  ID: 1,
                  Title: "task1"
                }, {
                  ID: 2,
                  Title: "task2"
                }
              ]
            }, {
              ID: 2,
              Title: "list2",
              ListTasks: []
            }
          ],
          noneList: {ID: 0, ListTasks: [] },
          listOptions: [
            {
              ID: 1,
              Action: "close"
            }, {
              ID: 2,
              Action: "open"
            }
          ],
          selectedListOption: {
            ID: 1,
            Action: "close"
          },
          selectedList: null,
          project: {
            Title: "testProject",
            Description: "description",
            RepositoryID: 1,
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
      it('should render list edit modal with action', () => {
        const { output, props } = setup(state)
        let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectEditModal, projectTitleWrapper, items, noneList ] = output.props.children
        expect(listEditModal.props.isOpen).toBe(true)

        let listForm = listEditModal.props.children
        let form = listForm.props.children
        let fieldset = form.props.children
        let [ legend, titleLabel, titleInput, colorLabel, colorInput, actionWrapper, formAction ] = fieldset.props.children
        let [ actionLabel, actionSelect ] = actionWrapper.props.children
        expect(actionSelect.props.value).toBe(state.ListReducer.selectedListOption.ID)
      })
    })
  })

  context('when showIssue is false', () => {
    let state = {
      ListReducer: {
        isListModalOpen: false,
        isTaskModalOpen: false,
        isListEditModalOpen: false,
        isProjectEditModalOpen: false,
        newList: {title: "", color: "0effff"},
        newTask: {title: ""},
        lists: [
          {
            ID: 1,
            Title: "list1",
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
        listOptions: [],
        selectedListOption: null,
        selectedList: null,
        project: {
          Title: "testProject",
          Description: "description",
          RepositoryID: 1,
          ShowIssues: false,
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
      closeFlash: expect.createSpy(),
      taskDrop: expect.createSpy(),
      openNewListModal: expect.createSpy(),
      taskDragStart: expect.createSpy(),
      openNewTaskModal: expect.createSpy()
    }
    it('should not render issues', () => {
      const { output, props } = setup(state)

      let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectEditModal, projectTitleWrapper, items, noneList ] = output.props.children
      let [list, button ] = items.props.children
      let [ firstListMenu, firstListTitle, firstTasks ] = list[0].props.children
      let [ firstListTasks, firstNewTask ] = firstTasks.props.children
      let [ task1, task2 ] = firstListTasks
      expect(task1.props['data-id']).toBe(1)
      expect(task2).toBe(undefined)
      let [ secondListMenu, secondListTitle, secondTasks ] = list[1].props.children
      let [ secondListTasks, secondNewTask ] = secondTasks.props.children
      let [ tasks, newTask ] = noneList.props.children.props.children
      expect(tasks[0]).toBe(undefined)
    })
  })

  context('when showPullRequest is false', () => {
    let state = {
      ListReducer: {
        isListModalOpen: false,
        isTaskModalOpen: false,
        isListEditModalOpen: false,
        isProjectEditModalOpen: false,
        newList: {title: "", color: "0effff"},
        newTask: {title: ""},
        lists: [
          {
            ID: 1,
            Title: "list1",
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
        listOptions: [],
        selectedListOption: null,
        selectedList: null,
        project: {
          Title: "testProject",
          Description: "description",
          RepositoryID: 1,
          ShowIssues: true,
          ShowPullRequests: false
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
      closeFlash: expect.createSpy(),
      taskDrop: expect.createSpy(),
      openNewListModal: expect.createSpy(),
      taskDragStart: expect.createSpy(),
      openNewTaskModal: expect.createSpy()
    }
    it('should not render pull requests', () => {
      const { output, props } = setup(state)

      let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectEditModal, projectTitleWrapper, items, noneList ] = output.props.children
      let [list, button ] = items.props.children
      let [ firstListMenu, firstListTitle, firstTasks ] = list[0].props.children
      let [ firstListTasks, firstNewTask ] = firstTasks.props.children
      let [ task1, task2 ] = firstListTasks
      expect(task1).toBe(undefined)
      expect(task2.props['data-id']).toBe(2)
      let [ secondListMenu, secondListTitle, secondTasks ] = list[1].props.children
      let [ secondListTasks, secondNewTask ] = secondTasks.props.children
      let [ tasks, newTask ] = noneList.props.children.props.children
      expect(tasks[0].props['data-id']).toBe(3)
    })
  })

  describe('github action buttons', () => {
    context('when project does not have repository', () => {
      let state = {
        ListReducer: {
          isListModalOpen: false,
          isTaskModalOpen: false,
          isListEditModalOpen: false,
          isProjectEditModalOpen: false,
          newList: {title: "", color: "0effff"},
          newTask: {title: ""},
          lists: [],
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
          listOptions: [],
          selectedListOption: null,
          selectedList: null,
          project: {
            Title: "testProject",
            Description: "description",
            RepositoryID: 0,
            ShowIssues: true,
            ShowPullRequests: false
          },
          selectedProject: { Title: "", Description: "", RepositoryID: 0, ShowIssues: true, ShowPullRequests: true},
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
        closeFlash: expect.createSpy()
      }
      it('should not render github action buttons', () => {
        const { output, props } = setup(state)

        let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectEditModal, projectTitleWrapper, items, noneList ] = output.props.children
        let [ operation, title ] = projectTitleWrapper.props.children
        expect(operation.props.children.props.children).toBe(undefined)
      })
    })
    context('when project has repository', () => {
      let state = {
        ListReducer: {
          isListModalOpen: false,
          isTaskModalOpen: false,
          isListEditModalOpen: false,
          isProjectEditModalOpen: false,
          newList: {title: "", color: "0effff"},
          newTask: {title: ""},
          lists: [],
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
          listOptions: [],
          selectedListOption: null,
          selectedList: null,
          project: {
            Title: "testProject",
            Description: "description",
            RepositoryID: 1,
            ShowIssues: true,
            ShowPullRequests: false
          },
          selectedProject: { Title: "", Description: "", RepositoryID: 0, ShowIssues: true, ShowPullRequests: true},
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
        closeFlash: expect.createSpy()
      }
      it('should not render github action buttons', () => {
        const { output, props } = setup(state)

        let [ wholeLoading, flash, listModal, taskModal, listEditModal, projectEditModal, projectTitleWrapper, items, noneList ] = output.props.children
        let [ operation, title ] = projectTitleWrapper.props.children
        expect(operation.props.children.props.children.length).toBe(3)
      })
    })
  })
})
