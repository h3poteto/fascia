import axios from "axios";
import { ErrorHandler, ServerError } from "../ErrorHandler";
import { startLoading, stopLoading } from "../Loading";

export const CLOSE_SHOW_TASK = "CLOSE_SHOW_TASK";
export function closeShowTaskModal() {
  return {
    type: CLOSE_SHOW_TASK
  };
}

export const CHANGE_EDIT_MODE = "CHANGE_EDIT_MODE";
export function changeEditMode(task) {
  return {
    type: CHANGE_EDIT_MODE,
    task: task
  };
}

export const REQUEST_UPDATE_TASK = "REQUEST_UPDATE_TASK";
function requestUpdateTask() {
  return {
    type: REQUEST_UPDATE_TASK
  };
}

export const RECEIVE_UPDATE_TASK = "RECEIVE_UPDATE_TASK";
function receiveUpdateTask(lists) {
  return {
    type: RECEIVE_UPDATE_TASK,
    lists: lists.Lists,
    noneList: lists.NoneList
  };
}

export function fetchUpdateTask(params) {
  return (dispatch, getState) => {
    const {
      ListReducer: {
        project: { ID: projectID }
      }
    } = getState();
    const {
      ListReducer: {
        selectedTask: { ListID: listID }
      }
    } = getState();
    const {
      ListReducer: {
        selectedTask: { ID: taskID }
      }
    } = getState();
    dispatch(startLoading());
    dispatch(requestUpdateTask());
    return axios
      .patch(
        `/api/projects/${projectID}/lists/${listID}/tasks/${taskID}`,
        params
      )
      .then(res => {
        dispatch(stopLoading());
        dispatch(receiveUpdateTask(res.data));
      })
      .catch(err => {
        dispatch(stopLoading());
        ErrorHandler(err)
          .then()
          .catch(error => {
            dispatch(ServerError(error));
          });
      });
  };
}

export const REQUEST_DELETE_TASK = "REQUEST_DELETE_TASK";
function requestDeleteTask() {
  return {
    type: REQUEST_DELETE_TASK
  };
}

export const RECEIVE_DELETE_TASK = "RECEIVE_DELETE_TASK";
function receiveDeleteTask(lists) {
  return {
    type: RECEIVE_DELETE_TASK,
    lists: lists.Lists,
    noneList: lists.NoneList
  };
}

export function fetchDeleteTask() {
  return (dispatch, getState) => {
    const {
      ListReducer: {
        project: { ID: projectID }
      }
    } = getState();
    const {
      ListReducer: {
        selectedTask: { ListID: listID }
      }
    } = getState();
    const {
      ListReducer: {
        selectedTask: { ID: taskID }
      }
    } = getState();
    dispatch(startLoading());
    dispatch(requestDeleteTask());
    return axios
      .delete(`/api/projects/${projectID}/lists/${listID}/tasks/${taskID}`)
      .then(res => {
        dispatch(stopLoading());
        dispatch(receiveDeleteTask(res.data));
      })
      .catch(err => {
        dispatch(stopLoading());
        ErrorHandler(err)
          .then()
          .catch(error => {
            dispatch(ServerError(error));
          });
      });
  };
}
