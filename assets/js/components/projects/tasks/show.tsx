import React from 'react'
import { Modal, Button } from 'react-bootstrap'
import { ThunkDispatch } from 'redux-thunk'
import { RouteComponentProps } from 'react-router-dom'

import Actions, { getTask } from '@/actions/projects/tasks/show'
import { RootStore } from '@/reducers/index'
import styles from './show.scss'

type Props = {
  dispatch: ThunkDispatch<any, any, Actions>
} & RootStore & RouteComponentProps<{ project_id: string, list_id: string, task_id: string }>

class Show extends React.Component<Props> {
  componentDidMount() {
    const projectID = this.props.match.params.project_id
    const listID = this.props.match.params.list_id
    const taskID = this.props.match.params.task_id
    this.props.dispatch(getTask(parseInt(projectID), parseInt(listID), parseInt(taskID)))
  }

  render() {
    const projectID = this.props.match.params.project_id
    const listID = this.props.match.params.list_id
    const taskID = this.props.match.params.task_id

    const hide = () => {
      this.props.history.push(`/projects/${projectID}`)
    }
    const edit = () => {
      this.props.history.push(`/projects/${projectID}/lists/${listID}/tasks/${taskID}/edit`)
    }
    const { task } = this.props.task
    const show = true
    return (
      <Modal
        show={show}
        onHide={hide}
        size="lg"
        centered
      >
        <Modal.Header closeButton>
          <Modal.Title id="contained-modal-title-vcenter">
            {task ? task.title : ''}
          </Modal.Title>
          <div className={styles.header}>
            <i title="Edit" className="fa fa-pencil" onClick={edit}></i>
          </div>
        </Modal.Header>
        <Modal.Body>
          <p>
            {task ? task.description : '' }
          </p>
        </Modal.Body>
        <Modal.Footer>
          <Button onClick={hide}>Close</Button>
        </Modal.Footer>
      </Modal>
    )
  }
}

export default Show
