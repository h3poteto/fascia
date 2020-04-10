import React from 'react'
import { RouteComponentProps } from 'react-router-dom'
import { ThunkDispatch } from 'redux-thunk'

import Actions, { updateTask } from '@/actions/projects/tasks/edit'
import TaskActions, { getTask} from '@/actions/projects/tasks/show'
import { RootStore } from '@/reducers/index'
import TaskForm from './form.tsx'

export type Props = {
  dispatch: ThunkDispatch<any, any, Actions | TaskActions>
} & RootStore & RouteComponentProps<{ project_id: string, list_id: string, task_id: string }>

class Edit extends React.Component<Props> {
  componentDidMount() {
    const projectID = this.props.match.params.project_id
    const listID = this.props.match.params.list_id
    const taskID = this.props.match.params.task_id
    this.props.dispatch(getTask(parseInt(projectID), parseInt(listID), parseInt(taskID)))
  }

  render() {
    const update = (params: any) => {
      const projectID = parseInt(this.props.match.params.project_id)
      const listID = parseInt(this.props.match.params.list_id)
      const taskID = parseInt(this.props.match.params.task_id)
      this.props.dispatch(updateTask(projectID, listID, taskID, params))
    }

    const hide = () => {
      const projectID = parseInt(this.props.match.params.project_id)
      this.props.history.push(`/projects/${projectID}`)
    }

    return (
      <div>
        <TaskForm hide={hide} onSubmit={update} operation="Edit" task={this.props.task.task} />
      </div>
    )
  }
}

export default Edit
