import { combineReducers } from 'redux';
import ProjectReducer from './ProjectReducer';
import ListReducer from './ListReducer';
import { routerStateReducer } from 'redux-router';

const rootReducer = combineReducers({
  router: routerStateReducer,
  ProjectReducer,
  ListReducer
});

export default rootReducer;
