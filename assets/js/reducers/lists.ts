import Actions, {
  List,
  RequestGetLists,
  ReceiveGetLists,
  Project,
  ReceiveGetProject,
  OpenDelete,
  CloseDelete,
  ReceiveDeleteProject
} from '@/actions/projects/lists'
import { Reducer } from 'redux'

export type State = {
  loading: boolean
  errors: Error | null
  lists: Array<List>
  project: Project | null
  deleteModal: boolean
}

const initState: State = {
  loading: false,
  errors: null,
  lists: [],
  project: null,
  deleteModal: false
}

const reducer: Reducer<State, Actions> = (state: State = initState, action: Actions): State => {
  switch (action.type) {
    case RequestGetLists:
      return {
        ...state,
        loading: true
      }
    case ReceiveGetLists:
      return {
        ...state,
        loading: false,
        lists: action.payload
      }
    case ReceiveGetProject:
      return {
        ...state,
        project: action.payload
      }
    case OpenDelete:
      return {
        ...state,
        deleteModal: true
      }
    case CloseDelete:
      return {
        ...state,
        deleteModal: false
      }
    case ReceiveDeleteProject:
      return {
        ...state,
        deleteModal: false
      }
    default:
      return state
  }
}

export default reducer
