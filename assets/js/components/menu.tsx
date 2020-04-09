import React from 'react'
import { Navbar, NavDropdown, Nav } from 'react-bootstrap'
import { ThunkDispatch } from 'redux-thunk'

import Actions, { logout } from '@/actions/menu'

type Props = {
  dispatch: ThunkDispatch<any, any, Actions>
}

class Menu extends React.Component<Props> {
  render() {
    const handleLogout = () => {
      this.props.dispatch(logout())
    }

    return (
      <div>
        <header className="top-header">
          <Navbar bg="primary" variant="dark" expand="lg">
            <Navbar.Brand href="/">
              <img
                alt=""
                src="/images/fascia-icon.png"
                width="30"
                height="30"
                className="d-inline-block align-top"
              />{' '}
              Fascia</Navbar.Brand>
              <Navbar.Toggle aria-controls="basic-navbar-nav" />
              <Navbar.Collapse id="basic-navbar-nav">
                <Nav className="mr-auto">
                  <Nav.Link href="/">Projects</Nav.Link>
                  <Nav.Link href="/inquiries/new">Contact</Nav.Link>
                </Nav>
                <Nav>
                  <NavDropdown title="Accounts" id="basic-nav-dropdown">
                    <NavDropdown.Item href="#action/settings">Settings</NavDropdown.Item>
                    <NavDropdown.Divider />
                    <NavDropdown.Item href="#" onClick={handleLogout} >Logout</NavDropdown.Item>
                  </NavDropdown>
                </Nav>
              </Navbar.Collapse>
          </Navbar>
        </header>
        <div className="contents">
          {this.props.children}
        </div>
      </div>
    )
  }
}

export default Menu
