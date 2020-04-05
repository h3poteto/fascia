import React from 'react'
import { Modal, Button, Form } from 'react-bootstrap'
import { reduxForm, Field, InjectedFormProps } from 'redux-form'

import { Task } from '@/actions/projects/tasks/show'

type Props = {
  hide: Function
  onSubmit: Function
  operation: string
  task: Task | null
}

type FormData = {
  title: string
  description: string
}

const renderField = (params: {
  input: any
  type: string
  placeholder: string
  meta: {
    touched: boolean
    error: string
  }
}) => (
  <div>
    <Form.Control {...params.input} type={params.type} placeholder={params.placeholder} />
    {(params.meta.touched || params.type === 'hidden') && params.meta.error &&
     <span className="text-danger">{params.meta.error}</span>}
  </div>
)

const renderTextarea = (params: {
  input: any
  placeholder: string
}) => (
  <Form.Control {...params.input} as="textarea" placeholder={params.placeholder} />
)

class TaskForm extends React.Component<Props & InjectedFormProps<FormData, Props>> {
  componentDidMount() {
    if (this.props.task) {
      this.handleInitialize(this.props.task)
    }
  }

  componentDidUpdate(prevProps: Props) {
    if (prevProps.task === null && this.props.task) {
      this.handleInitialize(this.props.task)
    }
  }

  handleInitialize(task: Task) {
    this.props.initialize({
      title: task.title,
      description: task.description
    })
  }

  render() {
    const show = true

    const { pristine, submitting, handleSubmit } = this.props

    const hide = () => {
      this.props.hide()
    }

    return (
      <Modal
        show={show}
        onHide={hide}
        size="lg"
        centered
      >
        <Modal.Header closeButton>
          <Modal.Title>
            {this.props.operation} Task
          </Modal.Title>
        </Modal.Header>
        <Form onSubmit={handleSubmit}>
          <Modal.Body>
            <Form.Group controlId="title">
              <Form.Label>Title</Form.Label>
              <Field component={renderField} name="title" id="title" type="text" placeholder="task name" />
            </Form.Group>
            <Form.Group controlId="description">
              <Form.Label>Description</Form.Label>
              <Field component={renderTextarea} name="description" placeholder="description" />
            </Form.Group>
          </Modal.Body>
          <Modal.Footer>
            <Button type="submit" disabled={pristine || submitting}>Submit</Button>
          </Modal.Footer>
        </Form>
      </Modal>
    )
  }
}

const validate = (values: any) => {
  let errors = {}
  if (!values.title) {
    errors = {
      title: 'title is required'
    }
  }
  return errors
}

export default reduxForm<FormData, Props>({ form: 'new-task-form', validate})(TaskForm)
