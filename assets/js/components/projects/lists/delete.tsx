import React from 'react'
import { Modal, Form, Button } from 'react-bootstrap'
import { Field, reduxForm, InjectedFormProps } from 'redux-form'
import { ThunkDispatch } from 'redux-thunk'

import { Project } from '@/entities/project'
import Actions, { deleteProject } from '@/actions/projects/delete'


type Props = {
  project: Project | null
  open: boolean
  close: Function
  dispatch: ThunkDispatch<any, any, Actions>
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

class DeleteComponent extends React.Component<InjectedFormProps<{}, Props> & Props> {
  render() {
    const hide = () => {
      this.props.close()
    }

    const del = () => {
      if (!this.props.project) {
        return
      }
      this.props.dispatch(deleteProject(this.props.project.id))
    }

    const { handleSubmit, pristine, submitting } = this.props

    return (
      <Modal
        show={this.props.open}
        onHide={hide}
        size="lg"
        centered
      >
        <Modal.Header closeButton>
          <Modal.Title>
            Delete Project
          </Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <h5>Are you sure?</h5>
          <p>This action can not be undone. This will permanently delete all tasks and lists in this project. Issues, pull requests, and labels in the repository never changes at all with this action, when this project is associated with a github repository.</p>
          <p>Please type in the name of the project to confirm.</p>
          <Form onSubmit={handleSubmit(del)}>
            <Form.Group controlId="name">
              <Field component={renderField} name="name" id="name" type="text" />
            </Form.Group>
            <Button type="submit" variant="danger" disabled={pristine || submitting}>Yes, I understand, and delete this project</Button>
          </Form>
        </Modal.Body>
      </Modal>
    )
  }
}


function validate(values: any, props: Props) {
  let errors = {}
  if (!values.name) {
    errors = {
      name: 'Required'
    }
  } else if (!props.project) {
    errors = {
      name: 'Project is invalid'
    }
  } else if (values.name != props.project.title) {
    errors = {
      name: 'Invalid project name'
    }
  }
  return errors
}

export default reduxForm<{}, Props>({form: 'delete-project-form', validate})(DeleteComponent)
