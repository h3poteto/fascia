import TestUtils from 'react-addons-test-utils'
import expect from 'expect'
import React from 'react'
import WholeLoading from '../../../frontend/javascripts/components/ListView/WholeLoading.jsx'

function setup(props) {
  let renderer = TestUtils.createRenderer()
  renderer.render(<WholeLoading {...props} />)
  let output = renderer.getRenderOutput()

  return {
    props,
    output,
    renderer
  }
}

describe('ListView::WholeLoading', () => {
  context('when whole loading is open', () => {
    let state = {
      isLoading: true
    }
    it('should render loading window', () => {
      const { output } = setup(state)

      expect(output.type).toBe('div')
      expect(output.props.className).toBe('whole-loading')
    })
  })
})
