import React from 'react';
import Modal from 'react-modal';

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
};

class ProjectView extends React.Component {
  constructor(props) {
    super(props);
  }

  componentWillMount() {
    this.props.fetchProjects();
  }

  componentDidMount() {
    this.props.fetchRepositories();
  }

  render() {
    const { isModalOpen, newProject, projects, repositories, selectedRepository } = this.props.ProjectReducer
    return (
      <div id="projects">
        <Modal
          isOpen={isModalOpen}
          onRequestClose={this.props.closeNewProjectModal}
          style={customStyles}
        >
          <div className="project-form">
            <form className="pure-form pure-form-stacked">
              <fieldset>
                <legend>Create Project</legend>
                <label htmlFor="title">Title</label>
                <input id="title" name="title" type="text" value={newProject.title} onChange={this.props.updateNewProjectTitle} placeholder="Project Name" className="form-control" />
                <label htmlFor="description">Description</label>
                <input id="description" name="description" type="text" value={newProject.description} onChange={this.props.updateNewProjectDescription} placeholder="Description" className="form-control" />
                <label htmlFor="repositories">GitHub</label>
                <select id="repositories" name="repositories" onChange={this.props.changeSelectedRepository} className="form-control">
                  <option value="0">--</option>
                  {repositories.map(function(repo, index) {
                    if (repo.id == selectedRepository) {
                      return <option value={repo.id} selected>{repo.full_name}</option>;
                    } else {
                      return<option value={repo.id}>{repo.full_name}</option>;
                    }
                   }, this)}
                </select>
                <div className="form-action">
                  <button onClick={e => this.props.fetchCreateProject(newProject.title, newProject.description, selectedRepository)} className="pure-button pure-button-primary" type="button">CreateProject</button>
                </div>
              </fieldset>
            </form>
          </div>
        </Modal>
        <div className="items">
          {projects.map(function(item, index) {
            return (
              <div className="fascia-card pure-button button-secondary" data-id={item.Id}>
              <span className="card-title">{item.Title}</span>
              <span className="card-description">{item.Description}</span>
              </div>
            );
           }, this)}
              <button onClick={this.props.openNewProjectModal} className="pure-button button-large fascia-new-project" type="button">New</button>
        </div>
      </div>
    );
  }
}

export default ProjectView;
