import React from 'react';
import { Router, Route, Link, IndexRoute } from 'react-router';
import { reduxReactRouter, routerStateReducer, ReduxRouter } from 'redux-router';
import { Provider } from 'react-redux';
import configureStore from './store/configStore';
import { history } from 'history';
import projectContainer from './containers/ProjectContainer';
import listContainer from './containers/ListContainer';
import menuContainer from './containers/MenuContainer';

const store = configureStore();

React.render(
  <Provider store={store}>
    {() =>
      <ReduxRouter>
        <Route history={history}>
          <Route path="/" component={menuContainer}>
            <Route path="/projects/:projectId" component={listContainer} />
            <IndexRoute component={projectContainer} />
          </Route>
        </Route>
      </ReduxRouter>
    }
  </Provider>,
  document.getElementById("content"));
