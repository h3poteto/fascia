import { combineReducers } from 'redux'
import ProjectReducer from './ProjectReducer'
import ListReducer from './ListReducer'
import { routerReducer } from 'react-router-redux'
import { reducer as formReducer } from 'redux-form'

const rootReducer = combineReducers({
  routing: routerReducer,
  form: formReducer,
  ProjectReducer,
  ListReducer
})

export default rootReducer
