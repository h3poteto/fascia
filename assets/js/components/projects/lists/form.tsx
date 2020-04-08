import React from 'react'
import { Modal, Form, Button } from 'react-bootstrap'
import { GithubPicker } from 'react-color'
import { Field, reduxForm, InjectedFormProps } from 'redux-form'

import { List } from '@/entities/list'
import { Project } from '@/entities/project'
import { ListOption } from '@/entities/list_option'
import styles from './new.scss'

type Props = {
  hide: Function
  onSubmit: Function
  list: List | null
  project: Project | null
  list_options: Array<ListOption>
}

type FormData = {
  title: string
  color: string
  option_id: number
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
  meta: {
    touched: boolean
    error: string
  }
}) => (
  <div className={styles.colorForm}>
    <div className={styles.realColor} style={{backgroundColor: `#${params.input.value}`}}>　</div>
    <Form.Control {...params.input} placeholder={params.placeholder} type={params.type} />
    {(params.meta.touched || params.type === 'hidden') && params.meta.error && <span className="text-error">{params.meta.error}</span>}
  </div>
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

class ListForm extends React.Component<Props & InjectedFormProps<FormData, Props>> {
  componentDidMount() {
    if (this.props.list) {
      this.handleInitialize(this.props.list)
    }
  }

  componentDidUpdate(prevProps: Props) {
    if (this.props.list && prevProps.list !== this.props.list) {
      this.handleInitialize(this.props.list)
    }
  }

  handleInitialize(list: List) {
    this.props.initialize({
      title: list.title,
      color: list.color,
      option_id: list.list_option_id
    })
  }

  optionForm() {
    if (this.props.list && this.props.project && this.props.project.repositoryID !== 0) {
    return (
      <Form.Group controlId="option_id">
        <Form.Label>Action</Form.Label>
        <Field component={renderSelect} name="option_id">
          <option value="0">nothing</option>
          {this.props.list_options.map(o => (
            <option key={o.id} value={o.id}>{o.action}</option>
          ))}
        </Field>
      </Form.Group>
    )
    } else {
      return null
    }
  }

  render() {
    const hide = () => {
      this.props.hide()
    }

    const { pristine, submitting, handleSubmit } = this.props

    return (
      <Modal
        show={true}
        onHide={hide}
        size="lg"
        centered
      >
        <Modal.Header closeButton>
          <Modal.Title>
            Edit List
          </Modal.Title>
        </Modal.Header>
        <Form onSubmit={handleSubmit}>
          <Modal.Body>
            <Form.Group controlId="title">
              <Form.Label>Title</Form.Label>
              <Field component={renderField} name="title" id="title" type="text" placeholder="list name" />
            </Form.Group>
            <Form.Group controlId="color">
              <Form.Label>Color</Form.Label>
              <Field component={renderColorField} name="color" type="text" placeholder="008ed4" />
              <GithubPicker
                onChangeComplete={(color) => {
                  this.props.change('color', color.hex.replace(/#/g, ''))
                }
                }
              />
            </Form.Group>
            {this.optionForm()}
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

export default reduxForm<FormData, Props>({ form: 'edit-list-form', validate})(ListForm)
