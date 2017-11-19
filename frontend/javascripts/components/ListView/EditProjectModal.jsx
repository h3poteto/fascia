import React from 'react'
import PropTypes from 'prop-types'
import Modal from 'react-modal'
import { Field, reduxForm } from 'redux-form'

import { RenderField, validate } from '../projectForm'

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

class EditProjectModal extends React.Component {
  componentWillReceiveProps(nextProps) {
    if (!nextProps.isProjectEditModalOpen) {
      this.handleInitialize(nextProps)
    }
  }

  handleInitialize(props) {
    const initData = {
      'title': props.project.Title,
      'description': props.project.Description,
    }

    this.props.initialize(initData)
  }

  webhookButton(project) {
    if (project == null || project.RepositoryID == undefined || project.RepositoryID == null || project.RepositoryID == 0) {
      return null
    } else {
      return (
        <button onClick={() => this.props.createWebhook(this.props.projectID)} className="pure-button button-secondary" type="button">Update Webhook</button>
      )
    }
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
      project,
    } = this.props
    return (
      <Modal
          isOpen={this.props.isProjectEditModalOpen}
          onRequestClose={onRequestClose}
          style={customStyles}
          contentLabel="EditProjectModal"
      >
        <div className="project-form">
          <form className="pure-form pure-form-stacked" onSubmit={handleSubmit((values) => { action(projectID, values) })}>
            <fieldset>
              <legend>Edit Project</legend>
              <label htmlFor="title">Title</label>
              <Field component={RenderField} name="title" id="title" type="text" />
              <label htmlFor="description">Description</label>
              <Field name="description" id="description" component="textarea" placeholder="Description" className="form-control" />
              <div className="form-action">
                <div>
                  {this.webhookButton(project)}
                </div>
                <div>
                  <button type="reset" className="pure-button pure-button-default" disabled={pristine || submitting} onClick={reset}>Reset</button>
                  <button type="submit" className="pure-button pure-button-primary" disabled={pristine || submitting}>Update Project</button>
                </div>
              </div>
            </fieldset>
          </form>
        </div>
      </Modal>
    )
  }
}

EditProjectModal.propTypes = {
  initialize: PropTypes.func.isRequired,
  handleSubmit: PropTypes.func.isRequired,
  pristine: PropTypes.bool,
  reset: PropTypes.func.isRequired,
  submitting: PropTypes.bool.isRequired,
  onRequestClose: PropTypes.func.isRequired,
  action: PropTypes.func.isRequired,
  projectID: PropTypes.string.isRequired,
  project: PropTypes.object,
  isProjectEditModalOpen: PropTypes.bool.isRequired,
  createWebhook: PropTypes.func.isRequired,
}

export default reduxForm({
  form: 'edit-project-form',
  validate,
})(EditProjectModal)
