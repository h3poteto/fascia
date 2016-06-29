import expect from 'expect'

export function initState() {
  return {
    isListModalOpen: false,
    newList: {title: "", color: "0effff"},
    projectID: 1,
    closeNewListModal: expect.createSpy(),
    updateNewListTitle: expect.createSpy(),
    updateNewListColor: expect.createSpy(),
    fetchCreateList: expect.createSpy
  }
}

export function openNewListModalState() {
  let state = initState()
  state["isListModalOpen"] = true
  return state
}
