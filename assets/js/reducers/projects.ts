import Actions, {
  Project,
  RequestGetProjects,
  ReceiveGetProjects,
  OpenNew,
  CloseNew,
  ReceiveGetRepositories,
  Repository
} from '../actions/projects'
import NewActions, { ReceiveCreateProject } from '@/actions/projects/new'
import { Reducer } from 'redux'

export type State = {
  loading: boolean
  errors: Error | null
  projects: Array<Project>
  newModal: boolean
  repositories: Array<Repository>
}

const initState: State = {
  loading: false,
  errors: null,
  projects: [],
  newModal: false,
  repositories: []
}

const reducer: Reducer<State, Actions> = (state: State = initState, action: Actions | NewActions): State => {
  switch (action.type) {
    case RequestGetProjects:
      return {
        ...state,
        loading: true
      }
    case ReceiveGetProjects:
      return {
        ...state,
        loading: false,
        projects: action.payload
      }
    case OpenNew:
      return {
        ...state,
        newModal: true
      }
    case CloseNew:
      return {
        ...state,
        newModal: false
      }
    case ReceiveGetRepositories:
      return {
        ...state,
        repositories: action.payload
      }
    case ReceiveCreateProject:
      return {
        ...state,
        newModal: false
      }
    default:
      return state
  }
}
export default reducer
