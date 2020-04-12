import React from 'react'
import { ThunkDispatch } from 'redux-thunk'

import Actions, { updatePassword, getSession } from '@/actions/settings'
import Form from './settings/form.tsx'
import styles from './settings.scss'
import { RootStore } from '@/reducers'

type Props = {
  dispatch: ThunkDispatch<any, any, Actions>
} & RootStore

class Settings extends React.Component<Props> {
  componentDidMount() {
    this.props.dispatch(getSession())
  }

  render() {
    const update = (params: any) => {
      this.props.dispatch(updatePassword(params))
    }

    return (
      <div className={styles.settings}>
        <h3>Account Settings</h3>
        <Form onSubmit={update} user={this.props.settings.user} />
      </div>
    )
  }
}

export default Settings
