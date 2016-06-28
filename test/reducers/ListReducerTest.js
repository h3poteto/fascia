import expect from 'expect'
import ListReducer from '../../frontend/javascripts/reducers/ListReducer'
import * as listActions from '../../frontend/javascripts/actions/ListAction'
import * as newListModalActions from '../../frontend/javascripts/actions/ListAction/NewListModalAction'
import * as editListModalActions from '../../frontend/javascripts/actions/ListAction/EditListModalAction'
import * as newTaskModalActions from '../../frontend/javascripts/actions/ListAction/NewTaskModalAction'
import * as editProjectModalActions from '../../frontend/javascripts/actions/ListAction/EditProjectModalAction'
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
    noneList: {ID: 0, ListTasks: []},
    listOptions: [],
    selectedListOption: null,
    selectedList: null,
    project: null,
    selectedProject: {Title: "", Description: "", RepositoryID: 0, ShowIssues: true, ShowPullRequests: true},
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

  context('listActions', () => {
    describe('SERVER_ERROR', () => {
      it('should return server error', () => {
        expect(
          ListReducer(null, {
            type: listActions.SERVER_ERROR
          })
        ).toEqual({
          error: "Internal Server Error",
          isLoading: false
        })
      })
    })
    describe('CLOSE_FLASH', () => {
      it('should close flash', () => {
        expect(
          ListReducer({
            error: "Internal Server Error"
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
              ListOptionID: 1
            }
          })
        ).toEqual({
          isListEditModalOpen: true,
          selectedList: {
            ListOptionID: 1
          },
          selectedListOption: {
            ID: 1
          }
        })
      })
    })

    describe('RECEIVE_LISTS', () => {
      context('when received lists is empty', () => {
        it('should return empty lists', () => {
          expect(
            ListReducer({
              lists: null,
              noneList: {ID: 0, ListTasks: []},
              isLoading: true
            }, {
              type: listActions.RECEIVE_LISTS,
              lists: null,
              noneList: {ID: 1, ListTasks: []}
            })
          ).toEqual({
            lists: [],
            noneList: {ID: 1, ListTasks: []},
            isLoading: false
          })
        })
      })
      context('when receive lists and empty listTasks', () => {
        it('should return lists and empty listTasks', () => {
          expect(
            ListReducer({
              lists: [],
              noneList: {ID: 0, ListTasks: []}
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
            noneList: {ID: 0, ListTasks: []},
            isLoading: false
          })
        })
      })
      context('when receive list and listTasks', () => {
        it('should return lists and empty listTasks', () => {
          expect(
            ListReducer({
              lists: [],
              noneList: {ID: 0, ListTasks: []}
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
                }
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
              }
            ],
            noneList: {ID: 0, ListTasks: []},
            isLoading: false
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
                  ID: 1,
                  Title: "list1",
                  ListTasks: [
                    { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                    { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                  ]
                }, {
                  ID: 2,
                  Title: "list2",
                  ListTasks: []
                }],
                noneList: {
                  ID: 3,
                  ListTasks: [
                    { ID: 3, ListID: 3, Title: "task3", Description: "hogehoge" },
                    { ID: 4, ListID: 3, Title: "task4", Description: "hogehoge" }
                  ]
                },
                taskDraggingTo: null,
                taskDraggingFrom: {
                  fromList: {
                    ID: 1,
                    Title: "list1",
                    ListTasks: [
                      { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                      { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                    ]
                  },
                  fromTask: {
                    ID: 2,
                    ListID: 1,
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
                  ID: 3,
                  ListTasks: [
                    { ID: 3, ListID: 3, Title: "task3", Description: "hogehoge" },
                    { ID: 4, ListID: 3, Title: "task4", Description: "hogehoge" },
                    { draggedOn: true }
                  ]
                },
                prevToTask: null
              },
              lists: [{
                ID: 1,
                Title: "list1",
                ListTasks: [
                  { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                  { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                ]
              }, {
                ID: 2,
                Title: "list2",
                ListTasks: []
              }],
              noneList: {
                ID: 3,
                ListTasks: [
                  { ID: 3, ListID: 3, Title: "task3", Description: "hogehoge" },
                  { ID: 4, ListID: 3, Title: "task4", Description: "hogehoge" },
                  { draggedOn: true }
                ]
              },
              taskDraggingFrom: {
                fromList: {
                  ID: 1,
                  Title: "list1",
                  ListTasks: [
                    { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                    { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                  ]
                },
                fromTask: {
                  ID: 2,
                  ListID: 1,
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
                  ID: 1,
                  Title: "list1",
                  ListTasks: [
                    { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                    { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                  ]
                }, {
                  ID: 2,
                  Title: "list2",
                  ListTasks: [
                    { ID: 3, ListID: 2, Title: "task3", Description: "hogehoge" },
                    { ID: 4, ListID: 2, Title: "task4", Description: "hogehoge" }]
                }],
                noneList: {ID: 0, ListTasks: [] },
                taskDraggingTo: null,
                taskDraggingFrom: {
                  fromList: {
                    ID: 1,
                    Title: "list1",
                    ListTasks: [
                      { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                      { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                    ]
                  },
                  fromTask: {
                    ID: 2,
                    ListID: 1,
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
                  ID: 2,
                  Title: "list2",
                  ListTasks: [
                    { ID: 3, ListID: 2, Title: "task3", Description: "hogehoge" },
                    { ID: 4, ListID: 2, Title: "task4", Description: "hogehoge" },
                    { draggedOn: true }
                  ]
                },
                prevToTask: null
              },
              lists: [{
                ID: 1,
                Title: "list1",
                ListTasks: [
                  { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                  { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                ]
              }, {
                ID: 2,
                Title: "list2",
                ListTasks: [
                  { ID: 3, ListID: 2, Title: "task3", Description: "hogehoge" },
                  { ID: 4, ListID: 2, Title: "task4", Description: "hogehoge" },
                  { draggedOn: true }
                ]
              }],
              noneList: {ID: 0, ListTasks: [] },
              taskDraggingFrom: {
                fromList: {
                  ID: 1,
                  Title: "list1",
                  ListTasks: [
                    { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                    { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                  ]
                },
                fromTask: {
                  ID: 2,
                  ListID: 1,
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
                  ID: 1,
                  Title: "list1",
                  ListTasks: [
                    { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                    { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                  ]
                }, {
                  ID: 2,
                  Title: "list2",
                  ListTasks: []
                }],
                noneList: {
                  ID: 3,
                  ListTasks: [
                    { ID: 3, ListID: 3, Title: "task3", Description: "hogehoge" },
                    { ID: 4, ListID: 3, Title: "task4", Description: "hogehoge" }
                  ]
                },
                taskDraggingTo: null,
                taskDraggingFrom: {
                  fromList: {
                    ID: 1,
                    Title: "list1",
                    ListTasks: [
                      { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                      { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                    ]
                  },
                  fromTask: {
                    ID: 2,
                    ListID: 1,
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
                  ID: 3,
                  ListTasks: [
                    { ID: 3, ListID: 3, Title: "task3", Description: "hogehoge" },
                    { draggedOn: true },
                    { ID: 4, ListID: 3, Title: "task4", Description: "hogehoge" }
                  ]
                },
                prevToTask: { ID: 4, ListID: 3, Title: "task4", Description: "hogehoge" }
              },
              lists: [{
                ID: 1,
                Title: "list1",
                ListTasks: [
                  { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                  { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                ]
              }, {
                ID: 2,
                Title: "list2",
                ListTasks: []
              }],
              noneList: {
                ID: 3,
                ListTasks: [
                  { ID: 3, ListID: 3, Title: "task3", Description: "hogehoge" },
                  { draggedOn: true },
                  { ID: 4, ListID: 3, Title: "task4", Description: "hogehoge" }
                ]
              },
              taskDraggingFrom: {
                fromList: {
                  ID: 1,
                  Title: "list1",
                  ListTasks: [
                    { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                    { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                  ]
                },
                fromTask: {
                  ID: 2,
                  ListID: 1,
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
                  ID: 1,
                  Title: "list1",
                  ListTasks: [
                    { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                    { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                  ]
                }, {
                  ID: 2,
                  Title: "list2",
                  ListTasks: [
                    { ID: 3, ListID: 2, Title: "task3", Description: "hogehoge" },
                    { ID: 4, ListID: 2, Title: "task4", Description: "hogehoge" }]
                }],
                noneList: {ID: 0, ListTasks: [] },
                taskDraggingTo: null,
                taskDraggingFrom: {
                  fromList: {
                    ID: 1,
                    Title: "list1",
                    ListTasks: [
                      { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                      { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                    ]
                  },
                  fromTask: {
                    ID: 2,
                    ListID: 1,
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
                  ID: 2,
                  Title: "list2",
                  ListTasks: [
                    { ID: 3, ListID: 2, Title: "task3", Description: "hogehoge" },
                    { draggedOn: true },
                    { ID: 4, ListID: 2, Title: "task4", Description: "hogehoge" }
                  ]
                },
                prevToTask: { ID: 4, ListID: 2, Title: "task4", Description: "hogehoge" }
              },
              lists: [{
                ID: 1,
                Title: "list1",
                ListTasks: [
                  { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                  { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                ]
              }, {
                ID: 2,
                Title: "list2",
                ListTasks: [
                  { ID: 3, ListID: 2, Title: "task3", Description: "hogehoge" },
                  { draggedOn: true },
                  { ID: 4, ListID: 2, Title: "task4", Description: "hogehoge" }
                ]
              }],
              noneList: {ID: 0, ListTasks: [] },
              taskDraggingFrom: {
                fromList: {
                  ID: 1,
                  Title: "list1",
                  ListTasks: [
                    { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                    { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                  ]
                },
                fromTask: {
                  ID: 2,
                  ListID: 1,
                  Title: "task2",
                  Description: "hogehoge"
                }
              }
            })
          })
        })
      })
    })

    describe('TASK_DRAG_START', () => {
      context('when drag from list is noneList', () => {
        it('should return updated lists', () => {
          expect(
            ListReducer({
              lists: [{
                ID: 1,
                Title: "list1",
                ListTasks: [
                  { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                  { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                ]
              }, {
                ID: 2,
                Title: "list2",
                ListTasks: []
              }],
              noneList: {
                ID: 3,
                ListTasks: [
                  { ID: 3, ListID: 3, Title: "task3", Description: "hogehoge" }
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
              ID: 1,
              Title: "list1",
              ListTasks: [
                { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              ID: 2,
              Title: "list2",
              ListTasks: []
            }],
            noneList: {
              ID: 3,
              ListTasks: [
                { ID: 3, ListID: 3, Title: "task3", Description: "hogehoge" }
              ]
            },
            taskDraggingFrom: {
              fromList: {
                ID: 3,
                ListTasks: [
                  { ID: 3, ListID: 3, Title: "task3", Description: "hogehoge" }
                ]
              },
              fromTask: {
                ID: 3,
                ListID: 3,
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
                ID: 1,
                Title: "list1",
                ListTasks: [
                  { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                  { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                ]
              }, {
                ID: 2,
                Title: "list2",
                ListTasks: []
              }],
              noneList: {ID: 3, ListTasks: [] },
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
              ID: 1,
              Title: "list1",
              ListTasks: [
                { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              ID: 2,
              Title: "list2",
              ListTasks: []
            }],
            noneList: {ID: 3, ListTasks: [] },
            taskDraggingFrom: {
              fromList: {
                ID: 1,
                Title: "list1",
                ListTasks: [
                  { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                  { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                ]
              },
              fromTask: {
                ID: 2,
                ListID: 1,
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
                ID: 1,
                Title: "list1",
                ListTasks: [
                  { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                  { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                ]
              }, {
                ID: 2,
                Title: "list2",
                ListTasks: []
              }],
              noneList: {
                ID: 3,
                ListTasks: [
                  { draggedOn: true }
                ]
              },
              isTaskDraggingOver: true,
              taskDraggingTo: {
                toList: {
                  ID: 3,
                  ListTasks: []
                },
                prevToTask: null
              }
            }, {
              type: listActions.TASK_DRAG_LEAVE
            })
          ).toEqual({
            isTaskDraggingOver: false,
            taskDraggingTo: null,
            lists: [{
              ID: 1,
              Title: "list1",
              ListTasks: [
                { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              ID: 2,
              Title: "list2",
              ListTasks: []
            }],
            noneList: {
              ID: 3,
              ListTasks: []
            }
          })
        })
      })
      context('when target list is not noneList', () => {
        it('should return lists do not contain arrow', () => {
          expect(
            ListReducer({
              lists: [{
                ID: 1,
                Title: "list1",
                ListTasks: [
                  { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                  { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                ]
              }, {
                ID: 2,
                Title: "list2",
                ListTasks: [ { draggedOn: true } ]
              }],
              noneList: {ID: 0, ListTasks: [] },
              isTaskDraggingOver: true,
              taskDraggingTo: {
                toList: {
                  ID: 2,
                  Title: "list2",
                  ListTasks: []
                },
                prevToTask: null
              }
            }, {
              type: listActions.TASK_DRAG_LEAVE
            })
          ).toEqual({
            isTaskDraggingOver: false,
            taskDraggingTo: null,
            lists: [{
              ID: 1,
              Title: "list1",
              ListTasks: [
                { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              ID: 2,
              Title: "list2",
              ListTasks: []
            }],
            noneList: {ID: 0, ListTasks: [] }
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
                ID: 1,
                Title: "list1",
                ListTasks: [
                  { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                  { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                ]
              }, {
                ID: 2,
                Title: "list2",
                ListTasks: []
              }],
              noneList: {
                ID: 3,
                ListTasks: [ {draggedOn: true } ]
              },
              isTaskDraggingOver: true,
              taskDraggingTo: {
                toList: {
                  ID: 3,
                  ListTasks: []
                },
                prevToTask: null
              },
              taskDraggingFrom: {
                fromList: {
                  ID: 1,
                  Title: "list1",
                  ListTasks: [
                    { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                    { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                  ]
                },
                fromTask: {
                  ID: 2,
                  ListID: 1,
                  Title: "task2",
                  Description: "hogehoge"
                }
              }
            }, {
              type: listActions.REQUEST_MOVE_TASK
            })
          ).toEqual({
            isTaskDraggingOver: false,
            taskDraggingTo: null,
            taskDraggingFrom: null,
            lists: [{
              ID: 1,
              Title: "list1",
              ListTasks: [
                { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
              ],
              isLoading: true
            }, {
              ID: 2,
              Title: "list2",
              ListTasks: []
            }],
            noneList: {
              ID: 3,
              ListTasks: []
            }
          })
        })
      })
      context('when target list is not noneList', () => {
        it('should return do not contain arrow and contain isLoading flag', () => {
          expect(
            ListReducer({
              lists: [{
                ID: 1,
                Title: "list1",
                ListTasks: [
                  { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                  { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                ]
              }, {
                ID: 2,
                Title: "list2",
                ListTasks: [ { draggedOn: true } ]
              }],
              noneList: {ID: 0, ListTasks: [] },
              isTaskDraggingOver: true,
              taskDraggingTo: {
                toList: {
                  ID: 2,
                  Title: "list2",
                  ListTasks: []
                },
                prevToTask: null
              },
              taskDraggingFrom: {
                fromList: {
                  ID: 1,
                  Title: "list1",
                  ListTasks: [
                    { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                    { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                  ]
                },
                fromTask: {
                  ID: 2,
                  ListID: 1,
                  Title: "task2",
                  Description: "hogehoge"
                }
              }
            }, {
              type: listActions.REQUEST_MOVE_TASK
            })
          ).toEqual({
            isTaskDraggingOver: false,
            taskDraggingTo: null,
            taskDraggingFrom: null,
            lists: [{
              ID: 1,
              Title: "list1",
              ListTasks: [
                { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
              ],
              isLoading: true
            }, {
              ID: 2,
              Title: "list2",
              ListTasks: [],
              isLoading: true
            }],
            noneList: {ID: 0, ListTasks: [] }
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
                ID: 1,
                Title: "list1",
                ListTasks: [
                  { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                  { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                ]
              }, {
                ID: 2,
                Title: "list2",
                ListTasks: []
              }],
              noneList: {
                ID: 3,
                ListTasks: [ { draggedOn: true } ]
              },
              isTaskDraggingOver: true,
              taskDraggingTo: {
                toList: {
                  ID: 3,
                  ListTasks: []
                },
                prevToTask: null
              },
              taskDraggingFrom: {
                fromList: {
                  ID: 1,
                  Title: "list1",
                  ListTasks: [
                    { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                    { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                  ]
                },
                fromTask: {
                  ID: 2,
                  ListID: 1,
                  Title: "task2",
                  Description: "hogehoge"
                }
              }
            }, {
              type: listActions.TASK_DROP
            })
          ).toEqual({
            isTaskDraggingOver: false,
            taskDraggingTo: null,
            taskDraggingFrom: null,
            lists: [{
              ID: 1,
              Title: "list1",
              ListTasks: [
                { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              ID: 2,
              Title: "list2",
              ListTasks: []
            }],
            noneList: {
              ID: 3,
              ListTasks: []
            }
          })
        })
      })
      context('when target list is not noneList', () => {
        it('should return do not contain arrow and isLoading', () => {
          expect(
            ListReducer({
              lists: [{
                ID: 1,
                Title: "list1",
                ListTasks: [
                  { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                  { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                ]
              }, {
                ID: 2,
                Title: "list2",
                ListTasks: [ { draggedOn: true } ]
              }],
              noneList: {ID: 0, ListTasks: [] },
              isTaskDraggingOver: true,
              taskDraggingTo: {
                toList: {
                  ID: 2,
                  Title: "list2",
                  ListTasks: []
                },
                prevToTask: null
              },
              taskDraggingFrom: {
                fromList: {
                  ID: 1,
                  Title: "list1",
                  ListTasks: [
                    { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                    { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
                  ]
                },
                fromTask: {
                  ID: 2,
                  ListID: 1,
                  Title: "task2",
                  Description: "hogehoge"
                }
              }
            }, {
              type: listActions.TASK_DROP
            })
          ).toEqual({
            isTaskDraggingOver: false,
            taskDraggingTo: null,
            taskDraggingFrom: null,
            lists: [{
              ID: 1,
              Title: "list1",
              ListTasks: [
                { ID: 1, ListID: 1, Title: "task1", Description: "hogehoge" },
                { ID: 2, ListID: 1, Title: "task2", Description: "hogehoge" }
              ]
            }, {
              ID: 2,
              Title: "list2",
              ListTasks: []
            }],
            noneList: {ID: 0, ListTasks: [] }
          })
        })
      })
    })
  })



  context('newTaskModalActions', () => {
    describe('CLOSE_NEW_TASK', () => {
      it('should close new task modal', () => {
        expect(
          ListReducer({
            isTaskModalOpen: true,
            selectedList: "newList"
          }, {
            type: newTaskModalActions.CLOSE_NEW_TASK,
            isTaskModalOpen: false
          })
        ).toEqual({
          isTaskModalOpen: false,
          selectedList: null
        })
      })
    })
    describe('UPDATE_NEW_TASK_TITLE', () => {
      it('should update new task title', () => {
        expect(
          ListReducer({
            newTask: { title: "" }
          }, {
            type: newTaskModalActions.UPDATE_NEW_TASK_TITLE,
            title: "newTask"
          })
        ).toEqual({
          newTask: { title: "newTask" }
        })
      })
    })

    describe('REQUEST_CREATE_TASK', () => {
      it('should open whole loading window', () => {
        expect(
          ListReducer({
            isLoading: false
          }, {
            type: newTaskModalActions.REQUEST_CREATE_TASK
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
              ID: 1,
              Title: "list1",
              ListTasks: [
                { ListID: 1, Title: "task1", Description: "fugafuga" }
              ]
            }, {
              ID: 2,
              Title: "list2",
              ListTasks: []
            }],
            noneList: {ID: 0, ListTasks: [] },
            newTask: { title: "task2", description: "hogehoge" },
            isTaskModalOpen: true,
            isLoading: true
          }, {
            type: newTaskModalActions.RECEIVE_CREATE_TASK,
            task: { ListID: 1, Title: "task2", Description: "hogehoge" }
          })
        ).toEqual({
          lists: [{
            ID: 1,
            Title: "list1",
            ListTasks: [
              { ListID: 1, Title: "task1", Description: "fugafuga" },
              { ListID: 1, Title: "task2", Description: "hogehoge" }
            ]
          }, {
            ID: 2,
            Title: "list2",
            ListTasks: []
          }],
          noneList: {ID: 0, ListTasks: [] },
          newTask: { title: "", description: "" },
          isTaskModalOpen: false,
          isLoading: false
        })
      })
    })
  })

  context('editListModalActions', () => {
    describe('CLOSE_EDIT_LIST', () => {
      it('should close edit list modal', () => {
        expect(
          ListReducer({
            isListEditModalOpen: true,
            selectedList: "editList"
          }, {
            type: editListModalActions.CLOSE_EDIT_LIST,
            isListEditModalOpen: false
          })
        ).toEqual({
          isListEditModalOpen: false,
          selectedList: null,
          selectedListOption: null
        })
      })
    })
    describe('UPDATE_SELECTED_LIST_TITLE', () => {
      it('should update selected list title', () => {
        expect(
          ListReducer({
            selectedList: { Title: "selectedL", Color: "" }
          }, {
            type: editListModalActions.UPDATE_SELECTED_LIST_TITLE,
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
            type: editListModalActions.UPDATE_SELECTED_LIST_COLOR,
            color: "30bef"
          })
        ).toEqual({
          selectedList: { Title: "selectedList", Color: "30bef" }
        })
      })
    })

    describe('REQUEST_UPDATE_LIST', () => {
      it('should open whole loading window', () => {
        expect(
          ListReducer({
            isLoading: false
          }, {
            type: editListModalActions.REQUEST_UPDATE_LIST
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
              ID: 1,
              Title: "list1",
              ListTasks: [
                { ListID: 1, Title: "task1", Description: "hogehoge" }
              ]
            }, {
              ID: 2,
              Title: "list2",
              ListTasks: []
            }],
            isListEditModalOpen: true,
            isLoading: true
          }, {
            type: editListModalActions.RECEIVE_UPDATE_LIST,
            list: {
              ID: 1,
              Title: "updateList1",
              ListTasks: [
                { ListID: 1, Title: "task1", Description: "hogehoge" }
              ]
            }
          })
        ).toEqual({
          lists: [{
            ID: 1,
            Title: "updateList1",
            ListTasks: [
              { ListID: 1, Title: "task1", Description: "hogehoge" }
            ]
          }, {
            ID: 2,
            Title: "list2",
            ListTasks: []
          }],
          isListEditModalOpen: false,
          isLoading: false
        })
      })
    })

    describe('CHANGE_SELECTED_LIST_OPTION', () => {
      it('should change selected list option', () => {
        expect(
          ListReducer({
            listOptions: [
              {
                ID: 1,
                Action: "close",
              }, {
                ID: 2,
                Action: "open"
              }
            ],
            selectedListOption: null
          }, {
            type: editListModalActions.CHANGE_SELECTED_LIST_OPTION,
            selectEvent: {
              value: 1
            }
          })
        ).toEqual({
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
          }
        })
      })
    })
  })

  context('newListModalActions', () => {
    describe('UPDATE_NEW_LIST_TITLE', () => {
      it('should update list title', () => {
        expect(
          ListReducer({
            newList: { title: "newL", color: "" }
          }, {
            type: newListModalActions.UPDATE_NEW_LIST_TITLE,
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
            type: newListModalActions.UPDATE_NEW_LIST_COLOR,
            color: "30bfe"
          })
        ).toEqual({
          newList: { title: "newList", color: "30bfe" }
        })
      })
    })
    describe('REQUEST_CREATE_LIST', () => {
      it('should open whole loading window', () => {
        expect(
          ListReducer({
            isLoading: false
          }, {
            type: newListModalActions.REQUEST_CREATE_LIST
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
              type: newListModalActions.RECEIVE_CREATE_LIST,
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
              type: newListModalActions.RECEIVE_CREATE_LIST,
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
  })
})
