import Actions, {
  List,
  RequestGetLists,
  ReceiveGetLists,
  Project,
  ReceiveGetProject,
  OpenDelete,
  CloseDelete,
  ReceiveDeleteProject,
  OpenNewList,
  CloseNewList,
  ReceiveNoneList
} from '@/actions/projects/lists'
import NewActions, { ReceiveCreateList } from '@/actions/projects/lists/new'
import { Reducer } from 'redux'

export type State = {
  loading: boolean
  errors: Error | null
  lists: Array<List>
  noneList: List | null
  project: Project | null
  deleteModal: boolean
  newListModal: boolean
  defaultColor: string
}

const initState: State = {
  loading: false,
  errors: null,
  lists: [],
  noneList: null,
  project: null,
  deleteModal: false,
  newListModal: false,
  defaultColor: '008ed4'
}

const reducer: Reducer<State, Actions | NewActions> = (state: State = initState, action: Actions | NewActions): State => {
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
    case ReceiveNoneList:
      return {
        ...state,
        loading: false,
        noneList: action.payload
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
    case OpenNewList:
      return {
        ...state,
        newListModal: true
      }
    case CloseNewList:
      return {
        ...state,
        newListModal: false
      }
    case ReceiveCreateList:
      return {
        ...state,
        newListModal: false
      }
    default:
      return state
  }
}

export default reducer
