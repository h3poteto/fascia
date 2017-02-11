import TestUtils from 'react-addons-test-utils'
import expect from 'expect'
import React from 'react'
import ListView from '../../frontend/javascripts/components/ListView.jsx'
import * as ListViewFixture from '../fixtures/components/ListViewFixture'

function setup(props) {
  let renderer = TestUtils.createRenderer()
  renderer.render(<ListView {...props} />)
  let output = renderer.getRenderOutput()

  return {
    props,
    output,
    renderer
  }
}

describe('ListView', () => {
  context('when no error', () => {
    let state = ListViewFixture.initState()
    it('should render correctly', () => {
      const { output, props } = setup(state)

      expect(output.type).toBe('div')
      expect(output.props.id).toBe('lists')

      let [ , , , , , , , , projectTitleWrapper, items, noneList ] = output.props.children
      let [ , projectTitle ] = projectTitleWrapper.props.children
      let [ title ] = projectTitle.props.children
      expect(title).toBe('testProject')

      expect(items.props.className).toBe('items')

      let [ list, button ] = items.props.children
      expect(list[0].props['data-id']).toBe(1)
      expect(list[1].props['data-id']).toBe(2)

      list[0].props.onDrop()
      expect(props.listActions.taskDrop.calls.length).toBe(1)

      button.props.onClick()
      expect(props.listActions.openNewListModal.calls.length).toBe(1)

      // list which have tasks
      let [ , , firstTasks ] = list[0].props.children
      let [ firstListTasks ] = firstTasks.props.children
      let [ task1, task2 ] = firstListTasks
      expect(task1.props['data-id']).toBe(1)
      expect(task1.props.children).toBe('task1')
      expect(task2.props['data-id']).toBe(2)
      expect(task2.props.children).toBe('task2')

      task2.props.onDragStart()
      expect(props.listActions.taskDragStart.calls.length).toBe(1)

      // list which do not have tasks
      let [ , , secondTasks ] = list[1].props.children
      let [ , secondNewTask ] = secondTasks.props.children
      expect(secondNewTask.props.className).toBe('new-task')

      secondNewTask.props.onClick()
      expect(props.listActions.openNewTaskModal.calls.length).toBe(1)

      let [ tasks, newTask ] = noneList.props.children.props.children
      expect(tasks[0].props['data-id']).toBe(3)
      expect(newTask.props.className).toBe('new-task pure-button button-blue')

      newTask.props.onClick()
      expect(props.listActions.openNewTaskModal.calls.length).toBe(2)
    })
  })

  context('when one error, not modal', () => {
    let state = ListViewFixture.errorState()
    it('should render error', () => {
      const { output } = setup(state)
      let [ , flash, , , , , , ] = output.props.children
      expect(flash.props.children).toBe('Server Error')
    })
  })


  context('when showIssue is false', () => {
    let state = ListViewFixture.hideIssueState()
    it('should not render issues', () => {
      const { output } = setup(state)

      let [ , , , , , , , , , items, noneList ] = output.props.children
      let [list ] = items.props.children
      let [ , , firstTasks ] = list[0].props.children
      let [ firstListTasks ] = firstTasks.props.children
      let [ task1, task2 ] = firstListTasks
      expect(task1.props['data-id']).toBe(1)
      expect(task2).toBe(undefined)
      let [ tasks ] = noneList.props.children.props.children
      expect(tasks[0]).toBe(undefined)
    })
  })

  context('when showPullRequest is false', () => {
    let state = ListViewFixture.showIssueState()
    it('should not render pull requests', () => {
      const { output } = setup(state)

      let [ , , , , , , , , , items, noneList ] = output.props.children
      let [list ] = items.props.children
      let [ , , firstTasks ] = list[0].props.children
      let [ firstListTasks ] = firstTasks.props.children
      let [ task1, task2 ] = firstListTasks
      expect(task1).toBe(undefined)
      expect(task2.props['data-id']).toBe(2)
      let [ tasks ] = noneList.props.children.props.children
      expect(tasks[0].props['data-id']).toBe(3)
    })
  })

  describe('github action buttons', () => {
    context('when project does not have repository', () => {
      let state = ListViewFixture.noRepositoryProjectState()
      it('should not render github action buttons', () => {
        const { output } = setup(state)

        let [ , , , , , , , , projectTitleWrapper ] = output.props.children
        let [ operation ] = projectTitleWrapper.props.children
        expect(operation.props.children.props.children).toBe(undefined)
      })
    })
    context('when project has repository', () => {
      let state = ListViewFixture.repositoryProjectState()
      it('should not render github action buttons', () => {
        const { output } = setup(state)

        let [ , , , , , , , , projectTitleWrapper ] = output.props.children
        let [ operation ] = projectTitleWrapper.props.children
        expect(operation.props.children.props.children.length).toBe(3)
      })
    })
  })

  describe('hide and display', () => {
    context('when a list is hidden', () => {
      let state = ListViewFixture.hiddenListState()
      it('should hide a list', () => {
        const { output } = setup(state)

        let [ , , , , , , , , , items ] = output.props.children
        let [ list ] = items.props.children
        let [ listMenu, listTitle, listLoading ] = list[0].props.children

        expect(listMenu.props.className).toBe('fascia-list-menu')
        expect(listTitle.props.className).toBe('list-title')
        expect(listLoading.props.type).toNotBe('ul')
      })
    })
    context('when a list is displeyd', () => {
      let state = ListViewFixture.initState()
      it('should display a list', () => {
        const { output } = setup(state)

        let [ , , , , , , , , , items ] = output.props.children
        let [ list ] = items.props.children
        let [ listMenu, listTitle, tasks ] = list[0].props.children
        expect(listMenu.props.className).toBe('fascia-list-menu')
        expect(listTitle.props.className).toBe('list-title')
        expect(tasks.type).toBe('ul')
        expect(tasks.props.children.length).toNotBe(0)
      })
    })
  })
})
