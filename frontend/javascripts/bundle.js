import React from 'react';
import {Router, Route, Link} from 'react-router';
import Request from 'superagent';
import BoardViewModel from './BoardViewModel';
import MenuViewModel from './MenuViewModel';

// TODO: ある程度できたらreduxで状態管理する

var routes = (
    <Router>
    <Route path="/" component={BoardViewModel}>
    </Route>
    </Router>
);

React.render(routes, document.getElementById("board"));

var menuRoutes = (
    <Router>
    <Route path="/" component={MenuViewModel}>
    </Route>
    </Router>
);

React.render(menuRoutes, document.getElementById("top"));
