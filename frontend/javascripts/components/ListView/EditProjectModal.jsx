import React from 'react'
import Modal from 'react-modal'

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

export default class EditProjectModal extends React.Component {
  constructor(props) {
    super(props)
  }

  webhookButton(project) {
    if (project == null || project.RepositoryID == undefined || project.RepositoryID == null || project.RepositoryID == 0) {
      return null
    } else {
      return (
        <button onClick={e => this.props.createWebhook(this.props.projectID)} className="pure-button button-secondary" type="button">Update Webhook</button>
      )
    }
  }

  render() {
    return (
      <Modal
          isOpen={this.props.isProjectEditModalOpen}
          onRequestClose={this.props.closeEditProjectModal}
          style={customStyles}
      >
        <div className="project-form">
          <form className="pure-form pure-form-stacked">
            <fieldset>
              <legend>Edit Project</legend>
              <label htmlFor="title">Title</label>
              <input id="title" name="title" type="text" value={this.props.selectedProject.Title} onChange={this.props.updateEditProjectTitle} className="form-control" />
              <label htmlFor="description">Description</label>
              <textarea id="description" name="description" value={this.props.selectedProject.Description} onChange={this.props.updateEditProjectDescription} className="form-control" />
              <div className="form-action">
                {this.webhookButton(this.props.project)}&nbsp;
                <button onClick={e => this.props.fetchUpdateProject(this.props.projectID, this.props.selectedProject)} className="pure-button pure-button-primary" type="button">Update Project</button>
              </div>
            </fieldset>
          </form>
        </div>
      </Modal>
    )
  }
}
