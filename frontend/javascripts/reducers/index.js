import { combineReducers } from 'redux';
import BoardReducer from './BoardReducer';
import { routerStateReducer } from 'redux-router';

const rootReducer = combineReducers({
  router: routerStateReducer,
  BoardReducer
});

export default rootReducer;
