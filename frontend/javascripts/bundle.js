import React from 'react';
import {Router, Route, Link, IndexRoute} from 'react-router';
import { Provider } from 'react-redux';
import BoardContainer from './containers/BoardContainer';
import configureStore from './store/configStore';
import { fetchProjects } from './actions/BoardAction';

const store = configureStore();
store.dispatch(fetchProjects());
React.render(
  <Provider store={store}>
    {() => <BoardContainer />}
  </Provider>,
  document.getElementById("content"));
