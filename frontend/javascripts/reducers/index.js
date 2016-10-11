import { combineReducers } from 'redux'
import ProjectReducer from './ProjectReducer'
import ListReducer from './ListReducer'

const rootReducer = combineReducers({
  ProjectReducer,
  ListReducer
})

export default rootReducer
