import React from 'react';

var MenuView = function(props, selected) {
  return (
    <div>
      <header className="top-nav">
        <div className="pure-menu pure-menu-horizontal">
          <span className="pure-menu-heading fascia-menu-heading">fascia</span>
          <ul className="pure-menu-list fascia-menu-list">
            <li className="pure-menu-item fascia-menu-item"><a href="/" className="pure-menu-link">projects</a></li>
          </ul>
          <ul className="pure-menu-list fascia-menu-list right-align-list">
            <li className="pure-menu-item fascia-menu-item pure-menu-has-children pure-menu-allow-hover"><a href="#" className="pure-menu-link">settings</a>
              <ul className="pure-menu-children">
                <li className="pure-menu-item fascia-menu-item"><a href="#" className="pure-menu-link">profile</a></li>
                <li className="pure-menu-item fascia-menu-item"><a href="#" className="pure-menu-link">account</a></li>
              </ul>
            </li>
          </ul>
        </div>
      </header>
      <div id="main_board">
        {props.children}
      </div>
    </div>
  );

}

export default MenuView;
