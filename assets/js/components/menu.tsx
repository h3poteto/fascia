import React from 'react'
import { Navbar, NavDropdown, Nav, Container } from 'react-bootstrap'
import { ThunkDispatch } from 'redux-thunk'
import { Link } from 'react-router-dom'

import Actions, { logout } from '@/actions/menu'
import styles from './menu.scss'

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
            <Container className={styles.container}>
              <Navbar.Brand href="/">
                <img alt="" src="/lp/images/fascia-icon.png" width="30" height="30" className="d-inline-block align-top" /> Fascia
              </Navbar.Brand>
              <Navbar.Toggle aria-controls="basic-navbar-nav" />
              <Navbar.Collapse id="basic-navbar-nav">
                <Nav className="me-auto">
                  <Nav.Link href="/">Projects</Nav.Link>
                  <Nav.Link href="/inquiries/new">Contact</Nav.Link>
                </Nav>
              </Navbar.Collapse>
              <Navbar.Collapse className="justify-content-end">
                <Nav>
                  <NavDropdown title="Accounts" id="basic-nav-dropdown" className={styles.accounts}>
                    <NavDropdown.Item>
                      <Link to="/settings">Settings</Link>
                    </NavDropdown.Item>
                    <NavDropdown.Divider />
                    <NavDropdown.Item onClick={handleLogout}>Logout</NavDropdown.Item>
                  </NavDropdown>
                </Nav>
              </Navbar.Collapse>
            </Container>
          </Navbar>
        </header>
        <div className="contents">{this.props.children}</div>
      </div>
    )
  }
}

export default Menu
