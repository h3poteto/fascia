import React from 'react'
import ReactDOM from 'react-dom'
import { Router, Route, Link, IndexRoute } from 'react-router'
import { Provider } from 'react-redux'
import configureStore from './store/configStore'
import { createHistory } from 'history'
import { syncReduxAndRouter } from 'redux-simple-router'
import projectContainer from './containers/ProjectContainer'
import listContainer from './containers/ListContainer'
import menuContainer from './containers/MenuContainer'

const store = configureStore()
const history = createHistory()
syncReduxAndRouter(history, store)

ReactDOM.render(
  <Provider store={store}>
    <Router history={history}>
      <Route path="/" component={menuContainer}>
        <Route path="/projects/:projectID" component={listContainer} />
        <IndexRoute component={projectContainer} />
      </Route>
    </Router>
  </Provider>,
  document.getElementById("content")
);
