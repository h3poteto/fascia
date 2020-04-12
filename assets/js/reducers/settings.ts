import { Reducer } from 'redux'

import { User } from '@/entities/user'
import Actions, { ReceiveGetSession } from '@/actions/settings'

export type State = {
  user: User | null
}

const initState: State = {
  user: null
}

const reducer: Reducer<State, Actions> = (state: State = initState, action: Actions): State => {
  switch (action.type) {
    case ReceiveGetSession:
      return {
        ...state,
        user: action.payload
      }
    default:
      return state
  }
}

export default reducer
