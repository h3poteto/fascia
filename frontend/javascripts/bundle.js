import React from 'react';
import {Router, Route, Link} from 'react-router';
import Request from 'superagent';
import BoardViewModel from './BoardViewModel';


// TODO: ある程度できたらreduxで状態管理する

var routes = (
    <Router>
    <Route path="/" component={BoardViewModel}>
    </Route>
    </Router>
);

React.render(routes, document.getElementById("board"));


var Menu = React.createClass({
  getInitialState: function() {
    return {
      selected: "projects"
    };
  },
  selectMenu: function(selectedClass) {
    this.setState({
      selected: selectedClass
    });
  },
  render: function() {
    return <div className="top-nav">
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
      </div>;
  }
});

var menuRoutes = (
    <Router>
    <Route path="/" component={Menu}>
    </Route>
    </Router>
);

React.render(menuRoutes, document.getElementById("top"));
