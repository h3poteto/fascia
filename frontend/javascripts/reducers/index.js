import { combineReducers } from 'redux';
import ProjectReducer from './ProjectReducer';
import ListReducer from './ListReducer';
import { routerReducer } from 'react-router-redux'

const rootReducer = combineReducers({
  routing: routerReducer,
  ProjectReducer,
  ListReducer
});

export default rootReducer;
