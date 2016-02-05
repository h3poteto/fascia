import TestUtils from 'react-addons-test-utils'
import expect from 'expect'
import React from 'react'
import MenuView from '../../frontend/javascripts/components/MenuView.jsx'

function setup(props) {
  let renderer = TestUtils.createRenderer()
  renderer.render(<MenuView {...props} />)
  let output = renderer.getRenderOutput()

  return {
    props,
    output,
    renderer
  }
}

describe('MenuView', () => {
  it('should render correctly', () => {
    let state = {
      logout: expect.createSpy()
    }
    const { output, props } = setup(state)

    expect(output.type).toBe('div')

    let [ header, mainBoard ] = output.props.children
    expect(header.props.className).toBe('top-nav')
    expect(mainBoard.props.id).toBe('main_board')

    let menuHorizontal = header.props.children
    let [ heading, menuList, controlList ] = menuHorizontal.props.children
    let [ settings, settingsList ] = controlList.props.children.props.children
    let logout = settingsList.props.children
    let logoutLink = logout.props.children.props.children

    expect(logout.props.children.props.action).toBe("/sign_out")
    logoutLink.props.onClick()
    expect(props.logout.calls.length).toBe(1)
  })
})
