import { combineReducers } from 'redux'
import { History } from 'history'
import { connectRouter, RouterState } from 'connected-react-router'
import { reducer as formReducer, FormStateMap } from 'redux-form'

import projectsReducer, { State as ProjectsState } from './projects'
import listsReducer, { State as ListsState } from './lists'
import taskReducer, { State as TaskState } from './projects/tasks/show'
import listReducer, { State as ListState } from './projects/lists/show'
import settingsReducer, { State as SettingsState } from './settings'

export type RootStore = {
  router: RouterState
  projects: ProjectsState
  lists: ListsState
  task: TaskState
  list: ListState
  settings: SettingsState
  form: FormStateMap
}

const reducers = (history: History) =>
  combineReducers<RootStore>({
    router: connectRouter(history),
    projects: projectsReducer,
    lists: listsReducer,
    task: taskReducer,
    list: listReducer,
    settings: settingsReducer,
    form: formReducer
  })

export default reducers
