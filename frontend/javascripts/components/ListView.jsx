import React from 'react'
import PropTypes from 'prop-types'
import NewListModal from './ListView/NewListModal.jsx'
import NewTaskModal from './ListView/NewTaskModal.jsx'
import EditListModal from './ListView/EditListModal.jsx'
import EditProjectModal from './ListView/EditProjectModal.jsx'
import WholeLoading from './ListView/WholeLoading.jsx'
import ListLoading from './ListView/ListLoading.jsx'
import ShowTaskModal from './ListView/ShowTaskModal.jsx'
import DeleteProjectModal from './ListView/DeleteProjectModal.jsx'

class ListView extends React.Component {
  constructor(props) {
    super(props)
  }

  componentWillMount() {
    this.props.listActions.fetchLists(this.props.match.params.projectID)
    this.props.listActions.fetchProject(this.props.match.params.projectID)
    this.props.listActions.fetchListOptions()
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.ListReducer.error != null || nextProps.ListReducer.error != null) {
      setTimeout(() => {
        this.props.listActions.closeFlash()
      }, 3000)
    }
  }

  componentDidMount() {
    let maxHeight = window.innerHeight * 0.7
    let stylesheet = document.styleSheets.item(2)
    var idx = stylesheet.cssRules.length
    stylesheet.insertRule('#lists .fascia-task { max-height: ' + maxHeight + 'px; }', idx)

  }

  flash(error) {
    if (error != null) {
      return <div className="flash flash-error">{error}</div>
    }
  }

  projectOperations(project) {
    if (project == null || project.RepositoryID == undefined || project.RepositoryID == null || project.RepositoryID == 0) {
      return <span></span>
    } else {
      return (
        <span>
          <span className={project.ShowPullRequests ? 'pull-request-select select' : 'pull-request-select'} onClick={() => this.props.listActions.showPullRequests(this.props.match.params.projectID, project.ShowIssues, project.ShowPullRequests)}><i title="Switch visibility of pull requests" className="octicon octicon-git-pull-request"></i></span>
          <span className={project.ShowIssues ? 'pull-request-select select' : 'pull-request-select'} onClick={() => this.props.listActions.showIssues(this.props.match.params.projectID, project.ShowIssues, project.ShowPullRequests)}><i title="Switch visibility of issues" className="octicon octicon-issue-opened"></i></span>
          <i title="Reload all lists and tasks from github" className="fa fa-repeat" onClick={() => this.props.listActions.fetchProjectGithub(this.props.match.params.projectID)}></i>
        </span>
      )
    }
  }

  listEditButton(list) {
    if (list.IsInitList) {
      return null
    } else {
      return <i title="Edit list" className="fa fa-pencil" onClick={() => this.props.listActions.openEditListModal(list)} data-dropped-depth="2"></i>
    }
  }

  listClass(list) {
    if (list.isDraggingOver === true) {
      return 'fascia-list fascia-list-dragging-over'
    } else {
      return 'fascia-list'
    }
  }

  listItem(index, list, project, taskDraggingFrom, taskDraggingTo) {
    if (list.IsHidden) {
      return (
        <div key={index} className={this.listClass(list)} data-dropped-depth="0" data-id={list.ID} onDragOver={this.props.listActions.taskDragOver} onDrop={() => this.props.listActions.taskDrop(project.ID, taskDraggingFrom, taskDraggingTo)} onDragLeave={this.props.listActions.taskDragLeave}>
          <div className="fascia-list-menu" data-dropped-depth="1">
            <i title="Show tasks" className="fa fa-eye" onClick={() => this.props.listActions.displayList(project.ID, list.ID)} data-dropped-depth="2"></i>
            {this.listEditButton(list)}
          </div>
          <span className="list-title" data-dropped-depth="1">{list.Title}</span>
          <ListLoading isLoading={list.isLoading} />
        </div>
      )
    } else {
      return (
        <div key={index} className={this.listClass(list)} data-dropped-depth="0" data-id={list.ID} onDragOver={this.props.listActions.taskDragOver} onDrop={() => this.props.listActions.taskDrop(project.ID, taskDraggingFrom, taskDraggingTo)} onDragLeave={this.props.listActions.taskDragLeave}>
          <div className="fascia-list-menu" data-dropped-depth="1">
            <i title="Hide tasks" className="fa fa-eye-slash" onClick={() => this.props.listActions.hideList(project.ID, list.ID)} data-dropped-depth="2"></i>
            {this.listEditButton(list)}
          </div>
          <span className="list-title" data-dropped-depth="1">{list.Title}</span>
          <ul className="fascia-task" data-dropped-depth="1">
            {list.ListTasks.map(function(task, index) {
               if (task.draggedOn) {
                 return <li key={index} className="arrow"></li>
               } else if(project != null && project.ShowIssues && !task.PullRequest || project != null && project.ShowPullRequests && task.PullRequest) {
                 return <li key={index} style={{'borderLeft': `solid 6px #${list.Color}`}} className="task" draggable="true" data-dropped-depth="2" data-id={task.ID} onDragStart={this.props.listActions.taskDragStart} onClick={() => this.props.listActions.openShowTaskModal(task)} >{task.Title}</li>
               }
             }, this)}
            <li className="new-task" data-dropped-depth="2" style={{'borderLeft': `solid 6px #${list.Color}`}} onClick={() => this.props.listActions.openNewTaskModal(list)}>
              <i className="fa fa-plus" data-dropped-depth="3"></i>
            </li>
          </ul>
          <ListLoading isLoading={list.isLoading} />
        </div>
      )
    }
  }

  render() {
    const {
      isLoading,
      isListModalOpen,
      isTaskModalOpen,
      isListEditModalOpen,
      isProjectEditModalOpen,
      isTaskShowModalOpen,
      isEditTaskModalVisible,
      isDeleteProjectModalOpen,
      lists,
      listOptions,
      noneList,
      project,
      selectedList,
      taskDraggingFrom,
      taskDraggingTo,
      color,
      selectedTask,
      error
    } = this.props.ListReducer

    return (
      <div id="lists">
        <WholeLoading isLoading={isLoading} />
        {this.flash(error)}
        <EditProjectModal
          isProjectEditModalOpen={isProjectEditModalOpen}
          project={project}
          onRequestClose={this.props.editProjectModalActions.closeEditProjectModal}
          action={this.props.editProjectModalActions.fetchUpdateProject}
          createWebhook={this.props.editProjectModalActions.createWebhook}
          flashMessage={error}
        />
        <DeleteProjectModal
          isDeleteProjectModalOpen={isDeleteProjectModalOpen}
          onRequestClose={this.props.deleteProjectModalActions.closeDeleteProjectModal}
          project={project}
          action={this.props.deleteProjectModalActions.fetchDeleteProject}
          flashMessage={error}
        />
        <NewListModal
          isListModalOpen={isListModalOpen}
          onRequestClose={this.props.newListModalActions.closeNewListModal}
          action={this.props.newListModalActions.fetchCreateList}
          changeColor={this.props.newListModalActions.changeColor}
          color={color}
          flashMessage={error}
        />
        <EditListModal
          isListEditModalOpen={isListEditModalOpen}
          list={selectedList}
          project={project}
          listOptions={listOptions}
          onRequestClose={this.props.editListModalActions.closeEditListModal}
          action={this.props.editListModalActions.fetchUpdateList}
          changeColor={this.props.editListModalActions.changeColor}
          color={color}
          flashMessage={error}
        />
        <NewTaskModal
          isTaskModalOpen={isTaskModalOpen}
          onRequestClose={this.props.newTaskModalActions.closeNewTaskModal}
          action={this.props.newTaskModalActions.fetchCreateTask}
          flashMessage={error}
        />
        <ShowTaskModal
          isShowTaskModalOpen={isTaskShowModalOpen}
          isEditTaskModalVisible={isEditTaskModalVisible}
          project={project}
          task={selectedTask}
          onRequestClose={this.props.showTaskModalActions.closeShowTaskModal}
          changeEditMode={this.props.showTaskModalActions.changeEditMode}
          action={this.props.showTaskModalActions.fetchUpdateTask}
          fetchDeleteTask={this.props.showTaskModalActions.fetchDeleteTask}
          flashMessage={error}
        />
        <div className="title-wrapper">
          <div className="project-operation">
            {this.projectOperations(project)}
          </div>
          <h3 className="project-title">{project != null ? project.Title : ''}<span className="fascia-project-menu" onClick={()=> this.props.listActions.openDeleteProjectModal()}><i title="Delete project" className="fa fa-trash"></i></span><span className="fascia-project-menu" onClick={() => this.props.listActions.openEditProjectModal(project)}><i title="Edit project" className="fa fa-pencil"></i></span>
          </h3>
        </div>
        <div className="items">
          {lists.map(function(list, index) {
             return this.listItem(index, list, project, taskDraggingFrom, taskDraggingTo)
           }, this)}
          <button onClick={this.props.listActions.openNewListModal} className="pure-button button-large fascia-new-list pure-button-primary" type="button">New</button>
          <div className="clearfix"></div>
        </div>
        <div className="none-list-tasks" data-dropped-depth="0" data-id={noneList.ID} onDragOver={this.props.listActions.taskDragOver} onDrop={() => this.props.listActions.taskDrop(project.ID, taskDraggingFrom, taskDraggingTo)} onDragLeave={this.props.listActions.taskDragLeave}>
          <ul className="fascia-none-list-tasks" data-dropped-depth="1">
            {noneList.ListTasks.map(function(task, index) {
               if (task.draggedOn) {
                 return <li key={index} className="arrow"></li>
               } else if( project != null && project.ShowIssues && !task.PullRequest || project != null && project.ShowPullRequests && task.PullRequest) {
                 return <li key={index} className="button-green task" draggable="true" data-dropped-depth="2" data-id={task.ID} onDragStart={this.props.listActions.taskDragStart} onClick={() => this.props.listActions.openShowTaskModal(task)}>{task.Title}</li>
               }
             }, this)}
            <li onClick={() => this.props.listActions.openNewTaskModal(noneList)} className="new-task pure-button button-blue" data-dropped-depth="2">
              <i className="fa fa-plus" data-dropped-depth="3"></i>
            </li>
          </ul>
        </div>
      </div>
    )
  }
}

