import { combineReducers } from 'redux';
import ProjectReducer from './ProjectReducer';
import ListReducer from './ListReducer';
import { routeReducer } from 'redux-simple-router'

const rootReducer = combineReducers({
  routing: routeReducer,
  ProjectReducer,
  ListReducer
});

export default rootReducer;
