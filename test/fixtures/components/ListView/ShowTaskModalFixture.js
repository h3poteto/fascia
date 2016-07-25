import expect from 'expect'

export function initState() {
  return {
    isShowTaskModalOpen: false,
    isEditTaskModalVisible: false,
    task: {Title: "", Description: "", IssueNumber: 0},
    editTask: { Title: "", Description: ""},
    closeShowTaskModal: expect.createSpy()
  }
}

export function openShowTaskModalState() {
  let state = initState()
  state["isShowTaskModalOpen"] = true
  return state
}

export function visibleEditTaskModalState() {
  let state = openShowTaskModalState()
  state["isEditTaskModalVisible"] = true
  return state
}
