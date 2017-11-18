import React from 'react'
import { Field, reduxForm } from 'redux-form'

export function validate(values) {
  const errors = {}
  if (!values.title) {
    errors.title = 'title is required'
  }
  if (!values.color) {
    errors.color = 'color is required'
  }
  return errors
}

export const RenderField = ({ name, input, placeholder, type, meta: { touched, error } }) => (
  <div>
    <input {...input} placeholder={placeholder} type={type} className="form-control" />
    {(touched || type === "hidden") && error && <span className="text-error">{error}</span>}
  </div>
)

export const RenderColorField = ({ name, input, placeholder, type, color, meta: { touched, error } }) => (
  <div className="color-control-group">
    <div className="real-color" style={{backgroundColor: `#${color}`}}>　</div>
    <input {...input} placeholder={placeholder} type={type} />
    {(touched || type === "hidden") && error && <span className="text-error">{error}</span>}
  </div>
)
