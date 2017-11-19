import React from 'react'
import { reduxForm } from 'redux-form'

export function validate(values) {
  const errors = {}
  if (!values.title) {
    errors.title = 'title is required'
  }
  return errors
}

export const RenderField = ({ name, input, placeholder, type, meta: { touched, error } }) => (
  <div>
    <input {...input} placeholder={placeholder} type={type} className="form-control" />
    {(touched || type === "hidden") && error && <span className="text-error">{error}</span>}
  </div>
)
