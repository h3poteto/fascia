import React from 'react'
import { Link } from 'react-router'
import Modal from 'react-modal'
import truncate from 'html-truncate'
import NewProjectModal from './ProjectView/NewProjectModal.jsx'
import WholeLoading from './ProjectView/WholeLoading.jsx'

class ProjectView extends React.Component {
  constructor(props) {
    super(props)
  }

  componentWillMount() {
    this.props.projectActions.fetchProjects()
  }

  componentDidMount() {
    this.props.projectActions.fetchRepositories()
    this.props.projectActions.fetchSession()
  }

  componentWillReceiveProps(nextProps) {
    // TODO: これnextだけ見ればいいのでは？
    if (this.props.ProjectReducer.error != null || nextProps.ProjectReducer.error != null) {
      setTimeout(() => {
        this.props.projectActions.closeFlash()
      }, 3000)
    }
  }

  render() {
    const { isModalOpen, newProject, projects, selectedRepository, isLoading, error } = this.props.ProjectReducer
    var { repositories } = this.props.ProjectReducer

    if (repositories == null ) {
      repositories = []
    }

    var flash
    if (error != null) {
      flash = <div className="flash flash-error">{error}</div>
    }
    return (
      <div id="projects">
        <WholeLoading isLoading={isLoading} />
        {flash}
        <NewProjectModal
            isModalOpen={isModalOpen}
            newProject={newProject}
            repositories={repositories}
            selectedRepository={selectedRepository}
            closeNewProjectModal={this.props.newProjectModalActions.closeNewProjectModal}
            updateNewProjectTitle={this.props.newProjectModalActions.updateNewProjectTitle}
            updateNewProjectDescription={this.props.newProjectModalActions.updateNewProjectDescription}
            changeSelectedRepository={this.props.newProjectModalActions.changeSelectedRepository}
            fetchCreateProject={this.props.newProjectModalActions.fetchCreateProject}
        />
        <div className="items">
          {projects.map(function(item, index) {
            return (
              <Link key={index} to={`/projects/${item.ID}`}>
                <div className="fascia-project pure-button button-secondary" data-id={item.ID}>
                  <div className="project-title">{item.Title}</div>
                  <div className="project-description">{truncate(item.Description, 52)}</div>
                </div>
              </Link>
            )
           }, this)}
              <button onClick={this.props.projectActions.openNewProjectModal} className="pure-button button-large fascia-new-project button-primary" type="button">New</button>
        </div>
      </div>
    )
  }
}

export default ProjectView
