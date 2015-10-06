import React from 'react';
import {Router, Route, Link, IndexRoute} from 'react-router';
import { Provider } from 'react-redux';
import BoardContainer from './containers/BoardContainer';
import configureStore from './store/configStore';

const store = configureStore();

React.render(
  <Provider store={store}>
    {() => <BoardContainer />}
  </Provider>,
  document.getElementById("content"));
