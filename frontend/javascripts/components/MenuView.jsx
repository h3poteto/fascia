import React from 'react'
import PropTypes from 'prop-types'
import { Link } from 'react-router'

class MenuView extends React.Component {
  constructor(props) {
    super(props)
  }
  render() {
    return (
      <div>
        <header className="top-nav">
          <div className="pure-menu pure-menu-horizontal">
            <span className="pure-menu-heading fascia-menu-heading">fascia</span>
            <ul className="pure-menu-list fascia-menu-list">
              <li className="pure-menu-item fascia-menu-item"><Link to='/' className="pure-menu-link">projects</Link></li>
            </ul>
            <ul className="pure-menu-list fascia-menu-list right-align-list">
              <li className="pure-menu-item fascia-menu-item pure-menu-has-children pure-menu-allow-hover">
                <a href="#" className="pure-menu-link">account</a>
                <ul className="pure-menu-children">
                  <li className="pure-menu-item fascia-menu-item">
                    <a href="#" className="pure-menu-link" onClick={this.props.menuActions.signOut}>Sign Out</a>
                  </li>
                </ul>
              </li>
            </ul>
          </div>
        </header>
        <div id="main_board">
          {this.props.children}
        </div>
      </div>
    )
  }
}

MenuView.propTypes = {
  children: PropTypes.object.isRequired,
  logout: PropTypes.func.isRequired,
}

export default MenuView
