import React from 'react'
import { ThunkDispatch } from 'redux-thunk'

import Actions, { updatePassword } from '@/actions/settings'
import Form from './settings/form.tsx'
import styles from './settings.scss'

type Props = {
  dispatch: ThunkDispatch<any, any, Actions>
}

class Settings extends React.Component<Props> {
  render() {
    const update = (params: any) => {
      this.props.dispatch(updatePassword(params))
    }

    return (
      <div className={styles.settings}>
        <h3>Account Settings</h3>
        <Form onSubmit={update} />
      </div>
    )
  }
}

export default Settings
