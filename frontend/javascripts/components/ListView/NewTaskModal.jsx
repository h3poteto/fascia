import React from 'react'
import PropTypes from 'prop-types'
import Modal from 'react-modal'
import { Field, reduxForm } from 'redux-form'

import { RenderField, validate } from './taskForm'

const customStyles = {
  overlay : {
    position          : 'fixed',
    top               : 0,
    left              : 0,
    right             : 0,
    bottom            : 0,
    backgroundColor   : 'rgba(255, 255, 255, 0.5)'
  },
  content : {
    position : 'fixed',
    top : '50%',
    left : '50%',
    right : 'auto',
    bottom : 'auto',
    marginRight : '-50%',
    transform : 'translate(-50%, -50%)'
  }
}

class NewTaskModal extends React.Component {
  componentWillReceiveProps(nextProps) {
    if (!nextProps.isTaskModalOpen) {
      this.handleInitialize()
    }
  }

  handleInitialize() {
    this.props.initialize({})
  }

  render() {
    const {
      handleSubmit,
      pristine,
      reset,
      submitting,
      onRequestClose,
      action,
      flashMessage
    } = this.props
    return (
      <Modal
        isOpen={this.props.isTaskModalOpen}
        onRequestClose={onRequestClose}
        style={customStyles}
        contentLabel="NewTaskModal"
      >
        <div className="task-form">
          <form className="pure-form pure-form-stacked" onSubmit={handleSubmit(action)} >
            <fieldset>
              <legend>Create Task</legend>
              <div className="flash flash-error">{flashMessage}</div>
              <label htmlFor="title">Title</label>
              <Field component={RenderField} name="title" id="title" type="text" placeholder="Task Name" />
              <label htmlFor="description">Description</label>
              <Field name="description" id="description" component="textarea" placeholder="Task description" className="form-control" />
              <div className="form-action">
                <button type="reset" className="pure-button pure-button-default" disabled={pristine || submitting} onClick={reset}>Reset</button>
                <button type="submit" className="pure-button pure-button-primary" disabled={pristine || submitting}>Create Task</button>
              </div>
            </fieldset>
          </form>
        </div>
      </Modal>
    )
  }
}

NewTaskModal.propTypes = {
  initialize: PropTypes.func.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  pristine: PropTypes.bool,
  reset: PropTypes.func.isRequired,
  submitting: PropTypes.bool.isRequired,
  onRequestClose: PropTypes.func.isRequired,
  action: PropTypes.func.isRequired,
  isTaskModalOpen: PropTypes.bool.isRequired,
  flashMessage: PropTypes.string,
}

export default reduxForm({
  form: 'new-task-form',
  validate,
})(NewTaskModal)
