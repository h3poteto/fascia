import Actions, { Task, RequestGetTask, ReceiveGetTask } from '@/actions/projects/tasks/show'
import { Reducer } from 'redux'

export type State = {
  task: Task | null
}

const initState: State = {
  task: null
}

const reducer: Reducer<State, Actions> = (state: State = initState, action: Actions): State => {
  switch (action.type) {
    case RequestGetTask:
      return state
    case ReceiveGetTask:
      return {
        ...state,
        task: action.payload
      }
    default:
      return state
  }
}

export default reducer
