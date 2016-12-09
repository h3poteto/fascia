import React from 'react'
import { Link } from 'react-router'
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
    const {
      isModalOpen,
      projects,
      repositories,
      isLoading,
      error
    } = this.props.ProjectReducer

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
            onRequestClose={this.props.newProjectModalActions.closeNewProjectModal}
            action={this.props.newProjectModalActions.fetchCreateProject}
            repositories={repositories}
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

ProjectView.propTypes = {
  projectActions: React.PropTypes.shape({
    fetchProjects: React.PropTypes.func.isRequired,
    fetchRepositories: React.PropTypes.func.isRequired,
    fetchSession: React.PropTypes.func.isRequired,
    closeFlash: React.PropTypes.func.isRequired,
    openNewProjectModal: React.PropTypes.func.isRequired,
  }),
  newProjectModalActions: React.PropTypes.shape({
    closeNewProjectModal: React.PropTypes.func.isRequired,
    fetchCreateProject: React.PropTypes.func.isRequired,
  }),
  ProjectReducer: React.PropTypes.shape({
    isModalOpen: React.PropTypes.bool.isRequired,
    projects: React.PropTypes.array,
    isLoading: React.PropTypes.bool.isRequired,
    error: React.PropTypes.string,
    repositories: React.PropTypes.array,
  }),
}

export default ProjectView
