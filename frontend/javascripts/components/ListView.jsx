import React from 'react'
import NewListModal from './ListView/NewListModal.jsx'
import NewTaskModal from './ListView/NewTaskModal.jsx'
import EditListModal from './ListView/EditListModal.jsx'
import EditProjectModal from './ListView/EditProjectModal.jsx'
import WholeLoading from './ListView/WholeLoading.jsx'
import ListLoading from './ListView/ListLoading.jsx'

export default class ListView extends React.Component {
  constructor(props) {
    super(props)
  }

  componentWillMount() {
    this.props.fetchLists(this.props.params.projectID)
    this.props.fetchProject(this.props.params.projectID)
    this.props.fetchListOptions()
  }

  componentWillReceiveProps(nextProps) {
    if (this.props.ListReducer.error != null || nextProps.ListReducer.error != null) {
      setTimeout(() => {
        this.props.closeFlash()
      }, 3000)
    }
  }

  componentDidMount() {
    let maxHeight = window.innerHeight * 0.7
    let stylesheet = document.styleSheets.item(4)
    var idx = document.styleSheets[4].cssRules.length
    stylesheet.insertRule("#lists .fascia-list { max-height: " + maxHeight + "px; }", idx)

  }

  flash(error) {
    if (error != null) {
      return <div className="flash flash-error">{error}</div>
    }
  }

  projectOperations(project, selectedProject) {
    if (project == null || project.RepositoryID == undefined || project.RepositoryID == null || project.RepositoryID == 0) {
      return <span></span>
    } else {
      return (
        <span>
          <span className={project.ShowPullRequests ? "pull-request-select select" : "pull-request-select"} onClick={e => this.props.showPullRequests(this.props.params.projectID, selectedProject.ShowIssues, selectedProject.ShowPullRequests)}><i className="octicon octicon-git-pull-request"></i></span>
          <span className={project.ShowIssues ? "pull-request-select select" : "pull-request-select"} onClick={e => this.props.showIssues(this.props.params.projectID, selectedProject.ShowIssues, selectedProject.ShowPullRequests)}><i className="octicon octicon-issue-opened"></i></span>
          <i className="fa fa-repeat" onClick={e => this.props.fetchProjectGithub(this.props.params.projectID)}></i>
        </span>
      )
    }
  }

  listEditButton(list) {
    if (list.IsInitList) {
      return null
    } else {
      return <i className="fa fa-pencil" onClick={e => this.props.openEditListModal(list)} data-dropped-depth="2"></i>
    }
  }

  listItem(index, list, project, taskDraggingFrom, taskDraggingTo) {
    if (list.IsHidden) {
      return (
        <div key={index} className="fascia-list fascia-list-hidden" data-dropped-depth="0" data-id={list.ID} onDragOver={this.props.taskDragOver} onDrop={e=> this.props.taskDrop(project.ID, taskDraggingFrom, taskDraggingTo)} onDragLeave={this.props.taskDragLeave}>
          <div className="fascia-list-menu" data-dropped-depth="1">
            <i className="fa fa-eye" onClick={e => this.props.displayList(project.ID, list.ID)} data-dropped-depth="2"></i>
            {this.listEditButton(list)}
          </div>
          <span className="list-title" data-dropped-depth="1">{list.Title}</span>
          <ListLoading isLoading={list.isLoading} />
        </div>
      )
    } else {
      return (
        <div key={index} className="fascia-list" data-dropped-depth="0" data-id={list.ID} onDragOver={this.props.taskDragOver} onDrop={e=> this.props.taskDrop(project.ID, taskDraggingFrom, taskDraggingTo)} onDragLeave={this.props.taskDragLeave}>
          <div className="fascia-list-menu" data-dropped-depth="1">
            <i className="fa fa-eye-slash" onClick={e => this.props.hideList(project.ID, list.ID)} data-dropped-depth="2"></i>
            {this.listEditButton(list)}
          </div>
          <span className="list-title" data-dropped-depth="1">{list.Title}</span>
          <ul className="fascia-task" data-dropped-depth="1">
            {list.ListTasks.map(function(task, index) {
               if (task.draggedOn) {
                 return <li key={index} className="arrow"></li>
               } else if(project != null && project.ShowIssues && !task.PullRequest || project != null && project.ShowPullRequests && task.PullRequest) {
                 return <li key={index} style={{"borderLeft": `solid 6px #${list.Color}`}} className="task" draggable="true" data-dropped-depth="2" data-id={task.ID} onDragStart={this.props.taskDragStart}>{task.Title}</li>
               }
             }, this)}
            <li className="new-task" data-dropped-depth="2" style={{"borderLeft": `solid 6px #${list.Color}`}} onClick={e => this.props.openNewTaskModal(list)}>
              <i className="fa fa-plus" data-dropped-depth="3"></i>
            </li>
          </ul>
          <ListLoading isLoading={list.isLoading} />
        </div>
      )
    }
  }

