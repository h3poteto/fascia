import { createStore, applyMiddleware, compose } from 'redux'
import rootReducer from '../reducers'
import thunk from 'redux-thunk'
import { createLogger } from 'redux-logger'

const logger = createLogger()
const createAppStore = compose(
  applyMiddleware(thunk, logger)
)(createStore)


export default function configureStore() {
  return createAppStore(rootReducer)
}
