import expect from 'expect'

export function initState() {
  return {
    isListModalOpen: false,
    newList: {title: "", color: "0effff"},
    projectID: 1,
  }
}

export function openNewListModalState() {
  let state = initState()
  state["isListModalOpen"] = true
  return state
}
