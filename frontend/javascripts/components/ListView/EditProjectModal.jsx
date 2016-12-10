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
      >
        <div className="project-form">
          <form className="pure-form pure-form-stacked" onSubmit={handleSubmit((values) => { action(projectID, values) })}>
            <fieldset>
              <legend>Edit Project</legend>
              <label htmlFor="title">Title</label>
              <Field name="title" id="title" component="input" type="text" className="form-control" />
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
  initialize: React.PropTypes.func.isRequired,
  handleSubmit: React.PropTypes.func.isRequired,
  pristine: React.PropTypes.bool,
  reset: React.PropTypes.func.isRequired,
  submitting: React.PropTypes.bool.isRequired,
  onRequestClose: React.PropTypes.func.isRequired,
  action: React.PropTypes.func.isRequired,
  projectID: React.PropTypes.string.isRequired,
  project: React.PropTypes.object,
  isProjectEditModalOpen: React.PropTypes.bool.isRequired,
  createWebhook: React.PropTypes.func.isRequired,
}

export default reduxForm({
  form: 'edit-project-form',
})(EditProjectModal)
