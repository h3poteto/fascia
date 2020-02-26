import * as React from 'react'
import * as ReactDOM from 'react-dom'
import { Provider } from 'react-redux'
import { BrowserRouter, Route } from 'react-router-dom'
import { ConnectedRouter } from 'connected-react-router'
import 'bootstrap/dist/css/bootstrap.min.css'

import Menu from './containers/menu'
import projects from './containers/projects'
import store, { history } from './store'

ReactDOM.render(
  <Provider store={store}>
    <ConnectedRouter history={history}>
      <div>
        <BrowserRouter>
          <div>
            <Menu>
              <Route exact path="/" component={projects} />
            </Menu>
          </div>
        </BrowserRouter>
      </div>
    </ConnectedRouter>
  </Provider>,
  document.getElementById('app')
)

