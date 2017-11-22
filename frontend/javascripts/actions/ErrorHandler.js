import { SubmissionError } from 'redux-form'

// const BAD_REQUEST_ERROR = 400
const AUTHENTICATE_ERROR = 401
const FORBIDDEN_ERROR = 403
const NOT_FOUND_ERROR = 404
// const METHOD_NOT_ALLOWED_ERROR = 405
// const NOT_ACCEPTABLE_ERROR = 406
const REQUEST_TIMEOUT_ERROR = 408
const UNPROCESSABLE_ENTITY_ERROR = 422
const INTERNAL_SERVER_ERROR = 500
const BAD_GATEWAY_ERROR = 502
const SERVICE_UNAVAILABLE_ERROR = 503
const GATEWAY_TIMEOUT_ERROR = 504

export const SERVER_ERROR = 'SERVER_ERROR'
export function ServerError(error) {
  return {
    type: SERVER_ERROR,
    message: error.message,
    status: error.status,
  }
}

class Error {
  constructor(status, message) {
    this.status = status
    this.message = message
  }
}

export function ErrorHandler(err) {
  if (err.response.status === UNPROCESSABLE_ENTITY_ERROR) {
    throw new SubmissionError(err.response.data)
  }

  return handler(err)
}

export function ErrorHandlerWithoutSubmission(err) {
  return handler(err)
}

export function handler(err) {
  return new Promise((resolve, reject) => {
    switch (err.response.status) {
      case AUTHENTICATE_ERROR:
        // ログインページはreact管理ではない
        // そのためbrowserHistoryでの移動ができないので，locationを直接書き換えてリダイレクトさせる
        window.location.pathname = '/sign_in'
        resolve(err)
        return
      case FORBIDDEN_ERROR:
      case NOT_FOUND_ERROR:
        reject(new Error(err.response.status, 'The requested URL was not found.'))
        return
      case REQUEST_TIMEOUT_ERROR:
        reject(new Error(err.response.status, 'Request timeout.'))
        return
      case UNPROCESSABLE_ENTITY_ERROR:
        reject(new Error(err.response.status, 'Validation error.'))
        return
      case BAD_GATEWAY_ERROR:
      case SERVICE_UNAVAILABLE_ERROR:
      case GATEWAY_TIMEOUT_ERROR:
        reject(new Error(err.response.status, 'Could not connect the server.'))
        return
      case INTERNAL_SERVER_ERROR:
        reject(new Error(err.response.status, 'Server error. Sorry, but something went wrong.'))
        return
      default:
        reject(new Error(err.response.status, 'Unknown error'))
    }
  })
}
