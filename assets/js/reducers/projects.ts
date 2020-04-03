import Actions, { Project, RequestGetProjects, ReceiveGetProjects, OpenNew, CloseNew } from '../actions/projects'
import { Reducer } from 'redux'

export type State = {
  loading: boolean
  errors: Error | null
  projects: Array<Project>
  newModal: boolean
}

const initState: State = {
  loading: false,
  errors: null,
  projects: [],
  newModal: false
}

const reducer: Reducer<State, Actions> = (state: State = initState, action: Actions): State => {
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
    default:
      return state
  }
}
export default reducer
