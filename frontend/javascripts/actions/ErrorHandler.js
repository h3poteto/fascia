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

export function ErrorHandler(err) {
  if (err.response.status === UNPROCESSABLE_ENTITY_ERROR) {
    throw new SubmissionError(err.response.data)
  }
}
