import expect from 'expect'
import ListReducer from '../../frontend/javascripts/reducers/ListReducer'
import * as listActions from '../../frontend/javascripts/actions/ListAction'
import 'babel-polyfill'

// shared examples
function sharedExampleInitState(action) {
  expect(
    ListReducer(undefined, action)
  ).toEqual({
    isListModalOpen: false,
    isTaskModalOpen: false,
    isListEditModalOpen: false,
    isLoading: false,
    newList: {title: "", color: "0effff"},
    newTask: {title: ""},
    lists: [],
    selectedList: null,
    project: null,
    isTaskDraggingOver: false,
    taskDraggingFrom: null,
    taskDraggingTo: null,
    error: null
  })
}


describe('ListReducer', () => {
  describe('initState', () => {
    it('should return initial state', () => {
      sharedExampleInitState({})
    })
  })
  describe('SERVER_ERROR', () => {
    it('should return server error', () => {
      expect(
        ListReducer(null, {
          type: listActions.SERVER_ERROR
        })
      ).toEqual({
        error: "Server Error",
        isLoading: false,
      })
    })
  })
  describe('CLOSE_FLASH', () => {
    it('should close flash', () => {
      expect(
        ListReducer({
          error: "Server Error"
        }, {
          type: listActions.CLOSE_FLASH
        })
      ).toEqual({
        error: null
      })
    })
  })
  describe('OPEN_NEW_LIST', () => {
    it('should open list modal', () => {
      expect(
        ListReducer({
          isListModalOpen: false
        }, {
          type: listActions.OPEN_NEW_LIST,
          isListModalOpen: true
        })
      ).toEqual({
        isListModalOpen: true
      })
    })
  })
  describe('CLOSE_NEW_LIST', () => {
    it('should close list modal', () => {
      expect(
        ListReducer({
          isListModalOpen: true
        }, {
          type: listActions.OPEN_NEW_LIST,
          isListModalOpen: false
        })
      ).toEqual({
        isListModalOpen: false
      })
    })
  })
  describe('OPEN_NEW_TASK', () => {
    it('should open new task modal', () => {
      expect(
        ListReducer({
          isTaskModalOpen: false,
          selectedList: null
        }, {
          type: listActions.OPEN_NEW_TASK,
          isTaskModalOpen: true,
          list: "newList"
        })
      ).toEqual({
        isTaskModalOpen: true,
        selectedList: "newList"
      })
    })
  })
  describe('CLOSE_NEW_TASK', () => {
    it('should close new task modal', () => {
      expect(
        ListReducer({
          isTaskModalOpen: true,
          selectedList: "newList"
        }, {
          type: listActions.CLOSE_NEW_TASK,
          isTaskModalOpen: false
        })
      ).toEqual({
        isTaskModalOpen: false,
        selectedList: null
      })
    })
  })
  describe('OPEN_EDIT_LIST', () => {
    it('should open edit list modal', () => {
      expect(
        ListReducer({
          isListEditModalOpen: false,
          selectedList: null
        }, {
          type: listActions.OPEN_EDIT_LIST,
          isListEditModalOpen: true,
          list: "editList"
        })
      ).toEqual({
        isListEditModalOpen: true,
        selectedList: "editList"
      })
    })
  })
  describe('CLOSE_EDIT_LIST', () => {
    it('should close edit list modal', () => {
      expect(
        ListReducer({
          isListEditModalOpen: true,
          selectedList: "editList"
        }, {
          type: listActions.CLOSE_EDIT_LIST,
          isListEditModalOpen: false
        })
      ).toEqual({
        isListEditModalOpen: false,
        selectedList: null
      })
    })
  })
  describe('UPDATE_NEW_LIST_TITLE', () => {
    it('should update list title', () => {
      expect(
        ListReducer({
          newList: { title: "newL", color: "" }
        }, {
          type: listActions.UPDATE_NEW_LIST_TITLE,
          title: "newList"
        })
      ).toEqual({
        newList: { title: "newList", color: "" }
      })
    })
  })
  describe('UPDATE_NEW_LIST_COLOR', () => {
    it('should update list color', () => {
      expect(
        ListReducer({
          newList: { title: "newList", color: "30b" }
        }, {
          type: listActions.UPDATE_NEW_LIST_COLOR,
          color: "30bfe"
        })
      ).toEqual({
        newList: { title: "newList", color: "30bfe" }
      })
    })
  })
  describe('UPDATE_SELECTED_LIST_TITLE', () => {
    it('should update selected list title', () => {
      expect(
        ListReducer({
          selectedList: { Title: "selectedL", Color: "" }
        }, {
          type: listActions.UPDATE_SELECTED_LIST_TITLE,
          title: "selectedList"
        })
      ).toEqual({
        selectedList: { Title: "selectedList", Color: "" }
      })
    })
  })
  describe('UPDATE_SELECTED_LIST_COLOR', () => {
    it('should update selected list color', () => {
      expect(
        ListReducer({
          selectedList: { Title: "selectedList", Color: "30b" }
        }, {
          type: listActions.UPDATE_SELECTED_LIST_COLOR,
          color: "30bef"
        })
      ).toEqual({
        selectedList: { Title: "selectedList", Color: "30bef" }
      })
    })
  })
  describe('UPDATE_NEW_TASK_TITLE', () => {
    it('should update new task title', () => {
      expect(
        ListReducer({
          newTask: { title: "" }
        }, {
          type: listActions.UPDATE_NEW_TASK_TITLE,
          title: "newTask"
        })
      ).toEqual({
        newTask: { title: "newTask" }
      })
    })
  })
  describe('RECEIVE_LISTS', () => {
    context('when received lists is empty', () => {
      it('should return empty lists', () => {
        expect(
          ListReducer({
            lists: null
          }, {
            type: listActions.RECEIVE_LISTS,
            lists: null
          })
        ).toEqual({
          lists: [],
          isLoading: false
        })
      })
    })
    context('when receive lists and empty listTasks', () => {
      it('should return lists and empty listTasks', () => {
        expect(
          ListReducer({
            lists: []
          }, {
            type: listActions.RECEIVE_LISTS,
            lists: [
              { title: "list1", ListTasks: null },
              { title: "list2", ListTasks: null }
            ]
          })
        ).toEqual({
          lists: [
            { title: "list1", ListTasks: [] },
            { title: "list2", ListTasks: [] }
          ],
          isLoading: false
        })
      })
    })
    context('when receive list and listTasks', () => {
      it('should return lists and empty listTasks', () => {
        expect(
          ListReducer({
            lists: []
          }, {
            type: listActions.RECEIVE_LISTS,
            lists: [
              {
                title: "list1",
                ListTasks: [
                  { title: "task1" }
                ]
              },
              {
                title: "list2",
                ListTasks: [
                  { title: "task2" }
                ]
              },
            ]
          })
        ).toEqual({
          lists: [
            {
              title: "list1",
              ListTasks: [
                { title: "task1" }
              ]
            },
            {
              title: "list2",
              ListTasks: [
                { title: "task2" }
              ]
            },
          ],
          isLoading: false
        })
      })
    })
  })
  describe('RECEIVE_CREATE_LIST', () => {
    context('when receive list and empty ListTasks', () => {
      it('should return list and empty ListTasks', () => {
        expect(
          ListReducer({
            lists: [
              { title: "list1",
                ListTasks: [
                  { title: "task1" },
                  { title: "task2" }
                ]
              },
              { title: "list2",
                ListTasks: []
              }
            ],
            isListModalOpen: true,
            newList: { title: "list3", color: "ffffff" }
          }, {
            type: listActions.RECEIVE_CREATE_LIST,
            list: { title: "list3", ListTasks: null }
          })
        ).toEqual({
          lists: [
            { title: "list1",
              ListTasks: [
                { title: "task1" },
                { title: "task2" }
              ]
            },
            { title: "list2",
              ListTasks: []
            },
            { title: "list3",
              ListTasks: []
            }
          ],
          isListModalOpen: false,
          newList: { title: "", color: "0effff" }
        })
      })
    })
    context('when receive list and ListTasks', () => {
      it('should return list and empty ListTasks', () => {
        expect(
          ListReducer({
            lists: [
              { title: "list1",
                ListTasks: [
                  { title: "task1" },
                  { title: "task2" }
                ]
              },
              { title: "list2",
                ListTasks: []
              }
            ],
            isListModalOpen: false,
            newList: { title: "", color: "ffffff" }
          }, {
            type: listActions.RECEIVE_CREATE_LIST,
            list: { title: "list3", ListTasks: [ { title: "task3" } ] }
          })
        ).toEqual({
          lists: [
            { title: "list1",
              ListTasks: [
                { title: "task1" },
                { title: "task2" }
              ]
            }, {
              title: "list2",
              ListTasks: []
            }, {
              title: "list3",
              ListTasks: [
                { title: "task3" }
              ]
            }
          ],
          isListModalOpen: false,
          newList: { title: "", color: "0effff" }
        })
      })
    })
  })
  describe('RECEIVE_CREATE_TASK', () => {
    it('should return lists contain new task', () => {
      expect(
        ListReducer({
          lists: [{
            Id: 1,
            Title: "list1",
            ListTasks: [
              { ListId: 1, Title: "task1", Color: "0effff" }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }],
          newTask: { title: "task2", color: "ffffff" },
          isTaskModalOpen: true
        }, {
          type: listActions.RECEIVE_CREATE_TASK,
          task: { ListId: 1, Title: "task2", Color: "ffffff" }
        })
      ).toEqual({
        lists: [{
          Id: 1,
          Title: "list1",
          ListTasks: [
            { ListId: 1, Title: "task1", Color: "0effff" },
            { ListId: 1, Title: "task2", Color: "ffffff" }
          ]
        }, {
          Id: 2,
          Title: "list2",
          ListTasks: []
        }],
        newTask: { title: "", color: "0effff" },
        isTaskModalOpen: false
      })
    })
  })
  describe('RECEIVE_UPDATE_LIST', () => {
    it('should return lists with new list', () => {
      expect(
        ListReducer({
          lists: [{
            Id: 1,
            Title: "list1",
            ListTasks: [
              { ListId: 1, Title: "task1", Color: "0effff" }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }],
          isListEditModalOpen: true
        }, {
          type: listActions.RECEIVE_UPDATE_LIST,
          list: {
            Id: 1,
            Title: "updateList1",
            ListTasks: [
              { ListId: 1, Title: "task1", Color: "0effff" }
            ]
          }
        })
      ).toEqual({
        lists: [{
          Id: 1,
          Title: "updateList1",
          ListTasks: [
            { ListId: 1, Title: "task1", Color: "0effff" }
          ]
        }, {
          Id: 2,
          Title: "list2",
          ListTasks: []
        }],
        isListEditModalOpen: false
      })
    })
  })
  describe('TASK_DRAG_START', () => {
    it('should return updated lists', () => {
      expect(
        ListReducer({
          lists: [{
            Id: 1,
            Title: "list1",
            ListTasks: [
              { Id: 1, ListId: 1, Title: "task1", Color: "0effff" },
              { Id: 2, ListId: 1, Title: "task2", Color: "0effff" }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }],
          taskDraggingFrom: null
        }, {
          type: listActions.TASK_DRAG_START,
          taskDragFromList: {
            dataset: {
              id: 1
            }
          },
          taskDragTarget: {
            dataset: {
              id: 2
            }
          }
        })
      ).toEqual({
        lists: [{
          Id: 1,
          Title: "list1",
          ListTasks: [
            { Id: 1, ListId: 1, Title: "task1", Color: "0effff" },
            { Id: 2, ListId: 1, Title: "task2", Color: "0effff" }
          ]
        }, {
          Id: 2,
          Title: "list2",
          ListTasks: []
        }],
        taskDraggingFrom: {
          fromList: {
            Id: 1,
            Title: "list1",
            ListTasks: [
              { Id: 1, ListId: 1, Title: "task1", Color: "0effff" },
              { Id: 2, ListId: 1, Title: "task2", Color: "0effff" }
            ]
          },
          fromTask: {
            Id: 2,
            ListId: 1,
            Title: "task2",
            Color: "0effff"
          }
        }
      })
    })
  })
  describe('TASK_DRAG_LEAVE', () => {
    it('should return lists do not contain arrow', () => {
      expect(
        ListReducer({
          lists: [{
            Id: 1,
            Title: "list1",
            ListTasks: [
              { Id: 1, ListId: 1, Title: "task1", Color: "0effff" },
              { Id: 2, ListId: 1, Title: "task2", Color: "0effff" }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: [ { draggedOn: true } ]
          }],
          isTaskDraggingOver: true,
          taskDraggingTo: {
            toList: {
              Id: 2,
              Title: "list2",
              ListTasks: []
            },
            prevToTask: null
          }
        }, {
          type: listActions.TASK_DRAG_LEAVE,
        })
      ).toEqual({
        isTaskDraggingOver: false,
        taskDraggingTo: null,
        lists: [{
          Id: 1,
          Title: "list1",
          ListTasks: [
            { Id: 1, ListId: 1, Title: "task1", Color: "0effff" },
            { Id: 2, ListId: 1, Title: "task2", Color: "0effff" }
          ]
        }, {
          Id: 2,
          Title: "list2",
          ListTasks: []
        }]
      })
    })
  })
  describe('TASK_DROP', () => {
    it('should return do not contain arrow', () => {
      expect(
        ListReducer({
          lists: [{
            Id: 1,
            Title: "list1",
            ListTasks: [
              { Id: 1, ListId: 1, Title: "task1", Color: "0effff" },
              { Id: 2, ListId: 1, Title: "task2", Color: "0effff" }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: [ { draggedOn: true } ]
          }],
          isTaskDraggingOver: true,
          taskDraggingTo: {
            toList: {
              Id: 2,
              Title: "list2",
              ListTasks: []
            },
            prevToTask: null
          },
          taskDraggingFrom: {
            fromList: {
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Color: "0effff" },
                { Id: 2, ListId: 1, Title: "task2", Color: "0effff" }
              ]
            },
            fromTask: {
              Id: 2,
              ListId: 1,
              Title: "task2",
              Color: "0effff"
            }
          }
        }, {
          type: listActions.TASK_DROP,
        })
      ).toEqual({
        isTaskDraggingOver: false,
        taskDraggingTo: null,
        taskDraggingFrom: null,
        lists: [{
          Id: 1,
          Title: "list1",
          ListTasks: [
            { Id: 1, ListId: 1, Title: "task1", Color: "0effff" },
            { Id: 2, ListId: 1, Title: "task2", Color: "0effff" }
          ]
        }, {
          Id: 2,
          Title: "list2",
          ListTasks: []
        }]
      })
    })
  })
  describe('TASK_DRAG_OVER', () => {
    context('when drag to last of tasks', () => {
      it('should return list s with arrow', () => {
        expect(
          ListReducer({
            lists: [{
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Color: "0effff" },
                { Id: 2, ListId: 1, Title: "task2", Color: "0effff" }
              ]
            }, {
              Id: 2,
              Title: "list2",
              ListTasks: [
                { Id: 3, ListId: 2, Title: "task3", Color: "0effff" },
                { Id: 4, ListId: 2, Title: "task4", Color: "0effff" }]
            }],
            taskDraggingTo: null,
            taskDraggingFrom: {
              fromList: {
                Id: 1,
                Title: "list1",
                ListTasks: [
                  { Id: 1, ListId: 1, Title: "task1", Color: "0effff" },
                  { Id: 2, ListId: 1, Title: "task2", Color: "0effff" }
                ]
              },
              fromTask: {
                Id: 2,
                ListId: 1,
                Title: "task2",
                Color: "0effff"
              }
            }
          }, {
            type: listActions.TASK_DRAG_OVER,
            taskDragToList: {
              dataset: {
                id: 2
              }
            },
            taskDragToTask: {
              className: null
            }
          })
        ).toEqual({
          isTaskDraggingOver: true,
          taskDraggingTo: {
            toList: {
              Id: 2,
              Title: "list2",
              ListTasks: [
                { Id: 3, ListId: 2, Title: "task3", Color: "0effff" },
                { Id: 4, ListId: 2, Title: "task4", Color: "0effff" },
                { draggedOn: true }
              ]
            },
            prevToTask: null
          },
          lists: [{
            Id: 1,
            Title: "list1",
            ListTasks: [
              { Id: 1, ListId: 1, Title: "task1", Color: "0effff" },
              { Id: 2, ListId: 1, Title: "task2", Color: "0effff" }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: [
              { Id: 3, ListId: 2, Title: "task3", Color: "0effff" },
              { Id: 4, ListId: 2, Title: "task4", Color: "0effff" },
              { draggedOn: true }
            ]
          }],
          taskDraggingFrom: {
            fromList: {
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Color: "0effff" },
                { Id: 2, ListId: 1, Title: "task2", Color: "0effff" }
              ]
            },
            fromTask: {
              Id: 2,
              ListId: 1,
              Title: "task2",
              Color: "0effff"
            }
          }
        })
      })
    })
    context('when drag to half way of tasks', () => {
      it('should return lists with arrow', () => {
        expect(
          ListReducer({
            lists: [{
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Color: "0effff" },
                { Id: 2, ListId: 1, Title: "task2", Color: "0effff" }
              ]
            }, {
              Id: 2,
              Title: "list2",
              ListTasks: [
                { Id: 3, ListId: 2, Title: "task3", Color: "0effff" },
                { Id: 4, ListId: 2, Title: "task4", Color: "0effff" }]
            }],
            taskDraggingTo: null,
            taskDraggingFrom: {
              fromList: {
                Id: 1,
                Title: "list1",
                ListTasks: [
                  { Id: 1, ListId: 1, Title: "task1", Color: "0effff" },
                  { Id: 2, ListId: 1, Title: "task2", Color: "0effff" }
                ]
              },
              fromTask: {
                Id: 2,
                ListId: 1,
                Title: "task2",
                Color: "0effff"
              }
            }
          }, {
            type: listActions.TASK_DRAG_OVER,
            taskDragToList: {
              dataset: {
                id: 2
              }
            },
            taskDragToTask: {
              className: "task",
              dataset: {
                id: 4
              }
            }
          })
        ).toEqual({
          isTaskDraggingOver: true,
          taskDraggingTo: {
            toList: {
              Id: 2,
              Title: "list2",
              ListTasks: [
                { Id: 3, ListId: 2, Title: "task3", Color: "0effff" },
                { draggedOn: true },
                { Id: 4, ListId: 2, Title: "task4", Color: "0effff" }
              ]
            },
            prevToTask: { Id: 4, ListId: 2, Title: "task4", Color: "0effff" }
          },
          lists: [{
            Id: 1,
            Title: "list1",
            ListTasks: [
              { Id: 1, ListId: 1, Title: "task1", Color: "0effff" },
              { Id: 2, ListId: 1, Title: "task2", Color: "0effff" }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: [
              { Id: 3, ListId: 2, Title: "task3", Color: "0effff" },
              { draggedOn: true },
              { Id: 4, ListId: 2, Title: "task4", Color: "0effff" }
            ]
          }],
          taskDraggingFrom: {
            fromList: {
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Color: "0effff" },
                { Id: 2, ListId: 1, Title: "task2", Color: "0effff" }
              ]
            },
            fromTask: {
              Id: 2,
              ListId: 1,
              Title: "task2",
              Color: "0effff"
            }
          }
        })

      })
    })
  })
})
