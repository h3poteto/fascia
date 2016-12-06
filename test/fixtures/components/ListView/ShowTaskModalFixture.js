import expect from 'expect'

export function initState() {
  return {
    isShowTaskModalOpen: false,
    isEditTaskModalVisible: false,
    task: {Title: "", Description: "", IssueNumber: 0},
    editTask: { Title: "", Description: ""},
  }
}

export function openShowTaskModalState() {
  let state = initState()
  state["isShowTaskModalOpen"] = true
  return state
}

