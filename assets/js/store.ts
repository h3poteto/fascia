import { createStore, applyMiddleware, compose } from 'redux'
import { createBrowserHistory } from 'history'
import { routerMiddleware } from 'connected-react-router'

import reducers from './reducers/index'

export const history = createBrowserHistory()
const middlewares = [routerMiddleware(history)]

const store = createStore(reducers(history), {}, compose(applyMiddleware(...middlewares)))

export default store
