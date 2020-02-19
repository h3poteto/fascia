import React from 'react'
import { Field, reduxForm } from 'redux-form'

class UserSettings extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    return (
      <div className="user-settings-form">
        <form className="pure-form pure-form-stacked">
          <fieldset>
            <legend>Change Password</legend>
            <label htmlFor="username">Username</label>
            <Field name="username" id="username" component="input" className="form-control" disabled />
            <label htmlFor="password">New Password</label>
            <Field name="password" id="password" component="input" type="password" className="form-control" />
            <div className="form-action">
              <button type="submit" className="pure-button pure-button-primary">Submit</button>
            </div>
          </fieldset>
        </form>
      </div>
    )
  }
}

export default reduxForm({
  form: 'user-settings-form',
})(UserSettings)
