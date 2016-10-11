import React from 'react'
import ReactDOM from 'react-dom'
import { Router, Route, IndexRoute, hashHistory } from 'react-router'
import { Provider } from 'react-redux'
import configureStore from './store/configStore'
import projectContainer from './containers/ProjectContainer'
import listContainer from './containers/ListContainer'
import menuContainer from './containers/MenuContainer'

const store = configureStore()

ReactDOM.render(
  <Provider store={store}>
    <Router history={hashHistory}>
      <Route path="/" component={menuContainer}>
        <Route path="/projects/:projectID" component={listContainer} />
        <IndexRoute component={projectContainer} />
      </Route>
    </Router>
  </Provider>,
  document.getElementById("content")
)
