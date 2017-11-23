import React from 'react'
import PropTypes from 'prop-types'
import Modal from 'react-modal'
import { Field, reduxForm } from 'redux-form'

import { RenderField, validate } from '../projectForm'

const customStyles = {
  overlay: {
    position          : 'fixed',
    top               : 0,
    left              : 0,
    right             : 0,
    bottom            : 0,
    backgroundColor   : 'rgba(255, 255, 255, 0.5)'
  },
  content: {
    position : 'fixed',
    top : '50%',
    left : '50%',
    right : 'auto',
    bottom : 'auto',
    marginRight : '-50%',
    transform : 'translate(-50%, -50%)'
  }
}

class NewProjectModal extends React.Component {
  componentWillReceiveProps(nextProps) {
    if (!nextProps.isModalOpen) {
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
      repositories,
    } = this.props
    return (
      <Modal
          isOpen={this.props.isModalOpen}
          onRequestClose={onRequestClose}
          style={customStyles}
          contentLabel="NewProjectModal"
      >
        <div className="project-form">
          <form className="pure-form pure-form-stacked" onSubmit={handleSubmit(action)}>
            <fieldset>
              <legend>Create Project</legend>
              <label htmlFor="title">Title</label>
              <Field component={RenderField} name="title" id="title" type="text" placeholder="Project name" />
              <label htmlFor="description">Description</label>
              <Field name="description" id="description" component="textarea" placeholder="Description" className="form-control" />
              <label htmlFor="repository_id">GitHub</label>
              <Field name="repository_id" id="repository_id" component="select" className="form-control">
                <option value="0">--</option>
                {repositories.map(function(repo, index) {
                  return <option key={index} value={repo.id}>{repo.full_name}</option>
                 }, this)}
              </Field>
              <div className="form-action">
                <button type="reset" className="pure-button pure-button-default" disabled={pristine || submitting} onClick={reset}>Reset</button>
                <button type="submit" className="pure-button pure-button-primary" disabled={pristine || submitting}>Create Project</button>
              </div>
            </fieldset>
          </form>
        </div>
      </Modal>
    )
  }
}

NewProjectModal.propTypes = {
  initialize: PropTypes.func.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  pristine: PropTypes.bool,
  reset: PropTypes.func.isRequired,
  submitting: PropTypes.bool.isRequired,
  onRequestClose: PropTypes.func.isRequired,
  action: PropTypes.func.isRequired,
  repositories: PropTypes.array,
  isModalOpen: PropTypes.bool.isRequired,
}

export default reduxForm({
  form: 'new-project-form',
  validate,
})(NewProjectModal)
