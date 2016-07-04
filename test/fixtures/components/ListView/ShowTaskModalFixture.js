import expect from 'expect'

export function initState() {
  return {
    isShowTaskModalOpen: false,
    task: {Title: "", Description: "", IssueNumber: 0},
    closeShowTaskModal: expect.createSpy()
  }
}

export function openShowTaskModalState() {
  let state = initState()
  state["isShowTaskModalOpen"] = true
  return state
}
