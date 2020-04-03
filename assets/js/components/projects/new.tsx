import React from 'react'
import { Modal, Button, Form } from 'react-bootstrap'
import { Field, reduxForm, InjectedFormProps } from 'redux-form'

type Props = {
  open: boolean,
  close: Function
}

class NewComponent extends React.Component<InjectedFormProps<{}> & Props> {
  render() {
    const hide = () => {
      this.props.close()
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
        <Form.Control type={params.type} placeholder={params.placeholder} {...params.input} />
        {(params.meta.touched || params.type === 'hidden') && params.meta.error && <span className="text-error">{params.meta.error}</span>}
      </div>
    )

    const renderTextarea = (params: {
      placeholder: string
    }) => (
      <Form.Control as="textarea" placeholder={params.placeholder} />
    )

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
        <Modal.Body>
          <Form>
            <Form.Group controlId="title">
              <Form.Label>Title</Form.Label>
              <Field component={renderField} name="title" id="title" type="text" placeholder="project name" />
            </Form.Group>
            <Form.Group controlId="description">
              <Form.Label>Description</Form.Label>
              <Field component={renderTextarea} name="description" placeholder="description" />
            </Form.Group>
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <Button onClick={hide}>Close</Button>
        </Modal.Footer>
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
