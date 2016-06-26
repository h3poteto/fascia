import Request from 'superagent'

export const UNAUTHORIZED = 'UNAUTHORIZED'
function unauthorized() {
  window.location.pathname = "/sign_in"
  return {
    type: UNAUTHORIZED
  }
}

export const NOT_FOUND = 'NOT_FOUND'
function notFound() {
  return {
    type: NOT_FOUND
  }
}

export const SERVER_ERROR = 'SERVER_ERROR'
function serverError() {
  return {
    type: SERVER_ERROR
  }
}

export const CLOSE_EDIT_LIST = 'CLOSE_EDIT_LIST'
export function closeEditListModal() {
  return {
    type: CLOSE_EDIT_LIST,
    isListEditModalOpen: false
  }
}

export const UPDATE_SELECTED_LIST_TITLE = 'UPDATE_SELECTED_LIST_TITLE'
export function updateSelectedListTitle(ev) {
  return {
    type: UPDATE_SELECTED_LIST_TITLE,
    title: ev.target.value
  }
}

export const UPDATE_SELECTED_LIST_COLOR = 'UPDATE_SELECTED_LIST_COLOR'
export function updateSelectedListColor(ev) {
  return {
    type: UPDATE_SELECTED_LIST_COLOR,
    color: ev.target.value
  }
}

export const REQUEST_UPDATE_LIST = 'REQUEST_UPDATE_LIST'
function requestUpdateList() {
  return {
    type: REQUEST_UPDATE_LIST
  }
}

export const RECEIVE_UPDATE_LIST = 'RECEIVE_UPDATE_LIST'
function receiveUpdateList(list) {
  return {
    type: RECEIVE_UPDATE_LIST,
    list: {ID: list.ID, ProjectID: list.ProjectID, Title: list.Title, Color: list.Color, ListTasks: list.ListTasks, ListOptionID: list.ListOptionID}
  }
}

export function fetchUpdateList(projectID, list, option) {
  var action
  if (option != undefined && option != null) {
    action = option.Action
  }
  return dispatch => {
    dispatch(requestUpdateList())
    return Request
      .post(`/projects/${projectID}/lists/${list.ID}`)
      .type('form')
      .send({title: list.Title, color: list.Color, action: action})
      .end((err, res)=> {
        if(res.ok) {
          dispatch(receiveUpdateList(res.body))
        } else if (res.unauthorized) {
          dispatch(unauthorized())
        } else if (res.notFound) {
          dispatch(notFound())
        } else {
          dispatch(serverError())
        }
      })
  }
}


export const CHANGE_SELECTED_LIST_OPTION = 'CHANGE_SELECTED_LIST_OPTION'
export function changeSelectedListOption(ev) {
  return {
    type: CHANGE_SELECTED_LIST_OPTION,
    selectEvent: ev.target
  }
}
