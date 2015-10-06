import { combineReducers } from 'redux';
import posts from './BoardReducer';

const rootReducer = combineReducers({
  posts
});

export default rootReducer;
