import React from 'react';
import {Router, Route, Link, IndexRoute} from 'react-router';
import Request from 'superagent';
import BoardViewModel from './BoardViewModel';
import MenuViewModel from './MenuViewModel';

// TODO: ある程度できたらreduxで状態管理する

var routes = (
    <Router>
      <Route path="/" component={MenuViewModel}>
        <IndexRoute component={BoardViewModel}/>
      </Route>
    </Router>
);

React.render(routes, document.getElementById("content"));
