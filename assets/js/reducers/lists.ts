import Actions, { List, RequestGetLists, ReceiveGetLists } from '@/actions/projects/lists'
import { Reducer } from 'redux'

export type State = {
  loading: boolean
  errors: Error | null
  lists: Array<List>
}

const initState: State = {
  loading: false,
  errors: null,
  lists: []
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
    default:
      return state
  }
}

export default reducer
