import { createStore, applyMiddleware, compose } from 'redux'
import { createBrowserHistory } from 'history'
import { routerMiddleware } from 'connected-react-router'
import thunk from 'redux-thunk'
import { createLogger } from 'redux-logger'

import reducers from './reducers/index'

export const history = createBrowserHistory()
const logger = createLogger()
const middlewares = [thunk, logger, routerMiddleware(history)]

const store = createStore(reducers(history), {}, compose(applyMiddleware(...middlewares)))

export default store
