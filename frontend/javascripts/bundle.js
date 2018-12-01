import React from 'react'
import { Provider } from 'react-redux'
import ReactDOM from 'react-dom'
import { Route, Switch } from 'react-router-dom'
import { loadProgressBar } from 'axios-progress-bar'
import ProjectContainer from './containers/ProjectContainer'
import ListContainer from './containers/ListContainer'
import MenuContainer from './containers/MenuContainer'
import { history, store } from './store/configStore'
import { ConnectedRouter } from 'connected-react-router'

loadProgressBar()

ReactDOM.render(
  <Provider store={store}>
    <ConnectedRouter history={history}>
      <div>
        <MenuContainer>
          <Switch>
            <Route exact path="/" component={ProjectContainer}></Route>
            <Route path="/projects/:projectID" component={ListContainer}></Route>
          </Switch>
        </MenuContainer>
      </div>
    </ConnectedRouter>
  </Provider>,
  document.getElementById('content')
)
