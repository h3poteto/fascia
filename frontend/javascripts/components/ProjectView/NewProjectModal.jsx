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

export default class NewProjectModal extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    return (
      <Modal
          isOpen={this.props.isModalOpen}
          onRequestClose={this.props.closeNewProjectModal}
          style={customStyles}
      >
        <div className="project-form">
          <form className="pure-form pure-form-stacked">
            <fieldset>
              <legend>Create Project</legend>
              <label htmlFor="title">Title</label>
              <input id="title" name="title" type="text" value={this.props.newProject.title} onChange={this.props.updateNewProjectTitle} placeholder="Project Name" className="form-control" />
              <label htmlFor="description">Description</label>
              <textarea id="description" name="description" value={this.props.newProject.description} onChange={this.props.updateNewProjectDescription} placeholder="Description" className="form-control" />
              <label htmlFor="repositories">GitHub</label>
              <select id="repositories" name="repositories" onChange={this.props.changeSelectedRepository} className="form-control">
                <option value="0">--</option>
                {this.props.repositories.map(function(repo, index) {
                   if (this.props.selectedRepository != null && repo.id == this.props.selectedRepository.id) {
                     return <option key={index} value={repo.id} selected>{repo.full_name}</option>
                   } else {
                     return <option key={index} value={repo.id}>{repo.full_name}</option>
                   }
                 }, this)}
              </select>
              <div className="form-action">
                <button onClick={e => this.props.fetchCreateProject(this.props.newProject.title, this.props.newProject.description, this.props.selectedRepository)} className="pure-button pure-button-primary" type="button">CreateProject</button>
              </div>
            </fieldset>
          </form>
        </div>
      </Modal>
    )
  }
}
