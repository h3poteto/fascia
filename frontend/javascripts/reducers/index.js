import { combineReducers } from 'redux';
import ProjectReducer from './ProjectReducer';
import { routerStateReducer } from 'redux-router';

const rootReducer = combineReducers({
  router: routerStateReducer,
  ProjectReducer
});

export default rootReducer;