ListView.propTypes = {
  listActions: PropTypes.shape({
    fetchLists: PropTypes.func.isRequired,
    fetchProject: PropTypes.func.isRequired,
    fetchListOptions: PropTypes.func.isRequired,
    closeFlash: PropTypes.func.isRequired,
    showPullRequests: PropTypes.func.isRequired,
    showIssues: PropTypes.func.isRequired,
    fetchProjectGithub: PropTypes.func.isRequired,
    openEditListModal: PropTypes.func.isRequired,
    taskDragOver: PropTypes.func.isRequired,
    taskDrop: PropTypes.func.isRequired,
    taskDragLeave: PropTypes.func.isRequired,
    displayList: PropTypes.func.isRequired,
    hideList: PropTypes.func.isRequired,
    taskDragStart: PropTypes.func.isRequired,
    openShowTaskModal: PropTypes.func.isRequired,
    openNewTaskModal: PropTypes.func.isRequired,
    openEditProjectModal: PropTypes.func.isRequired,
    openNewListModal: PropTypes.func.isRequired,
    openDeleteProjectModal: PropTypes.func.isRequired,
  }),
  newListModalActions: PropTypes.shape({
    closeNewListModal: PropTypes.func.isRequired,
    fetchCreateList: PropTypes.func.isRequired,
    changeColor: PropTypes.func.isRequired,
  }),
  newTaskModalActions: PropTypes.shape({
    closeNewTaskModal: PropTypes.func.isRequired,
    fetchCreateTask: PropTypes.func.isRequired,
  }),
  editListModalActions: PropTypes.shape({
    closeEditListModal: PropTypes.func.isRequired,
    fetchUpdateList: PropTypes.func.isRequired,
    changeColor: PropTypes.func.isRequired,
  }),
  editProjectModalActions: PropTypes.shape({
    closeEditProjectModal: PropTypes.func.isRequired,
    fetchUpdateProject: PropTypes.func.isRequired,
    createWebhook: PropTypes.func.isRequired,
  }),
  showTaskModalActions: PropTypes.shape({
    closeShowTaskModal: PropTypes.func.isRequired,
    changeEditMode: PropTypes.func.isRequired,
    fetchUpdateTask: PropTypes.func.isRequired,
    fetchDeleteTask: PropTypes.func.isRequired,
  }),
  deleteProjectModalActions: PropTypes.shape({
    closeDeleteProjectModal: PropTypes.func.isRequired,
    fetchDeleteProject: PropTypes.func.isRequired,
  }),
  params: PropTypes.shape({
    projectID: PropTypes.string.isRequired,
  }),
  ListReducer: PropTypes.shape({
    isLoading: PropTypes.bool.isRequired,
    isListModalOpen: PropTypes.bool.isRequired,
    isTaskModalOpen: PropTypes.bool.isRequired,
    isListEditModalOpen: PropTypes.bool.isRequired,
    isProjectEditModalOpen: PropTypes.bool.isRequired,
    isTaskShowModalOpen: PropTypes.bool.isRequired,
    isEditTaskModalVisible: PropTypes.bool.isRequired,
    isDeleteProjectModalOpen: PropTypes.bool.isRequired,
    lists: PropTypes.array,
    listOptions: PropTypes.array,
    noneList: PropTypes.object,
    project: PropTypes.object,
    selectedList: PropTypes.object,
    taskDraggingFrom: PropTypes.object,
    taskDraggingTo: PropTypes.object,
    selectedTask: PropTypes.object,
    error: PropTypes.string,
    color: PropTypes.string,
  }),

}

export default ListView
