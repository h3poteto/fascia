import axios from "axios";
import { ErrorHandler, ServerError } from "../ErrorHandler";
import { startLoading, stopLoading } from "../Loading";

export const CLOSE_NEW_TASK = "CLOSE_NEW_TASK";
export function closeNewTaskModal() {
  return {
    type: CLOSE_NEW_TASK,
    isTaskModalOpen: false
  };
}

export const REQUEST_CREATE_TASK = "REQUEST_CREATE_TASK";
function requestCreateTask() {
  return {
    type: REQUEST_CREATE_TASK
  };
}

export const RECEIVE_CREATE_TASK = "RECEIVE_CREATE_TASK";
function receiveCreateTask(lists) {
  return {
    type: RECEIVE_CREATE_TASK,
    lists: lists.Lists,
    noneList: lists.NoneList
  };
}

export function fetchCreateTask(params) {
  return (dispatch, getState) => {
    const {
      ListReducer: {
        project: { ID: projectID }
      }
    } = getState();
    const {
      ListReducer: {
        selectedList: { ID: listID }
      }
    } = getState();
    dispatch(startLoading());
    dispatch(requestCreateTask());
    return axios
      .post(`/api/projects/${projectID}/lists/${listID}/tasks`, params)
      .then(res => {
        dispatch(stopLoading());
        dispatch(receiveCreateTask(res.data));
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
