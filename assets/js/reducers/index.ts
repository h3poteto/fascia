import { combineReducers } from 'redux'
import { History } from 'history'
import { connectRouter, RouterState } from 'connected-react-router'

import projectsReducer, { State as ProjectsState } from './projects'
import listsReducer, { State as ListsState } from './lists'
import taskReducer, { State as TaskState } from './projects/tasks/show'

export type RootStore = {
  router: RouterState
  projects: ProjectsState
  lists: ListsState
  task: TaskState
}

const reducers = (history: History) =>
  combineReducers<RootStore>({
    router: connectRouter(history),
    projects: projectsReducer,
    lists: listsReducer,
    task: taskReducer
  })

export default reducers
