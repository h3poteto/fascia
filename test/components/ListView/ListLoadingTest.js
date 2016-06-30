import TestUtils from 'react-addons-test-utils'
import expect from 'expect'
import React from 'react'
import ListLoading from '../../../frontend/javascripts/components/ListView/ListLoading.jsx'
//import * as  from '../../fixtures/components/ListView/EditProjectModalFixture'

function setup(props) {
  let renderer = TestUtils.createRenderer()
  renderer.render(<ListLoading {...props} />)
  let output = renderer.getRenderOutput()

  return {
    props,
    output,
    renderer
  }
}

describe('ListView::ListLoading', () => {
  context('when list is loading', () => {
    let state = {
      isLoading: true
    }
    it('should render loading', () => {
      const { output } = setup(state)
      expect(output.type).toBe('div')
      expect(output.props.className).toBe('list-loading')
    })
  })
  context('when list is not loading', () => {
    let state = {
      isLoading: false
    }
    it('should not render loading', () => {
      const { output } = setup(state)
      expect(output.type).toBe('span')
    })
  })
  context('when loading is not defined', () => {
    let state = {}
    it('should not render loading', () => {
      const { output } = setup(state)
      expect(output.type).toBe('span')
    })
  })
})
