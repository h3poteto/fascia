import expect from 'expect'

export function initState() {
  return {
    isTaskModalOpen: false,
    newTask: {title: ""},
    selectedList: 1,
    projectID: 1,
  }
}

export function openNewTaskModalState() {
  let state = initState()
  state["isTaskModalOpen"] = true
  return state
}
