import React from 'react';
import { Router, Route, Link, IndexRoute } from 'react-router';
import { reduxReactRouter, routerStateReducer, ReduxRouter } from 'redux-router';
import { Provider } from 'react-redux';
import BoardContainer from './containers/BoardContainer';
import configureStore from './store/configStore';
import { history } from 'history';
import BoardView from './components/BoardView.jsx';

const store = configureStore();

React.render(
  <Provider store={store}>
    {() =>
      <ReduxRouter>
        <Route history={history}>
          <Route path="/" component={BoardContainer}>
            <IndexRoute component={BoardView} />
          </Route>
        </Route>
      </ReduxRouter>
    }
  </Provider>,
  document.getElementById("content"));
