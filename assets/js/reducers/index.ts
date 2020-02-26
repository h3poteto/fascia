import { combineReducers } from 'redux'
import { History } from 'history'
import { connectRouter, RouterState } from 'connected-react-router'

import projectsReducer, { State as ProjectState } from './projects'

export type RootStore = {
  router: RouterState
  projects: ProjectState
}

const reducers = (history: History) =>
  combineReducers<RootStore>({
    router: connectRouter(history),
    projects: projectsReducer
  })

export default reducers
