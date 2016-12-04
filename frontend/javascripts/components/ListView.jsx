import React from 'react'
import NewListModal from './ListView/NewListModal.jsx'
import NewTaskModal from './ListView/NewTaskModal.jsx'
import EditListModal from './ListView/EditListModal.jsx'
import EditProjectModal from './ListView/EditProjectModal.jsx'
import WholeLoading from './ListView/WholeLoading.jsx'
import ListLoading from './ListView/ListLoading.jsx'
import ShowTaskModal from './ListView/ShowTaskModal.jsx'

export default class ListView extends React.Component {
  constructor(props) {
    super(props)
  }

  componentWillMount() {
    this.props.listActions.fetchLists(this.props.params.projectID)
    this.props.listActions.fetchProject(this.props.params.projectID)
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
    let stylesheet = document.styleSheets.item(4)
    var idx = document.styleSheets[4].cssRules.length
    stylesheet.insertRule("#lists .fascia-task { max-height: " + maxHeight + "px; }", idx)

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
          <span className={project.ShowPullRequests ? "pull-request-select select" : "pull-request-select"} onClick={e => this.props.listActions.showPullRequests(this.props.params.projectID, project.ShowIssues, project.ShowPullRequests)}><i title="Switch visibility of pull requests" className="octicon octicon-git-pull-request"></i></span>
          <span className={project.ShowIssues ? "pull-request-select select" : "pull-request-select"} onClick={e => this.props.listActions.showIssues(this.props.params.projectID, project.ShowIssues, project.ShowPullRequests)}><i title="Switch visibility of issues" className="octicon octicon-issue-opened"></i></span>
          <i title="Reload all lists and tasks from github" className="fa fa-repeat" onClick={e => this.props.listActions.fetchProjectGithub(this.props.params.projectID)}></i>
        </span>
      )
    }
  }

  listEditButton(list) {
    if (list.IsInitList) {
      return null
    } else {
      return <i title="Edit list" className="fa fa-pencil" onClick={e => this.props.listActions.openEditListModal(list)} data-dropped-depth="2"></i>
    }
  }

  listClass(list) {
    if (list.isDraggingOver === true) {
      return "fascia-list fascia-list-dragging-over"
    } else {
      return "fascia-list"
    }
  }

  listItem(index, list, project, taskDraggingFrom, taskDraggingTo) {
    if (list.IsHidden) {
      return (
        <div key={index} className={this.listClass(list)} data-dropped-depth="0" data-id={list.ID} onDragOver={this.props.listActions.taskDragOver} onDrop={e=> this.props.listActions.taskDrop(project.ID, taskDraggingFrom, taskDraggingTo)} onDragLeave={this.props.listActions.taskDragLeave}>
          <div className="fascia-list-menu" data-dropped-depth="1">
            <i title="Show tasks" className="fa fa-eye" onClick={e => this.props.listActions.displayList(project.ID, list.ID)} data-dropped-depth="2"></i>
            {this.listEditButton(list)}
          </div>
          <span className="list-title" data-dropped-depth="1">{list.Title}</span>
          <ListLoading isLoading={list.isLoading} />
        </div>
      )
    } else {
      return (
        <div key={index} className={this.listClass(list)} data-dropped-depth="0" data-id={list.ID} onDragOver={this.props.listActions.taskDragOver} onDrop={e=> this.props.listActions.taskDrop(project.ID, taskDraggingFrom, taskDraggingTo)} onDragLeave={this.props.listActions.taskDragLeave}>
          <div className="fascia-list-menu" data-dropped-depth="1">
            <i title="Hide tasks" className="fa fa-eye-slash" onClick={e => this.props.listActions.hideList(project.ID, list.ID)} data-dropped-depth="2"></i>
            {this.listEditButton(list)}
          </div>
          <span className="list-title" data-dropped-depth="1">{list.Title}</span>
          <ul className="fascia-task" data-dropped-depth="1">
            {list.ListTasks.map(function(task, index) {
               if (task.draggedOn) {
                 return <li key={index} className="arrow"></li>
               } else if(project != null && project.ShowIssues && !task.PullRequest || project != null && project.ShowPullRequests && task.PullRequest) {
                 return <li key={index} style={{"borderLeft": `solid 6px #${list.Color}`}} className="task" draggable="true" data-dropped-depth="2" data-id={task.ID} onDragStart={this.props.listActions.taskDragStart} onClick={e => this.props.listActions.openShowTaskModal(task)} >{task.Title}</li>
               }
             }, this)}
            <li className="new-task" data-dropped-depth="2" style={{"borderLeft": `solid 6px #${list.Color}`}} onClick={e => this.props.listActions.openNewTaskModal(list)}>
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
      lists,
      listOptions,
      noneList,
      project,
      isTaskModalOpen,
      selectedList,
      isListEditModalOpen,
      isProjectEditModalOpen,
      taskDraggingFrom,
      taskDraggingTo,
      isTaskShowModalOpen,
      isEditTaskModalVisible,
      selectedTask,
      editTask,
      error
    } = this.props.ListReducer

    return (
      <div id="lists">
        <WholeLoading isLoading={isLoading} />
        {this.flash(error)}
        <NewListModal
            isListModalOpen={isListModalOpen}
            projectID={this.props.params.projectID}
            onRequestClose={this.props.newListModalActions.closeNewListModal}
            action={this.props.newListModalActions.fetchCreateList}
        />
        <NewTaskModal
            isTaskModalOpen={isTaskModalOpen}
            listID={selectedList.ID}
            projectID={this.props.params.projectID}
            onRequestClose={this.props.newTaskModalActions.closeNewTaskModal}
            action={this.props.newTaskModalActions.fetchCreateTask}
        />
        <EditListModal
            isListEditModalOpen={isListEditModalOpen}
            list={selectedList}
            project={project}
            listOptions={listOptions}
            onRequestClose={this.props.editListModalActions.closeEditListModal}
            action={this.props.editListModalActions.fetchUpdateList}
        />
        <EditProjectModal
            isProjectEditModalOpen={isProjectEditModalOpen}
            projectID={this.props.params.projectID}
            project={project}
            onRequestClose={this.props.editProjectModalActions.closeEditProjectModal}
            action={this.props.editProjectModalActions.fetchUpdateProject}
            createWebhook={this.props.editProjectModalActions.createWebhook}
        />
        <ShowTaskModal
            isShowTaskModalOpen={isTaskShowModalOpen}
            isEditTaskModalVisible={isEditTaskModalVisible}
            projectID={this.props.params.projectID}
            task={selectedTask}
            editTask={editTask}
            closeShowTaskModal={this.props.showTaskModalActions.closeShowTaskModal}
            changeEditMode={this.props.showTaskModalActions.changeEditMode}
            updateEditTaskTitle={this.props.showTaskModalActions.updateEditTaskTitle}
            updateEditTaskDescription={this.props.showTaskModalActions.updateEditTaskDescription}
            fetchUpdateTask={this.props.showTaskModalActions.fetchUpdateTask}
            fetchDeleteTask={this.props.showTaskModalActions.fetchDeleteTask}
        />
        <div className="title-wrapper">
          <div className="project-operation">
            {this.projectOperations(project)}
          </div>
          <h3 className="project-title">{project != null ? project.Title : ''}<span className="fascia-project-menu" onClick={e => this.props.listActions.openEditProjectModal(project)}><i title="Edit project" className="fa fa-pencil"></i></span></h3>
        </div>
        <div className="items">
          {lists.map(function(list, index) {
            return this.listItem(index, list, project, taskDraggingFrom, taskDraggingTo)
           }, this)}
           <button onClick={this.props.listActions.openNewListModal} className="pure-button button-large fascia-new-list pure-button-primary" type="button">New</button>
           <div className="clearfix"></div>
        </div>
        <div className="none-list-tasks" data-dropped-depth="0" data-id={noneList.ID} onDragOver={this.props.listActions.taskDragOver} onDrop={e => this.props.listActions.taskDrop(project.ID, taskDraggingFrom, taskDraggingTo)} onDragLeave={this.props.listActions.taskDragLeave}>
          <ul className="fascia-none-list-tasks" data-dropped-depth="1">
            {noneList.ListTasks.map(function(task, index) {
               if (task.draggedOn) {
                 return <li key={index} className="arrow"></li>
               } else if( project != null && project.ShowIssues && !task.PullRequest || project != null && project.ShowPullRequests && task.PullRequest) {
                 return <li key={index} className="button-green task" draggable="true" data-dropped-depth="2" data-id={task.ID} onDragStart={this.props.listActions.taskDragStart} onClick={e => this.props.listActions.openShowTaskModal(task)}>{task.Title}</li>
               }
             }, this)}
            <li onClick={e => this.props.listActions.openNewTaskModal(noneList)} className="new-task pure-button button-blue" data-dropped-depth="2">
              <i className="fa fa-plus" data-dropped-depth="3"></i>
            </li>
          </ul>
        </div>
      </div>
    )
  }
}
