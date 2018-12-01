import { combineReducers } from 'redux'
import ProjectReducer from './ProjectReducer'
import ListReducer from './ListReducer'
import { reducer as formReducer } from 'redux-form'
import { connectRouter } from 'connected-react-router'

export default (history) => combineReducers({
  router: connectRouter(history),
  form: formReducer,
  ProjectReducer,
  ListReducer
})
