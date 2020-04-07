import { Reducer } from 'redux'

import Actions, {
  RequestGetLists,
  ReceiveGetLists,
  OpenDelete,
  CloseDelete,
  OpenNewList,
  CloseNewList,
  ReceiveNoneList,
  OpenEditProject,
  CloseEditProject
} from '@/actions/projects/lists'
import NewActions, { ReceiveCreateList } from '@/actions/projects/lists/new'
import EditProjectActions, { ReceiveUpdateProject } from '@/actions/projects/edit'
import DeleteProjectActions, { ReceiveDeleteProject } from '@/actions/projects/delete'
import ProjectActions, { ReceiveGetProject } from '@/actions/projects/show'
import { List } from '@/entities/list'
import { Project } from '@/entities/project'

export type State = {
  loading: boolean
  errors: Error | null
  lists: Array<List>
  noneList: List | null
  project: Project | null
  deleteModal: boolean
  newListModal: boolean
  defaultColor: string
  editProjectModal: boolean
}

const initState: State = {
  loading: false,
  errors: null,
  lists: [],
  noneList: null,
  project: null,
  deleteModal: false,
  newListModal: false,
  defaultColor: '008ed4',
  editProjectModal: false
}

const reducer: Reducer<State, Actions | NewActions> = (
  state: State = initState,
  action: Actions | NewActions | EditProjectActions | DeleteProjectActions | ProjectActions
): State => {
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
    case ReceiveCreateList:
      return {
        ...state,
        newListModal: false
      }
    case OpenEditProject:
      return {
        ...state,
        editProjectModal: true
      }
    case CloseEditProject:
    case ReceiveUpdateProject:
      return {
        ...state,
        editProjectModal: false
      }
    default:
      return state
  }
}

export default reducer
