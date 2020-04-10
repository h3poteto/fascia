import Actions, { RequestGetList, ReceiveGetList, ReceiveGetListOptions } from '@/actions/projects/lists/show'
import ProjectActions, { ReceiveGetProject } from '@/actions/projects/show'
import { List } from '@/entities/list'
import { Project } from '@/entities/project'
import { ListOption } from '@/entities/list_option'
import { Reducer } from 'redux'

export type State = {
  list: List | null
  project: Project | null
  list_options: Array<ListOption>
}

const initState: State = {
  list: null,
  project: null,
  list_options: []
}

const reducer: Reducer<State, Actions> = (state: State = initState, action: Actions | ProjectActions): State => {
  switch (action.type) {
    case RequestGetList:
      return state
    case ReceiveGetList:
      return {
        ...state,
        list: action.payload
      }
    case ReceiveGetProject:
      return {
        ...state,
        project: action.payload
      }
    case ReceiveGetListOptions:
      return {
        ...state,
        list_options: action.payload
      }
    default:
      return state
  }
}

export default reducer
