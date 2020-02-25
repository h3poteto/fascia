import React from 'react'
import { Navbar, NavDropdown, Nav } from 'react-bootstrap'

type Props = {}

const menu: React.FC<Props> = ({ children }) => {
  return (
    <div>
      <header className="top-header">
        <Navbar bg="primary" variant="dark" expand="lg">
          <Navbar.Brand href="/">Fascia</Navbar.Brand>
          <Navbar.Toggle aria-controls="basic-navbar-nav" />
          <Navbar.Collapse id="basic-navbar-nav">
            <Nav className="mr-auto">
              <Nav.Link href="#home">Projects</Nav.Link>
              <Nav.Link href="#link">Contact</Nav.Link>
              <NavDropdown title="Accounts" id="basic-nav-dropdown">
                <NavDropdown.Item href="#action/3.1">Settings</NavDropdown.Item>
                <NavDropdown.Divider />
                <NavDropdown.Item href="#action/3.4">Logout</NavDropdown.Item>
              </NavDropdown>
            </Nav>
          </Navbar.Collapse>
        </Navbar>
      </header>
      <div className="contents">
        {children}
      </div>
    </div>
  )
}

export default menu
