import React from 'react'
import { RouteComponentProps } from 'react-router-dom'
import { ThunkDispatch } from 'redux-thunk'

import Actions, { createTask } from '@/actions/projects/tasks/new'
import { RootStore } from '@/reducers/index'
import TaskForm from './form.tsx'

export type Props = {
  dispatch: ThunkDispatch<any, any, Actions>
} & RootStore & RouteComponentProps<{ project_id: string, list_id: string }>

class New extends React.Component<Props> {
  render() {
    const create = (params: any) => {
      const projectID = parseInt(this.props.match.params.project_id)
      const listID = parseInt(this.props.match.params.list_id)
      this.props.dispatch(createTask(projectID, listID, params))
    }

    const hide = () => {
      const projectID = parseInt(this.props.match.params.project_id)
      this.props.history.push(`/projects/${projectID}`)
    }

    return (
      <div>
        <TaskForm hide={hide} onSubmit={create} />
      </div>
    )
  }
}

export default New
