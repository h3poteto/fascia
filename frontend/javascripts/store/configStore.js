import { createStore, applyMiddleware } from 'redux';
import rootReducer from '../reducers';
import thunk from 'redux-thunk';
import createLogger from 'redux-logger';

const logger = createLogger();
const createStoreWithMiddleware = applyMiddleware(
  thunk,
  logger
)(createStore);

export default function configureStore() {
  return createStoreWithMiddleware(rootReducer);
}
