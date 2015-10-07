import { createStore, applyMiddleware, compose } from 'redux';
import rootReducer from '../reducers';
import thunk from 'redux-thunk';
import createLogger from 'redux-logger';
import { createHistory } from 'history';
import { routerStateReducer, reduxReactRouter } from 'redux-router';

const logger = createLogger();
const createAppStore = compose(
  applyMiddleware(thunk, logger),
  reduxReactRouter({createHistory})
)(createStore);


export default function configureStore() {
  return createAppStore(rootReducer);
}
