import React from 'react'
import Modal from 'react-modal'
import { Field, reduxForm } from 'redux-form'

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
      projectID,
      listID,
    } = this.props
    return (
      <Modal
          isOpen={this.props.isTaskModalOpen}
          onRequestClose={onRequestClose}
          style={customStyles}
          contentLabel="NewTaskModal"
      >
        <div className="task-form">
          <form className="pure-form pure-form-stacked" onSubmit={handleSubmit((values) => { action(projectID, listID, values) })} >
            <fieldset>
              <legend>Create Task</legend>
              <label htmlFor="title">Title</label>
              <Field name="title" id="title" component="input" type="text" placeholder="Task Name" className="form-control" />
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
  initialize: React.PropTypes.func.isRequired,
  handleSubmit: React.PropTypes.func.isRequired,
  pristine: React.PropTypes.bool,
  reset: React.PropTypes.func.isRequired,
  submitting: React.PropTypes.bool.isRequired,
  onRequestClose: React.PropTypes.func.isRequired,
  action: React.PropTypes.func.isRequired,
  projectID: React.PropTypes.string.isRequired,
  listID: React.PropTypes.number,
  isTaskModalOpen: React.PropTypes.bool.isRequired,
}

export default reduxForm({
  form: 'new-task-form',
})(NewTaskModal)
