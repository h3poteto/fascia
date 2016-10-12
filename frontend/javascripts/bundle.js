import React from 'react'
import ReactDOM from 'react-dom'
import { Router, Route, Link, IndexRoute, browserHistory } from 'react-router'
import { Provider } from 'react-redux'
import configureStore from './store/configStore'
import { syncHistoryWithStore } from 'react-router-redux'
import projectContainer from './containers/ProjectContainer'
import listContainer from './containers/ListContainer'
import menuContainer from './containers/MenuContainer'

const store = configureStore()
const history = syncHistoryWithStore(browserHistory, store)

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
)
