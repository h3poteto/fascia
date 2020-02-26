import Actions, { Project, RequestGetProjects, ReceiveGetProjects } from '../actions/projects'
import { Reducer } from 'redux'

export type State = {
  loading: boolean
  errors: Error | null
  projects: Array<Project>
}

const initState: State = {
  loading: false,
  errors: null,
  projects: []
}

const reducer: Reducer<State, Actions> = (state: State = initState, action: Actions): State => {
  switch (action.type) {
    case RequestGetProjects:
      return {
        ...state,
        loading: false
      }
    case ReceiveGetProjects:
      return {
        ...state,
        loading: false,
        projects: action.payload
      }
    default:
      return state
  }
}
export default reducer
