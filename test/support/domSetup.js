import { JSDOM } from 'jsdom'

const dom = new JSDOM('<!doctype html><html><body><div id="content"></div></body></html>')
global.document = dom.window.document
global.window = document.defaultView
global.navigator = global.window.navigator
