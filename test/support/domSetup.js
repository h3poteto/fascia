import { JSDOM } from 'jsdom'

global.document = new JSDOM('<!doctype html><html><body></body></html>')
global.window = document.defaultView
//global.navigator = global.window.navigator
