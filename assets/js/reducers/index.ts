import { combineReducers } from 'redux'
import { History } from 'history'
import { connectRouter } from 'connected-react-router'

export type Store = {}

const reducers = (history: History) =>
  combineReducers<Store>({
    router: connectRouter(history)
  })

export default reducers
