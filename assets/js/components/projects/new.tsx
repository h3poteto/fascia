import React from 'react'
import { Modal, Button, Form } from 'react-bootstrap'
import { Field, reduxForm, InjectedFormProps } from 'redux-form'
import { ThunkDispatch } from 'redux-thunk'

import Actions, { Repository, createProject } from '@/actions/projects'

type Props = {
  open: boolean,
  close: Function,
  repositories: Array<Repository>
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

const renderTextarea = (params: {
  input: any
  placeholder: string
}) => (
  <Form.Control {...params.input} as="textarea" placeholder={params.placeholder} />
)

const renderSelect = (params: {
  input: any
  meta: {
    touched: boolean
    error: string
  }
  children: any
}) => (
  <Form.Control {...params.input} as="select">
    {params.children}
  </Form.Control>
)

class NewComponent extends React.Component<InjectedFormProps<{}, Props> & Props> {
  render() {
    const hide = () => {
      this.props.close()
    }

    const { handleSubmit } = this.props

    const create = (params: any) => {
      console.log(params)
      this.props.dispatch(createProject(params))
    }

    return (
      <Modal
        show={this.props.open}
        onHide={hide}
        size="lg"
        centered
      >
        <Modal.Header closeButton>
          <Modal.Title id="contained-modal-title-vcenter">
            New Project
          </Modal.Title>
        </Modal.Header>
        <Form onSubmit={handleSubmit(create)} noValidate={true}>
          <Modal.Body>
            <Form.Group controlId="title">
              <Form.Label>Title</Form.Label>
              <Field component={renderField} name="title" id="title" type="text" placeholder="project name" />
            </Form.Group>
            <Form.Group controlId="description">
              <Form.Label>Description</Form.Label>
              <Field component={renderTextarea} name="description" placeholder="description" />
            </Form.Group>
            <Form.Group controlId="repository_id">
              <Form.Label>GitHub</Form.Label>
              <Field component={renderSelect} name="repository_id">
                <option value="0">--</option>
                {this.props.repositories.map(r => (
                  <option key={r.id} value={r.id}>{r.full_name}</option>
                ))}
              </Field>
            </Form.Group>
          </Modal.Body>
          <Modal.Footer>
            <Button type="submit">Submit</Button>
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

export default reduxForm<{}, Props>({form: 'new-project-form', validate})(NewComponent)
