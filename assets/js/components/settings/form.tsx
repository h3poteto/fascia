import React from 'react'
import { Form, Button } from 'react-bootstrap'
import { reduxForm, Field, InjectedFormProps } from 'redux-form'

type Props = {}

type FormData = {
  password: string
  passwordConfirm: string
}

const renderField = (params: {
  input: any
  type: string
  placeholder: string
  disabled: boolean
  meta: {
    touched: boolean
    error: string
  }
}) => (
  <div>
    <Form.Control {...params.input} type={params.type} placeholder={params.placeholder} disabled={params.disabled} />
    {(params.meta.touched || params.type === 'hidden') && params.meta.error &&
     <span className="text-danger">{params.meta.error}</span>}
  </div>
)

class SettingsForm extends React.Component<InjectedFormProps<FormData, Props> & Props> {
  render() {
    const { pristine, submitting, handleSubmit } = this.props

    return (
      <div>
        <Form onSubmit={handleSubmit}>
          <Form.Group controlId="username">
            <Form.Label>Username</Form.Label>
            <Field component={renderField} name="username" id="username" type="text" disabled={true} />
          </Form.Group>
          <Form.Group controlId="password">
            <Form.Label>Password</Form.Label>
            <Field component={renderField} name="password" id="password" type="password" disabled={false} />
          </Form.Group>
          <Form.Group controlId="passwordConfirm">
            <Form.Label>Password Confirm</Form.Label>
            <Field component={renderField} name="passwordConfirm" id="passwordConfirm" type="password" disabled={false} />
          </Form.Group>
          <Button type="submit" disabled={pristine || submitting}>Submit</Button>
        </Form>
      </div>
    )
  }
}

const validate = (values: FormData) => {
  let errors = {}
  if (!values.password) {
    errors = Object.assign(errors, {
      password: 'password is required'
    })
  }
  if (values.password && values.password.length < 12) {
    errors = Object.assign(errors, {
      password: 'password must be over 12 characters'
    })
  }
  if (values.password !== values.passwordConfirm) {
    errors = Object.assign(errors, {
      passwordConfirm: 'password is not matched'
    })
  }
  return errors
}

export default reduxForm<FormData, Props>({form: 'settings-form', validate})(SettingsForm)