import React from 'react'
import { Modal, Form, Button } from 'react-bootstrap'
import { reduxForm, Field, InjectedFormProps } from 'redux-form'
import { GithubPicker } from 'react-color'
import { ThunkDispatch } from 'redux-thunk'

import styles from './new.scss'
import Actions, { createList } from '@/actions/projects/lists/new'

type Props = {
  open: boolean
  close: Function
  color: string
  projectID: number
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

const renderColorField = (params: {
  input: any
  placeholder: string
  type: string
  color: string
  meta: {
    touched: boolean
    error: string
  } }) => (
  <div className={styles.colorForm}>
    <div className={styles.realColor} style={{backgroundColor: `#${params.input.value || params.color}`}}>　</div>
    <Form.Control {...params.input} placeholder={params.placeholder} type={params.type} />
    {(params.meta.touched || params.type === 'hidden') && params.meta.error && <span className="text-error">{params.meta.error}</span>}
  </div>
)

class NewComponent extends React.Component<InjectedFormProps<{}, Props> & Props> {
  render() {
    const hide = () => {
      this.props.close()
    }

    const create = (params: any) => {
      this.props.dispatch(createList(this.props.projectID, params))
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
            New List
          </Modal.Title>
        </Modal.Header>
        <Form onSubmit={handleSubmit(create)}>
          <Modal.Body>
            <Form.Group controlId="title">
              <Form.Label>Title</Form.Label>
              <Field component={renderField} name="title" id="title" type="text" placeholder="list name" />
            </Form.Group>
            <Form.Group controlId="color">
              <Form.Label>Color</Form.Label>
              <Field component={renderColorField} name="color" type="text" placeholder="008ed4" color={this.props.color} />
              <GithubPicker
                onChangeComplete={(color) => {
                  this.props.change('color', color.hex.replace(/#/g, ''))
                }
                }
              />
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
  if (!values.color) {
    errors = {
      color: 'color is required'
    }
  }
  return errors
}


export default reduxForm<{}, Props>({form: 'new-list-form', validate})(NewComponent)
