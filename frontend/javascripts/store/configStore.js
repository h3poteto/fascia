import { createStore, applyMiddleware, compose } from 'redux'
import rootReducer from '../reducers'
import thunk from 'redux-thunk'
import { createLogger } from 'redux-logger'

const logger = createLogger()
let middleware = [thunk]
if (process.env.NODE_ENV !== 'production') {
  middleware = [...middleware, logger]
}
const createAppStore = compose(
  applyMiddleware(...middleware)
)(createStore)


export default function configureStore() {
  return createAppStore(rootReducer)
}
