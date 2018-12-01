import { createStore, applyMiddleware, compose } from 'redux'
import rootReducer from '../reducers'
import thunk from 'redux-thunk'
import { createLogger } from 'redux-logger'
import { createBrowserHistory } from 'history'
import { routerMiddleware } from 'connected-react-router'

const logger = createLogger()
export const history = createBrowserHistory()

let middleware = [thunk]
if (process.env.NODE_ENV !== 'production') {
  middleware = [...middleware, logger]
}
middleware = [...middleware, routerMiddleware(history)]

export const store = createStore(
  rootReducer(history),
  {},
  compose(
    applyMiddleware(...middleware),
  ),
)
