import React from 'react';
import { Router, Route, Link, IndexRoute } from 'react-router';
import { reduxReactRouter, routerStateReducer, ReduxRouter } from 'redux-router';
import { Provider } from 'react-redux';
import configureStore from './store/configStore';
import { history } from 'history';
import boardContainer from './containers/BoardContainer';
import menuContainer from './containers/MenuContainer';

const store = configureStore();

React.render(
  <Provider store={store}>
    {() =>
      <ReduxRouter>
        <Route history={history}>
          <Route path="/" component={menuContainer}>
            <IndexRoute component={boardContainer} />
          </Route>
        </Route>
      </ReduxRouter>
    }
  </Provider>,
  document.getElementById("content"));
