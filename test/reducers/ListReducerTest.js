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
    isProjectEditModalOpen: false,
    isLoading: false,
    newList: {title: "", color: "008ed4"},
    newTask: {title: "", description: ""},
    lists: [],
    noneList: {Id: 0, ListTasks: []},
    listOptions: [],
    selectedListOption: null,
    selectedList: null,
    project: null,
    selectedProject: {Title: "", Description: "", RepositoryID: 0},
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
  describe('REQUEST_FETCH_GITHUB', () => {
    it('should render whole loading window', () => {
      expect(
        ListReducer({
          isLoading: false
        }, {
          type: listActions.REQUEST_FETCH_GITHUB
        })
      ).toEqual({
        isLoading: true
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
          list: {
            Title: "newList"
          }
        })
      ).toEqual({
        isTaskModalOpen: true,
        selectedList: {
          Title: "newList"
        }
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
          selectedList: null,
          selectedListOption: null
        }, {
          type: listActions.OPEN_EDIT_LIST,
          isListEditModalOpen: true,
          list: {
            ListOptionId: 1
          }
        })
      ).toEqual({
        isListEditModalOpen: true,
        selectedList: {
          ListOptionId: 1
        },
        selectedListOption: {
          Id: 1
        }
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
        selectedList: null,
        selectedListOption: null
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
            lists: null,
            noneList: {Id: 0, ListTasks: []},
            isLoading: true
          }, {
            type: listActions.RECEIVE_LISTS,
            lists: null,
            noneList: {Id: 1, ListTasks: []}
          })
        ).toEqual({
          lists: [],
          noneList: {Id: 1, ListTasks: []},
          isLoading: false
        })
      })
    })
    context('when receive lists and empty listTasks', () => {
      it('should return lists and empty listTasks', () => {
        expect(
          ListReducer({
            lists: [],
            noneList: {Id: 0, ListTasks: []}
          }, {
            type: listActions.RECEIVE_LISTS,
            lists: [
              { title: "list1", ListTasks: null },
              { title: "list2", ListTasks: null }
            ],
            noneList: null
          })
        ).toEqual({
          lists: [
            { title: "list1", ListTasks: [] },
            { title: "list2", ListTasks: [] }
          ],
          noneList: {Id: 0, ListTasks: []},
          isLoading: false
        })
      })
    })
    context('when receive list and listTasks', () => {
      it('should return lists and empty listTasks', () => {
        expect(
          ListReducer({
            lists: [],
            noneList: {Id: 0, ListTasks: []}
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
            ],
            noneList: null
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
          noneList: {Id: 0, ListTasks: []},
          isLoading: false
        })
      })
    })
  })
  describe('REQUEST_CREATE_LIST', () => {
    it('should open whole loading window', () => {
      expect(
        ListReducer({
          isLoading: false
        }, {
          type: listActions.REQUEST_CREATE_LIST
        })
      ).toEqual({
        isLoading: true
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
            newList: { title: "list3", color: "ffffff" },
            isLoading: true
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
          newList: { title: "", color: "008ed4" },
          isLoading: false
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
            newList: { title: "", color: "ffffff" },
            isLoading: true
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
          newList: { title: "", color: "008ed4" },
          isLoading: false
        })
      })
    })
  })
  describe('REQUEST_CREATE_TASK', () => {
    it('should open whole loading window', () => {
      expect(
        ListReducer({
          isLoading: false
        }, {
          type: listActions.REQUEST_CREATE_TASK
        })
      ).toEqual({
        isLoading: true
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
              { ListId: 1, Title: "task1", Description: "fugafuga" }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }],
          noneList: {Id: 0, ListTasks: [] },
          newTask: { title: "task2", description: "hogehoge" },
          isTaskModalOpen: true,
          isLoading: true
        }, {
          type: listActions.RECEIVE_CREATE_TASK,
          task: { ListId: 1, Title: "task2", Description: "hogehoge" }
        })
      ).toEqual({
        lists: [{
          Id: 1,
          Title: "list1",
          ListTasks: [
            { ListId: 1, Title: "task1", Description: "fugafuga" },
            { ListId: 1, Title: "task2", Description: "hogehoge" }
          ]
        }, {
          Id: 2,
          Title: "list2",
          ListTasks: []
        }],
        noneList: {Id: 0, ListTasks: [] },
        newTask: { title: "", description: "" },
        isTaskModalOpen: false,
        isLoading: false
      })
    })
  })
  describe('REQUEST_UPDATE_LIST', () => {
    it('should open whole loading window', () => {
      expect(
        ListReducer({
          isLoading: false
        }, {
          type: listActions.REQUEST_UPDATE_LIST
        })
      ).toEqual({
        isLoading: true
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
              { ListId: 1, Title: "task1", Description: "hogehoge" }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }],
          isListEditModalOpen: true,
          isLoading: true
        }, {
          type: listActions.RECEIVE_UPDATE_LIST,
          list: {
            Id: 1,
            Title: "updateList1",
            ListTasks: [
              { ListId: 1, Title: "task1", Description: "hogehoge" }
            ]
          }
        })
      ).toEqual({
        lists: [{
          Id: 1,
          Title: "updateList1",
          ListTasks: [
            { ListId: 1, Title: "task1", Description: "hogehoge" }
          ]
        }, {
          Id: 2,
          Title: "list2",
          ListTasks: []
        }],
        isListEditModalOpen: false,
        isLoading: false
      })
    })
  })
  describe('TASK_DRAG_START', () => {
    context('when drag from list is noneList', () => {
      it('should return updated lists', () => {
        expect(
          ListReducer({
            lists: [{
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              Id: 2,
              Title: "list2",
              ListTasks: []
            }],
            noneList: {
              Id: 3,
              ListTasks: [
                { Id: 3, ListId: 3, Title: "task3", Description: "hogehoge" }
              ]
            },
            taskDraggingFrom: null
          }, {
            type: listActions.TASK_DRAG_START,
            taskDragFromList: {
              dataset: {
                id: 3
              }
            },
            taskDragTarget: {
              dataset: {
                id: 3
              }
            }
          })
        ).toEqual({
          lists: [{
            Id: 1,
            Title: "list1",
            ListTasks: [
              { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
              { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }],
          noneList: {
            Id: 3,
            ListTasks: [
              { Id: 3, ListId: 3, Title: "task3", Description: "hogehoge" }
            ]
          },
          taskDraggingFrom: {
            fromList: {
              Id: 3,
              ListTasks: [
                { Id: 3, ListId: 3, Title: "task3", Description: "hogehoge" }
              ]
            },
            fromTask: {
              Id: 3,
              ListId: 3,
              Title: "task3",
              Description: "hogehoge"
            }
          }
        })
      })
    })
    context('when drag from list is not noneList', () => {
      it('should return updated lists', () => {
        expect(
          ListReducer({
            lists: [{
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              Id: 2,
              Title: "list2",
              ListTasks: []
            }],
            noneList: {Id: 3, ListTasks: [] },
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
              { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
              { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }],
          noneList: {Id: 3, ListTasks: [] },
          taskDraggingFrom: {
            fromList: {
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
              ]
            },
            fromTask: {
              Id: 2,
              ListId: 1,
              Title: "task2",
              Description: "hogehoge"
            }
          }
        })
      })
    })
  })
  describe('TASK_DRAG_LEAVE', () => {
    context('when target list is noneList', () => {
      it('should return lists do not contain arrow', () => {
        expect(
          ListReducer({
            lists: [{
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              Id: 2,
              Title: "list2",
              ListTasks: []
            }],
            noneList: {
              Id: 3,
              ListTasks: [
                { draggedOn: true }
              ]
            },
            isTaskDraggingOver: true,
            taskDraggingTo: {
              toList: {
                Id: 3,
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
              { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
              { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }],
          noneList: {
            Id: 3,
            ListTasks: []
          },
        })
      })
    })
    context('when target list is not noneList', () => {
      it('should return lists do not contain arrow', () => {
        expect(
          ListReducer({
            lists: [{
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              Id: 2,
              Title: "list2",
              ListTasks: [ { draggedOn: true } ]
            }],
            noneList: {Id: 0, ListTasks: [] },
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
              { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
              { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }],
          noneList: {Id: 0, ListTasks: [] },
        })
      })
    })
  })
  describe('REQUEST_MOVE_TASK', () => {
    context('when target list is noneList', () => {
      it('should return do not contain arrow and contain isLoading flag', () => {
        expect(
          ListReducer({
            lists: [{
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              Id: 2,
              Title: "list2",
              ListTasks: []
            }],
            noneList: {
              Id: 3,
              ListTasks: [ {draggedOn: true } ]
            },
            isTaskDraggingOver: true,
            taskDraggingTo: {
              toList: {
                Id: 3,
                ListTasks: []
              },
              prevToTask: null
            },
            taskDraggingFrom: {
              fromList: {
                Id: 1,
                Title: "list1",
                ListTasks: [
                  { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                  { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
                ]
              },
              fromTask: {
                Id: 2,
                ListId: 1,
                Title: "task2",
                Description: "hogehoge"
              }
            }
          }, {
            type: listActions.REQUEST_MOVE_TASK,
          })
        ).toEqual({
          isTaskDraggingOver: false,
          taskDraggingTo: null,
          taskDraggingFrom: null,
          lists: [{
            Id: 1,
            Title: "list1",
            ListTasks: [
              { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
              { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
            ],
            isLoading: true
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }],
          noneList: {
            Id: 3,
            ListTasks: []
          },
        })
      })
    })
    context('when target list is not noneList', () => {
      it('should return do not contain arrow and contain isLoading flag', () => {
        expect(
          ListReducer({
            lists: [{
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              Id: 2,
              Title: "list2",
              ListTasks: [ { draggedOn: true } ]
            }],
            noneList: {Id: 0, ListTasks: [] },
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
                  { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                  { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
                ]
              },
              fromTask: {
                Id: 2,
                ListId: 1,
                Title: "task2",
                Description: "hogehoge"
              }
            }
          }, {
            type: listActions.REQUEST_MOVE_TASK,
          })
        ).toEqual({
          isTaskDraggingOver: false,
          taskDraggingTo: null,
          taskDraggingFrom: null,
          lists: [{
            Id: 1,
            Title: "list1",
            ListTasks: [
              { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
              { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
            ],
            isLoading: true
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: [],
            isLoading: true
          }],
          noneList: {Id: 0, ListTasks: [] },
        })
      })
    })
  })
  describe('TASK_DROP', () => {
    context('when target list is noneList', () => {
      it('should return do not contain arrow and isLoading', () => {
        expect(
          ListReducer({
            lists: [{
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              Id: 2,
              Title: "list2",
              ListTasks: []
            }],
            noneList: {
              Id: 3,
              ListTasks: [ { draggedOn: true } ]
            },
            isTaskDraggingOver: true,
            taskDraggingTo: {
              toList: {
                Id: 3,
                ListTasks: []
              },
              prevToTask: null
            },
            taskDraggingFrom: {
              fromList: {
                Id: 1,
                Title: "list1",
                ListTasks: [
                  { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                  { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
                ]
              },
              fromTask: {
                Id: 2,
                ListId: 1,
                Title: "task2",
                Description: "hogehoge"
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
              { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
              { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }],
          noneList: {
            Id: 3,
            ListTasks: []
          },
        })
      })
    })
    context('when target list is not noneList', () => {
      it('should return do not contain arrow and isLoading', () => {
        expect(
          ListReducer({
            lists: [{
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              Id: 2,
              Title: "list2",
              ListTasks: [ { draggedOn: true } ]
            }],
            noneList: {Id: 0, ListTasks: [] },
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
                  { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                  { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
                ]
              },
              fromTask: {
                Id: 2,
                ListId: 1,
                Title: "task2",
                Description: "hogehoge"
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
              { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
              { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
            ]
          }, {
            Id: 2,
            Title: "list2",
            ListTasks: []
          }],
          noneList: {Id: 0, ListTasks: [] },
        })
      })
    })
  })
  describe('TASK_DRAG_OVER', () => {
    context('when drag to last of tasks', () => {
      context('when target list is noneList', () => {
        it('should return list s with arrow', () => {
          expect(
            ListReducer({
              lists: [{
                Id: 1,
                Title: "list1",
                ListTasks: [
                  { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                  { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
                ]
              }, {
                Id: 2,
                Title: "list2",
                ListTasks: []
              }],
              noneList: {
                Id: 3,
                ListTasks: [
                  { Id: 3, ListId: 3, Title: "task3", Description: "hogehoge" },
                  { Id: 4, ListId: 3, Title: "task4", Description: "hogehoge" }
                ]
              },
              taskDraggingTo: null,
              taskDraggingFrom: {
                fromList: {
                  Id: 1,
                  Title: "list1",
                  ListTasks: [
                    { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                    { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
                  ]
                },
                fromTask: {
                  Id: 2,
                  ListId: 1,
                  Title: "task2",
                  Description: "hogehoge"
                }
              }
            }, {
              type: listActions.TASK_DRAG_OVER,
              taskDragToList: {
                dataset: {
                  id: 3
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
                Id: 3,
                ListTasks: [
                  { Id: 3, ListId: 3, Title: "task3", Description: "hogehoge" },
                  { Id: 4, ListId: 3, Title: "task4", Description: "hogehoge" },
                  { draggedOn: true }
                ]
              },
              prevToTask: null
            },
            lists: [{
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              Id: 2,
              Title: "list2",
              ListTasks: []
            }],
            noneList: {
              Id: 3,
              ListTasks: [
                { Id: 3, ListId: 3, Title: "task3", Description: "hogehoge" },
                { Id: 4, ListId: 3, Title: "task4", Description: "hogehoge" },
                { draggedOn: true }
              ]
            },
            taskDraggingFrom: {
              fromList: {
                Id: 1,
                Title: "list1",
                ListTasks: [
                  { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                  { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
                ]
              },
              fromTask: {
                Id: 2,
                ListId: 1,
                Title: "task2",
                Description: "hogehoge"
              }
            }
          })
        })
      })
      context('when target list is not noneList', () => {
        it('should return list s with arrow', () => {
          expect(
            ListReducer({
              lists: [{
                Id: 1,
                Title: "list1",
                ListTasks: [
                  { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                  { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
                ]
              }, {
                Id: 2,
                Title: "list2",
                ListTasks: [
                  { Id: 3, ListId: 2, Title: "task3", Description: "hogehoge" },
                  { Id: 4, ListId: 2, Title: "task4", Description: "hogehoge" }]
              }],
              noneList: {Id: 0, ListTasks: [] },
              taskDraggingTo: null,
              taskDraggingFrom: {
                fromList: {
                  Id: 1,
                  Title: "list1",
                  ListTasks: [
                    { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                    { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
                  ]
                },
                fromTask: {
                  Id: 2,
                  ListId: 1,
                  Title: "task2",
                  Description: "hogehoge"
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
                  { Id: 3, ListId: 2, Title: "task3", Description: "hogehoge" },
                  { Id: 4, ListId: 2, Title: "task4", Description: "hogehoge" },
                  { draggedOn: true }
                ]
              },
              prevToTask: null
            },
            lists: [{
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              Id: 2,
              Title: "list2",
              ListTasks: [
                { Id: 3, ListId: 2, Title: "task3", Description: "hogehoge" },
                { Id: 4, ListId: 2, Title: "task4", Description: "hogehoge" },
                { draggedOn: true }
              ]
            }],
            noneList: {Id: 0, ListTasks: [] },
            taskDraggingFrom: {
              fromList: {
                Id: 1,
                Title: "list1",
                ListTasks: [
                  { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                  { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
                ]
              },
              fromTask: {
                Id: 2,
                ListId: 1,
                Title: "task2",
                Description: "hogehoge"
              }
            }
          })
        })
      })
    })
    context('when drag to half way of tasks', () => {
      context('when target list is noneList', () => {
        it('should return lists with arrow', () => {
          expect(
            ListReducer({
              lists: [{
                Id: 1,
                Title: "list1",
                ListTasks: [
                  { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                  { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
                ]
              }, {
                Id: 2,
                Title: "list2",
                ListTasks: []
              }],
              noneList: {
                Id: 3,
                ListTasks: [
                  { Id: 3, ListId: 3, Title: "task3", Description: "hogehoge" },
                  { Id: 4, ListId: 3, Title: "task4", Description: "hogehoge" }
                ]
              },
              taskDraggingTo: null,
              taskDraggingFrom: {
                fromList: {
                  Id: 1,
                  Title: "list1",
                  ListTasks: [
                    { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                    { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
                  ]
                },
                fromTask: {
                  Id: 2,
                  ListId: 1,
                  Title: "task2",
                  Description: "hogehoge"
                }
              }
            }, {
              type: listActions.TASK_DRAG_OVER,
              taskDragToList: {
                dataset: {
                  id: 3
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
                Id: 3,
                ListTasks: [
                  { Id: 3, ListId: 3, Title: "task3", Description: "hogehoge" },
                  { draggedOn: true },
                  { Id: 4, ListId: 3, Title: "task4", Description: "hogehoge" }
                ]
              },
              prevToTask: { Id: 4, ListId: 3, Title: "task4", Description: "hogehoge" }
            },
            lists: [{
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              Id: 2,
              Title: "list2",
              ListTasks: []
            }],
            noneList: {
              Id: 3,
              ListTasks: [
                { Id: 3, ListId: 3, Title: "task3", Description: "hogehoge" },
                { draggedOn: true },
                { Id: 4, ListId: 3, Title: "task4", Description: "hogehoge" }
              ]
            },
            taskDraggingFrom: {
              fromList: {
                Id: 1,
                Title: "list1",
                ListTasks: [
                  { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                  { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
                ]
              },
              fromTask: {
                Id: 2,
                ListId: 1,
                Title: "task2",
                Description: "hogehoge"
              }
            }
          })
        })
      })
      context('when target list is not noneList', () => {
        it('should return lists with arrow', () => {
          expect(
            ListReducer({
              lists: [{
                Id: 1,
                Title: "list1",
                ListTasks: [
                  { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                  { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
                ]
              }, {
                Id: 2,
                Title: "list2",
                ListTasks: [
                  { Id: 3, ListId: 2, Title: "task3", Description: "hogehoge" },
                  { Id: 4, ListId: 2, Title: "task4", Description: "hogehoge" }]
              }],
              noneList: {Id: 0, ListTasks: [] },
              taskDraggingTo: null,
              taskDraggingFrom: {
                fromList: {
                  Id: 1,
                  Title: "list1",
                  ListTasks: [
                    { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                    { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
                  ]
                },
                fromTask: {
                  Id: 2,
                  ListId: 1,
                  Title: "task2",
                  Description: "hogehoge"
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
                  { Id: 3, ListId: 2, Title: "task3", Description: "hogehoge" },
                  { draggedOn: true },
                  { Id: 4, ListId: 2, Title: "task4", Description: "hogehoge" }
                ]
              },
              prevToTask: { Id: 4, ListId: 2, Title: "task4", Description: "hogehoge" }
            },
            lists: [{
              Id: 1,
              Title: "list1",
              ListTasks: [
                { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              Id: 2,
              Title: "list2",
              ListTasks: [
                { Id: 3, ListId: 2, Title: "task3", Description: "hogehoge" },
                { draggedOn: true },
                { Id: 4, ListId: 2, Title: "task4", Description: "hogehoge" }
              ]
            }],
            noneList: {Id: 0, ListTasks: [] },
            taskDraggingFrom: {
              fromList: {
                Id: 1,
                Title: "list1",
                ListTasks: [
                  { Id: 1, ListId: 1, Title: "task1", Description: "hogehoge" },
                  { Id: 2, ListId: 1, Title: "task2", Description: "hogehoge" }
                ]
              },
              fromTask: {
                Id: 2,
                ListId: 1,
                Title: "task2",
                Description: "hogehoge"
              }
            }
          })
        })
      })
    })
  })

  describe('CHANGE_SELECTED_LIST_OPTION', () => {
    it('should change selected list option', () => {
      expect(
        ListReducer({
          listOptions: [
            {
              Id: 1,
              Action: "close",
            }, {
              Id: 2,
              Action: "open"
            }
          ],
          selectedListOption: null
        }, {
          type: listActions.CHANGE_SELECTED_LIST_OPTION,
          selectEvent: {
            value: 1
          }
        })
      ).toEqual({
          listOptions: [
            {
              Id: 1,
              Action: "close",
            }, {
              Id: 2,
              Action: "open"
            }
          ],
          selectedListOption: {
            Id: 1,
            Action: "close"
          }
      })
    })
  })
})
