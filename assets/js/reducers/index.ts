import { combineReducers } from 'redux'
import { History } from 'history'
import { connectRouter, RouterState } from 'connected-react-router'

import projectsReducer, { State as ProjectsState } from './projects'
import listsReducer, { State as ListsState } from './lists'

export type RootStore = {
  router: RouterState
  projects: ProjectsState
  lists: ListsState
}

const reducers = (history: History) =>
  combineReducers<RootStore>({
    router: connectRouter(history),
    projects: projectsReducer,
    lists: listsReducer
  })

export default reducers
