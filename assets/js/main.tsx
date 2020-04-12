import * as React from 'react'
import * as ReactDOM from 'react-dom'
import { Provider } from 'react-redux'
import {  Route, Switch } from 'react-router-dom'
import { ConnectedRouter } from 'connected-react-router'
import { loadProgressBar } from 'axios-progress-bar'
import 'bootstrap/dist/css/bootstrap.min.css'
import 'axios-progress-bar/dist/nprogress.css'

import Menu from './containers/menu'
import projects from './containers/projects'
import lists from './containers/projects/lists'
import store, { history } from './store'
import Task from '@/containers/projects/tasks/show'
import NewTask from '@/containers/projects/tasks/new'
import EditTask from '@/containers/projects/tasks/edit'
import EditList from '@/containers/projects/lists/edit'
import Settings from '@/containers/settings'
import './axios-progress-bar.css'

loadProgressBar()

ReactDOM.render(
  <Provider store={store}>
    <ConnectedRouter history={history}>
      <div>
        <Menu>
          <Route exact path="/" component={projects} />
          <Route exact path="/settings" component={Settings} />
          <Route path="/projects/:project_id" component={lists} />
          <Switch>
            <Route exact path="/projects/:project_id/lists/:list_id/tasks/new" component={NewTask} />
            <Route exact path="/projects/:project_id/lists/:list_id/tasks/:task_id/edit" component={EditTask} />
            <Route exact path="/projects/:project_id/lists/:list_id/tasks/:task_id" component={Task} />
            <Route exact path="/projects/:project_id/lists/:list_id/edit" component={EditList} />
          </Switch>
        </Menu>
      </div>
    </ConnectedRouter>
  </Provider>,
  document.getElementById('app')
)
