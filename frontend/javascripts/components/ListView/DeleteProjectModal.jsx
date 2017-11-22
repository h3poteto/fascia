import React from 'react'
import PropTypes from 'prop-types'
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

function validate(values, props) {
  const errors = {}
  if (!values.name) {
    errors.name = 'Required'
  } else if (values.name != props.project.Title) {
    errors.name = 'Invalid project name'
  }
  return errors
}

const renderField = ({input, type, meta: { touched, error } }) => {
  return (
    <div>
      <input type={type} {...input} className="form-control" />
      {touched && error && <span className="text-error">{error}</span>}
    </div>
  )
}

class DeleteProjectModal extends React.Component {
  render() {
    const {
      handleSubmit,
      pristine,
      submitting,
      onRequestClose,
      action,
      project,
      isDeleteProjectModalOpen,
    } = this.props

    return (
      <Modal
          isOpen={isDeleteProjectModalOpen}
          onRequestClose={onRequestClose}
          style={customStyles}
          contentLabel="DeleteProjectModal"
      >
        <div className="delete-project-form">
          <form className="pure-form pure-form-stacked" onSubmit={handleSubmit(action)}>
            <fieldset>
              <legend>Are you sure?</legend>
              <div className="delete-project-message pure-form-message">This action can not be undone.
                This will permanently delete all tasks and lists in this project.
                Issues, pull requests, and labels in the repository never changes at all with this action, when this project is associated with a github repository.
                <div className="confirm">Please type in the name of the project to confirm.</div>
              </div>
              <Field name="name" id="name" component={renderField} type="text" />
              <div className="form-action">
                <button type="submit" className="pure-button button-error" disabled={pristine || submitting}>Yes, I understand, and delete this project</button>
              </div>
            </fieldset>
          </form>
        </div>
      </Modal>
    )
  }
}

DeleteProjectModal.propTypes = {
  handleSubmit: PropTypes.func.isRequired,
  pristine: PropTypes.bool,
  reset: PropTypes.func.isRequired,
  submitting: PropTypes.bool,
  onRequestClose: PropTypes.func.isRequired,
  isDeleteProjectModalOpen: PropTypes.bool.isRequired,
  action: PropTypes.func.isRequired,
  project: PropTypes.object,
}

renderField.propTypes = {
  input: PropTypes.shape().isRequired,
  type: PropTypes.string.isRequired,
  meta: PropTypes.shape().isRequired,
}

export default reduxForm({
  form: 'delete-project-form',
  validate,
})(DeleteProjectModal)
