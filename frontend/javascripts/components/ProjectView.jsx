import React from 'react';
import { Link } from 'react-router';
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

  componentWillReceiveProps(nextProps) {
    if (this.props.ProjectReducer.error != null || nextProps.ProjectReducer.error != null) {
      setTimeout(() => {
        this.props.closeFlash()
      }, 3000);
    }
  }

  render() {
    const { isModalOpen, newProject, projects, selectedRepository, error } = this.props.ProjectReducer
    var { repositories } = this.props.ProjectReducer

    if (repositories == null ) {
      repositories = []
    }

    var flash;
    if (error != null) {
      flash = <div className="flash flash-error">{error}</div>;
    }
    return (
      <div id="projects">
        {flash}
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
                <textarea id="description" name="description" value={newProject.description} onChange={this.props.updateNewProjectDescription} placeholder="Description" className="form-control" />
                <label htmlFor="repositories">GitHub</label>
                <select id="repositories" name="repositories" onChange={this.props.changeSelectedRepository} className="form-control">
                  <option value="0">--</option>
                  {repositories.map(function(repo, index) {
                    if (selectedRepository != null && repo.id == selectedRepository.id) {
                      return <option key={index} value={repo.id} selected>{repo.full_name}</option>;
                    } else {
                      return <option key={index} value={repo.id}>{repo.full_name}</option>;
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
              <Link key={index} to={`/projects/${item.Id}`}>
                <div className="fascia-project pure-button button-secondary" data-id={item.Id}>
                  <span className="project-title">{item.Title}</span>
                  <span className="project-description">{item.Description}</span>
                </div>
              </Link>
            );
           }, this)}
              <button onClick={this.props.openNewProjectModal} className="pure-button button-large fascia-new-project button-primary" type="button">New</button>
        </div>
      </div>
    );
  }
}

export default ProjectView;
