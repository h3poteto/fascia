import React from 'react'
import PropTypes from 'prop-types'
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
  projectActions: PropTypes.shape({
    fetchProjects: PropTypes.func.isRequired,
    fetchRepositories: PropTypes.func.isRequired,
    fetchSession: PropTypes.func.isRequired,
    closeFlash: PropTypes.func.isRequired,
    openNewProjectModal: PropTypes.func.isRequired,
  }),
  newProjectModalActions: PropTypes.shape({
    closeNewProjectModal: PropTypes.func.isRequired,
    fetchCreateProject: PropTypes.func.isRequired,
  }),
  ProjectReducer: PropTypes.shape({
    isModalOpen: PropTypes.bool.isRequired,
    projects: PropTypes.array,
    isLoading: PropTypes.bool.isRequired,
    error: PropTypes.string,
    repositories: PropTypes.array,
  }),
}

export default ProjectView
