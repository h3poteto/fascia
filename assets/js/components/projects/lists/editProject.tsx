import React from 'react'
import { Modal, Button, Form } from 'react-bootstrap'
import { reduxForm, Field, InjectedFormProps } from 'redux-form'
import { ThunkDispatch } from 'redux-thunk'

import { Project } from '@/actions/projects'
import Actions, { updateProject } from '@/actions/projects/edit'

type Props = {
  open: boolean
  close: Function
  project: Project | null
  dispatch: ThunkDispatch<any, any, Actions>
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


class Edit extends React.Component<InjectedFormProps<FormData, Props> & Props> {
  componentDidMount() {
    if (this.props.project) {
      this.handleInitialize(this.props.project)
    }
  }

  componentDidUpdate(prevProps: Props) {
    if (!prevProps.project && this.props.project) {
      this.handleInitialize(this.props.project)
    }
  }

  handleInitialize(project: Project) {
    this.props.initialize({
      title: project.title,
      description: project.description
    })
  }

  render() {
    const hide = () => {
      this.props.close()
    }

    const update = (params: any) => {
      if (!this.props.project) {
        return
      }
      this.props.dispatch(updateProject(this.props.project.id, params))
    }

    const { pristine, submitting, handleSubmit } = this.props

    return (
      <Modal
        show={this.props.open}
        onHide={hide}
        size="lg"
        centered
      >
        <Modal.Header closeButton>
          <Modal.Title>
            Edit Project
          </Modal.Title>
        </Modal.Header>
        <Form onSubmit={handleSubmit(update)}>
          <Modal.Body>
            <Form.Group controlId="title">
              <Form.Label>Title</Form.Label>
              <Field component={renderField} name="title" id="title" type="text" placeholder="project name" />
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

export default reduxForm<FormData, Props>({form: 'edit-project-form', validate})(Edit)