  render() {
    const { isLoading, isListModalOpen, newList, lists, listOptions, noneList, project, isTaskModalOpen, newTask, selectedList, selectedListOption, isListEditModalOpen, isProjectEditModalOpen, taskDraggingFrom, taskDraggingTo, selectedProject, error } = this.props.ListReducer

    return (
      <div id="lists">
        <WholeLoading isLoading={isLoading} />
        {this.flash(error)}
        <NewListModal
            isListModalOpen={isListModalOpen}
            newList={newList}
            projectID={this.props.params.projectID}
            closeNewListModal={this.props.closeNewListModal}
            updateNewListTitle={this.props.updateNewListTitle}
            updateNewListColor={this.props.updateNewListColor}
            fetchCreateList={this.props.fetchCreateList}
        />
        <NewTaskModal
            isTaskModalOpen={isTaskModalOpen}
            newTask={newTask}
            selectedList={selectedList}
            projectID={this.props.params.projectID}
            closeNewTaskModal={this.props.closeNewTaskModal}
            updateNewTaskTitle={this.props.updateNewTaskTitle}
            updateNewTaskDescription={this.props.updateNewTaskDescription}
            fetchCreateList={this.props.fetchCreateList}
        />
        <EditListModal
            isListEditModalOpen={isListEditModalOpen}
            selectedList={selectedList}
            selectedListOption={selectedListOption}
            project={project}
            listOptions={listOptions}
            projectID={this.props.params.projectID}
            closeEditListModal={this.props.closeEditListModal}
            updateSelectedListTitle={this.props.updateSelectedListTitle}
            updateSelectedListColor={this.props.updateSelectedListColor}
            changeSelectedListOption={this.props.changeSelectedListOption}
            fetchUpdateList={this.props.fetchUpdateList}
        />
        <EditProjectModal
            isProjectEditModalOpen={isProjectEditModalOpen}
            selectedProject={selectedProject}
            project={project}
            projectID={this.props.params.projectID}
            closeEditProjectModal={this.props.closeEditProjectModal}
            updateEditProjectTitle={this.props.updateEditProjectTitle}
            updateEditProjectDescription={this.props.updateEditProjectDescription}
            fetchUpdateProject={this.props.fetchUpdateProject}
            createWebhook={this.props.createWebhook}
        />
        <div className="title-wrapper">
          <div className="project-operation">
            {this.projectOperations(project, selectedProject)}
          </div>
          <h3 className="project-title">{project != null ? project.Title : ''}<span className="fascia-project-menu" onClick={e => this.props.openEditProjectModal(project)}><i className="fa fa-pencil"></i></span></h3>
        </div>
        <div className="items">
          {lists.map(function(list, index) {
            return this.listItem(index, list, project, taskDraggingFrom, taskDraggingTo)
           }, this)}
           <button onClick={this.props.openNewListModal} className="pure-button button-large fascia-new-list pure-button-primary" type="button">New</button>
           <div className="clearfix"></div>
        </div>
        <div className="none-list-tasks" data-dropped-depth="0" data-id={noneList.ID} onDragOver={this.props.taskDragOver} onDrop={e => this.props.taskDrop(project.ID, taskDraggingFrom, taskDraggingTo)} onDragLeave={this.props.taskDragLeave}>
          <ul className="fascia-none-list-tasks" data-dropped-depth="1">
            {noneList.ListTasks.map(function(task, index) {
               if (task.draggedOn) {
                 return <li key={index} className="arrow"></li>
               } else if( project != null && project.ShowIssues && !task.PullRequest || project != null && project.ShowPullRequests && task.PullRequest) {
                 return <li key={index} className="button-green task" draggable="true" data-dropped-depth="2" data-id={task.ID} onDragStart={this.props.taskDragStart}>{task.Title}</li>
               }
             }, this)}
            <li onClick={e => this.props.openNewTaskModal(noneList)} className="new-task pure-button button-blue" data-dropped-depth="2">
              <i className="fa fa-plus" data-dropped-depth="3"></i>
            </li>
          </ul>
        </div>
      </div>
    )
  }
}
